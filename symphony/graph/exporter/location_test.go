// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestEmptyLocationDataExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, locationsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)

	req.Header.Set(tenantHeader, "fb-test")
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
			"\ufeffLocation ID",
			"External ID",
			"Latitude",
			"Longitude",
		}, ln)
	}
}

func TestLocationsExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, locationsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	req.Header.Set(tenantHeader, "fb-test")

	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareData(ctx, t, *r)
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
		case ln[1] == locTypeNameL:
			require.EqualValues(t, []string{
				"\ufeffLocation ID",
				locTypeNameL,
				locTypeNameM,
				locTypeNameS,
				"External ID",
				"Latitude",
				"Longitude",
				propNameStr,
				propNameBool,
				propNameDate,
			}, ln)
		case ln[4] == externalIDL:
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				"",
				"",
				externalIDL,
				fmt.Sprintf("%f", lat),
				fmt.Sprintf("%f", long),
				"",
				"",
				"",
			})
		case ln[4] == externalIDM:
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				parentLocation,
				"",
				externalIDM,
				fmt.Sprintf("%f", lat),
				fmt.Sprintf("%f", long),
				"",
				"",
				"",
			})
		case ln[3] == childLocation:
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				parentLocation,
				childLocation,
				"",
				fmt.Sprintf("%f", 0.0),
				fmt.Sprintf("%f", 0.0),
				"override",
				"true",
				"1988-03-29",
			})
		default:
			require.Fail(t, "line does not match")
		}
	}
}

func TestExportLocationWithFilters(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log
	ctx := viewertest.NewContext(context.Background(), r.client)
	e := &exporter{log, locationsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	prepareData(ctx, t, *r)

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set(tenantHeader, "fb-test")

	f, err := json.Marshal([]locationsFilterInput{
		{
			Name:      "LOCATION_INST_HAS_EQUIPMENT",
			Operator:  "IS",
			BoolValue: pointer.ToBool(false),
		},
	})
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("filters", string(f))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	linesCount := 0
	for {
		ln, err := reader.Read()
		if err == io.EOF {
			break
		}
		linesCount++
		require.NoError(t, err, "error reading row")
		switch ln[4] {
		case externalIDL:
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				"",
				"",
				externalIDL,
				fmt.Sprintf("%f", lat),
				fmt.Sprintf("%f", long),
			})
		case externalIDM:
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				parentLocation,
				"",
				externalIDM,
				fmt.Sprintf("%f", lat),
				fmt.Sprintf("%f", long),
			})
		default:
			if ln[1] == locTypeNameL {
				continue
			} else {
				require.Fail(t, "unknown line %s", ln)
			}
		}
	}
	require.Equal(t, 3, linesCount)
}

func TestExportLocationWithPropertyFilters(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log
	ctx := viewertest.NewContext(context.Background(), r.client)
	e := &exporter{log, locationsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	prepareData(ctx, t, *r)

	req, err := http.NewRequest("GET", server.URL, nil)
	require.NoError(t, err)
	req.Header.Set(tenantHeader, "fb-test")

	f, err := json.Marshal([]locationsFilterInput{

		{
			Name:     "PROPERTY",
			Operator: "IS",
			PropertyValue: models.PropertyTypeInput{
				Name:        propNameStr,
				Type:        "string",
				StringValue: pointer.ToString("override"),
			},
		},
	})
	require.NoError(t, err)
	q := req.URL.Query()
	q.Add("filters", string(f))
	req.URL.RawQuery = q.Encode()

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	defer res.Body.Close()

	reader := csv.NewReader(res.Body)
	linesCount := 0
	for {
		ln, err := reader.Read()
		if err == io.EOF {
			break
		}
		linesCount++
		require.NoError(t, err, "error reading row")
		if ln[3] == childLocation {
			require.EqualValues(t, ln[1:], []string{
				grandParentLocation,
				parentLocation,
				childLocation,
				"",
				fmt.Sprintf("%f", 0.0),
				fmt.Sprintf("%f", 0.0),
				"override",
				"true",
				"1988-03-29",
			})
		}
	}
	require.Equal(t, 2, linesCount)
}
