package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type LoggerConfig struct {
	Level      string
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

const (
	EnvTest = "test"
	EnvProd = "prod"
)

var (
	sugar    *zap.SugaredLogger
	config   LoggerConfig
	initOnce sync.Once
	env      string
)

// InitLogger initializes the logger configuration.
func InitLogger(loggerConfig LoggerConfig) {
	config = loggerConfig
	Init()
}

// Init initializes the logger. It ensures that the logger is initialized only once.
func Init() {
	initOnce.Do(func() {
		env = os.Getenv("TELE_MODE")
		if env == "" {
			env = EnvTest // Default to test if no environment variable is set
		}

		switch env {
		case EnvProd:
			InitProductionLogger()
		default:
			InitTestLogger()
		}
	})
}

// InitTestLogger initializes the logger for the testing environment.
func InitTestLogger() {
	initLogger(zapcore.AddSync(os.Stdout))
}

// InitProductionLogger initializes the logger for the production environment.
func InitProductionLogger() {
	initLogger(nil)
}

// initLogger initializes the SugaredLogger instance based on the provided configuration parameters.
func initLogger(output zapcore.WriteSyncer) {
	if config == (LoggerConfig{}) {
		panic("logger configuration is not set")
	}

	encoderConfig := buildEncoderConfig()
	logWriter := getLogWriter(output)

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter),
		parseLogLevel(config.Level),
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	sugar = logger.Sugar()
}

// buildEncoderConfig builds the encoder configuration for the logger.
func buildEncoderConfig() zapcore.EncoderConfig {
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
	return encoderConfig
}

// getLogWriter returns the log writer based on the output provided.
func getLogWriter(output zapcore.WriteSyncer) zapcore.WriteSyncer {
	if output == nil {
		return zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,    // megabytes after which new file is created
			MaxBackups: config.MaxBackups, // number of backups to keep
			MaxAge:     config.MaxAge,     // days to keep the log files
			Compress:   config.Compress,   // whether to compress the log files
		})
	}
	return output
}

// parseLogLevel converts the log level string to zapcore.Level.
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
	Init() // Ensure initialization
	return sugar
}
