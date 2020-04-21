// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package importer

import (
	"bytes"
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipmentposition"
	"github.com/facebookincubator/symphony/graph/ent/equipmentpositiondefinition"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/graphql/models"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"

	"github.com/stretchr/testify/require"
)

var locStruct map[string]*ent.Location

func importEquipmentExportedData(ctx context.Context, t *testing.T, r *TestImporterResolver) int {
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)
	err := bw.WriteField("skip_lines", "[5,6]")
	require.NoError(t, err)
	file, err := os.Open("testdata/exportedEquipmentData.csv")
	require.NoError(t, err)
	fileWriter, err := bw.CreateFormFile("file_0", file.Name())
	require.NoError(t, err)
	_, err = io.Copy(fileWriter, file)
	require.NoError(t, err)
	contentType := bw.FormDataContentType()
	require.NoError(t, bw.Close())

	th := viewer.TenancyHandler(http.HandlerFunc(r.importer.processExportedEquipment), viewer.NewFixedTenancy(r.client))
	h := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		th.ServeHTTP(w, r.WithContext(ctx))
	})
	server := httptest.NewServer(h)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL, &buf)
	require.Nil(t, err)

	viewertest.SetDefaultViewerHeaders(req)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	code := resp.StatusCode
	resp.Body.Close()
	return code
}

func createLocationTypes(ctx context.Context, t *testing.T, r *TestImporterResolver) {
	cnt := r.client.LocationType.Create().SetName("Country").SaveX(ctx)
	cty := r.client.LocationType.Create().SetName("City").SaveX(ctx)
	bld := r.client.LocationType.Create().SetName("Building").SaveX(ctx)
	_, err := r.importer.r.Mutation().EditLocationTypesIndex(ctx, []*models.LocationTypeIndex{
		{
			LocationTypeID: cnt.ID,
			Index:          0,
		},
		{
			LocationTypeID: cty.ID,
			Index:          1,
		},
		{
			LocationTypeID: bld.ID,
			Index:          2,
		},
	})
	require.NoError(t, err)
}

func createEquipmentTypes(ctx context.Context, r *TestImporterResolver) {
	proptypes := []*models.PropertyTypeInput{{
		Name: "prop1Str",
		Type: "string",
	}, {
		Name: "prop1Int",
		Type: "int",
	},
	}
	postypes := []*models.EquipmentPositionInput{{
		Name: "pos1",
	}}
	m := r.importer.r.Mutation()
	_, _ = m.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "EquipType1",
		Properties: proptypes,
	})
	_, _ = m.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:       "EquipType2",
		Properties: proptypes,
	})
	_, _ = m.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name:      "EquipType3",
		Positions: postypes,
	})
	_, _ = m.AddEquipmentType(ctx, models.AddEquipmentTypeInput{
		Name: "EquipType4",
	})
}

