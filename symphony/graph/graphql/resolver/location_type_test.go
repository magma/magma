// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// nolint: goconst
package resolver

import (
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/locationtype"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestAddLocationType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	mapType := "map"
	mapZoomLvl := 12

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_type",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
	})
	require.NoError(t, err)

	fetchedLocType, _ := qr.LocationType(ctx, locType.ID)
	require.Equal(t, fetchedLocType.Name, "example_type", "verifying location type name")
	require.Equal(t, fetchedLocType.MapType, mapType, "verifying location type map type")
	require.Equal(t, fetchedLocType.MapZoomLevel, mapZoomLvl, "verifying location type zoom level")
}

func TestAddLocationTypes(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	_, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "a"})
	require.NoError(t, err)
	_, err = mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "b"})
	require.NoError(t, err)

	types, _ := qr.LocationTypes(ctx, nil, nil, nil, nil)
	require.Len(t, types.Edges, 2, "verify the added location types are fetched properly")
}

func TestAddLocationTypeWithProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	strValue, strIndex := "Foo", 7
	intValue, intIndex := 5, 12

	strPropType := models.PropertyTypeInput{
		Name:        "str_prop",
		Type:        "string",
		Index:       &strIndex,
		StringValue: &strValue,
	}
	intPropType := models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     "int",
		Index:    &intIndex,
		IntValue: &intValue,
	}
	propTypeInputs := []*models.PropertyTypeInput{&strPropType, &intPropType}
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:       "example_type_a",
		Properties: propTypeInputs,
	})
	require.NoError(t, err)

	intProp := locType.QueryPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	strProp := locType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)

	require.Equal(t, "int_prop", intProp.Name, "verifying int property type's name")
	require.Equal(t, "", intProp.StringVal, "verifying int property type's string value (default as this is an int property)")
	require.Equal(t, intValue, intProp.IntVal, "verifying int property type's int value")
	require.Equal(t, intIndex, intProp.Index, "verifying int property type's index")
	require.Equal(t, "str_prop", strProp.Name, "verifying string property type's name")
	require.Equal(t, strValue, strProp.StringVal, "verifying string property type's String value")
	require.Equal(t, 0, strProp.IntVal, "verifying int property type's int value")
	require.Equal(t, strIndex, strProp.Index, "verifying string property type's index")
}

func TestAddLocationTypeWithEquipmentProperty(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	lt, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "location_type"})
	require.NoError(t, err)
	l, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "location_name",
		Type: lt.ID,
	})
	require.NoError(t, err)

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "equipment_type",
	})
	require.NoError(t, err)

	_, err = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "equipment_name",
		Type:     equipmentType.ID,
		Location: &l.ID,
	})
	require.NoError(t, err)

	index := 0
	eqPropType := models.PropertyTypeInput{
		Name:  "eq_prop",
		Type:  "equipment",
		Index: &index,
	}
	propTypeInputs := []*models.PropertyTypeInput{&eqPropType}
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:       "example_type",
		Properties: propTypeInputs,
	})
	require.NoError(t, err)

	eqProp := locType.QueryPropertyTypes().Where(propertytype.Type("equipment")).OnlyX(ctx)

	require.Equal(t, "eq_prop", eqProp.Name)
	require.Equal(t, "equipment", eqProp.Type)
}

func TestAddLocationTypeWithSurveyTemplate(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	question := models.SurveyTemplateQuestionInput{
		QuestionTitle:       "What is the power rating?",
		QuestionDescription: "Tell me more about this question",
		QuestionType:        models.SurveyQuestionTypeText,
		Index:               0,
	}

	category := models.SurveyTemplateCategoryInput{
		CategoryTitle:           "Power",
		CategoryDescription:     "Description",
		SurveyTemplateQuestions: []*models.SurveyTemplateQuestionInput{&question},
	}

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:                     "example_type_a",
		SurveyTemplateCategories: []*models.SurveyTemplateCategoryInput{&category},
	})
	require.NoError(t, err)

	categories, _ := locType.QuerySurveyTemplateCategories().All(ctx)
	require.Equal(t, categories[0].CategoryTitle, category.CategoryTitle)
	require.Equal(t, categories[0].CategoryDescription, category.CategoryDescription)

	questions, _ := categories[0].QuerySurveyTemplateQuestions().All(ctx)
	require.Equal(t, questions[0].QuestionTitle, question.QuestionTitle)
	require.Equal(t, questions[0].QuestionDescription, question.QuestionDescription)
}

