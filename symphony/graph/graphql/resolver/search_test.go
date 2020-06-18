// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/pkg/ent"
	"github.com/facebookincubator/symphony/pkg/ent/user"
	"github.com/facebookincubator/symphony/pkg/viewer"
	"github.com/facebookincubator/symphony/pkg/viewer/viewertest"

	"github.com/99designs/gqlgen/client"
	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	woCountQuery = `query($filters: [WorkOrderFilterInput!]!) {
		workOrderSearch(filters:$filters) {
			count
		}
	}`

	woAllQuery = `query($filters: [WorkOrderFilterInput!]!) {
		workOrderSearch(filters:$filters) {
			count
			workOrders {
				id
			}
		}
	}`
)

type (
	equipmentSearchDataModels struct {
		locType1  int
		locType2  int
		loc1      int
		loc2      int
		equType   int
		equ2ExtID string
	}

	woSearchDataModels struct {
		loc1        int
		woType1     int
		assignee1   int
		wo1         int
		owner       int
		installDate time.Time
	}

	woSearchResult struct {
		WorkOrderSearch struct {
			Count      int
			WorkOrders []struct {
				ID string
			}
		}
	}
)

func prepareEquipmentData(ctx context.Context, r *TestResolver, name string, props []*models.PropertyInput) equipmentSearchDataModels {
	mr := r.Mutation()
	locType1, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "loc_type1",
	})
	locType2, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "loc_type2",
	})

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst1",
		Type: locType1.ID,
	})
	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst2",
		Type: locType2.ID,
	})
	propType := models.PropertyTypeInput{
		Name: "Owner",
		Type: models.PropertyKindString,
	}
	equType, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       name + "eq_type",
		Properties: []*models.PropertyTypeInput{&propType},
	})
	if len(props) != 0 {
		props[0].PropertyTypeID = equType.QueryPropertyTypes().OnlyXID(ctx)
	}
	_, _ = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:       name + "eq_inst1",
		Type:       equType.ID,
		Location:   &loc1.ID,
		Properties: props,
	})
	extID := name + "123"
	equ2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:       name + "eq_inst2",
		Type:       equType.ID,
		Location:   &loc2.ID,
		Properties: props,
		ExternalID: &extID,
	})
	return equipmentSearchDataModels{
		locType1.ID,
		locType2.ID,
		loc1.ID,
		loc2.ID,
		equType.ID,
		equ2.ExternalID,
	}
}

func prepareWOData(ctx context.Context, r *TestResolver, name string) woSearchDataModels {
	mr := r.Mutation()
	locType1, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "loc_type1",
	})
	locType2, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "loc_type2",
	})

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst1",
		Type: locType1.ID,
	})
	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst2",
		Type: locType2.ID,
	})

	woType1, _ := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "wo_type_a"})
	woType2, _ := mr.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "wo_type_b"})
	assigneeName1 := "user1@fb.com"
	assigneeName2 := "user2@fb.com"
	assignee1 := viewer.MustGetOrCreateUser(ctx, assigneeName1, user.RoleOWNER)
	assignee2 := viewer.MustGetOrCreateUser(ctx, assigneeName2, user.RoleOWNER)
	desc := "random description"

	wo1, _ := mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name + "wo_1",
		Description:     &desc,
		WorkOrderTypeID: woType1.ID,
		LocationID:      &loc1.ID,
		AssigneeID:      &assignee1.ID,
	})
	_, _ = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name + "wo_2",
		Description:     &desc,
		WorkOrderTypeID: woType1.ID,
		AssigneeID:      &assignee1.ID,
	})
	_, _ = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name + "wo_3",
		Description:     &desc,
		WorkOrderTypeID: woType2.ID,
		LocationID:      &loc1.ID,
		AssigneeID:      &assignee2.ID,
	})
	_, _ = mr.AddWorkOrder(ctx, models.AddWorkOrderInput{
		Name:            name + "wo_4",
		Description:     &desc,
		WorkOrderTypeID: woType2.ID,
		LocationID:      &loc2.ID,
	})

	installDate := time.Now()
	ownerName := "owner"
	owner := viewer.MustGetOrCreateUser(ctx, ownerName, user.RoleOWNER)
	_, _ = mr.EditWorkOrder(ctx, models.EditWorkOrderInput{
		ID:          wo1.ID,
		Name:        wo1.Name,
		OwnerID:     &owner.ID,
		InstallDate: &installDate,
		Status:      models.WorkOrderStatusDone,
		Priority:    models.WorkOrderPriorityHigh,
		LocationID:  &loc1.ID,
		AssigneeID:  &assignee1.ID,
	})

	return woSearchDataModels{
		loc1.ID,
		woType1.ID,
		assignee1.ID,
		wo1.ID,
		owner.ID,
		installDate,
	}
}

