// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNumOfProjects(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, ptr := r.Mutation(), r.ProjectType()

	pType, err := mr.CreateProjectType(ctx, models.AddProjectTypeInput{Name: "example_type"})
	require.NoError(t, err)

	numWO, err := ptr.NumberOfProjects(ctx, pType)
	require.NoError(t, err)
	require.Equal(t, 0, numWO)

	workOrder, err := mr.CreateProject(ctx, models.AddProjectInput{
		Name: "foo", Type: pType.ID,
	})
	require.NoError(t, err)

	numWO, err = ptr.NumberOfProjects(ctx, pType)
	require.NoError(t, err)
	require.Equal(t, 1, numWO)

	_, err = mr.DeleteProject(ctx, workOrder.ID)
	require.NoError(t, err)

	numWO, err = ptr.NumberOfProjects(ctx, pType)
	require.NoError(t, err)
	require.Equal(t, 0, numWO)
}

func TestProjectQuery(t *testing.T) {
	resolver, ctx := resolverctx(t)

	typ, err := resolver.Mutation().CreateProjectType(
		ctx, models.AddProjectTypeInput{Name: "test", Description: pointer.ToString("foobar")},
	)
	require.NoError(t, err)

	node, err := resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok := node.(*ent.ProjectType)
	require.True(t, ok)
	assert.Equal(t, typ.Name, rtyp.Name)
	assert.Equal(t, typ.Description, rtyp.Description)

	proj, err := resolver.Mutation().CreateProject(
		ctx, models.AddProjectInput{
			Name:        "test-project",
			Type:        typ.ID,
			Description: pointer.ToString("baz"),
		},
	)
	require.NoError(t, err)
	node, err = resolver.Query().Node(ctx, proj.ID)
	require.NoError(t, err)
	rproj, ok := node.(*ent.Project)
	require.True(t, ok)
	assert.Equal(t, proj.Name, rproj.Name)
	assert.Equal(t, proj.Description, rproj.Description)
}

func TestProjectWithWorkOrders(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()
	ctx := viewertest.NewContext(context.Background(), resolver.client)
	mutation := resolver.Mutation()

	woType, err := mutation.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type_a"})
	require.NoError(t, err)
	woDef := models.WorkOrderDefinitionInput{Type: woType.ID, Index: pointer.ToInt(1)}

	typ, err := resolver.Mutation().CreateProjectType(
		ctx, models.AddProjectTypeInput{
			Name:        "test",
			Description: pointer.ToString("foobar"),
			WorkOrders:  []*models.WorkOrderDefinitionInput{&woDef},
		},
	)
	require.NoError(t, err)
	node, err := resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok := node.(*ent.ProjectType)
	require.True(t, ok)
	woDefs, err := rtyp.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(woDefs))

	location := createLocation(ctx, t, *resolver)
	input := models.AddProjectInput{Name: "test", Type: typ.ID, Location: &location.ID}
	proj, err := mutation.CreateProject(ctx, input)
	require.NoError(t, err)
	wos, err := proj.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Len(t, wos, 1)
	wo := wos[0]
	assert.EqualValues(t, wo.Name, woType.Name)
	assert.EqualValues(t, wo.Index, *woDef.Index)
	assert.EqualValues(t, wo.QueryLocation().FirstXID(ctx), location.ID)
}

