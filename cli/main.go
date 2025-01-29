package main

import (
	"echelon_cli/commands"
	"echelon_cli/utils"
	"fmt"
	"log"
	"os"
)

// ParseCLIInput обрабатывает ввод из командной строки и возвращает флаг асинхронности и список ссылок.
func ParseCLIInput(logger utils.Logger) (bool, []string) {
	parserCLI := commands.NewParserConsole(logger)
	logger.Info("Command line parser initialized")

	// Парсинг флагов
	if err := parserCLI.ParseFlags(); err != nil {
		logger.Error("Failed to parse flags", err)
		fmt.Println("Check logs for details")
		os.Exit(1)
	}

	return parserCLI.GetParsedInput()
}

func main() {

	config, err := utils.LoadConfig("utils/config.json")
	if err != nil {
		log.Fatalf("Ошибка загрузки конфигурации: %v", err)
	}

	// Инициализация логгера
	logger, err := utils.NewZapLogger(config.LogFilePath)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	// Парсинг ввода из командной строки
	async, links := ParseCLIInput(logger)
	logger.Info("Command line parsing completed successfully")

	// Тестовый вывод
	fmt.Printf("Async: %t\n", async)
	fmt.Printf("Links: %v\n", links)

	// Адрес gRPC-сервера
	serverAddress := config.GrpcServerAddress

	// Создание клиента
	client := &utils.GRPCTransportSender{
		Logger: logger,
	}
	logger.Info("Client created successfully")

	// Подключение к серверу
	if err := client.Connect(serverAddress); err != nil {
		logger.Error("Failed to connect to server", err)
		return
	}
	defer func() {
		if err := client.Close(); err != nil {
			logger.Error("Failed to close connection", err)
		} else {
			logger.Info("Connection closed successfully")
		}
	}()

	// Отправка данных
	logger.Info("Starting data transmission to server")
	if err := client.SendData(async, links); err != nil {
		logger.Error("Failed to send data", err)
		return
	}
	logger.Info("Data sent successfully")
}
