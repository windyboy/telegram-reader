package queue

import (
	"fmt"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
	"gzzn.com/airport/serial/internal/config"
	"gzzn.com/airport/serial/internal/logger"
)

var (
	nc *nats.Conn
	mu sync.Mutex
)

// Connect establishes a connection to the NATS server with retry logic.
func Connect(cfg config.NATSConfig) error {
	sugar := logger.GetLogger()
	sugar.Infof("Connecting to NATS: %+v", cfg)

	opts := []nats.Option{
		nats.Name("Serial Telegram Publisher"),
		nats.UserInfo(cfg.Username, cfg.Password),
		nats.MaxReconnects(5),               // Set maximum reconnect attempts
		nats.ReconnectWait(2 * time.Second), // Wait time between reconnect attempts
	}

	var err error
	nc, err = nats.Connect(cfg.URLS, opts...)
	if err != nil {
		sugar.Errorf("Error connecting to NATS: %v", err)
		return err
	}

	sugar.Infof("Connected to NATS: %s", cfg.URLS)
	return nil
}

// Publish sends a message to the specified NATS subject.
func Publish(msg string) error {
	sugar := logger.GetLogger()
	mu.Lock()
	defer mu.Unlock()
	if nc == nil {
		err := fmt.Errorf("NATS connection is not initialized")
		sugar.Error(err)
		return err
	}

	subject := config.GetParameter().NATS.Subject
	if err := nc.Publish(subject, []byte(msg)); err != nil {
		sugar.Errorf("Error publishing message to NATS: %v", err)
		return err
	}

	sugar.Infof("Published message to NATS subject '%s'", subject)
	return nil
}

// Close gracefully closes the NATS connection.
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
