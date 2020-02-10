// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strconv"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createPort() models.EquipmentPortInput {
	visibleLabel := "Eth1"
	bandwidth := "10/100/1000BASE-T"
	portInput := models.EquipmentPortInput{
		Name:         "Port 1",
		VisibleLabel: &visibleLabel,
		Bandwidth:    &bandwidth,
	}
	return portInput
}

func createLocation(ctx context.Context, t *testing.T, r TestResolver) *ent.Location {
	return createLocationWithName(ctx, t, r, "location_name_1")
}

func createLocationWithName(ctx context.Context, t *testing.T, r TestResolver, name string) *ent.Location {
	mr := r.Mutation()
	locationType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "_type",
	})
	require.NoError(t, err)
	location, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name,
		Type: locationType.ID,
	})
	require.NoError(t, err)
	return location
}

func createPosition() *models.EquipmentPositionInput {
	label1 := "label1"
	position := models.EquipmentPositionInput{
		Name:         "Position 1",
		VisibleLabel: &label1,
	}
	return &position
}

func createWorkOrder(ctx context.Context, t *testing.T, r TestResolver, name string) *ent.WorkOrder {
	mr := r.Mutation()
	location := createLocationWithName(ctx, t, r, name+"location")
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: name + "type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     nil,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)
	return workOrder
}

func executeWorkOrder(ctx context.Context, t *testing.T, mr generated.MutationResolver, workOrder ent.WorkOrder) (*models.WorkOrderExecutionResult, error) {
	_, err := mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerName:   workOrder.OwnerName,
		InstallDate: &workOrder.InstallDate,
		Status:      models.WorkOrderStatusDone,
		Priority:    models.WorkOrderPriorityNone,
		Assignee:    &workOrder.Assignee,
	})
	require.NoError(t, err)
	return mr.ExecuteWorkOrder(ctx, workOrder.ID)
}

const (
	longWorkOrderName     = "long_work_order"
	longWorkOrderDesc     = "long_work_order_description"
	longWorkOrderAssignee = "long_work_order_Assignee"
)

func TestAddWorkOrderWithLocation(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr, wr := r.Mutation(), r.Query(), r.WorkOrder()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)

	fetchedLocation, err := wr.Location(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Equal(t, fetchedLocation.ID, location.ID)
	assert.Equal(t, fetchedLocation.Name, location.Name)
}

func TestAddWorkOrderWithType(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr, wr := r.Mutation(), r.Query(), r.WorkOrder()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)

	fetchedWorkOrderType, err := wr.WorkOrderType(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Equal(t, fetchedWorkOrderType.ID, woType.ID)
	assert.Equal(t, fetchedWorkOrderType.Name, woType.Name)
}

func TestAddWorkOrderWithAssignee(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr, wr := r.Mutation(), r.Query(), r.WorkOrder()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	assignee := longWorkOrderAssignee
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)
	require.Empty(t, workOrder.Assignee)

	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerName:   workOrder.OwnerName,
		Status:      models.WorkOrderStatusPending,
		Priority:    models.WorkOrderPriorityNone,
		Assignee:    &assignee,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)
	require.Equal(t, &workOrder.Assignee, &assignee)

	fetchedWorkOrderType, err := wr.WorkOrderType(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Equal(t, fetchedWorkOrderType.ID, woType.ID)
	assert.Equal(t, fetchedWorkOrderType.Name, woType.Name)
}

func TestAddWorkOrderInvalidType(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()
	name := longWorkOrderName
	description := longWorkOrderDesc
	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: "123",
		LocationID:      nil,
	})
	require.Error(t, err)
}

func TestEditInvalidWorkOrder(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	_, err = r.Mutation().EditWorkOrderType(ctx, models.EditWorkOrderTypeInput{
		ID:   "234",
		Name: "foo",
	})
	require.Error(t, err)
}

func TestAddWorkOrderWithDescription(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)

	assert.Equal(t, fetchedWorkOrder.Name, name)
	assert.Equal(t, fetchedWorkOrder.Description, description)
	assert.Equal(t, location.ID, workOrder.QueryLocation().OnlyX(ctx).ID)
}

func TestAddWorkOrderWithPriority(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()
	name := longWorkOrderName
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	pri := models.WorkOrderPriorityLow
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		WorkOrderTypeID: woType.ID,
		Priority:        &pri,
	})
	require.NoError(t, err)
	require.Equal(t, workOrder.Assignee, "")
	require.EqualValues(t, pri, workOrder.Priority)

	input := models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerName:   workOrder.OwnerName,
		Status:      models.WorkOrderStatusPending,
		Priority:    models.WorkOrderPriorityHigh,
		Index:       pointer.ToInt(42),
	}

	workOrder, err = mr.EditWorkOrder(ctx, input)
	require.NoError(t, err)
	require.EqualValues(t, input.Priority, workOrder.Priority)
	require.Equal(t, *input.Index, workOrder.Index)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	workOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)
	require.EqualValues(t, input.Priority, workOrder.Priority)
	require.Equal(t, *input.Index, workOrder.Index)
}

