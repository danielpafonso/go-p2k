package main

import (
	"fmt"
	"strings"
)

type Metric struct {
	Name      string
	Labels    map[string]string
	Value     float64
	Timestamp int64
}

type Metrics struct {
	LastPubsub Metric
	LastKafka  Metric
	Updates    int
}

func InitiateMetrics() *Metrics {
	return &Metrics{
		LastPubsub: Metric{
			Name: "last_pubsub_message",
		},
		LastKafka: Metric{
			Name:   "last_kafka_message",
			Labels: map[string]string{"service": "kafka"},
		},
	}
}

func (mtc *Metric) Text() string {
	builder := strings.Builder{}

	// Name
	builder.WriteString(mtc.Name)

	// Compile labels
	if len(mtc.Labels) != 0 {
		labels := make([]string, 0)
		for key, value := range mtc.Labels {
			labels = append(labels, key+"=\""+value+"\"")
		}
		builder.WriteString("{")
		builder.WriteString(strings.Join(labels, ","))
		builder.WriteString("}")
	}

	// Value
	builder.WriteString(fmt.Sprintf(" %f", mtc.Value))

	// Timestamp
	if mtc.Timestamp != 0 {
		builder.WriteString(fmt.Sprintf(" %d", mtc.Timestamp))
	}
	return builder.String()
}

func (mtc *Metrics) Print() string {
	return fmt.Sprintf(
		`# go-P2K metrics endpoint
# update version: %d

# Last timestamps
%s

%s`,
		mtc.Updates,
		mtc.LastPubsub.Text(),
		mtc.LastKafka.Text(),
	)
}
