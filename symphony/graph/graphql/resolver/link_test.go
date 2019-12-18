// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint: goconst, ineffassign
package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAddLink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr, lr := r.Mutation(), r.Query(), r.EquipmentPort(), r.Link()

	locationType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "location_type"})
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})

	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
	})
	assert.Nil(t, err)
	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	assert.Equal(t, fetchedPortA.QueryParent().OnlyXID(ctx), equipmentA.ID)
	assert.Equal(t, fetchedPortB.QueryParent().OnlyXID(ctx), equipmentB.ID)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	fetchedPorts, err := lr.Ports(ctx, createdLink)
	require.NoError(t, err)
	assert.Len(t, fetchedPorts, 2)
}

func TestAddLinkWithProperties(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr, lr := r.Mutation(), r.Query(), r.EquipmentPort(), r.Link()

	locationType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "location_type"})
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	linkStrValue := "Foo"
	linkStrPropType := models.PropertyTypeInput{
		Name:        "link_str_prop",
		Type:        models.PropertyKindString,
		StringValue: &linkStrValue,
	}
	linkPropTypeInput := []*models.PropertyTypeInput{&linkStrPropType}
	portType, err := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:           "example_type_a",
		LinkProperties: linkPropTypeInput,
	})
	assert.Nil(t, err)

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
		PortTypeID:   &portType.ID,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)

	linkVal := "Bar"
	linkPropTypeID := portType.QueryLinkPropertyTypes().FirstXID(ctx)
	linkProp := models.PropertyInput{
		StringValue:    &linkVal,
		PropertyTypeID: linkPropTypeID,
	}
	propInput := []*models.PropertyInput{&linkProp}
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		Properties: propInput,
	})
	assert.Nil(t, err)
	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	assert.Equal(t, fetchedPortA.QueryParent().OnlyXID(ctx), equipmentA.ID)
	assert.Equal(t, fetchedPortB.QueryParent().OnlyXID(ctx), equipmentB.ID)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	fetchedPorts, err := lr.Ports(ctx, createdLink)
	require.NoError(t, err)
	assert.Len(t, fetchedPorts, 2)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	propA := linkA.QueryProperties().FirstX(ctx)
	propZ := linkB.QueryProperties().FirstX(ctx)

	assert.Equal(t, propA.StringVal, linkVal)
	assert.Equal(t, propZ.StringVal, linkVal)
}

func TestEditLinkWithProperties(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr, lr := r.Mutation(), r.Query(), r.EquipmentPort(), r.Link()

	locationType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "location_type",
	})
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	linkStrValue := "Foo"
	linkStrPropType := models.PropertyTypeInput{
		Name:        "link_str_prop",
		Type:        models.PropertyKindString,
		StringValue: &linkStrValue,
	}
	linkPropTypeInput := []*models.PropertyTypeInput{&linkStrPropType}
	portType, _ := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name:           "example_type_a",
		LinkProperties: linkPropTypeInput,
	})

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
		PortTypeID:   &portType.ID,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)

	linkVal := "Bar"
	linkPropTypeID := portType.QueryLinkPropertyTypes().FirstXID(ctx)
	linkProp := models.PropertyInput{
		StringValue:    &linkVal,
		PropertyTypeID: linkPropTypeID,
	}
	propInput := []*models.PropertyInput{&linkProp}
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		Properties: propInput,
	})
	assert.NoError(t, err)
	propID := createdLink.QueryProperties().FirstXID(ctx)

	editedLinkVal := "Baz"
	editedLinkProp := models.PropertyInput{
		ID:             &propID,
		StringValue:    &editedLinkVal,
		PropertyTypeID: linkPropTypeID,
	}
	editedPropInput := []*models.PropertyInput{&editedLinkProp}
	editedLink, err := mr.EditLink(ctx, models.EditLinkInput{
		ID:         createdLink.ID,
		Properties: editedPropInput,
	})
	assert.Nil(t, err)
	assert.Equal(t, editedLink.ID, createdLink.ID)

	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	assert.Equal(t, fetchedPortA.QueryParent().OnlyXID(ctx), equipmentA.ID)
	assert.Equal(t, fetchedPortB.QueryParent().OnlyXID(ctx), equipmentB.ID)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	fetchedPorts, err := lr.Ports(ctx, createdLink)
	require.NoError(t, err)
	assert.Len(t, fetchedPorts, 2)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	propA := linkA.QueryProperties().FirstX(ctx)
	propZ := linkB.QueryProperties().FirstX(ctx)

	assert.Equal(t, propA.StringVal, editedLinkVal)
	assert.Equal(t, propZ.StringVal, editedLinkVal)
}

func TestRemoveLink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr := r.Mutation(), r.Query(), r.EquipmentPort()
	locationType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "location_type",
	})
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name",
		Type: locationType.ID,
	})
	require.NoError(t, err)

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	link, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
	})
	assert.Nil(t, err)
	assert.NotNil(t, link)

	_, _ = mr.RemoveLink(ctx, link.ID, nil)

	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.Nil(t, linkA)
	assert.Nil(t, linkB)
}

func TestAddLinkWithWorkOrder(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr, wor := r.Mutation(), r.Query(), r.EquipmentPort(), r.WorkOrder()

	workOrder := createWorkOrder(ctx, t, *r, "work_order_name_102")
	location := workOrder.QueryLocation().FirstX(ctx)

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		WorkOrder: &workOrder.ID,
	})
	assert.NoError(t, err)
	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	assert.Equal(t, fetchedPortA.QueryParent().OnlyXID(ctx), equipmentA.ID)
	assert.Equal(t, fetchedPortB.QueryParent().OnlyXID(ctx), equipmentB.ID)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.Equal(t, linkA.ID, createdLink.ID)
	assert.Equal(t, linkB.ID, createdLink.ID)

	fetchedWorkOrder, err := qr.WorkOrder(ctx, workOrder.ID)
	require.NoError(t, err)

	linksToRemove, err := wor.LinksToRemove(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Len(t, linksToRemove, 0)

	linksToAdd, err := wor.LinksToAdd(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Len(t, linksToAdd, 1)
	assert.Equal(t, linksToAdd[0].ID, createdLink.ID)
}

func TestRemoveLinkWithWorkOrder(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr, wor := r.Mutation(), r.Query(), r.EquipmentPort(), r.WorkOrder()

	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)

	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
	}
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "parent_equipment_type",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_a",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	equipmentB, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_b",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	link, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
	})
	require.NoError(t, err)
	assert.NotNil(t, link)

	_, _ = mr.RemoveLink(ctx, link.ID, &workOrder.ID)

	fetchedEquipmentA, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedEquipmentB, _ := qr.Equipment(ctx, equipmentB.ID)
	fetchedPortA := fetchedEquipmentA.QueryPorts().OnlyX(ctx)
	fetchedPortB := fetchedEquipmentB.QueryPorts().OnlyX(ctx)

	linkA, _ := pr.Link(ctx, fetchedPortA)
	linkB, _ := pr.Link(ctx, fetchedPortB)

	assert.NotNil(t, linkA)
	assert.NotNil(t, linkB)

	fetchedWorkOrder, err := qr.WorkOrder(ctx, workOrder.ID)
	require.NoError(t, err)

	linksToRemove, err := wor.LinksToRemove(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Len(t, linksToRemove, 1)
	assert.Equal(t, linksToRemove[0].ID, link.ID)

	linksToAdd, err := wor.LinksToAdd(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Len(t, linksToAdd, 0)
}
