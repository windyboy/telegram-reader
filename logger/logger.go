package logger

import (
	"fmt"
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gzzn.com/airport/serial/config"
)

var (
	sugar     *zap.SugaredLogger
	parameter *config.Parameter
	initOnce  sync.Once
)

// SetParameter sets the configuration parameters for the logger.
func SetParameter(param *config.Parameter) {
	parameter = param
}

func InitLoggerWithMode() {
	env := os.Getenv(config.TELEGRAM_MODE)
	if env == "" {
		env = "prod" // Default to production if no environment variable is set
	}

	switch env {
	case "test":
		InitTestLogger()
	default:
		InitLogger(nil)
	}
}

func InitTestLogger() {
	fmt.Println("InitTestLogger")
	InitLogger(zapcore.AddSync(os.Stdout))
}

// InitLogger initializes the SugaredLogger instance based on the provided configuration parameters.
func InitLogger(output zapcore.WriteSyncer) {
	if parameter == nil {
		panic("logger parameter is not set")
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.LevelKey = "level"
	encoderConfig.NameKey = "logger"
	encoderConfig.CallerKey = "caller"
	encoderConfig.MessageKey = "msg"
	encoderConfig.StacktraceKey = "stacktrace"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder

	var logWriter zapcore.WriteSyncer
	if output == nil {
		logWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   parameter.Logger.Filename,
			MaxSize:    parameter.Logger.MaxSize,    // megabytes after which new file is created
			MaxBackups: parameter.Logger.MaxBackups, // number of backups to keep
			MaxAge:     parameter.Logger.MaxAge,     // days to keep the log files
			Compress:   parameter.Logger.Compress,   // whether to compress the log files
		})
	} else {
		logWriter = output
	}

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter),
		parseLogLevel(parameter.Logger.Level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	sugar = logger.Sugar()
}

// parseLogLevel converts the log level string to zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "dpanic":
		return zapcore.DPanicLevel
	case "panic":
		return zapcore.PanicLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// SugaredLogger returns the initialized SugaredLogger instance.
// It ensures that the logger is initialized only once.
func SugaredLogger() *zap.SugaredLogger {
	initOnce.Do(func() { InitLoggerWithMode() }) // Default initialization with mode
	return sugar
}
