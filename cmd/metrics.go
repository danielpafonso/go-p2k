package main

import (
	"fmt"
	"strings"
	"time"
)

type Metric struct {
	Name      string
	Labels    map[string]string
	Value     float64
	Timestamp int64
}

type TimeWindowMetric struct {
	Name      string
	Labels    map[string]string
	Value     []int64
	Timestamp int64
}

type Metrics struct {
	LastPubsub Metric
	LastKafka  Metric
	ValidMsg   TimeWindowMetric
}

func (twm *TimeWindowMetric) AddTime() {
	now := time.Now().UTC().UnixMilli()
	twm.Timestamp = now
	twm.Value = append(twm.Value, now)
	limit := now - 60000 // 1m
	for len(twm.Value) > 0 {
		if twm.Value[0] < limit {
			twm.Value = twm.Value[1:]
		} else {
			break
		}
	}
}

func (twm *TimeWindowMetric) Text() string {
	builder := strings.Builder{}

	// Name
	builder.WriteString(twm.Name)

	// Compile labels
	if len(twm.Labels) != 0 {
		labels := make([]string, 0)
		for key, value := range twm.Labels {
			labels = append(labels, key+"=\""+value+"\"")
		}
		builder.WriteString("{")
		builder.WriteString(strings.Join(labels, ","))
		builder.WriteString("}")
	}

	// Value
	limit := time.Now().Add(-1 * time.Minute).UTC().UnixMilli()
	for len(twm.Value) > 0 {
		if twm.Value[0] < limit {
			twm.Value = twm.Value[1:]
		} else {
			break
		}
	}
	builder.WriteString(fmt.Sprintf(" %d", len(twm.Value)))

	// Timestamp
	if twm.Timestamp != 0 {
		builder.WriteString(fmt.Sprintf(" %d", twm.Timestamp))
	}
	return builder.String()
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
		ValidMsg: TimeWindowMetric{
			Name:   "num_valid_msgs",
			Labels: map[string]string{"window": "1min"},
			Value:  make([]int64, 0),
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

# Last timestamps
%s

%s

# Time Window metrics
%s
`,
		mtc.LastPubsub.Text(),
		mtc.LastKafka.Text(),
		mtc.ValidMsg.Text(),
	)
}
