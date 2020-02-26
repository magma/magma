// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

var (
	locationFixedDataHeader = []string{"External ID", "Latitude", "Longitude"}
	locationIDHeader        = []string{"Location ID"}
)

func prepareLocationTypesWithProperties(ctx context.Context, t *testing.T, r TestImporterResolver) locTypeIDs {
	mr := r.importer.r.Mutation()
	strDefVal := propDefValue
	propDefInput1 := models.PropertyTypeInput{
		Name:        propName1,
		Type:        "string",
		StringValue: &strDefVal,
	}
	propDefInput2 := models.PropertyTypeInput{
		Name: propName2,
		Type: "int",
	}
	propDefInput3 := models.PropertyTypeInput{
		Name: propName3,
		Type: "date",
	}
	propDefInput4 := models.PropertyTypeInput{
		Name: propName4,
		Type: "bool",
	}
	propDefInput5 := models.PropertyTypeInput{
		Name: propName5,
		Type: "range",
	}
	propDefInput6 := models.PropertyTypeInput{
		Name: propName6,
		Type: "gps_location",
	}

	l, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameL, Properties: []*models.PropertyTypeInput{&propDefInput5, &propDefInput6}})
	require.NoError(t, err)
	m, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameM, Properties: []*models.PropertyTypeInput{&propDefInput3, &propDefInput4}})
	require.NoError(t, err)
	s, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameS, Properties: []*models.PropertyTypeInput{&propDefInput1, &propDefInput2}})
	require.NoError(t, err)

	_, err = r.importer.r.Mutation().EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: l.ID,
			Index:          0,
		},
		{
			LocationTypeID: m.ID,
			Index:          1,
		},
		{
			LocationTypeID: s.ID,
			Index:          2,
		},
	})
	require.NoError(t, err)
	return locTypeIDs{
		locTypeIDS: s.ID,
		locTypeIDM: m.ID,
		locTypeIDL: l.ID,
	}
}

func TestLocationTitleInputValidation(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()

	ctx := newImportContext(viewertest.NewContext(r.client))
	prepareBasicData(ctx, t, *r)

	header, _ := NewImportHeader([]string{"aa"}, ImportEntityLocation)
	err := importer.inputValidationsLocation(ctx, header)
	require.Error(t, err)

	header, _ = NewImportHeader(locationIDHeader, ImportEntityLocation)
	err = importer.inputValidationsLocation(ctx, header)
	require.Error(t, err)

	locationTypeNotInOrder := append(append(locationIDHeader, []string{locTypeNameS, locTypeNameM, locTypeNameL}...), locationFixedDataHeader...)

	header, _ = NewImportHeader(locationTypeNotInOrder, ImportEntityLocation)
	err = importer.inputValidationsLocation(ctx, header)
	require.Error(t, err)

	locationTypeInOrder := append(append(locationIDHeader, []string{locTypeNameL, locTypeNameM, locTypeNameS}...), locationFixedDataHeader...)

	header, _ = NewImportHeader(locationTypeInOrder, ImportEntityLocation)
	err = importer.inputValidationsLocation(ctx, header)
	require.NoError(t, err)
}

func TestImportLocationHierarchy(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))

	ids := prepareBasicData(ctx, t, *r)

	var (
		test1 = []string{"", "locNameL", "", "", "external_1", "32", "33"}
		test2 = []string{"", "locNameL", "locNameM", "locNameS", "external_2", "32", "33"}
		test3 = []string{"", "", "locNameM", "", "external_3", "32", "33"}
		test4 = []string{"", "locNameL", "", "locNameS", "external_3", "32", "33"}
	)
	locationTypeInOrder := append(append(locationIDHeader, []string{locTypeNameL, locTypeNameM, locTypeNameS}...), locationFixedDataHeader...)
	title, _ := NewImportHeader(locationTypeInOrder, ImportEntityLocation)
	err := importer.inputValidationsLocation(ctx, title)
	require.NoError(t, err)

	rec1 := NewImportRecord(test1, title)
	parentIndex, err := importer.getParentOfLocationIndex(ctx, rec1)
	require.NoError(t, err)
	require.Equal(t, parentIndex, -1)
	currIndex, err := importer.getCurrentLocationIndex(ctx, rec1)
	require.NoError(t, err)
	require.Equal(t, currIndex, 1)

	rec2 := NewImportRecord(test2, title)
	parentIndex, err = importer.getParentOfLocationIndex(ctx, rec2)
	require.NoError(t, err)
	require.Equal(t, parentIndex, 2)
	currIndex, err = importer.getCurrentLocationIndex(ctx, rec2)
	require.NoError(t, err)
	require.Equal(t, currIndex, 3)

	parentLoc2, err := importer.verifyOrCreateLocationHierarchy(ctx, rec2, true, &parentIndex)
	require.NoError(t, err)
	require.Equal(t, parentLoc2.Name, "locNameM")
	require.Equal(t, parentLoc2.QueryType().OnlyXID(ctx), ids.locTypeIDM)
	require.Equal(t, parentLoc2.QueryParent().OnlyX(ctx).Name, "locNameL")

	rec3 := NewImportRecord(test3, title)
	parentIndex, err = importer.getParentOfLocationIndex(ctx, rec3)
	require.NoError(t, err)
	require.Equal(t, parentIndex, -1)
	currIndex, err = importer.getCurrentLocationIndex(ctx, rec3)
	require.NoError(t, err)
	require.Equal(t, currIndex, 2)

	rec4 := NewImportRecord(test4, title)
	parentIndex, err = importer.getParentOfLocationIndex(ctx, rec4)
	require.NoError(t, err)
	require.Equal(t, parentIndex, 1)
	currIndex, err = importer.getCurrentLocationIndex(ctx, rec4)
	require.NoError(t, err)
	require.Equal(t, currIndex, 3)

	parentLoc3, err := importer.verifyOrCreateLocationHierarchy(ctx, rec2, true, &parentIndex)
	require.NoError(t, err)
	require.Equal(t, parentLoc3.Name, "locNameL")
	require.Equal(t, parentLoc3.QueryType().OnlyXID(ctx), ids.locTypeIDL)
	require.False(t, parentLoc3.QueryParent().ExistX(ctx))
}

