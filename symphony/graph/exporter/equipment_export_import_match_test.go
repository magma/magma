// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"bytes"
	"context"
	"encoding/csv"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/facebookincubator/symphony/graph/authz"
	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/ent/equipment"
	"github.com/facebookincubator/symphony/graph/ent/location"
	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/graph/viewer/viewertest"
	"github.com/facebookincubator/symphony/pkg/log/logtest"

	"github.com/stretchr/testify/require"
)

// TODO (T59270743): Move this file to importer folder and refactor similar code with exported_service_integration_test.go
func writeModifiedCSV(t *testing.T, r *csv.Reader, method method, withVerify bool) (*bytes.Buffer, string) {
	var newLine []string
	var lines = make([][]string, 3)
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)

	if withVerify {
		_ = bw.WriteField("verify_before_commit", "true")
	}
	fileWriter, err := bw.CreateFormFile("file_0", "name1")
	require.Nil(t, err)
	for i := 0; ; i++ {
		line, err := r.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			require.Nil(t, err)
		}
		if i == 0 {
			lines[0] = line
		} else {
			switch method {
			case MethodAdd:
				newLine = append([]string{""}, line[1:]...)
			case MethodEdit:
				newLine = line
				if line[1] == currEquip {
					newLine[13] = "str-prop-value" + strconv.FormatInt(int64(i), 10)
					newLine[1] = "newName" + strconv.FormatInt(int64(i), 10)
					newLine[14] = "10" + strconv.FormatInt(int64(i), 10)
					newLine[15] = "new-prop-value" + strconv.FormatInt(int64(i), 10)
				}
			default:
				require.Fail(t, "method should be add or edit")
			}
			// parent row must exist before child row
			if line[1] == parentEquip {
				lines[1] = newLine
			} else {
				lines[2] = newLine
			}
		}
	}

	if withVerify {
		failLine := make([]string, len(lines[1]))
		copy(failLine, lines[1])
		lines = append(lines, failLine)
		lines[3][1] = "this"
		lines[3][2] = "should"
		lines[3][3] = "fail"
	}

	for _, l := range lines {
		stringLine := strings.Join(l, ",")
		_, _ = io.WriteString(fileWriter, stringLine+"\n")
	}
	ct := bw.FormDataContentType()
	require.NoError(t, bw.Close())
	return &buf, ct
}

