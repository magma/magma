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
	"strings"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent"
	"github.com/facebookincubator/symphony/graph/event"
	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/facebookincubator/symphony/graph/viewer"
	"github.com/facebookincubator/symphony/pkg/log/logtest"
	"github.com/stretchr/testify/require"
)

func writeModifiedLocationsCSV(t *testing.T, r *csv.Reader, method method, withVerify, skipLines bool) (*bytes.Buffer, string) {
	var newLine []string
	var lines = make([][]string, 4)
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)
	if skipLines {
		err := bw.WriteField("skip_lines", "[2,3]")
		require.NoError(t, err)
	}
	if withVerify {
		err := bw.WriteField("verify_before_commit", "true")
		require.NoError(t, err)
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
				if line[3] == childLocation {
					newLine[4] = "new-external-id"
					newLine[5] = "44"
					newLine[6] = "55"
					newLine[7] = "new-str"
					newLine[8] = "false"
					newLine[9] = "1988-01-01"
				}
			default:
				require.Fail(t, "method should be add or edit")
			}
			lines[i] = newLine
		}
	}
	if withVerify {
		lines[2][0] = "this"
		lines[3][5] = "should"
		lines[3][6] = "fail"
	}
	for _, l := range lines {
		stringLine := strings.Join(l, ",")
		fileWriter.Write([]byte(stringLine + "\n"))
	}
	ct := bw.FormDataContentType()
	require.NoError(t, bw.Close())
	return &buf, ct
}

func importLocationsFile(t *testing.T, client *ent.Client, r io.Reader, method method, withVerify, skipLines bool) {
	readr := csv.NewReader(r)
	buf, contentType := writeModifiedLocationsCSV(t, readr, method, withVerify, skipLines)

	h, _ := importer.NewHandler(
		importer.Config{
			Logger:     logtest.NewTestLogger(t),
			Emitter:    event.NewNopEmitter(),
			Subscriber: event.NewNopSubscriber(),
		},
	)
	th := viewer.TenancyHandler(h, viewer.NewFixedTenancy(client))
	server := httptest.NewServer(th)
	defer server.Close()

	req, err := http.NewRequest(http.MethodPost, server.URL+"/export_locations", buf)
	require.Nil(t, err)

	req.Header.Set(tenantHeader, "fb-test")
	req.Header.Set("Content-Type", contentType)

	resp, err := http.DefaultClient.Do(req)
	require.Nil(t, err)
	require.Equal(t, resp.StatusCode, http.StatusOK)
	resp.Body.Close()
}

func TestExportAndEditLocations(t *testing.T) {
	for _, withVerify := range []bool{true, false} {
		for _, skipLines := range []bool{true, false} {
			r := newExporterTestResolver(t)
			log := r.exporter.log
			e := &exporter{log, locationsRower{log}}
			ctx, res := prepareHandlerAndExport(t, r, e)
			importLocationsFile(t, r.client, res.Body, MethodEdit, withVerify, skipLines)
			res.Body.Close()

			locations, err := r.Query().LocationSearch(ctx, nil, nil)
			require.NoError(t, err)
			switch {
			case skipLines || withVerify:
				require.Equal(t, 3, locations.Count)
				require.Len(t, locations.Locations, 3)
				for _, loc := range locations.Locations {
					if loc.Name == childLocation {
						require.Empty(t, loc.ExternalID)
						require.Zero(t, loc.Latitude)
						require.Zero(t, loc.Longitude)
					}
				}
			default:
				require.Equal(t, 3, locations.Count)
				for _, loc := range locations.Locations {
					props := loc.QueryProperties().AllX(ctx)
					if loc.Name == childLocation {
						require.Equal(t, "new-external-id", loc.ExternalID)
						require.Equal(t, 44.0, loc.Latitude)
						require.Equal(t, 55.0, loc.Longitude)
						for _, prop := range props {
							switch prop.QueryType().OnlyX(ctx).Name {
							case propNameDate:
								require.Equal(t, "1988-01-01", prop.StringVal)
							case propNameBool:
								require.Equal(t, false, prop.BoolVal)
							case propNameStr:
								require.Equal(t, "new-str", prop.StringVal)
							}
						}
					}
				}
			}
		}
	}
}

func TestExportAndAddLocations(t *testing.T) {
	for _, withVerify := range []bool{true, false} {
		for _, skipLines := range []bool{true, false} {
			r := newExporterTestResolver(t)
			log := r.exporter.log
			e := &exporter{log, locationsRower{log}}
			ctx, res := prepareHandlerAndExport(t, r, e)

			locs := r.client.Location.Query().AllX(ctx)
			require.Len(t, locs, 3)
			// Deleting link and of side's equipment to verify it creates it on import
			deleteLocationsForReImport(ctx, t, r)
			locs = r.client.Location.Query().AllX(ctx)
			require.Len(t, locs, 0)

			importLocationsFile(t, r.client, res.Body, MethodAdd, withVerify, skipLines)
			res.Body.Close()

			locations, err := r.Query().LocationSearch(ctx, nil, nil)
			require.NoError(t, err)
			switch {
			case !skipLines && withVerify:
				require.Zero(t, 0, locations.Count)
				require.Empty(t, locations.Locations)
			case skipLines || withVerify:
				require.Equal(t, 1, locations.Count)
				require.Len(t, locations.Locations, 1)
			default:
				require.Equal(t, 3, locations.Count)
				for _, loc := range locations.Locations {
					props := loc.QueryProperties().AllX(ctx)
					for _, prop := range props {
						switch prop.QueryType().OnlyX(ctx).Name {
						case propNameDate:
							require.Equal(t, "1988-03-29", prop.StringVal)
						case propNameBool:
							require.Equal(t, true, prop.BoolVal)
						case propNameStr:
							require.Equal(t, "override", prop.StringVal)
						}
					}
				}
			}
		}
	}
}

func deleteLocationsForReImport(ctx context.Context, t *testing.T, r *TestExporterResolver) {
	locs := r.client.Location.Query().AllX(ctx)
	for _, loc := range locs {
		err := r.client.Location.DeleteOne(loc).Exec(ctx)
		require.NoError(t, err)
	}
}