func TestSearchEquipmentByName(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	c := r.GraphClient()

	mr := r.Mutation()

	locationType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "location_type",
	})
	location, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loCaTIon_name",
		Type: locationType.ID,
	})

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type",
	})
	require.NoError(t, err)
	e, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "EqUiPmEnT",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})

	var rsp struct {
		SearchForNode struct {
			Edges []struct {
				Node struct {
					Name string
				}
			}
		}
	}
	c.MustPost(
		`query($name: String!) { searchForNode(name: $name, first: 10) { edges { node { ... on Equipment { name } } } } }`,
		&rsp,
		client.Var("name", "equip"),
	)
	require.Len(t, rsp.SearchForNode.Edges, 1)
	assert.Equal(t, e.Name, rsp.SearchForNode.Edges[0].Node.Name)
	require.NoError(t, err)

	_, _ = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipMENT",
		Type:     equipmentType.ID,
		Location: &location.ID,
	})

	c.MustPost(
		`query($name: String!) { searchForNode(name: $name, first: 10) { edges { node { ... on Equipment { name } } } } }`,
		&rsp,
		client.Var("name", "ment"),
	)
	require.Len(t, rsp.SearchForNode.Edges, 2)

	c.MustPost(
		`query($name: String!) { searchForNode(name: $name, first: 10) { edges { node { ... on Location { name } } } } }`,
		&rsp,
		client.Var("name", "cation"),
	)
	require.Len(t, rsp.SearchForNode.Edges, 1)
	assert.Equal(t, location.Name, rsp.SearchForNode.Edges[0].Node.Name)
	require.NoError(t, err)
}

func TestEquipmentSearch(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	owner := "Ted"
	prop := models.PropertyInput{
		StringValue: &owner,
	}

	model1 := prepareEquipmentData(ctx, r, "A", []*models.PropertyInput{&prop}) // two locations on same type. each has one equipment.
	model2 := prepareEquipmentData(ctx, r, "B", nil)                            // two locations on same type. each has one equipment.
	/*
		helper: data now is of type:
		loctype1:
			inst1
				eq1 (typeA, name "A_") + prop
			inst2
				eq2 (typeA, name "A_") + prop
		loctype2:
			inst1
				eq1 (typeB, name "B_")
			inst2
				eq2 (typeB, name "B_")
	*/
	qr := r.Query()
	limit := 100
	all, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{}, &limit)
	require.NoError(t, err)
	require.Len(t, all.Equipment, 4)
	require.Equal(t, all.Count, 4)

	maxDepth := 5
	f1 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{model1.loc1, model2.loc1},
		MaxDepth:   &maxDepth,
	}
	res1, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Equipment, 2)
	require.Equal(t, res1.Count, 2)

	f2 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeEquipmentType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{model1.equType},
		MaxDepth:   &maxDepth,
	}
	res2, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f1, &f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Equipment, 1)
	require.Equal(t, res2.Count, 1)

	fetchedPropType := res2.Equipment[0].QueryType().QueryPropertyTypes().OnlyX(ctx)
	f3 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeProperty,
		Operator:   models.FilterOperatorIs,
		PropertyValue: &models.PropertyTypeInput{
			Name:        fetchedPropType.Name,
			Type:        models.PropertyKind(fetchedPropType.Type),
			StringValue: &owner,
		},
		MaxDepth: &maxDepth,
	}
	res3, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f3}, &limit)
	require.NoError(t, err)

	require.Len(t, res3.Equipment, 2)
	require.Equal(t, res3.Count, 2)

	subst := "inst1"
	f4 := models.EquipmentFilterInput{
		FilterType:  models.EquipmentFilterTypeEquipInstName,
		Operator:    models.FilterOperatorContains,
		StringValue: &subst,
		MaxDepth:    &maxDepth,
	}
	res4, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f3, &f4}, &limit)
	require.NoError(t, err)
	require.Len(t, res4.Equipment, 1)
	require.Equal(t, res4.Count, 1)

	f5 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{model2.loc1},
		MaxDepth:   &maxDepth,
	}
	res5, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f3, &f4, &f5}, &limit)
	require.NoError(t, err)
	require.Empty(t, res5.Equipment)
	require.Zero(t, res5.Count)

	f6 := models.EquipmentFilterInput{
		FilterType:  models.EquipmentFilterTypeEquipInstExternalID,
		Operator:    models.FilterOperatorIs,
		StringValue: &model1.equ2ExtID,
		MaxDepth:    &maxDepth,
	}
	res6, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f6}, &limit)
	require.NoError(t, err)
	require.Len(t, res6.Equipment, 1)
	require.Equal(t, res6.Count, 1)
}

