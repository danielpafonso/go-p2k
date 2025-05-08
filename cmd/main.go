package main

import (
	"context"
	"flag"
	"fmt"
	"log"

	"cloud.google.com/go/pubsub"
	"github.com/IBM/sarama"

	"go-p2k/internal"
)

var (
	kafkaConn sarama.SyncProducer
)

func p2kMainLoop(sub *pubsub.Subscription, topic string) error {
	ctx, cancel := context.WithCancel(context.Background())
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(string(m.Data)),
		}
		_, _, err := kafkaConn.SendMessage(msg)
		if err != nil {
			cancel()
		}
		m.Ack()
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var configPath string

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file")
	flag.Parse()

	log.Println("Load configurations")
	configs, err := internal.LoadConfigurations(configPath)
	if err != nil {
		panic(err)
	}

	log.Println("Create Pub/Sub Client")
	pubsubClient, err := pubsub.NewClient(context.Background(), configs.Pubsub.Project)
	if err != nil {
		panic(err)
	}
	defer pubsubClient.Close()
	pubsubSubscriper := pubsubClient.Subscription(configs.Pubsub.Subscription)

	log.Println("Create Kafka Client")
	kafkaConfig := sarama.NewConfig()
	kafkaConfig.Producer.Return.Successes = true
	kafkaConn, err = sarama.NewSyncProducer(configs.Kafka.Endpoints, kafkaConfig)
	if err != nil {
		panic(err)
	}
	defer kafkaConn.Close()

	log.Println("Start main loop")
	err = p2kMainLoop(pubsubSubscriper, configs.Kafka.Topic)

	fmt.Println("End")
}
