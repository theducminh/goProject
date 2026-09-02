package grpcdelivery

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"strconv"

	"goproject/domain"
	"goproject/usecase"

	"github.com/go-playground/validator/v10"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CheckStockRequest struct {
	ProductID uint `json:"product_id" validate:"required,gt=0"`
	Quantity  uint `json:"quantity" validate:"required,gt=0"`
}

type CheckStockResponse struct {
	Available bool `json:"available"`
	Stock     uint `json:"stock"`
}

type WarehouseServer struct {
	stock    usecase.StockUsecase
	validate *validator.Validate
}

func NewWarehouseServer(stock usecase.StockUsecase) *WarehouseServer {
	return &WarehouseServer{stock: stock, validate: validator.New()}
}

func (server *WarehouseServer) CheckStock(ctx context.Context, request *CheckStockRequest) (*CheckStockResponse, error) {
	if request == nil || server.validate.Struct(request) != nil {
		return nil, status.Error(codes.InvalidArgument, "product_id and quantity must be positive")
	}
	available, stock, err := server.stock.Check(ctx, request.ProductID, request.Quantity)
	if err != nil {
		if errors.Is(err, domain.ErrNotFound) {
			return nil, status.Error(codes.NotFound, "product not found")
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}
	return &CheckStockResponse{Available: available, Stock: stock}, nil
}

func RegisterWarehouseServer(grpcServer *grpc.Server, server *WarehouseServer) {
	grpcServer.RegisterService(&grpc.ServiceDesc{
		ServiceName: "warehouse.v1.WarehouseService",
		HandlerType: (*warehouseService)(nil),
		Methods:     []grpc.MethodDesc{{MethodName: "CheckStock", Handler: checkStockHandler}},
	}, server)
}

type warehouseService interface {
	CheckStock(context.Context, *CheckStockRequest) (*CheckStockResponse, error)
}

func checkStockHandler(server interface{}, ctx context.Context, decode func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	request := new(CheckStockRequest)
	if err := decode(request); err != nil {
		return nil, err
	}
	call := func(ctx context.Context, request interface{}) (interface{}, error) {
		return server.(warehouseService).CheckStock(ctx, request.(*CheckStockRequest))
	}
	if interceptor == nil {
		return call(ctx, request)
	}
	return interceptor(ctx, request, &grpc.UnaryServerInfo{FullMethod: "/warehouse.v1.WarehouseService/CheckStock"}, call)
}

type jsonCodec struct{}

func (jsonCodec) Name() string                                   { return "warehouse-json" }
func (jsonCodec) Marshal(value interface{}) ([]byte, error)      { return json.Marshal(value) }
func (jsonCodec) Unmarshal(data []byte, value interface{}) error { return json.Unmarshal(data, value) }

func NewServer(stock usecase.StockUsecase) *grpc.Server {
	server := grpc.NewServer(grpc.ForceServerCodec(jsonCodec{}))
	RegisterWarehouseServer(server, NewWarehouseServer(stock))
	return server
}

func ListenAndServe(ctx context.Context, server *grpc.Server, address string) error {
	listener, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	go func() {
		<-ctx.Done()
		server.GracefulStop()
	}()
	return server.Serve(listener)
}

func Address(port int) string { return ":" + strconv.Itoa(port) }
