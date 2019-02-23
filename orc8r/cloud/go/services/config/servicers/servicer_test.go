/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers_test

import (
	"context"
	"errors"
	"sort"
	"testing"

	"magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/services/config/servicers"
	"magma/orc8r/cloud/go/services/config/storage"
	"magma/orc8r/cloud/go/services/config/storage/mocks"

	"github.com/stretchr/testify/assert"
)

func TestConfigService_GetConfig(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("GetConfig", "network", "type", "key").
		Return(&storage.ConfigValue{Value: []byte("value1"), Version: 42}, nil)
	store.On("GetConfig", "network", "type2", "key2").
		Return(nil, errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.GetOrDeleteConfigRequest{NetworkId: "network", Type: "type", Key: "key"}
	actual, err := service.GetConfig(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, []byte("value1"), actual.Value)

	req.Type = "type2"
	req.Key = "key2"
	_, err = service.GetConfig(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error retrieving config: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_GetConfigs(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("GetConfigs", "network", &storage.FilterCriteria{Type: "type"}).
		Return(
			map[storage.TypeAndKey]*storage.ConfigValue{
				{Type: "type", Key: "key1"}: {Value: []byte("value1"), Version: 42},
				{Type: "type", Key: "key2"}: {Value: []byte("value2"), Version: 1},
			}, nil,
		)
	store.On("GetConfigs", "network", &storage.FilterCriteria{Type: "type2", Key: "key2"}).
		Return(nil, errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.GetOrDeleteConfigsRequest{NetworkId: "network", Filter: &protos.ConfigFilter{Type: "type"}}
	actual, err := service.GetConfigs(ctx, req)
	assert.NoError(t, err)
	sort.Slice(actual.Configs, func(i, j int) bool { return actual.Configs[i].Key < actual.Configs[j].Key })
	expected := []*protos.Config{
		{Type: "type", Key: "key1", Value: []byte("value1")},
		{Type: "type", Key: "key2", Value: []byte("value2")},
	}
	assert.Equal(t, expected, actual.Configs)

	req.Filter = &protos.ConfigFilter{Type: "type2", Key: "key2"}
	_, err = service.GetConfigs(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error retrieving configs: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_ListKeysForType(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("ListKeysForType", "network", "type").
		Return([]string{"key1", "key2"}, nil)
	store.On("ListKeysForType", "network", "type2").
		Return(nil, errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.ListKeysForTypeRequest{NetworkId: "network", Type: "type"}
	actual, err := service.ListKeysForType(ctx, req)
	assert.NoError(t, err)
	assert.Equal(t, []string{"key1", "key2"}, actual.Keys)

	req.Type = "type2"
	_, err = service.ListKeysForType(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error retrieving keys for type: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_CreateConfig(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("CreateConfig", "network", "type", "key", []byte("value")).
		Return(nil)
	store.On("CreateConfig", "network", "type2", "key2", []byte("value2")).
		Return(errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.CreateOrUpdateConfigRequest{NetworkId: "network", Type: "type", Key: "key", Value: []byte("value")}
	_, err := service.CreateConfig(ctx, req)
	assert.NoError(t, err)

	req.Type = "type2"
	req.Key = "key2"
	req.Value = []byte("value2")
	_, err = service.CreateConfig(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error creating config: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_UpdateConfig(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("UpdateConfig", "network", "type", "key", []byte("value")).
		Return(nil)
	store.On("UpdateConfig", "network", "type2", "key2", []byte("value2")).
		Return(errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.CreateOrUpdateConfigRequest{NetworkId: "network", Type: "type", Key: "key", Value: []byte("value")}
	_, err := service.UpdateConfig(ctx, req)
	assert.NoError(t, err)

	req.Type = "type2"
	req.Key = "key2"
	req.Value = []byte("value2")
	_, err = service.UpdateConfig(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error updating config: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_DeleteConfig(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("DeleteConfig", "network", "type", "key").
		Return(nil)
	store.On("DeleteConfig", "network", "type2", "key2").
		Return(errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.GetOrDeleteConfigRequest{NetworkId: "network", Type: "type", Key: "key"}
	_, err := service.DeleteConfig(ctx, req)
	assert.NoError(t, err)

	req.Type = "type2"
	req.Key = "key2"
	_, err = service.DeleteConfig(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error deleting config: Mock store error")
	store.AssertExpectations(t)
}

func TestConfigService_DeleteConfigs(t *testing.T) {
	store := &mocks.ConfigurationStorage{}
	service := servicers.NewConfigService(store)

	store.On("DeleteConfigs", "network", &storage.FilterCriteria{Type: "type"}).
		Return(nil)
	store.On("DeleteConfigs", "network", &storage.FilterCriteria{Type: "type2", Key: "key2"}).
		Return(errors.New("Mock store error"))

	ctx := context.Background()
	req := &protos.GetOrDeleteConfigsRequest{NetworkId: "network", Filter: &protos.ConfigFilter{Type: "type"}}
	_, err := service.DeleteConfigs(ctx, req)
	assert.NoError(t, err)

	req.Filter = &protos.ConfigFilter{Type: "type2", Key: "key2"}
	_, err = service.DeleteConfigs(ctx, req)
	assert.EqualError(t, err, "rpc error: code = Aborted desc = Error deleting configs: Mock store error")
	store.AssertExpectations(t)
}
