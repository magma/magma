// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

const (
	serviceNameTitle        = "Service Name"
	serviceTypeTitle        = "Service Type"
	serviceExternalIDTitle  = "Service External ID"
	customerNameTitle       = "Customer Name"
	customerExternalIDTitle = "Customer External ID"
	statusTitle             = "Status"
	strPropTitle            = "service_str_prop"
	intPropTitle            = "service_int_prop"
	boolPropTitle           = "service_bool_prop"
	floatPropTitle          = "service_float_prop"
)

func pointerToServiceStatus(status models.ServiceStatus) *models.ServiceStatus {
	return &status
}

func preparePropertyTypes() []*models.PropertyTypeInput {
	serviceStrPropType := models.PropertyTypeInput{
		Name:        strPropTitle,
		Type:        "string",
		StringValue: pointer.ToString("Foo is the best"),
	}
	serviceIntPropType := models.PropertyTypeInput{
		Name: intPropTitle,
		Type: "int",
	}
	serviceBoolPropType := models.PropertyTypeInput{
		Name: boolPropTitle,
		Type: "bool",
	}
	serviceFloatPropType := models.PropertyTypeInput{
		Name: floatPropTitle,
		Type: "float",
	}

	return []*models.PropertyTypeInput{
		&serviceStrPropType,
		&serviceIntPropType,
		&serviceBoolPropType,
		&serviceFloatPropType,
	}
}

func prepareServiceData(ctx context.Context, t *testing.T, r TestExporterResolver) {
	mr := r.Mutation()

	serviceType1, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "L2 Service", HasCustomer: false})
	require.NoError(t, err)
	serviceType2, err := mr.AddServiceType(ctx, models.ServiceTypeCreateData{Name: "L3 Service", HasCustomer: true, Properties: preparePropertyTypes()})
	require.NoError(t, err)

	strType, _ := serviceType2.QueryPropertyTypes().Where(propertytype.Name(strPropTitle)).Only(ctx)
	intType, _ := serviceType2.QueryPropertyTypes().Where(propertytype.Name(intPropTitle)).Only(ctx)
	boolType, _ := serviceType2.QueryPropertyTypes().Where(propertytype.Name(boolPropTitle)).Only(ctx)
	floatType, _ := serviceType2.QueryPropertyTypes().Where(propertytype.Name(floatPropTitle)).Only(ctx)

	customer1, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name:       "Customer 1",
		ExternalID: pointer.ToString("AD123"),
	})
	require.NoError(t, err)

	customer2, err := mr.AddCustomer(ctx, models.AddCustomerInput{
		Name: "Customer 2",
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "L2 S1",
		ExternalID:    pointer.ToString("XS542"),
		ServiceTypeID: serviceType1.ID,
		Status:        pointerToServiceStatus(models.ServiceStatusInService),
	})
	require.NoError(t, err)

	strProp := models.PropertyInput{
		PropertyTypeID: strType.ID,
		StringValue:    pointer.ToString("Foo"),
	}
	intProp := models.PropertyInput{
		PropertyTypeID: intType.ID,
		IntValue:       pointer.ToInt(10),
	}

	boolProp := models.PropertyInput{
		PropertyTypeID: boolType.ID,
		BooleanValue:   pointer.ToBool(false),
	}

	floatProp := models.PropertyInput{
		PropertyTypeID: floatType.ID,
		FloatValue:     pointer.ToFloat64(3.5),
	}

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "L3 S1",
		ServiceTypeID: serviceType2.ID,
		CustomerID:    &customer1.ID,
		Properties:    []*models.PropertyInput{&strProp, &intProp, &boolProp},
		Status:        pointerToServiceStatus(models.ServiceStatusMaintenance),
	})
	require.NoError(t, err)

	_, err = mr.AddService(ctx, models.ServiceCreateData{
		Name:          "L3 S2",
		ServiceTypeID: serviceType2.ID,
		CustomerID:    &customer2.ID,
		Properties:    []*models.PropertyInput{&floatProp},
		Status:        pointerToServiceStatus(models.ServiceStatusDisconnected),
	})
	require.NoError(t, err)
}

func TestEmptyServicesDataExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, servicesRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	viewertest.SetDefaultViewerHeaders(req)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	for {
		ln, err := reader.Read()
		if err == io.EOF {
			break
		}
		require.NoError(t, err, "error reading row")
		require.EqualValues(t, []string{
			"\ufeffService ID",
			serviceNameTitle,
			serviceTypeTitle,
			serviceExternalIDTitle,
			customerNameTitle,
			customerExternalIDTitle,
			statusTitle,
		}, ln)
	}
}

func TestServicesExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, servicesRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	viewertest.SetDefaultViewerHeaders(req)

	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareServiceData(ctx, t, *r)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	for {
		ln, err := reader.Read()
		if err == io.EOF {
			break
		}
		require.NoError(t, err, "error reading row")
		switch {
		case ln[1] == serviceNameTitle:
			require.EqualValues(t, []string{
				"\ufeffService ID",
				serviceNameTitle,
				serviceTypeTitle,
				serviceExternalIDTitle,
				customerNameTitle,
				customerExternalIDTitle,
				statusTitle,
				strPropTitle,
				intPropTitle,
				boolPropTitle,
				floatPropTitle,
			}, ln)
		case ln[1] == "L2 S1":
			require.EqualValues(t, ln[1:], []string{
				"L2 S1",
				"L2 Service",
				"XS542",
				"",
				"",
				models.ServiceStatusInService.String(),
				"",
				"",
				"",
				"",
			})
		case ln[1] == "L3 S1":
			require.EqualValues(t, ln[1:], []string{
				"L3 S1",
				"L3 Service",
				"",
				"Customer 1",
				"AD123",
				models.ServiceStatusMaintenance.String(),
				"Foo",
				"10",
				"false",
				"0.000",
			})
		case ln[1] == "L3 S2":
			require.EqualValues(t, ln[1:], []string{
				"L3 S2",
				"L3 Service",
				"",
				"Customer 2",
				"",
				models.ServiceStatusDisconnected.String(),
				"Foo is the best",
				"0",
				"false",
				"3.500",
			})
		default:
			require.Fail(t, "line does not match")
		}
	}
}

func TestServiceWithFilters(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, servicesRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	viewertest.SetDefaultViewerHeaders(req)

	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareServiceData(ctx, t, *r)
	require.NoError(t, err)
	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	f1, err := json.Marshal([]servicesFilterInput{
		{
			Name:        models.ServiceFilterTypeServiceInstCustomerName,
			Operator:    models.FilterOperatorContains,
			StringValue: "Customer 1",
		},
	})
	require.NoError(t, err)
	f2, err := json.Marshal([]servicesFilterInput{
		{
			Name:        models.ServiceFilterTypeServiceInstExternalID,
			Operator:    models.FilterOperatorIs,
			StringValue: "XS542",
		},
	})
	require.NoError(t, err)

	for i, filter := range [][]byte{f1, f2} {
		req, err := http.NewRequest("GET", server.URL, nil)
		require.NoError(t, err)
		viewertest.SetDefaultViewerHeaders(req)

		q := req.URL.Query()
		q.Add("filters", string(filter))
		req.URL.RawQuery = q.Encode()

		res, err := http.DefaultClient.Do(req)
		require.NoError(t, err)

		reader := csv.NewReader(res.Body)
		linesCount := 0
		for {
			ln, err := reader.Read()
			if err == io.EOF {
				break
			}
			linesCount++
			require.NoError(t, err, "error reading row")
			if i == 0 {
				if ln[1] != serviceNameTitle {
					require.EqualValues(t, ln[1:], []string{
						"L3 S1",
						"L3 Service",
						"",
						"Customer 1",
						"AD123",
						models.ServiceStatusMaintenance.String(),
						"Foo",
						"10",
						"false",
						"0.000",
					})
				}
			}
			if i == 1 {
				if ln[1] != serviceNameTitle {
					require.EqualValues(t, ln[1:], []string{
						"L2 S1",
						"L2 Service",
						"XS542",
						"",
						"",
						models.ServiceStatusInService.String(),
					})
				}
			}
		}
		err = res.Body.Close()
		require.NoError(t, err)
	}
}
