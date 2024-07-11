package queue

import (
	"fmt"
	"sync"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

var (
	nc         *nats.Conn
	initOnce   sync.Once
	mu         sync.Mutex
	natsConfig config.NATSConfig
	sugar      *zap.SugaredLogger
)

// InitNATS initializes the NATS connection with the provided URL.
func Init() error {
	var err error
	initOnce.Do(func() {
		if param, err := config.GetParameter(); err != nil {
			panic(err)
		} else {
			natsConfig = param.NATS
		}
		opts := []nats.Option{nats.Name("Serial Telegram Publisher")}
		opts = append(opts, nats.UserInfo(natsConfig.Username, natsConfig.Password))
		urls := natsConfig.URLS
		sugar = logger.SugaredLogger()
		nc, err = nats.Connect(urls, opts...)
		if err != nil {
			sugar.Errorf("Error connecting to NATS: %v", err)
		} else {
			sugar.Infof("Connected to NATS: %s", urls)
		}
	})
	return err
}

// Publish sends a message to the specified subject.
func Publish(msg []byte) error {
	mu.Lock()
	defer mu.Unlock()

	if nc == nil {
		err := fmt.Errorf("NATS connection is not initialized")
		sugar.Errorf(err.Error())
		return err
	}
	subject := natsConfig.Subject
	sugar.Debugf("Publishing message to subject: %s", subject)
	return nc.Publish(subject, []byte(msg))
}

// Close terminates the NATS connection if it is initialized.
func Close() {
	if nc != nil {
		nc.Close()
		sugar.Info("NATS connection closed")
	}
}
