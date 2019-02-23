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
	"sync/atomic"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"

	"magma/orc8r/cloud/go/services/materializer/gateways/metrics"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

// StateRecorder listens to the aggregated magma.public.gatewaystates topic for any gateway updates, and updates
// the materialized view store accordingly.
type StateRecorder struct {
	consumerFactory ConsumerFactory
	decoder         Decoder
	store           storage.GatewayViewStorage
	stopFlag        *atomic.Value
}

func NewStateRecorder(
	consumerFactory ConsumerFactory,
	decoder Decoder,
	store storage.GatewayViewStorage,
) *StateRecorder {
	recorder := &StateRecorder{
		consumerFactory: consumerFactory,
		decoder:         decoder,
		store:           store,
		stopFlag:        &atomic.Value{},
	}
	recorder.stopFlag.Store(false)
	return recorder
}

func NewStateRecorderConsumer() (Consumer, error) {
	clusterServers, ok := os.LookupEnv("CLUSTER_SERVERS")
	if !ok {
		return nil, fmt.Errorf("CLUSTER_SERVERS was not defined")
	}
	return kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":   clusterServers,
		"group.id":            "gateway-state-state-recorder-consumers",
		"auto.offset.reset":   "earliest",
		"enable.auto.commit":  false,
		"metadata.max.age.ms": 30000,
	})
}

// Run is the main function for StateRecorder. It does not return unless an error is encountered.
func (recorder *StateRecorder) Run() error {
	glog.V(2).Info("Starting stream recorder in gateways materializer")

	consumer, err := recorder.consumerFactory()
	if err != nil {
		return err
	}
	defer consumer.Close()

	err = consumer.Subscribe(aggregatedTopic, nil)
	if err != nil {
		return err
	}
	for !recorder.stopFlag.Load().(bool) {
		message, err := consumer.ReadMessage(-1)
		if err != nil {
			return fmt.Errorf("Error reading message from StateRecorder: %s", err)
		}
		metrics.UpdatesProcessed.Inc()

		update, err := recorder.decoder.GetUpdateFromStateRecorderMessage(message)
		if err != nil {
			return fmt.Errorf("Error parsing aggregated update from message: %s", err)
		}

		err = update.ApplyUpdate(recorder.store, int64(message.TopicPartition.Offset))
		if err != nil {
			return err
		}

		_, err = consumer.Commit()
		if err != nil {
			return fmt.Errorf("Materializer state recorder failed to commit offset: %s", err)
		}
	}
	return nil
}

func (recorder *StateRecorder) Stop() {
	recorder.stopFlag.Store(true)
}
