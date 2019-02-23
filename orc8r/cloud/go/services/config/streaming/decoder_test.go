/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming_test

import (
	"testing"

	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/config/streaming"
	upgrade_protos "magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/stretchr/testify/assert"
)

func TestDecoderImpl_GetUpdateFromMessage_Configs(t *testing.T) {
	registry.ClearRegistryForTesting()
	registry.RegisterConfigManager(&testConfigManager{})

	decoder := streaming.NewDecoder()

	topic := "magma.public.network_name_1234_configurations"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte(ConfigCreatedValueFixture),
	}

	actual, err := decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected := &streaming.ConfigUpdate{
		NetworkId:  "network_name_1234",
		ConfigType: "test_config",
		ConfigKey:  "key",
		NewValue:   "value",
		Operation:  streaming.CreateOperation,
	}
	assert.Equal(t, expected, actual)

	msg.Value = []byte(ConfigDeletedValueFixture)
	actual, err = decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected = &streaming.ConfigUpdate{
		NetworkId:  "network_name_1234",
		Operation:  streaming.DeleteOperation,
		ConfigType: "wifi_network",
		ConfigKey:  "n1",
	}
	assert.Equal(t, expected, actual)

	topic = "magma.public.network_configurations_configurations_configurations"
	msg.TopicPartition.Topic = &topic
	msg.Value = []byte(ConfigCreatedValueFixture)

	actual, err = decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected = &streaming.ConfigUpdate{
		NetworkId:  "network_configurations_configurations",
		ConfigType: "test_config",
		ConfigKey:  "key",
		NewValue:   "value",
		Operation:  streaming.CreateOperation,
	}
	assert.Equal(t, expected, actual)
}

func TestDecoderImpl_GetUpdateFromMessage_Gateways(t *testing.T) {
	decoder := streaming.NewDecoder()

	topic := "magma.public.network_name_1234_gatewayrecords"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte(GatewayCreatedValueFixture),
	}

	actual, err := decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected := &streaming.GatewayUpdate{
		NetworkId: "network_name_1234",
		GatewayId: "vm1",
		Operation: streaming.CreateOperation,
	}
	assert.Equal(t, expected, actual)

	msg.Value = []byte(GatewayDeletedValueFixture)
	actual, err = decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected = &streaming.GatewayUpdate{
		NetworkId: "network_name_1234",
		Operation: streaming.DeleteOperation,
		GatewayId: "deleteme",
	}
	assert.Equal(t, expected, actual)
}

func TestDecoderImpl_GetUpdateFromMessage_Tier(t *testing.T) {
	decoder := streaming.NewDecoder()

	topic := "magma.public.network_name_1234_tierversions"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte(TierCreatedValueFixture),
	}

	actual, err := decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected := &streaming.TierUpdate{
		NetworkId:   "network_name_1234",
		Operation:   streaming.CreateOperation,
		TierId:      "default",
		TierVersion: "1.2.3-4",
		TierImages:  []*upgrade_protos.ImageSpec{{Name: "img1", Order: 1}},
	}
	assert.Equal(t, expected, actual)

	msg.Value = []byte(TierDeletedValueFixture)
	actual, err = decoder.GetUpdateFromMessage(msg)
	assert.NoError(t, err)
	expected = &streaming.TierUpdate{
		NetworkId: "network_name_1234",
		Operation: streaming.DeleteOperation,
		TierId:    "deleteme",
	}
	assert.Equal(t, expected, actual)
}

func TestDecoderImpl_GetUpdateFromMessage_InvalidTopic(t *testing.T) {
	topic := "magma.public.network_sometable"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic},
		Value:          []byte(ConfigCreatedValueFixture),
	}

	decoder := streaming.NewDecoder()
	_, err := decoder.GetUpdateFromMessage(msg)
	assert.EqualError(t, err, "Could not parse network and table name from topic name magma.public.network_sometable")
}

type testConfigManager struct{}

func (*testConfigManager) GetConfigType() string {
	return "test_config"
}

func (*testConfigManager) GetGatewayIdsForConfig(networkId string, configKey string) ([]string, error) {
	return []string{configKey}, nil
}

func (*testConfigManager) MarshalConfig(config interface{}) ([]byte, error) {
	return []byte(config.(string)), nil
}

func (*testConfigManager) UnmarshalConfig(message []byte) (interface{}, error) {
	return string(message), nil
}

