// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"context"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

const (
	serviceTypeName  = "serviceType"
	serviceType2Name = "serviceType2"
	serviceType3Name = "serviceType3"
	serviceType4Name = "serviceType4"
)

type serviceIds struct {
	serviceTypeID  string
	serviceTypeID2 string
	serviceTypeID3 string
	serviceTypeID4 string
}

func prepareServiceTypeData(ctx context.Context, t *testing.T, r TestImporterResolver) serviceIds {
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
	propDefInput7 := models.PropertyTypeInput{
		Name: propName7,
		Type: "location",
	}
	propDefInput8 := models.PropertyTypeInput{
		Name: propName8,
		Type: "service",
	}

	serviceType1, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:       serviceTypeName,
		Properties: []*models.PropertyTypeInput{&propDefInput1, &propDefInput2},
	})
	require.NoError(t, err)
	serviceType2, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:       serviceType2Name,
		Properties: []*models.PropertyTypeInput{&propDefInput3, &propDefInput4},
	})
	require.NoError(t, err)
	serviceType3, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:       serviceType3Name,
		Properties: []*models.PropertyTypeInput{&propDefInput5, &propDefInput6},
	})
	require.NoError(t, err)
	serviceType4, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{
		Name:       serviceType4Name,
		Properties: []*models.PropertyTypeInput{&propDefInput7, &propDefInput8},
	})
	require.NoError(t, err)
	return serviceIds{
		serviceTypeID:  serviceType1.ID,
		serviceTypeID2: serviceType2.ID,
		serviceTypeID3: serviceType3.ID,
		serviceTypeID4: serviceType4.ID,
	}
}

func TestValidatePropertiesForServiceType(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	q := r.importer.r.Query()
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	data := prepareServiceTypeData(ctx, t, *r)

	mr := r.importer.r.Mutation()
	locType, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: "City"})
	require.NoError(t, err)
	loc, err := mr.AddLocation(ctx, models.AddLocationInput{Name: "New York", Type: locType.ID})
	require.NoError(t, err)

	serviceType, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "L2 Access", HasCustomer: false})
	require.NoError(t, err)
	service, err := mr.AddService(ctx, models.ServiceCreateData{
		Name:          "Service23",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)

	var (
		dataHeader = [...]string{"Service ID", "Service Name", "Service Type", "Service External ID", "Customer Name", "Customer External ID", "Status"}
		row1       = []string{"", "s1", serviceTypeName, "M123", "", "", "IN_SERVICE", "strVal", "54", "", "", "", "", "", ""}
		row2       = []string{"", "s2", serviceType2Name, "M456", "", "", "MAINTENANCE", "", "", "29/03/88", "false", "", "", "", ""}
		row3       = []string{"", "s3", serviceType3Name, "M789", "", "", "DISCONNECTED", "", "", "", "", "30.23-50", "45.8,88.9", "", ""}
		row4       = []string{"", "s3", serviceType4Name, "M789", "", "", "DISCONNECTED", "", "", "", "", "", "", loc.ID, service.Name}
	)

	titleWithProperties := append(dataHeader[:], propName1, propName2, propName3, propName4, propName5, propName6, propName7, propName8)
	fl, _ := NewImportHeader(titleWithProperties, ImportEntityService)
	r1 := NewImportRecord(row1, fl)
	require.NoError(t, err)
	styp1, err := q.ServiceType(ctx, data.serviceTypeID)
	require.NoError(t, err)
	ptypes, err := importer.validatePropertiesForServiceType(ctx, r1, styp1)
	require.NoError(t, err)
	require.Len(t, ptypes, 2)
	require.NotEqual(t, ptypes[0].PropertyTypeID, ptypes[1].PropertyTypeID)
	for _, value := range ptypes {
		ptyp := styp1.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
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
	styp2, err := q.ServiceType(ctx, data.serviceTypeID2)
	require.NoError(t, err)

	r2 := NewImportRecord(row2, fl)
	ptypes2, err := importer.validatePropertiesForServiceType(ctx, r2, styp2)
	require.NoError(t, err)
	require.Len(t, ptypes2, 2)
	for _, value := range ptypes2 {
		ptyp := styp2.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
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

	styp3, err := q.ServiceType(ctx, data.serviceTypeID3)
	require.NoError(t, err)

	r3 := NewImportRecord(row3, fl)
	ptypes3, err := importer.validatePropertiesForServiceType(ctx, r3, styp3)
	require.NoError(t, err)
	require.Len(t, ptypes3, 2)
	require.NotEqual(t, ptypes3[0].PropertyTypeID, ptypes3[1].PropertyTypeID)
	for _, value := range ptypes3 {
		ptyp := styp3.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName5:
			require.Equal(t, *value.RangeFromValue, 30.23)
			require.EqualValues(t, *value.RangeToValue, 50)
			require.Equal(t, ptyp.Type, "range")
		case propName6:
			require.Equal(t, *value.LatitudeValue, 45.8)
			require.Equal(t, *value.LongitudeValue, 88.9)
			require.Equal(t, ptyp.Type, "gps_location")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}

	styp4, err := q.ServiceType(ctx, data.serviceTypeID4)
	require.NoError(t, err)

	r4 := NewImportRecord(row4, fl)
	ptypes4, err := importer.validatePropertiesForServiceType(ctx, r4, styp4)
	require.NoError(t, err)
	require.Len(t, ptypes4, 2)
	require.NotEqual(t, ptypes4[0].PropertyTypeID, ptypes4[1].PropertyTypeID)
	for _, value := range ptypes4 {
		ptyp := styp4.QueryPropertyTypes().Where(propertytype.ID(value.PropertyTypeID)).OnlyX(ctx)
		switch ptyp.Name {
		case propName7:
			require.Equal(t, *value.LocationIDValue, loc.ID)
			require.Equal(t, ptyp.Type, "location")
		case propName8:
			require.Equal(t, *value.ServiceIDValue, service.ID)
			require.Equal(t, ptyp.Type, "service")
		default:
			require.Fail(t, "property type name should be one of the two")
		}
	}
}

func TestValidateForExistingService(t *testing.T) {
	r := newImporterTestResolver(t)
	importer := r.importer
	defer r.drv.Close()
	ctx := newImportContext(viewertest.NewContext(r.client))
	prepareServiceTypeData(ctx, t, *r)

	titleWithProperties := []string{"Service ID", "Service Name", "Service Type", "Service External ID", "Customer Name", "Customer External ID", "Status"}
	title, _ := NewImportHeader(titleWithProperties, ImportEntityService)

	serviceType, err := importer.r.Mutation().AddServiceType(ctx, models.ServiceTypeCreateData{
		Name: "type1",
	})
	require.NoError(t, err)
	service, err := importer.r.Mutation().AddService(ctx, models.ServiceCreateData{
		Name:          "myService",
		ServiceTypeID: serviceType.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusPending),
	})
	require.NoError(t, err)
	var (
		test = []string{service.ID, "myService", "type1", "", "", "", "", models.ServiceStatusPending.String()}
	)
	_, err = importer.validateLineForExistingService(ctx, service.ID, NewImportRecord(test, title))
	require.NoError(t, err)
}
