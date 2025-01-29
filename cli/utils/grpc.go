package utils

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"echelon_cli/transport"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// TransportSender описывает интерфейс для отправки данных (флага и ссылок) микросервису
type TransportSender interface {
	// Connect устанавливает соединение с сервером
	Connect(address string) error
	// SendData отправляет флаг и ссылки на сервер
	SendData(flag bool, links []string) error
	// Close закрывает соединение с сервером
	Close() error
}

// GRPCTransportSender структура для работы с gRPC транспортом
type GRPCTransportSender struct {
	Logger Logger
	conn   *grpc.ClientConn
}

// Connect устанавливает соединение с gRPC-сервером
func (gc *GRPCTransportSender) Connect(address string) error {
	// Установление соединения
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()), grpc.WithBlock())
	if err != nil {
		return fmt.Errorf("не удалось подключиться к серверу: %w", err)
	}
	// Привязка соединения к структуре
	gc.conn = conn
	gc.Logger.Info("Коннект прошел успешно")
	// Проверка состояния соединения
	if conn.GetState().String() != "READY" {
		return fmt.Errorf("сервер недоступен: состояние %s", conn.GetState().String())
	}
	return nil
}

// SendData отправляет данные на сервер
func (gc *GRPCTransportSender) SendData(flag bool, links []string) error {
	// Создание gRPC клиента
	client := transport.NewTransportServiceClient(gc.conn)
	// Формирование запроса
	req := &transport.SendDataRequest{
		Flag:  flag,
		Links: links,
	}
	// Отправка запроса
	resp, err := client.SendData(context.Background(), req)
	if err != nil {
		return fmt.Errorf("ошибка при отправке данных: %w", err)
	}
	// Логирование ответа
	log.Printf("Ответ от сервера: %s", resp.Status)

	// Директория для сохранения файлов
	saveDir := "./thumbnails"
	err = os.MkdirAll(saveDir, os.ModePerm)
	if err != nil {
		return fmt.Errorf("не удалось создать директорию для сохранения файлов: %w", err)
	}

	for i, img := range resp.Images {
		videoID := extractVideoID(links[i])
		fileName := fmt.Sprintf("%s.jpeg", videoID)
		filePath := filepath.Join(saveDir, fileName)
		err := os.WriteFile(filePath, img, 0644)
		if err != nil {
			log.Printf("Ошибка сохранения картинки %s: %v", filePath, err)
			continue
		}
		log.Printf("Картинка сохранена: %s", filePath)
	}
	return nil
}

// Close закрывает соединение
func (gc *GRPCTransportSender) Close() error {
	if gc.conn != nil {
		err := gc.conn.Close()
		if err != nil {
			return fmt.Errorf("не удалось закрыть соединение: %w", err)
		}
		gc.Logger.Info("Соединение закрыто")
	}
	return nil
}

// extractVideoID извлекает идентификатор видео из стандартной ссылки YouTube.
func extractVideoID(videoURL string) string {
	// Парсим URL
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		log.Printf("Не удалось распарсить URL: %v", err)
		return ""
	}

	var videoID string

	// Обработка стандартных ссылок типа "youtube.com/watch?v=..."
	if parsedURL.Host == "www.youtube.com" || parsedURL.Host == "youtube.com" {
		queryParams := parsedURL.Query()
		if id := queryParams.Get("v"); id != "" {
			videoID = id
		} else {
			log.Println("Не удалось найти ID видео в URL")
			return ""
		}
	} else if parsedURL.Host == "youtu.be" {
		// Обработка коротких ссылок типа "youtu.be/..."
		pathSegments := strings.Split(parsedURL.Path, "/")
		if len(pathSegments) > 1 {
			videoID = pathSegments[1]
		} else {
			log.Println("Не удалось найти ID видео в короткой ссылке")
			return ""
		}
	} else {
		log.Println("Неподдерживаемый формат URL")
		return ""
	}

	return videoID
}
