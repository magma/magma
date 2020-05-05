// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strconv"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent/user"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/checklistcategory"
	"github.com/facebookincubator/symphony/graph/ent/checklistitem"
	"github.com/facebookincubator/symphony/graph/ent/file"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/generated"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/99designs/gqlgen/client"
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
	var ownerID *int
	owner, _ := workOrder.QueryOwner().Only(ctx)
	if owner != nil {
		ownerID = &owner.ID
	}
	var assigneeID *int
	assignee, _ := workOrder.QueryAssignee().Only(ctx)
	if assignee != nil {
		assigneeID = &assignee.ID
	}
	_, err := mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerID:     ownerID,
		InstallDate: &workOrder.InstallDate,
		Status:      models.WorkOrderStatusDone,
		Priority:    models.WorkOrderPriorityNone,
		AssigneeID:  assigneeID,
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr, wr := r.Mutation(), r.Query(), r.WorkOrder()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	assigneeName := longWorkOrderAssignee
	assignee := viewer.MustGetOrCreateUser(ctx, assigneeName, user.RoleOWNER)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.NoError(t, err)
	require.False(t, workOrder.QueryAssignee().ExistX(ctx))

	var ownerID *int
	owner, _ := workOrder.QueryOwner().Only(ctx)
	if owner != nil {
		ownerID = &owner.ID
	}

	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerID:     ownerID,
		Status:      models.WorkOrderStatusPending,
		Priority:    models.WorkOrderPriorityNone,
		AssigneeID:  &assignee.ID,
	})
	require.NoError(t, err)

	node, err := qr.Node(ctx, workOrder.ID)
	require.NoError(t, err)
	fetchedWorkOrder, ok := node.(*ent.WorkOrder)
	require.True(t, ok)
	require.Equal(t, workOrder.QueryAssignee().OnlyXID(ctx), assignee.ID)

	fetchedWorkOrderType, err := wr.WorkOrderType(ctx, fetchedWorkOrder)
	require.NoError(t, err)
	assert.Equal(t, fetchedWorkOrderType.ID, woType.ID)
	assert.Equal(t, fetchedWorkOrderType.Name, woType.Name)
}

func TestAddWorkOrderWithDefaultAutomationOwner(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := ent.NewContext(context.Background(), r.client)
	v := viewer.NewAutomation(
		viewertest.DefaultTenant,
		viewertest.DefaultUser,
		viewertest.DefaultRole,
		viewer.WithFeatures(viewer.FeatureReadOnly, viewer.FeatureUserManagementDev))
	ctx = viewer.NewContext(ctx, v)
	ctx = authz.NewContext(ctx, authz.FullPermissions())
	mr := r.Mutation()
	name := longWorkOrderName
	description := longWorkOrderDesc
	location := createLocation(ctx, t, *r)
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type"})
	require.NoError(t, err)
	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: woType.ID,
		LocationID:      &location.ID,
	})
	require.Contains(t, err.Error(), "could not be executed in automation")
}

func TestAddWorkOrderInvalidType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()
	name := longWorkOrderName
	description := longWorkOrderDesc
	_, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name,
		Description:     &description,
		WorkOrderTypeID: 123,
		LocationID:      nil,
	})
	require.Error(t, err)
}

func TestEditInvalidWorkOrder(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	_, err := r.Mutation().EditWorkOrderType(ctx, models.EditWorkOrderTypeInput{
		ID:   234,
		Name: "foo",
	})
	require.Error(t, err)
}

