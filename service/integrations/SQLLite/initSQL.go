package database

import (
	"database/sql"
	"shelon_server/utilss/logger"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"go.uber.org/zap"
)

// SQLiteDatabase структура для работы с базой данных SQLite.
type SQLiteDatabase struct {
	Logger  logger.Logger
	DB      *sqlx.DB
	Builder squirrel.StatementBuilderType
}

// NewSQLiteDatabase создает новый экземпляр базы данных.
// logger: экземпляр интерфейса logger.Logger для логирования действий.
// dbName: имя файла базы данных.
func NewSQLiteDatabase(logger logger.Logger, dbName string) (*SQLiteDatabase, error) {
	db, err := sqlx.Open("sqlite3", dbName)
	if err != nil {
		return nil, err
	}
	return &SQLiteDatabase{
		Logger:  logger,
		DB:      db,
		Builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question),
	}, nil
}

// InitDatabase инициализирует базу данных, создавая таблицы при необходимости.
func (s *SQLiteDatabase) InitDatabase() error {
	createTableQuery := `
    CREATE TABLE IF NOT EXISTS resources (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        url TEXT NOT NULL,
        photo BLOB
    );`
	_, err := s.DB.Exec(createTableQuery)
	if err != nil {
		s.Logger.Error("Failed to initialize database", zap.Error(err))
		return err
	}
	s.Logger.Info("Database successfully initialized")
	return nil
}

// InsertResource добавляет новый ресурс в базу данных.
func (s *SQLiteDatabase) InsertResource(url string, photo []byte) error {
	query, args, err := s.Builder.
		Insert("resources").
		Columns("url", "photo").
		Values(url, photo).
		ToSql()
	if err != nil {
		s.Logger.Error("Failed to build insert query", zap.Error(err))
		return err
	}
	_, execErr := s.DB.Exec(query, args...)
	if execErr != nil {
		s.Logger.Error("Failed to execute insert query", zap.Error(execErr))
		return execErr
	}
	s.Logger.Info("Resource with URL added successfully", zap.String("url", url))
	return nil
}

// ResourceExists проверяет, существует ли ресурс с заданным URL.
func (s *SQLiteDatabase) ResourceExists(url string) (bool, error) {
	query, args, err := s.Builder.
		Select("COUNT(*)").
		From("resources").
		Where(squirrel.Eq{"url": url}).
		ToSql()
	if err != nil {
		s.Logger.Error("Failed to build existence check query", zap.Error(err))
		return false, err
	}
	var count int
	err = s.DB.Get(&count, query, args...)
	if err != nil {
		s.Logger.Error("Failed to execute existence check query", zap.Error(err))
		return false, err
	}
	return count > 0, nil
}

// Close закрывает соединение с базой данных.
func (s *SQLiteDatabase) Close() error {
	err := s.DB.Close()
	if err != nil {
		s.Logger.Error("Failed to close database connection", zap.Error(err))
		return err
	}
	s.Logger.Info("Database connection closed successfully")
	return nil
}

// GetPhotoByUrl получает фото по URL из базы данных.
func (s *SQLiteDatabase) GetPhotoByUrl(url string) ([]byte, error) {
	query, args, err := s.Builder.
		Select("photo").
		From("resources").
		Where(squirrel.Eq{"url": url}).
		ToSql()
	if err != nil {
		s.Logger.Error("Failed to build select query", zap.Error(err))
		return nil, err
	}
	var photo []byte
	err = s.DB.Get(&photo, query, args...)
	if err == sql.ErrNoRows {
		s.Logger.Info("No photo found for the given URL", zap.String("url", url))
		return nil, nil // Если фото не найдено, возвращаем nil
	} else if err != nil {
		s.Logger.Error("Failed to execute select query", zap.Error(err))
		return nil, err
	}
	s.Logger.Info("Photo retrieved successfully", zap.String("url", url))
	return photo, nil
}

/*
NewSQLiteDatabase создает новый экземпляр базы данных.
logger: экземпляр интерфейса logger.Logger для логирования действий.
dbName: имя файла базы данных.

InitDatabase инициализирует базу данных, создавая таблицы при необходимости.

InsertResource добавляет новый ресурс в базу данных.
url: URL ресурса.
photo: массив байтов фотографии.

ResourceExists проверяет, существует ли ресурс с заданным URL.
url: URL ресурса.

Close закрывает соединение с базой данных.

GetPhotoByUrl получает фото по URL из базы данных.
url: URL ресурса.
*/
