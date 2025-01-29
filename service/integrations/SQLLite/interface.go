package database

// Database определяет интерфейс для взаимодействия с базой данных.
type Database interface {
	InitDatabase() error
	InsertResource(url string, photo []byte) error
	ResourceExists(url string) (bool, error)
	GetPhotoByUrl(url string) ([]byte, error)
	Close() error
}