func TestAddWorkOrderWithDescription(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	require.False(t, workOrder.QueryAssignee().ExistX(ctx))
	require.EqualValues(t, pri, workOrder.Priority)

	var ownerID *int
	owner, _ := workOrder.QueryOwner().Only(ctx)
	if owner != nil {
		ownerID = &owner.ID
	}

	input := models.EditWorkOrderInput{
		ID:          workOrder.ID,
		Name:        workOrder.Name,
		Description: &workOrder.Description,
		OwnerID:     ownerID,
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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

	var ownerID *int
	owner, _ := workOrder.QueryOwner().Only(ctx)
	if owner != nil {
		ownerID = &owner.ID
	}

	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:      workOrder.ID,
		Name:    workOrder.Name,
		OwnerID: ownerID,
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()
	w := createWorkOrder(ctx, t, *r, "Foo")

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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedWorkOrderNode, err := qr.Node(ctx, workOrderEquipment.ID)
	require.NoError(t, err)
	fetchedWorkOrderEquipment, ok := fetchedWorkOrderNode.(*ent.Equipment)
	require.True(t, ok)
	assert.Empty(t, fetchedWorkOrderEquipment.FutureState)

	wo, err := fetchedWorkOrderEquipment.QueryWorkOrder().FirstID(ctx)
	require.Error(t, err)
	assert.Empty(t, wo)
}

func TestExecuteWorkOrderRemoveEquipment(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedParentNode, err := qr.Node(ctx, parentEquipment.ID)
	assert.NoError(t, err)
	fetchedParentEquipment, ok := fetchedParentNode.(*ent.Equipment)
	assert.True(t, ok)
	fetchedPosition := fetchedParentEquipment.QueryPositions().OnlyX(ctx)

	updatedPosition, err := mr.RemoveEquipmentFromPosition(ctx, fetchedPosition.ID, &workOrder.ID)
	require.NoError(t, err)
	assert.NotNil(t, updatedPosition.QueryParent().OnlyX(ctx)) // equipment isn't removed yet, only when workOrder is executed

	fetchedWorkOrderNode, err := qr.Node(ctx, childEquipment.ID)
	require.NoError(t, err)
	fetchedWorkOrderEquipment, ok := fetchedWorkOrderNode.(*ent.Equipment)
	require.True(t, ok)
	assert.Equal(t, models.FutureStateRemove.String(), fetchedWorkOrderEquipment.FutureState)
	assert.Equal(t, workOrder.ID, fetchedWorkOrderEquipment.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedRemovedWorkOrderNode, err := qr.Node(ctx, childEquipment.ID)
	require.NoError(t, err)
	assert.Nil(t, fetchedRemovedWorkOrderNode)

	fetchedParentNodeAfterExecution, err := qr.Node(ctx, parentEquipment.ID)
	assert.NoError(t, err)
	fetchedParentEquipmentAfterExecution, ok := fetchedParentNodeAfterExecution.(*ent.Equipment)
	assert.True(t, ok)

	fetchedPositionAfterExecution := fetchedParentEquipmentAfterExecution.QueryPositions().OnlyX(ctx)
	_, err = mr.RemoveEquipmentFromPosition(ctx, fetchedPositionAfterExecution.ID, &workOrder.ID)
	require.NoError(t, err)
	eq, err := fetchedPositionAfterExecution.QueryAttachment().Only(ctx)
	require.Error(t, err)
	assert.Nil(t, eq)
}

func TestExecuteWorkOrderInstallLink(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedNode, err := qr.Node(ctx, equipmentA.ID)
	assert.NoError(t, err)
	fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
	assert.True(t, ok)
	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, _ := pr.Link(ctx, fetchedPort)
	assert.Equal(t, createdLink.ID, fetchedLink.ID)
	assert.Empty(t, fetchedLink.FutureState)
}

func TestExecuteWorkOrderRemoveLink(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedNode, err := qr.Node(ctx, equipmentA.ID)
	require.NoError(t, err)
	fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
	require.True(t, ok)
	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, err := pr.Link(ctx, fetchedPort)
	require.NoError(t, err)

	assert.Equal(t, models.FutureStateRemove.String(), fetchedLink.FutureState)
	assert.Equal(t, workOrder.ID, fetchedLink.QueryWorkOrder().OnlyXID(ctx))

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedNodeAfterExecution, err := qr.Node(ctx, equipmentA.ID)
	require.NoError(t, err)
	fetchedEquipmentAfterExecution, ok := fetchedNodeAfterExecution.(*ent.Equipment)
	require.True(t, ok)
	fetchedPortAfterExecution, err := fetchedEquipmentAfterExecution.QueryPorts().Only(ctx)
	require.NoError(t, err)
	fetchedLinkAfterExecution, err := pr.Link(ctx, fetchedPortAfterExecution)
	assert.Nil(t, fetchedLinkAfterExecution)
	assert.NoError(t, err)
}

func TestExecuteWorkOrderInstallDependantEquipmentAndLink(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedNode, err := qr.Node(ctx, equipmentA.ID)
	assert.NoError(t, err)
	fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
	assert.True(t, ok)

	fetchedPort := fetchedEquipment.QueryPorts().OnlyX(ctx)
	fetchedLink, _ := pr.Link(ctx, fetchedPort)
	assert.Equal(t, createdLink.ID, fetchedLink.ID)
	assert.Empty(t, fetchedLink.FutureState)
}

func TestExecuteWorkOrderInstallEquipmentMultilayer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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
		fetchedNode, err := qr.Node(ctx, equipment.ID)
		assert.NoError(t, err)
		fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
		assert.True(t, ok)
		assert.Empty(t, fetchedEquipment.FutureState)
		assert.Nil(t, fetchedEquipment.QueryWorkOrder().FirstX(ctx))
	}
}

