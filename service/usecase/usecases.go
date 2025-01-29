package usecase

import (
	database "shelon_server/integrations/SQLLite"
	youtubeclient "shelon_server/integrations/youtubeCLient"
	"shelon_server/utilss/logger"
	"sync"

	"go.uber.org/zap"
)

// BusinessLogic представляет слой бизнес-логики, который включает в себя зависимости:
// - Logger: логирование событий и ошибок.
// - YouTubeService: клиент для взаимодействия с API YouTube.
// - Sqlite: интерфейс для работы с базой данных SQLite.
type BusinessLogic struct {
	Logger         logger.Logger
	YouTubeService youtubeclient.YouTubeClient
	Sqlite         database.Database
}

// NewBusinessLogic создает и инициализирует объект BusinessLogic с переданными зависимостями.
// logger: экземпляр интерфейса logger.Logger для логирования действий.
// sqlite: экземпляр интерфейса database.Database для работы с базой данных.
// youTubeService: экземпляр интерфейса youtubeclient.YouTubeClient для взаимодействия с YouTube API.
func NewBusinessLogic(logger logger.Logger, sqlite database.Database, youTubeService youtubeclient.YouTubeClient) *BusinessLogic {
	return &BusinessLogic{
		Logger:         logger,
		Sqlite:         sqlite,
		YouTubeService: youTubeService,
	}
}

// ProcessData управляет обработкой списка ссылок. Если флаг "flag" установлен, данные обрабатываются асинхронно.
// Возвращает массив фотографий или ошибку.
func (bl *BusinessLogic) ProcessData(flag bool, links []string) ([][]byte, error) {
	bl.Logger.Info("Starting data processing", zap.Bool("Async", flag), zap.Int("Links count", len(links)))
	if flag {
		return bl.processAsync(links)
	} else {
		return bl.process(links)
	}
}

// processAsync обрабатывает ссылки в асинхронном режиме с использованием горутин и каналов.
func (bl *BusinessLogic) processAsync(links []string) ([][]byte, error) {
	bl.Logger.Info("Starting asynchronous processing of links")
	ch := make(chan []byte, len(links))
	errCh := make(chan error, len(links))

	var wg sync.WaitGroup
	for _, link := range links {
		wg.Add(1)
		go func(link string) {
			defer wg.Done()
			bl.Logger.Info("Processing link in goroutine", zap.String("Link", link))
			photo, err := bl.getPhotoOrFetch(link)
			if err != nil {
				bl.Logger.Error("Error in goroutine", zap.String("Link", link), zap.Error(err))
				errCh <- err
				return
			}
			ch <- photo
		}(link)
	}

	go func() {
		wg.Wait()
		close(ch)
		close(errCh)
	}()

	var photos [][]byte
	for {
		select {
		case photo, ok := <-ch:
			if !ok {
				ch = nil
			} else {
				bl.Logger.Info("Photo successfully processed in goroutine")
				photos = append(photos, photo)
			}
		case err, ok := <-errCh:
			if !ok {
				errCh = nil
			} else {
				bl.Logger.Error("Error during asynchronous processing", zap.Error(err))
				return nil, err
			}
		}
		if ch == nil && errCh == nil {
			break
		}
	}

	bl.Logger.Info("Asynchronous processing completed")
	return photos, nil
}

// process обрабатывает ссылки в синхронном режиме.
func (bl *BusinessLogic) process(links []string) ([][]byte, error) {
	bl.Logger.Info("Starting synchronous processing of links")
	var photos [][]byte
	for _, link := range links {
		bl.Logger.Info("Processing link", zap.String("Link", link))
		photo, err := bl.getPhotoOrFetch(link)
		if err != nil {
			bl.Logger.Error("Error processing link", zap.String("Link", link), zap.Error(err))
			continue
		}
		if photo != nil {
			bl.Logger.Info("Photo successfully loaded", zap.String("Link", link))
			photos = append(photos, photo)
		}
	}
	bl.Logger.Info("Synchronous processing completed")
	return photos, nil
}

// getPhotoOrFetch проверяет наличие фотографии в базе данных и возвращает её.
// Если фото отсутствует, обращается к YouTubeService и сохраняет результат в базе.
func (bl *BusinessLogic) getPhotoOrFetch(link string) ([]byte, error) {
	bl.Logger.Info("Checking photo in the database", zap.String("Link", link))
	// Проверяем наличие в базе и возвращаем фото, если оно есть
	photo, err := bl.Sqlite.GetPhotoByUrl(link)
	if err != nil {
		bl.Logger.Error("Error checking photo in the database", zap.String("Link", link), zap.Error(err))
		return nil, err
	}
	if photo != nil {
		bl.Logger.Info("Photo found in the database", zap.String("Link", link))
		return photo, nil
	}

	bl.Logger.Info("Photo not found in the database, fetching from YouTube API", zap.String("Link", link))
	// Если в базе нет, обращаемся к YouTube-клиенту
	photo, err = bl.YouTubeService.FetchThumbnail(link)
	if err != nil {
		bl.Logger.Error("Error fetching from YouTube API", zap.String("Link", link), zap.Error(err))
		return nil, err
	}

	// Сохраняем фото в базу
	bl.Logger.Info("Saving photo to the database", zap.String("Link", link))
	err = bl.Sqlite.InsertResource(link, photo)
	if err != nil {
		bl.Logger.Error("Error saving photo to the database for link", zap.String("Link", link), zap.Error(err))
	}

	bl.Logger.Info("Photo successfully processed and saved", zap.String("Link", link))
	return photo, nil
}

/*
NewBusinessLogic создает новый экземпляр BusinessLogic с предоставленными зависимостями.
logger: экземпляр интерфейса logger.Logger для логирования действий.
sqlite: экземпляр интерфейса database.Database для работы с базой данных.
youTubeService: экземпляр интерфейса youtubeclient.YouTubeClient для взаимодействия с YouTube API.

ProcessData управляет обработкой списка ссылок. Если флаг "flag" установлен, данные обрабатываются асинхронно.
Возвращает массив фотографий или ошибку.

processAsync обрабатывает ссылки в асинхронном режиме с использованием горутин и каналов.

process обрабатывает ссылки в синхронном режиме.

getPhotoOrFetch проверяет наличие фотографии в базе данных и возвращает её.
Если фото отсутствует, обращается к YouTubeService и сохраняет результат в базе.
*/
