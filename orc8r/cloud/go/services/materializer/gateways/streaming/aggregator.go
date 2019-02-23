/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"encoding/json"
	"fmt"
	"os"
	"sync/atomic"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

const (
	consumerTopicsRegex = "^magma\\.public\\.(.+)_(configurations|gwstatus|gatewayrecords)$"
	aggregatedTopic     = "magma.public.gatewaystates"
)

// StreamAggregator is what listens to all topics with information pertaining to gateway state and aggregates
// them into a single topic of gateway state updates
type StreamAggregator struct {
	consumerFactory ConsumerFactory
	decoder         Decoder
	producerFactory ProducerFactory
	stopFlag        *atomic.Value
}

func NewStreamAggregatorConsumer() (Consumer, error) {
	clusterServers, ok := os.LookupEnv("CLUSTER_SERVERS")
	if !ok {
		return nil, fmt.Errorf("CLUSTER_SERVERS was not defined")
	}

	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   clusterServers,
		"group.id":            "gateway-state-stream-aggregator-consumers",
		"auto.offset.reset":   "earliest",
		"enable.auto.commit":  false,
		"metadata.max.age.ms": 30000,
	})
}

func NewStreamAggregatorProducer() (Producer, error) {
	clusterServers, ok := os.LookupEnv("CLUSTER_SERVERS")
	if !ok {
		return nil, fmt.Errorf("CLUSTER_SERVERS was not defined")
	}

	// Default partitioner is consistent_random, which assigns the same keys
	// to the same partition, and randomly assigns null/empty keys
	return kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": clusterServers,
		"group.id":          "gateway-state-stream-aggregator-producers",
	})
}

func NewStreamAggregator(
	consumerFactory ConsumerFactory,
	decoder Decoder,
	producerFactory ProducerFactory,
) *StreamAggregator {
	aggregator := &StreamAggregator{
		consumerFactory: consumerFactory,
		decoder:         decoder,
		producerFactory: producerFactory,
		stopFlag:        &atomic.Value{},
	}
	aggregator.stopFlag.Store(false)
	return aggregator
}

// Run is the function which runs the aggregator, which listens to all relevant topics and combines them into a single
// gateway state topic. Does not return unless an error is encountered.
func (aggregator *StreamAggregator) Run() error {
	glog.V(2).Info("Starting stream aggregator in gateways materializer")

	consumer, err := aggregator.consumerFactory()
	if err != nil {
		return fmt.Errorf("Error creating Kafka consumer for stream aggregator: %s", err)
	}
	defer consumer.Close()

	producer, err := aggregator.producerFactory()
	if err != nil {
		return fmt.Errorf("Error creating Kafka producer for stream aggregator: %s", err)
	}
	defer producer.Close()

	err = consumer.SubscribeTopics([]string{consumerTopicsRegex}, nil)
	if err != nil {
		return fmt.Errorf("Error subscribing to aggregator topics: %s", err)
	}

	for !aggregator.stopFlag.Load().(bool) {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			return fmt.Errorf("Error reading message from aggregator consumer: %s", err)
		}
		gatewayUpdate, err := aggregator.decoder.GetUpdateFromStreamAggregatorMessage(message)
		if err != nil {
			return fmt.Errorf("Error converting message to aggregated update: %s", err)
		}

		// Sometimes Kafka will send an empty event -- just ignore those
		if gatewayUpdate == nil {
			glog.Warning("Empty kafka event encountered in materializer gateways stream aggregator -- ignoring")
			continue
		}

		err = aggregator.pushMessage(gatewayUpdate, producer)
		if err != nil {
			return fmt.Errorf("Error pushing message to aggregated topic: %s", err)
		}
		// If no errors occurred while pushing the update, commit
		_, err = consumer.Commit()
		if err != nil {
			return fmt.Errorf("Materializer aggregator failed to commit offset: %s", err)
		}
	}
	return nil
}

func (aggregator *StreamAggregator) Stop() {
	aggregator.stopFlag.Store(true)
}

func (aggregator *StreamAggregator) pushMessage(update *KafkaGatewayUpdate, producer Producer) error {
	value, err := json.Marshal(update)
	if err != nil {
		return fmt.Errorf("Error marshaling aggregated update to JSON: %s", err)
	}
	topic := aggregatedTopic
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic: &topic,
			// The default partitioner is used, which just ensures messages with the same key
			// are on the same partition
			Partition: kafka.PartitionAny,
		},
		Key:   []byte(update.NetworkID),
		Value: value,
	}
	deliveryChan := make(chan kafka.Event)
	err = producer.Produce(message, deliveryChan)
	if err != nil {
		return fmt.Errorf("Error producing aggregated message: %s", err)
	}
	return verifyDelivery(deliveryChan)
}

func verifyDelivery(deliveryChan chan kafka.Event) error {
	e, err := getDeliveryReportEvent(deliveryChan)
	if err != nil {
		return err
	}
	m := e.(*kafka.Message)
	if m.TopicPartition.Error != nil {
		return fmt.Errorf("Error delivering aggregated message: %s", err)
	}
	return nil
}

func getDeliveryReportEvent(deliveryChan chan kafka.Event) (kafka.Event, error) {
	select {
	case e := <-deliveryChan:
		return e, nil
	case <-time.After(10 * time.Second):
		return nil, fmt.Errorf("Aggregator timed out waiting for delivery report from producer")
	}
}
