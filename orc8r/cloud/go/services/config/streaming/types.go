/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package streaming

import (
	"fmt"
	"time"

	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/protos/mconfig"
	"magma/orc8r/cloud/go/services/config/blacklist"
	"magma/orc8r/cloud/go/services/config/registry"
	"magma/orc8r/cloud/go/services/config/streaming/storage"
	"magma/orc8r/cloud/go/services/magmad"
	upgradeprotos "magma/orc8r/cloud/go/services/upgrade/protos"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/any"
)

type ChangeOperation string

const CreateOperation ChangeOperation = "c"
const DeleteOperation ChangeOperation = "d"
const ReadOperation ChangeOperation = "r"
const UpdateOperation ChangeOperation = "u"

// Wraps the methods of a kafka stream consumer that we care about for easy
// test mocking
type StreamConsumer interface {
	Commit() ([]kafka.TopicPartition, error)
	Close() error
	ReadMessage(timeout time.Duration) (*kafka.Message, error)
	SubscribeTopics(topics []string, cb kafka.RebalanceCb) error
}

type StreamConsumerFactory func() (StreamConsumer, error)

// Interface type for a stream update that can be applied to stored mconfigs
type ApplyableUpdate interface {
	// Applies this update to stored mconfigs given the offset of the stream
	// element that this update was computed from
	Apply(store storage.MconfigStorage, offset int64) error
}

// The following types implement ApplyableUpdate

// Wraps an update to a config
type ConfigUpdate struct {
	Operation  ChangeOperation
	NetworkId  string
	ConfigType string
	ConfigKey  string
	NewValue   interface{}
}

// Wraps the parts of an update to the gateway records table that we care about
type GatewayUpdate struct {
	Operation ChangeOperation
	NetworkId string
	GatewayId string
}

// Wraps an update to a tier
type TierUpdate struct {
	Operation   ChangeOperation
	NetworkId   string
	TierId      string
	TierVersion string
	TierImages  []*upgradeprotos.ImageSpec
}

// A no-op update from the stream. Debezium will for example send a nil record
// to a table's stream whenever we delete a row for whatever reason. That
// record should be consumed and no-op'd.
type NoOpUpdate struct{}

func (*NoOpUpdate) Apply(store storage.MconfigStorage, offset int64) error {
	return nil
}

func (configUpdate *ConfigUpdate) Apply(store storage.MconfigStorage, offset int64) error {
	if blacklist.IsConfigBlacklisted(configUpdate.ConfigType) {
		glog.Warningf("Ignoring stream record for blacklisted config type %s", configUpdate.ConfigType)
		return nil
	}

	oldMconfigs, err := configUpdate.getOldMconfigs(store)
	if err != nil {
		return err
	}

	oldMconfigsInput := getMconfigsByGatewayId(oldMconfigs)
	updatedMconfigsByGatewayId, err := ApplyMconfigUpdate(configUpdate, oldMconfigsInput)
	if err != nil {
		return fmt.Errorf("Error applying mconfig update: %s", err)
	}
	updateCriteria := createUpdateCriteria(updatedMconfigsByGatewayId, offset)

	return store.CreateOrUpdateMconfigs(configUpdate.NetworkId, updateCriteria)
}

func (configUpdate *ConfigUpdate) getOldMconfigs(store storage.MconfigStorage) (map[string]*storage.StoredMconfig, error) {
	relatedGatewayIds, err := registry.GetGatewayIdsForConfig(configUpdate.ConfigType, configUpdate.NetworkId, configUpdate.ConfigKey)
	if err != nil {
		return map[string]*storage.StoredMconfig{}, fmt.Errorf("Could not get gateway IDs for config update: %s", err)
	}

	mconfigsById, err := store.GetMconfigs(configUpdate.NetworkId, relatedGatewayIds)
	if err != nil {
		return map[string]*storage.StoredMconfig{}, err
	}

	// Insert empty configs for any gateways that didn't have views already
	for _, gwId := range relatedGatewayIds {
		_, exists := mconfigsById[gwId]
		if !exists {
			mconfigsById[gwId] = &storage.StoredMconfig{
				NetworkId: configUpdate.NetworkId,
				GatewayId: gwId,
				Mconfig:   &protos.GatewayConfigs{ConfigsByKey: map[string]*any.Any{}},
				Offset:    0,
			}
		}
	}
	return mconfigsById, nil
}

func (gatewayUpdate *GatewayUpdate) Apply(store storage.MconfigStorage, offset int64) error {
	networkId := gatewayUpdate.NetworkId
	gatewayId := gatewayUpdate.GatewayId

	switch gatewayUpdate.Operation {
	case ReadOperation, UpdateOperation:
		// No-op for read and update, we only care about create/delete events
		return nil
	case CreateOperation:
		newMconfig, err := GetMconfigForNewGateway(networkId, gatewayId)
		if err != nil {
			return err
		}

		// Offset of stored mconfig only refers to the offset in the config stream,
		// so just persist -1
		update := &storage.MconfigUpdateCriteria{GatewayId: gatewayId, NewMconfig: newMconfig, Offset: -1}
		return store.CreateOrUpdateMconfigs(networkId, []*storage.MconfigUpdateCriteria{update})
	case DeleteOperation:
		return store.DeleteMconfigs(networkId, []string{gatewayId})
	default:
		return fmt.Errorf("Unrecognized stream operation %s", gatewayUpdate.Operation)
	}
}

