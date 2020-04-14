// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"testing"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/ent/equipmentport"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpoint"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func pointerToServiceStatus(status models.ServiceStatus) *models.ServiceStatus {
	return &status
}

func TestAddServiceWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	serviceTypeStrValue := "Foo"
	serviceStrPropType := models.PropertyTypeInput{
		Name:        "service_str_prop",
		Type:        "string",
		StringValue: &serviceTypeStrValue,
	}
	servicePropTypeInput := []*models.PropertyTypeInput{&serviceStrPropType}

	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "Internet Access", HasCustomer: false, Properties: servicePropTypeInput})
	require.NoError(t, err)

	propertyType, err := serviceType.QueryPropertyTypes().Only(ctx)
	require.NoError(t, err)

	serviceStrValue := "Bar"
	serviceStrProp := models.PropertyInput{PropertyTypeID: propertyType.ID, StringValue: &serviceStrValue}

	servicePropInput := []*models.PropertyInput{&serviceStrProp}

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 201",
		ServiceTypeID: serviceType.ID,
		Properties:    servicePropInput,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	fetchedProperty, err := service.QueryProperties().Only(ctx)
	require.NoError(t, err)

	assert.Equal(t, fetchedProperty.StringVal, serviceStrValue)
}

func TestAddServiceWithExternalIdUnique(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "Internet Access",
		HasCustomer: false,
		Properties:  []*models.PropertyTypeInput{},
	})
	require.NoError(t, err)

	externalID1 := "S121"
	externalID2 := "S122"

	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 201",
		ServiceTypeID: serviceType.ID,
		ExternalID:    &externalID1,
		Properties:    []*models.PropertyInput{},
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	assert.Equal(t, *s1.ExternalID, externalID1)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 202",
		ServiceTypeID: serviceType.ID,
		ExternalID:    &externalID1,
		Properties:    []*models.PropertyInput{},
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.Error(t, err)

	s2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 203",
		ServiceTypeID: serviceType.ID,
		ExternalID:    &externalID2,
		Properties:    []*models.PropertyInput{},
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	assert.Equal(t, *s2.ExternalID, externalID2)
}

func TestAddServiceWithCustomer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "Internet Access",
		HasCustomer: true,
	})
	require.NoError(t, err)

	customerID := "S3213"
	customer, err := mr.AddCustomer(ctx, models.AddCustomerInput{Name: "Donald", ExternalID: &customerID})
	require.NoError(t, err)

	s, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 201",
		ServiceTypeID: serviceType.ID,
		CustomerID:    &customer.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	fetchedService, err := qr.Service(ctx, s.ID)
	require.NoError(t, err)

	customer = fetchedService.QueryCustomer().OnlyX(ctx)

	assert.Equal(t, customer.Name, "Donald")
	assert.Equal(t, *customer.ExternalID, customerID)
}

func TestServiceTopologyReturnsCorrectLinksAndEquipment(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()

	locType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "Room",
	})

	eqt, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "Router",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
			{Name: "typ1_p2"},
		},
	})

	loc, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "Room2",
		Type: locType.ID,
	})

	eq1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router1",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	eq2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router2",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	eq3, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router3",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	equipmentType := r.client.EquipmentType.GetX(ctx, eqt.ID)
	defs := equipmentType.QueryPortDefinitions().AllX(ctx)
	ep1 := eq1.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs[0].ID))).OnlyX(ctx)

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

	st, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "Internet Access",
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("CONSUMER"),
				Index:           0,
				EquipmentTypeID: eqt.ID,
			},
		},
	})

	s, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Internet Access Room 2",
		ServiceTypeID: st.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddServiceLink(ctx, s.ID, l1.ID)
	require.NoError(t, err)
	_, err = mr.AddServiceLink(ctx, s.ID, l2.ID)
	require.NoError(t, err)

	ept := st.QueryEndpointDefinitions().OnlyX(ctx)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s.ID,
		EquipmentID: eq1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	res, err := r.Service().Topology(ctx, s)
	require.NoError(t, err)

	require.Len(t, res.Nodes, 3)
	require.Len(t, res.Links, 2)
}