func TestAddWorkOrderWithProject(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, pr := r.Mutation(), r.Project()

	input := models.AddProjectTypeInput{Name: "test", Description: pointer.ToString("test desc")}
	ltyp, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type",
	})
	require.NoError(t, err)
	loc, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_name",
		Type: ltyp.ID,
	})
	require.NoError(t, err)
	typ, err := mr.CreateProjectType(ctx, input)
	require.NoError(t, err)
	pinput := models.AddProjectInput{Name: "test", Description: pointer.ToString("desc"), Type: typ.ID, Location: &loc.ID}
	project, err := mr.CreateProject(ctx, pinput)
	require.NoError(t, err)
	woNum, err := pr.NumberOfWorkOrders(ctx, project)
	require.NoError(t, err)
	require.Zero(t, woNum)

	name := longWorkOrderName
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		WorkOrderTypeID: woType.ID,
		ProjectID:       &project.ID,
	})
	require.NoError(t, err)
	require.Equal(t, workOrder.QueryProject().OnlyX(ctx).ID, project.ID)
	woNum, err = pr.NumberOfWorkOrders(ctx, project)
	require.NoError(t, err)
	require.Equal(t, 1, woNum)
	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:        workOrder.ID,
		Name:      workOrder.Name,
		OwnerName: workOrder.OwnerName,
	})
	require.NoError(t, err)
	fetchProject, err := workOrder.QueryProject().Only(ctx)
	require.Error(t, err)
	require.Nil(t, fetchProject)
	woNum, err = pr.NumberOfWorkOrders(ctx, project)
	require.NoError(t, err)
	require.Zero(t, woNum)
}

func TestAddWorkOrderWithComment(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()
	w := createWorkOrder(ctx, t, *r, "Foo")
	require.NoError(t, err)

	node, err := qr.Node(ctx, w.ID)
	require.NoError(t, err)
	w, ok := node.(*ent.WorkOrder)
	require.True(t, ok)
	comments, err := w.QueryComments().All(ctx)
	require.NoError(t, err)
	assert.Len(t, comments, 0)

	ctxt := "Bar"
	c, err := mr.AddComment(ctx, models.CommentInput{
		ID:         w.ID,
		EntityType: "WORK_ORDER",
		Text:       ctxt,
	})
	require.NoError(t, err)
	assert.Equal(t, ctxt, c.Text)

	node, err = qr.Node(ctx, w.ID)
	require.NoError(t, err)
	w, ok = node.(*ent.WorkOrder)
	require.True(t, ok)
	comments, err = w.QueryComments().All(ctx)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
	assert.Equal(t, ctxt, comments[0].Text)
}

func TestAddWorkOrderNoDescription(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	name := "short_work_order"
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)

	assert.Equal(t, fetchedWorkOrder.Name, name)
	assert.Empty(t, fetchedWorkOrder.Description)
}

func TestFetchWorkOrder(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, wor := r.Mutation(), r.Query(), r.WorkOrder()
	name := "example_work_order"
	workOrder := createWorkOrder(ctx, t, *r, name)
	location := workOrder.QueryLocation().FirstX(ctx)

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type_name_1",
	})
	require.NoError(t, err)

	equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "equipment_name_1",
		Type:      equipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)
	assert.Equal(t, fetchedWorkOrder.Name, name)

	installedEquipment, err := wor.EquipmentToAdd(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Len(t, installedEquipment, 1)

	fetchedEquipment := installedEquipment[0]
	assert.Equal(t, equipment.ID, fetchedEquipment.ID)
	assert.Equal(t, equipment.Name, fetchedEquipment.Name)
}

func TestFetchWorkOrders(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	for i := 0; i < 2; i++ {
		_, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
			Name:            "example_work_order_" + strconv.Itoa(i),
			WorkOrderTypeID: woType.ID,
			LocationID:      &location.ID,
		})
		require.NoError(t, err)
	}

	trueVal := true
	types, err := qr.WorkOrders(ctx, nil, nil, nil, nil, &trueVal)
	require.NoError(t, err)
	assert.Len(t, types.Edges, 2)
}

