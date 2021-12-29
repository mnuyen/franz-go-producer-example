package app

import (
	"context"
	"crypto/tls"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/twmb/franz-go/pkg/kgo"
	"github.com/twmb/franz-go/pkg/sasl/plain"
)

var (
	saslUser = flag.String("sasl-user", "$ConnectionString", "if non-empty, username to use for sasl (must specify all options)")
)

// local address of your kafka broker
var seeds = []string{"localhost:29092"}

// if you use an eventhub. paste here your connection string in
var saslPass = "someConnectionString"

var FanzProducerClient *kgo.Client

func StartProducer() {

	opts := []kgo.Opt{
		kgo.SeedBrokers(seeds...),
		kgo.RequiredAcks(kgo.AllISRAcks()),
		kgo.AllowAutoTopicCreation(),
	}

	eventhubConf, err := strconv.ParseBool(os.Getenv("EVENTHUB"))

	if err != nil {
		eventhubConf = false
	}

	if eventhubConf {
		tlsDialer := &tls.Dialer{NetDialer: &net.Dialer{Timeout: 10 * time.Second}}
		opts = append(opts,
			kgo.SASL(plain.Auth{User: *saslUser, Pass: saslPass}.AsMechanism()),
			// Configure TLS. Uses SystemCertPool for RootCAs by default.
			kgo.Dialer(tlsDialer.DialContext))
	}

	cl, err := kgo.NewClient(opts...)

	if err != nil {
		logrus.Error("Client could not be created cause", err)
		panic(err)
	}

	FanzProducerClient = cl
}

func SendKafkaMessage(topic string, trackingId string, message string) {
	record := &kgo.Record{
		Topic: topic,
		Value: []byte(message),
	}
	logrus.Info(fmt.Sprintf("Sending Kafka Messages with Franz Client to Kakfa on topic: [%v]", topic))

	FanzProducerClient.Produce(context.Background(), record, func(_ *kgo.Record, err error) {
		if err != nil {
			logrus.Error(fmt.Sprintf("[TrackindId %s] Failed to produce message %s with error:\n %v\n", trackingId, message, err))
		}
	})
}
