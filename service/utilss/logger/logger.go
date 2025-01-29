package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// ZapLogger реализует интерфейс Logger на базе zap.Logger.
type ZapLogger struct {
	logger *zap.Logger // Логгер zap для работы с логированием.
}

// NewZapLogger создает новый экземпляр ZapLogger, оборачивающий zap.Logger.
// Логи записываются в указанный файл.
func NewZapLogger(filePath string) (Logger, error) {
	// Создаем zap.Logger, записывающий в файл.
	zapLogger, err := NewFileLogger(filePath)
	if err != nil {
		return nil, err
	}
	// Возвращаем обертку ZapLogger, реализующую интерфейс Logger.
	return &ZapLogger{logger: zapLogger}, nil
}

// NewFileLogger создает zap.Logger, который записывает логи в указанный файл.
func NewFileLogger(filePath string) (*zap.Logger, error) {
	// Открываем файл для записи логов, создавая его, если он не существует.
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Создаем WriteSyncer для записи в файл.
	fileSyncer := zapcore.AddSync(file)

	// Настраиваем формат логов (JSON с читаемой датой).
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // Формат времени ISO 8601.

	// Создаем core для записи логов в файл с уровнем INFO.
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig), // JSON-формат логов.
		fileSyncer,                            // Запись в файл.
		zapcore.InfoLevel,                     // Уровень логирования.
	)

	// Возвращаем экземпляр zap.Logger с вызовами записывающего метода.
	return zap.New(core, zap.AddCaller()), nil
}

// Info логирует сообщение уровня INFO с произвольными полями.
func (zl *ZapLogger) Info(message string, fields ...interface{}) {
	zl.logger.Sugar().Infof(message, fields...)
}

// Warn логирует сообщение уровня WARN с произвольными полями.
func (zl *ZapLogger) Warn(message string, fields ...interface{}) {
	zl.logger.Sugar().Warnf(message, fields...)
}

// Error логирует сообщение уровня ERROR с произвольными полями.
func (zl *ZapLogger) Error(message string, fields ...interface{}) {
	zl.logger.Sugar().Errorf(message, fields...)
}
