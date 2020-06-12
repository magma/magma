// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"sort"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddEquipmentTypesSameName(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "example_type_name",
	})
	require.NoError(t, err)
	assert.Equal(t, "example_type_name", equipmentType.Name)
	_, err = mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "example_type_name",
	})
	require.Error(t, err)
}

func TestQueryEquipmentTypes(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	for _, suffix := range []string{"a", "b"} {
		_, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
			Name:     "example_type_" + suffix,
			Category: pointer.ToString("example_type"),
		})
		require.NoError(t, err)
	}
	types, _ := qr.EquipmentTypes(ctx, nil, nil, nil, nil)
	require.Len(t, types.Edges, 2)

	var (
		names      = make([]string, len(types.Edges))
		categories = make([]*ent.EquipmentCategory, len(types.Edges))
	)
	for i, v := range types.Edges {
		names[i] = v.Node.Name
		category, err := v.Node.QueryCategory().Only(ctx)
		require.NoError(t, err)
		categories[i] = category
		require.Equal(t, "example_type", category.Name)
	}
	require.Len(t, categories, 2)
	assert.Equal(t, categories[0].ID, categories[1].ID)
	sort.Strings(names)
	assert.Equal(t, names, []string{"example_type_a", "example_type_b"})
}

func TestAddEquipmentTypeWithPositions(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	position1 := models.EquipmentPositionInput{
		Name: "Position 1",
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "equipment_type_name_1",
		Positions: []*models.EquipmentPositionInput{&position1},
	})
	require.NoError(t, err)
	fetchedNode, err := qr.Node(ctx, equipmentType.ID)
	require.NoError(t, err)
	fetchedEquipmentType, ok := fetchedNode.(*ent.EquipmentType)
	require.True(t, ok)

	require.Equal(t, equipmentType.ID, fetchedEquipmentType.ID, "Verifying saved equipment type vs fetched equipmenttype : ID")
	require.Equal(t, equipmentType.Name, fetchedEquipmentType.Name, "Verifying saved equipment type  vs fetched equipment type : Name")
	require.Equal(t, equipmentType.QueryPositionDefinitions().OnlyXID(ctx), fetchedEquipmentType.QueryPositionDefinitions().OnlyXID(ctx), "Verifying saved equipment type  vs fetched equipment type: position definition")
}

func TestAddEquipmentTypeWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr, etr := r.Mutation(), r.Query(), r.EquipmentType()
	extID := "12345"
	ptype := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        "string",
		Index:       pointer.ToInt(5),
		StringValue: pointer.ToString("Foo"),
		ExternalID:  &extID,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "example_type_a",
		Properties: []*models.PropertyTypeInput{&ptype},
	})
	require.NoError(t, err)

	fetchedNode, err := qr.Node(ctx, equipmentType.ID)
	require.NoError(t, err)
	fetchedEquipmentType, ok := fetchedNode.(*ent.EquipmentType)
	require.True(t, ok)
	fetchedPropertyTypes, _ := etr.PropertyTypes(ctx, fetchedEquipmentType)
	require.Len(t, fetchedPropertyTypes, 1)
	assert.Equal(t, fetchedPropertyTypes[0].Name, "str_prop")
	assert.Equal(t, fetchedPropertyTypes[0].Type, "string")
	assert.Equal(t, fetchedPropertyTypes[0].Index, 5)
	assert.Equal(t, fetchedPropertyTypes[0].ExternalID, extID)
}

func TestAddEquipmentTypeWithoutPositionNames(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type_name_1",
	})
	require.NoError(t, err)
	positions, err := equipmentType.QueryPositionDefinitions().All(ctx)
	require.NoError(t, err)
	assert.Len(t, positions, 0)
}

