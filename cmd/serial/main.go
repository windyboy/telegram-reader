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
	"gzzn.com/airport/serial/internal/config"
	"gzzn.com/airport/serial/internal/logger"
	nats "gzzn.com/airport/serial/internal/queue"
	internalSerial "gzzn.com/airport/serial/internal/serial"
	"gzzn.com/airport/serial/internal/telegram"
)

const (
	NatsUrl     = "nats-url"
	NatsSubject = "nats-subj"
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
						Name:    NatsUrl,
						Aliases: []string{"n"},
						Value:   "nats://localhost:4222",
						Usage:   "NATS server URL",
						EnvVars: []string{"NATS_URL"},
					},
					&cli.StringFlag{
						Name:    NatsSubject,
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
	parameter = config.GetParameter()
	sugar := logger.GetLogger()
	overwriteParameter(c)
	go startMetricsServer(parameter, sugar)

	dataChannel := make(chan []byte)
	go readFromPort(dataChannel, parameter, sugar)

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

func startMetricsServer(parameter *config.Parameter, log *zap.SugaredLogger) {
	addr := parameter.Prometheus.Address
	// log := logger.GetLogger()
	log.Infof("Starting metrics server on %s", addr)
	http.Handle("/metrics", promhttp.Handler())
	panic(http.ListenAndServe(addr, nil))
}

func readFromPort(dataChannel chan<- []byte, parameter *config.Parameter, log *zap.SugaredLogger) {
	// telegram.Init(parameter)
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

func overwriteParameter(c *cli.Context) {
	parameter = config.GetParameter()
	if c.IsSet(NatsUrl) {
		parameter.NATS.URLS = c.String(NatsUrl)
	}

	if c.IsSet(NatsSubject) {
		parameter.NATS.Subject = c.String(NatsSubject)
	}
}