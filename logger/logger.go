package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gzzn.com/airport/serial/config"
)

var sugar *zap.SugaredLogger

var parameter *config.Parameter

func SetParameter(param *config.Parameter) {
	parameter = param
}

func InitLogger() {
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

	logWriter := zapcore.AddSync(&lumberjack.Logger{
		Filename:   parameter.Logger.Filename,
		MaxSize:    parameter.Logger.MaxSize,    // megabytes after which new file is created
		MaxBackups: parameter.Logger.MaxBackups, // number of backups to keep
		MaxAge:     parameter.Logger.MaxAge,     // days to keep the log files
		Compress:   parameter.Logger.Compress,   // whether to compress the log files
	})

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter),
		zapcore.InfoLevel,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	defer logger.Sync()

	sugar = logger.Sugar()
}

// SugaredLogger returns the initialized SugaredLogger instance
func SugaredLogger() *zap.SugaredLogger {
	if sugar == nil {
		InitLogger()
	}
	return sugar
}
