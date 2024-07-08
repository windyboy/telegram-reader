package config

import (
	"os"

	"github.com/BurntSushi/toml"
	serial "go.bug.st/serial"
)

const FILE_NAME = "../config.toml"

type Parameter struct {
	Serial   SerialConfig
	NATS     NATSConfig
	Telegram TelegramConfig
	Logger   LoggerConfig
}
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

type NATSConfig struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
	Subject  string `toml:"subject"`
}

type TelegramConfig struct {
	EndTag string `toml:"end_tag"`
	SeqTag string `toml:"seq_tag"`
}

type LoggerConfig struct {
	Level      string `toml:"level"`
	Filename   string `toml:"filename"`
	MaxSize    int    `toml:"max_size"`
	MaxBackups int    `toml:"max_backups"`
	MaxAge     int    `toml:"max_age"`
	Compress   bool   `toml:"compress"`
}

func LoadConfig(filename string) (Parameter, error) {
	if filename == "" {
		filename = FILE_NAME
	}
	// sugar := logger.SugaredLogger()
	// sugar.Infof("Loading config from %s", filename)
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
		// return Parameter{}, err
	}
	defer file.Close()

	var config Parameter
	if _, err := toml.NewDecoder(file).Decode(&config); err != nil {
		panic(err)
	}
	// sugar.Infof("Loaded config: %v", config)
	return config, nil
}

func ReadSerialConfig(serrialConfig SerialConfig) (*serial.Mode, string) {
	mode := &serial.Mode{
		BaudRate: serrialConfig.Baud,
		DataBits: serrialConfig.Size,
		Parity:   parseParity(serrialConfig.Parity),
		StopBits: parseStopBits(serrialConfig.StopBits),
		// ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Millisecond,
	}
	return mode, serrialConfig.Name
}

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
