package main

import (
	"os"

	"github.com/urfave/cli/v2"
	"go.uber.org/zap"
	internalConfig "gzzn.com/airport/serial/internal/config"
	"gzzn.com/airport/serial/internal/serial"
)

func main() {
	logger, _ := zap.NewProduction()
	defer logger.Sync() // flushes buffer, if any
	sugar := logger.Sugar()
	app := &cli.App{
		Name:  "serial-read",
		Usage: "A serial port reading CLI application",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "config",
				Aliases: []string{"c"},
				Value:   "config.toml", // updated to utilize the TOML config file
				Usage:   "Path to the configuration file",
				EnvVars: []string{"CONFIG_PATH"},
			},
		},
		Action: func(c *cli.Context) error {
			configFile := c.String("config")
			mode, portName, err := internalConfig.LoadConfig(configFile)
			sugar.Infof("Mode: %v", mode)
			if err != nil {
				sugar.Fatalf("Error loading configuration: %v", err)
			}

			if err := serial.ReadFromPort(mode, portName); err != nil {
				sugar.Fatalf("Error reading from serial port: %v", err)
			}

			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		sugar.Fatal(err)
	}
}