func TestServiceTopologyWithSlots(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()

	locType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "Room",
	})

	router, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "Router",

		Positions: []*models.EquipmentPositionInput{
			{Name: "slot1"},
		},
	})

	card, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "Card",
		Ports: []*models.EquipmentPortInput{
			{Name: "port1"},
		},
	})

	loc, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "Room2",
		Type: locType.ID,
	})

	posDefs := router.QueryPositionDefinitions().AllX(ctx)

	router1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router1",
		Type:     router.ID,
		Location: &loc.ID,
	})

	card1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "Card1",
		Type:               card.ID,
		Parent:             &router1.ID,
		PositionDefinition: &posDefs[0].ID,
	})

	router2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router2",
		Type:     router.ID,
		Location: &loc.ID,
	})

	card2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "Card2",
		Type:               card.ID,
		Parent:             &router2.ID,
		PositionDefinition: &posDefs[0].ID,
	})

	portDefs := card.QueryPortDefinitions().AllX(ctx)

	l, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: card1.ID, Port: portDefs[0].ID},
			{Equipment: card2.ID, Port: portDefs[0].ID},
		},
	})

	ep1 := card1.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(portDefs[0].ID))).OnlyX(ctx)

	st, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "Internet Access",
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("CONSUMER"),
				Index:           0,
				EquipmentTypeID: card.ID,
			},
		},
	})

	s, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Internet Access Room 2",
		ServiceTypeID: st.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddServiceLink(ctx, s.ID, l.ID)
	require.NoError(t, err)

	ept := st.QueryEndpointDefinitions().OnlyX(ctx)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s.ID,
		EquipmentID: card1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	res, err := r.Service().Topology(ctx, s)
	require.NoError(t, err)

	require.Len(t, res.Nodes, 2)
	require.Len(t, res.Links, 1)

	source, err := res.Links[0].Source.Node(ctx)
	require.NoError(t, err)
	require.Contains(t, []int{router1.ID, router2.ID}, source.ID)
	target, err := res.Links[0].Target.Node(ctx)
	require.NoError(t, err)
	require.Contains(t, []int{router1.ID, router2.ID}, target.ID)
}

func TestEditService(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name",
	})
	require.NoError(t, err)
	require.Equal(t, "service_type_name", serviceType.Name)

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_name",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	newService, err := mr.EditService(ctx, models.ServiceEditData{
		ID:   service.ID,
		Name: pointer.ToString("new_service_name"),
	})
	require.NoError(t, err)
	require.Equal(t, "new_service_name", newService.Name)

	fetchedService, _ := qr.Service(ctx, service.ID)
	require.Equal(t, newService.Name, fetchedService.Name)
}

func TestEditServiceWithExternalID(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name",
	})
	require.NoError(t, err)
	require.Equal(t, "service_type_name", serviceType.Name)

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_name",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	fetchedService, _ := qr.Service(ctx, service.ID)
	require.Nil(t, fetchedService.ExternalID)

	externalID1 := "externalID1"
	_, err = mr.EditService(ctx, models.ServiceEditData{
		ID:         service.ID,
		Name:       pointer.ToString(service.Name),
		ExternalID: &externalID1,
	})
	require.NoError(t, err)
	fetchedService, _ = qr.Service(ctx, service.ID)
	require.Equal(t, externalID1, *fetchedService.ExternalID)

	externalID2 := "externalID2"
	_, err = mr.EditService(ctx, models.ServiceEditData{
		ID:         service.ID,
		Name:       pointer.ToString(service.Name),
		ExternalID: &externalID2,
	})
	require.NoError(t, err)
	fetchedService, _ = qr.Service(ctx, service.ID)
	require.Equal(t, externalID2, *fetchedService.ExternalID)
}

func TestEditServiceWithCustomer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name",
	})
	require.NoError(t, err)
	require.Equal(t, "service_type_name", serviceType.Name)

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_name",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	fetchedService, _ := qr.Service(ctx, service.ID)
	exist := fetchedService.QueryCustomer().ExistX(ctx)
	require.Equal(t, false, exist)

	donald, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name: "Donald Duck",
	})
	require.NoError(t, err)

	dafi, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name: "Dafi Duck",
	})
	require.NoError(t, err)

	_, err = mr.EditService(ctx, models.ServiceEditData{
		ID:         service.ID,
		Name:       pointer.ToString(service.Name),
		CustomerID: &donald.ID,
	})
	require.NoError(t, err)
	fetchedService, _ = qr.Service(ctx, service.ID)
	fetchedCustomer := fetchedService.QueryCustomer().OnlyX(ctx)
	require.Equal(t, donald.ID, fetchedCustomer.ID)

	_, err = mr.EditService(ctx, models.ServiceEditData{
		ID:         service.ID,
		Name:       pointer.ToString(service.Name),
		CustomerID: &dafi.ID,
	})
	require.NoError(t, err)
	fetchedService, _ = qr.Service(ctx, service.ID)
	fetchedCustomer = fetchedService.QueryCustomer().OnlyX(ctx)
	require.Equal(t, dafi.ID, fetchedCustomer.ID)
}

func TestEditServiceWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	pTypes := models.PropertyTypeInput{
		Name: "str_prop",
		Type: "string",
	}

	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:       "type_name_1",
		Properties: []*models.PropertyTypeInput{&pTypes},
	})
	require.NoError(t, err)
	pTypeID := serviceType.QueryPropertyTypes().OnlyXID(ctx)

	strValue := "Foo"
	strProp := models.PropertyInput{
		PropertyTypeID: pTypeID,
		StringValue:    &strValue,
	}
	strValue2 := "Bar"
	strProp2 := models.PropertyInput{
		PropertyTypeID: pTypeID,
		StringValue:    &strValue2,
	}

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "inst_name_1",
		ServiceTypeID: serviceType.ID,
		Properties:    []*models.PropertyInput{&strProp, &strProp2},
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	fetchedService, _ := qr.Service(ctx, service.ID)
	fetchedProps, _ := fetchedService.QueryProperties().All(ctx)

	// Property[] -> PropertyInput[]
	var propInputClone []*models.PropertyInput
	for _, v := range fetchedProps {
		var strValue = v.StringVal + "-2"
		propInput := &models.PropertyInput{
			ID:             &v.ID,
			PropertyTypeID: v.QueryType().OnlyXID(ctx),
			StringValue:    &strValue,
		}
		propInputClone = append(propInputClone, propInput)
	}

	_, err = mr.EditService(ctx, models.ServiceEditData{
		ID:         service.ID,
		Name:       pointer.ToString("service_name_1"),
		Properties: propInputClone,
	})
	require.NoError(t, err, "Editing service")

	newFetchedService, err := qr.Service(ctx, service.ID)
	require.NoError(t, err)
	existA := newFetchedService.QueryProperties().Where(property.StringVal("Foo-2")).ExistX(ctx)
	require.NoError(t, err)
	require.True(t, existA, "Property with the new name should exist on service")
	existB := newFetchedService.QueryProperties().Where(property.StringVal("Bar-2")).ExistX(ctx)
	require.NoError(t, err)
	require.True(t, existB, "Property with the new name should exist on service")
	existC := newFetchedService.QueryProperties().Where(property.StringVal("Bar")).ExistX(ctx)
	require.NoError(t, err)
	require.False(t, existC, "Property with the old name should not exist on service")
}

func TestAddEndpointsToService(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type_name",
	})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst_name",
		Type: locType.ID,
	})
	require.NoError(t, err)

	eqType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type_name",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
		},
	})
	require.NoError(t, err)

	defs := eqType.QueryPortDefinitions().AllX(ctx)

	eq1, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst_name_1",
		Type:     eqType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	ep1 := eq1.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs[0].ID))).OnlyX(ctx)

	eq2, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst_name_2",
		Type:     eqType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	ep2 := eq2.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs[0].ID))).OnlyX(ctx)

	eq3, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst_name_3",
		Type:     eqType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	ep3 := eq3.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs[0].ID))).OnlyX(ctx)

	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name",
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("CONSUMER"),
				EquipmentTypeID: eqType.ID,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "service_type_name", serviceType.Name)

	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_name",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	ept := serviceType.QueryEndpointDefinitions().OnlyX(ctx)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service.ID,
		EquipmentID: eq1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service.ID,
		EquipmentID: eq2.ID,
		PortID:      pointer.ToInt(ep2.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	fetchedService, _ := qr.Service(ctx, service.ID)
	endpoints := fetchedService.QueryEndpoints().QueryPort().IDsX(ctx)
	require.Len(t, endpoints, 2)
	require.NotContains(t, endpoints, eq3.ID)

	e1 := fetchedService.QueryEndpoints().Where(serviceendpoint.HasPortWith(equipmentport.ID(ep1.ID))).OnlyX(ctx)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service.ID,
		EquipmentID: eq3.ID,
		PortID:      pointer.ToInt(ep3.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	_, err = mr.RemoveServiceEndpoint(ctx, e1.ID)
	require.NoError(t, err)

	require.NoError(t, err)
	fetchedService, _ = qr.Service(ctx, service.ID)
	endpoints = fetchedService.QueryEndpoints().QueryPort().IDsX(ctx)
	require.Len(t, endpoints, 2)
	require.Contains(t, endpoints, ep3.ID)
	require.NotContains(t, endpoints, ep1.ID)
}

func TestServicesOfEquipment(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	locType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "Room",
	})

	eqt, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "Router",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
			{Name: "typ1_p2"},
		},
	})

	loc, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "Room2",
		Type: locType.ID,
	})

	eq1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router1",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	eq2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router2",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	eq3, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "Router3",
		Type:     eqt.ID,
		Location: &loc.ID,
	})

	equipmentType := r.client.EquipmentType.GetX(ctx, eqt.ID)
	defs := equipmentType.QueryPortDefinitions().AllX(ctx)

	l1, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: eq1.ID, Port: defs[0].ID},
			{Equipment: eq2.ID, Port: defs[0].ID},
		},
	})

	ep1 := eq1.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs[0].ID))).OnlyX(ctx)

	l2, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: eq2.ID, Port: defs[1].ID},
			{Equipment: eq3.ID, Port: defs[1].ID},
		},
	})

	st, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "Internet Access",
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("CONSUMER"),
				EquipmentTypeID: eqt.ID,
			},
		}})

	s1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Internet Access Room 2a",
		ServiceTypeID: st.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	ept := st.QueryEndpointDefinitions().OnlyX(ctx)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s1.ID,
		EquipmentID: eq1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	s2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Internet Access Room 2b",
		ServiceTypeID: st.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddServiceLink(ctx, s2.ID, l1.ID)
	require.NoError(t, err)
	_, err = mr.AddServiceLink(ctx, s2.ID, l2.ID)
	require.NoError(t, err)
	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          s2.ID,
		EquipmentID: eq1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept.ID,
	})
	require.NoError(t, err)

	eq1Services, err := r.Equipment().Services(ctx, eq1)
	require.NoError(t, err)
	require.Len(t, eq1Services, 2)

	eq2Services, err := r.Equipment().Services(ctx, eq2)
	require.NoError(t, err)
	require.Len(t, eq2Services, 1)
	require.Equal(t, s2.ID, eq2Services[0].ID)
}

