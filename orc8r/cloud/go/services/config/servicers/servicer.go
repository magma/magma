/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package servicers

import (
	"errors"
	"sort"

	"magma/orc8r/cloud/go/protos"
	config_protos "magma/orc8r/cloud/go/services/config/protos"
	"magma/orc8r/cloud/go/services/config/storage"
	mstore "magma/orc8r/cloud/go/storage"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type ConfigService struct {
	store storage.ConfigurationStorage
}

func NewConfigService(store storage.ConfigurationStorage) config_protos.ConfigServiceServer {
	return &ConfigService{store: store}
}

func (service *ConfigService) GetConfig(context context.Context, req *config_protos.GetOrDeleteConfigRequest) (*config_protos.Config, error) {
	ret := &config_protos.Config{}
	if err := config_protos.ValidateGetOrDeleteConfigRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	configValue, err := service.store.GetConfig(req.GetNetworkId(), req.GetType(), req.GetKey())
	if err != nil {
		msgFormat := "Error retrieving config: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return &config_protos.Config{Type: req.Type, Key: req.Key, Value: configValue.Value}, nil
}

func (service *ConfigService) GetConfigs(context context.Context, req *config_protos.GetOrDeleteConfigsRequest) (*config_protos.GetConfigsResponse, error) {
	ret := &config_protos.GetConfigsResponse{}
	if err := config_protos.ValidateGetOrDeleteConfigsRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	configs, err := service.store.GetConfigs(req.GetNetworkId(), protoFilterToStorageFilter(req.GetFilter()))
	if err != nil {
		msgFormat := "Error retrieving configs: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return configValueMapToGetConfigsResponse(configs), nil
}

func (service *ConfigService) ListKeysForType(context context.Context, req *config_protos.ListKeysForTypeRequest) (*config_protos.ListKeysForTypeResponse, error) {
	ret := &config_protos.ListKeysForTypeResponse{}
	if err := config_protos.ValidateListKeysForTypeRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	keys, err := service.store.ListKeysForType(req.GetNetworkId(), req.GetType())
	if err != nil {
		msgFormat := "Error retrieving keys for type: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	// Return keys in a deterministic order
	sort.Strings(keys)
	ret.Keys = keys
	return ret, nil
}

func (service *ConfigService) CreateConfig(context context.Context, req *config_protos.CreateOrUpdateConfigRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if err := config_protos.ValidateCreateOrUpdateConfigRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	err := service.store.CreateConfig(req.GetNetworkId(), req.GetType(), req.GetKey(), req.GetValue())
	if err != nil {
		msgFormat := "Error creating config: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return ret, nil
}

func (service *ConfigService) UpdateConfig(context context.Context, req *config_protos.CreateOrUpdateConfigRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if err := config_protos.ValidateCreateOrUpdateConfigRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	err := service.store.UpdateConfig(req.GetNetworkId(), req.GetType(), req.GetKey(), req.GetValue())
	if err != nil {
		msgFormat := "Error updating config: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return ret, nil
}

func (service *ConfigService) DeleteConfig(context context.Context, req *config_protos.GetOrDeleteConfigRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if err := config_protos.ValidateGetOrDeleteConfigRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	err := service.store.DeleteConfig(req.GetNetworkId(), req.GetType(), req.GetKey())
	if err != nil {
		msgFormat := "Error deleting config: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return ret, nil
}

func (service *ConfigService) DeleteConfigs(context context.Context, req *config_protos.GetOrDeleteConfigsRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if err := config_protos.ValidateGetOrDeleteConfigsRequest(req); err != nil {
		return ret, protos.NewGrpcValidationError(err)
	}

	err := service.store.DeleteConfigs(req.GetNetworkId(), protoFilterToStorageFilter(req.GetFilter()))
	if err != nil {
		msgFormat := "Error deleting configs: %s"
		glog.Errorf(msgFormat, err)
		return ret, status.Errorf(codes.Aborted, msgFormat, err)
	}
	return ret, nil
}

func (service *ConfigService) DeleteAllConfigsForNetwork(context context.Context, req *config_protos.NetworkIdRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if req.GetNetworkId() == "" {
		return ret, protos.NewGrpcValidationError(errors.New("Network ID cannot be empty"))
	}

	err := service.store.DeleteConfigsForNetwork(req.GetNetworkId())
	if err != nil {
		msgFormat := "Error deleting configs for network %s: %s"
		glog.Errorf(msgFormat, req.GetNetworkId(), err)
		return ret, status.Errorf(codes.Aborted, msgFormat, req.GetNetworkId(), err)
	}
	return ret, nil
}

func protoFilterToStorageFilter(in *config_protos.ConfigFilter) *storage.FilterCriteria {
	return &storage.FilterCriteria{
		Type: in.GetType(),
		Key:  in.GetKey(),
	}
}

func configValueMapToGetConfigsResponse(in map[mstore.TypeAndKey]*storage.ConfigValue) *config_protos.GetConfigsResponse {
	ret := make([]*config_protos.Config, 0, len(in))
	for typeAndKey, val := range in {
		ret = append(ret, &config_protos.Config{Type: typeAndKey.Type, Key: typeAndKey.Key, Value: val.Value})
	}
	return &config_protos.GetConfigsResponse{Configs: ret}
}
