/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"encoding/base64"
	"fmt"
	"reflect"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config/blacklist"
	"magma/orc8r/cloud/go/services/config/registry"
	magmadprotos "magma/orc8r/cloud/go/services/magmad/protos"
	"magma/orc8r/cloud/go/services/materializer/gateways/storage"
)

// KafkaGatewayUpdate is the object which is serialized into JSON and pushed onto the common gatewaystate topic stream
type KafkaGatewayUpdate struct {
	// UpdateType is the type of this update: either "gwstatus", "gatewayrecords", or "configurations"
	UpdateType UpdateType
	// Operation is the operation type of this update: either "c", "r", "u", or "d",
	// for create, read, update, and delete, respectively
	Operation string
	// NetworkID is the ID of the network this update applies to
	NetworkID string
	// The payload object of the update
	Payload UpdatePayload
}

// ApplyUpdate applies the update specified by this KafkaGatewayUpdate to the store
func (update *KafkaGatewayUpdate) ApplyUpdate(store storage.GatewayViewStorage, offset int64) error {
	return update.Payload.Apply(update.Operation, update.NetworkID, offset, store)
}

type GatewayConfigUpdate struct {
	ConfigKey   string
	ConfigType  string
	ConfigBytes string
}

func (configUpdate *GatewayConfigUpdate) Apply(
	operation string,
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	switch operation {
	case "c", "r", "u":
		return configUpdate.applyUpdateOrCreate(networkID, offset, store)
	case "d":
		return configUpdate.applyDelete(networkID, offset, store)
	default:
		return fmt.Errorf("Unrecognized operation: %s", operation)
	}
}

func (configUpdate *GatewayConfigUpdate) applyUpdateOrCreate(
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	configType := configUpdate.ConfigType
	if blacklist.IsConfigBlacklisted(configType) {
		return nil
	}
	configBytes, err := getBytesFromBase64(configUpdate.ConfigBytes)
	if err != nil {
		return fmt.Errorf("Error decoding config bytes: %s", err)
	}
	configObj, err := registry.UnmarshalConfig(configType, configBytes)
	if err != nil {
		return err
	}

	gatewayIDs, err := registry.GetGatewayIdsForConfig(configType, networkID, configUpdate.ConfigKey)
	if err != nil {
		return err
	}
	// Only apply gateway-level configs
	if len(gatewayIDs) > 1 {
		return nil
	}

	updateParam := &storage.GatewayUpdateParams{
		NewConfig: map[string]interface{}{
			configType: configObj,
		},
		Offset: offset,
	}

	updateParams := make(map[string]*storage.GatewayUpdateParams)
	for _, gatewayID := range gatewayIDs {
		updateParams[gatewayID] = updateParam
	}
	return store.UpdateOrCreateGatewayViews(networkID, updateParams)
}

func (configUpdate *GatewayConfigUpdate) applyDelete(
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	updateParam := &storage.GatewayUpdateParams{
		ConfigsToDelete: []string{configUpdate.ConfigType},
		Offset:          offset,
	}
	gatewayIDs, err := registry.GetGatewayIdsForConfig(configUpdate.ConfigType, networkID, configUpdate.ConfigKey)
	if err != nil {
		return err
	}
	updateParams := make(map[string]*storage.GatewayUpdateParams)
	for _, gatewayID := range gatewayIDs {
		updateParams[gatewayID] = updateParam
	}
	return store.UpdateOrCreateGatewayViews(networkID, updateParams)
}

type GatewayStatusUpdate struct {
	GatewayID   string
	StatusBytes string
}

func (statusUpdate *GatewayStatusUpdate) Apply(
	operation string,
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	switch operation {
	case "c", "r", "u":
		return statusUpdate.applyUpdateOrCreate(networkID, offset, store)
	case "d":
		return nil
	default:
		return fmt.Errorf("Unrecognized operation: %s", operation)
	}
}

func (statusUpdate *GatewayStatusUpdate) applyUpdateOrCreate(
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	statusBytes, err := getBytesFromBase64(statusUpdate.StatusBytes)
	if err != nil {
		return fmt.Errorf("Error decoding status bytes: %s", err)
	}
	status := &protos.GatewayStatus{}
	err = protos.Unmarshal(statusBytes, status)
	if err != nil {
		return err
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		statusUpdate.GatewayID: &storage.GatewayUpdateParams{
			NewStatus: status,
			Offset:    offset,
		},
	}
	return store.UpdateOrCreateGatewayViews(networkID, updateParams)
}

type GatewayRecordUpdate struct {
	GatewayID   string
	RecordBytes string
}

func (recordUpdate *GatewayRecordUpdate) Apply(
	operation string,
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	switch operation {
	case "c", "r", "u":
		return recordUpdate.applyUpdateOrCreate(networkID, offset, store)
	case "d":
		return recordUpdate.applyDelete(networkID, offset, store)
	default:
		return fmt.Errorf("Unrecognized operation: %s", operation)
	}
}

func (recordUpdate *GatewayRecordUpdate) applyUpdateOrCreate(
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	recordBytes, err := getBytesFromBase64(recordUpdate.RecordBytes)
	if err != nil {
		return fmt.Errorf("Error decoding record bytes: %s", err)
	}
	record := &magmadprotos.AccessGatewayRecord{}
	err = protos.Unmarshal(recordBytes, record)
	if err != nil {
		return err
	}
	updateParams := map[string]*storage.GatewayUpdateParams{
		recordUpdate.GatewayID: {
			NewRecord: record,
			Offset:    offset,
		},
	}
	return store.UpdateOrCreateGatewayViews(networkID, updateParams)
}

func getBytesFromBase64(base64Str string) ([]byte, error) {
	base64Bytes := []byte(base64Str)
	decoded := make([]byte, base64.StdEncoding.DecodedLen(len(base64Bytes)))
	numBytes, err := base64.StdEncoding.Decode(decoded, base64Bytes)
	if err != nil {
		return nil, fmt.Errorf("Error decoding bytes from base 64: %s", err)
	}
	return decoded[:numBytes], nil
}

func (recordUpdate *GatewayRecordUpdate) applyDelete(
	networkID string,
	offset int64,
	store storage.GatewayViewStorage,
) error {
	return store.DeleteGatewayViews(networkID, []string{recordUpdate.GatewayID})
}

// Type registry

var updateTypeRegistry = map[UpdateType]reflect.Type{
	Configurations: reflect.TypeOf(GatewayConfigUpdate{}),
	Statuses:       reflect.TypeOf(GatewayStatusUpdate{}),
	Records:        reflect.TypeOf(GatewayRecordUpdate{}),
}
