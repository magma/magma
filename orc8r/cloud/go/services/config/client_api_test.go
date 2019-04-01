/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config_test

import (
	"strings"
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/storage"
	"magma/orc8r/cloud/go/services/config/test_init"

	"github.com/stretchr/testify/assert"
)

func TestConfigAPI(t *testing.T) {
	test_init.StartTestService(t)
	err := serde.RegisterSerdes(&type1Manager{}, &type2Manager{})
	assert.NoError(t, err)

	// Check contract for empty network
	actualConfigs, err := config.GetConfigsByType("network0", "type1")
	assert.NoError(t, err)
	assert.Equal(t, map[storage.TypeAndKey]interface{}{}, actualConfigs)

	actual, err := config.GetConfig("network0", "type1", "key")
	assert.NoError(t, err)
	assert.Nil(t, actual)

	// Create configs on 2 networks
	// network1: (type1, type2) X (key1, key2)
	err = config.CreateConfig("network1", "type1", "key1", "value1")
	assert.NoError(t, err)
	err = config.CreateConfig("network1", "type1", "key2", "value2")
	assert.NoError(t, err)
	err = config.CreateConfig("network1", "type2", "key1", "value3")
	assert.NoError(t, err)
	err = config.CreateConfig("network1", "type2", "key2", "value4")
	assert.NoError(t, err)

	// network2: (type1, type2) X (key1, key2)
	err = config.CreateConfig("network2", "type1", "key1", "value1")
	assert.NoError(t, err)
	err = config.CreateConfig("network2", "type1", "key2", "value2")
	assert.NoError(t, err)
	err = config.CreateConfig("network2", "type2", "key1", "value3")
	assert.NoError(t, err)
	err = config.CreateConfig("network2", "type2", "key2", "value4")
	assert.NoError(t, err)

	// Read back
	actualKeys, err := config.ListKeysForType("network1", "type1")
	assert.NoError(t, err)
	assert.Equal(t, []string{"key1", "key2"}, actualKeys)

	actualConfigs, err = config.GetConfigsByType("network2", "type2")
	assert.NoError(t, err)
	expectedConfigs := map[storage.TypeAndKey]interface{}{
		{Type: "type2", Key: "key1"}: "VALUE3",
		{Type: "type2", Key: "key2"}: "VALUE4",
	}
	assert.Equal(t, expectedConfigs, actualConfigs)

	actualConfigs, err = config.GetConfigsByKey("network1", "key1")
	assert.NoError(t, err)
	expectedConfigs = map[storage.TypeAndKey]interface{}{
		{Type: "type1", Key: "key1"}: "value1",
		{Type: "type2", Key: "key1"}: "VALUE3",
	}
	assert.Equal(t, expectedConfigs, actualConfigs)

	actualConfig, err := config.GetConfig("network1", "type2", "key1")
	assert.NoError(t, err)
	assert.Equal(t, "VALUE3", actualConfig)

	// Update-read
	err = config.UpdateConfig("network2", "type2", "key2", "value42")
	assert.NoError(t, err)
	actualConfig, err = config.GetConfig("network2", "type2", "key2")
	assert.NoError(t, err)
	assert.Equal(t, "VALUE42", actualConfig)

	// Update nonexisting
	err = config.UpdateConfig("network1", "type1", "key4", "value42")
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error updating config: Updating nonexistent config")

	// Delete single
	err = config.DeleteConfig("network2", "type2", "key1")
	assert.NoError(t, err)
	actualKeys, err = config.ListKeysForType("network2", "type2")
	assert.NoError(t, err)
	assert.Equal(t, []string{"key2"}, actualKeys)

	// Delete multiple
	err = config.DeleteConfigsByKey("network1", "key1")
	assert.NoError(t, err)
	actualConfigs, err = config.GetConfigsByKey("network1", "key1")
	assert.NoError(t, err)
	assert.Equal(t, map[storage.TypeAndKey]interface{}{}, actualConfigs)
}

type type1Manager struct{}

func (*type1Manager) GetDomain() string {
	return config.SerdeDomain
}

func (*type1Manager) GetType() string {
	return "type1"
}

func (*type1Manager) Serialize(config interface{}) ([]byte, error) {
	return []byte(config.(string)), nil
}

func (*type1Manager) Deserialize(message []byte) (interface{}, error) {
	return string(message), nil
}

type type2Manager struct{}

func (*type2Manager) GetDomain() string {
	return config.SerdeDomain
}

func (*type2Manager) GetType() string {
	return "type2"
}

func (*type2Manager) Serialize(config interface{}) ([]byte, error) {
	return []byte(strings.ToUpper(config.(string))), nil
}

func (*type2Manager) Deserialize(message []byte) (interface{}, error) {
	return string(message), nil
}
