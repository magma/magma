// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"encoding/csv"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

func TestEmptyLocationDataExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

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
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, locationsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	req.Header.Set(tenantHeader, "fb-test")

	ctx := viewertest.NewContext(r.client)
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
		case ln[4] == "":
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
