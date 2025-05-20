package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/IBM/sarama"

	"go-p2k/internal"
)

var (
	kafkaConn sarama.SyncProducer
)

func p2kMainLoop(sub *pubsub.Subscription, topic string, metrics *Metrics) error {
	ctx, cancel := context.WithCancel(context.Background())
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		now := time.Now().UnixMilli()
		metrics.LastPubsub.Value = float64(now)
		metrics.LastPubsub.Timestamp = now
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(string(m.Data)),
		}
		_, _, err := kafkaConn.SendMessage(msg)
		if err != nil {
			cancel()
		}
		metrics.LastKafka.Value = float64(time.Now().UnixMilli())
		metrics.ValidMsg.AddTime()
		m.Ack()
	})
	if err != nil {
		return err
	}
	return nil
}

func main() {
	var configPath string
	var serverAddr string
	var serverPort string

	flag.StringVar(&configPath, "c", "config.json", "Path to configuration file")
	flag.StringVar(&serverAddr, "a", "localhost", "Server Address")
	flag.StringVar(&serverPort, "p", "8080", "Server Port")
	flag.Parse()

	log.Println("Load configurations")
	configs, err := internal.LoadConfigurations(configPath)
	if err != nil {
		panic(err)
	}

	// runtime variables
	running := false
	metrics := InitiateMetrics()

	// http endpoints
	rootHandler := http.NewServeMux()
	rootHandler.HandleFunc("/", generalHandler)
	rootHandler.HandleFunc("/health", healthHandler(&running))
	rootHandler.HandleFunc("/metrics", metricsHandler(metrics))

	server := http.Server{
		Addr:    serverAddr + ":" + serverPort,
		Handler: rootHandler,
	}
	log.Printf("Start server at %s\n", server.Addr)
	// launch server in a go routine
	go server.ListenAndServe()

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

	time.Sleep(5 * time.Second)

	running = true

	log.Println("Start main loop")
	err = p2kMainLoop(pubsubSubscriper, configs.Kafka.Topic, metrics)

	log.Println("End")
}