func TestExecuteWorkOrderInstallEquipment(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()

	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type_name_1",
	})
	require.NoError(t, err)

	workOrderEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "work_order_equipment",
		Type:      equipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	assert.Equal(t, models.FutureStateInstall.String(), workOrderEquipment.FutureState)
	assert.Equal(t, workOrder.ID, workOrderEquipment.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedWorkOrderEquipment, err := qr.Equipment(ctx, workOrderEquipment.ID)
	require.NoError(t, err)
	assert.Empty(t, fetchedWorkOrderEquipment.FutureState)

	wo, err := fetchedWorkOrderEquipment.QueryWorkOrder().FirstID(ctx)
	require.Error(t, err)
	assert.Empty(t, wo)
}

func TestExecuteWorkOrderRemoveEquipment(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	position1 := createPosition()
	parentEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "parent_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)

	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "parent_equipment",
		Type:     parentEquipmentType.ID,
		Location: &location.ID,
	})
	assert.NoError(t, err)

	childEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "child_equipment_type",
	})
	assert.NoError(t, err)

	posDefID := parentEquipmentType.QueryPositionDefinitions().FirstXID(ctx)
	childEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child_equipment",
		Type:               childEquipmentType.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDefID,
	})
	assert.NoError(t, err)

	fetchedParentEquipment, err := qr.Equipment(ctx, parentEquipment.ID)
	assert.NoError(t, err)
	fetchedPosition := fetchedParentEquipment.QueryPositions().OnlyX(ctx)

	updatedPosition, err := mr.RemoveEquipmentFromPosition(ctx, fetchedPosition.ID, &workOrder.ID)
	require.NoError(t, err)
	assert.NotNil(t, updatedPosition.QueryParent().OnlyX(ctx)) // equipment isn't removed yet, only when workOrder is executed

	fetchedWorkOrderEquipment, err := qr.Equipment(ctx, childEquipment.ID)
	require.NoError(t, err)
	assert.Equal(t, models.FutureStateRemove.String(), fetchedWorkOrderEquipment.FutureState)
	assert.Equal(t, workOrder.ID, fetchedWorkOrderEquipment.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedRemovedWorkOrderEquipment, err := qr.Equipment(ctx, childEquipment.ID)
	require.NoError(t, err)
	assert.Nil(t, fetchedRemovedWorkOrderEquipment)

	fetchedParentEquipmentAfterExecution, err := qr.Equipment(ctx, parentEquipment.ID)
	assert.NoError(t, err)

	fetchedPositionAfterExecution := fetchedParentEquipmentAfterExecution.QueryPositions().OnlyX(ctx)
	_, err = mr.RemoveEquipmentFromPosition(ctx, fetchedPositionAfterExecution.ID, &workOrder.ID)
	require.NoError(t, err)
	eq, err := fetchedPositionAfterExecution.QueryAttachment().Only(ctx)
	require.Error(t, err)
	assert.Nil(t, eq)
}

func TestExecuteWorkOrderInstallLink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr := r.Mutation(), r.Query(), r.EquipmentPort()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	portInput := createPort()

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "equipment_type_name_1",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)

	equipmentA, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment1",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)
	equipmentB, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment2",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		WorkOrder: &workOrder.ID,
	})
	assert.NoError(t, err)

	assert.Equal(t, models.FutureStateInstall.String(), createdLink.FutureState)
	assert.Equal(t, workOrder.ID, createdLink.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedEquipment, _ := qr.Equipment(ctx, equipmentA.ID)
	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, _ := pr.Link(ctx, fetchedPort)
	assert.Equal(t, createdLink.ID, fetchedLink.ID)
	assert.Empty(t, fetchedLink.FutureState)
}

func TestExecuteWorkOrderRemoveLink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr := r.Mutation(), r.Query(), r.EquipmentPort()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	portInput := createPort()

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "equipment_type_name_1",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)

	equipmentA, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment1",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)
	equipmentB, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment2",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
	})
	assert.NoError(t, err)

	_, err = mr.RemoveLink(ctx, createdLink.ID, &workOrder.ID)
	assert.Nil(t, err)

	fetchedEquipment, err := qr.Equipment(ctx, equipmentA.ID)
	require.NoError(t, err)
	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, err := pr.Link(ctx, fetchedPort)
	require.NoError(t, err)

	assert.Equal(t, models.FutureStateRemove.String(), fetchedLink.FutureState)
	assert.Equal(t, workOrder.ID, fetchedLink.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedEquipmentAfterExecution, err := qr.Equipment(ctx, equipmentA.ID)
	require.NoError(t, err)
	fetchedPortAfterExecution, err := fetchedEquipmentAfterExecution.QueryPorts().Only(ctx)
	require.NoError(t, err)
	fetchedLinkAfterExecution, err := pr.Link(ctx, fetchedPortAfterExecution)
	assert.Nil(t, fetchedLinkAfterExecution)
	assert.NoError(t, err)
}

