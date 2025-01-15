package config

import (
	"fmt"
	"os"
	"sync"

	"github.com/BurntSushi/toml"
	serial "go.bug.st/serial"
	// "gzzn.com/airport/serial/internal/logger"
)

const (
	ProdConfigFile     = "config.toml"
	TestConfigFile     = "config.test.toml"
	DevConfigFile      = "config.dev.toml"
	TelegramModeEnvVar = "TELE_MODE"
)

var (
	parameter *Parameter
	once      sync.Once
)

// Parameter holds the configuration for the application.
type Parameter struct {
	Serial     SerialConfig
	NATS       NATSConfig
	Prometheus PrometheusConfig
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

// PrometheusConfig holds the Prometheus configuration.
type PrometheusConfig struct {
	Address string `toml:"address"`
}

// getEnvironment retrieves the TELE_MODE environment variable or defaults to "test".
func getEnvironment() string {
	env := os.Getenv(TelegramModeEnvVar)
	if env == "" {
		env = "dev"
	}
	return env
}

// load initializes the configuration by loading the appropriate config file.
func load() error {
	var err error
	once.Do(func() {
		env := getEnvironment()
		parameter, err = loadConfig(env)
	})
	return err
}

// GetParameter returns the initialized global parameter.
// It loads the configuration if it hasn't been loaded yet.
func GetParameter() *Parameter {
	if parameter == nil {
		if err := load(); err != nil {
			fmt.Printf("Failed to load configuration: %v\n", err)
			return nil
		}
	}
	return parameter
}

// loadConfig loads the configuration based on the provided environment.
// It returns a Parameter and an error if loading fails.
func loadConfig(env string) (*Parameter, error) {
	configFile := getConfigFileForEnv(env)
	// currentDir, err := os.Getwd()
	// if err != nil {
	// 	return nil, fmt.Errorf("error getting current directory: %w", err)
	// }
	// configPath := filepath.Join(currentDir, "internal/config", configFile)

	// file, err := os.Open(configPath)
	file, err := os.Open(configFile)
	fmt.Println("configFile: ", configFile)
	if err != nil {
		return nil, fmt.Errorf("error opening config file: %w", err)
	}
	defer file.Close()

	var config Parameter
	if _, err := toml.NewDecoder(file).Decode(&config); err != nil {
		return nil, fmt.Errorf("error decoding config file: %w", err)
	}

	// Set default values if necessary
	setDefaultValues(&config)

	return &config, nil
}

// setDefaultValues sets default values for the configuration parameters.
func setDefaultValues(config *Parameter) {
	if config.Serial.Name == "" {
		config.Serial.Name = "COM1" // Default value
	}
	if config.Serial.Baud == 0 {
		config.Serial.Baud = 9600 // Default value
	}
	// Set other defaults as needed
}

// getConfigFileForEnv returns the appropriate configuration file for the given environment.
func getConfigFileForEnv(env string) string {
	fmt.Println("env: ", env)
	switch env {
	case "prod":
		return ProdConfigFile
	case "test":
		return TestConfigFile
	default:
		return DevConfigFile
	}
}

// ReadSerialConfig returns the serial mode and port name from the SerialConfig.
func ReadSerialConfig(serialConfig SerialConfig) (mode *serial.Mode, portName string) {
	portName = serialConfig.Name

	// Convert string parity to serial.Parity
	var parity serial.Parity
	switch serialConfig.Parity {
	case "none":
		parity = serial.NoParity
	case "odd":
		parity = serial.OddParity
	case "even":
		parity = serial.EvenParity
	default:
		parity = serial.NoParity
	}

	mode = &serial.Mode{
		BaudRate: serialConfig.Baud,
		DataBits: serialConfig.Size,
		Parity:   parity,
		StopBits: serial.StopBits(serialConfig.StopBits),
		// Add flow control if needed
		// FlowControl: serialConfig.FlowControl,
	}
	return
}
