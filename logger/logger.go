package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	TestConfigFileName = "logger.test.json"
	ProdConfigFileName = "logger.json"
	EnvTest            = "test"
	EnvProd            = "prod"
)

// LoggerConfig represents the configuration for the logger.
type LoggerConfig struct {
	ZapConfig        zap.Config       `json:"zapConfig"`
	LumberjackConfig LumberjackConfig `json:"lumberjackConfig"`
}

// LumberjackConfig represents the configuration for lumberjack logging.
type LumberjackConfig struct {
	Filename   string `json:"filename"`
	MaxSize    int    `json:"maxSize"`
	MaxBackups int    `json:"maxBackups"`
	MaxAge     int    `json:"maxAge"`
	Compress   bool   `json:"compress"`
}

var (
	sugar    *zap.SugaredLogger
	initOnce sync.Once
)

// InitLogger initializes the logger configuration.
func InitLogger(loggerConfig LoggerConfig) {
	Init()
}

// Init initializes the logger. It ensures that the logger is initialized only once.
func Init() {
	initOnce.Do(func() {
		env := getEnv()
		configFile, err := getConfigFile(env)
		if err != nil {
			fmt.Printf("Error finding config file: %v\n", err)
			return
		}
		if err := InitLoggerFromFile(configFile, env); err != nil {
			fmt.Printf("Error initializing logger: %v\n", err)
		}
	})
}

// getEnv retrieves the logging environment from the TELE_MODE environment variable.
func getEnv() string {
	env := os.Getenv("TELE_MODE")
	if env == "" {
		return EnvTest // Default to test if no environment variable is set.
	}
	return env
}

// getConfigFile determines the configuration file path based on the environment.
func getConfigFile(env string) (string, error) {
	var configFileName string
	if env == EnvProd {
		configFileName = ProdConfigFileName
	} else {
		configFileName = TestConfigFileName
	}

	// Get the current working directory
	absPath, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("cannot determine current working directory: %w", err)
	}

	// Combine the current working directory with the config file name
	configFile := filepath.Join(absPath, configFileName)
	if _, err := os.Stat(configFile); err != nil {
		return "", fmt.Errorf("cannot find config file: %w", err)
	}

	return configFile, nil
}

// InitLoggerFromFile initializes the logger using the provided configuration file.
func InitLoggerFromFile(configFile, env string) error {
	// fmt.Printf("Loading config from file: %s\n", configFile)
	config, err := loadConfig(configFile)
	if err != nil {
		return fmt.Errorf("error loading config: %w", err)
	}

	var logWriter zapcore.WriteSyncer
	if env == EnvProd {
		logWriter = zapcore.AddSync(&lumberjack.Logger{
			Filename:   config.LumberjackConfig.Filename,
			MaxSize:    config.LumberjackConfig.MaxSize,
			MaxBackups: config.LumberjackConfig.MaxBackups,
			MaxAge:     config.LumberjackConfig.MaxAge,
			Compress:   config.LumberjackConfig.Compress,
		})
	} else {
		logWriter = zapcore.AddSync(os.Stdout)
	}

	encoder := zapcore.NewJSONEncoder(config.ZapConfig.EncoderConfig)
	level := parseLogLevel(config.ZapConfig.Level.String())

	core := zapcore.NewCore(
		encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), logWriter),
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	if logger == nil {
		return fmt.Errorf("failed to create logger instance")
	}

	sugar = logger.Sugar()
	return nil
}

// loadConfig loads the logger configuration from the given file.
func loadConfig(configFile string) (LoggerConfig, error) {
	// fmt.Printf("loading config: %s\n", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		fmt.Printf("error opening file: %v\n", err)
		return LoggerConfig{}, err
	}
	defer file.Close()

	var config LoggerConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return LoggerConfig{}, err
	}
	return config, nil
}

// SugaredLogger returns the initialized SugaredLogger instance.
func SugaredLogger() *zap.SugaredLogger {
	Init() // Ensure initialization
	return sugar
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