func TestExecuteWorkOrderRemoveEquipmentMultilayer(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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
		fetchedNode, err := qr.Node(ctx, equipment.ID)
		if i == 0 {
			assert.NoError(t, err)
			fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
			assert.True(t, ok)
			assert.Empty(t, fetchedEquipment.FutureState)

			fetchedEquipmentWorkOrder, _ := fetchedEquipment.QueryWorkOrder().Only(ctx)
			assert.Nil(t, fetchedEquipmentWorkOrder)
		} else {
			assert.Nil(t, fetchedNode)
		}
	}
}

func TestExecuteWorkOrderInstallChildOnUninstalledParent(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedWorkOrderNode, err := qr.Node(ctx, childEquipment.ID)
	assert.NoError(t, err)
	fetchedWorkOrderEquipment, ok := fetchedWorkOrderNode.(*ent.Equipment)
	assert.True(t, ok)
	assert.Equal(t, models.FutureStateInstall.String(), fetchedWorkOrderEquipment.FutureState)
	equipmentWorkOrder, err := fetchedWorkOrderEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, equipmentWorkOrder.ID)

	returnedWorkOrder, _ := executeWorkOrder(ctx, t, mr, *workOrder)
	assert.Nil(t, returnedWorkOrder)

	fetchedChildNode, err := qr.Node(ctx, childEquipment.ID)
	require.NoError(t, err)
	fetchedChildEquipment, ok := fetchedChildNode.(*ent.Equipment)
	require.True(t, ok)
	equipmentWorkOrder, err = fetchedChildEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)

	// the child wasn't installed because the parent isn't installed yet
	assert.Equal(t, models.FutureStateInstall.String(), fetchedChildEquipment.FutureState)
	assert.Equal(t, workOrder.ID, equipmentWorkOrder.ID)
}

