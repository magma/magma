/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
	"magma/orc8r/cloud/go/services/materializer/gateways/streaming/mocks"
)

func TestStateRecorderRun(t *testing.T) {
	mockConsumer := &mocks.Consumer{}
	mockDecoder := &mocks.Decoder{}

	mockConsumer.On("Close").Return(nil)
	mockConsumer.On("Subscribe", "magma.public.gatewaystates", mock.Anything).Return(nil)
	mockConsumer.On("Commit").Return(nil, nil)

	// First iteration: normal update
	message := &kafka.Message{
		Key: []byte("net1"),
		TopicPartition: kafka.TopicPartition{
			Offset: 12345,
		},
	}

	mockPayload := &mocks.UpdatePayload{}
	mockPayload.On("Apply", "c", "net1", int64(12345), nil).Once().Return(nil)
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Statuses,
		Operation:  "c",
		NetworkID:  "net1",
		Payload:    mockPayload,
	}
	mockConsumer.On("ReadMessage", time.Duration(-1)).Once().Return(message, nil)
	mockDecoder.On("GetUpdateFromStateRecorderMessage", message).Once().Return(update, nil)

	// Second iteration: delete record
	update = &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Records,
		Operation:  "d",
		NetworkID:  "net2",
		Payload:    mockPayload,
	}
	message = &kafka.Message{
		Key: []byte("net2"),
		TopicPartition: kafka.TopicPartition{
			Offset: 12345,
		},
	}
	mockConsumer.On("ReadMessage", time.Duration(-1)).Once().Return(message, nil)
	mockDecoder.On("GetUpdateFromStateRecorderMessage", message).Once().Return(update, nil)
	mockPayload.On("Apply", "d", "net2", int64(12345), nil).Once().Return(nil)

	// End test
	mockConsumer.On("ReadMessage", time.Duration(-1)).Once().Return(nil, fmt.Errorf("test complete"))

	consumerFactory := func() (streaming.Consumer, error) { return mockConsumer, nil }
	recorder := streaming.NewStateRecorder(consumerFactory, mockDecoder, nil)
	err := recorder.Run()
	assert.EqualError(t, err, "Error reading message from StateRecorder: test complete")

	mockConsumer.AssertExpectations(t)
	mockDecoder.AssertExpectations(t)
	mockPayload.AssertExpectations(t)

	mockConsumer.AssertNumberOfCalls(t, "Subscribe", 1)
	mockConsumer.AssertNumberOfCalls(t, "Close", 1)
	mockConsumer.AssertNumberOfCalls(t, "ReadMessage", 3)
	mockConsumer.AssertNumberOfCalls(t, "Commit", 2)
	mockDecoder.AssertNumberOfCalls(t, "GetUpdateFromStateRecorderMessage", 2)
	mockPayload.AssertNumberOfCalls(t, "Apply", 2)
}
