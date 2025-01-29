package server

import "context"

// GRPCServerInterface определяет методы для работы с gRPC сервером.
type GRPCServerInterface interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context) error
}