// IMPORTANT: eventually, we will want to pull the tier part of upgrade service
// into config service to make this logic a whole lot simpler. Also it would
// give us sequential consistency in upgrade and config applications.

// We should also probably add a reverse-mapping between gateways and tiers per
// network to the state store for this stream processor. That would make
// tier update applications even faster and simpler.

// This function queries magmad for all gateways in the tier's network, then
// grabs all those gateways' mconfigs from the state store. For each mconfig,
// if the tier ID in magmad config matches this tier's ID, update the version
// and images and write back the mconfig.
func (tierUpdate *TierUpdate) Apply(store storage.MconfigStorage, offset int64) error {
	gatewayIds, err := magmad.ListGateways(tierUpdate.NetworkId)
	if err != nil {
		return err
	}
	storedMconfigs, err := store.GetMconfigs(tierUpdate.NetworkId, gatewayIds)
	if err != nil {
		return err
	}

	switch tierUpdate.Operation {
	case DeleteOperation:
		// If a tier is deleted, propagate 0.0.0-0 and no images
		updates, err := getUpdatesForTierUpdate(storedMconfigs, tierUpdate.TierId, "0.0.0-0", []*upgradeprotos.ImageSpec{})
		if err != nil {
			return err
		}
		return store.CreateOrUpdateMconfigs(tierUpdate.NetworkId, updates)
	case CreateOperation, UpdateOperation:
		updates, err := getUpdatesForTierUpdate(storedMconfigs, tierUpdate.TierId, tierUpdate.TierVersion, tierUpdate.TierImages)
		if err != nil {
			return err
		}
		return store.CreateOrUpdateMconfigs(tierUpdate.NetworkId, updates)
	case ReadOperation:
		// no-op for a snapshot read
		return nil
	default:
		return fmt.Errorf("Unrecognized stream operation %s", tierUpdate.Operation)
	}
}

func getUpdatesForTierUpdate(
	storedMconfigs map[string]*storage.StoredMconfig,
	tierId string,
	newVersion string,
	newImages []*upgradeprotos.ImageSpec,
) ([]*storage.MconfigUpdateCriteria, error) {
	ret := make([]*storage.MconfigUpdateCriteria, 0, len(storedMconfigs))
	for _, storedMconfig := range storedMconfigs {
		updated, err := updateVersionAndImagesIfNecessary(tierId, newVersion, newImages, storedMconfig.Mconfig)
		if err != nil {
			return []*storage.MconfigUpdateCriteria{}, fmt.Errorf("Error updating mconfigs for tier update: %s", err)
		}

		if updated {
			ret = append(
				ret,
				&storage.MconfigUpdateCriteria{
					GatewayId:  storedMconfig.GatewayId,
					NewMconfig: storedMconfig.Mconfig,
					Offset:     storedMconfig.Offset,
				},
			)
		}
	}
	return ret, nil
}

func updateVersionAndImagesIfNecessary(tierId string, newVersion string, newImages []*upgradeprotos.ImageSpec, mconfigOut *protos.GatewayConfigs) (bool, error) {
	magmadAny, exists := mconfigOut.ConfigsByKey["magmad"]
	if !exists {
		// No magmad config, no-op
		return false, nil
	}

	magmadConfig := &mconfig.MagmaD{}
	err := ptypes.UnmarshalAny(magmadAny, magmadConfig)
	if err != nil {
		return false, err
	}

	if magmadConfig.TierId == tierId {
		magmadConfig.PackageVersion = newVersion
		magmadConfig.Images = imagesProtosToMconfigImages(newImages)
		newAny, err := ptypes.MarshalAny(magmadConfig)
		if err != nil {
			return false, err
		}
		mconfigOut.ConfigsByKey["magmad"] = newAny
		return true, nil
	} else {
		// Not the tier we're interested in, no-op
		return false, nil
	}
}

func imagesProtosToMconfigImages(images []*upgradeprotos.ImageSpec) []*mconfig.ImageSpec {
	ret := make([]*mconfig.ImageSpec, 0, len(images))
	for _, img := range images {
		mconfigImg := &mconfig.ImageSpec{}
		protos.FillIn(img, mconfigImg)
		ret = append(ret, mconfigImg)
	}
	return ret
}

func getMconfigsByGatewayId(oldStoredMconfigs map[string]*storage.StoredMconfig) map[string]*protos.GatewayConfigs {
	ret := map[string]*protos.GatewayConfigs{}
	for k, storedMconfig := range oldStoredMconfigs {
		ret[k] = storedMconfig.Mconfig
	}
	return ret
}

func createUpdateCriteria(
	updatedMconfigsByGatewayId map[string]*protos.GatewayConfigs,
	offset int64,
) []*storage.MconfigUpdateCriteria {
	ret := make([]*storage.MconfigUpdateCriteria, 0, len(updatedMconfigsByGatewayId))
	for gatewayId, mcfg := range updatedMconfigsByGatewayId {
		ret = append(ret, &storage.MconfigUpdateCriteria{GatewayId: gatewayId, NewMconfig: mcfg, Offset: offset})
	}
	return ret
}
