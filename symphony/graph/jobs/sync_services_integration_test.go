package jobs

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"

	"github.com/AlekSi/pointer"
	"github.com/stretchr/testify/require"

	"github.com/facebookincubator/symphony/graph/ent/serviceendpointdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
)

func TestAddNewService(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)

	sCount := r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)
	syncServicesRequest(t, r)
	services := r.client.Service.Query().AllX(ctx)
	require.Len(t, services, 1)
	s := services[0]

	require.Equal(t, sData.st1.ID, s.QueryType().OnlyXID(ctx))
	for _, ep := range s.QueryEndpoints().AllX(ctx) {
		e := ep.QueryEquipment().OnlyX(ctx)
		switch ep.QueryDefinition().OnlyX(ctx).Index {
		case 0:
			require.Equal(t, eData.equ1.ID, e.ID)
		case 1:
			require.Equal(t, eData.equ2.ID, e.ID)
		case 2:
			require.Equal(t, eData.equ3.ID, e.ID)
		default:
			require.Fail(t, "no valid index")
		}
	}
}

func TestRemoveLinkAndAddNewLink(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	ldata := prepareLinksData(ctx, *r, eData)

	sCount := r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Equal(t, 1, sCount)
	_, err := r.jobsRunner.r.Mutation().RemoveLink(ctx, ldata.l1.ID, nil)
	// links now : eq2 => eq3 => eq2 => eq1
	require.NoError(t, err)
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Equal(t, 1, sCount)

	_, err = r.jobsRunner.r.Mutation().RemoveLink(ctx, ldata.l3.ID, nil)
	require.NoError(t, err)
	// links now : eq3 => eq2 => eq1

	_, err = r.jobsRunner.r.Mutation().RemoveLink(ctx, ldata.l4.ID, nil)
	require.NoError(t, err)
	// links now : eq3 => eq2
	require.NoError(t, err)
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)

	_, err = r.jobsRunner.r.Mutation().AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{
				Equipment: eData.equ3.ID,
				Port: eData.equ3.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p6")).
					OnlyXID(ctx),
			},
			{
				Equipment: eData.equ1.ID,
				Port: eData.equ1.QueryType().QueryPortDefinitions().
					Where(equipmentportdefinition.Name("p6")).
					OnlyXID(ctx),
			},
		}})
	require.NoError(t, err)
	// links now : eq2 => eq3 => eq1
	syncServicesRequest(t, r)
	services := r.client.Service.Query().AllX(ctx)
	require.Len(t, services, 1)
	s := services[0]

	require.Equal(t, sData.st2.ID, s.QueryType().OnlyXID(ctx))
	for _, ep := range s.QueryEndpoints().AllX(ctx) {
		e := ep.QueryEquipment().OnlyX(ctx)
		switch ep.QueryDefinition().OnlyX(ctx).Index {
		case 0:
			require.Equal(t, eData.equ2.ID, e.ID)
		case 1:
			require.Equal(t, eData.equ3.ID, e.ID)
		case 2:
			require.Equal(t, eData.equ1.ID, e.ID)
		default:
			require.Fail(t, "no valid index")
		}
	}
}

func TestEditServiceTypeEndpointDefinitionsOrder(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)
	syncServicesRequest(t, r)
	sCount := r.client.Service.Query().CountX(ctx)
	require.Equal(t, 1, sCount)

	_, err := r.jobsRunner.r.Mutation().EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          sData.st1.ID,
		Name:        sData.st1.Name,
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type1")).OnlyXID(ctx)),
				Role:            pointer.ToString("PROVIDER2"),
				Index:           0,
				EquipmentTypeID: eData.equType2.ID,
			},
			{
				Index:           1,
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type3")).OnlyXID(ctx)),
				Name:            "endpoint type3",
				Role:            pointer.ToString("CONSUMER2"),
				EquipmentTypeID: eData.equType1.ID,
			},
			{
				Index:           2,
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type2")).OnlyXID(ctx)),
				Name:            "endpoint type2",
				Role:            pointer.ToString("MIDDLE2"),
				EquipmentTypeID: eData.equType3.ID,
			},
		},
	})
	require.NoError(t, err)
	// service type 1 now: eqType2 => eqType1 => eqType3
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)

	_, err = r.jobsRunner.r.Mutation().EditServiceType(ctx, models.ServiceTypeEditData{
		ID:          sData.st1.ID,
		Name:        sData.st1.Name,
		HasCustomer: false,
		Endpoints: []*models.ServiceEndpointDefinitionInput{
			{
				Name:            "endpoint type1",
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type1")).OnlyXID(ctx)),
				Role:            pointer.ToString("PROVIDER2"),
				Index:           0,
				EquipmentTypeID: eData.equType3.ID,
			},
			{
				Index:           1,
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type3")).OnlyXID(ctx)),
				Name:            "endpoint type3",
				Role:            pointer.ToString("CONSUMER2"),
				EquipmentTypeID: eData.equType2.ID,
			},
			{
				Index:           2,
				ID:              pointer.ToInt(sData.st1.QueryEndpointDefinitions().Where(serviceendpointdefinition.Name("endpoint type2")).OnlyXID(ctx)),
				Name:            "endpoint type2",
				Role:            pointer.ToString("MIDDLE2"),
				EquipmentTypeID: eData.equType1.ID,
			},
		},
	})
	require.NoError(t, err)
	// service type 1 now: eqType3 => eqType2 => eqType1
	syncServicesRequest(t, r)
	services := r.client.Service.Query().AllX(ctx)
	require.Len(t, services, 1)
	s := services[0]

	require.Equal(t, sData.st1.ID, s.QueryType().OnlyXID(ctx))
	for _, ep := range s.QueryEndpoints().AllX(ctx) {
		e := ep.QueryEquipment().OnlyX(ctx)
		switch ep.QueryDefinition().OnlyX(ctx).Index {
		case 0:
			require.Equal(t, eData.equ3.ID, e.ID)
		case 1:
			require.Equal(t, eData.equ2.ID, e.ID)
		case 2:
			require.Equal(t, eData.equ1.ID, e.ID)
		default:
			require.Fail(t, "no valid index")
		}
	}
}

func TestDeletedServiceType(t *testing.T) {
	r := newJobsTestResolver(t)
	defer r.drv.Close()

	ctx := newServicesContext(viewertest.NewContext(context.Background(), r.client))
	eData := prepareEquipmentData(ctx, *r, "A")
	sData := prepareServiceTypeData(ctx, *r, eData)
	prepareLinksData(ctx, *r, eData)

	sCount := r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)
	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Equal(t, 1, sCount)

	r.client.ServiceType.UpdateOneID(sData.st1.ID).SetIsDeleted(true).ExecX(ctx)

	syncServicesRequest(t, r)
	sCount = r.client.Service.Query().CountX(ctx)
	require.Zero(t, sCount)
}
