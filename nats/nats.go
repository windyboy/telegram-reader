package nats

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"gzzn.com/airport/serial/logger"
)

var (
	nc       *nats.Conn
	initOnce sync.Once
	mu       sync.Mutex
)

// InitNATS initializes the NATS connection with the provided URL.
func InitNATS(url string) error {
	sugar := logger.SugaredLogger()

	var err error
	initOnce.Do(func() {
		sugar.Infof("Connecting to NATS: %s", url)
		nc, err = nats.Connect(url)
		if err != nil {
			sugar.Errorf("Error connecting to NATS: %v", err)
		} else {
			sugar.Infof("Connected to NATS: %s", url)
		}
	})

	return err
}

// Publish sends a message to the specified subject.
func Publish(subject, msg string) error {
	sugar := logger.SugaredLogger()
	mu.Lock()
	defer mu.Unlock()

	if nc == nil {
		sugar.Errorf("NATS connection is not initialized")
		return fmt.Errorf("NATS connection is not initialized")
	}

	sugar.Debugf("Publishing message to subject: %s", subject)
	return nc.Publish(subject, []byte(msg))
}

// Close terminates the NATS connection if it is initialized.
func Close() {
	sugar := logger.SugaredLogger()
	mu.Lock()
	defer mu.Unlock()

	if nc != nil {
		sugar.Infof("Closing NATS connection")
		nc.Close()
		nc = nil
	} else {
		sugar.Infof("NATS connection is already closed or was never initialized")
	}
}