func TestExecuteWorkOrderInstallDependantEquipmentAndLink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr := r.Mutation(), r.Query(), r.EquipmentPort()

	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	portInput := createPort()
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "equipment_type_name_1",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	equipmentA, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment1",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)
	equipmentB, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment2",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})
	require.NoError(t, err)

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		WorkOrder: &workOrder.ID,
	})
	assert.Nil(t, err)

	assert.Equal(t, models.FutureStateInstall.String(), createdLink.FutureState)
	assert.Equal(t, workOrder.ID, createdLink.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedEquipment, _ := qr.Equipment(ctx, equipmentA.ID)

	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, _ := pr.Link(ctx, fetchedPort)
	assert.Equal(t, createdLink.ID, fetchedLink.ID)
	assert.Empty(t, fetchedLink.FutureState)
}

func TestExecuteWorkOrderInstallEquipmentMultilayer(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)

	position1 := createPosition()
	rootEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "root_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)
	rootEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "root_equipment",
		Type:     rootEquipmentType.ID,
		Location: &location.ID,
	})
	assert.NoError(t, err)

	var equipments []ent.Equipment
	equipments = append(equipments, *rootEquipment)
	for i := 0; i < 10; i++ {
		prevEquipmentPosition, err := equipments[i].QueryPositions().Only(ctx)
		require.NoError(t, err)
		defID := prevEquipmentPosition.QueryDefinition().OnlyXID(ctx)
		parentID := prevEquipmentPosition.QueryParent().OnlyXID(ctx)
		require.NoError(t, err)
		equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
			Name:               string(i),
			Type:               rootEquipmentType.ID,
			Parent:             &parentID,
			PositionDefinition: &defID,
			WorkOrder:          &workOrder.ID,
		})
		assert.NoError(t, err)
		equipments = append(equipments, *equipment)
	}

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	for _, equipment := range equipments {
		fetchedEquipment, err := qr.Equipment(ctx, equipment.ID)
		assert.NoError(t, err)
		assert.Empty(t, fetchedEquipment.FutureState)
		assert.Nil(t, fetchedEquipment.QueryWorkOrder().FirstX(ctx))
	}
}

func TestExecuteWorkOrderRemoveEquipmentMultilayer(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)
	position1 := createPosition()

	rootEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "root_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)
	rootEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "root_equipment",
		Type:     rootEquipmentType.ID,
		Location: &location.ID,
	})
	assert.NoError(t, err)

	var equipments []ent.Equipment
	equipments = append(equipments, *rootEquipment)
	for i := 0; i < 10; i++ {
		position, err := equipments[i].QueryPositions().Only(ctx)
		require.NoError(t, err)
		defID := position.QueryDefinition().OnlyXID(ctx)
		parentID := position.QueryParent().OnlyXID(ctx)
		equipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
			Name:               string(i),
			Type:               rootEquipmentType.ID,
			Parent:             &parentID,
			PositionDefinition: &defID,
		})
		assert.NoError(t, err)
		equipments = append(equipments, *equipment)
	}

	for i := 8; i >= 0; i-- {
		position, err := equipments[i].QueryPositions().Only(ctx)
		require.NoError(t, err)
		_, err = mr.RemoveEquipmentFromPosition(ctx, position.ID, &workOrder.ID)
		assert.NoError(t, err)
	}

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	for i, equipment := range equipments {
		fetchedEquipment, err := qr.Equipment(ctx, equipment.ID)
		if i == 0 {
			assert.NoError(t, err)
			assert.Empty(t, fetchedEquipment.FutureState)

			fetchedEquipmentWorkOrder, _ := fetchedEquipment.QueryWorkOrder().Only(ctx)
			assert.Nil(t, fetchedEquipmentWorkOrder)
		} else {
			assert.Nil(t, fetchedEquipment)
		}
	}
}

func TestExecuteWorkOrderInstallChildOnUninstalledParent(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	futureWorkOrder := createWorkOrder(ctx, t, *r, "example_work_order_2")
	location := workOrder.QueryLocation().FirstX(ctx)

	position1 := createPosition()
	parentEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "parent_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)
	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "parent_equipment",
		Type:      parentEquipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &futureWorkOrder.ID,
	})
	assert.NoError(t, err)

	childEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "child_equipment_type",
	})
	assert.NoError(t, err)

	posDefID := parentEquipmentType.QueryPositionDefinitions().FirstXID(ctx)
	require.NoError(t, err)
	childEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child_equipment",
		Type:               childEquipmentType.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDefID,
		WorkOrder:          &workOrder.ID,
	})
	assert.NoError(t, err)

	fetchedWorkOrderEquipment, _ := qr.Equipment(ctx, childEquipment.ID)
	assert.Equal(t, models.FutureStateInstall.String(), fetchedWorkOrderEquipment.FutureState)
	equipmentWorkOrder, err := fetchedWorkOrderEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, equipmentWorkOrder.ID)

	returnedWorkOrder, _ := executeWorkOrder(ctx, t, mr, *workOrder)
	assert.Nil(t, returnedWorkOrder)

	fetchedChildEquipment, _ := qr.Equipment(ctx, childEquipment.ID)
	equipmentWorkOrder, err = fetchedChildEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)

	// the child wasn't installed because the parent isn't installed yet
	assert.Equal(t, models.FutureStateInstall.String(), fetchedChildEquipment.FutureState)
	assert.Equal(t, workOrder.ID, equipmentWorkOrder.ID)
}

