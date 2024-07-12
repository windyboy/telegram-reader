package main

import (
	"fmt"
	"os"

	"github.com/urfave/cli/v2"
	"go.bug.st/serial"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
	nats "gzzn.com/airport/serial/queue"
	internalSerial "gzzn.com/airport/serial/serial"
	"gzzn.com/airport/serial/telegram"
)

var (
	parameter *config.Parameter
	// sugar     *zap.SugaredLogger
)

// listAvailablePorts returns a list of available serial ports
func listAvailablePorts() ([]string, error) {
	enumerator, err := serial.GetPortsList()
	if err != nil {
		return nil, err
	}
	if enumerator == nil {
		return nil, fmt.Errorf("no serial ports enumerator available")
	}
	return enumerator, nil
}

// setupApp initializes the CLI application and its commands
func setupApp() *cli.App {
	app := &cli.App{
		Name:  "serial-read",
		Usage: "A serial port reading CLI application",
		Before: func(c *cli.Context) error {
			parameter = config.GetParameter()
			// sugar = logger.SugaredLogger()
			return nil
		},
		Commands: []*cli.Command{
			// Add "read" command configuration
			{
				Name:  "read",
				Usage: "Read data from a serial port",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "config",
						Aliases: []string{"c"},
						Value:   "config.toml",
						Usage:   "Path to the configuration file",
						EnvVars: []string{"CONFIG_PATH"},
					},
					&cli.StringFlag{
						Name:    "nats-url",
						Aliases: []string{"n"},
						Value:   "nats://localhost:4222",
						Usage:   "NATS server URL",
						EnvVars: []string{"NATS_URL"},
					},
					&cli.StringFlag{
						Name:    "nats-subj",
						Aliases: []string{"s"},
						Value:   "serial.data",
						Usage:   "NATS subject to publish data to",
						EnvVars: []string{"NATS_SUBJECT"},
					},
				},
				Action: func(c *cli.Context) error {
					return executeReadCommand()
				},
			},
			// Add "list" command configuration
			{
				Name:  "list",
				Usage: "List available local serial ports",
				Action: func(c *cli.Context) error {
					return executeListCommand()
				},
			},
		},
	}
	return app
}

// executeReadCommand handles the logic for the "read" command
func executeReadCommand() error {
	sugar := logger.GetLogger()
	telegram.Init()
	// nats.InitNATS()
	// defer nats.Close()
	mode, portName := config.ReadSerialConfig(parameter.Serial)
	// sugar.Infof("Opening port: %s with mode: %+v", portName, mode)

	dataChannel := make(chan []byte)
	if err := nats.Connect(parameter.NATS); err != nil {
		sugar.Fatalf("Error connecting to NATS: %v", err)
	}

	go func() {
		if err := internalSerial.ReadFromPort(mode, portName, parameter.Serial.BufferSize, dataChannel); err != nil {
			sugar.Fatalf("Error reading from port: %v", err)
		}
	}()

	for data := range dataChannel {
		processReceivedData(data)
	}
	return nil
}

// executeListCommand handles the logic for the "list" command
func executeListCommand() error {
	sugar := logger.GetLogger()
	sugar.Infof("Listing available serial ports")
	ports, err := listAvailablePorts()
	if err != nil {
		sugar.Fatalf("Error listing serial ports: %v", err)
	}
	for _, port := range ports {
		fmt.Println(port)
	}
	return nil
}

// processReceivedData processes received data from the serial port and publishes to NATS
func processReceivedData(data []byte) {
	sugar := logger.GetLogger()
	if telegrams := telegram.Append(string(data)); len(telegrams) > 0 {
		// sugar.Debugf("Got %d telegrams to publish", len(telegrams))
		for _, telegramData := range telegrams {
			if sequence := telegram.GetTelegramSequence(telegramData); sequence != "" {
				sugar.Infof("Publishing telegram: %s", sequence)
			}
			if err := nats.Publish(telegramData); err != nil {
				sugar.Errorf("Error publishing to NATS: %v", err)
			}
		}
	}
}

func main() {
	app := setupApp()
	if err := app.Run(os.Args); err != nil {
		panic(err)
	}
	nats.Close()
}
