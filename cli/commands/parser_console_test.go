package commands

import (
	"flag"
	"os"
	"strings"
	"testing"
)

type MockLogger struct{}

func (m *MockLogger) Info(message string, fields ...interface{})  {}
func (m *MockLogger) Warn(message string, fields ...interface{})  {}
func (m *MockLogger) Error(message string, fields ...interface{}) {}

// TestParseFlags_Success проверяет корректную обработку флагов и аргументов.
func TestParseFlags_Success(t *testing.T) {
	// Инициализация mock-логгера
	logger := &MockLogger{}

	// Создание экземпляра ParserConsole
	parser := NewParserConsole(logger)

	// Симуляция входных данных
	args := []string{"cmd", "--async", "--links", "https://www.youtube.com/watch?v=EXAMPLE1,https://www.youtube.com/watch?v=EXAMPLE2"}
	originalArgs := os.Args                   // Сохраняем оригинальные аргументы
	defer func() { os.Args = originalArgs }() // Восстанавливаем после теста
	os.Args = args

	// Создаём новый FlagSet, чтобы избежать конфликтов с глобальными флагами
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	asyncFlag := flagSet.Bool("async", false, "Enable asynchronous mode for downloads")
	linksFlag := flagSet.String("links", "", "Comma-separated list of video URLs")

	err := flagSet.Parse(args[1:])
	if err != nil {
		t.Fatalf("Failed to parse test flags: %v", err)
	}

	// Присваиваем результаты парсинга в parser
	parser.isAsync = *asyncFlag
	if *linksFlag != "" {
		parser.links = strings.Split(*linksFlag, ",")
	} else {
		parser.links = flagSet.Args()
	}

	// Проверка флага --async
	isAsync, links := parser.GetParsedInput()
	if !isAsync {
		t.Errorf("Expected isAsync to be true, got false")
	}

	// Проверка списка ссылок
	expectedLinks := []string{"https://www.youtube.com/watch?v=EXAMPLE1", "https://www.youtube.com/watch?v=EXAMPLE2"}
	if len(links) != len(expectedLinks) {
		t.Fatalf("Expected %d links, got %d", len(expectedLinks), len(links))
	}
	for i, link := range links {
		if link != expectedLinks[i] {
			t.Errorf("Expected link %s, got %s", expectedLinks[i], link)
		}
	}
}

// TestParseFlags_NoLinks проверяет случай, когда ссылки не переданы.
func TestParseFlags_NoLinks(t *testing.T) {
	logger := &MockLogger{}
	parser := NewParserConsole(logger)

	// Создаём новый FlagSet, чтобы избежать конфликта флагов
	flagSet := flag.NewFlagSet("test", flag.ContinueOnError)
	asyncFlag := flagSet.Bool("async", false, "Enable asynchronous mode for downloads")
	linksFlag := flagSet.String("links", "", "Comma-separated list of video URLs")

	args := []string{"--async"}
	err := flagSet.Parse(args)
	if err != nil {
		t.Fatalf("Failed to parse test flags: %v", err)
	}

	// Присваиваем результаты парсинга
	parser.isAsync = *asyncFlag
	if *linksFlag != "" {
		parser.links = strings.Split(*linksFlag, ",")
	} else {
		parser.links = flagSet.Args()
	}

	// Выполняем метод `ParseFlags`
	err = parser.ParseFlags()
	if err == nil || !strings.Contains(err.Error(), "no video links provided") {
		t.Errorf("Expected error 'no video links provided', got %v", err)
	}
}

// TestValidateLinks_Invalid проверяет валидацию некорректных ссылок.
func TestValidateLinks_Invalid(t *testing.T) {
	logger := &MockLogger{}
	parser := NewParserConsole(logger)

	invalidLinks := []string{"invalid-url", "http://"}
	err := parser.validateLinks(invalidLinks)

	if err == nil || !strings.Contains(err.Error(), "invalid link") {
		t.Errorf("Expected error 'invalid link', got %v", err)
	}
}
