/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

// UpdatePayload represents the possible payloads of a KafkaGatewayUpdate message
type UpdatePayload interface {
	// Apply applies the update specified by this payload to the store
	Apply(operation string, networkID string, offset int64, store storage.GatewayViewStorage) error
}

// Decoder is represents the object responsible for parsing Kafka messages to KafkaGatewayUpdates
type Decoder interface {
	// GetUpdateFromStreamAggregatorMessage parses a KafkaGatewayUpdate from Kafka message of the format seen
	// by the StreamAggregator
	GetUpdateFromStreamAggregatorMessage(message *kafka.Message) (*KafkaGatewayUpdate, error)
	// GetUpdateFromStateRecorderMessage parses a KafkaGatewayUpdate from a Kafka message of the format seen in
	// the StateRecorder
	GetUpdateFromStateRecorderMessage(message *kafka.Message) (*KafkaGatewayUpdate, error)
}

// Consumer is an interface implemented by the Kafka consumer class, created for dependency injection
type Consumer interface {
	Subscribe(topic string, rebalanceCb kafka.RebalanceCb) error
	SubscribeTopics(topics []string, rebalanceCb kafka.RebalanceCb) error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	Close() error
	Commit() ([]kafka.TopicPartition, error)
}

// Producer is an interface implemented by the Kafka producer class, created for dependency injection
type Producer interface {
	Produce(message *kafka.Message, deliveryChan chan kafka.Event) error
	Close()
}

type ConsumerFactory func() (Consumer, error)
type ProducerFactory func() (Producer, error)

// UpdateType is an enumerated type of possible classes of updates - configurations, statuses, or records
type UpdateType string

const (
	Configurations UpdateType = "configurations"
	Statuses       UpdateType = "gwstatus"
	Records        UpdateType = "gatewayrecords"
)
