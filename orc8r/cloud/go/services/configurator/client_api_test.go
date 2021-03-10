/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package configurator_test

import (
	"fmt"
	"testing"

	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/test_init"
	"magma/orc8r/cloud/go/storage"

	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/assert"
)

const (
	networkID1 = "network_id_1"
	networkID2 = "network_id_2"
)

func TestConfiguratorService(t *testing.T) {
	test_init.StartTestService(t)

	networkSerdes := serde.NewRegistry(
		&mockSerde{domain: configurator.NetworkConfigSerdeDomain, serdeType: "foo"},
		&mockSerde{domain: configurator.NetworkConfigSerdeDomain, serdeType: "bar"},
	)
	entitySerdes := serde.NewRegistry(
		&mockSerde{domain: configurator.NetworkEntitySerdeDomain, serdeType: "foo"},
		&mockSerde{domain: configurator.NetworkEntitySerdeDomain, serdeType: "bar"},
	)

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
	_, err := configurator.CreateNetworks([]configurator.Network{network1}, networkSerdes)
	assert.NoError(t, err)

	networks, notFound, err := configurator.LoadNetworks([]string{networkID1}, true, true, networkSerdes)
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

	err = configurator.UpdateNetworks([]configurator.NetworkUpdateCriteria{updateCriteria1}, networkSerdes)
	assert.NoError(t, err)
	networks, notFound, err = configurator.LoadNetworks([]string{networkID1}, true, true, networkSerdes)
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
	_, err = configurator.CreateNetworks([]configurator.Network{network2}, networkSerdes)
	assert.NoError(t, err)

	networkIDs, err := configurator.ListNetworkIDs()
	assert.NoError(t, err)
	assert.Equal(t, 2, len(networkIDs))
	assert.Equal(t, networkID2, networkIDs[1])

	// Delete, Load
	err = configurator.DeleteNetworks([]string{network2.ID})
	assert.NoError(t, err)

	networks, notFound, err = configurator.LoadNetworks([]string{networkID2}, true, true, networkSerdes)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(networks))
	assert.Equal(t, 1, len(notFound))

	// Create Networks With Type
	createdTypedLteNetworks, err := configurator.CreateNetworks(
		[]configurator.Network{
			{
				Name: "lte network 1",
				Type: "lte",
				ID:   "test_network3",
			},
			{
				Name: "lte network 2",
				Type: "lte",
				ID:   "test_network4",
			},
		},
		networkSerdes,
	)
	assert.NoError(t, err)

	createdTypedLteNetworks[0].Name = ""
	createdTypedLteNetworks[1].Name = ""

	networks, err = configurator.LoadNetworksOfType("lte", false, false, networkSerdes)
	assert.NoError(t, err)
	assert.Equal(t, createdTypedLteNetworks, networks)

	// Test Basic Entity Interface
	entityID1 := storage.TypeAndKey{Type: "foo", Key: "bar"}
	entity1 := configurator.NetworkEntity{
		Type:        "foo",
		Key:         "bar",
		Name:        "foobar",
		Description: "ent: foobar",
		PhysicalID:  "1234",
		Config:      "hello",
	}
	entityID2 := storage.TypeAndKey{Type: "foo", Key: "boo"}
	entity2 := configurator.NetworkEntity{
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
	_, err = configurator.CreateEntities(networkID1, []configurator.NetworkEntity{entity1, entity2}, entitySerdes)
	assert.NoError(t, err)

	entities, entitiesNotFound, err := configurator.LoadEntities(
		networkID1,
		nil, nil, nil,
		[]storage.TypeAndKey{entityID1, entityID2},
		fullEntityLoad,
		entitySerdes,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)

	// LoadAllPerType
	entities, _, err = configurator.LoadAllEntitiesOfType(networkID1, "foo", fullEntityLoad, entitySerdes)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, "foobar", entities[0].Name)
	assert.Equal(t, "fooboo", entities[1].Name)

	// Load paginated entities
	paginatedEntityLoad := fullEntityLoad
	paginatedEntityLoad.PageSize = 1
	paginatedEntities, nextPageToken, err := configurator.LoadAllEntitiesOfType(networkID1, "foo", paginatedEntityLoad, entitySerdes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(paginatedEntities))
	assert.Equal(t, "foobar", paginatedEntities[0].Name)
	assert.Equal(t, "CgNiYXI=", nextPageToken)

	paginatedEntityLoad.PageToken = nextPageToken
	paginatedEntities, nextPageToken, err = configurator.LoadAllEntitiesOfType(networkID1, "foo", paginatedEntityLoad, entitySerdes)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(paginatedEntities))
	assert.Equal(t, "fooboo", paginatedEntities[0].Name)
	assert.Equal(t, "CgNib28=", nextPageToken)

	// Update, Load add an association from foobar to fooboo
	newPhysID := "4321"
	entityUpdateCriteria := configurator.EntityUpdateCriteria{
		Type:              entityID1.Type,
		Key:               entityID1.Key,
		NewPhysicalID:     &newPhysID,
		AssociationsToAdd: []storage.TypeAndKey{entityID2},
	}

	_, err = configurator.UpdateEntities(networkID1, []configurator.EntityUpdateCriteria{entityUpdateCriteria}, entitySerdes)
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil, nil, nil,
		fullEntityLoad,
		entitySerdes,
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
	assert.Equal(t, entityID1.Key, entities[1].ParentAssociations[0].Key)

	// Update foobar, create foobaz, add association fooboo -> foobaz  in 1
	// client call
	err = configurator.WriteEntities(
		networkID1,
		[]configurator.EntityWriteOperation{
			configurator.EntityUpdateCriteria{Type: entityID1.Type, Key: entityID1.Key, NewDescription: swag.String("newnewnew")},
			configurator.NetworkEntity{Type: "foo", Key: "baz"},
			configurator.EntityUpdateCriteria{
				Type: entityID2.Type, Key: entityID2.Key,
				AssociationsToAdd: []storage.TypeAndKey{{Type: "foo", Key: "baz"}},
			},
		},
		entitySerdes,
	)
	assert.NoError(t, err)

	entities, _, err = configurator.LoadEntities(
		networkID1,
		swag.String("foo"), nil,
		nil, nil,
		fullEntityLoad,
		entitySerdes,
	)
	assert.NoError(t, err)
	expected := configurator.NetworkEntities{
		{
			NetworkID: networkID1, Type: entityID1.Type, Key: entityID1.Key,
			Name: "foobar", Description: "newnewnew",
			PhysicalID:   "4321",
			Config:       "hello",
			GraphID:      "2",
			Associations: []storage.TypeAndKey{entityID2},
			Version:      2,
		},
		{
			NetworkID: networkID1, Type: "foo", Key: "baz",
			GraphID:            "2",
			ParentAssociations: []storage.TypeAndKey{entityID2},
		},
		{
			NetworkID: networkID1, Type: entityID2.Type, Key: entityID2.Key,
			Name: "fooboo", Description: "ent: fooboo",
			PhysicalID:         "5678",
			Config:             "bye",
			GraphID:            "2",
			Associations:       []storage.TypeAndKey{{Type: "foo", Key: "baz"}},
			ParentAssociations: []storage.TypeAndKey{entityID1},
			Version:            1,
		},
	}
	assert.Equal(t, expected, entities)

	// Delete, Load
	err = configurator.DeleteEntities(networkID1, []storage.TypeAndKey{entityID2})
	assert.NoError(t, err)
	entities, entitiesNotFound, err = configurator.LoadEntities(
		networkID1,
		strPointer("foo"),
		nil, nil, nil,
		fullEntityLoad,
		entitySerdes,
	)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))
	assert.Equal(t, 0, len(entitiesNotFound))
	assert.Equal(t, "foobar", entities[0].Name)
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