func TestAddEquipmentTypeWithPorts(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portDef := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: pointer.ToString("Eth1"),
		Bandwidth:    pointer.ToString("10/100/1000BASE-T"),
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "example_type_a",
		Ports: []*models.EquipmentPortInput{&portDef},
	})
	require.NoError(t, err)
	fetchedNode, err := qr.Node(ctx, equipmentType.ID)
	require.NoError(t, err)
	fetchedEquipmentType, ok := fetchedNode.(*ent.EquipmentType)
	require.True(t, ok)
	ports := fetchedEquipmentType.QueryPortDefinitions().AllX(ctx)
	require.Len(t, ports, 1)

	assert.Equal(t, ports[0].Name, "Port 1")
	assert.Equal(t, ports[0].VisibilityLabel, visibleLabel)
	assert.Equal(t, ports[0].Bandwidth, bandwidth)
}

func TestRemoveEquipmentTypeWithExistingEquipments(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "example_type_a",
	})
	require.NoError(t, err)

	locationType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "location_type_name_1"})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name_1",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	_, err = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_name_1",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	_, err = mr.RemoveEquipmentType(ctx, equipmentType.ID)
	require.Error(t, err)

	fetchedNode, err := qr.Node(ctx, equipmentType.ID)
	require.NoError(t, err)
	fetchedEquipmentType, ok := fetchedNode.(*ent.EquipmentType)
	require.True(t, ok)
	assert.Equal(t, fetchedEquipmentType.ID, equipmentType.ID)
}

func TestRemoveEquipmentType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	portDef := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: pointer.ToString("Eth1"),
		Bandwidth:    pointer.ToString("10/100/1000BASE-T"),
	}
	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        models.PropertyKindString,
		StringValue: pointer.ToString("Foo"),
	}
	position1 := models.EquipmentPositionInput{
		Name: "Position 1",
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "example_type_a",
		Positions:  []*models.EquipmentPositionInput{&position1},
		Ports:      []*models.EquipmentPortInput{&portDef},
		Properties: []*models.PropertyTypeInput{&strPropType},
	})
	require.NoError(t, err)

	_, err = mr.RemoveEquipmentType(ctx, equipmentType.ID)
	require.NoError(t, err)

	deletedNode, err := qr.Node(ctx, equipmentType.ID)
	require.NoError(t, err)
	assert.Nil(t, deletedNode)

	propertyTypes := equipmentType.QueryPropertyTypes().AllX(ctx)
	assert.Empty(t, propertyTypes)
}

func TestEditEquipmentType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	eqType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "example_type_name",
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", eqType.Name)
	c, _ := eqType.QueryCategory().Only(ctx)
	require.Nil(t, c)

	newType, err := mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:       eqType.ID,
		Name:     "example_type_name_edited",
		Category: pointer.ToString("example_type"),
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited equipment type name")
	c, _ = newType.QueryCategory().Only(ctx)
	require.Equal(t, "example_type", c.Name)

	eqType, err = mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "example_type_name_2",
	})
	require.NoError(t, err)
	_, err = mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:   eqType.ID,
		Name: "example_type_name_edited",
	})
	require.Error(t, err, "duplicate names")

	types, err := qr.EquipmentTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 2)

	node, err := qr.Node(ctx, eqType.ID)
	require.NoError(t, err)
	typ, ok := node.(*ent.EquipmentType)
	require.True(t, ok)
	require.Equal(t, "example_type_name_2", typ.Name)
}

func TestEditEquipmentTypeRemoveCategory(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	eqType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:     "example_type_name",
		Category: pointer.ToString("example_type"),
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", eqType.Name)
	c, _ := eqType.QueryCategory().Only(ctx)
	require.Equal(t, "example_type", c.Name)

	newType, err := mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:   eqType.ID,
		Name: "example_type_name_edited",
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited equipment type name")
	c, _ = newType.QueryCategory().Only(ctx)
	require.Nil(t, c)
}

func TestEditEquipmentTypeWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        models.PropertyKindString,
		StringValue: pointer.ToString("Foo"),
	}
	propTypeInput := []*models.PropertyTypeInput{&strPropType}
	eqType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "example_type_a",
		Properties: propTypeInput,
	})
	require.NoError(t, err)

	strProp := eqType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	strPropType = models.PropertyTypeInput{
		ID:          &strProp.ID,
		Name:        "str_prop_new",
		Type:        models.PropertyKindString,
		StringValue: pointer.ToString("Foo - edited"),
	}
	intPropType := models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     models.PropertyKindInt,
		IntValue: pointer.ToInt(5),
	}
	editedPropTypeInput := []*models.PropertyTypeInput{&strPropType, &intPropType}
	newType, err := mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:         eqType.ID,
		Name:       "example_type_a",
		Properties: editedPropTypeInput,
	})
	require.NoError(t, err)
	require.Equal(t, eqType.Name, newType.Name, "successfully edited equipment type name")

	strProp = eqType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	require.Equal(t, "str_prop_new", strProp.Name, "successfully edited prop type name")
	require.Equal(t, "Foo - edited", strProp.StringVal, "successfully edited prop type string value")

	intProp := eqType.QueryPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	require.Equal(t, "int_prop", intProp.Name, "successfully edited prop type name")
	require.Equal(t, 5, intProp.IntVal, "successfully edited prop type int value")

	intPropType = models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     models.PropertyKindInt,
		IntValue: pointer.ToInt(6),
	}
	editedPropTypeInput = []*models.PropertyTypeInput{&intPropType}
	_, err = mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:         eqType.ID,
		Name:       "example_type_a",
		Properties: editedPropTypeInput,
	})
	require.Error(t, err, "duplicate property type names")
}

func TestEditEquipmentTypeWithPortsAndPositions(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	bandwidth := "b1"
	label := "v1"
	strPortType := models.EquipmentPortInput{
		Name:         "str_prop",
		Bandwidth:    &bandwidth,
		VisibleLabel: &label,
	}
	posTypeA := models.EquipmentPositionInput{
		Name:         "str_prop",
		VisibleLabel: &label,
	}
	posTypeInput := []*models.EquipmentPositionInput{&posTypeA}
	portTypeInput := []*models.EquipmentPortInput{&strPortType}
	eqType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "example_type_a",
		Positions: posTypeInput,
		Ports:     portTypeInput,
	})
	require.NoError(t, err)

	bandwidth = "b1 edited"
	label = "v1 edited"
	strPort := eqType.QueryPortDefinitions().OnlyX(ctx)
	strPortType = models.EquipmentPortInput{
		ID:           &strPort.ID,
		Name:         "str_port_edited",
		Bandwidth:    &bandwidth,
		VisibleLabel: &label,
	}
	bandwidthInt := "b2 new"
	labelInt := "v2 new"
	intPortType := models.EquipmentPortInput{
		Name:         "int_port",
		Bandwidth:    &bandwidthInt,
		VisibleLabel: &labelInt,
	}
	portTypeInput = []*models.EquipmentPortInput{&strPortType, &intPortType}

	strPos := eqType.QueryPositionDefinitions().OnlyX(ctx)
	posTypeA = models.EquipmentPositionInput{
		ID:           &strPos.ID,
		Name:         "str_pos_edited",
		VisibleLabel: &label,
	}
	posTypeB := models.EquipmentPositionInput{
		Name:         "str_pos_new",
		VisibleLabel: &label,
	}
	posTypeInput = []*models.EquipmentPositionInput{&posTypeA, &posTypeB}

	newType, err := mr.EditEquipmentType(ctx, models.EditEquipmentTypeInput{
		ID:        eqType.ID,
		Name:      "example_type_a",
		Positions: posTypeInput,
		Ports:     portTypeInput,
	})
	require.NoError(t, err)
	require.Equal(t, eqType.Name, newType.Name, "successfully edited equipment type name")

	pos1 := eqType.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name("str_pos_edited")).OnlyX(ctx)
	require.Equal(t, label, pos1.VisibilityLabel, "successfully edited prop type string value")

	pos2 := eqType.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name("str_pos_new")).OnlyX(ctx)
	require.Equal(t, label, pos2.VisibilityLabel, "successfully edited prop type string value")
}
