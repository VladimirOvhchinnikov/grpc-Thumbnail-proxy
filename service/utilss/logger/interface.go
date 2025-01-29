package logger

// Logger описывает интерфейс для логирования действий в приложении.
type Logger interface {
	// Info записывает информационное сообщение (уровень INFO).
	Info(message string, fields ...interface{})
	// Warn записывает предупреждение (уровень WARN).
	Warn(message string, fields ...interface{})
	// Error записывает сообщение об ошибке (уровень ERROR).
	Error(message string, fields ...interface{})
}
