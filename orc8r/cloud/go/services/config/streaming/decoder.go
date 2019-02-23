/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/blacklist"
	"magma/orc8r/cloud/go/services/config/registry"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
)

type Decoder interface {
	GetUpdateFromMessage(message *kafka.Message) (ApplyableUpdate, error)
}

type debeziumUpdateEvent struct {
	Schema  map[string]interface{} `json:"schema"`
	Payload *debeziumUpdatePayload `json:"payload"`
}

type debeziumUpdatePayload struct {
	Op        string                 `json:"op"`
	Timestamp int64                  `json:"ts_ms"`
	Before    *debeziumRow           `json:"before"`
	After     *debeziumRow           `json:"after"`
	Source    map[string]interface{} `json:"source"`
}

// Type will be specified for config service table updates, otherwise it will
// be left blank
type debeziumRow struct {
	Type             string `json:"type"`
	Key              string `json:"key"`
	Value            string `json:"value"`
	Version          int    `json:"version"`
	GenerationNumber int    `json:"generation_number"`
}

type decoderImpl struct{}

func NewDecoder() Decoder {
	return &decoderImpl{}
}

func (d *decoderImpl) GetUpdateFromMessage(message *kafka.Message) (ApplyableUpdate, error) {
	nwId, tableName, err := getNetworkIdAndTableNameFromMessage(message)
	if err != nil {
		return nil, err
	}

	kafkaUpdate := &debeziumUpdateEvent{}
	err = json.Unmarshal(message.Value, kafkaUpdate)
	if err != nil {
		return nil, fmt.Errorf("Could not deserialize event: %s", err)
	}

	// No-op if there is no payload
	if kafkaUpdate.Payload == nil {
		return &NoOpUpdate{}, nil
	}

	switch tableName {
	case "configurations":
		return d.getConfigUpdateFromMessage(kafkaUpdate, nwId)
	case "gatewayrecords":
		return d.getGatewayUpdateFromMessage(kafkaUpdate, nwId)
	case "tierversions":
		return d.getTierUpdateFromMessage(kafkaUpdate, nwId)
	default:
		return nil, fmt.Errorf("Unsupported topic: %s", *message.TopicPartition.Topic)
	}
}

func (*decoderImpl) getConfigUpdateFromMessage(kafkaUpdate *debeziumUpdateEvent, networkId string) (ApplyableUpdate, error) {
	changeOp := ChangeOperation(kafkaUpdate.Payload.Op)
	switch changeOp {
	case CreateOperation, ReadOperation, UpdateOperation:
		if blacklist.IsConfigBlacklisted(kafkaUpdate.Payload.After.Type) {
			// Blacklisted config types may not have registered config managers
			glog.Warningf("Ignoring config update for blacklisted type %s", kafkaUpdate.Payload.After.Type)
			return &NoOpUpdate{}, nil
		}

		valueDecoded, err := base64.StdEncoding.DecodeString(kafkaUpdate.Payload.After.Value)
		if err != nil {
			return nil, err
		}
		valueIface, err := registry.UnmarshalConfig(kafkaUpdate.Payload.After.Type, valueDecoded)
		if err != nil {
			return nil, err
		}

		return &ConfigUpdate{
			NetworkId:  networkId,
			ConfigType: kafkaUpdate.Payload.After.Type,
			ConfigKey:  kafkaUpdate.Payload.After.Key,
			Operation:  changeOp,
			NewValue:   valueIface,
		}, nil
	case DeleteOperation:
		if blacklist.IsConfigBlacklisted(kafkaUpdate.Payload.Before.Type) {
			glog.Warningf("Ignoring config deletion for blacklisted type %s", kafkaUpdate.Payload.Before.Type)
			return &NoOpUpdate{}, nil
		}

		return &ConfigUpdate{
			NetworkId:  networkId,
			ConfigType: kafkaUpdate.Payload.Before.Type,
			ConfigKey:  kafkaUpdate.Payload.Before.Key,
			Operation:  changeOp,
			NewValue:   nil,
		}, nil
	default:
		return nil, fmt.Errorf("Unrecognized stream operation %s", changeOp)
	}
}

func (*decoderImpl) getGatewayUpdateFromMessage(kafkaUpdate *debeziumUpdateEvent, networkId string) (*GatewayUpdate, error) {
	changeOp := ChangeOperation(kafkaUpdate.Payload.Op)
	switch changeOp {
	case CreateOperation, ReadOperation, UpdateOperation:
		return &GatewayUpdate{
			NetworkId: networkId,
			GatewayId: kafkaUpdate.Payload.After.Key,
			Operation: changeOp,
		}, nil
	case DeleteOperation:
		return &GatewayUpdate{
			NetworkId: networkId,
			GatewayId: kafkaUpdate.Payload.Before.Key,
			Operation: changeOp,
		}, nil
	default:
		return nil, fmt.Errorf("Unrecognized stream operation %s", changeOp)
	}
}

func (*decoderImpl) getTierUpdateFromMessage(kafkaUpdate *debeziumUpdateEvent, networkId string) (*TierUpdate, error) {
	changeOp := ChangeOperation(kafkaUpdate.Payload.Op)
	switch changeOp {
	case CreateOperation, ReadOperation, UpdateOperation:
		tierInfo := &upgrade_protos.TierInfo{}
		valueDecoded, err := base64.StdEncoding.DecodeString(kafkaUpdate.Payload.After.Value)
		if err != nil {
			return nil, err
		}
		err = protos.Unmarshal(valueDecoded, tierInfo)
		if err != nil {
			return nil, err
		}

		return &TierUpdate{
			NetworkId:   networkId,
			TierId:      kafkaUpdate.Payload.After.Key,
			Operation:   changeOp,
			TierVersion: tierInfo.Version,
			TierImages:  tierInfo.Images,
		}, nil
	case DeleteOperation:
		return &TierUpdate{
			NetworkId: networkId,
			TierId:    kafkaUpdate.Payload.Before.Key,
			Operation: changeOp,
		}, nil
	default:
		return nil, fmt.Errorf("Unrecognized stream operation %s", changeOp)
	}
}

func getNetworkIdAndTableNameFromMessage(message *kafka.Message) (string, string, error) {
	topicName := *message.TopicPartition.Topic
	matches := streamTopicRe.FindStringSubmatch(topicName)
	if len(matches) < 3 {
		return "", "", fmt.Errorf("Could not parse network and table name from topic name %s", topicName)
	}
	return matches[1], matches[2], nil
}
