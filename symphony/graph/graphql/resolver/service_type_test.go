// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

func TestEditServiceType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "example_type_name", HasCustomer: false})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", serviceType.Name)

	newType, err := mr.EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          serviceType.ID,
		Name:        "example_type_name_edited",
		HasCustomer: true,
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited service type name")
	require.Equal(t, true, newType.HasCustomer)

	serviceType, err = mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "example_type_name_2"})
	require.NoError(t, err)
	_, err = mr.EditServiceType(ctx, models.ServiceTypeEditData{
		ID:   serviceType.ID,
		Name: "example_type_name_edited",
	})
	require.Error(t, err, "duplicate names")

	types, err := qr.ServiceTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 2)

	typ, err := qr.ServiceType(ctx, serviceType.ID)
	require.NoError(t, err)
	require.Equal(t, "example_type_name_2", typ.Name)
}

func TestEditServiceTypeWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	strValue := "Foo"
	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        "string",
		StringValue: &strValue,
	}
	propTypeInput := []*models.PropertyTypeInput{&strPropType}
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "example_type_a", HasCustomer: true, Properties: propTypeInput})
	require.NoError(t, err)

	strProp := serviceType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
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
	newType, err := mr.EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          serviceType.ID,
		Name:        "example_type_a",
		Properties:  editedPropTypeInput,
		HasCustomer: true,
	})
	require.NoError(t, err)
	require.Equal(t, serviceType.Name, newType.Name, "successfully edited service type name")

	strProp = serviceType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	require.Equal(t, "str_prop_new", strProp.Name, "successfully edited prop type name")
	require.Equal(t, strValue, strProp.StringVal, "successfully edited prop type string value")

	intProp := serviceType.QueryPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	require.Equal(t, "int_prop", intProp.Name, "successfully edited prop type name")
	require.Equal(t, intValue, intProp.IntVal, "successfully edited prop type int value")

	intValue = 6
	intPropType = models.PropertyTypeInput{
		ID:       &intProp.ID,
		Name:     "int_prop",
		Type:     "int",
		IntValue: &intValue,
	}
	editedPropTypeInput = []*models.PropertyTypeInput{&intPropType}
	serviceType, err = mr.EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          serviceType.ID,
		Name:        "example_type_a",
		Properties:  editedPropTypeInput,
		HasCustomer: true,
	})
	require.NoError(t, err)
	intProp = serviceType.QueryPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	require.Equal(t, "int_prop", intProp.Name, "successfully edited prop type name")
	require.Equal(t, intValue, intProp.IntVal, "successfully edited prop type int value")
}
