// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"
)

type locationSearchDataModels struct {
	loc1     *ent.Location
	loc2     *ent.Location
	locType1 *ent.LocationType
	locType2 *ent.LocationType
}

// nolint: errcheck
func prepareLocationData(ctx context.Context, r *TestResolver) locationSearchDataModels {
	mr := r.Mutation()
	locType1, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type1",
		Properties: []*models.PropertyTypeInput{
			{
				Name: "date_established",
				Type: models.PropertyKindDate,
			},
			{
				Name: "stringProp",
				Type: models.PropertyKindString,
			},
		},
	})
	datePropDef := locType1.QueryPropertyTypes().Where(propertytype.Name("date_established")).OnlyX(ctx)
	strPropDef := locType1.QueryPropertyTypes().Where(propertytype.Name("stringProp")).OnlyX(ctx)

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: "loc_inst1",
		Type: locType1.ID,
		Properties: []*models.PropertyInput{
			{
				PropertyTypeID: datePropDef.ID,
				StringValue:    pointer.ToString("1988-03-29"),
			},
			{
				PropertyTypeID: strPropDef.ID,
				StringValue:    pointer.ToString("testProp"),
			},
		},
	})

	locType2, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: "loc_type2",
	})

	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   "loc_inst2",
		Type:   locType2.ID,
		Parent: pointer.ToInt(loc1.ID),
	})

	equType, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "eq_type",
	})
	if _, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "eq_inst",
		Type:     equType.ID,
		Location: &loc1.ID,
	}); err != nil {
		panic(err)
	}
	return locationSearchDataModels{
		loc1,
		loc2,
		locType1,
		locType2,
	}
}

func TestSearchLocationAncestors(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	data := prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1):
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()
	limit := 100
	all, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{}, &limit)
	require.NoError(t, err)
	require.Len(t, all.Locations, 2)
	require.Equal(t, all.Count, 2)
	maxDepth := 2
	f1 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.loc1.ID},
		MaxDepth:   &maxDepth,
	}
	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, &limit)
	require.NoError(t, err)
	require.Len(t, res.Locations, 2)
	require.Equal(t, res.Count, 2)

	f2 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.loc2.ID},
		MaxDepth:   &maxDepth,
	}
	res, err = qr.LocationSearch(ctx, []*models.LocationFilterInput{&f2}, &limit)
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}

func TestSearchLocationByName(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	data := prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1):
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()

	f1 := models.LocationFilterInput{
		FilterType:  models.LocationFilterTypeLocationInstName,
		Operator:    models.FilterOperatorIs,
		StringValue: &data.loc2.Name,
	}
	resAll, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, resAll.Locations, 2)
	require.Equal(t, resAll.Count, 2)

	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}

func TestSearchLocationByType(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	data := prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1):
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()
	f1 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.locType2.ID},
	}
	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}

func TestSearchLocationHasEquipment(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1):
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()
	f1 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationInstHasEquipment,
		Operator:   models.FilterOperatorIs,
		BoolValue:  pointer.ToBool(true),
	}
	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)

	f2 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationInstHasEquipment,
		Operator:   models.FilterOperatorIs,
		BoolValue:  pointer.ToBool(false),
	}
	res, err = qr.LocationSearch(ctx, []*models.LocationFilterInput{&f2}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}

func TestSearchMultipleFilters(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	data := prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1):
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()
	f1 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationInst,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.loc1.ID},
		MaxDepth:   pointer.ToInt(2),
	}
	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 2)
	require.Equal(t, res.Count, 2)

	f2 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeLocationType,
		Operator:   models.FilterOperatorIsOneOf,
		IDSet:      []int{data.locType2.ID},
	}
	res, err = qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1, &f2}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}

func TestSearchLocationProperties(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(r.client)

	prepareLocationData(ctx, r)
	/*
		helper: data now is of type:
		 loc1 (loc_type1): - properties
			eq_inst (eq_type)
			loc2 (loc_type2)
	*/
	qr := r.Query()
	f1 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeProperty,
		Operator:   models.FilterOperatorDateLessThan,
		PropertyValue: &models.PropertyTypeInput{
			Type:        models.PropertyKindDate,
			Name:        "date_established",
			StringValue: pointer.ToString("2019-11-15"),
		},
	}

	res, err := qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)

	f2 := models.LocationFilterInput{
		FilterType: models.LocationFilterTypeProperty,
		Operator:   models.FilterOperatorIs,
		PropertyValue: &models.PropertyTypeInput{
			Type:        models.PropertyKindString,
			Name:        "stringProp",
			StringValue: pointer.ToString("testProp"),
		},
	}
	res, err = qr.LocationSearch(ctx, []*models.LocationFilterInput{&f1, &f2}, pointer.ToInt(100))
	require.NoError(t, err)
	require.Len(t, res.Locations, 1)
	require.Equal(t, res.Count, 1)
}
