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

	"github.com/facebookincubator/symphony/graph/importer"
	"github.com/stretchr/testify/require"
)

func writeModifiedLinksCSV(t *testing.T, r *csv.Reader) (*bytes.Buffer, string) {
	var newLine []string
	var lines = make([][]string, 3)
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)

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
			if line[2] == portName1 {
				newLine[14] = "new-prop-value"
				newLine[15] = "true"
				newLine[16] = "10"
			}
			lines[i] = newLine
		}
	}

	for _, l := range lines {
		stringLine := strings.Join(l, ",")
		fileWriter.Write([]byte(stringLine + "\n"))
	}
	ct := bw.FormDataContentType()
	require.NoError(t, bw.Close())
	return &buf, ct
}

func TestImportAndEditLinks(t *testing.T) {
	r, err := newExporterTestResolver(t)
	require.NoError(t, err)
	log := r.exporter.log
	e := &exporter{log, linksRower{log}}
	ctx, res := prepareLinksPortsAndExport(t, r, e)
	defer res.Body.Close()
	importLinksPortsFile(t, r.client, res.Body, importer.ImportEntityLink)

	locs := r.client.Location.Query().AllX(ctx)
	require.Len(t, locs, 3)
	links, err := r.Query().LinkSearch(ctx, nil, nil)
	require.NoError(t, err)
	require.Equal(t, 1, links.Count)
	for _, link := range links.Links {
		props := link.QueryProperties().AllX(ctx)
		for _, prop := range props {
			switch prop.QueryType().OnlyX(ctx).Name {
			case propNameInt:
				require.Equal(t, prop.IntVal, 10)
			case propNameBool:
				require.Equal(t, prop.BoolVal, true)
			case propNameStr:
				require.Equal(t, prop.StringVal, "new-prop-value")
			}
		}
	}
}