func TestExecuteWorkOrderInstallLinkOnUninstalledEquipment(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedNode, err := qr.Node(ctx, equipmentA.ID)
	require.NoError(t, err)
	fetchedEquipment, ok := fetchedNode.(*ent.Equipment)
	require.True(t, ok)

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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedRootNode, err := qr.Node(ctx, rootEquipment.ID)
	assert.NoError(t, err)
	fetchedRootEquipment, ok := fetchedRootNode.(*ent.Equipment)
	assert.True(t, ok)
	fetchedPosition, err := fetchedRootEquipment.QueryPositions().Only(ctx)
	require.NoError(t, err)

	updatedPosition, err := mr.RemoveEquipmentFromPosition(ctx, fetchedPosition.ID, &workOrder.ID)
	require.NoError(t, err)

	attachedEquipment, err := updatedPosition.QueryAttachment().Only(ctx)
	require.NoError(t, err)

	assert.NotNil(t, attachedEquipment)

	fetchedWorkOrderNode, err := qr.Node(ctx, parentEquipment.ID)
	assert.NoError(t, err)
	fetchedWorkOrderEquipment, ok := fetchedWorkOrderNode.(*ent.Equipment)
	assert.True(t, ok)

	assert.Equal(t, models.FutureStateRemove.String(), fetchedWorkOrderEquipment.FutureState)
	fetchedWorkOrderEquipmentWorkOrder, err := fetchedWorkOrderEquipment.QueryWorkOrder().Only(ctx)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, fetchedWorkOrderEquipmentWorkOrder.ID)

	returnedWorkOrder, err := executeWorkOrder(ctx, t, mr, *workOrder)
	require.NoError(t, err)
	assert.Equal(t, workOrder.ID, returnedWorkOrder.ID)

	fetchedParentWorkOrderNode, err := qr.Node(ctx, parentEquipment.ID)

	assert.Nil(t, err)

	assert.Nil(t, fetchedParentWorkOrderNode)
	fetchedPChildNode, _ := qr.Node(ctx, childEquipment.ID)
	assert.Nil(t, fetchedPChildNode)
}

func TestAddAndDeleteWorkOrderHyperlink(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	fetchedParentWorkOrderNode, _ := qr.Node(ctx, parentEquipment.ID)
	assert.Nil(t, fetchedParentWorkOrderNode)

	fetchedChildWorkOrderNode, _ := qr.Node(ctx, childEquipment.ID)
	assert.Nil(t, fetchedChildWorkOrderNode)

	fetchedConnectedWorkOrderNode, _ := qr.Node(ctx, connectedEquipment.ID)
	assert.Nil(t, fetchedConnectedWorkOrderNode)
}

func TestAddWorkOrderWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

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

	failProp := models.PropertyInput{PropertyTypeID: -1}
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
	_, err = mr.EditWorkOrder(ctx, editInput)
	require.NoError(t, err)

	newStrFixedValue := "updated FixedFoo"
	newStrFixedProp := models.PropertyInput{
		PropertyTypeID: strFixedProp.PropertyTypeID,
		StringValue:    &newStrFixedValue,
	}
	editFixedPropInput := models.EditWorkOrderInput{
		ID:         wo.ID,
		Name:       newWorkOrderName,
		Properties: []*models.PropertyInput{&newStrFixedProp},
	}
	updatedP, err := mr.EditWorkOrder(ctx, editFixedPropInput)
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

	notUpdatedFixedProp := updatedWO.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_fixed_prop"))).OnlyX(ctx)
	require.Equal(t, notUpdatedFixedProp.StringVal, *strFixedProp.StringValue, "Comparing not changed fixed property: string value")
	require.Equal(t, notUpdatedFixedProp.QueryType().OnlyXID(ctx), strFixedProp.PropertyTypeID, "Comparing updated properties: PropertyType value")
}

func TestAddWorkOrderWithInvalidProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr := r.Mutation()
	latlongPropType := models.PropertyTypeInput{
		Name:        "lat_long_prop",
		Type:        "gps_location",
		IsMandatory: pointer.ToBool(true),
	}
	propTypeInputs := []*models.PropertyTypeInput{&latlongPropType}
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type", Properties: propTypeInputs})
	require.NoError(t, err)

	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "should_fail",
		WorkOrderTypeID: woType.ID,
	})
	require.Error(t, err, "Adding work order instance with missing mandatory properties")

	latlongProp := models.PropertyInput{
		PropertyTypeID: woType.QueryPropertyTypes().Where(propertytype.Name("lat_long_prop")).OnlyXID(ctx),
		LatitudeValue:  pointer.ToFloat64(32.6),
		LongitudeValue: pointer.ToFloat64(34.7),
	}
	propInputs := []*models.PropertyInput{&latlongProp}
	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "location_name_3",
		WorkOrderTypeID: woType.ID,
		Properties:      propInputs,
	})
	require.NoError(t, err)

	// mandatory is deleted - should be ok
	latlongPropType.IsDeleted = pointer.ToBool(true)
	woType, err = mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type2", Properties: propTypeInputs})
	deletedPropType := woType.QueryPropertyTypes().OnlyX(ctx)
	require.NoError(t, err)

	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "should_not_fail",
		WorkOrderTypeID: woType.ID,
	})
	require.NoError(t, err, "Adding work order instance of template with mandatory properties but deleted")

	_, err = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "should_fail",
		WorkOrderTypeID: woType.ID,
		Properties: []*models.PropertyInput{{
			PropertyTypeID: deletedPropType.ID,
			StringValue:    pointer.ToString("new"),
		}},
	})
	require.Errorf(t, err, "deleted property types")

	// not mandatory props
	notMandatoryProp := &models.PropertyTypeInput{
		Name: "lat_long_prop",
		Type: "gps_location",
	}
	props := []*models.PropertyTypeInput{notMandatoryProp}

	woType, err = mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type3", Properties: props})
	require.NoError(t, err)

	wo, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            "should_pass",
		WorkOrderTypeID: woType.ID,
	})
	require.NoError(t, err, "Adding work order instance with missing mandatory properties")
	require.Len(t, wo.QueryProperties().AllX(ctx), 1)
}

func TestAddWorkOrderWithCheckListCategory(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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

	barCLC := models.CheckListCategoryInput{
		Title:     "Bar",
		CheckList: clInputs,
	}

	clcInputs := []*models.CheckListCategoryInput{&barCLC}
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:                longWorkOrderName,
		WorkOrderTypeID:     woType.ID,
		CheckListCategories: clcInputs,
	})
	require.NoError(t, err)
	cls := workOrder.QueryCheckListCategories().AllX(ctx)
	require.Len(t, cls, 1)

	barCLCFetched := workOrder.QueryCheckListCategories().Where(checklistcategory.Title("Bar")).OnlyX(ctx)
	fooCLFetched := barCLCFetched.QueryCheckListItems().Where(checklistitem.Type("simple")).OnlyX(ctx)
	require.Equal(t, "Foo", fooCLFetched.Title, "verifying checklist name")

	clcs, err := wr.CheckListCategories(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, clcs, 1)

	clcr := r.CheckListCategory()
	cl, err := clcr.CheckList(ctx, barCLCFetched)
	require.NoError(t, err)
	require.Len(t, cl, 1)
}

func TestEditWorkOrderWithCheckListCategory(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, wr := r.Mutation(), r.WorkOrder()
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name: "example_type_a",
	})
	require.NoError(t, err)
	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            longWorkOrderName,
		WorkOrderTypeID: woType.ID,
	})
	require.NoError(t, err)
	indexValue := 1
	fooCL := models.CheckListItemInput{
		Title: "Foo",
		Type:  "simple",
		Index: &indexValue,
	}
	clInputs := []*models.CheckListItemInput{&fooCL}

	barCLC := models.CheckListCategoryInput{
		Title:     "Bar",
		CheckList: clInputs,
	}

	clcInputs := []*models.CheckListCategoryInput{&barCLC}
	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:                  workOrder.ID,
		Name:                longWorkOrderName,
		CheckListCategories: clcInputs,
	})
	require.NoError(t, err)
	cls := workOrder.QueryCheckListCategories().AllX(ctx)
	require.Len(t, cls, 1)

	barCLCFetched := workOrder.QueryCheckListCategories().Where(checklistcategory.Title("Bar")).OnlyX(ctx)
	fooCLFetched := barCLCFetched.QueryCheckListItems().Where(checklistitem.Type("simple")).OnlyX(ctx)
	require.Equal(t, "Foo", fooCLFetched.Title, "verifying checklist name")

	clcs, err := wr.CheckListCategories(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, clcs, 1)

	clcr := r.CheckListCategory()
	cl, err := clcr.CheckList(ctx, barCLCFetched)
	require.NoError(t, err)
	require.Len(t, cl, 1)
}

