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
	"reflect"
	"regexp"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type decoderImpl struct{}

func NewDecoder() Decoder {
	return &decoderImpl{}
}

type debeziumUpdateEvent struct {
	Payload *debeziumPayload `json:"payload"`
}

type debeziumPayload struct {
	Op     string       `json:"op"`
	Before *debeziumRow `json:"before"`
	After  *debeziumRow `json:"after"`
}

type debeziumRow struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

func (*decoderImpl) GetUpdateFromStreamAggregatorMessage(
	message *kafka.Message,
) (*KafkaGatewayUpdate, error) {
	networkID, updateType, err := getNetworkIDAndUpdateType(*message.TopicPartition.Topic)
	if err != nil {
		return nil, err
	}
	updateEvent, err := getUpdateEvent(message)
	if err != nil {
		return nil, err
	}
	// Kafka will sometimes throw in an empty event, so just return null and ignore it in the aggregator
	if updateEvent.Payload == nil {
		return nil, nil
	}
	row, err := getRow(updateEvent)
	if err != nil {
		return nil, err
	}
	payload, err := getUpdatePayload(updateType, row)
	if err != nil {
		return nil, err
	}
	update := &KafkaGatewayUpdate{
		UpdateType: updateType,
		Operation:  updateEvent.Payload.Op,
		NetworkID:  networkID,
		Payload:    payload,
	}
	return update, nil
}

func getNetworkIDAndUpdateType(topic string) (string, UpdateType, error) {
	re := regexp.MustCompile(consumerTopicsRegex)
	matches := re.FindStringSubmatch(topic)
	if len(matches) < 3 {
		return "", UpdateType(""), fmt.Errorf("Unexpected topic name: %s", topic)
	}
	return matches[1], UpdateType(matches[2]), nil
}

func getUpdateEvent(message *kafka.Message) (*debeziumUpdateEvent, error) {
	updateEvent := &debeziumUpdateEvent{}
	err := json.Unmarshal(message.Value, updateEvent)
	return updateEvent, err
}

func getRow(updateEvent *debeziumUpdateEvent) (*debeziumRow, error) {
	switch updateEvent.Payload.Op {
	case "c", "r", "u":
		return updateEvent.Payload.After, nil
	case "d":
		return updateEvent.Payload.Before, nil
	default:
		return nil, fmt.Errorf("Unrecognized op: %s", updateEvent.Payload.Op)
	}
}

func getUpdatePayload(updateType UpdateType, row *debeziumRow) (UpdatePayload, error) {
	switch updateType {
	case Configurations:
		return getConfigPayload(row), nil
	case Statuses:
		return getStatusPayload(row), nil
	case Records:
		return getRecordPayload(row), nil
	default:
		return nil, fmt.Errorf("Unrecognized update type: %s", string(updateType))
	}
}

func getConfigPayload(row *debeziumRow) *GatewayConfigUpdate {
	return &GatewayConfigUpdate{
		ConfigKey:   row.Key,
		ConfigType:  row.Type,
		ConfigBytes: row.Value,
	}
}

func getStatusPayload(row *debeziumRow) *GatewayStatusUpdate {
	return &GatewayStatusUpdate{
		GatewayID:   row.Key,
		StatusBytes: row.Value,
	}
}

func getRecordPayload(row *debeziumRow) *GatewayRecordUpdate {
	return &GatewayRecordUpdate{
		GatewayID:   row.Key,
		RecordBytes: row.Value,
	}
}

type kafkaUpdateUnmarshalTarget struct {
	UpdateType UpdateType
	Operation  string
	NetworkID  string
	Payload    map[string]interface{}
}

func (*decoderImpl) GetUpdateFromStateRecorderMessage(
	message *kafka.Message,
) (*KafkaGatewayUpdate, error) {
	update := &kafkaUpdateUnmarshalTarget{}
	err := json.Unmarshal(message.Value, update)
	if err != nil {
		return nil, err
	}
	payload, err := getPayloadFromUnmarshalTarget(update)
	if err != nil {
		return nil, err
	}
	return &KafkaGatewayUpdate{
		UpdateType: update.UpdateType,
		Operation:  update.Operation,
		NetworkID:  update.NetworkID,
		Payload:    payload,
	}, nil
}

func getPayloadFromUnmarshalTarget(update *kafkaUpdateUnmarshalTarget) (UpdatePayload, error) {
	payloadInterface, err := getPayload(update.UpdateType, update.Payload)
	if err != nil {
		return nil, err
	}
	payload, ok := payloadInterface.(UpdatePayload)
	if !ok {
		return nil, fmt.Errorf("Payload parsed into type which does not implement UpdatePayload")
	}
	return payload, nil
}

func getPayload(updateType UpdateType, payloadMap map[string]interface{}) (interface{}, error) {
	payloadObj := reflect.New(updateTypeRegistry[updateType]).Elem()
	for fieldName, fieldValue := range payloadMap {
		payloadFieldValue := payloadObj.FieldByName(fieldName)
		if !payloadFieldValue.IsValid() {
			return nil, fmt.Errorf("Error parsing update payload: no such field %s", fieldName)
		}
		if !payloadFieldValue.CanSet() {
			return nil, fmt.Errorf("Error parsing update payload: cannot set field %s", fieldName)
		}
		val := reflect.ValueOf(fieldValue)
		if !val.Type().AssignableTo(payloadFieldValue.Type()) {
			return nil, fmt.Errorf("Error parsing update payload: value not assignable to struct field %s", fieldName)
		}
		payloadFieldValue.Set(val)
	}
	return payloadObj.Addr().Interface(), nil
}