func TestAddServiceWithServiceProperty(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type", HasCustomer: false})
	require.NoError(t, err)

	service1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_1",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	index := 0
	servicePropType := models.PropertyTypeInput{
		Name:  "service_prop",
		Type:  "service",
		Index: &index,
	}

	propTypeInputs := []*models.PropertyTypeInput{&servicePropType}
	serviceTypeWithServiceProp, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "service_type_with_service_prop",
		HasCustomer: true,
		Properties:  propTypeInputs,
	})
	require.NoError(t, err)

	propType := serviceTypeWithServiceProp.QueryPropertyTypes().OnlyX(ctx)
	servicePropInput := models.PropertyInput{
		PropertyTypeID: propType.ID,
		ServiceIDValue: &service1.ID,
	}

	service2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "service_2",
		ServiceTypeID: serviceTypeWithServiceProp.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
		Properties:    []*models.PropertyInput{&servicePropInput},
	})
	require.NoError(t, err)

	serviceProp := service2.QueryProperties().Where(property.HasTypeWith(propertytype.Name("service_prop"))).OnlyX(ctx)
	serviceValue := serviceProp.QueryServiceValue().OnlyX(ctx)

	require.Equal(t, "service_1", serviceValue.Name)
}

func TestAddServiceEndpointType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type_name",
	})
	require.NoError(t, err)

	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst_name",
		Type: locType.ID,
	})
	require.NoError(t, err)

	eqType1, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type_name",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
		},
	})
	require.NoError(t, err)
	defs1 := eqType1.QueryPortDefinitions().AllX(ctx)

	eqType2, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type_name2",
		Ports: []*models.EquipmentPortInput{
			{Name: "typ1_p1"},
		},
	})
	require.NoError(t, err)
	defs2 := eqType2.QueryPortDefinitions().AllX(ctx)

	eq1, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst_name_1",
		Type:     eqType1.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	ep1 := eq1.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs1[0].ID))).OnlyX(ctx)

	eq2, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst_name_2",
		Type:     eqType2.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	ep2 := eq2.QueryPorts().Where(equipmentport.HasDefinitionWith(equipmentportdefinition.ID(defs2[0].ID))).OnlyX(ctx)

	serviceType1, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name",
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("CONSUMER"),
				EquipmentTypeID: eqType1.ID,
			},
		},
	})
	require.NoError(t, err)
	require.Equal(t, "service_type_name", serviceType1.Name)

	ept1 := serviceType1.QueryEndpointDefinitions().OnlyX(ctx)

	service1, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent building, room 201",
		ServiceTypeID: serviceType1.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service1.ID,
		EquipmentID: eq1.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept1.ID,
	})
	require.NoError(t, err)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service1.ID,
		EquipmentID: eq2.ID,
		PortID:      pointer.ToInt(ep2.ID),
		Definition:  ept1.ID,
	})
	require.Error(t, err, "port equipment is of different type than service type")

	serviceType2, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "service_type_name2",
	})
	require.NoError(t, err)
	service2, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Kent2 building, room 401",
		ServiceTypeID: serviceType2.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	_, err = mr.AddServiceEndpoint(ctx, models.AddServiceEndpointInput{
		ID:          service2.ID,
		EquipmentID: eq2.ID,
		PortID:      pointer.ToInt(ep1.ID),
		Definition:  ept1.ID,
	})
	require.Error(t, err, "service is of different type than service endpoint type")
}