func TestUnsupportedEquipmentSearch(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	qr := r.Query()
	limit := 100

	maxDepth := 5
	f := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeLocationInst,
		Operator:   models.FilterOperatorContains,
		MaxDepth:   &maxDepth,
	}
	_, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f}, &limit)
	require.Error(t, err)

	f = models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeProperty,
		Operator:   models.FilterOperatorContains,
		MaxDepth:   &maxDepth,
	}
	_, err = qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f}, &limit)
	require.Error(t, err)

	f = models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeEquipmentType,
		Operator:   models.FilterOperatorContains,
		MaxDepth:   &maxDepth,
	}
	_, err = qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f}, &limit)
	require.Error(t, err)
}

func TestQueryEquipmentPossibleProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()

	namePropType := models.PropertyTypeInput{
		Name: "Name",
		Type: "string",
	}

	widthPropType := models.PropertyTypeInput{
		Name: "Width",
		Type: "number",
	}

	_, _ = mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "example_type_a",
		Properties: []*models.PropertyTypeInput{&namePropType, &widthPropType},
	})

	propDefs, err := qr.PossibleProperties(ctx, models.PropertyEntityEquipment)
	require.NoError(t, err)
	for _, propDef := range propDefs {
		assert.True(t, propDef.Name == "Name" || propDef.Name == "Width")
	}

	assert.Len(t, propDefs, 2)
}

func TestSearchEquipmentByLocation(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	locType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type1",
	})

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst1",
		Type: locType.ID,
	})
	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "loc_inst2",
		Type:   locType.ID,
		Parent: &loc1.ID,
	})
	eqType, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type",
	})

	_, _ = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst1",
		Type:     eqType.ID,
		Location: &loc1.ID,
	})
	_, _ = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst2",
		Type:     eqType.ID,
		Location: &loc2.ID,
	})

	maxDepth := 2
	limit := 100
	f1 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{loc1.ID},
		MaxDepth:   &maxDepth,
	}
	res1, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Equipment, 2)

	f2 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{loc2.ID},
		MaxDepth:   &maxDepth,
	}
	res2, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Equipment, 1)
}

func TestSearchEquipmentByDate(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr := r.Mutation(), r.Query()
	locType, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type1",
	})

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst1",
		Type: locType.ID,
	})
	date := "2020-01-01"
	propType := models.PropertyTypeInput{
		Name:        "install_date",
		Type:        models.PropertyKindDate,
		StringValue: &date,
	}
	eqType, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "eq_type",
		Properties: []*models.PropertyTypeInput{&propType},
	})
	date = "2010-01-01"
	ptypeID := eqType.QueryPropertyTypes().OnlyXID(ctx)

	prop1 := models.PropertyInput{
		PropertyTypeID: ptypeID,
		StringValue:    &date,
	}
	e1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:       "eq_inst1",
		Type:       eqType.ID,
		Location:   &loc1.ID,
		Properties: []*models.PropertyInput{&prop1},
	})
	e2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst2",
		Type:     eqType.ID,
		Location: &loc1.ID,
	})
	_ = e2
	limit := 100
	date = "2015-05-05"
	f1 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeProperty,
		Operator:   models.FilterOperatorDateGreaterThan,
		PropertyValue: &models.PropertyTypeInput{
			Name:        "install_date",
			Type:        models.PropertyKindDate,
			StringValue: &date,
		},
	}

	res1, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res1.Equipment, 1)
	require.Equal(t, res1.Equipment[0].ID, e2.ID)

	f2 := models.EquipmentFilterInput{
		FilterType: models.EquipmentFilterTypeProperty,
		Operator:   models.FilterOperatorDateLessThan,
		PropertyValue: &models.PropertyTypeInput{
			Name:        "install_date",
			Type:        models.PropertyKindDate,
			StringValue: &date,
		},
	}
	res2, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res2.Equipment, 1)
	require.Equal(t, res2.Equipment[0].ID, e1.ID)

	res3, err := qr.EquipmentSearch(ctx, []*models.EquipmentFilterInput{&f1, &f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res3.Equipment, 0)
}

