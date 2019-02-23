/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
	"magma/orc8r/cloud/go/services/materializer/gateways/streaming/mocks"
)

func TestAggregatorRun(t *testing.T) {
	mockConsumer := &mocks.Consumer{}
	mockProducer := &mocks.Producer{}
	mockDecoder := &mocks.Decoder{}

	mockConsumer.On(
		"SubscribeTopics",
		[]string{"^magma\\.public\\.(.+)_(configurations|gwstatus|gatewayrecords)$"},
		mock.Anything,
	).Return(nil)

	message := &kafka.Message{}
	mockConsumer.On("ReadMessage", time.Duration(-1)).Twice().Return(message, nil)
	mockConsumer.On("ReadMessage", time.Duration(-1)).Once().Return(nil, fmt.Errorf("Test completed."))
	mockConsumer.On("Commit").Return(nil, nil)
	gatewayUpdate := &streaming.KafkaGatewayUpdate{
		UpdateType: "gatewayrecords",
		Operation:  "c",
		NetworkID:  "net1",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID:   "gw1",
			RecordBytes: "value",
		},
	}
	mockDecoder.On("GetUpdateFromStreamAggregatorMessage", message).Return(gatewayUpdate, nil)

	gatewayUpdateBytes, err := json.Marshal(gatewayUpdate)
	assert.NoError(t, err)
	topicName := "magma.public.gatewaystates"
	aggregatedMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topicName,
			Partition: kafka.PartitionAny,
		},
		Key:   []byte("net1"),
		Value: gatewayUpdateBytes,
	}
	mockProducer.On("Produce", aggregatedMessage, mock.Anything).Return(
		func(message *kafka.Message, deliveryChan chan kafka.Event) error {
			go func() {
				deliveryChan <- &kafka.Message{}
			}()
			return nil
		},
	)

	mockConsumer.On("Close").Return(nil)
	mockProducer.On("Close").Return()

	consumerFactory := func() (streaming.Consumer, error) { return mockConsumer, nil }
	producerFactory := func() (streaming.Producer, error) { return mockProducer, nil }
	streamAggregator := streaming.NewStreamAggregator(consumerFactory, mockDecoder, producerFactory)
	err = streamAggregator.Run()
	assert.EqualError(t, err, "Error reading message from aggregator consumer: Test completed.")

	mockConsumer.AssertExpectations(t)
	mockProducer.AssertExpectations(t)
	mockDecoder.AssertExpectations(t)

	mockConsumer.AssertNumberOfCalls(t, "SubscribeTopics", 1)
	mockConsumer.AssertNumberOfCalls(t, "ReadMessage", 3)
	mockConsumer.AssertNumberOfCalls(t, "Commit", 2)
	mockDecoder.AssertNumberOfCalls(t, "GetUpdateFromStreamAggregatorMessage", 2)
	mockProducer.AssertNumberOfCalls(t, "Produce", 2)
	mockConsumer.AssertNumberOfCalls(t, "Close", 1)
	mockProducer.AssertNumberOfCalls(t, "Close", 1)
}