func TestExecuteWorkOrderInstallLinkOnUninstalledEquipment(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, pr := r.Mutation(), r.Query(), r.EquipmentPort()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	futureWorkOrder := createWorkOrder(ctx, t, *r, "example_work_order_2")
	location := workOrder.QueryLocation().FirstX(ctx)

	portInput := createPort()
	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  "equipment_type_name_1",
		Ports: []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)

	equipmentA, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "equipment1",
		Type:      equipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	equipmentB, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "equipment2",
		Type:      equipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &futureWorkOrder.ID,
	})
	require.NoError(t, err)

	portDef := equipmentType.QueryPortDefinitions().OnlyX(ctx)
	createdLink, err := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: equipmentA.ID, Port: portDef.ID},
			{Equipment: equipmentB.ID, Port: portDef.ID},
		},
		WorkOrder: &workOrder.ID,
	})
	assert.NoError(t, err)

	createLinkWorkOrder, err := createdLink.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)

	assert.Equal(t, models.FutureStateInstall.String(), createdLink.FutureState)
	assert.Equal(t, workOrder.ID, createLinkWorkOrder.ID)

	returnedWorkOrder, _ := executeWorkOrder(ctx, t, mr, *workOrder)
	require.Nil(t, returnedWorkOrder)

	fetchedEquipment, _ := qr.Equipment(ctx, equipmentA.ID)

	fetchedPort, err := fetchedEquipment.QueryPorts().Only(ctx)
	require.NoError(t, err)

	fetchedLink, _ := pr.Link(ctx, fetchedPort)
	assert.Equal(t, createdLink.ID, fetchedLink.ID)

	fetchedLinkWorkOrder, err := fetchedLink.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)

	// the link wasn't installed because equipmentB is not installed
	assert.Equal(t, models.FutureStateInstall.String(), fetchedLink.FutureState)
	assert.Equal(t, workOrder.ID, fetchedLinkWorkOrder.ID)

	_, err = fetchedEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)
}

func TestExecuteWorkOrderRemoveParentEquipment(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)

	position1 := createPosition()
	rootEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "root_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)
	rootEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "root_equipment",
		Type:     rootEquipmentType.ID,
		Location: &location.ID,
	})
	assert.NoError(t, err)

	parentEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "parent_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
	})
	assert.NoError(t, err)

	require.NoError(t, err)
	posDefID := rootEquipmentType.QueryPositionDefinitions().FirstXID(ctx)

	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "parent_equipment",
		Type:               parentEquipmentType.ID,
		Parent:             &rootEquipment.ID,
		PositionDefinition: &posDefID,
	})
	assert.NoError(t, err)

	childEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "child_equipment_type",
	})
	assert.NoError(t, err)
	posDefID = parentEquipmentType.QueryPositionDefinitions().FirstXID(ctx)
	require.NoError(t, err)
	childEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child_equipment",
		Type:               childEquipmentType.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDefID,
	})
	assert.NoError(t, err)
	assert.NotNil(t, childEquipment)

	fetchedRootEquipment, err := qr.Equipment(ctx, rootEquipment.ID)
	assert.NoError(t, err)
	fetchedPosition, err := fetchedRootEquipment.QueryPositions().Only(ctx)
	require.NoError(t, err)

	updatedPosition, err := mr.RemoveEquipmentFromPosition(ctx, fetchedPosition.ID, &workOrder.ID)
	require.NoError(t, err)

	attachedEquipment, err := updatedPosition.QueryAttachment().Only(ctx)
	require.NoError(t, err)

	assert.NotNil(t, attachedEquipment)

	fetchedWorkOrderEquipment, _ := qr.Equipment(ctx, parentEquipment.ID)

	assert.Equal(t, models.FutureStateRemove.String(), fetchedWorkOrderEquipment.FutureState)
	fetchedWorkOrderEquipmentWorkOrder, err := fetchedWorkOrderEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, fetchedWorkOrderEquipmentWorkOrder.ID)

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedParentWorkOrderEquipment, err := qr.Equipment(ctx, parentEquipment.ID)

	assert.Nil(t, err)

	assert.Nil(t, fetchedParentWorkOrderEquipment)
	fetchedPChildEquipment, _ := qr.Equipment(ctx, childEquipment.ID)
	assert.Nil(t, fetchedPChildEquipment)
}

