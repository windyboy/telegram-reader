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

var parameter *Parameter
var initialized = false

// Parameter holds the configuration for the application.
type Parameter struct {
	Serial   SerialConfig
	NATS     NATSConfig
	Telegram TelegramConfig
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
	URLS      string `toml:"urls"`
	Username  string `toml:"username"`
	Password  string `toml:"password"`
	Subject   string `toml:"subject"`
	ClusterId string `toml:"cluster_id"`
	ClientId  string `toml:"client_id"`
}

// TelegramConfig holds the Telegram configuration.
type TelegramConfig struct {
	EndTag       string `toml:"end_tag"`
	SeqTag       string `toml:"seq_tag"`
	PatternSplit string `toml:"pattern_split"`
}

// This function is idempotent, meaning it can be called multiple times without reinitializing the configuration.
// Init initializes the configuration for the application.
// It loads the configuration based on the environment variable TelegramModeEnvVar.
// If the environment variable is not set, it defaults to "test".
// The loaded configuration is stored in the parameter variable.
// This function is idempotent, meaning it can be called multiple times without reinitializing the configuration.
func Init() {
	// Check if already initialized
	if initialized {
		return
	}

	// Get the environment variable
	env := os.Getenv(TelegramModeEnvVar)
	if env == "" {
		env = "test"
	}
	fmt.Printf("Environment : %s\n", env)
	// Load the configuration based on the environment
	param, err := loadConfig(env)
	if err != nil {
		// Handle error if failed to load config
		_ = fmt.Errorf("failed to load config: %w", err)
	} else {

		fmt.Printf("Config : %+v\n", param)

		// Set the global parameter
		parameter = param

		// Mark as initialized
		initialized = true
	}
}

// GetParameter returns the initialized global parameter.
func GetParameter() *Parameter {
	Init()
	return parameter
}

// LoadConfigFromEnv loads the configuration parameters based on the TELE_MODE environment variable.
func LoadConfigFromEnv() (*Parameter, error) {
	env := os.Getenv(TelegramModeEnvVar)
	if env == "" {
		env = "test"
	}
	return loadConfig(env)
}

// loadConfig loads the configuration parameters based on the specified environment.
func loadConfig(env string) (*Parameter, error) {
	configFile := getConfigFileForEnv(env)
	fmt.Printf("Config File : %s\n", configFile)
	file, err := os.Open(configFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var config Parameter
	if _, err := toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}

// getConfigFileForEnv returns the appropriate configuration file for the given environment.
func getConfigFileForEnv(env string) string {
	switch env {
	case "prod":
		return ProdConfigFile
	case "test":
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
