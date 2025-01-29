package handlers

import (
	"context"
	"fmt"

	pb "shelon_server/proto"
	"shelon_server/usecase"
	"shelon_server/utilss/logger"

	"go.uber.org/zap"
)

// DataHandler структура для обработки данных.
type DataHandler struct {
	logger        logger.Logger
	BusinessLogic usecase.DataProcessorUsecase
}

// NewDataHandler создает новый экземпляр DataHandler с предоставленными зависимостями.
// logger: экземпляр интерфейса logger.Logger для логирования действий.
// businessLogic: экземпляр бизнес-логики.
func NewDataHandler(logger logger.Logger, businessLogic usecase.DataProcessorUsecase) *DataHandler {
	return &DataHandler{
		logger:        logger,
		BusinessLogic: businessLogic,
	}
}

// HandleSendData обрабатывает запрос на отправку данных.
// ctx: контекст выполнения.
// req: запрос на отправку данных в формате proto.
// Возвращает ответ на отправку данных в формате proto и ошибку, если она возникла.
func (dh *DataHandler) HandleSendData(ctx context.Context, req *pb.SendDataRequest) (*pb.SendDataResponse, error) {
	// Логируем входные данные
	dh.logger.Info("Received SendData request")
	dh.logger.Info("Flag", zap.Bool("flag", req.Flag))
	dh.logger.Info("Links", zap.Strings("links", req.Links))

	// Вызываем бизнес-логику
	result, err := dh.BusinessLogic.ProcessData(req.Flag, req.Links)
	if err != nil {
		dh.logger.Error("Failed to process data", zap.Error(err))
		return nil, fmt.Errorf("failed to process data: %w", err)
	}

	// Конвертируем результат в формат, который клиент сможет обработать
	images := append([][]byte(nil), result...)

	// Формируем успешный ответ с массивом картинок
	dh.logger.Info("Forming successful response with image array")

	return &pb.SendDataResponse{
		Status: "success",
		Images: images,
	}, nil
}

/*
NewDataHandler создает новый экземпляр DataHandler с предоставленными зависимостями.
logger: экземпляр интерфейса logger.Logger для логирования действий.
businessLogic: экземпляр бизнес-логики.

HandleSendData обрабатывает запрос на отправку данных.
ctx: контекст выполнения.
req: запрос на отправку данных в формате proto.
Возвращает ответ на отправку данных в формате proto и ошибку, если она возникла.
*/
