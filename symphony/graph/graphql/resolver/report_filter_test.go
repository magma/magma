// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package resolver

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

const (
	filterName       = "filter1"
	filterNameEdited = "filter1 - edited"
	substring        = "substring"
)

func validateEmptyFilters(t *testing.T, r *TestResolver) {
	qr := r.Query()
	ctx := viewertest.NewContext(context.Background(), r.client)
	for _, entity := range models.AllFilterEntity {
		filters, err := qr.ReportFilters(ctx, entity)
		require.NoError(t, err)
		require.Empty(t, filters)
	}
}

func TestAddReportFilter(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)

	mr, qr, rfr := r.Mutation(), r.Query(), r.ReportFilter()
	validateEmptyFilters(t, r)
	data := prepareEquipmentData(ctx, r, "A", nil)
	for _, entity := range models.AllFilterEntity {
		inp := models.ReportFilterInput{
			Name:   filterName,
			Entity: entity,
			Filters: []*models.GeneralFilterInput{
				{
					FilterType: models.EquipmentFilterTypeLocationInst.String(),
					Operator:   models.FilterOperatorIsOneOf,
					Key:        "for-ui-purposes",
					IDSet:      []int{data.loc1, data.loc2},
				},
				{
					FilterType:  models.EquipmentFilterTypeEquipInstName.String(),
					Operator:    models.FilterOperatorContains,
					Key:         "for-ui-purposes",
					StringValue: pointer.ToString(substring),
				},
			}}

		newFilter, err := mr.AddReportFilter(ctx, inp)
		switch entity {
		case models.FilterEntityEquipment:
			require.NoError(t, err)
			fetchedFilters, err := qr.ReportFilters(ctx, models.FilterEntityEquipment)
			require.NoError(t, err)
			require.Len(t, fetchedFilters, 1)
			fetchedReportFilter := fetchedFilters[0]
			require.Equal(t, newFilter.ID, fetchedReportFilter.ID)
			require.Equal(t, filterName, fetchedReportFilter.Name)

			fetchEnt, err := rfr.Entity(ctx, fetchedReportFilter)
			require.NoError(t, err)
			require.Equal(t, entity, fetchEnt)

			actualFilters, err := rfr.Filters(ctx, fetchedReportFilter)
			require.NoError(t, err)
			require.Len(t, actualFilters, 2)
			for _, f := range actualFilters {
				switch f.FilterType {
				case models.EquipmentFilterTypeLocationInst.String():
					require.Len(t, f.IDSet, 2)
					require.Contains(t, f.IDSet, data.loc1)
					require.Contains(t, f.IDSet, data.loc2)
					require.Equal(t, models.FilterOperatorIsOneOf, f.Operator)
				case models.EquipmentFilterTypeEquipInstName.String():
					require.Equal(t, f.StringValue, pointer.ToString(substring))
					require.Equal(t, models.FilterOperatorContains, f.Operator)
				default:
					require.Fail(t, "unsupported filter type %s", f.FilterType)
				}
			}
		default:
			require.Error(t, err, "entities does not match")
		}
	}
}

func TestAddInvalidReportFilters(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr := r.Mutation()
	validateEmptyFilters(t, r)
	data := prepareEquipmentData(ctx, r, "A", nil)
	inp := models.ReportFilterInput{
		Name:   filterName,
		Entity: models.FilterEntityEquipment,
		Filters: []*models.GeneralFilterInput{
			{
				FilterType: models.EquipmentFilterTypeLocationInst.String(),
				Operator:   models.FilterOperatorIsOneOf,
				Key:        "for-ui-purposes",
				IDSet:      []int{data.loc1, data.loc2},
			},
		}}

	_, err := mr.AddReportFilter(ctx, inp)
	require.NoError(t, err)
	duplInput := models.ReportFilterInput{
		Name:   filterName,
		Entity: models.FilterEntityEquipment,
		Filters: []*models.GeneralFilterInput{
			{
				FilterType:  models.EquipmentFilterTypeEquipInstName.String(),
				Operator:    models.FilterOperatorContains,
				Key:         "for-ui-purposes",
				StringValue: pointer.ToString(substring),
			},
		},
	}
	_, err = mr.AddReportFilter(ctx, duplInput)
	require.Error(t, err)

	inp = models.ReportFilterInput{
		Name:   filterName,
		Entity: models.FilterEntityEquipment,
		Filters: []*models.GeneralFilterInput{
			{
				FilterType:  "InvalidType",
				Operator:    models.FilterOperatorContains,
				Key:         "for-ui-purposes",
				StringValue: pointer.ToString(substring),
			},
		},
	}
	_, err = mr.AddReportFilter(ctx, inp)
	require.Error(t, err)
	inp = models.ReportFilterInput{
		Name:   filterName,
		Entity: models.FilterEntityEquipment,
		Filters: []*models.GeneralFilterInput{
			{
				FilterType:  models.EquipmentFilterTypeEquipInstName.String(),
				Operator:    "invalidOperator",
				Key:         "for-ui-purposes",
				StringValue: pointer.ToString(substring),
			},
		},
	}
	_, err = mr.AddReportFilter(ctx, inp)
	require.Error(t, err)
	inp = models.ReportFilterInput{
		Name:   filterName,
		Entity: "no entity",
		Filters: []*models.GeneralFilterInput{
			{
				FilterType:  models.EquipmentFilterTypeEquipInstName.String(),
				Operator:    models.FilterOperatorContains,
				Key:         "for-ui-purposes",
				StringValue: pointer.ToString(substring),
			},
		},
	}
	_, err = mr.AddReportFilter(ctx, inp)
	require.Error(t, err)
}

func TestEditReportFilters(t *testing.T) {
	r := newTestResolver(t)
	defer r.drv.Close()
	ctx := viewertest.NewContext(context.Background(), r.client)
	mr, qr := r.Mutation(), r.Query()
	validateEmptyFilters(t, r)
	data := prepareEquipmentData(ctx, r, "A", nil)
	inp := models.ReportFilterInput{
		Name:   filterName,
		Entity: models.FilterEntityEquipment,
		Filters: []*models.GeneralFilterInput{
			{
				FilterType: models.EquipmentFilterTypeLocationInst.String(),
				Operator:   models.FilterOperatorIsOneOf,
				Key:        "for-ui-purposes",
				IDSet:      []int{data.loc1, data.loc2},
			},
		}}
	newFilter, err := mr.AddReportFilter(ctx, inp)
	require.NoError(t, err)
	fetchedFilters, err := qr.ReportFilters(ctx, models.FilterEntityEquipment)
	require.NoError(t, err)
	require.Len(t, fetchedFilters, 1)
	fetchedReportFilter := fetchedFilters[0]
	require.Equal(t, newFilter.ID, fetchedReportFilter.ID)
	require.Equal(t, filterName, fetchedReportFilter.Name)

	editInput := models.EditReportFilterInput{
		ID:   newFilter.ID,
		Name: filterNameEdited,
	}
	newFilter, err = mr.EditReportFilter(ctx, editInput)
	require.NoError(t, err)
	fetchedFilters, err = qr.ReportFilters(ctx, models.FilterEntityEquipment)
	require.NoError(t, err)
	require.Len(t, fetchedFilters, 1)
	fetchedReportFilter = fetchedFilters[0]
	require.Equal(t, newFilter.ID, fetchedReportFilter.ID)
	require.Equal(t, filterNameEdited, fetchedReportFilter.Name)
}
