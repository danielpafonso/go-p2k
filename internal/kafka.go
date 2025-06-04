package internal

import (
	"encoding/json"
	"errors"
)

type KafkaConfig struct {
	Topic   string
	Cluster string
}

// func ProcessMsg(data []byte) map[string]interface{} {
func ProcessMsg(data []byte) (*KafkaConfig, []byte, error) {
	msg := make(map[string]interface{})
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, make([]byte, 0), errors.New("message not json formated")
	}
	kCfg := KafkaConfig{
		Cluster: "all",
	}
	if kafka, ok := msg["_kafka"]; ok {
		kConfig := kafka.(map[string]interface{})

		// check topic
		if topic, ok := kConfig["topic"]; ok {
			kCfg.Topic = topic.(string)
		} else {
			return nil, make([]byte, 0), errors.New("no 'topic' field")
		}
		// check cluster
		if cluster, ok := kConfig["cluster"]; ok {
			kCfg.Cluster = cluster.(string)
		}

		// remove _kafka field
		delete(msg, "_kafka")

		bytesMsg, _ := json.Marshal(msg)

		return &kCfg, bytesMsg, nil
	} else {
		return nil, make([]byte, 0), errors.New("no '_kafka' field")
	}
}
