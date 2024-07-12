package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli/v2"
	"go.bug.st/serial"
	"go.uber.org/zap"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
	nats "gzzn.com/airport/serial/queue"
	internalSerial "gzzn.com/airport/serial/serial"
	"gzzn.com/airport/serial/telegram"
)

var (
	parameter *config.Parameter

	totalBytes = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "serial",
		Name:      "serial_read_bytes_total",
		Help:      "The total number of bytes read from the serial port",
	})
	totalTelegrams = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "serial",
		Name:      "telegram_total",
		Help:      "The total number of telegrams received",
	})
)

func main() {
	app := setupApp()
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
	nats.Close()
}

func setupApp() *cli.App {
	app := &cli.App{
		Name:  "serial-read",
		Usage: "A serial port reading CLI application",
		Before: func(c *cli.Context) error {
			parameter = config.GetParameter()
			return nil
		},
		Commands: []*cli.Command{
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
				Action: executeReadCommand,
			},
			{
				Name:   "list",
				Usage:  "List available local serial ports",
				Action: executeListCommand,
			},
		},
	}
	return app
}

func executeReadCommand(c *cli.Context) error {
	sugar := logger.GetLogger()
	go startMetricsServer(sugar)

	dataChannel := make(chan []byte)
	go readFromPort(dataChannel, sugar)

	for data := range dataChannel {
		totalBytes.Add(float64(len(data)))
		processReceivedData(data)
	}
	return nil
}

func executeListCommand(c *cli.Context) error {
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

func startMetricsServer(log *zap.SugaredLogger) {
	addr := config.GetParameter().Prometheus.Address
	// log := logger.GetLogger()
	log.Infof("Starting metrics server on %s", addr)
	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(addr, nil))
}

func readFromPort(dataChannel chan<- []byte, log *zap.SugaredLogger) {
	// logger := logger.GetLogger()
	telegram.Init()
	mode, portName := config.ReadSerialConfig(parameter.Serial)
	if err := nats.Connect(parameter.NATS); err != nil {
		log.Fatalf("Error connecting to NATS: %v", err)
	}
	if err := internalSerial.ReadFromPort(mode, portName, parameter.Serial.BufferSize, dataChannel); err != nil {
		log.Fatalf("Error reading from port: %v", err)
	}
}

func processReceivedData(data []byte) {
	sugar := logger.GetLogger()
	for _, b := range data {
		if telegrams := telegram.Append(b); len(telegrams) > 0 {
			for _, telegramData := range telegrams {
				if sequence := telegram.GetSequence(telegramData); sequence != "" {
					sugar.Infof("Publishing telegram: %s", sequence)
				}
				if err := nats.Publish(telegramData); err != nil {
					sugar.Errorf("Error publishing to NATS: %v", err)
				}
			}
			totalTelegrams.Add(float64(len(telegrams)))
		}
	}
	sugar.Infof("Total bytes: %d, Total telegrams: %d", totalBytes, totalTelegrams)
}

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
