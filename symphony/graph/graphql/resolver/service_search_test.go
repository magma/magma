// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type serviceSearchDataModels struct {
	st1         string
	st2         string
	strType     string
	intType     string
	boolType    string
	floatType   string
	locRoom     string
	locBuilding string
	eqt         string
}

func preparePropertyTypes() []*models.PropertyTypeInput {
	serviceStrPropType := models.PropertyTypeInput{
		Name:        "service_str_prop",
		Type:        "string",
		StringValue: pointer.ToString("Foo is the best"),
	}
	serviceIntPropType := models.PropertyTypeInput{
		Name: "service_int_prop",
		Type: "int",
	}
	serviceBoolPropType := models.PropertyTypeInput{
		Name: "service_bool_prop",
		Type: "bool",
	}
	serviceFloatPropType := models.PropertyTypeInput{
		Name: "service_float_prop",
		Type: "float",
	}

	return []*models.PropertyTypeInput{
		&serviceStrPropType,
		&serviceIntPropType,
		&serviceBoolPropType,
		&serviceFloatPropType,
	}
}

func prepareServiceData(ctx context.Context, r *TestResolver) serviceSearchDataModels {
	mr := r.Mutation()

	props := preparePropertyTypes()

	st1, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "Internet Access", HasCustomer: false, Properties: props})

	st2, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "Internet Access 2", HasCustomer: false, Properties: []*models.PropertyTypeInput{}})

	strType, _ := st1.QueryPropertyTypes().Where(propertytype.Name("service_str_prop")).Only(ctx)
	intType, _ := st1.QueryPropertyTypes().Where(propertytype.Name("service_int_prop")).Only(ctx)
	boolType, _ := st1.QueryPropertyTypes().Where(propertytype.Name("service_bool_prop")).Only(ctx)
	floatType, _ := st1.QueryPropertyTypes().Where(propertytype.Name("service_float_prop")).Only(ctx)

	locBuilding, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "building",
	})

	locRoom, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "room",
	})

	eqt, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
			{Name: "typ1_p2"},
		},
	})

	return serviceSearchDataModels{
		st1.ID,
		st2.ID,
		strType.ID,
		intType.ID,
		boolType.ID,
		floatType.ID,
		locRoom.ID,
		locBuilding.ID,
		eqt.ID,
	}
}

func TestSearchServicesByName(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)

	data := prepareServiceData(ctx, r)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 2010",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstName,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("Room"),
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 3)

	f2 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstName,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("Room 201"),
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 2)
}

func TestSearchServicesByStatus(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)

	data := prepareServiceData(ctx, r)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusMaintenance),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusInService),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 2010",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusInService),
	})
	require.NoError(t, err)

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceStatus,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{models.ServiceStatusMaintenance.String()},
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 1)

	f2 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceStatus,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{models.ServiceStatusInService.String()},
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 2)

	f3 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceStatus,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{models.ServiceStatusPending.String()},
	}
	res3, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f3}, &limit)
	require.NoError(t, err)
	require.Len(t, res3.Services, 0)
}

func TestSearchServicesByType(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)

	data := prepareServiceData(ctx, r)

	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	s2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st2,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{data.st1},
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 1)
	assert.Equal(t, res1.Services[0].ID, s1.ID)

	f2 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{data.st2},
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 1)
	assert.Equal(t, res2.Services[0].ID, s2.ID)

	f3 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{data.st1, data.st2},
	}
	res3, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f3}, &limit)
	require.NoError(t, err)
	require.Len(t, res3.Services, 2)
}

func TestSearchServicesByExternalID(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)
	data := prepareServiceData(ctx, r)

	externalID1 := "S1111"
	externalID2 := "S2222"
	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		ExternalID:    &externalID1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	s2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st2,
		ExternalID:    &externalID2,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 203",
		ServiceTypeID: data.st2,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstExternalID,
		Operator:    models.FilterOperatorIs,
		StringValue: &externalID1,
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 1)
	assert.Equal(t, res1.Services[0].ID, s1.ID)

	f2 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstExternalID,
		Operator:    models.FilterOperatorIs,
		StringValue: &externalID2,
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 1)
	assert.Equal(t, res2.Services[0].ID, s2.ID)
}

