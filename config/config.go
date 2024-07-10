package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	serial "go.bug.st/serial"
)

const (
	ProdConfigFile     = "./config.toml"
	TestConfigFile     = "../config.test.toml"
	TelegramModeEnvVar = "TELE_MODE"
)

var (
	parameter *Parameter
)

// Parameter holds the configuration for the application.
type Parameter struct {
	Serial   SerialConfig
	NATS     NATSConfig
	Telegram TelegramConfig
	// Logger   LoggerConfig
}

// SerialConfig holds the serial port configuration.
type SerialConfig struct {
	Name        string `toml:"name"`
	Baud        int    `toml:"baud"`
	ReadTimeout int    `toml:"read_timeout"`
	Size        int    `toml:"size"`
	Parity      string `toml:"parity"`
	StopBits    int    `toml:"stop_bits"`
	FlowControl string `toml:"flow_control"`
	BufferSize  int    `toml:"buffer_size"`
}

// NATSConfig holds the NATS server configuration.
type NATSConfig struct {
	URL       string `toml:"url"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Subject   string `toml:"subject"`
	ClusterId string `toml:"cluster_id"`
	ClientId  string `toml:"client_id"`
}

// TelegramConfig holds the Telegram configuration.
type TelegramConfig struct {
	EndTag string `toml:"end_tag"`
	SeqTag string `toml:"seq_tag"`
}

// LoggerConfig holds the logger configuration.
type LoggerConfig struct {
	Level      string `toml:"level"`
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

// InitParameter initializes the global parameter by loading the config based on the environment.
func InitParameter() error {
	env := getEnv(TelegramModeEnvVar, "test")
	param, err := loadConfig(env)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	parameter = &param
	// fmt.Printf("Loaded config: %+v\n", param)
	return nil
}

// GetParameter returns the initialized global parameter.
func GetParameter() *Parameter {
	if parameter == nil {
		if err := InitParameter(); err != nil {
			panic(fmt.Sprintf("Failed to initialize parameters: %v", err))
		}
	}
	return parameter
}

// LoadConfigFromEnv loads the configuration parameters based on the TELE_MODE environment variable.
func LoadConfigFromEnv() (Parameter, error) {
	env := getEnv(TelegramModeEnvVar, "test")
	return loadConfig(env)
}

// getEnv retrieves the value of an environment variable or returns a default value if it's not set.
func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

// loadConfig loads the configuration parameters based on the specified environment.
func loadConfig(env string) (Parameter, error) {
	configFile := getConfigFileForEnv(env)

	file, err := os.Open(configFile)
	if err != nil {
		return Parameter{}, err
	}
	defer file.Close()

	var config Parameter
	if _, err := toml.NewDecoder(file).Decode(&config); err != nil {
		return Parameter{}, err
	}
	return config, nil
}

// getConfigFileForEnv returns the appropriate configuration file for the given environment.
func getConfigFileForEnv(env string) string {
	switch env {
	case "prod":
		return ProdConfigFile
	case "test":
		// fmt.Println("Loading test environment config")
		return TestConfigFile
	default:
		return ProdConfigFile
	}
}

// ReadSerialConfig converts the SerialConfig to a serial.Mode and returns the name of the serial port.
func ReadSerialConfig(serialConfig SerialConfig) (*serial.Mode, string) {
	return &serial.Mode{
		BaudRate: serialConfig.Baud,
		DataBits: serialConfig.Size,
		Parity:   parseParity(serialConfig.Parity),
		StopBits: parseStopBits(serialConfig.StopBits),
	}, serialConfig.Name
}

// parseParity converts the parity string to a corresponding serial.Parity value.
func parseParity(parity string) serial.Parity {
	switch parity {
	case "N":
		return serial.NoParity
	case "O":
		return serial.OddParity
	case "E":
		return serial.EvenParity
	default:
		return serial.NoParity
	}
}

// parseStopBits converts the stop bits int value to the corresponding serial.StopBits value.
func parseStopBits(stopBits int) serial.StopBits {
	switch stopBits {
	case 1:
		return serial.OneStopBit
	case 2:
		return serial.TwoStopBits
	default:
		return serial.OneStopBit
	}
}