func TestEditProjectTypeWorkOrders(t *testing.T) {
	resolver, ctx := resolverctx(t)
	mutation := resolver.Mutation()

	woType, err := mutation.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{Name: "example_type_a"})
	require.NoError(t, err)
	woDef := models.WorkOrderDefinitionInput{Type: woType.ID, Index: pointer.ToInt(1)}

	typ, err := resolver.Mutation().CreateProjectType(
		ctx, models.AddProjectTypeInput{
			Name:        "test",
			Description: pointer.ToString("foobar"),
			WorkOrders:  []*models.WorkOrderDefinitionInput{&woDef},
		},
	)
	require.NoError(t, err)
	node, err := resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok := node.(*ent.ProjectType)
	require.True(t, ok)
	woDefs, err := rtyp.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(woDefs))

	woDef = models.WorkOrderDefinitionInput{ID: &woDefs[0].ID, Type: woType.ID, Index: pointer.ToInt(2)}
	typ, err = resolver.Mutation().EditProjectType(
		ctx, models.EditProjectTypeInput{
			ID:          typ.ID,
			Name:        "test",
			Description: pointer.ToString("foobar"),
			WorkOrders:  []*models.WorkOrderDefinitionInput{&woDef},
		},
	)
	require.NoError(t, err)
	node, err = resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok = node.(*ent.ProjectType)
	require.True(t, ok)
	woDefs, err = rtyp.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(woDefs))
	assert.Equal(t, *woDef.ID, woDefs[0].ID)

	woDef2 := models.WorkOrderDefinitionInput{Type: woType.ID, Index: pointer.ToInt(3)}
	typ, err = resolver.Mutation().EditProjectType(
		ctx, models.EditProjectTypeInput{
			ID:          typ.ID,
			Name:        "test",
			Description: pointer.ToString("foobar"),
			WorkOrders:  []*models.WorkOrderDefinitionInput{&woDef2},
		},
	)
	require.NoError(t, err)
	node, err = resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok = node.(*ent.ProjectType)
	require.True(t, ok)
	woDefs, err = rtyp.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(woDefs))
	assert.NotEqual(t, *woDef.ID, woDefs[0].ID)
}

func TestProjectMutation(t *testing.T) {
	mutation, ctx := mutationctx(t)
	input := models.AddProjectTypeInput{Name: "test", Description: pointer.ToString("test desc")}
	ltyp, err := mutation.AddLocationType(ctx, models.AddLocationTypeInput{Name: "loc_type"})
	require.NoError(t, err)
	loc, err := mutation.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_name",
		Type: ltyp.ID,
	})
	require.NoError(t, err)
	typ, err := mutation.CreateProjectType(ctx, input)
	require.NoError(t, err)
	assert.Equal(t, input.Name, typ.Name)
	assert.EqualValues(t, input.Description, typ.Description)
	_, err = mutation.CreateProjectType(ctx, models.AddProjectTypeInput{})
	assert.Error(t, err, "project type name cannot be empty")
	_, err = mutation.CreateProjectType(ctx, input)
	assert.Error(t, err, "project type name must be unique")

	var project *ent.Project
	{
		input := models.AddProjectInput{
			Name:        "test",
			Description: pointer.ToString("desc"),
			Type:        typ.ID,
			Location:    &loc.ID,
		}
		project, err = mutation.CreateProject(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, input.Name, project.Name)
		assert.Equal(t, *input.Location, project.QueryLocation().OnlyX(ctx).ID)

		_, err = mutation.CreateProject(ctx, input)
		assert.Error(t, err, "project name must be unique under type")
		_, err = mutation.CreateProject(ctx, models.AddProjectInput{Type: input.Type})
		assert.Error(t, err, "project name cannot be empty")
		_, err = mutation.CreateProject(ctx, models.AddProjectInput{Name: "another", Type: 42424242})
		assert.Error(t, err, "project type id must be valid")
	}

	deleted, err := mutation.DeleteProjectType(ctx, typ.ID)
	assert.Error(t, err, "project type cannot be deleted with associated projects")
	assert.False(t, deleted)
	deleted, err = mutation.DeleteProject(ctx, project.ID)
	assert.NoError(t, err)
	assert.True(t, deleted)
	deleted, err = mutation.DeleteProject(ctx, project.ID)
	assert.EqualError(t, err, errNoProject.Error(), "project cannot be deleted twice")
	assert.False(t, deleted)

	deleted, err = mutation.DeleteProjectType(ctx, typ.ID)
	assert.NoError(t, err)
	assert.True(t, deleted)
	deleted, err = mutation.DeleteProjectType(ctx, typ.ID)
	assert.EqualError(t, err, errNoProjectType.Error(), "project type cannot be deleted twice")
	assert.False(t, deleted)
}