func TestAddAndDeleteWorkOrderHyperlink(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, wor := r.Mutation(), r.WorkOrder()

	workOrderType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name: "work_order_type_name_1",
	})
	require.NoError(t, err)
	require.Equal(t, "work_order_type_name_1", workOrderType.Name)

	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "work_order_name_1",
		WorkOrderTypeID: workOrderType.ID,
	})
	require.NoError(t, err)
	require.Equal(t, workOrderType.ID, workOrder.QueryType().OnlyXID(ctx))

	category := "TSS"
	url := "http://some.url"
	displayName := "link to some url"
	hyperlink, err := mr.AddHyperlink(ctx, models.AddHyperlinkInput{
		EntityType:  models.ImageEntityWorkOrder,
		EntityID:    workOrder.ID,
		URL:         url,
		DisplayName: &displayName,
		Category:    &category,
	})
	require.NoError(t, err)
	require.Equal(t, url, hyperlink.URL, "verifying hyperlink url")
	require.Equal(t, displayName, hyperlink.Name, "verifying hyperlink display name")
	require.Equal(t, category, hyperlink.Category, "verifying 1st hyperlink category")

	hyperlinks, err := wor.Hyperlinks(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, hyperlinks, 1, "verifying has 1 hyperlink")

	deletedHyperlink, err := mr.DeleteHyperlink(ctx, hyperlink.ID)
	require.NoError(t, err)
	require.Equal(t, hyperlink.ID, deletedHyperlink.ID, "verifying return id of deleted hyperlink")

	hyperlinks, err = wor.Hyperlinks(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, hyperlinks, 0, "verifying no hyperlinks remained")
}

func TestDeleteWorkOrderWithAttachmentAndLinksAdded(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr := r.Mutation(), r.Query()
	workOrder := createWorkOrder(ctx, t, *r, "example_work_order")
	location := workOrder.QueryLocation().FirstX(ctx)

	position1 := createPosition()
	portInput := createPort()

	parentEquipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "parent_equipment_type",
		Positions: []*models.EquipmentPositionInput{position1},
		Ports:     []*models.EquipmentPortInput{&portInput},
	})
	require.NoError(t, err)
	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "parent_equipment",
		Type:      parentEquipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	posDefID := parentEquipmentType.QueryPositionDefinitions().FirstXID(ctx)

	childEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               "child_equipment",
		Type:               parentEquipmentType.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDefID,
		WorkOrder:          &workOrder.ID,
	})
	assert.NoError(t, err)

	connectedEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:      "connected_equipment",
		Type:      parentEquipmentType.ID,
		Location:  &location.ID,
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	portDef := parentEquipmentType.QueryPortDefinitions().OnlyX(ctx)

	_, err = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: parentEquipment.ID, Port: portDef.ID},
			{Equipment: connectedEquipment.ID, Port: portDef.ID},
		},
		WorkOrder: &workOrder.ID,
	})
	require.NoError(t, err)

	_, err = mr.RemoveWorkOrder(ctx, workOrder.ID)
	require.NoError(t, err)

	fetchedParentWorkOrderEquipment, _ := qr.Equipment(ctx, parentEquipment.ID)
	assert.Nil(t, fetchedParentWorkOrderEquipment)

	fetchedChildWorkOrderEquipment, _ := qr.Equipment(ctx, childEquipment.ID)
	assert.Nil(t, fetchedChildWorkOrderEquipment)

	fetchedConnectedWorkOrderEquipment, _ := qr.Equipment(ctx, connectedEquipment.ID)
	assert.Nil(t, fetchedConnectedWorkOrderEquipment)
}