func TestSearchServicesByCustomerName(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)
	data := prepareServiceData(ctx, r)

	customerA, err := mr.AddCustomer(ctx, models.AddCustomerInput{Name: "Donald"})
	require.NoError(t, err)

	customerB, err := mr.AddCustomer(ctx, models.AddCustomerInput{Name: "Mia", ExternalID: pointer.ToString("4242")})
	require.NoError(t, err)

	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		CustomerID:    &customerA.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st2,
		CustomerID:    &customerB.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Lobby",
		ServiceTypeID: data.st2,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstCustomerName,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("Donald"),
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 1)
	assert.Equal(t, res1.Services[0].ID, s1.ID)

	f2 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeServiceInstCustomerName,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("a"),
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 2)
}

func TestSearchServicesByProperties(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	qr, mr := r.Query(), r.Mutation()
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	data := prepareServiceData(ctx, r)

	serviceStrProp := models.PropertyInput{
		PropertyTypeID: data.strType,
		StringValue:    pointer.ToString("Bar is the best"),
	}
	servicePropInput := []*models.PropertyInput{&serviceStrProp}

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 201",
		ServiceTypeID: data.st1,
		Properties:    servicePropInput,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Room 202",
		ServiceTypeID: data.st1,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	limit := 100
	all, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{}, &limit)
	require.NoError(t, err)
	require.Len(t, all.Services, 2)
	f := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeServiceInstProperty,
		Operator:   models.FilterOperatorIs,
		PropertyValue: &models.PropertyTypeInput{
			Name:        "service_str_prop",
			Type:        models.PropertyKind("string"),
			StringValue: pointer.ToString("Foo is the best"),
		},
	}
	res, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f}, &limit)
	require.NoError(t, err)
	require.Len(t, res.Services, 1)
}

func TestSearchServicesByLocations(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)
	data := prepareServiceData(ctx, r)

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst1",
		Type: data.locBuilding,
	})
	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "loc_inst2",
		Type:   data.locRoom,
		Parent: &loc1.ID,
	})
	loc3, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "loc_inst3",
		Type:   data.locRoom,
		Parent: &loc1.ID,
	})

	eq1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst2",
		Type:     data.eqt,
		Location: &loc2.ID,
	})
	eq2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst3",
		Type:     data.eqt,
		Location: &loc3.ID,
	})

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 201",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq1.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 202",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq2.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 203",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq1.ID, eq2.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	maxDepth := 2
	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{loc1.ID},
		MaxDepth:   &maxDepth,
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 3)

	f2 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{loc2.ID},
		MaxDepth:   &maxDepth,
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 2)

	f3 := models.ServiceFilterInput{
		FilterType: models.ServiceFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []string{loc2.ID, loc3.ID},
		MaxDepth:   &maxDepth,
	}
	res3, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f3}, &limit)
	require.NoError(t, err)
	require.Len(t, res3.Services, 3)
}

func TestSearchServicesByEquipmentInside(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	qr, mr := r.Query(), r.Mutation()
	ctx := viewertest.NewContext(r.client)
	data := prepareServiceData(ctx, r)

	loc, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst",
		Type: data.locRoom,
	})

	eq1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst1",
		Type:     data.eqt,
		Location: &loc.ID,
	})

	eq2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst2",
		Type:     data.eqt,
		Location: &loc.ID,
	})

	eq3, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst3",
		Type:     data.eqt,
		Location: &loc.ID,
	})

	equipmentType := r.client.EquipmentType.GetX(ctx, data.eqt)
	defs := equipmentType.QueryPortDefinitions().AllX(ctx)

	l1, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: eq1.ID, Port: defs[0].ID},
			{Equipment: eq2.ID, Port: defs[0].ID},
		},
	})
	l2, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: eq2.ID, Port: defs[1].ID},
			{Equipment: eq3.ID, Port: defs[1].ID},
		},
	})

	s1, _ := mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 201",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq1.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})
	_, _ = mr.AddServiceLink(ctx, s1.ID, l1.ID)
	_, _ = mr.AddServiceLink(ctx, s1.ID, l2.ID)

	s2, _ := mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 202",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq1.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})
	_, _ = mr.AddServiceLink(ctx, s2.ID, l1.ID)

	_, _ = mr.AddService(ctx, models.ServiceCreateData{
		Name:                "Room 203",
		ServiceTypeID:       data.st1,
		TerminationPointIds: []string{eq1.ID},
		Status:              pointerToServiceStatus(models.ServiceStatusPending),
	})

	limit := 100
	f1 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeEquipmentInService,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("eq_"),
	}
	res1, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Services, 3)

	f2 := models.ServiceFilterInput{
		FilterType:  models.ServiceFilterTypeEquipmentInService,
		Operator:    models.FilterOperatorContains,
		StringValue: pointer.ToString("eq_inst3"),
	}
	res2, err := qr.ServiceSearch(ctx, []*models.ServiceFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Services, 1)
}
