package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"goproject/domain"
)

var (
	DB     *gorm.DB
	ES     *elasticsearch.Client
	RDB    *redis.Client
	Ctx    = context.Background()
	AppCfg *AppConfig
)

type AppConfig struct {
	AppEnv        string
	Port          string
	GinMode       string
	DBHost        string
	DBUser        string
	DBPassword    string
	DBName        string
	DBPort        string
	ESAddress     string
	ESUsername    string
	ESPassword    string
	RedisAddr     string
	RedisPassword string
	KafkaBrokers  string
	KafkaUsername string
	KafkaPassword string
	KafkaTLS      bool
	KafkaTopic    string
	KafkaGroupID  string
	GRPCPort      string
}

func getEnv(key, defaultValue string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return defaultValue
	}
	return value
}

func LoadConfig() (*AppConfig, error) {
	cfg := &AppConfig{
		AppEnv:        strings.ToLower(getEnv("APP_ENV", "development")),
		Port:          getEnv("PORT", "8080"),
		GinMode:       strings.ToLower(getEnv("GIN_MODE", "debug")),
		DBHost:        getEnv("DB_HOST", "localhost"),
		DBUser:        getEnv("DB_USER", "admin"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        getEnv("DB_NAME", "godb"),
		DBPort:        getEnv("DB_PORT", "5432"),
		ESAddress:     getEnv("ES_ADDRESS", "http://localhost:9200"),
		ESUsername:    os.Getenv("ES_USERNAME"),
		ESPassword:    os.Getenv("ES_PASSWORD"),
		RedisAddr:     getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		KafkaBrokers:  getEnv("KAFKA_BROKERS", "localhost:9092"),
		KafkaUsername: os.Getenv("KAFKA_USERNAME"),
		KafkaPassword: os.Getenv("KAFKA_PASSWORD"),
		KafkaTLS:      strings.EqualFold(getEnv("KAFKA_TLS", "true"), "true"),
		KafkaTopic:    getEnv("KAFKA_ORDER_TOPIC", "warehouse.order.created"),
		KafkaGroupID:  getEnv("KAFKA_CONSUMER_GROUP", "warehouse-service"),
		GRPCPort:      getEnv("GRPC_PORT", "9090"),
	}

	if cfg.AppEnv == "production" && cfg.GinMode == "debug" {
		cfg.GinMode = "release"
	}

	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		return nil, fmt.Errorf("invalid database config: DB_HOST, DB_USER and DB_NAME are required")
	}
	if cfg.AppEnv == "production" && strings.TrimSpace(cfg.DBPassword) == "" {
		return nil, fmt.Errorf("DB_PASSWORD is required in production")
	}

	return cfg, nil
}

func SetAppConfig(cfg *AppConfig) {
	AppCfg = cfg
}

func (c *AppConfig) DSN() string {
	sslMode := "disable"
	if c.AppEnv == "production" {
		sslMode = "require"
	}
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=Asia/Ho_Chi_Minh",
		c.DBHost,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBPort,
		sslMode,
	)
}

func ConnectDB() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi cấu hình database: %v", err)
	}
	SetAppConfig(cfg)

	DB, err = gorm.Open(postgres.Open(cfg.DSN()), &gorm.Config{})
	if err != nil {
		log.Fatalf("Lỗi Postgres: %v", err)
	}

	if err := DB.AutoMigrate(
		&domain.User{}, &domain.Role{}, &domain.UserRole{}, &domain.Product{},
		&domain.Order{}, &domain.OrderItem{}, &domain.Supplier{}, &domain.Customer{},
	); err != nil {
		log.Fatalf("AutoMigrate Postgres thất bại: %v", err)
	}

	log.Println("✅ Postgres kết nối thành công")
}

func ConnectES() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi cấu hình Elasticsearch: %v", err)
	}

	ES, err = elasticsearch.NewClient(elasticsearch.Config{Addresses: []string{cfg.ESAddress}, Username: cfg.ESUsername, Password: cfg.ESPassword})
	if err != nil {
		log.Fatalf("Lỗi Elasticsearch: %v", err)
	}

	log.Println("✅ Elasticsearch kết nối thành công")
}

func ConnectRedis() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Lỗi cấu hình Redis: %v", err)
	}

	RDB = redis.NewClient(&redis.Options{Addr: cfg.RedisAddr, Password: cfg.RedisPassword})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := RDB.Ping(ctx).Err(); err != nil {
		log.Fatalf("Lỗi Redis: %v", err)
	}

	log.Println("✅ Redis kết nối thành công")
}
