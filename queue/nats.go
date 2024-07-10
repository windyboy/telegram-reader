package queue

import (
	"crypto/rand"
	"time"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-nats/pkg/nats"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/nats-io/stan.go"
	"github.com/oklog/ulid"
	"gzzn.com/airport/serial/config"
	"gzzn.com/airport/serial/logger"
)

var (
	natsConfig config.NATSConfig
	publisher  *nats.StreamingPublisher
)

func Init() {
	if err := config.InitParameter(); err != nil {
		panic(err)
	}
	natsConfig = config.GetParameter().NATS
	publisherConfig := nats.StreamingPublisherConfig{
		ClusterID: natsConfig.ClusterId,
		ClientID:  natsConfig.ClientId,
		StanOptions: []stan.Option{
			stan.NatsURL(natsConfig.URL),
		},
		Marshaler: nats.GobMarshaler{},
	}
	var err error
	publisher, err = nats.NewStreamingPublisher(publisherConfig, watermill.NewStdLogger(false, false))
	if err != nil {
		panic(err)
	}

}

// Publish sends a message to the specified subject.
func Publish(payload []byte) error {
	sugar := logger.SugaredLogger()
	uuid := ulid.MustNew(ulid.Timestamp(time.Now()), rand.Reader)
	msg := message.NewMessage(uuid.String(), payload)
	if err := publisher.Publish(natsConfig.Subject, msg); err != nil {
		sugar.Errorf("Error publishing message: %v", err)
		return err
	}
	sugar.Debugf("Published message: %s", uuid.String())
	return nil
}
