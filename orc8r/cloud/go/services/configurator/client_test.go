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
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"

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
	config := map[string]interface{}{
		"foo": "world",
	}
	// Create, Load
	network1 := configurator.Network{
		ID:          networkID1,
		Name:        "test_network",
		Description: "description",
		Configs:     config,
	}
	_, err = configurator.CreateNetworks([]configurator.Network{network1})
	assert.NoError(t, err)

	networks, notFound, err := configurator.LoadNetworks([]string{networkID1}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	assert.Equal(t, 1, len(networks))

	// Update, Load
	newDesc := "Should be updated now"
	toAddOrUpdate := map[string]interface{}{}
	toAddOrUpdate["bar"] = "hello"
	toDelete := []string{"foo"}
	updateCriteria1 := configurator.NetworkUpdateCriteria{
		ID:                   networkID1,
		NewDescription:       &newDesc,
		ConfigsToAddOrUpdate: toAddOrUpdate,
		ConfigsToDelete:      toDelete,
	}

	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria1})
	assert.NoError(t, err)
	networks, notFound, err = configurator.LoadNetworks([]string{networkID1}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	assert.Equal(t, 1, len(networks))
	assert.Equal(t, newDesc, networks[0].Description)
	_, fooPresent := networks[0].Configs["foo"]
	assert.False(t, fooPresent)
	assert.Equal(t, "hello", networks[0].Configs["bar"])

	// Create, Load
	network2 := configurator.Network{
		ID:          networkID2,
		Name:        "test_network2",
		Description: "description2",
	}
	_, err = configurator.CreateNetworks([]configurator.Network{network2})
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
	foobarID := storage.TypeAndKey{Type: "foo", Key: "bar"}
	foobarEnt := configurator.NetworkEntity{
		Type:        "foo",
		Key:         "bar",
		Name:        "foobar",
		Description: "ent: foobar",
		PhysicalID:  "1234",
		Config:      "hello",
	}
	foobooID := storage.TypeAndKey{Type: "foo", Key: "boo"}
	foobooEnt := configurator.NetworkEntity{
		Type:        "foo",
		Key:         "boo",
		Name:        "fooboo",
		Description: "ent: fooboo",
		PhysicalID:  "5678",
		Config:      "bye",
	}
	fullEntityLoad := configurator.EntityLoadCriteria{
		LoadMetadata:       true,
		LoadAssocsToThis:   true,
		LoadAssocsFromThis: true,
		LoadConfig:         true,
	}

	// Create, Load
	_, err = configurator.CreateEntities(networkID1, []configurator.NetworkEntity{foobarEnt, foobooEnt})
	assert.NoError(t, err)

	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID1,
		nil, nil, nil,
		[]storage.TypeAndKey{foobarID, foobooID},
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

	// Update, Load add an association from foobar to fooboo
	newPhysID := "4321"
	entityUpdateCriteria := configurator.EntityUpdateCriteria{
		Type:              foobarID.Type,
		Key:               foobarID.Key,
		NewPhysicalID:     &newPhysID,
		AssociationsToAdd: []storage.TypeAndKey{foobooID},
	}

	_, err = configurator.UpdateEntities(networkID1, []configurator.EntityUpdateCriteria{entityUpdateCriteria})
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil, nil, nil,
		fullEntityLoad,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)
	assert.Equal(t, "4321", entities[0].PhysicalID)
	assert.Equal(t, 1, len(entities[0].Associations))
	assert.Equal(t, foobooID.Type, entities[0].Associations[0].Type)
	assert.Equal(t, foobooID.Key, entities[0].Associations[0].Key)
	assert.Equal(t, foobarID.Key, entities[1].ParentAssociations[0].Key)

	// Delete, Load
	err = configurator.DeleteEntities(networkID1, []storage.TypeAndKey{foobooID})
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil, nil, nil,
		fullEntityLoad,
	)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)

	// Create an entitiy in network2 and load all entities
	foostartEnt := configurator.NetworkEntity{
		Type:        "foo",
		Key:         "star",
		Name:        "foostar",
		Description: "ent: foostar",
		PhysicalID:  "5678",
		Config:      "bye",
	}
	_, err = configurator.CreateEntity(networkID2, foostartEnt)
	assert.NoError(t, err)

	entities, err = configurator.LoadAllEntities(strPointer("foo"), nil, nil, nil, fullEntityLoad)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, foobarEnt.Name, entities[0].Name)
	assert.Equal(t, networkID1, entities[0].NetworkID)
	assert.Equal(t, foostartEnt.Name, entities[1].Name)
	assert.Equal(t, networkID2, entities[1].NetworkID)
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