func TestAddLocationTypesSameName(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", locType.Name)
	_, err = mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_name"})
	require.Error(t, err, "adding location type with an existing location type name yields an error")
	types, _ := qr.LocationTypes(ctx, nil, nil, nil, nil)
	require.Len(t, types.Edges, 1, "two attempts to create location types with same name will create one location type")
}

func TestRemoveLocationType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", locType.Name)

	types, err := qr.LocationTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 1)

	id, err := mr.RemoveLocationType(ctx, locType.ID)
	require.NoError(t, err)
	require.Equal(t, locType.ID, id, "successfully remove location type")
	types, err = qr.LocationTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Empty(t, types.Edges, "no location types exist after deletion")
}

func TestEditLocationType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_name"})
	require.NoError(t, err)
	require.Equal(t, "example_type_name", locType.Name)

	newType, err := mr.EditLocationType(ctx, models.EditLocationTypeInput{
		ID: locType.ID, Name: "example_type_name_edited",
	})
	require.NoError(t, err)
	require.Equal(t, "example_type_name_edited", newType.Name, "successfully edited location type name")

	locType, err = mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "example_type_name_2"})
	require.NoError(t, err)
	_, err = mr.EditLocationType(ctx, models.EditLocationTypeInput{
		ID: locType.ID, Name: "example_type_name_edited",
	})
	require.Error(t, err, "duplicate names")

	types, err := qr.LocationTypes(ctx, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, types.Edges, 2)

	typ, err := qr.LocationType(ctx, locType.ID)
	require.NoError(t, err)
	require.Equal(t, "example_type_name_2", typ.Name)
}

func TestEditLocationTypeWithSurveyTemplate(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	question := models.SurveyTemplateQuestionInput{
		QuestionTitle:       "What is the power rating?",
		QuestionDescription: "Tell me more about this question",
		QuestionType:        models.SurveyQuestionTypeText,
	}

	category := models.SurveyTemplateCategoryInput{
		CategoryTitle:           "Power",
		CategoryDescription:     "Description",
		SurveyTemplateQuestions: []*models.SurveyTemplateQuestionInput{&question},
	}

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:                     "example_type_name",
		SurveyTemplateCategories: []*models.SurveyTemplateCategoryInput{&category},
	})

	require.NoError(t, err)
	require.Equal(t, "example_type_name", locType.Name)

	categories, _ := locType.QuerySurveyTemplateCategories().All(ctx)
	questions, _ := categories[0].QuerySurveyTemplateQuestions().All(ctx)

	updatedQuestion := models.SurveyTemplateQuestionInput{
		ID:                  &questions[0].ID,
		QuestionTitle:       "New Title",
		QuestionDescription: "New Description",
		QuestionType:        models.SurveyQuestionTypeText,
	}

	updatedCategory := models.SurveyTemplateCategoryInput{
		ID:                      &categories[0].ID,
		CategoryTitle:           "New Power",
		CategoryDescription:     "Updated Description",
		SurveyTemplateQuestions: []*models.SurveyTemplateQuestionInput{&updatedQuestion},
	}

	_, err = mr.EditLocationTypeSurveyTemplateCategories(ctx, locType.ID, []*models.SurveyTemplateCategoryInput{&updatedCategory})
	require.NoError(t, err)

	categories, _ = locType.QuerySurveyTemplateCategories().All(ctx)
	require.Equal(t, len(categories), 1)
	require.Equal(t, categories[0].CategoryTitle, updatedCategory.CategoryTitle)
	require.Equal(t, categories[0].CategoryDescription, updatedCategory.CategoryDescription)

	questions, _ = categories[0].QuerySurveyTemplateQuestions().All(ctx)
	require.Equal(t, len(questions), 1)
	require.Equal(t, questions[0].QuestionTitle, updatedQuestion.QuestionTitle)
	require.Equal(t, questions[0].QuestionDescription, updatedQuestion.QuestionDescription)

	updatedCategory = models.SurveyTemplateCategoryInput{
		ID:                      &categories[0].ID,
		CategoryTitle:           "New Power",
		CategoryDescription:     "Updated Description",
		SurveyTemplateQuestions: []*models.SurveyTemplateQuestionInput{},
	}

	categories, err = mr.EditLocationTypeSurveyTemplateCategories(ctx, locType.ID, []*models.SurveyTemplateCategoryInput{&updatedCategory})
	require.NoError(t, err)

	questions, err = categories[0].QuerySurveyTemplateQuestions().All(ctx)
	require.NoError(t, err)
	require.Equal(t, len(categories), 1)
	require.Equal(t, len(questions), 0)

	_, err = mr.EditLocationTypeSurveyTemplateCategories(ctx, locType.ID, []*models.SurveyTemplateCategoryInput{})
	require.NoError(t, err)
	categories, _ = locType.QuerySurveyTemplateCategories().All(ctx)
	require.Equal(t, len(categories), 0)
}

