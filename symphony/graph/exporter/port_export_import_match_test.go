// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package exporter

import (
	"bytes"
	"encoding/csv"
	"io"
	"mime/multipart"
	"strings"
	"testing"

	"github.com/facebookincubator/symphony/graph/ent/property"
	"github.com/facebookincubator/symphony/graph/ent/propertytype"
	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/stretchr/testify/require"
)

func writeModifiedPortsCSV(t *testing.T, r *csv.Reader, skipLines, withVerify bool) (*bytes.Buffer, string) {
	var newLine []string
	var lines = make([][]string, 3)
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)
	if skipLines {
		_ = bw.WriteField("skip_lines", "[2]")
	}
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
			newLine = line
			if line[1] == portName1 {
				newLine[17] = "new-prop-value"
				newLine[18] = "new-prop-value2"
			} else if line[1] == portName2 {
				newLine[16] = ""
			}
			lines[i] = newLine
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
		fileWriter.Write([]byte(stringLine + "\n"))
	}
	ct := bw.FormDataContentType()
	require.NoError(t, bw.Close())
	return &buf, ct
}

func TestImportAndEditPorts(t *testing.T) {
	for _, withVerify := range []bool{true, false} {
		for _, skipLines := range []bool{true, false} {
			r := newExporterTestResolver(t)
			log := r.exporter.log
			e := &exporter{log, portsRower{log}}
			ctx, res := prepareHandlerAndExport(t, r, e)

			importLinksPortsFile(t, r.client, res.Body, importer.ImportEntityPort, MethodEdit, skipLines, withVerify)
			res.Body.Close()
			locs := r.client.Location.Query().AllX(ctx)
			require.Len(t, locs, 3)
			ports, err := r.Query().PortSearch(ctx, nil, nil)
			require.NoError(t, err)
			require.Equal(t, 2, ports.Count)
			for _, port := range ports.Ports {
				def := port.QueryDefinition().OnlyX(ctx)
				if def.Name == portName1 {
					typ := def.QueryEquipmentPortType().OnlyX(ctx)
					propTyps := typ.QueryPropertyTypes().AllX(ctx)
					require.Len(t, propTyps, 2)

					props := port.QueryProperties().AllX(ctx)
					if withVerify {
						require.Empty(t, props)
					} else {
						require.Len(t, props, 2)

						p1 := typ.QueryPropertyTypes().Where(propertytype.Name(propStr)).OnlyX(ctx)
						p2 := typ.QueryPropertyTypes().Where(propertytype.Name(propStr2)).OnlyX(ctx)

						require.Equal(t, port.QueryProperties().Where(property.HasTypeWith(propertytype.ID(p1.ID))).OnlyX(ctx).StringVal, "new-prop-value")
						require.Equal(t, port.QueryProperties().Where(property.HasTypeWith(propertytype.ID(p2.ID))).OnlyX(ctx).StringVal, "new-prop-value2")
					}
				}
			}
		}
	}
}