func TestAddWorkOrderWithCheckList(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	enumValues := "val1,val2,val3"
	selectionMode := models.CheckListItemEnumSelectionModeMultiple
	selectedValues := "val2,val3"
	multiCL := models.CheckListItemInput{
		Title:              "Multi",
		Type:               "enum",
		Index:              pointer.ToInt(2),
		EnumValues:         &enumValues,
		EnumSelectionMode:  &selectionMode,
		SelectedEnumValues: &selectedValues,
	}
	yesNoResponse := models.YesNoResponse("YES")
	yesNoCL := models.CheckListItemInput{
		Title:         "Yes/No",
		Type:          "yes_no",
		Index:         pointer.ToInt(3),
		YesNoResponse: &yesNoResponse,
	}
	clInputs = []*models.CheckListItemInput{&barCL, &multiCL, &yesNoCL}
	workOrder, err = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:        workOrder.ID,
		Name:      longWorkOrderName,
		CheckList: clInputs,
	})
	require.NoError(t, err)
	cls := workOrder.QueryCheckListItems().AllX(ctx)
	require.Len(t, cls, 3)

	fooCLFetched := workOrder.QueryCheckListItems().Where(checklistitem.Type("simple")).OnlyX(ctx)
	require.Equal(t, "Bar", fooCLFetched.Title, "verifying check list name")

	multiCLFetched := workOrder.QueryCheckListItems().Where(checklistitem.Type("enum")).OnlyX(ctx)
	require.Equal(t, "Multi", multiCLFetched.Title)
	require.Equal(t, selectionMode.String(), multiCLFetched.EnumSelectionMode)
	require.Equal(t, enumValues, multiCLFetched.EnumValues)
	require.Equal(t, selectedValues, multiCLFetched.SelectedEnumValues)

	yesNoCLFetched := workOrder.QueryCheckListItems().Where(checklistitem.Type("yes_no")).OnlyX(ctx)
	require.Equal(t, "YES", yesNoCLFetched.YesNoVal.String())

	cl, err := wr.CheckList(ctx, workOrder)
	require.NoError(t, err)
	require.Len(t, cl, 3)
}

