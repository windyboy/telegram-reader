package nats

import (
	"sync"

	"github.com/nats-io/nats.go"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

var (
	nc         *nats.Conn
	initOnce   sync.Once
	mu         sync.Mutex
	natsConfig config.NATSConfig
)

// InitNATS initializes the NATS connection with the provided URL.
// func InitNATS() {

// 	initOnce.Do(func() {
// 		sugar := logger.GetLogger()
// 		natsConfig = config.GetParameter().NATS
// 		logger.GetLogger().Infof("Connecting to NATS: %+v", natsConfig)
// 		opts := []nats.Option{nats.Name("Serial Telegram Publisher")}
// 		opts = append(opts, nats.UserInfo(natsConfig.Username, natsConfig.Password))
// 		url := natsConfig.URLS
// 		var err error
// 		if nc, err = nats.Connect(url, opts...); err != nil {
// 			sugar.Errorf("Error connecting to NATS: %v", err)
// 		} else {
// 			defer Close()
// 			sugar.Infof("Connected to NATS: %s", url)
// 		}
// 	})
// }

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
	// defer Close()
	sugar.Infof("Connected to NATS: %s", url)
	return nil
}

// Publish sends a message to the specified subject.
func Publish(msg string) error {
	sugar := logger.GetLogger()
	mu.Lock()
	defer mu.Unlock()

	if nc == nil {
		sugar.Fatalf("NATS connection is not initialized")
	}
	subject := config.GetParameter().NATS.Subject
	// sugar.Debugf("Publishing message to subject: %s", subject)
	return nc.Publish(subject, []byte(msg))
}

// Close terminates the NATS connection if it is initialized.
func Close() {
	sugar := logger.GetLogger()
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
