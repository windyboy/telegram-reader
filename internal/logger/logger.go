package logger

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

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
	sugar *zap.SugaredLogger
	log   *zap.Logger
	// once  sync.Once
)

// load initializes the logger. It ensures that the logger is initialized only once.
func load() {
	if log == nil {
		env := getEnv()
		fmt.Printf("Enviroment : %s\n", env)
		configFile, err := getConfigFile(env)
		// fmt.Printf("Config File : %s\n", configFile)
		if err != nil {
			fmt.Printf("Error finding config file: %v\n", err)
			return
		}
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		// fmt.Printf("Loading config from file: %s\n", configFile)
		// config, err := loadConfig(configFile)
		file := filepath.Join(currentDir, "internal/logger", configFile)
		config, err := loadConfig(file)

		if err != nil {
			fmt.Printf("Error loading config: %v\n", err)
			return
		}

		var logWriter zapcore.WriteSyncer
		if env == EnvProd {
			// Ensure log directory exists
			logDir := filepath.Dir(config.LumberjackConfig.Filename)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				fmt.Printf("Error creating log directory: %v\n", err)
				return
			}

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
			logWriter,
			level,
		)

		log = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
		sugar = log.Sugar()
		// fmt.Println("Logger initialized")
	}
}

// getEnv retrieves the logging environment from the TELE_MODE environment variable.
func getEnv() string {
	env := os.Getenv("TELE_MODE")
	if env == "" {
		env = EnvTest
	}
	return env
}

// getConfigFile determines the configuration file path based on the environment.
func getConfigFile(env string) (string, error) {
	switch env {
	case EnvTest:
		return TestConfigFileName, nil
	case EnvProd:
		return ProdConfigFileName, nil
	default:
		return "", fmt.Errorf("unknown environment: %s", env)
	}
}

// loadConfig loads the logger configuration from the given file.
func loadConfig(configFile string) (LoggerConfig, error) {
	file, err := os.Open(configFile)
	if err != nil {
		return LoggerConfig{}, fmt.Errorf("error opening file: %v", err)
	}
	defer file.Close()

	var config LoggerConfig
	if err := json.NewDecoder(file).Decode(&config); err != nil {
		return LoggerConfig{}, fmt.Errorf("error decoding config: %v", err)
	}
	return config, nil
}

// GetLogger returns the initialized SugaredLogger instance.
func GetLogger() *zap.SugaredLogger {
	if sugar == nil {
		load()
	}
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