func TestAddWorkOrderWithProperties(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr, qr, wr := r.Mutation(), r.Query(), r.WorkOrder()

	strPropType := models.PropertyTypeInput{
		Name: "str_prop",
		Type: "string",
	}
	strFixedValue := "FixedFoo"
	strFixedPropType := models.PropertyTypeInput{
		Name:               "str_fixed_prop",
		Type:               "string",
		IsInstanceProperty: pointer.ToBool(false),
		StringValue:        &strFixedValue,
	}
	intPropType := models.PropertyTypeInput{
		Name: "int_prop",
		Type: "int",
	}
	rangePropType := models.PropertyTypeInput{
		Name: "rng_prop",
		Type: "range",
	}
	propTypeInputs := []*models.PropertyTypeInput{&strPropType, &strFixedPropType, &intPropType, &rangePropType}
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type", Properties: propTypeInputs})
	require.NoError(t, err, "Adding location type")

	strValue := "Foo"
	strProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("str_prop")).OnlyXID(ctx),
		StringValue:    &strValue,
	}

	strFixedProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("str_fixed_prop")).OnlyXID(ctx),
		StringValue:    &strFixedValue,
	}

	intValue := 5
	intProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("int_prop")).OnlyXID(ctx),
		StringValue:    nil,
		IntValue:       &intValue,
	}
	fl1, fl2 := 5.5, 7.8
	rngProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("rng_prop")).OnlyXID(ctx),
		RangeFromValue: &fl1,
		RangeToValue:   &fl2,
	}
	propInputs := []*models.PropertyInput{&strProp, &strFixedProp, &intProp, &rngProp}
	wo, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "location_name_1",
		WorkOrderTypeID: woType.ID,
		Properties:      propInputs,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, wo.ID)
	require.NoError(t, err)
	fetchedWo, ok := node.(*ent.WorkOrder)
	require.True(t, ok)

	intFetchProp := fetchedWo.QueryProperties().Where(property.HasTypeWith(propertytype.Name("int_prop"))).OnlyX(ctx)
	require.Equal(t, intFetchProp.IntVal, *intProp.IntValue, "Comparing properties: int value")
	require.Equal(t, intFetchProp.QueryType().OnlyXID(ctx), intProp.PropertyTypeID, "Comparing properties: PropertyType value")

	strFetchProp := fetchedWo.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_prop"))).OnlyX(ctx)
	require.Equal(t, strFetchProp.StringVal, *strProp.StringValue, "Comparing properties: string value")
	require.Equal(t, strFetchProp.QueryType().OnlyXID(ctx), strProp.PropertyTypeID, "Comparing properties: PropertyType value")

	fixedStrFetchProp := fetchedWo.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_fixed_prop"))).OnlyX(ctx)
	require.Equal(t, fixedStrFetchProp.StringVal, *strFixedProp.StringValue, "Comparing properties: fixed string value")
	require.Equal(t, fixedStrFetchProp.QueryType().OnlyXID(ctx), strFixedProp.PropertyTypeID, "Comparing properties: PropertyType value")

	rngFetchProp := fetchedWo.QueryProperties().Where(property.HasTypeWith(propertytype.Name("rng_prop"))).OnlyX(ctx)
	require.Equal(t, rngFetchProp.RangeFromVal, *rngProp.RangeFromValue, "Comparing properties: range value")
	require.Equal(t, rngFetchProp.RangeToVal, *rngProp.RangeToValue, "Comparing properties: range value")
	require.Equal(t, rngFetchProp.QueryType().OnlyXID(ctx), rngProp.PropertyTypeID, "Comparing properties: PropertyType value")

	fetchedProps, err := wr.Properties(ctx, fetchedWo)
	require.NoError(t, err)
	require.Equal(t, len(propInputs), len(fetchedProps))

	failProp := models.PropertyInput{
		PropertyTypeID: "someFakeTypeID",
	}
	failEditInput := models.EditWorkOrderInput{
		ID:         wo.ID,
		Name:       "test",
		Properties: []*models.PropertyInput{&failProp},
	}
	_, err = mr.EditWorkOrder(ctx, failEditInput)
	require.Error(t, err, "editing Work Order instance property with wrong property type id")

	failProp2 := models.PropertyInput{
		ID:             &strFetchProp.ID,
		PropertyTypeID: intProp.PropertyTypeID,
	}
	failEditInput2 := models.EditWorkOrderInput{
		ID:         wo.ID,
		Name:       "test",
		Properties: []*models.PropertyInput{&failProp2},
	}
	_, err = mr.EditWorkOrder(ctx, failEditInput2)
	require.Error(t, err, "editing Work Order instance property when id and property type id mismach")

	newStrValue := "Foo"
	prop := models.PropertyInput{
		PropertyTypeID: strProp.PropertyTypeID,
		StringValue:    &newStrValue,
	}
	newWorkOrderName := "updated test"
	editInput := models.EditWorkOrderInput{
		ID:         wo.ID,
		Name:       newWorkOrderName,
		Properties: []*models.PropertyInput{&prop},
	}
	updatedP, err := mr.EditWorkOrder(ctx, editInput)
	require.NoError(t, err)

	updatedNode, err := qr.Node(ctx, updatedP.ID)
	require.NoError(t, err, "querying updated Work Order node")
	updatedWO, ok := updatedNode.(*ent.WorkOrder)
	require.True(t, ok, "casting updated Work Order instance")

	require.Equal(t, updatedWO.Name, newWorkOrderName, "Comparing updated Work Order name")

	fetchedProps, _ = wr.Properties(ctx, updatedWO)
	require.Equal(t, len(propInputs), len(fetchedProps), "number of properties should remain he same")

	updatedProp := updatedWO.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_prop"))).OnlyX(ctx)
	require.Equal(t, updatedProp.StringVal, *prop.StringValue, "Comparing updated properties: string value")
	require.Equal(t, updatedProp.QueryType().OnlyXID(ctx), prop.PropertyTypeID, "Comparing updated properties: PropertyType value")
}