func TestEditCheckListItemFiles(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()
	woType, err := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name: "example_type",
	})
	require.NoError(t, err)
	indexValue := 0

	workOrder, err := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            longWorkOrderName,
		WorkOrderTypeID: woType.ID,
		CheckListCategories: []*models.CheckListCategoryInput{{
			Title: "Category1",
			CheckList: []*models.CheckListItemInput{{
				Title: "Files",
				Type:  "files",
				Index: &indexValue,
				Files: []*models.FileInput{
					{
						FileName: "File1",
						StoreKey: "File1StoreKey",
					},
					{
						FileName: "File2",
						StoreKey: "File2StoreKey",
					},
				},
			}},
		}},
	})
	require.NoError(t, err)

	queriedFiles, err := workOrder.QueryCheckListCategories().QueryCheckListItems().QueryFiles().All(ctx)
	require.NoError(t, err)
	require.Len(t, queriedFiles, 2)
	file1, err := workOrder.QueryCheckListCategories().QueryCheckListItems().QueryFiles().Where(file.Name("File1")).Only(ctx)
	require.NoError(t, err)
	require.NotNil(t, file1)

	checklistCategoryID, err := workOrder.QueryCheckListCategories().OnlyID(ctx)
	require.NoError(t, err)
	filesItemID, err := workOrder.QueryCheckListCategories().QueryCheckListItems().OnlyID(ctx)
	require.NoError(t, err)
	updatedWorkOrder, err := mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:   workOrder.ID,
		Name: longWorkOrderName,
		CheckListCategories: []*models.CheckListCategoryInput{{
			ID: &checklistCategoryID,
			CheckList: []*models.CheckListItemInput{{
				ID:    &filesItemID,
				Title: "Files",
				Type:  "files",
				Index: &indexValue,
				Files: []*models.FileInput{
					{
						ID:       &file1.ID,
						FileName: "File1 Renamed",
						StoreKey: "File1StoreKey",
					},
					{
						FileName: "File3",
						StoreKey: "File3StoreKey",
					},
				},
			}},
		}},
	})
	require.NoError(t, err)

	queriedUpdatedFiles, err := updatedWorkOrder.QueryCheckListCategories().QueryCheckListItems().QueryFiles().All(ctx)
	require.NoError(t, err)
	require.Len(t, queriedUpdatedFiles, 2)

	file2Exists, err := workOrder.QueryCheckListCategories().QueryCheckListItems().QueryFiles().Where(file.Name("File2")).Exist(ctx)
	require.NoError(t, err)
	require.False(t, file2Exists)

	updatedFile1Exists, err := workOrder.QueryCheckListItems().QueryFiles().Where(file.Name("File1 Renamed")).Exist(ctx)
	require.NoError(t, err)
	require.False(t, updatedFile1Exists)

	file3Exists, err := workOrder.QueryCheckListCategories().QueryCheckListItems().QueryFiles().Where(file.Name("File3")).Exist(ctx)
	require.NoError(t, err)
	require.True(t, file3Exists)
}

func TestEditWorkOrderLocation(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
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
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()

	w := createWorkOrder(ctx, t, *r, "Foo")
	w, err := mr.TechnicianWorkOrderCheckIn(ctx, w.ID)
	require.NoError(t, err)

	assert.Equal(t, w.Status, models.WorkOrderStatusPending.String())
	comments, err := w.QueryComments().All(ctx)
	require.NoError(t, err)
	assert.Len(t, comments, 1)
}

