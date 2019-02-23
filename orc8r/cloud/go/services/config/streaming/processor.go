/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"fmt"
	"os"
	"regexp"
	"sync/atomic"

	"magma/orc8r/cloud/go/services/config/streaming/storage"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

var streamTopicRe = regexp.MustCompile("^magma\\.public\\.([a-zA-Z0-9_]+)_(configurations|gatewayrecords|tierversions)$")

type StreamProcessor struct {
	store           storage.MconfigStorage
	decoder         Decoder
	consumerFactory StreamConsumerFactory
	stopFlag        *atomic.Value
}

func ConsumerFactoryImpl() (StreamConsumer, error) {
	clusterServers, ok := os.LookupEnv("CLUSTER_SERVERS")
	if !ok {
		return nil, fmt.Errorf("CLUSTER_SERVERS was not defined")
	}

	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   clusterServers,
		"group.id":            "mconfig-streamers",
		"auto.offset.reset":   "earliest",
		"enable.auto.commit":  false,
		"metadata.max.age.ms": 30000,
	})
}

func NewStreamProcessor(store storage.MconfigStorage, decoder Decoder, consumerFactory StreamConsumerFactory) *StreamProcessor {
	processor := &StreamProcessor{
		store:           store,
		decoder:         decoder,
		consumerFactory: consumerFactory,
		stopFlag:        &atomic.Value{},
	}
	processor.stopFlag.Store(false)
	return processor
}

func (processor *StreamProcessor) Run() error {
	c, err := processor.consumerFactory()
	if err != nil {
		return err
	}
	defer c.Close()

	// Subscribe to 3 topics per network - configs, gateway records (to know
	// when a gateway is registered or deleted), and tier versions to know
	// when to update a gateway's software version
	err = c.SubscribeTopics([]string{streamTopicRe.String()}, nil)
	if err != nil {
		return fmt.Errorf("Could not subscribe to Kafka topics: %s", err)
	}

	for {
		if processor.stopFlag.Load().(bool) {
			return nil
		}

		message, err := c.ReadMessage(-1)
		if err != nil {
			glog.Errorf("Consumer error: %v (%v)\n", err, message)
			return err
		}

		// Decode and apply an update
		update, err := processor.decoder.GetUpdateFromMessage(message)
		if err != nil {
			return fmt.Errorf("Mconfig streamer error while decoding stream update: %s\n%v", err, update)
		}
		offset := int64(message.TopicPartition.Offset)

		err = update.Apply(processor.store, offset)
		if err != nil {
			return fmt.Errorf("Mconfig streamer error while applying update: %s\n%v", err, update)
		}

		// If no errors occurred while persisting updates, commit the new offset
		_, err = c.Commit()
		if err != nil {
			glog.Errorf("Mconfig streamer error while committing offsets: %s", err)
			return err
		}
	}
}

func (processor *StreamProcessor) Stop() {
	processor.stopFlag.Store(true)
}
