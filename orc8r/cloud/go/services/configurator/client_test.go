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
	"magma/orc8r/cloud/go/services/configurator/protos"
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
	network1 := &protos.Network{
		Id:          networkID1,
		Name:        "test_network",
		Description: "description",
		Configs:     config,
	}
	_, err = configurator.CreateNetworks([]*protos.Network{network1})
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
	updateCriteria1 := &protos.NetworkUpdateCriteria{
		Id:                   networkID1,
		NewDescription:       strToStringValue(newDesc),
		ConfigsToAddOrUpdate: toAddOrUpdate,
		ConfigsToDelete:      toDelete,
	}

	err = configurator.UpdateNetworks([]*protos.NetworkUpdateCriteria{updateCriteria1})
	networks, notFound, err = configurator.LoadNetworks([]string{networkID1}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(notFound))
	assert.Equal(t, 1, len(networks))
	assert.Equal(t, newDesc, networks[networkID1].Description)
	_, fooPresent := networks[networkID1].Configs["foo"]
	assert.False(t, fooPresent)
	assert.Equal(t, []byte("hello"), networks[networkID1].Configs["bar"])

	// Create, Load
	network2 := &protos.Network{
		Id:          networkID2,
		Name:        "test_network2",
		Description: "description2",
	}
	_, err = configurator.CreateNetworks([]*protos.Network{network2})
	assert.NoError(t, err)

	networkIDs, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(networkIDs))
	assert.Equal(t, networkID2, networkIDs[1])

	// Delete, Load
	err = configurator.DeleteNetworks([]string{network2.Id})
	assert.NoError(t, err)

	networks, notFound, err = configurator.LoadNetworks([]string{networkID2}, true, true)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(networks))
	assert.Equal(t, 1, len(notFound))

	// Test Basic Entity Interface
	entityID1 := &protos.EntityID{Type: "foo", Id: "bar"}
	entity1 := &protos.NetworkEntity{
		Type:        "foo",
		Id:          "bar",
		Name:        "foobar",
		Description: "ent: foobar",
		PhysicalId:  "1234",
		Config:      []byte("hello"),
	}
	entityID2 := &protos.EntityID{Type: "foo", Id: "boo"}
	entity2 := &protos.NetworkEntity{
		Type:        "foo",
		Id:          "boo",
		Name:        "fooboo",
		Description: "ent: fooboo",
		PhysicalId:  "5678",
		Config:      []byte("bye"),
	}
	fullEntityLoad := &protos.EntityLoadCriteria{
		LoadMetadata:    true,
		LoadAssocsTo:    true,
		LoadAssocsFrom:  true,
		LoadConfig:      true,
		LoadPermissions: true,
	}

	// Create, Load
	_, err = configurator.CreateEntities(networkID1, []*protos.NetworkEntity{entity1, entity2})
	assert.NoError(t, err)

	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID1,
		nil,
		nil,
		[]*protos.EntityID{entityID1, entityID2},
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
	entityUpdateCriteria := &protos.EntityUpdateCriteria{
		Type:              entityID1.Type,
		Key:               entityID1.Id,
		NewPhysicalID:     strToStringValue("4321"),
		AssociationsToAdd: []*protos.EntityID{entityID2},
	}

	_, err = configurator.UpdateEntities(networkID1, []*protos.EntityUpdateCriteria{entityUpdateCriteria})
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
	assert.Equal(t, "4321", entities[0].PhysicalId)
	assert.Equal(t, 1, len(entities[0].Assocs))
	assert.Equal(t, entityID2.Id, entities[0].Assocs[0].Id)

	// Delete, Load
	err = configurator.DeleteEntities(networkID1, []*protos.EntityID{entityID2})
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
