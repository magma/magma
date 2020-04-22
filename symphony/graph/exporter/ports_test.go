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
	"strconv"
	"testing"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/graphql/models"
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
const linkPName = "Linked Port Name"
const linkEID = "Linked Equipment ID"
const linkEName = "Linked Equipment"
const propStr = "propStr"
const propStr2 = "propStr2"
const servicesTitle = "Service Names"

func TestEmptyPortsDataExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, portsRower{log}}
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
			servicesTitle,
		}, ln)
	}
}

func TestPortsExport(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log

	e := &exporter{log, portsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	viewertest.SetDefaultViewerHeaders(req)

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
				servicesTitle,
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
				"S1;S2",
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
				"S1",
				"",
				"",
			})
		default:
			require.Fail(t, "line does not match")
		}
	}
}

func TestPortWithFilters(t *testing.T) {
	r := newExporterTestResolver(t)
	log := r.exporter.log
	ctx := viewertest.NewContext(context.Background(), r.client)
	e := &exporter{log, portsRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	prepareData(ctx, t, *r)
	loc := r.client.Location.Query().Where(location.Name(childLocation)).OnlyX(ctx)
	pDef2 := r.client.EquipmentPortDefinition.Query().Where(equipmentportdefinition.Name(portName2)).OnlyX(ctx)

	f1, err := json.Marshal([]portFilterInput{
		{
			Name:     "LOCATION_INST",
			Operator: "IS_ONE_OF",
			IDSet:    []string{strconv.Itoa(loc.ID)},
		},
		{
			Name:     "PORT_DEF",
			Operator: "IS_ONE_OF",
			IDSet:    []string{strconv.Itoa(pDef2.ID)},
		},
		{
			Name:      "PORT_INST_HAS_LINK",
			Operator:  "IS",
			BoolValue: true,
		},
	})
	require.NoError(t, err)
	f2, err := json.Marshal([]portFilterInput{
		{
			Name:     "PROPERTY",
			Operator: "IS",
			PropertyValue: models.PropertyTypeInput{
				ID:          pointer.ToInt(42),
				Name:        propStr,
				StringValue: pointer.ToString("t1"),
				Type:        "string",
			},
		},
	})
	require.NoError(t, err)

	f3, err := json.Marshal([]portFilterInput{
		{
			Name:      "PORT_INST_HAS_LINK",
			Operator:  "IS",
			BoolValue: false,
		},
	})
	require.NoError(t, err)

	for i, filter := range [][]byte{f1, f2, f3} {
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
				if ln[1] == portName1 {
					ln[12] = "--"
					ln[14] = "--"
					require.EqualValues(t, []string{
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
						"",
					}, ln[1:])
					require.Equal(t, 2, linesCount)
				}
			}
			if i == 1 {
				if ln[1] == portName1 {
					ln[12] = "--"
					ln[14] = "--"
					require.EqualValues(t, []string{
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
						"S1;S2",
						"t1",
						"",
					}, ln[1:])
					require.Equal(t, 2, linesCount)
				}
			}
			if i == 2 {
				require.Equal(t, 1, linesCount)
			}
		}
		err = res.Body.Close()
		require.NoError(t, err)
	}
}