func TestValidateLocationPropertiesForType(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	q := r.importer.r.Query()
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	data := prepareLocationTypesWithProperties(ctx, t, *r)

	var (
		test1 = []string{"", "locNameL", "locNameM", "locNameS", "external_2", "32", "33", "strVal", "54", "", "", "", ""}
		test2 = []string{"", "", "locNameM", "", "external_3", "32", "33", "", "", "29/03/88", "false", "", ""}
		test3 = []string{"", "locNameL", "", "", "external_1", "32", "33", "", "", "", "", "30.23-50", "45.8,88.9"}
	)

	firstRowLocations := append(append(locationIDHeader, []string{locTypeNameL, locTypeNameM, locTypeNameS}...), locationFixedDataHeader...)
	finalFirstRow := append(firstRowLocations, propName1, propName2, propName3, propName4, propName5, propName6)
	fl, _ := NewImportHeader(finalFirstRow, ImportEntityLocation)

	err := importer.inputValidationsLocation(ctx, fl)
	require.NoError(t, err)

	fl, _ = NewImportHeader(finalFirstRow, ImportEntityLocation)
	r1 := NewImportRecord(test1, fl)
	require.NoError(t, err)
	lType1, err := q.LocationType(ctx, data.locTypeIDS)
	require.NoError(t, err)

	ptypes, err := importer.validatePropertiesForLocationType(ctx, r1, lType1)
	require.NoError(t, err)
	require.Len(t, ptypes, 2)
	require.NotEqual(t, ptypes[0].PropertyTypeID, ptypes[1].PropertyTypeID)
	for _, value := range ptypes {
		ptyp := lType1.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName1:
			require.Equal(t, *value.StringValue, "strVal")
			require.Equal(t, ptyp.Type, "string")
		case propName2:
			require.Equal(t, *value.IntValue, 54)
			require.Equal(t, ptyp.Type, "int")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}
	lType2, err := q.LocationType(ctx, data.locTypeIDM)
	require.NoError(t, err)

	r2 := NewImportRecord(test2, fl)
	ptypes2, err := importer.validatePropertiesForLocationType(ctx, r2, lType2)
	require.NoError(t, err)
	require.Len(t, ptypes2, 2)
	for _, value := range ptypes2 {
		ptyp := lType2.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName3:
			require.Equal(t, *value.StringValue, "29/03/88")
			require.Equal(t, ptyp.Type, "date")
		case propName4:
			require.Equal(t, *value.BooleanValue, false)
			require.Equal(t, ptyp.Type, "bool")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}

	lType3, err := q.LocationType(ctx, data.locTypeIDL)
	require.NoError(t, err)

	r3 := NewImportRecord(test3, fl)
	ptypes3, err := importer.validatePropertiesForLocationType(ctx, r3, lType3)
	require.NoError(t, err)
	require.Len(t, ptypes3, 2)
	require.NotEqual(t, ptypes3[0].PropertyTypeID, ptypes3[1].PropertyTypeID)
	for _, value := range ptypes3 {
		ptyp := lType3.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName5:
			require.Equal(t, *value.RangeFromValue, 30.23)
			require.EqualValues(t, *value.RangeToValue, 50)
			require.Equal(t, ptyp.Type, "range")
		case ptyp.Name:
			require.Equal(t, *value.LatitudeValue, 45.8)
			require.Equal(t, *value.LongitudeValue, 88.9)
			require.Equal(t, ptyp.Type, "gps_location")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}
}