func TestEditProject(t *testing.T) {
	mutation, ctx := mutationctx(t)
	input := models.AddProjectTypeInput{Name: "test", Description: pointer.ToString("test desc")}
	ltyp, err := mutation.AddLocationType(ctx, models.AddLocationTypeInput{Name: "loc_type"})
	require.NoError(t, err)
	loc, err := mutation.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_name",
		Type: ltyp.ID,
	})
	require.NoError(t, err)
	typ, err := mutation.CreateProjectType(ctx, input)
	require.NoError(t, err)

	var project *ent.Project
	{
		u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
		input := models.AddProjectInput{
			Name:        "test",
			Description: pointer.ToString("desc"),
			Type:        typ.ID,
			Location:    &loc.ID,
			CreatorID:   &u.ID,
		}
		project, err = mutation.CreateProject(ctx, input)
		require.NoError(t, err)
		assert.Equal(t, input.Name, project.Name)
		assert.Equal(t, *input.Location, project.QueryLocation().OnlyX(ctx).ID)
		assert.Equal(t, *input.CreatorID, project.QueryCreator().OnlyXID(ctx))

		updateInput := models.EditProjectInput{
			ID:          project.ID,
			Name:        "new-test",
			Description: pointer.ToString("new-desc"),
			Type:        typ.ID,
		}
		project, err = mutation.EditProject(ctx, updateInput)
		require.NoError(t, err)
		assert.Equal(t, updateInput.Name, project.Name)
		assert.Equal(t, *updateInput.Description, *project.Description)
		assert.False(t, project.QueryCreator().ExistX(ctx))
	}
}

func TestEditProjectLocation(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()
	location := createLocation(ctx, t, *r)
	typ, err := mr.CreateProjectType(ctx, models.AddProjectTypeInput{Name: "example_type"})
	require.NoError(t, err)
	input := models.AddProjectInput{Name: "test", Type: typ.ID, Location: &location.ID}
	proj, err := mr.CreateProject(ctx, input)

	require.NoError(t, err)
	require.Equal(t, proj.QueryLocation().FirstXID(ctx), location.ID)

	location = createLocationWithName(ctx, t, *r, "location2")
	ei := models.EditProjectInput{ID: proj.ID, Name: "test", Type: typ.ID, Location: &location.ID}
	proj, err = mr.EditProject(ctx, ei)
	require.NoError(t, err)
	require.Equal(t, proj.QueryLocation().FirstXID(ctx), location.ID)

	ei = models.EditProjectInput{ID: proj.ID, Name: "test", Type: typ.ID}
	proj, err = mr.EditProject(ctx, ei)
	require.NoError(t, err)
	locEx, err := proj.QueryLocation().Exist(ctx)
	require.NoError(t, err)
	require.False(t, locEx)
}

func TestAddProjectWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	mutation, ctx := mutationctx(t)

	mr, qr, pr := r.Mutation(), r.Query(), r.Project()
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
	typ, err := mr.CreateProjectType(ctx, models.AddProjectTypeInput{Name: "example_type", Properties: propTypeInputs})
	require.NoError(t, err, "Adding project type")

	strValue := "Foo"
	strProp := models.PropertyInput{
		PropertyTypeID: typ.QueryProperties().Where(propertytype.Name("str_prop")).OnlyXID(ctx),
		StringValue:    &strValue,
	}
	strFixedProp := models.PropertyInput{
		PropertyTypeID: typ.QueryProperties().Where(propertytype.Name("str_fixed_prop")).OnlyXID(ctx),
		StringValue:    &strFixedValue,
	}
	intValue := 5
	intProp := models.PropertyInput{
		PropertyTypeID: typ.QueryProperties().Where(propertytype.Name("int_prop")).OnlyXID(ctx),
		StringValue:    nil,
		IntValue:       &intValue,
	}
	fl1, fl2 := 5.5, 7.8
	rngProp := models.PropertyInput{
		PropertyTypeID: typ.QueryProperties().Where(propertytype.Name("rng_prop")).OnlyXID(ctx),
		RangeFromValue: &fl1,
		RangeToValue:   &fl2,
	}
	propInputs := []*models.PropertyInput{&strProp, &strFixedProp, &intProp, &rngProp}
	u := viewer.FromContext(ctx).(*viewer.UserViewer).User()
	input := models.AddProjectInput{
		Name:        "test",
		Description: pointer.ToString("desc"),
		Type:        typ.ID,
		CreatorID:   &u.ID,
		Properties:  propInputs,
	}
	p, err := mutation.CreateProject(ctx, input)
	require.NoError(t, err, "adding project instance")

	node, err := qr.Node(ctx, p.ID)
	require.NoError(t, err, "querying project node")
	fetchedProj, ok := node.(*ent.Project)
	require.True(t, ok, "casting project instance")

	intFetchProp := fetchedProj.QueryProperties().Where(property.HasTypeWith(propertytype.Name("int_prop"))).OnlyX(ctx)
	require.Equal(t, intFetchProp.IntVal, *intProp.IntValue, "Comparing properties: int value")
	require.Equal(t, intFetchProp.QueryType().OnlyXID(ctx), intProp.PropertyTypeID, "Comparing properties: PropertyType value")

	strFetchProp := fetchedProj.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_prop"))).OnlyX(ctx)
	require.Equal(t, strFetchProp.StringVal, *strProp.StringValue, "Comparing properties: string value")
	require.Equal(t, strFetchProp.QueryType().OnlyXID(ctx), strProp.PropertyTypeID, "Comparing properties: PropertyType value")

	fixedStrFetchProp := fetchedProj.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_fixed_prop"))).OnlyX(ctx)
	require.Equal(t, fixedStrFetchProp.StringVal, *strFixedProp.StringValue, "Comparing properties: fixed string value")
	require.Equal(t, fixedStrFetchProp.QueryType().OnlyXID(ctx), strFixedProp.PropertyTypeID, "Comparing properties: PropertyType value")

	rngFetchProp := fetchedProj.QueryProperties().Where(property.HasTypeWith(propertytype.Name("rng_prop"))).OnlyX(ctx)
	require.Equal(t, rngFetchProp.RangeFromVal, *rngProp.RangeFromValue, "Comparing properties: range value")
	require.Equal(t, rngFetchProp.RangeToVal, *rngProp.RangeToValue, "Comparing properties: range value")
	require.Equal(t, rngFetchProp.QueryType().OnlyXID(ctx), rngProp.PropertyTypeID, "Comparing properties: PropertyType value")

	fetchedProps, err := pr.Properties(ctx, fetchedProj)
	require.NoError(t, err)
	require.Equal(t, len(propInputs), len(fetchedProps))

	failProp := models.PropertyInput{PropertyTypeID: -1}
	failEditInput := models.EditProjectInput{
		ID:         p.ID,
		Name:       "test",
		Properties: []*models.PropertyInput{&failProp},
	}
	_, err = mutation.EditProject(ctx, failEditInput)
	require.Error(t, err, "editing project instance property with wrong property type id")

	failProp2 := models.PropertyInput{
		ID:             &strFetchProp.ID,
		PropertyTypeID: intProp.PropertyTypeID,
	}
	failEditInput2 := models.EditProjectInput{
		ID:         p.ID,
		Name:       "test",
		Properties: []*models.PropertyInput{&failProp2},
	}
	_, err = mutation.EditProject(ctx, failEditInput2)
	require.Error(t, err, "editing project instance property when id and property type id mismach")

	newStrValue := "Foo"
	prop := models.PropertyInput{
		PropertyTypeID: strProp.PropertyTypeID,
		StringValue:    &newStrValue,
	}
	newProjectName := "updated test"
	editInput := models.EditProjectInput{
		ID:         p.ID,
		Name:       newProjectName,
		Properties: []*models.PropertyInput{&prop},
	}
	updatedP, err := mutation.EditProject(ctx, editInput)
	require.NoError(t, err)

	updatedNode, err := qr.Node(ctx, updatedP.ID)
	require.NoError(t, err, "querying updated project node")
	updatedProj, ok := updatedNode.(*ent.Project)
	require.True(t, ok, "casting updated project instance")

	require.Equal(t, updatedProj.Name, newProjectName, "Comparing updated project name")

	fetchedProps, _ = pr.Properties(ctx, updatedProj)
	require.Equal(t, len(propInputs), len(fetchedProps), "number of properties should remain he same")

	updatedProp := updatedProj.QueryProperties().Where(property.HasTypeWith(propertytype.Name("str_prop"))).OnlyX(ctx)
	require.Equal(t, updatedProp.StringVal, *prop.StringValue, "Comparing updated properties: string value")
	require.Equal(t, updatedProp.QueryType().OnlyXID(ctx), prop.PropertyTypeID, "Comparing updated properties: PropertyType value")
}

