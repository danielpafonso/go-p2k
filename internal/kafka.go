package internal

import (
	"encoding/json"
	"errors"
	"strings"
)

type KafkaConfig struct {
	Topic    string
	Clusters []string
}

// func ProcessMsg(data []byte) map[string]interface{} {
func ProcessMsg(data []byte) (*KafkaConfig, []byte, error) {
	msg := make(map[string]interface{})
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, make([]byte, 0), errors.New("message not json formated")
	}
	kCfg := KafkaConfig{
		Clusters: []string{"all"},
	}
	if kafka, ok := msg["_kafka"]; ok {
		kConfig := kafka.(map[string]interface{})

		// check topic
		if topic, ok := kConfig["topic"]; ok {
			kCfg.Topic = topic.(string)
		} else {
			return nil, make([]byte, 0), errors.New("no 'topic' field")
		}
		// check clusters
		if clusters, ok := kConfig["clusters"]; ok {
			kCfg.Clusters = strings.Split(strings.ReplaceAll(clusters.(string), ", ", ","), ",")
		}

		// remove _kafka field
		delete(msg, "_kafka")

		bytesMsg, _ := json.Marshal(msg)

		return &kCfg, bytesMsg, nil
	} else {
		return nil, make([]byte, 0), errors.New("no '_kafka' field")
	}
}