func TestTechnicianUploadDataToWorkOrder(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()
	c := newGraphClient(t, r)

	wo := createWorkOrder(ctx, t, *r, "Foo")
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	mimeType := "image/jpeg"
	sizeInBytes := 120
	wo, err := mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:         wo.ID,
		Name:       longWorkOrderName,
		AssigneeID: &u.ID,
		CheckListCategories: []*models.CheckListCategoryInput{{
			Title: "Bar",
			CheckList: []*models.CheckListItemInput{{
				Title:   "Foo",
				Type:    models.CheckListItemTypeSimple,
				Index:   pointer.ToInt(0),
				Checked: pointer.ToBool(false),
			},
				{
					Title: "CellScan",
					Type:  models.CheckListItemTypeCellScan,
					Index: pointer.ToInt(1),
				}, {
					Title: "Files",
					Type:  models.CheckListItemTypeFiles,
					Index: pointer.ToInt(2),
					Files: []*models.FileInput{
						{
							StoreKey:    "StoreKeyAlreadyIn",
							FileName:    "FileAlreadyInWorkOrder",
							SizeInBytes: &sizeInBytes,
							MimeType:    &mimeType,
						},
						{
							StoreKey:    "StoreKeyToBeDeleted",
							FileName:    "FileToBeDeleted",
							SizeInBytes: &sizeInBytes,
							MimeType:    &mimeType,
						},
					},
				}},
		}},
	})
	require.NoError(t, err)

	fooID, err := wo.QueryCheckListCategories().QueryCheckListItems().Where(checklistitem.TypeEQ("simple")).OnlyID(ctx)
	require.NoError(t, err)
	cellScanID, err := wo.QueryCheckListCategories().QueryCheckListItems().Where(checklistitem.TypeEQ("cell_scan")).OnlyID(ctx)
	require.NoError(t, err)
	filesID, err := wo.QueryCheckListCategories().QueryCheckListItems().Where(checklistitem.TypeEQ("files")).OnlyID(ctx)
	require.NoError(t, err)
	fileToKeepID, err := wo.QueryCheckListCategories().QueryCheckListItems().Where(checklistitem.TypeEQ("files")).QueryFiles().Where(file.StoreKey("StoreKeyAlreadyIn")).OnlyID(ctx)
	require.NoError(t, err)
	techInput := models.TechnicianWorkOrderUploadInput{
		WorkOrderID: wo.ID,
		Checklist: []*models.TechnicianCheckListItemInput{
			{
				ID:      fooID,
				Checked: pointer.ToBool(true),
			},
			{
				ID: cellScanID,
				CellData: []*models.SurveyCellScanData{{
					NetworkType:    models.CellularNetworkTypeLte,
					SignalStrength: -93,
				}},
			},
			{
				ID: filesID,
				FilesData: []*models.FileInput{ // Adding one new file, updating an existing file, deleting a file
					{
						StoreKey:    "StoreKeyToAdd",
						FileName:    "FileNameToAdd",
						SizeInBytes: &sizeInBytes,
						MimeType:    &mimeType,
					},
					{
						ID:          &fileToKeepID,
						StoreKey:    "StoreKeyAlreadyIn",
						FileName:    "FileAlreadyInWorkOrder",
						SizeInBytes: &sizeInBytes,
						MimeType:    &mimeType,
					},
				},
			},
		},
	}

	var rsp struct {
		TechnicianWorkOrderUploadData struct {
			ID                  string
			CheckListCategories []struct {
				CheckList []struct {
					ID       string
					Type     models.CheckListItemType
					Checked  *bool
					CellData []struct {
						NetworkType    string
						SignalStrength int
					}
					Files []struct {
						StoreKey    string
						FileName    string
						SizeInBytes int
						MimeType    string
						FileType    models.FileType
					}
				}
			}
		}
	}
	c.MustPost(
		`mutation($input: TechnicianWorkOrderUploadInput!) {
			technicianWorkOrderUploadData(input: $input) {
				id
				checkListCategories {
					checkList {
						id
						type
						checked
						cellData {
							networkType
							signalStrength
						}
						files {
							storeKey
							fileName
							sizeInBytes
							mimeType
							fileType
						}
					}
				}
			}
		}`,
		&rsp,
		client.Var("input", techInput),
	)

	require.Len(t, rsp.TechnicianWorkOrderUploadData.CheckListCategories, 1)
	require.Len(t, rsp.TechnicianWorkOrderUploadData.CheckListCategories[0].CheckList, 3)

	for _, item := range rsp.TechnicianWorkOrderUploadData.CheckListCategories[0].CheckList {
		switch item.Type {
		case models.CheckListItemTypeSimple:
			require.True(t, *item.Checked)
		case models.CheckListItemTypeCellScan:
			require.Equal(t, models.CellularNetworkTypeLte.String(), item.CellData[0].NetworkType)
			require.Equal(t, -93, item.CellData[0].SignalStrength)
		case models.CheckListItemTypeFiles:
			require.Equal(t, 2, len(item.Files))

			require.Equal(t, "StoreKeyAlreadyIn", item.Files[0].StoreKey)
			require.Equal(t, 120, item.Files[0].SizeInBytes)
			require.Equal(t, models.FileTypeImage, item.Files[0].FileType)

			require.Equal(t, "StoreKeyToAdd", item.Files[1].StoreKey)
		}

	}

}
