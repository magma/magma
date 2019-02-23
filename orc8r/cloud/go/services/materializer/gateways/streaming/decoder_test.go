/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"encoding/json"
	"testing"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"

	"magma/orc8r/cloud/go/services/materializer/gateways/streaming"
)

func TestGetUpdateFromStreamAggregatorMessage_Configurations(t *testing.T) {
	testCreateMessageConversionWithTopics(
		t,
		"magma.public.net1_configurations",
		streaming.Configurations,
		"net1",
		&streaming.GatewayConfigUpdate{
			ConfigType:  "test_type",
			ConfigKey:   "key",
			ConfigBytes: "value",
		},
	)
}

func TestGetUpdateFromStreamAggregatorMessage_Statuses(t *testing.T) {
	testCreateMessageConversionWithTopics(
		t,
		"magma.public.net1_gwstatus",
		streaming.Statuses,
		"net1",
		&streaming.GatewayStatusUpdate{
			GatewayID:   "key",
			StatusBytes: "value",
		},
	)
}

func TestGetUpdateFromStreamAggregatorMessage_Records(t *testing.T) {
	testCreateMessageConversionWithTopics(
		t,
		"magma.public.net1_gatewayrecords",
		streaming.Records,
		"net1",
		&streaming.GatewayRecordUpdate{
			GatewayID:   "key",
			RecordBytes: "value",
		},
	)
}

func testCreateMessageConversionWithTopics(
	t *testing.T,
	topic string,
	expectedType streaming.UpdateType,
	expectedNetworkID string,
	expectedPayload streaming.UpdatePayload,
) {
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(StreamValueFixtureCreate),
	}
	expectedAggregatedUpdate := &streaming.KafkaGatewayUpdate{
		UpdateType: expectedType,
		Operation:  "c",
		NetworkID:  expectedNetworkID,
		Payload:    expectedPayload,
	}
	decoder := streaming.NewDecoder()
	actualAggregatedUpdate, err := decoder.GetUpdateFromStreamAggregatorMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, expectedAggregatedUpdate, actualAggregatedUpdate)
}

func TestGetUpdateFromStreamAggregatorMessage_Deletion(t *testing.T) {
	topic := "magma.public.net1_gatewayrecords"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(StreamValueFixtureDelete),
	}
	expectedAggregatedUpdate := &streaming.KafkaGatewayUpdate{
		UpdateType: "gatewayrecords",
		Operation:  "d",
		NetworkID:  "net1",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID:   "key",
			RecordBytes: "value",
		},
	}
	decoder := streaming.NewDecoder()
	actualAggregatedUpdate, err := decoder.GetUpdateFromStreamAggregatorMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, expectedAggregatedUpdate, actualAggregatedUpdate)
}

func TestGetUpdateFromStateRecorderMessage_Empty(t *testing.T) {
	topic := "magma.public.net1_gatewayrecords"
	message := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte(StreamValueFixtureEmpty),
	}
	decoder := streaming.NewDecoder()
	update, err := decoder.GetUpdateFromStreamAggregatorMessage(message)
	assert.NoError(t, err)
	assert.Nil(t, update)
}

func TestGetUpdateFromStateRecorderMessage_Config(t *testing.T) {
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Configurations,
		Operation:  "operation",
		NetworkID:  "networkID",
		Payload: &streaming.GatewayConfigUpdate{
			ConfigKey:   "net1",
			ConfigType:  "test_type",
			ConfigBytes: "test",
		},
	}
	testGetUpdateFromStateRecorderMessage(t, update)
}

func TestGetUpdateFromStateRecorderMessage_Status(t *testing.T) {
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Statuses,
		Operation:  "operation",
		NetworkID:  "networkID",
		Payload: &streaming.GatewayStatusUpdate{
			GatewayID:   "gw1",
			StatusBytes: "value",
		},
	}
	testGetUpdateFromStateRecorderMessage(t, update)
}

func TestGetUpdateFromStateRecorderMessage_Record(t *testing.T) {
	update := &streaming.KafkaGatewayUpdate{
		UpdateType: streaming.Records,
		Operation:  "operation",
		NetworkID:  "networkID",
		Payload: &streaming.GatewayRecordUpdate{
			GatewayID:   "gw1",
			RecordBytes: "value",
		},
	}
	testGetUpdateFromStateRecorderMessage(t, update)
}

func testGetUpdateFromStateRecorderMessage(t *testing.T, update *streaming.KafkaGatewayUpdate) {
	updateBytes, err := json.Marshal(update)
	assert.NoError(t, err)
	message := &kafka.Message{
		Value: updateBytes,
	}
	decoder := streaming.NewDecoder()
	actualUpdate, err := decoder.GetUpdateFromStateRecorderMessage(message)
	assert.NoError(t, err)
	assert.Equal(t, update, actualUpdate)
}

