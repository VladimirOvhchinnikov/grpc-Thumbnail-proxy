package youtubeclient

// YouTubeClient определяет интерфейс клиента для обработки ссылок YouTube.
type YouTubeClient interface {
	ProcessLinks(links []string) error
	FetchThumbnail(link string) ([]byte, error)
}
