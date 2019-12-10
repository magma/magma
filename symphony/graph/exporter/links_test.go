// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"context"
	"encoding/csv"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"

	"github.com/AlekSi/pointer"
	"github.com/facebookincubator/symphony/graph/ent/equipmentportdefinition"
	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/stretchr/testify/require"
)

const (
	portAIDTitle        = "Port A ID"
	portANameTitle      = "Port A Name"
	portATypeTitle      = "Port A Type"
	portBIDTitle        = "Port B ID"
	portBNameTitle      = "Port B Name"
	portBTypeTitle      = "Port B Type"
	equipmentAIDTitle   = "Equipment A ID"
	equipmentANameTitle = "Equipment A Name"
	equipmentATypeTitle = "Equipment A Type"
	equipmentBIDTitle   = "Equipment B ID"
	equipmentBNameTitle = "Equipment B Name"
	equipmentBTypeTitle = "Equipment B Type"
)

func prepareLinkData(ctx context.Context, t *testing.T, r TestExporterResolver) {
	mr := r.Mutation()

	locTypeL, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameL})
	require.NoError(t, err)
	locTypeM, err := mr.AddLocationType(ctx, models.AddLocationTypeInput{Name: locTypeNameM})
	require.NoError(t, err)

	_, err = mr.EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: locTypeL.ID,
			Index:          0,
		},
		{
			LocationTypeID: locTypeM.ID,
			Index:          1,
		},
	})
	require.NoError(t, err)

	gpLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name: grandParentLocation,
		Type: locTypeL.ID,
	})
	require.NoError(t, err)
	pLocation, err := mr.AddLocation(ctx, models.AddLocationInput{
		Name:   parentLocation,
		Type:   locTypeM.ID,
		Parent: &gpLocation.ID,
	})
	require.NoError(t, err)
	strDefVal := propDefValue
	intDefVal := propDevValInt
	propDefInput1 := models.PropertyTypeInput{
		Name:        propNameStr,
		Type:        "string",
		StringValue: &strDefVal,
	}
	propDefInput2 := models.PropertyTypeInput{
		Name:     propNameInt,
		Type:     "int",
		IntValue: &intDefVal,
	}

	ptyp, _ := mr.AddEquipmentPortType(ctx, models.AddEquipmentPortTypeInput{
		Name: "portType1",
		LinkProperties: []*models.PropertyTypeInput{
			{
				Name:        propStr,
				Type:        "string",
				StringValue: pointer.ToString("t1"),
			},
			{
				Name: propStr2,
				Type: "string",
			},
		},
	})
	port1 := models.EquipmentPortInput{
		Name:       portName1,
		PortTypeID: &ptyp.ID,
	}

	equipmentType, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:  equipmentTypeName,
		Ports: []*models.EquipmentPortInput{&port1},
	})
	require.NoError(t, err)

	port2 := models.EquipmentPortInput{
		Name: portName2,
	}
	port3 := models.EquipmentPortInput{
		Name: portName3,
	}
	position1 := models.EquipmentPositionInput{
		Name: positionName,
	}
	equipmentType2, err := mr.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       equipmentType2Name,
		Properties: []*models.PropertyTypeInput{&propDefInput1, &propDefInput2},
		Ports:      []*models.EquipmentPortInput{&port2, &port3},
		Positions:  []*models.EquipmentPositionInput{&position1},
	})
	require.NoError(t, err)

	parentEquipment, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     parentEquip,
		Type:     equipmentType.ID,
		Location: &pLocation.ID,
	})
	require.NoError(t, err)

	posDef1 := equipmentType2.QueryPositionDefinitions().Where(equipmentpositiondefinition.Name(positionName)).OnlyX(ctx)

	childEquip1, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:               currEquip,
		Type:               equipmentType2.ID,
		Parent:             &parentEquipment.ID,
		PositionDefinition: &posDef1.ID,
	})
	require.NoError(t, err)

	childEquip2, err := mr.AddEquipment(ctx, models.AddEquipmentInput{
		Name:     currEquip2,
		Type:     equipmentType2.ID,
		Location: &pLocation.ID,
	})
	require.NoError(t, err)

	portDef1 := equipmentType.QueryPortDefinitions().Where(equipmentportdefinition.Name(portName1)).OnlyX(ctx)
	portDef2 := equipmentType2.QueryPortDefinitions().Where(equipmentportdefinition.Name(portName2)).OnlyX(ctx)
	portDef3 := equipmentType2.QueryPortDefinitions().Where(equipmentportdefinition.Name(portName3)).OnlyX(ctx)

	_, _ = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: parentEquipment.ID, Port: portDef1.ID},
			{Equipment: childEquip1.ID, Port: portDef2.ID},
		},
	})

	_, _ = mr.AddLink(ctx, models.AddLinkInput{
		Sides: []*models.LinkSide{
			{Equipment: childEquip1.ID, Port: portDef3.ID},
			{Equipment: childEquip2.ID, Port: portDef3.ID},
		},
	})
	require.NoError(t, err)
}

func TestEmptyLinksDataExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, linksRower{log}}
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
			"\ufeffLink ID",
			portAIDTitle,
			portANameTitle,
			portATypeTitle,
			equipmentAIDTitle,
			equipmentANameTitle,
			equipmentATypeTitle,
			portBIDTitle,
			portBNameTitle,
			portBTypeTitle,
			equipmentBIDTitle,
			equipmentBNameTitle,
			equipmentBTypeTitle,
		}, ln)
	}
}

func TestLinksExport(t *testing.T) {
	r, err := newExporterTestResolver(t)
	log := r.exporter.log
	require.NoError(t, err)

	e := &exporter{log, linksRower{log}}
	th := viewer.TenancyHandler(e, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest("GET", server.URL, nil)
	req.Header.Set(tenantHeader, "fb-test")

	ctx := viewertest.NewContext(r.client)
	prepareLinkData(ctx, t, *r)
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
		case ln[1] == portAIDTitle:
			require.EqualValues(t, []string{
				"\ufeffLink ID",
				portAIDTitle,
				portANameTitle,
				portATypeTitle,
				equipmentAIDTitle,
				equipmentANameTitle,
				equipmentATypeTitle,
				portBIDTitle,
				portBNameTitle,
				portBTypeTitle,
				equipmentBIDTitle,
				equipmentBNameTitle,
				equipmentBTypeTitle,
				propStr,
				propStr2,
			}, ln)
		case ln[2] == portName1:
			ln[4] = "--"
			ln[7] = "--"
			ln[10] = "--"
			require.EqualValues(t, ln[2:], []string{
				portName1,
				"portType1",
				"--",
				parentEquip,
				equipmentTypeName,
				"--",
				portName2,
				"",
				"--",
				currEquip,
				equipmentType2Name,
				"t1",
				"",
			})
		case ln[2] == portName3:
			ln[4] = "--"
			ln[7] = "--"
			ln[10] = "--"
			require.EqualValues(t, ln[2:], []string{
				portName3,
				"",
				"--",
				currEquip,
				equipmentType2Name,
				"--",
				portName3,
				"",
				"--",
				currEquip2,
				equipmentType2Name,
				"",
				"",
			})
		default:
			require.Fail(t, "line does not match", ln)
		}
	}
}
