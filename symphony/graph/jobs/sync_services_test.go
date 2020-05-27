// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestNextEquipmentInstances(t *testing.T) {
	r := newJobsTestResolver(t)
	mr := r.jobsRunner.r.Mutation()
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)

	newEq1, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "new1",
		Type:     eData.equType1.ID,
		Location: pointer.ToInt(eData.loc1.ID),
	})
	require.NoError(t, err)

	_, err = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "new2",
		Type:     eData.equType2.ID,
		Location: pointer.ToInt(eData.loc2.ID),
	})
	require.NoError(t, err)

	_, err = mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "new3",
		Type:     eData.equType3.ID,
		Location: pointer.ToInt(eData.loc3.ID),
	})
	require.NoError(t, err)

	equipArr, err := getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(0)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Empty(t, equipArr)
	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(1)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Len(t, equipArr, 1)
	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(2)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Empty(t, equipArr)

	_, err = r.jobsRunner.r.Mutation().AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: newEq1.ID,
				Port: newEq1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p1")).
					OnlyXID(ctx),
			},
			{
				Equipment: eData.equ1.ID,
				Port: eData.equ1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p4")).
					OnlyXID(ctx),
			},
		}})
	require.NoError(t, err)

	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(0)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Len(t, equipArr, 1)
	// link within the same equipment - isn't counted
	_, err = r.jobsRunner.r.Mutation().AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: eData.equ1.ID,
				Port: eData.equ1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p6")).
					OnlyXID(ctx),
			},
			{
				Equipment: eData.equ1.ID,
				Port: eData.equ1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p5")).
					OnlyXID(ctx),
			},
		}})
	require.NoError(t, err)

	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(0)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Len(t, equipArr, 1)

	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(1)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Len(t, equipArr, 1)
	equipArr, err = getNextEquipmentInstances(
		ctx,
		eData.equ1,
		sData.st1.
			QueryEndpointDefinitions().
			Where(serviceendpointdefinition.Index(2)).OnlyX(ctx),
		nil,
	)
	require.NoError(t, err)
	require.Empty(t, equipArr)
}

func TestGetServiceDetailsList(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()
	mr := r.jobsRunner.r.Mutation()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	longServiceData := prepareLongServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)
	toAdd, err := r.jobsRunner.getServicesDetailsList(ctx, longServiceData.stTooLong)
	require.NoError(t, err)
	require.Empty(t, toAdd)

	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, sData.st1)
	require.NoError(t, err)
	require.Len(t, toAdd, 1)

	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, sData.st2)
	require.NoError(t, err)
	require.Empty(t, toAdd)

	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, longServiceData.stLong)
	require.NoError(t, err)
	require.Len(t, toAdd, 1)

	// test that consecutive links in service needs to be different
	stWithReverse, err := mr.EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          longServiceData.stLong.ID,
		Name:        longServiceData.stLong.Name,
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				ID:              &longServiceData.stLong.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(0)).OnlyX(ctx).ID,
				Name:            "endpoint type1",
				Role:            pointer.ToString("PROVIDER2"),
				Index:           0,
				EquipmentTypeID: eData.equType1.ID,
			},
			{
				ID: &longServiceData.stLong.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(1)).OnlyX(ctx).ID,

				Index:           1,
				Name:            "endpoint type2",
				Role:            pointer.ToString("MIDDLE2"),
				EquipmentTypeID: eData.equType2.ID,
			},
			{
				ID: &longServiceData.stLong.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(2)).OnlyX(ctx).ID,

				Index:           2,
				Name:            "endpoint type3",
				Role:            pointer.ToString("CONSUMER2"),
				EquipmentTypeID: eData.equType3.ID,
			},
			{
				ID: &longServiceData.stLong.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(3)).OnlyX(ctx).ID,

				Index:           3,
				Name:            "endpoint type4",
				Role:            pointer.ToString("CONSUMER3"),
				EquipmentTypeID: longServiceData.eType4.ID,
			},
			{
				ID:              &longServiceData.stLong.QueryEndpointDefinitions().Where(serviceendpointdefinition.Index(4)).OnlyX(ctx).ID,
				Index:           4,
				Name:            "endpoint type5 - reverse",
				Role:            pointer.ToString("CONSUMER4"),
				EquipmentTypeID: eData.equType3.ID,
			},
		},
	})
	require.NoError(t, err)
	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, stWithReverse)
	require.NoError(t, err)
	require.Empty(t, toAdd)

	// test that consecutive links in service needs to be from  different equipment
	_, err = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: longServiceData.equ4.ID,
				Port: longServiceData.equ4.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p3")).
					OnlyXID(ctx),
			},
			{
				Equipment: eData.equ3.ID,
				Port: eData.equ3.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p4")).
					OnlyXID(ctx),
			},
		}})
	require.NoError(t, err)
	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, stWithReverse)
	require.NoError(t, err)
	require.Empty(t, toAdd)

	// test that consecutive links in service can be of same equipment type (if instance is different)
	newE3, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     "new e3",
		Type:     eData.equType3.ID,
		Location: pointer.ToInt(eData.loc1.ID),
	})
	require.NoError(t, err)
	_, err = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: longServiceData.equ4.ID,
				Port: longServiceData.equ4.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p4")).
					OnlyXID(ctx),
			},
			{
				Equipment: newE3.ID,
				Port: newE3.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p1")).
					OnlyXID(ctx),
			},
		}})
	require.NoError(t, err)

	toAdd, err = r.jobsRunner.getServicesDetailsList(ctx, stWithReverse)
	require.NoError(t, err)
	require.Len(t, toAdd, 1)
}
