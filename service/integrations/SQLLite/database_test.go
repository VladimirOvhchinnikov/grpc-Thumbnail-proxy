package database

import (
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
)

// MockLogger заглушка для логирования в тестах
type MockLogger struct{}

func (m *MockLogger) Info(message string, fields ...interface{})  {}
func (m *MockLogger) Warn(message string, fields ...interface{})  {}
func (m *MockLogger) Error(message string, fields ...interface{}) {}

// TestDatabaseOperations проверяет все операции с базой данных
func TestDatabaseOperations(t *testing.T) {
	// Используем временную базу данных в памяти
	dbFile := "test_database.db"
	defer os.Remove(dbFile) // Удаляем файл после тестов

	logger := &MockLogger{}
	db, err := NewSQLiteDatabase(logger, dbFile)
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	// Тест инициализации базы
	err = db.InitDatabase()
	if err != nil {
		t.Fatalf("Failed to initialize tables: %v", err)
	}

	// Тест добавления ресурса
	url := "https://www.example.com/image.jpg"
	photo := []byte{1, 2, 3, 4} // Заглушка фото
	err = db.InsertResource(url, photo)
	if err != nil {
		t.Fatalf("Failed to insert resource: %v", err)
	}

	// Тест проверки существования ресурса
	exists, err := db.ResourceExists(url)
	if err != nil {
		t.Fatalf("Failed to check resource existence: %v", err)
	}
	if !exists {
		t.Errorf("Expected resource to exist, but it does not")
	}

	// Тест получения фото
	retrievedPhoto, err := db.GetPhotoByUrl(url)
	if err != nil {
		t.Fatalf("Failed to retrieve photo: %v", err)
	}
	if len(retrievedPhoto) == 0 {
		t.Errorf("Expected non-empty photo, but got empty result")
	}

	// Тест закрытия базы данных
	err = db.Close()
	if err != nil {
		t.Errorf("Expected database to close successfully, got error: %v", err)
	}
}
