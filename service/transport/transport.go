package transport

import (
	"context"

	"shelon_server/handlers"
	pb "shelon_server/proto"
)

// TransportService представляет реализацию gRPC-сервиса для обработки данных.
type TransportService struct {
	pb.UnimplementedTransportServiceServer
	handler *handlers.DataHandler
}

// NewTransportService создает новый экземпляр TransportService с предоставленным обработчиком данных.
// handler: экземпляр обработчика данных.
func NewTransportService(handler *handlers.DataHandler) *TransportService {
	return &TransportService{handler: handler}
}

// SendData обрабатывает запрос на отправку данных через gRPC.
// ctx: контекст выполнения.
// req: запрос на отправку данных в формате proto.
// Возвращает ответ на отправку данных в формате proto и ошибку, если она возникла.
func (ts *TransportService) SendData(ctx context.Context, req *pb.SendDataRequest) (*pb.SendDataResponse, error) {
	return ts.handler.HandleSendData(ctx, req)
}

/*
NewTransportService создает новый экземпляр TransportService с предоставленным обработчиком данных.
handler: экземпляр обработчика данных.

SendData обрабатывает запрос на отправку данных через gRPC.
ctx: контекст выполнения.
req: запрос на отправку данных в формате proto.
Возвращает ответ на отправку данных в формате proto и ошибку, если она возникла.
*/