const StreamValueFixtureCreate = `
{
	"schema": {
		"type": "struct",
		"fields": [
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": false,
						"field": "type"
					},
					{
						"type": "string",
						"optional": false,
						"field": "key"
					},
					{
						"type": "bytes",
						"optional": true,
						"field": "value"
					},
					{
						"type": "int32",
						"optional": false,
						"field": "version"
					}
				],
				"optional": true,
				"name": "magma.public.n1_configurations.Value",
				"field": "before"
			},
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": false,
						"field": "type"
					},
					{
						"type": "string",
						"optional": false,
						"field": "key"
					},
					{
						"type": "bytes",
						"optional": true,
						"field": "value"
					},
					{
						"type": "int32",
						"optional": false,
						"field": "version"
					}
				],
				"optional": true,
				"name": "magma.public.n1_configurations.Value",
				"field": "after"
			},
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": true,
						"field": "version"
					},
					{
						"type": "string",
						"optional": false,
						"field": "name"
					},
					{
						"type": "int64",
						"optional": true,
						"field": "ts_usec"
					},
					{
						"type": "int32",
						"optional": true,
						"field": "txId"
					},
					{
						"type": "int64",
						"optional": true,
						"field": "lsn"
					},
					{
						"type": "boolean",
						"optional": true,
						"default": false,
						"field": "snapshot"
					},
					{
						"type": "boolean",
						"optional": true,
						"field": "last_snapshot_record"
					}
				],
				"optional": false,
				"name": "io.debezium.connector.postgresql.Source",
				"field": "source"
			},
			{
				"type": "string",
				"optional": false,
				"field": "op"
			},
			{
				"type": "int64",
				"optional": true,
				"field": "ts_ms"
			}
		],
		"optional": false,
		"name": "magma.public.n1_configurations.Envelope"
	},
	"payload": {
		"before": null,
		"after": {
			"type": "test_type",
			"key": "key",
			"value": "value",
			"version": 0
		},
		"source": {
			"version": "0.7.5",
			"name": "magma",
			"ts_usec": 1528753424430907000,
			"txId": 601,
			"lsn": 22496127,
			"snapshot": false,
			"last_snapshot_record": null
		},
		"op": "c",
		"ts_ms": 1528831778308
	}
}`

const StreamValueFixtureDelete = `
{
	"schema": {
		"type": "struct",
		"fields": [
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": false,
						"field": "type"
					},
					{
						"type": "string",
						"optional": false,
						"field": "key"
					},
					{
						"type": "bytes",
						"optional": true,
						"field": "value"
					},
					{
						"type": "int32",
						"optional": false,
						"field": "version"
					}
				],
				"optional": true,
				"name": "magma.public.n1_configurations.Value",
				"field": "before"
			},
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": false,
						"field": "type"
					},
					{
						"type": "string",
						"optional": false,
						"field": "key"
					},
					{
						"type": "bytes",
						"optional": true,
						"field": "value"
					},
					{
						"type": "int32",
						"optional": false,
						"field": "version"
					}
				],
				"optional": true,
				"name": "magma.public.n1_configurations.Value",
				"field": "after"
			},
			{
				"type": "struct",
				"fields": [
					{
						"type": "string",
						"optional": true,
						"field": "version"
					},
					{
						"type": "string",
						"optional": false,
						"field": "name"
					},
					{
						"type": "int64",
						"optional": true,
						"field": "ts_usec"
					},
					{
						"type": "int32",
						"optional": true,
						"field": "txId"
					},
					{
						"type": "int64",
						"optional": true,
						"field": "lsn"
					},
					{
						"type": "boolean",
						"optional": true,
						"default": false,
						"field": "snapshot"
					},
					{
						"type": "boolean",
						"optional": true,
						"field": "last_snapshot_record"
					}
				],
				"optional": false,
				"name": "io.debezium.connector.postgresql.Source",
				"field": "source"
			},
			{
				"type": "string",
				"optional": false,
				"field": "op"
			},
			{
				"type": "int64",
				"optional": true,
				"field": "ts_ms"
			}
		],
		"optional": false,
		"name": "magma.public.n1_configurations.Envelope"
	},
	"payload": {
		"before": {
			"key": "key",
			"value": "value",
			"version": 0
		},
		"after": null,
		"source": {
			"version": "0.7.5",
			"name": "magma",
			"ts_usec": 1528753424430907000,
			"txId": 601,
			"lsn": 22496127,
			"snapshot": false,
			"last_snapshot_record": null
		},
		"op": "d",
		"ts_ms": 1528831778308
	}
}`

const StreamValueFixtureEmpty = `
{
	"schema": null,
	"payload": null
}`
