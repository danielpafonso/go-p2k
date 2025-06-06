package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"cloud.google.com/go/pubsub"
	"github.com/IBM/sarama"

	"go-p2k/internal"
)

var (
	kafkaConn    sarama.SyncProducer
	kafkaClients map[string]sarama.SyncProducer
)

func p2kMainLoop(sub *pubsub.Subscription, metrics *Metrics) error {
	ctx, cancel := context.WithCancel(context.Background())
	err := sub.Receive(ctx, func(ctx context.Context, m *pubsub.Message) {
		now := time.Now().UnixMilli()
		metrics.LastPubsub.Value = float64(now)
		metrics.LastPubsub.Timestamp = now
		// process message
		kafkaConfig, bytesMsg, err := internal.ProcessMsg(m.Data)
		if err != nil {
			log.Println("error on processing message:", err)
			log.Printf("DEADLETTER %s\n", string(m.Data))
		} else {
			if kafkaConfig.Topic == "" {
				fmt.Println("no topic in _kafka field")
			} else {
				// check clusters
				for _, cluster := range kafkaConfig.Clusters {
					if conn, ok := kafkaClients[cluster]; ok {
						msg := &sarama.ProducerMessage{
							Topic: kafkaConfig.Topic,
							Value: sarama.StringEncoder(string(bytesMsg)),
						}
						_, _, err = conn.SendMessage(msg)
						if err != nil {
							cancel()
						}
						metrics.LastKafka.Value = float64(time.Now().UnixMilli())
						metrics.ValidMsg.AddTime()
					} else if cluster == "all" {
						for _, conn := range kafkaClients {
							msg := &sarama.ProducerMessage{
								Topic: kafkaConfig.Topic,
								Value: sarama.StringEncoder(string(bytesMsg)),
							}
							_, _, err = conn.SendMessage(msg)
							if err != nil {
								cancel()
							}
						}
					} else {
						fmt.Printf("no configured cluster: %s\n", cluster)
						log.Printf("DEADLETTER %s\n", string(m.Data))
					}
				}
			}
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

	log.Println("Create Kafka Clients")
	// Create Kafka Clients
	kafkaClients = make(map[string]sarama.SyncProducer)

	for _, config := range configs.Kafka.Clusters {
		cfg := sarama.NewConfig()
		cfg.Producer.Return.Successes = true
		connection, err := sarama.NewSyncProducer(config.Endpoints, cfg)
		if err != nil {
			panic(err)
		}
		kafkaClients[config.Name] = connection

		// defer instead of creating a "defer" function
		defer connection.Close()
	}

	running = true

	log.Println("Start main loop")
	err = p2kMainLoop(pubsubSubscriper, metrics)

	log.Println("End")
}
