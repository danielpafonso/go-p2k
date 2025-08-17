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

func ProcessMsg(data []byte) (*KafkaConfig, []byte, error) {
	msg := make(map[string]any)
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, make([]byte, 0), errors.New("message not json formated")
	}
	kCfg := KafkaConfig{
		Clusters: []string{"all"},
	}
	if kafka, ok := msg["_kafka"]; ok {
		kConfig := kafka.(map[string]any)

		// check topic
		if topic, ok := kConfig["topic"]; ok {
			if stopic, ok := topic.(string); ok {
				kCfg.Topic = stopic
			} else {
				return nil, make([]byte, 0), errors.New("'topic' field isn't a string")
			}
		} else {
			return nil, make([]byte, 0), errors.New("no 'topic' field")
		}
		// check clusters
		if clusters, ok := kConfig["clusters"]; ok {
			if sclusters, ok := clusters.(string); ok {
				kCfg.Clusters = strings.Split(strings.ReplaceAll(sclusters, ", ", ","), ",")
			} else {
				return nil, make([]byte, 0), errors.New("'clusters' fields isn't a string")
			}
		}

		// remove _kafka field
		delete(msg, "_kafka")

		bytesMsg, _ := json.Marshal(msg)

		return &kCfg, bytesMsg, nil
	} else {
		return nil, make([]byte, 0), errors.New("no '_kafka' field")
	}
}
