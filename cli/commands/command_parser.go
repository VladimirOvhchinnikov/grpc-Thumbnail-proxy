package commands

import (
	"echelon_cli/utils"
	"errors"
	"flag"
	"fmt"
	"net/url"
	"strings"

	"go.uber.org/zap"
)

// CommandParser определяет интерфейс для парсинга консольных флагов и аргументов.
type CommandParser interface {
	// ParseFlags парсит консольные флаги и аргументы, сохраняет результаты в структуре.
	ParseFlags() error
	// GetParsedInput возвращает сохраненные данные: флаг (асинхронность) и список ссылок.
	GetParsedInput() (bool, []string)
}

// ParserConsole реализует интерфейс CommandParser и обрабатывает консольные данные.
type ParserConsole struct {
	isAsync bool         // Указывает, включен ли асинхронный режим (--async).
	links   []string     // Список ссылок, переданных через консоль.
	logger  utils.Logger // Логгер для записи событий.
}

// NewParserConsole создает новый экземпляр ParserConsole с предоставленным логгером.
// logger: экземпляр интерфейса utils.Logger для логирования действий.
func NewParserConsole(logger utils.Logger) *ParserConsole {
	return &ParserConsole{logger: logger}
}

// ParseFlags парсит флаги и аргументы из консоли, заполняя поля структуры ParserConsole.
// Флаг --async включает асинхронный режим.
// Флаг --links позволяет передать список ссылок, разделенных запятой.
// Если ссылки не переданы через --links, они извлекаются из оставшихся аргументов.
// Возвращает ошибку, если список ссылок пуст.
func (pc *ParserConsole) ParseFlags() error {
	// Определение флагов
	asyncFlag := flag.Bool("async", false, "Enable asynchronous mode for downloads")
	linksFlag := flag.String("links", "", "Comma-separated list of video URLs")

	// Парсинг флагов
	flag.Parse()

	// Сохранение результатов в структуру
	pc.isAsync = *asyncFlag

	if *linksFlag != "" {
		pc.links = strings.Split(*linksFlag, ",")
	} else {
		pc.links = flag.Args()
	}

	// Проверка корректности ссылок
	if err := pc.validateLinks(pc.links); err != nil {
		pc.logger.Error("Failed to validate links", zap.Error(err))
		return fmt.Errorf("failed to validate links: %w", err)
	}

	// Проверка наличия ссылок
	if len(pc.links) == 0 {
		pc.logger.Error("No video links provided")
		return errors.New("no video links provided")
	}

	pc.logger.Info("Successfully parsed command line flags and arguments")
	return nil
}

// GetParsedInput возвращает результаты парсинга:
// - Флаг асинхронности (isAsync).
// - Список ссылок (links).
func (pc *ParserConsole) GetParsedInput() (bool, []string) {
	return pc.isAsync, pc.links
}

// validateLinks проверяет корректность URL-адресов.
// Возвращает ошибку, если хотя бы одна ссылка некорректна.
func (pc *ParserConsole) validateLinks(links []string) error {
	for _, link := range links {
		parsedURL, err := url.Parse(link)
		if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
			pc.logger.Error("Invalid link detected", zap.String("link", link), zap.Error(err))
			return fmt.Errorf("invalid link: %s", link)
		}
	}
	return nil
}

/*
NewParserConsole создает новый экземпляр ParserConsole с предоставленным логгером.
logger: экземпляр интерфейса utils.Logger для логирования действий.

ParseFlags парсит флаги и аргументы из консоли, заполняя поля структуры ParserConsole.
Флаг --async включает асинхронный режим.
Флаг --links позволяет передать список ссылок, разделенных запятой.
Если ссылки не переданы через --links, они извлекаются из оставшихся аргументов.
Возвращает ошибку, если список ссылок пуст.

GetParsedInput возвращает результаты парсинга:
- Флаг асинхронности (isAsync).
- Список ссылок (links).

validateLinks проверяет корректность URL-адресов.
Возвращает ошибку, если хотя бы одна ссылка некорректна.
*/
