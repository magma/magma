// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"encoding/csv"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

const serviceNameTitle = "Service Name"
const serviceTypeTitle = "Service Type"
const serviceExternalIDTitle = "Service External ID"
const customerNameTitle = "Customer Name"
const customerExternalIDTitle = "Customer External ID"

func TestEmptyServicesDataExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, servicesRower{log}}
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
			"\ufeffService ID",
			serviceNameTitle,
			serviceTypeTitle,
			serviceExternalIDTitle,
			customerNameTitle,
			customerExternalIDTitle,
		}, ln)
	}
}