func TestEditLocationTypeWithProperties(t *testing.T) {
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
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:       "example_type_a",
		Properties: propTypeInput,
	})
	require.NoError(t, err)

	strProp := locType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
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
	newType, err := mr.EditLocationType(ctx, models.EditLocationTypeInput{
		ID: locType.ID, Name: "example_type_a", Properties: editedPropTypeInput,
	})
	require.NoError(t, err)
	require.Equal(t, locType.Name, newType.Name, "successfully edited location type name")

	strProp = locType.QueryPropertyTypes().Where(propertytype.Type("string")).OnlyX(ctx)
	require.Equal(t, "str_prop_new", strProp.Name, "successfully edited prop type name")
	require.Equal(t, strValue, strProp.StringVal, "successfully edited prop type string value")

	intProp := locType.QueryPropertyTypes().Where(propertytype.Type("int")).OnlyX(ctx)
	require.Equal(t, "int_prop", intProp.Name, "successfully edited prop type name")
	require.Equal(t, intValue, intProp.IntVal, "successfully edited prop type int value")

	intValue = 6
	intPropType = models.PropertyTypeInput{
		Name:     "int_prop",
		Type:     "int",
		IntValue: &intValue,
	}
	editedPropTypeInput = []*models.PropertyTypeInput{&intPropType}
	_, err = mr.EditLocationType(ctx, models.EditLocationTypeInput{
		ID: locType.ID, Name: "example_type_a", Properties: editedPropTypeInput,
	})
	require.Error(t, err, "duplicate property type names")
}

func TestMarkLocationTypeAsSite(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr, qr := r.Mutation(), r.Query()

	mapType := "map"
	mapZoomLvl := 12
	isSite := true

	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_type",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
		IsSite:       &isSite,
	})
	require.NoError(t, err)

	fetchedLocType, _ := qr.LocationType(ctx, locType.ID)
	require.True(t, fetchedLocType.Site)
}

func TestEditLocationTypesIndex(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)
	mr := r.Mutation()

	mapType := "map"
	mapZoomLvl := 12
	locType1, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_type1",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
	})
	require.NoError(t, err)
	mapZoomLvl++
	locType2, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_type2",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
	})
	require.NoError(t, err)
	mapZoomLvl++
	locType3, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name:         "example_type3",
		MapType:      &mapType,
		MapZoomLevel: &mapZoomLvl,
	})
	require.NoError(t, err)

	_, err = mr.EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: locType1.ID,
			Index:          2,
		},
		{
			LocationTypeID: locType2.ID,
			Index:          0,
		},
		{
			LocationTypeID: locType3.ID,
			Index:          1,
		},
	})
	require.NoError(t, err)
	require.Equal(t, 2, r.client.LocationType.Query().Where(locationtype.ID(locType1.ID)).OnlyX(ctx).Index)
	require.Equal(t, 0, r.client.LocationType.Query().Where(locationtype.ID(locType2.ID)).OnlyX(ctx).Index)
	require.Equal(t, 1, r.client.LocationType.Query().Where(locationtype.ID(locType3.ID)).OnlyX(ctx).Index)
}
