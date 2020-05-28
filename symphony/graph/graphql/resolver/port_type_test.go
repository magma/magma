// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint: goconst
package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent/propertytype"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestAddEquipmentPortType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	_, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name: "example_type",
	})
	require.NoError(t, err)

	portTypes, err := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, portTypes.Edges, 1, "two attempts to create EquipmentPort types with same name will create one EquipmentPort type")
	require.Equal(t, portTypes.Edges[0].Node.Name, "example_type", "verifying EquipmentPort type name")
}

func TestAddEquipmentPortTypes(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	_, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_a"})
	require.NoError(t, err)
	_, err = mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_b"})
	require.NoError(t, err)

	types, _ := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.Len(t, types.Edges, 2, "verify the added EquipmentPort types are fetched properly")
}

func TestAddEquipmentPortTypeWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	strValue, strIndex := "Foo", 7
	intValue, intIndex := 5, 12

	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        models.PropertyKindString,
		Index:       &strIndex,
		StringValue: &strValue,
	}
	intPropType := models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     models.PropertyKindInt,
		Index:    &intIndex,
		IntValue: &intValue,
	}
	propTypeInputs := []*models.PropertyTypeInput{&strPropType, &intPropType}
	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:       "example_type_a",
		Properties: propTypeInputs,
	})
	require.NoError(t, err)

	intProp := portType.QueryPropertyTypes().Where(propertytype.Type(models.PropertyKindInt.String())).OnlyX(ctx)
	strProp := portType.QueryPropertyTypes().Where(propertytype.Type(models.PropertyKindString.String())).OnlyX(ctx)

	require.Equal(t, "int_prop", intProp.Name, "verifying int property type's name")
	require.Equal(t, "", intProp.StringVal, "verifying int property type's string value (default as this is an int property)")
	require.Equal(t, intValue, intProp.IntVal, "verifying int property type's int value")
	require.Equal(t, intIndex, intProp.Index, "verifying int property type's index")
	require.Equal(t, "str_prop", strProp.Name, "verifying string property type's name")
	require.Equal(t, strValue, strProp.StringVal, "verifying string property type's String value")
	require.Equal(t, 0, strProp.IntVal, "verifying int property type's int value")
	require.Equal(t, strIndex, strProp.Index, "verifying string property type's index")
}

func TestAddEquipmentPortTypeWithLinkProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	strValue, strIndex := "Foo", 7
	intValue, intIndex := 5, 12

	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        models.PropertyKindString,
		Index:       &strIndex,
		StringValue: &strValue,
	}
	intPropType := models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     models.PropertyKindInt,
		Index:    &intIndex,
		IntValue: &intValue,
	}
	propTypeInputs := []*models.PropertyTypeInput{&strPropType, &intPropType}
	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:           "example_type_a",
		LinkProperties: propTypeInputs,
	})
	require.NoError(t, err)

	intProp := portType.QueryLinkPropertyTypes().Where(propertytype.Type(models.PropertyKindInt.String())).OnlyX(ctx)
	strProp := portType.QueryLinkPropertyTypes().Where(propertytype.Type(models.PropertyKindString.String())).OnlyX(ctx)

	require.Equal(t, "int_prop", intProp.Name, "verifying int property type's name")
	require.Equal(t, "", intProp.StringVal, "verifying int property type's string value (default as this is an int property)")
	require.Equal(t, intValue, intProp.IntVal, "verifying int property type's int value")
	require.Equal(t, intIndex, intProp.Index, "verifying int property type's index")
	require.Equal(t, "str_prop", strProp.Name, "verifying string property type's name")
	require.Equal(t, strValue, strProp.StringVal, "verifying string property type's String value")
	require.Equal(t, 0, strProp.IntVal, "verifying int property type's int value")
	require.Equal(t, strIndex, strProp.Index, "verifying string property type's index")
}

func TestAddEquipmentPortTypesSameName(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", portType.Name)
	_, err = mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_name"})
	require.Error(t, err, "adding EquipmentPort type with an existing EquipmentPort type name yields an error")
	types, _ := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.Len(t, types.Edges, 1, "two attempts to create EquipmentPort types with same name will create one EquipmentPort type")
}

func TestRemoveEquipmentPortType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", portType.Name)

	types, err := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 1)

	id, err := mr.RemoveEquipmentPortType(ctx, portType.ID)
	require.NoError(t, err)
	require.Equal(t, portType.ID, id, "successfully remove EquipmentPort type")
	types, err = qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Empty(t, types.Edges, "no EquipmentPort types exist after deletion")
}

