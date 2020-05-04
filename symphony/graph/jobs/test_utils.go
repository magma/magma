// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package jobs

import (
	"context"

	"github.com/AlekSi/pointer"

	"github.com/facebookincubator/ent/dialect"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
)

type (
	equipmentDataModels struct {
		locType  *ent.LocationType
		loc1     *ent.Location
		loc2     *ent.Location
		loc3     *ent.Location
		equType1 *ent.EquipmentType
		equType2 *ent.EquipmentType
		equType3 *ent.EquipmentType
		equ1     *ent.Equipment
		equ2     *ent.Equipment
		equ3     *ent.Equipment
	}
	serviceTypeDataModels struct {
		st1 *ent.ServiceType
		st2 *ent.ServiceType
	}
	linkDataModels struct {
		l1 *ent.Link
		l2 *ent.Link
	}
)

//TestJobsResolver contains data for jobs resolver
type TestJobsResolver struct {
	drv        dialect.Driver
	client     *ent.Client
	jobsRunner jobs
}

func prepareEquipmentData(ctx context.Context, r TestJobsResolver, name string) equipmentDataModels {
	mr := r.jobsRunner.r.Mutation()
	locType1, _ := mr.AddLocationType(ctx, models.AddLocationTypeInput{
		Name: name + "loc_type1",
	})

	loc1, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst1",
		Type: locType1.ID,
	})
	loc2, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst2",
		Type: locType1.ID,
	})
	loc3, _ := mr.AddLocation(ctx, models.AddLocationInput{
		Name: name + "loc_inst3",
		Type: locType1.ID,
	})
	portTypes := []*models.EquipmentPortInput{
		{
			Name: "p1",
		},
		{
			Name: "p2",
		},
		{
			Name: "p3",
		},
	}
	equType1, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  name + "eq_type",
		Ports: portTypes,
	})

	equ1, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     name + "eq_inst1",
		Type:     equType1.ID,
		Location: &loc1.ID,
	})
	equType2, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  name + "eq_type2",
		Ports: portTypes,
	})

	equ2, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     name + "eq_inst2",
		Type:     equType2.ID,
		Location: &loc2.ID,
	})

	equType3, _ := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  name + "eq_type3",
		Ports: portTypes,
	})

	equ3, _ := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     name + "eq_inst3",
		Type:     equType3.ID,
		Location: &loc3.ID,
	})
	return equipmentDataModels{
		locType1,
		loc1,
		loc2,
		loc3,
		equType1,
		equType2,
		equType3,
		equ1,
		equ2,
		equ3,
	}
}

func prepareServiceTypeData(ctx context.Context, r TestJobsResolver, equipData equipmentDataModels) serviceTypeDataModels {
	mr := r.jobsRunner.r.Mutation()
	dm := models.DiscoveryMethodInventory
	srvType1, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:            "test service type1",
		HasCustomer:     false,
		DiscoveryMethod: &dm,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("PROVIDER1"),
				Index:           0,
				EquipmentTypeID: equipData.equType1.ID,
			},
			{
				Index:           1,
				Name:            "endpoint type2",
				Role:            pointer.ToString("MIDDLE1"),
				EquipmentTypeID: equipData.equType2.ID,
			},
			{
				Index:           2,
				Name:            "endpoint type3",
				Role:            pointer.ToString("CONSUMER1"),
				EquipmentTypeID: equipData.equType3.ID,
			},
		},
	})
	srvType2, _ := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:        "test service type2",
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				Role:            pointer.ToString("PROVIDER2"),
				Index:           0,
				EquipmentTypeID: equipData.equType2.ID,
			},
			{
				Index:           1,
				Name:            "endpoint type2",
				Role:            pointer.ToString("MIDDLE2"),
				EquipmentTypeID: equipData.equType3.ID,
			},
			{
				Index:           2,
				Name:            "endpoint type3",
				Role:            pointer.ToString("CONSUMER2"),
				EquipmentTypeID: equipData.equType1.ID,
			},
		},
	})
	return serviceTypeDataModels{
		st1: srvType1,
		st2: srvType2,
	}
}

func prepareLinksData(ctx context.Context, r TestJobsResolver, equipData equipmentDataModels) linkDataModels {
	mr := r.jobsRunner.r.Mutation()
	l1, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: equipData.equ1.ID,
				Port: equipData.equ1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p1")).
					OnlyXID(ctx),
			},
			{
				Equipment: equipData.equ2.ID,
				Port: equipData.equ2.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p1")).
					OnlyXID(ctx),
			},
		}})
	l2, _ := mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: equipData.equ2.ID,
				Port: equipData.equ2.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p2")).
					OnlyXID(ctx),
			},
			{
				Equipment: equipData.equ3.ID,
				Port: equipData.equ3.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p2")).
					OnlyXID(ctx),
			},
		}})
	return linkDataModels{
		l1: l1,
		l2: l2,
	}
}