func TestValidateForExistingLocation(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	ids := prepareLocationTypesWithProperties(ctx, t, *r)

	firstRowLocations := append(append(locationIDHeader, []string{locTypeNameL, locTypeNameM, locTypeNameS}...), locationFixedDataHeader...)
	finalFirstRow := append(firstRowLocations, propName1, propName2, propName3, propName4, propName5, propName6)
	fl, _ := NewImportHeader(finalFirstRow, ImportEntityLocation)
	err := importer.inputValidationsLocation(ctx, fl)
	require.NoError(t, err)

	loc1, err := importer.r.Mutation().AddLocation(ctx, models.AddLocationInput{
		Name: "loc1L",
		Type: ids.locTypeIDL,
		Properties: []*models.PropertyInput{{
			RangeFromValue: pointer.ToFloat64(30),
			RangeToValue:   pointer.ToFloat64(50.88),
			PropertyTypeID: r.client.PropertyType.Query().Where(propertytype.Name(propName5)).OnlyXID(ctx),
		}},
	})
	require.NoError(t, err)
	loc2, err := importer.r.Mutation().AddLocation(ctx, models.AddLocationInput{
		Name: "loc2M",
		Type: ids.locTypeIDM,
		Properties: []*models.PropertyInput{{
			StringValue:    pointer.ToString("10/11/88"),
			PropertyTypeID: r.client.PropertyType.Query().Where(propertytype.Name(propName3)).OnlyXID(ctx),
		}},
		Latitude:  pointer.ToFloat64(16),
		Longitude: pointer.ToFloat64(44),
	})
	require.NoError(t, err)
	loc3, err := importer.r.Mutation().AddLocation(ctx, models.AddLocationInput{
		Name: "loc3S",
		Type: ids.locTypeIDS,
		Properties: []*models.PropertyInput{{
			IntValue:       pointer.ToInt(16),
			PropertyTypeID: r.client.PropertyType.Query().Where(propertytype.Name(propName2)).OnlyXID(ctx),
		}},
		Parent:     pointer.ToString(loc1.ID),
		ExternalID: pointer.ToString("123"),
	})
	require.NoError(t, err)

	var (
		test1 = []string{loc1.ID, "loc1L", "", "", "external_2", "", "", "", "", "", "", "30.23-50", ""}
		test2 = []string{loc2.ID, "", "loc2M", "", "", "32", "33", "", "", "29/03/88", "", "", ""}
		test3 = []string{loc3.ID, "loc1L", "", "loc3S", "external_1", "32", "33", "abc", "19", "", "", "", ""}
	)

	rec1 := NewImportRecord(test1, fl)
	_, err = importer.validateLineForExistingLocation(ctx, loc1.ID, rec1)
	require.NoError(t, err)

	inputs, _, err := importer.getLocationPropertyInputs(ctx, rec1, loc1)
	require.NoError(t, err)
	propInput := inputs[0]

	require.Equal(t, propName5, r.client.PropertyType.Query().Where(propertytype.ID(propInput.PropertyTypeID)).OnlyX(ctx).Name)
	require.Equal(t, 30.23, *propInput.RangeFromValue)
	require.Equal(t, 50.0, *propInput.RangeToValue)

	rec2 := NewImportRecord(test2, fl)
	_, err = importer.validateLineForExistingLocation(ctx, loc2.ID, rec2)
	require.NoError(t, err)

	inputs, _, err = importer.getLocationPropertyInputs(ctx, rec2, loc2)
	require.NoError(t, err)
	propInput = inputs[0]
	require.Equal(t, propName3, r.client.PropertyType.Query().Where(propertytype.ID(propInput.PropertyTypeID)).OnlyX(ctx).Name)
	require.Equal(t, "29/03/88", *propInput.StringValue)

	rec3 := NewImportRecord(test3, fl)
	_, err = importer.validateLineForExistingLocation(ctx, loc3.ID, rec3)
	require.NoError(t, err)
	inputs, _, err = importer.getLocationPropertyInputs(ctx, rec3, loc3)
	require.NoError(t, err)
	require.Len(t, inputs, 2)
	for _, inp := range inputs {
		typeName := r.client.PropertyType.Query().
			Where(propertytype.ID(inp.PropertyTypeID)).OnlyX(ctx).Name
		switch typeName {
		case propName1:
			require.Equal(t, "abc", *inp.StringValue)
		case propName2:
			require.Equal(t, 19, *inp.IntValue)
		default:
			require.Fail(t, "must  be one of the two")
		}
	}
}