func TestEditEquipmentPortType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", portType.Name)

	newType, err := mr.EditEquipmentPortType(ctx, models.EditEquipmentPortTypeInput{
		ID:   portType.ID,
		Name: "example_type_name_edited",
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited EquipmentPort type name")

	portType, err = mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name: "example_type_name_2",
	})
	require.NoError(t, err)
	_, err = mr.EditEquipmentPortType(ctx, models.EditEquipmentPortTypeInput{
		ID:   portType.ID,
		Name: "example_type_name_edited",
	})
	require.Error(t, err, "Duplicate port type name")

	portTypes, err := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, portTypes.Edges, 2, "two attempts to create EquipmentPort types with same name will create one EquipmentPort type")
}

func TestEditEquipmentPortTypeWithLinkProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	strValue := "Foo"
	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        "string",
		StringValue: &strValue,
	}
	propTypeInput := []*models.PropertyTypeInput{&strPropType}
	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:           "example_type_a",
		LinkProperties: propTypeInput,
	})
	require.NoError(t, err)

	strProp := portType.QueryLinkPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	strValue = "Foo - edited"
	intValue := 5
	strPropType = models.PropertyTypeInput{
		ID:          &strProp.ID,
		Name:        "str_prop_new",
		Type:        "string",
		StringValue: &strValue,
	}
	intPropType := models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     "int",
		IntValue: &intValue,
	}
	editedPropTypeInput := []*models.PropertyTypeInput{&strPropType, &intPropType}
	newType, err := mr.EditEquipmentPortType(ctx, models.EditEquipmentPortTypeInput{
		ID:             portType.ID,
		Name:           "example_type_a",
		LinkProperties: editedPropTypeInput,
	})
	require.NoError(t, err)
	require.Equal(t, portType.Name, newType.Name, "successfully edited EquipmentPort type name")

	strProp = portType.QueryLinkPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	require.Equal(t, "str_prop_new", strProp.Name, "successfully edited prop type name")
	require.Equal(t, strValue, strProp.StringVal, "successfully edited prop type string value")

	intProp := portType.QueryLinkPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	require.Equal(t, "int_prop", intProp.Name, "successfully edited prop type name")
	require.Equal(t, intValue, intProp.IntVal, "successfully edited prop type int value")

	intValue = 6
	intPropType = models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     "int",
		IntValue: &intValue,
	}
	editedPropTypeInput = []*models.PropertyTypeInput{&intPropType}
	_, err = mr.EditEquipmentPortType(ctx, models.EditEquipmentPortTypeInput{
		ID:             portType.ID,
		Name:           "example_type_a",
		LinkProperties: editedPropTypeInput,
	})
	require.Error(t, err, "duplicate property type names")
}

func TestEditEquipmentPortTypeWithLinkPropertiesSameName(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	strValue := "Foo"
	strPropType := models.PropertyTypeInput{
		Name:        "foo_prop",
		Type:        "string",
		StringValue: &strValue,
	}
	propTypeInput := []*models.PropertyTypeInput{&strPropType}
	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:           "example_type_a",
		LinkProperties: propTypeInput,
	})
	require.NoError(t, err)

	intValue := 5
	intPropType := models.PropertyTypeInput{
		Name:     "foo_prop",
		Type:     "int",
		IntValue: &intValue,
	}
	editedPropTypeInput := []*models.PropertyTypeInput{&strPropType, &intPropType}
	_, err = mr.EditEquipmentPortType(ctx, models.EditEquipmentPortTypeInput{
		ID:             portType.ID,
		Name:           "example_type_a",
		LinkProperties: editedPropTypeInput,
	})
	require.Error(t, err)
}

func TestRemoveEquipmentPortTypeWithLinkedEquipmentType(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name: "example_type_name",
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", portType.Name)

	portDef := models.EquipmentPortInput{
		Name:       "Port 1",
		PortTypeID: &portType.ID,
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "example_type_a",
		Ports: []*models.EquipmentPortInput{&portDef},
	})
	require.NoError(t, err)

	types, err := qr.EquipmentPortTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 1)
	typ := types.Edges[0]
	require.Equal(t, 1, r.client.EquipmentPortDefinition.Query().CountX(ctx))
	require.Equal(t, 1, typ.Node.QueryPortDefinitions().CountX(ctx))
	_, err = mr.RemoveEquipmentType(ctx, equipmentType.ID)
	require.NoError(t, err)
	def, _ := r.client.EquipmentPortDefinition.Query().All(ctx)
	require.Equal(t, 0, len(def))
	require.Equal(t, 0, r.client.EquipmentPortDefinition.Query().CountX(ctx))
	require.Equal(t, 0, typ.Node.QueryPortDefinitions().CountX(ctx))
}