func TestAddWorkOrderWithInvalidProperties(t *testing.T) {
	t.Skip("skipping test until mandatory props are added - T57858029")
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	mr := r.Mutation()
	latlongPropType := models.PropertyTypeInput{
		Name: "lat_long_prop",
		Type: "gps_location",
	}
	propTypeInputs := []*models.PropertyTypeInput{&latlongPropType}
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type", Properties: propTypeInputs})
	require.NoError(t, err)

	latlongProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("lat_long_prop")).OnlyXID(ctx),
	}
	propInputs := []*models.PropertyInput{&latlongProp}
	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "location_name_3",
		WorkOrderTypeID: woType.ID,
		Properties:      propInputs,
	})
	require.Error(t, err, "Adding work order instance with invalid lat-long prop")
}

func TestAddWorkOrderWithCheckList(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, wr := r.Mutation(), r.WorkOrder()
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name: "example_type_a",
	})
	require.NoError(t, err)
	indexValue := 1
	fooCL := models.CheckListItemInput{
		Title: "Foo",
		Type:  "simple",
		Index: &indexValue,
	}
	clInputs := []*models.CheckListItemInput{&fooCL}
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            longWorkOrderName,
		WorkOrderTypeID: woType.ID,
		CheckList:       clInputs,
	})
	require.NoError(t, err)
	cls := workOrder.QueryCheckListItems().AllX(ctx)
	require.Len(t, cls, 1)

	fooCLFetched := workOrder.QueryCheckListItems().Where(checklistitem.Type("simple")).OnlyX(ctx)
	require.Equal(t, "Foo", fooCLFetched.Title, "verifying check list name")

	cl, err := wr.CheckList(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, cl, 1)
}

func TestEditWorkOrderWithCheckList(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, wr := r.Mutation(), r.WorkOrder()
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name: "example_type_a",
	})
	require.NoError(t, err)
	indexValue := 1
	fooCL := models.CheckListItemInput{
		Title: "Foo",
		Type:  "simple",
		Index: &indexValue,
	}
	clInputs := []*models.CheckListItemInput{&fooCL}
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            longWorkOrderName,
		WorkOrderTypeID: woType.ID,
		CheckList:       clInputs,
	})
	require.NoError(t, err)

	barCL := models.CheckListItemInput{
		Title: "Bar",
		Type:  "simple",
		Index: &indexValue,
	}
	clInputs = []*models.CheckListItemInput{&barCL}
	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:        workOrder.ID,
		Name:      longWorkOrderName,
		CheckList: clInputs,
	})
	require.NoError(t, err)
	cls := workOrder.QueryCheckListItems().AllX(ctx)
	require.Len(t, cls, 1)

	fooCLFetched := workOrder.QueryCheckListItems().Where(checklistitem.Type("simple")).OnlyX(ctx)
	require.Equal(t, "Bar", fooCLFetched.Title, "verifying check list name")

	cl, err := wr.CheckList(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, cl, 1)
}

func TestEditWorkOrderLocation(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()
	name := longWorkOrderName
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)
	require.Equal(t, workOrder.QueryLocation().FirstXID(ctx), location.ID)

	location = createLocationWithName(ctx, t, *r, "location2")
	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:         workOrder.ID,
		Name:       workOrder.Name,
		LocationID: &location.ID,
	})
	require.NoError(t, err)
	require.Equal(t, workOrder.QueryLocation().FirstXID(ctx), location.ID)

	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:   workOrder.ID,
		Name: workOrder.Name,
	})
	require.NoError(t, err)
	locEx, err := workOrder.QueryLocation().Exist(ctx)
	require.NoError(t, err)
	require.False(t, locEx)
}

func TestTechnicianCheckinToWorkOrder(t *testing.T) {
	r, err := newTestResolver(t)
	require.NoError(t, err)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	w := createWorkOrder(ctx, t, *r, "Foo")
	require.NoError(t, err)

	w, err = mr.TechnicianWorkOrderCheckIn(ctx, w.ID)
	require.NoError(t, err)

	assert.Equal(t, w.Status, models.WorkOrderStatusPending.String())
	comments, err := w.QueryComments().All(ctx)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
}
