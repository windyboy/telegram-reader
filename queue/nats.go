package nats

import (
	"sync"

	"github.com/nats-io/nats.go"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

var (
	nc *nats.Conn
	mu sync.Mutex
)

func Connect(config config.NATSConfig) error {
	sugar := logger.GetLogger()
	sugar.Infof("Connecting to NATS: %+v", config)
	opts := []nats.Option{nats.Name("Serial Telegram Publisher")}
	opts = append(opts, nats.UserInfo(config.Username, config.Password))
	url := config.URLS
	var err error
	if nc, err = nats.Connect(url, opts...); err != nil {
		sugar.Errorf("Error connecting to NATS: %v", err)
		return err
	}
	sugar.Infof("Connected to NATS: %s", url)
	return nil
}

func Publish(msg string) error {
	sugar := logger.GetLogger()
	mu.Lock()
	defer mu.Unlock()

	if nc == nil {
		sugar.Fatalf("NATS connection is not initialized")
	}
	subject := config.GetParameter().NATS.Subject
	return nc.Publish(subject, []byte(msg))
}

func Close() {
	sugar := logger.GetLogger()
	mu.Lock()
	defer mu.Unlock()

	if nc != nil {
		sugar.Infof("Closing NATS connection")
		nc.Close()
		nc = nil
	}
}