// Actual Kafka message sent by Debezium captured by the console consumer and
// piped into jq. Payload fields have been replaced with mock values.
// Note debezium b64-encodes bytea value fields
const ConfigCreatedValueFixture = `
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
      "type": "test_config",
      "key": "key",
      "value": "dmFsdWU=",
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

const ConfigDeletedValueFixture = `
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
      "type": "wifi_network",
      "key": "n1",
      "value": null,
      "version": 0
    },
    "after": null,
    "source": {
      "version": "0.7.5",
      "name": "magma",
      "ts_usec": 1529631694500288000,
      "txId": 957,
      "lsn": 24143660,
      "snapshot": true,
      "last_snapshot_record": true
    },
    "op": "d",
    "ts_ms": 1529631694837
  }
}
`

const GatewayCreatedValueFixture = `
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_gatewayrecords.Value",
        "field": "before"
      },
      {
        "type": "struct",
        "fields": [
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_gatewayrecords.Value",
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
    "name": "magma.public.n1_gatewayrecords.Envelope"
  },
  "payload": {
    "before": null,
    "after": {
      "key": "vm1",
      "value": "ewogImh3SWQiOiB7CiAgImlkIjogImY3Y2JhN2M1LTViMzctNDM3OC05YmJiLWJjYWY3OTBjY2IyNiIKIH0sCiAibmFtZSI6ICJTb3V0aCBQYXJrJ3MgQ3RQYSBUb3duIFRvd2VyIiwKICJrZXkiOiB7CiAgImtleVR5cGUiOiAiRUNITyIsCiAgImtleSI6IG51bGwKIH0sCiAiaXAiOiAiIiwKICJwb3J0IjogMAp9",
      "generation_number": 0,
      "deleted": false
    },
    "source": {
      "version": "0.7.5",
      "name": "magma",
      "ts_usec": 1529453940684682000,
      "txId": 610,
      "lsn": 23056257,
      "snapshot": true,
      "last_snapshot_record": true
    },
    "op": "c",
    "ts_ms": 1529528059074
  }
}
`

const GatewayDeletedValueFixture = `
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_gatewayrecords.Value",
        "field": "before"
      },
      {
        "type": "struct",
        "fields": [
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_gatewayrecords.Value",
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
    "name": "magma.public.n1_gatewayrecords.Envelope"
  },
  "payload": {
    "before": {
      "key": "deleteme",
      "value": null,
      "generation_number": 0,
      "deleted": false
    },
    "after": null,
    "source": {
      "version": "0.7.5",
      "name": "magma",
      "ts_usec": 1529631549067167000,
      "txId": 956,
      "lsn": 24143234,
      "snapshot": true,
      "last_snapshot_record": true
    },
    "op": "d",
    "ts_ms": 1529631549411
  }
}
`

// From debezium. Deserializes to tierInfo proto:
// {
//   "id": "default",
//   "name": "default",
//   "version": "1.2.3-4",
//   "images": [
//     {
//       "name": "img1",
//       "order": 1
//     }
//   ]
// }
const TierCreatedValueFixture = `
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_tierversions.Value",
        "field": "before"
      },
      {
        "type": "struct",
        "fields": [
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_tierversions.Value",
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
    "name": "magma.public.n1_tierversions.Envelope"
  },
  "payload": {
    "before": null,
    "after": {
      "key": "default",
      "value": "ewogIm5hbWUiOiAiZGVmYXVsdCIsCiAidmVyc2lvbiI6ICIxLjIuMy00IiwKICJpbWFnZXMiOiBbCiAgewogICAibmFtZSI6ICJpbWcxIiwKICAgIm9yZGVyIjogIjEiCiAgfQogXQp9",
      "generation_number": 0,
      "deleted": false
    },
    "source": {
      "version": "0.7.5",
      "name": "magma",
      "ts_usec": 1529566402761506000,
      "txId": 948,
      "lsn": 24082051,
      "snapshot": true,
      "last_snapshot_record": true
    },
    "op": "c",
    "ts_ms": 1529566402762
  }
}
`

// Fixture from a deleted tier
const TierDeletedValueFixture = `
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_tierversions.Value",
        "field": "before"
      },
      {
        "type": "struct",
        "fields": [
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
            "field": "generation_number"
          },
          {
            "type": "boolean",
            "optional": false,
            "field": "deleted"
          }
        ],
        "optional": true,
        "name": "magma.public.n1_tierversions.Value",
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
    "name": "magma.public.n1_tierversions.Envelope"
  },
  "payload": {
    "before": {
      "key": "deleteme",
      "value": null,
      "generation_number": 0,
      "deleted": false
    },
    "after": null,
    "source": {
      "version": "0.7.5",
      "name": "magma",
      "ts_usec": 1529631208057605000,
      "txId": 950,
      "lsn": 24140064,
      "snapshot": true,
      "last_snapshot_record": true
    },
    "op": "d",
    "ts_ms": 1529631208396
  }
}
`
