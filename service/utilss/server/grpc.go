package server

import (
	"context"
	"log"
	"net"

	pb "shelon_server/proto"
	"shelon_server/transport"
	"shelon_server/utilss/logger"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// GRPCServer представляет реализацию gRPC сервера.
type GRPCServer struct {
	logger           logger.Logger
	server           *grpc.Server
	port             string
	transportService *transport.TransportService
}

// NewGRPCServer создает новый экземпляр GRPC-сервера с предоставленными зависимостями.
// port: порт, на котором будет слушать сервер.
// logger: экземпляр интерфейса logger.Logger для логирования действий.
// transportService: экземпляр транспортного сервиса для обработки данных.
func NewGRPCServer(port string, logger logger.Logger, transportService *transport.TransportService) *GRPCServer {
	logger.Info("Creating new instance of GRPC server")
	return &GRPCServer{
		logger:           logger,
		server:           grpc.NewServer(),
		port:             port,
		transportService: transportService,
	}
}

// Start запускает gRPC сервер и регистрирует сервисы.
// ctx: контекст выполнения.
func (gs *GRPCServer) Start(ctx context.Context) error {
	gs.logger.Info("Initializing listener on port", zap.String("port", gs.port))
	listener, err := net.Listen("tcp", gs.port)
	if err != nil {
		gs.logger.Error("Failed to create listener", zap.Error(err))
		return err
	}
	gs.logger.Info("Listener successfully created on port", zap.String("port", gs.port))

	// Регистрируем сервис TransportService
	gs.logger.Info("Registering TransportService in gRPC server")
	pb.RegisterTransportServiceServer(gs.server, gs.transportService)
	gs.logger.Info("TransportService successfully registered")

	// Запускаем сервер в отдельной горутине
	go func() {
		gs.logger.Info("Starting gRPC server")
		if err := gs.server.Serve(listener); err != nil {
			gs.logger.Error("Failed to start gRPC server", zap.Error(err))
			log.Fatalf("failed to serve gRPC server: %v", err)
		}
	}()
	gs.logger.Info("gRPC server started and ready to accept requests")

	// Блокируем выполнение до завершения контекста
	<-ctx.Done()
	gs.logger.Info("Received shutdown signal, stopping server")
	return nil
}

// Stop останавливает gRPC сервер.
// ctx: контекст выполнения.
func (gs *GRPCServer) Stop(ctx context.Context) error {
	gs.logger.Info("Stopping gRPC server")
	gs.server.GracefulStop()
	gs.logger.Info("gRPC server successfully stopped")
	return nil
}

/*
NewGRPCServer создает новый экземпляр GRPC-сервера с предоставленными зависимостями.
port: порт, на котором будет слушать сервер.
logger: экземпляр интерфейса logger.Logger для логирования действий.
transportService: экземпляр транспортного сервиса для обработки данных.

Start запускает gRPC сервер и регистрирует сервисы.
ctx: контекст выполнения.

Stop останавливает gRPC сервер.
ctx: контекст выполнения.
*/