func verifyLocationsStructure(ctx context.Context, t *testing.T, r TestImporterResolver) {
	locs, err := r.importer.r.Query().Locations(ctx, nil, nil, nil, nil, nil, nil, nil, nil)
	require.NoError(t, err)
	require.Len(t, locs.Edges, 6)
	client := r.client

	locStruct = map[string]*ent.Location{
		"ISR":  client.Location.Query().Where(location.Name("Israel")).OnlyX(ctx),
		"JERU": client.Location.Query().Where(location.Name("Jerusalem")).OnlyX(ctx),
		"TLV":  client.Location.Query().Where(location.Name("TLV")).OnlyX(ctx),
		"b1":   client.Location.Query().Where(location.Name("b1")).OnlyX(ctx),
		"b2":   client.Location.Query().Where(location.Name("b2")).OnlyX(ctx),
		"b3":   client.Location.Query().Where(location.Name("b3")).OnlyX(ctx),
	}
	require.Equal(t, "Country", locStruct["ISR"].QueryType().OnlyX(ctx).Name)
	p, err := locStruct["ISR"].QueryParent().Only(ctx)
	require.Error(t, err)
	require.Nil(t, p)

	require.Equal(t, "City", locStruct["TLV"].QueryType().OnlyX(ctx).Name)
	require.Equal(t, locStruct["ISR"].ID, locStruct["TLV"].QueryParent().OnlyX(ctx).ID)

	require.Equal(t, "City", locStruct["JERU"].QueryType().OnlyX(ctx).Name)
	require.Equal(t, locStruct["ISR"].ID, locStruct["JERU"].QueryParent().OnlyX(ctx).ID)

	require.Equal(t, "Building", locStruct["b1"].QueryType().OnlyX(ctx).Name)
	require.Equal(t, locStruct["TLV"].ID, locStruct["b1"].QueryParent().OnlyX(ctx).ID)

	require.Equal(t, "Building", locStruct["b2"].QueryType().OnlyX(ctx).Name)
	require.Equal(t, locStruct["ISR"].ID, locStruct["b2"].QueryParent().OnlyX(ctx).ID)

	require.Equal(t, "Building", locStruct["b3"].QueryType().OnlyX(ctx).Name)
	require.Equal(t, locStruct["JERU"].ID, locStruct["b3"].QueryParent().OnlyX(ctx).ID)
}

func TestEquipmentImportData(t *testing.T) {
	r := newImporterTestResolver(t)
	ctx := newImportContext(viewertest.NewContext(context.Background(), r.client))
	code := importEquipmentExportedData(ctx, t, r)
	require.Equal(t, http.StatusBadRequest, code)
	q := r.importer.r.Query()

	createLocationTypes(ctx, t, r)
	createEquipmentTypes(ctx, r)
	code = importEquipmentExportedData(ctx, t, r)
	require.Equal(t, 200, code)

	verifyLocationsStructure(ctx, t, *r)
	equips, err := q.EquipmentSearch(ctx, nil, nil)
	require.NoError(t, err)
	require.Equal(t, 4, equips.Count)
	for _, equip := range equips.Equipment {
		switch equip.Name {
		case "A":
			require.Equal(t, "EquipType1", equip.QueryType().OnlyX(ctx).Name)
			require.Equal(t, "AA", equip.ExternalID)
			require.Equal(t, locStruct["b1"].ID, equip.QueryLocation().OnlyXID(ctx))
			require.Equal(t, "val1", equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name("prop1Str"))).OnlyX(ctx).StringVal)
			require.Equal(t, 12, equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name("prop1Int"))).OnlyX(ctx).IntVal)
		case "B":
			require.Equal(t, "EquipType2", equip.QueryType().OnlyX(ctx).Name)
			require.Equal(t, "BB", equip.ExternalID)
			require.Equal(t, locStruct["b2"].ID, equip.QueryLocation().OnlyXID(ctx))
			require.Equal(t, "val2", equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name("prop1Str"))).OnlyX(ctx).StringVal)
			require.Equal(t, 13, equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name("prop1Int"))).OnlyX(ctx).IntVal)
		case "C":
			require.Equal(t, "EquipType3", equip.QueryType().OnlyX(ctx).Name)
			require.Equal(t, "CC", equip.ExternalID)
			require.Equal(t, locStruct["b3"].ID, equip.QueryLocation().OnlyXID(ctx))
			require.True(t, equip.QueryPositions().Where(equipmentposition.HasDefinitionWith(equipmentpositiondefinition.Name("pos1"))).ExistX(ctx))
		case "D":
			require.Equal(t, "EquipType4", equip.QueryType().OnlyX(ctx).Name)
			require.Equal(t, "", equip.ExternalID)
			loc, err := equip.QueryLocation().Only(ctx)
			require.Error(t, err)
			require.Nil(t, loc)
			require.Equal(t, "pos1", equip.QueryParentPosition().QueryDefinition().OnlyX(ctx).Name)
			require.Equal(t, "C", equip.QueryParentPosition().QueryParent().OnlyX(ctx).Name)
		}
	}
}
