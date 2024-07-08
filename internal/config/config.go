package config

import (
	"os"

	"github.com/BurntSushi/toml"
	serial "go.bug.st/serial"
	"go.uber.org/zap"
)

type Config struct {
	Name        string `toml:"name"`
	Baud        int    `toml:"baud"`
	ReadTimeout int    `toml:"read_timeout"`
	Size        int    `toml:"size"`
	Parity      string `toml:"parity"`
	StopBits    int    `toml:"stop_bits"`
	FlowControl string `toml:"flow_control"`
}

func LoadConfig(filename string) (*serial.Mode, string, error) {
	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()
	file, err := os.Open(filename)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	var cfg Config
	if _, err := toml.NewDecoder(file).Decode(&cfg); err != nil {
		return nil, "", err
	}
	sugar.Infof("Loaded config: %v", cfg)

	mode := &serial.Mode{
		BaudRate: cfg.Baud,
		DataBits: cfg.Size,
		Parity:   parseParity(cfg.Parity),
		StopBits: parseStopBits(cfg.StopBits),
		// ReadTimeout: time.Duration(cfg.ReadTimeout) * time.Millisecond,
	}

	return mode, cfg.Name, nil
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
