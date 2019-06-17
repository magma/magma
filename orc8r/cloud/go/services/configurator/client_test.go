/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/
package configurator_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/storage"
	"magma/orc8r/cloud/go/services/configurator/test_init"

	"github.com/golang/protobuf/ptypes/wrappers"
	"github.com/stretchr/testify/assert"
)

const (
	networkID1 = "network_id_1"
	networkID2 = "network_id_2"
)

func TestConfiguratorService(t *testing.T) {
	test_init.StartTestService(t)
	err := serde.RegisterSerdes(&mockSerde{domain: configurator.NetworkConfigSerdeDomain, serdeType: "foo"})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&mockSerde{domain: configurator.NetworkEntitySerdeDomain, serdeType: "foo"})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&mockSerde{domain: configurator.NetworkConfigSerdeDomain, serdeType: "bar"})
	assert.NoError(t, err)
	err = serde.RegisterSerdes(&mockSerde{domain: configurator.NetworkEntitySerdeDomain, serdeType: "bar"})
	assert.NoError(t, err)

	// Test Basic Network Interface
	config := map[string][]byte{
		"foo": []byte("world"),
	}
	// Create, Load
	network1 := &storage.Network{
		ID:          networkID1,
		Name:        "test_network",
		Description: "description",
		Configs:     config,
	}
	_, err = configurator.CreateNetworks([]*storage.Network{network1})
	assert.NoError(t, err)

	networks, notFound, err := configurator.LoadNetworks([]string{networkID1}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	assert.Equal(t, 1, len(networks))

	// Update, Load
	newDesc := "Should be updated now"
	toAddOrUpdate := map[string][]byte{}
	toAddOrUpdate["bar"] = []byte("hello")
	toDelete := []string{"foo"}
	updateCriteria1 := &storage.NetworkUpdateCriteria{
		ID:                   networkID1,
		NewDescription:       strToStringValue(newDesc),
		ConfigsToAddOrUpdate: toAddOrUpdate,
		ConfigsToDelete:      toDelete,
	}

	err = configurator.UpdateNetworks([]*storage.NetworkUpdateCriteria{updateCriteria1})
	networks, notFound, err = configurator.LoadNetworks([]string{networkID1}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	assert.Equal(t, 1, len(networks))
	assert.Equal(t, newDesc, networks[0].Description)
	_, fooPresent := networks[0].Configs["foo"]
	assert.False(t, fooPresent)
	assert.Equal(t, []byte("hello"), networks[0].Configs["bar"])

	// Create, Load
	network2 := &storage.Network{
		ID:          networkID2,
		Name:        "test_network2",
		Description: "description2",
	}
	_, err = configurator.CreateNetworks([]*storage.Network{network2})
	assert.NoError(t, err)

	networkIDs, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(networkIDs))
	assert.Equal(t, networkID2, networkIDs[1])

	// Delete, Load
	err = configurator.DeleteNetworks([]string{network2.ID})
	assert.NoError(t, err)

	networks, notFound, err = configurator.LoadNetworks([]string{networkID2}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(networks))
	assert.Equal(t, 1, len(notFound))

	// Test Basic Entity Interface
	entityID1 := &storage.EntityID{Type: "foo", Key: "bar"}
	entity1 := &storage.NetworkEntity{
		Type:        "foo",
		Key:         "bar",
		Name:        "foobar",
		Description: "ent: foobar",
		PhysicalID:  "1234",
		Config:      []byte("hello"),
	}
	entityID2 := &storage.EntityID{Type: "foo", Key: "boo"}
	entity2 := &storage.NetworkEntity{
		Type:        "foo",
		Key:         "boo",
		Name:        "fooboo",
		Description: "ent: fooboo",
		PhysicalID:  "5678",
		Config:      []byte("bye"),
	}
	fullEntityLoad := &storage.EntityLoadCriteria{
		LoadMetadata:       true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
		LoadConfig:         true,
		LoadPermissions:    true,
	}

	// Create, Load
	_, err = configurator.CreateEntities(networkID1, []*storage.NetworkEntity{entity1, entity2})
	assert.NoError(t, err)

	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID1,
		nil,
		nil,
		[]*storage.EntityID{entityID1, entityID2},
		fullEntityLoad,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)

	// LoadAllPerType
	entities, err = configurator.LoadAllEntitiesInNetwork(networkID1, "foo", fullEntityLoad)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)

	// Update, Load
	entityUpdateCriteria := &storage.EntityUpdateCriteria{
		Type:              entityID1.Type,
		Key:               entityID1.Key,
		NewPhysicalID:     strToStringValue("4321"),
		AssociationsToAdd: []*storage.EntityID{entityID2},
	}

	_, err = configurator.UpdateEntities(networkID1, []*storage.EntityUpdateCriteria{entityUpdateCriteria})
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil,
		nil,
		fullEntityLoad,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)
	assert.Equal(t, "4321", entities[0].PhysicalID)
	assert.Equal(t, 1, len(entities[0].Associations))
	assert.Equal(t, entityID2.Type, entities[0].Associations[0].Type)
	assert.Equal(t, entityID2.Key, entities[0].Associations[0].Key)

	// Delete, Load
	err = configurator.DeleteEntities(networkID1, []*storage.EntityID{entityID2})
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil,
		nil,
		fullEntityLoad,
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
}

func strToStringValue(str string) *wrappers.StringValue {
	return &wrappers.StringValue{Value: str}
}

func strPointer(str string) *string {
	return &str
}

type mockSerde struct {
	domain, serdeType string
}

func (m *mockSerde) GetDomain() string {
	return m.domain
}

func (m *mockSerde) GetType() string {
	return m.serdeType
}

func (m *mockSerde) Serialize(in interface{}) ([]byte, error) {
	str, ok := in.(string)
	if !ok {
		return nil, fmt.Errorf("serialization error")
	}
	return []byte(str), nil
}

func (m *mockSerde) Deserialize(in []byte) (interface{}, error) {
	return string(in), nil
}
