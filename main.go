package main

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"goproject/config"
	grpcdelivery "goproject/delivery/grpc"
	_ "goproject/docs"
	"goproject/domain"
	"goproject/repository"
	elasticsearchRepository "goproject/repository/elasticsearch"
	kafkaRepository "goproject/repository/kafka"
	postgresRepository "goproject/repository/postgres"
	redisRepository "goproject/repository/redis"
	"goproject/routes"
	"goproject/usecase"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

// @title API Hệ thống tích hợp Postgres, ES, Redis
// @version 1.0
// @host localhost:8080
// @BasePath /
func main() {
	if err := godotenv.Load(); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("could not load .env: %v", err)
	}
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	config.SetAppConfig(cfg)
	gin.SetMode(cfg.GinMode)

	config.ConnectDB()
	config.ConnectES()
	config.ConnectRedis()

	productRepository := postgresRepository.NewProductRepository(config.DB)
	productCache := redisRepository.NewProductCache(config.RDB, 10*time.Minute)
	productSearch := elasticsearchRepository.NewProductSearchRepository(config.ES, "products")
	productUsecase := usecase.NewProductUsecase(productRepository, productCache, productSearch, 5*time.Second)
	orderRepository := postgresRepository.NewOrderRepository(config.DB)

	brokers := strings.Split(cfg.KafkaBrokers, ",")
	producer, err := kafkaRepository.NewProducer(brokers, cfg.KafkaUsername, cfg.KafkaPassword, cfg.KafkaTLS)
	if err != nil {
		log.Fatalf("Kafka producer configuration failed: %v", err)
	}
	consumer, err := kafkaRepository.NewConsumer(brokers, cfg.KafkaGroupID, cfg.KafkaTopic, cfg.KafkaUsername, cfg.KafkaPassword, cfg.KafkaTLS)
	if err != nil {
		_ = producer.Close()
		log.Fatalf("Kafka consumer configuration failed: %v", err)
	}
	orderUsecase := usecase.NewOrderUsecase(orderRepository, productRepository, producer, cfg.KafkaTopic, 5*time.Second)
	stockUsecase := usecase.NewStockUsecase(productRepository, 5*time.Second)
	r := routes.SetupRouter(productUsecase, orderUsecase, productRepository, productCache, productSearch)

	shutdownContext, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	go func() {
		if err := consumer.Consume(shutdownContext, handleOrderCreated); err != nil && shutdownContext.Err() == nil {
			log.Printf("Kafka consumer stopped: %v", err)
		}
	}()
	grpcServer := grpcdelivery.NewServer(stockUsecase)
	grpcErrors := make(chan error, 1)
	go func() { grpcErrors <- grpcdelivery.ListenAndServe(shutdownContext, grpcServer, ":"+cfg.GRPCPort) }()

	server := &http.Server{Addr: ":" + cfg.Port, Handler: r, ReadHeaderTimeout: 5 * time.Second}
	serverErrors := make(chan error, 1)
	go func() { serverErrors <- server.ListenAndServe() }()
	select {
	case err := <-serverErrors:
		if err != http.ErrServerClosed {
			log.Printf("Server failed to start: %v", err)
		}
	case <-shutdownContext.Done():
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		_ = server.Shutdown(shutdown)
	}
	select {
	case err := <-grpcErrors:
		if err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	case <-time.After(100 * time.Millisecond):
	}
	_ = consumer.Close()
	_ = producer.Close()
}

func handleOrderCreated(ctx context.Context, event repository.Event) error {
	var orderEvent domain.OrderCreated
	if err := json.Unmarshal(event.Value, &orderEvent); err != nil {
		return err
	}
	if orderEvent.EventType != domain.OrderCreatedEvent || orderEvent.Order.ID == 0 {
		return errors.New("invalid order created event")
	}
	log.Printf("received order created event %s for order %d", orderEvent.EventID, orderEvent.Order.ID)
	return nil
}