func TestEditProjectType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()

	pType, err := mr.CreateProjectType(ctx, models.AddProjectTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	newType, err := mr.EditProjectType(ctx, models.EditProjectTypeInput{
		ID:          pType.ID,
		Name:        "example_type_name_edited",
		Description: pointer.ToString("example_type_desc_edited"),
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited project name")
	require.Equal(t, "example_type_desc_edited", *newType.Description, "successfully edited project description")
	pType2, err := mr.CreateProjectType(ctx, models.AddProjectTypeInput{Name: "example_type_name_2"})
	require.NoError(t, err)
	_, err = mr.EditProjectType(ctx, models.EditProjectTypeInput{
		ID:   pType2.ID,
		Name: "example_type_name_edited",
	})
	require.Error(t, err, "duplicate names")

	types, err := qr.ProjectTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 2)

	node, err := qr.Node(ctx, pType.ID)
	require.NoError(t, err)
	typ, ok := node.(*ent.ProjectType)
	require.True(t, ok)
	assert.Equal(t, "example_type_name_edited", typ.Name)
}

func TestProjectWithWorkOrdersAndProperties(t *testing.T) {
	resolver := newTestResolver(t)
	defer resolver.drv.Close()
	ctx := viewertest.NewContext(context.Background(), resolver.client)
	mutation := resolver.Mutation()

	strPropType := models.PropertyTypeInput{
		Name: "str_prop",
		Type: "string",
	}
	intPropType := models.PropertyTypeInput{
		Name:        "int_prop",
		Type:        "int",
		IsMandatory: pointer.ToBool(true),
	}
	woType, err := mutation.AddWorkOrderType(ctx, models.AddWorkOrderTypeInput{
		Name:       "example_type_a",
		Properties: []*models.PropertyTypeInput{&strPropType, &intPropType},
	})
	require.NoError(t, err)
	woDef := models.WorkOrderDefinitionInput{Type: woType.ID, Index: pointer.ToInt(1)}

	typ, err := resolver.Mutation().CreateProjectType(
		ctx, models.AddProjectTypeInput{
			Name:        "test",
			Description: pointer.ToString("foobar"),
			WorkOrders:  []*models.WorkOrderDefinitionInput{&woDef},
		},
	)
	require.NoError(t, err)
	node, err := resolver.Query().Node(ctx, typ.ID)
	require.NoError(t, err)
	rtyp, ok := node.(*ent.ProjectType)
	require.True(t, ok)
	woDefs, err := rtyp.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Equal(t, 1, len(woDefs))

	location := createLocation(ctx, t, *resolver)
	input := models.AddProjectInput{Name: "test", Type: typ.ID, Location: &location.ID}
	proj, err := mutation.CreateProject(ctx, input)
	require.NoError(t, err)
	wos, err := proj.QueryWorkOrders().All(ctx)
	require.NoError(t, err)
	assert.Len(t, wos, 1)
	props := wos[0].QueryProperties().AllX(ctx)
	require.Len(t, props, 2)
}
