package internal

import (
	"encoding/json"
	"os"
	"strconv"
	"strings"
)

type PubsubConfigurations struct {
	Project      string `json:"project"`
	Subscription string `json:"subscription"`
}

type KafkaConfigurations struct {
	Endpoints []string `json:"endpoints"`
	UseTLS    bool     `json:"useTls"`
	CaFile    string   `json:"caFile"`
	CrtFile   string   `json:"crtFile"`
	KeyFile   string   `json:"keyFile"`
}

type Configurations struct {
	Pubsub PubsubConfigurations `json:"pubsub"`
	Kafka  KafkaConfigurations  `json:"kafka"`
}

func LoadConfigurations(filepath string) (*Configurations, error) {

	var configs Configurations

	// read configurations file if exists
	if _, err := os.Stat(filepath); err == nil {
		fdata, err := os.ReadFile(filepath)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(fdata, &configs)
		if err != nil {
			return nil, err
		}
	}

	// overwrite environment variables
	if value, ok := os.LookupEnv("PUBSUB_PROJECT"); ok {
		configs.Pubsub.Project = value
	}
	if value, ok := os.LookupEnv("PUBSUB_SUBSCRIPTION"); ok {
		configs.Pubsub.Subscription = value
	}
	if value, ok := os.LookupEnv("KAFKA_ENDPOINTS"); ok {
		configs.Kafka.Endpoints = strings.Split(value, ",")
	}
	if value, ok := os.LookupEnv("KAFKA_USE_TLS"); ok {
		configs.Kafka.UseTLS, _ = strconv.ParseBool(value)
	}
	if value, ok := os.LookupEnv("KAFKA_CA_FILE"); ok {
		configs.Kafka.CaFile = value
	}
	if value, ok := os.LookupEnv("KAFKA_CRT_FILE"); ok {
		configs.Kafka.CrtFile = value
	}
	if value, ok := os.LookupEnv("KAFKA_KEY_FILE"); ok {
		configs.Kafka.KeyFile = value
	}

	return &configs, nil
}
