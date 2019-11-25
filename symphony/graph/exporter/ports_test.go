// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"encoding/csv"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/stretchr/testify/require"
)

const portNameTitle = "Port Name"
const portTypeTitle = "Port Type"
const equipmentNameTitle = "Equipment Name"
const equipmentTypeTitle = "Equipment Type"
const p3Title = "Parent Equipment (3)"
const p2Title = "Parent Equipment (2)"
const p1Title = "Parent Equipment"
const positionTitle = "Equipment Position"
const linkPID = "Linked Port ID"
const linkPName = "Linked Port name"
const linkEID = "Linked Port Equipment ID"
const linkEName = "Linked Port Equipment"
const propStr = "propStr"
const propStr2 = "propStr2"

func TestEmptyPortsDataExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, portsRower{log}}
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
			"\ufeffPort ID",
			portNameTitle,
			portTypeTitle,
			equipmentNameTitle,
			equipmentTypeTitle,
			p3Title,
			p2Title,
			p1Title,
			positionTitle,
			linkPID,
			linkPName,
			linkEID,
			linkEName,
		}, ln)
	}
}

func TestPortsExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, portsRower{log}}
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
		case ln[1] == portNameTitle:
			require.EqualValues(t, []string{
				"\ufeffPort ID",
				portNameTitle,
				portTypeTitle,
				equipmentNameTitle,
				equipmentTypeTitle,
				"locTypeLarge",
				"locTypeMedium",
				"locTypeSmall",
				p3Title,
				p2Title,
				p1Title,
				positionTitle,
				linkPID,
				linkPName,
				linkEID,
				linkEName,
				propStr,
				propStr2,
			}, ln)
		case ln[1] == portName1:
			ln[12] = "--"
			ln[14] = "--"
			require.EqualValues(t, ln[1:], []string{
				portName1,
				"portType1",
				parentEquip,
				equipmentTypeName,
				grandParentLocation,
				parentLocation,
				childLocation,
				"",
				"",
				"",
				"",
				"--",
				"port2",
				"--",
				currEquip,
				"t1",
				"",
			})
		case ln[1] == portName2:
			ln[12] = "--"
			ln[14] = "--"
			require.EqualValues(t, ln[1:], []string{
				portName2,
				"",
				currEquip,
				equipmentType2Name,
				grandParentLocation,
				parentLocation,
				childLocation,
				"",
				"",
				parentEquip,
				positionName,
				"--",
				portName1,
				"--",
				parentEquip,
				"",
				"",
			})
		default:
			require.Fail(t, "line does not match")
		}
	}
}
