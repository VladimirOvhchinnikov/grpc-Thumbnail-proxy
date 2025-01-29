package youtubeclient

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"shelon_server/utilss/logger"

	"go.uber.org/zap"
)

// YouTubeService представляет реализацию YouTubeClient.
type YouTubeService struct {
	Logger logger.Logger
	Client *http.Client
}

// NewYouTubeService создает и настраивает YouTubeService с использованием прокси.
// logger: экземпляр интерфейса logger.Logger для логирования действий.
// proxyURL: URL прокси сервера.
// proxyUser: имя пользователя для аутентификации на прокси сервере.
// proxyPass: пароль для аутентификации на прокси сервере.
func NewYouTubeService(logger logger.Logger, proxyURL, proxyUser, proxyPass string) (*YouTubeService, error) {
	// Разбираем URL прокси
	parsedProxyURL, err := url.Parse(proxyURL)
	if err != nil {
		logger.Error("Invalid proxy URL format", zap.Error(err))
		return nil, errors.New("invalid proxy URL format")
	}
	// Устанавливаем аутентификацию, если логин и пароль заданы
	if proxyUser != "" && proxyPass != "" {
		parsedProxyURL.User = url.UserPassword(proxyUser, proxyPass)
	}
	// Создаем транспорт с прокси
	transport := &http.Transport{
		Proxy: http.ProxyURL(parsedProxyURL),
	}
	// Создаем HTTP-клиент с указанным транспортом
	client := &http.Client{
		Transport: transport,
	}
	return &YouTubeService{
		Logger: logger,
		Client: client,
	}, nil
}

// ProcessLinks выполняет обработку ссылок YouTube.
// links: список ссылок на видео YouTube.
func (ys *YouTubeService) ProcessLinks(links []string) error {
	for _, link := range links {
		_, err := ys.FetchThumbnail(link)
		if err != nil {
			ys.Logger.Error("Failed to process link", zap.String("link", link), zap.Error(err))
			return fmt.Errorf("failed to process link %s: %w", link, err)
		}
		ys.Logger.Info("Link processed successfully", zap.String("link", link))
	}
	return nil
}

// FetchThumbnail загружает картинку по указанной ссылке через прокси.
// link: URL изображения для загрузки.
func (ys *YouTubeService) FetchThumbnail(link string) ([]byte, error) {
	thumbnailLink, err := ys.GenerateThumbnailURL(link)
	if err != nil {
		ys.Logger.Error("Failed to generate thumbnail URL", zap.String("link", link), zap.Error(err))
		return nil, err
	}

	ys.Logger.Info("Starting thumbnail download", zap.String("link", thumbnailLink))
	// Выполняем HTTP-запрос
	resp, err := ys.Client.Get(thumbnailLink)
	if err != nil {
		ys.Logger.Error("Failed to download link", zap.String("link", thumbnailLink), zap.Error(err))
		return nil, fmt.Errorf("failed to connect to URL %s: %w", thumbnailLink, err)
	}
	defer resp.Body.Close()
	// Проверяем статус ответа
	if resp.StatusCode != http.StatusOK {
		ys.Logger.Error("Unexpected HTTP status", zap.String("status", resp.Status), zap.String("link", thumbnailLink))
		return nil, fmt.Errorf("unexpected HTTP status %d for URL %s", resp.StatusCode, thumbnailLink)
	}
	// Читаем содержимое ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		ys.Logger.Error("Failed to read response body", zap.Error(err), zap.String("link", thumbnailLink))
		return nil, fmt.Errorf("failed to read data from URL %s: %w", thumbnailLink, err)
	}
	ys.Logger.Info("Thumbnail downloaded successfully", zap.String("link", thumbnailLink))
	return body, nil
}

// GenerateThumbnailURL генерирует URL обложки для указанной ссылки YouTube.
func (ys *YouTubeService) GenerateThumbnailURL(videoURL string) (string, error) {
	videoID, err := extractVideoID(videoURL)
	if err != nil {
		ys.Logger.Error("Failed to extract video ID", zap.String("url", videoURL), zap.Error(err))
		return "", err
	}
	thumbnailURL := fmt.Sprintf("https://img.youtube.com/vi/%s/maxresdefault.jpg", videoID)
	ys.Logger.Info("Generated thumbnail URL", zap.String("videoID", videoID), zap.String("thumbnailURL", thumbnailURL))
	return thumbnailURL, nil
}

// extractVideoID извлекает идентификатор видео из стандартной ссылки YouTube.
func extractVideoID(videoURL string) (string, error) {
	// Парсим URL
	parsedURL, err := url.Parse(videoURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}

	var videoID string

	// Обработка стандартных ссылок типа "youtube.com/watch?v=..."
	if parsedURL.Host == "www.youtube.com" || parsedURL.Host == "youtube.com" {
		queryParams := parsedURL.Query()
		if id := queryParams.Get("v"); id != "" {
			videoID = id
		} else {
			return "", errors.New("could not find video ID in URL")
		}
	} else if parsedURL.Host == "youtu.be" {
		// Обработка коротких ссылок типа "youtu.be/..."
		pathSegments := strings.Split(parsedURL.Path, "/")
		if len(pathSegments) > 1 {
			videoID = pathSegments[1]
		} else {
			return "", errors.New("could not find video ID in short URL")
		}
	} else {
		return "", errors.New("unsupported URL format")
	}

	if videoID == "" {
		return "", errors.New("could not extract video ID from URL")
	}

	return videoID, nil
}

/*
NewYouTubeService создает и настраивает YouTubeService с использованием прокси.
logger: экземпляр интерфейса logger.Logger для логирования действий.
proxyURL: URL прокси сервера.
proxyUser: имя пользователя для аутентификации на прокси сервере.
proxyPass: пароль для аутентификации на прокси сервере.

ProcessLinks выполняет обработку ссылок YouTube.
links: список ссылок на видео YouTube.

FetchThumbnail загружает картинку по указанной ссылке через прокси.
link: URL изображения для загрузки.

GenerateThumbnailURL генерирует URL обложки для указанной ссылки YouTube.

extractVideoID извлекает идентификатор видео из стандартной ссылки YouTube.
*/