func TestSearchWO(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := ent.NewContext(viewertest.NewContext(context.Background(), r.client), r.client)

	c := r.GraphClient()
	data := prepareWOData(ctx, r, "A")
	/*
		helper: data now is of type:
		wo type a :
			Awo_1: loc1, assignee1. has "owner" and install date
			Awo_2: no loc, assignee1
		wo type b :
			Awo_3: loc1, assignee2
			Awo_4: loc2, no assignee
	*/

	var result woSearchResult
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{}),
	)
	require.Equal(t, 4, result.WorkOrderSearch.Count)
	require.Empty(t, result.WorkOrderSearch.WorkOrders)

	name := "_1"
	f1 := models.WorkOrderFilterInput{
		FilterType:  models.WorkOrderFilterTypeWorkOrderName,
		Operator:    models.FilterOperatorContains,
		StringValue: &name,
	}
	c.MustPost(
		woAllQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f1}),
	)
	require.Equal(t, 1, result.WorkOrderSearch.Count)
	require.Equal(t, strconv.Itoa(data.wo1), result.WorkOrderSearch.WorkOrders[0].ID)

	status := models.WorkOrderStatusPlanned.String()
	f2 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderStatus,
		Operator:   models.FilterOperatorIsOneOf,
		StringSet:  []string{status},
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f2}),
	)
	require.Equal(t, 3, result.WorkOrderSearch.Count)

	f3 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.woType1},
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f3}),
	)
	require.Equal(t, 2, result.WorkOrderSearch.Count)

	f4 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderAssignedTo,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.assignee1},
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f4}),
	)
	require.Equal(t, 2, result.WorkOrderSearch.Count)

	f5 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.loc1},
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f5}),
	)
	require.Equal(t, 2, result.WorkOrderSearch.Count)

	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f4, f5}),
	)
	require.Equal(t, 1, result.WorkOrderSearch.Count)

	f7 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderOwnedBy,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.owner},
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f7}),
	)
	require.Equal(t, 1, result.WorkOrderSearch.Count)

	f8 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderInstallDate,
		Operator:   models.FilterOperatorIs,
		StringValue: pointer.ToString(
			strconv.FormatInt(data.installDate.Unix(), 10),
		),
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f8}),
	)
	require.Equal(t, 1, result.WorkOrderSearch.Count)

	f9 := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderCreationDate,
		Operator:   models.FilterOperatorIs,
		StringValue: pointer.ToString(
			strconv.FormatInt(time.Now().Unix(), 10),
		),
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f9}),
	)
	require.Equal(t, 4, result.WorkOrderSearch.Count)
}

func TestSearchWOByPriority(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := ent.NewContext(viewertest.NewContext(context.Background(), r.client), r.client)
	data := prepareWOData(ctx, r, "B")
	c := r.GraphClient()

	var result woSearchResult
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{}),
	)
	require.Equal(t, 4, result.WorkOrderSearch.Count)

	f := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderPriority,
		Operator:   models.FilterOperatorIsOneOf,
		StringSet:  []string{models.WorkOrderPriorityHigh.String()},
	}
	c.MustPost(
		woAllQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f}),
	)
	require.Equal(t, 1, result.WorkOrderSearch.Count)
	require.Equal(t, strconv.Itoa(data.wo1), result.WorkOrderSearch.WorkOrders[0].ID)

	f.StringSet = []string{models.WorkOrderPriorityLow.String()}
	c.MustPost(
		woAllQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f}),
	)
	require.Zero(t, result.WorkOrderSearch.Count)
}

func TestSearchWOByLocation(t *testing.T) {
	r := newTestResolver(t)
	defer r.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	c := r.GraphClient()

	data := prepareWOData(ctx, r, "A")
	/*
		helper: data now is of type:
		wo type a :
			Awo_1: loc1, assignee1. has "owner" and install date
			Awo_2: no loc, assignee1
		wo type b :
			Awo_3: loc1, assignee2
			Awo_4: loc2, no assignee
	*/
	var result woSearchResult
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{}),
	)
	require.Equal(t, 4, result.WorkOrderSearch.Count)
	f := models.WorkOrderFilterInput{
		FilterType: models.WorkOrderFilterTypeWorkOrderLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.loc1},
		MaxDepth:   pointer.ToInt(2),
	}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f}),
	)
	require.Equal(t, 2, result.WorkOrderSearch.Count)

	f.IDSet = []int{-1}
	c.MustPost(
		woCountQuery,
		&result,
		client.Var("filters", []models.WorkOrderFilterInput{f}),
	)
	require.Zero(t, result.WorkOrderSearch.Count)
}