func importEquipmentFile(t *testing.T, client *ent.Client, r io.Reader, method method, withVerify bool) {
	readr := csv.NewReader(r)
	buf, contentType := writeModifiedCSV(t, readr, method, withVerify)

	h, _ := importer.NewHandler(
		importer.Config{
			Logger:     logtest.NewTestLogger(t),
			Subscriber: event.NewNopSubscriber(),
		},
	)
	auth := authz.AuthHandler{Handler: h, Logger: logtest.NewTestLogger(t)}
	th := viewer.TenancyHandler(auth, viewer.NewFixedTenancy(client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/export_equipment", buf)
	require.Nil(t, err)

	viewertest.SetDefaultViewerHeaders(req)
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	resp.Body.Close()
}

func deleteEquipmentData(ctx context.Context, t *testing.T, r *TestExporterResolver) {
	id := r.client.Equipment.Query().Where(equipment.Name(currEquip)).OnlyXID(ctx)
	_, err := r.Mutation().RemoveEquipment(ctx, id, nil)
	require.NoError(t, err)

	id = r.client.Equipment.Query().Where(equipment.Name(parentEquip)).OnlyXID(ctx)
	_, err = r.Mutation().RemoveEquipment(ctx, id, nil)
	require.NoError(t, err)

	id = r.client.Location.Query().Where(location.Name(childLocation)).OnlyXID(ctx)
	_, err = r.Mutation().RemoveLocation(ctx, id)
	require.NoError(t, err)

	id = r.client.Location.Query().Where(location.Name(parentLocation)).OnlyXID(ctx)
	_, err = r.Mutation().RemoveLocation(ctx, id)
	require.NoError(t, err)

	id = r.client.Location.Query().Where(location.Name(grandParentLocation)).OnlyXID(ctx)
	_, err = r.Mutation().RemoveLocation(ctx, id)
	require.NoError(t, err)
}

func prepareEquipmentAndExport(t *testing.T, r *TestExporterResolver) (context.Context, *http.Response) {
	log := r.exporter.log

	e := &exporter{log, equipmentRower{log}}
	auth := authz.AuthHandler{Handler: e, Logger: logtest.NewTestLogger(t)}
	th := viewer.TenancyHandler(auth, viewer.NewFixedTenancy(r.client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodGet, server.URL, nil)
	require.NoError(t, err)
	viewertest.SetDefaultViewerHeaders(req)

	ctx := viewertest.NewContext(context.Background(), r.client)
	prepareData(ctx, t, *r)
	locs := r.client.Location.Query().AllX(ctx)
	require.Len(t, locs, 3)

	res, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	return ctx, res
}

func TestEquipmentExportAndImportMatch(t *testing.T) {
	for _, verify := range []bool{true, false} {
		verify := verify
		t.Run("Verify/"+strconv.FormatBool(verify), func(t *testing.T) {
			r := newExporterTestResolver(t)
			ctx, res := prepareEquipmentAndExport(t, r)
			defer res.Body.Close()
			deleteEquipmentData(ctx, t, r)

			locs := r.client.Location.Query().AllX(ctx)
			require.Len(t, locs, 0)

			importEquipmentFile(t, r.client, res.Body, MethodAdd, verify)
			locs = r.client.Location.Query().AllX(ctx)
			if verify {
				require.Len(t, locs, 0)
			} else {
				require.Len(t, locs, 3)
			}
			for _, loc := range locs {
				switch loc.Name {
				case grandParentLocation:
					require.Equal(t, locTypeNameL, loc.QueryType().OnlyX(ctx).Name)
					require.Equal(t, parentLocation, loc.QueryChildren().OnlyX(ctx).Name)
				case parentLocation:
					require.Equal(t, locTypeNameM, loc.QueryType().OnlyX(ctx).Name)
					require.Equal(t, childLocation, loc.QueryChildren().OnlyX(ctx).Name)
				case childLocation:
					require.Equal(t, locTypeNameS, loc.QueryType().OnlyX(ctx).Name)
					require.Equal(t, parentEquip, loc.QueryEquipment().OnlyX(ctx).Name)
				}
			}
			equips, err := r.Query().EquipmentSearch(ctx, nil, nil)
			require.NoError(t, err)
			if verify {
				require.Equal(t, 0, equips.Count)
			} else {
				require.Equal(t, 2, equips.Count)
			}
			for _, equip := range equips.Equipment {
				switch equip.Name {
				case currEquip:
					require.Equal(t, equipmentType2Name, equip.QueryType().OnlyX(ctx).Name)
					pos := equip.QueryParentPosition().OnlyX(ctx)
					require.Equal(t, positionName, pos.QueryDefinition().OnlyX(ctx).Name)
					require.Equal(t, parentEquip, pos.QueryParent().OnlyX(ctx).Name)
					prop := equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(propNameStr))).OnlyX(ctx)
					require.Equal(t, propInstanceValue, prop.StringVal)

					prop = equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(propNameInt))).OnlyX(ctx)
					require.Equal(t, propDevValInt, prop.IntVal)

					prop = equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(newPropNameStr))).OnlyX(ctx)
					require.Equal(t, propDefValue2, prop.StringVal)

				case parentEquip:
					require.Equal(t, childLocation, equip.QueryLocation().OnlyX(ctx).Name)
					require.Equal(t, equipmentTypeName, equip.QueryType().OnlyX(ctx).Name)
				}
			}
		})
	}
}

func TestEquipmentImportAndEdit(t *testing.T) {
	for _, verify := range []bool{true, false} {
		verify := verify
		t.Run("Verify/"+strconv.FormatBool(verify), func(t *testing.T) {
			r := newExporterTestResolver(t)
			ctx, res := prepareEquipmentAndExport(t, r)
			defer res.Body.Close()

			importEquipmentFile(t, r.client, res.Body, MethodEdit, verify)
			locs := r.client.Location.Query().AllX(ctx)
			require.Len(t, locs, 3)
			equips, err := r.Query().EquipmentSearch(ctx, nil, nil)
			require.NoError(t, err)
			require.Equal(t, 2, equips.Count)
			for _, equip := range equips.Equipment {
				switch equip.Name {
				case parentEquip:
					require.Equal(t, equipmentTypeName, equip.QueryType().OnlyX(ctx).Name)
					pos := equip.QueryPositions().OnlyX(ctx)
					require.Equal(t, positionName, pos.QueryDefinition().OnlyX(ctx).Name)
				case "newName2":
					require.Equal(t, equipmentType2Name, equip.QueryType().OnlyX(ctx).Name)
					pos := equip.QueryParentPosition().OnlyX(ctx)
					require.Equal(t, positionName, pos.QueryDefinition().OnlyX(ctx).Name)
					require.Equal(t, parentEquip, pos.QueryParent().OnlyX(ctx).Name)
					prop := equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(propNameStr))).OnlyX(ctx)
					require.Equal(t, "str-prop-value2", prop.StringVal)

					prop = equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(propNameInt))).OnlyX(ctx)
					require.Equal(t, 102, prop.IntVal)

					prop = equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(newPropNameStr))).OnlyX(ctx)
					require.Equal(t, "new-prop-value2", prop.StringVal)
				case currEquip:
					// reach here on "verify mode" that failed (this is the name with no edit
					require.True(t, verify)
					require.Equal(t, equipmentType2Name, equip.QueryType().OnlyX(ctx).Name)
					prop := equip.QueryProperties().Where(property.HasTypeWith(propertytype.Name(propNameStr))).OnlyX(ctx)
					require.Equal(t, propInstanceValue, prop.StringVal)
				}
			}
		})
	}
}
