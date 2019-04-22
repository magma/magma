/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// Package config contains a client API for accessing the configuration
// service, which manages entities representing generic configurations for
// other entities.

// Configuration objects are uniquely identified by the combination of a type
// and a key. The configuration type is used to look up the corresponding
// configuration manager in config/registry to marshal/unmarshal the
// configuration object as it moves to/from the GRPC service.

// To use the configuration service to manage your own configurations, first
// register your configuration manager in the registry. Then you will be able
// to use this API to manage your configs.
package config

import (
	"context"

	"magma/orc8r/cloud/go/errors"
	service_registry "magma/orc8r/cloud/go/registry"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"google.golang.org/grpc"
)

const (
	ServiceName = "CONFIG"

	// SerdeDomain is the domain for config serdes
	SerdeDomain = "config_manager"
)

func getConfigServiceClient() (protos.ConfigServiceClient, *grpc.ClientConn, error) {
	conn, err := service_registry.GetConnection(ServiceName)
	if err != nil {
		initErr := errors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, nil, initErr
	}
	return protos.NewConfigServiceClient(conn), conn, err
}

// Retrieve a specific config.
func GetConfig(networkId string, configType string, key string) (interface{}, error) {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigRequest{NetworkId: networkId, Type: configType, Key: key}
	val, err := client.GetConfig(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return serde.Deserialize(SerdeDomain, configType, val.GetValue())
}

// Fetch all configs matching a type.
func GetConfigsByType(networkId string, configType string) (map[storage.TypeAndKey]interface{}, error) {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigsRequest{
		NetworkId: networkId,
		Filter:    &protos.ConfigFilter{Type: configType},
	}
	vals, err := client.GetConfigs(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return makeGetConfigsResult(vals.GetConfigs())
}

// Fetch all configs matching a key.
func GetConfigsByKey(networkId string, key string) (map[storage.TypeAndKey]interface{}, error) {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigsRequest{
		NetworkId: networkId,
		Filter:    &protos.ConfigFilter{Key: key},
	}
	vals, err := client.GetConfigs(context.Background(), req)
	if err != nil {
		return nil, err
	}
	return makeGetConfigsResult(vals.GetConfigs())
}

// List all configuration keys for a specific type.
func ListKeysForType(networkId string, configType string) ([]string, error) {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	req := &protos.ListKeysForTypeRequest{NetworkId: networkId, Type: configType}
	keys, err := client.ListKeysForType(context.Background(), req)
	if err != nil {
		return nil, err
	}
	if len(keys.Keys) == 0 {
		keys.Keys = []string{}
	}
	return keys.Keys, nil
}

// Create a new config. This will error out if there is an existing config
// with the same type and key.
func CreateConfig(networkId string, configType string, key string, value interface{}) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	marshaledValue, err := serde.Serialize(SerdeDomain, configType, value)
	if err != nil {
		return err
	}

	req := &protos.CreateOrUpdateConfigRequest{NetworkId: networkId, Type: configType, Key: key, Value: marshaledValue}
	_, err = client.CreateConfig(context.Background(), req)
	return err
}

// Update an existing config. This will error out if there is no existing
// config with a matching type and key.
func UpdateConfig(networkId string, configType string, key string, updatedValue interface{}) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	marshaledValue, err := serde.Serialize(SerdeDomain, configType, updatedValue)
	if err != nil {
		return err
	}

	req := &protos.CreateOrUpdateConfigRequest{NetworkId: networkId, Type: configType, Key: key, Value: marshaledValue}
	_, err = client.UpdateConfig(context.Background(), req)
	return err
}

// Delete an existing config. This will error out if there is no existing
// config with a matching type and key.
func DeleteConfig(networkId string, configType string, key string) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigRequest{NetworkId: networkId, Type: configType, Key: key}
	_, err = client.DeleteConfig(context.Background(), req)
	return err
}

// Delete all configs matching a type
func DeleteConfigsByType(networkId string, configType string) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigsRequest{
		NetworkId: networkId,
		Filter:    &protos.ConfigFilter{Type: configType},
	}
	_, err = client.DeleteConfigs(context.Background(), req)
	return err
}

// Delete all configs matching a key
func DeleteConfigsByKey(networkId string, key string) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	req := &protos.GetOrDeleteConfigsRequest{
		NetworkId: networkId,
		Filter:    &protos.ConfigFilter{Key: key},
	}
	_, err = client.DeleteConfigs(context.Background(), req)
	return err

}

// Delete all the configs for all entities in a network
func DeleteAllNetworkConfigs(networkId string) error {
	client, conn, err := getConfigServiceClient()
	if err != nil {
		return err
	}
	defer conn.Close()

	req := &protos.NetworkIdRequest{NetworkId: networkId}
	_, err = client.DeleteAllConfigsForNetwork(context.Background(), req)
	return err
}

func makeGetConfigsResult(configs []*protos.Config) (map[storage.TypeAndKey]interface{}, error) {
	ret := make(map[storage.TypeAndKey]interface{}, len(configs))
	for _, configProto := range configs {
		k := storage.TypeAndKey{Type: configProto.Type, Key: configProto.Key}
		unmarshaledVal, err := serde.Deserialize(SerdeDomain, configProto.Type, configProto.Value)
		if err != nil {
			return nil, err
		}
		ret[k] = unmarshaledVal
	}
	return ret, nil
}
