package main

import (
	"context"
	"log"
	"os"
	"shelon_server/handlers"
	database "shelon_server/integrations/SQLLite"
	youtubeclient "shelon_server/integrations/youtubeCLient"
	"shelon_server/transport"
	"shelon_server/usecase"
	"shelon_server/utilss/config"
	"shelon_server/utilss/logger"
	"shelon_server/utilss/server"

	"go.uber.org/zap"
)

func main() {
	config, err := config.LoadConfig("utilss/config/config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	loggerInstance, err := logger.NewZapLogger(config.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Инициализация базы данных
	sqliteDB, err := database.NewSQLiteDatabase(loggerInstance, config.Database.DataSourceName)
	if err != nil {
		loggerInstance.Error("Error initializing database", zap.Error(err))
		os.Exit(1)
	}
	if err := sqliteDB.InitDatabase(); err != nil {
		loggerInstance.Error("Error creating tables", zap.Error(err))
		os.Exit(1)
	}

	// Инициализация клиента YouTube
	youtubeConnect, err := youtubeclient.NewYouTubeService(
		loggerInstance,
		config.YoutubeClient.BaseURL,
		config.YoutubeClient.ClientID,
		config.YoutubeClient.ClientSecret,
	)
	if err != nil {
		loggerInstance.Error("Error creating YouTube client", zap.Error(err))
		os.Exit(1)
	}

	// Инициализация бизнес-логики
	businessLogic := usecase.NewBusinessLogic(loggerInstance, sqliteDB, youtubeConnect)

	// Инициализация обработчиков
	dataHandler := handlers.NewDataHandler(loggerInstance, businessLogic)

	// Инициализация транспортного сервиса
	transportService := transport.NewTransportService(dataHandler)

	// Инициализация gRPC сервера
	serverInstance := server.NewGRPCServer(config.GRPCServerAddress, loggerInstance, transportService)
	if err := serverInstance.Start(context.TODO()); err != nil {
		loggerInstance.Error("Error starting gRPC server", zap.Error(err))
		os.Exit(1)
	}

	loggerInstance.Info("Server started successfully")
}
