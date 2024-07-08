package nats

import (
	"fmt"

	"github.com/nats-io/nats.go"
	"gzzn.com/airport/serial/logger"
)

var Nc *nats.Conn

func InitNATS(url string) error {
	sugar := logger.SugaredLogger()
	var err error
	sugar.Infof("Connecting to NATS: %s", url)
	Nc, err = nats.Connect(url)
	if err != nil {
		sugar.Errorf("Error connecting to NATS: %v", err)
		return err
	}
	return nil
}

func Publish(subject, msg string) error {
	sugar := logger.SugaredLogger()
	if Nc == nil {
		sugar.Errorf("NATS connection is not initialized")
		return fmt.Errorf("NATS connection is not initialized")
	}
	sugar.Debugf("Publishing message to subject: %s", subject)
	return Nc.Publish(subject, []byte(msg))
}

func Close() {
	sugar := logger.SugaredLogger()
	sugar.Infof("Closing NATS connection")
	if Nc != nil {
		Nc.Close()
	}
}
