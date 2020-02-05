// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package tests

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"testing"

	"github.com/facebookincubator/symphony/graph/graphql/models"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

func TestImportLocations(t *testing.T) {
	organization := uuid.New().String()
	c := newClient(t, organization, "user@test.com")

	c.log.Debug("adding location types")
	addLocationTypes(t, c)

	c.log.Debug("importing locations")
	importLocations(t, organization, "ExampleLocation.csv")

	c.log.Debug("loading locations")
	locations, err := c.QueryLocations()
	require.NoError(t, err)
	c.log.Debug("loaded locations",
		zap.Int("count", len(locations.Edges)),
	)

	var casesFound int
	for i := range locations.Edges {
		node := locations.Edges[i].Node
		casesFound++
		c.log.Debug("inspecting node", zap.String("name", string(node.Name)))
		switch node.Name {
		case "Houston":
			assert.Len(t, node.Children, 5)
			assert.EqualValues(t, "Texas", node.Parent.Name)
		case "F1":
			require.NotZero(t, len(node.Properties))
			property := node.Properties[0]
			assert.EqualValues(t, "200 sq ft", property.Value)
			assert.EqualValues(t, "Floor Size", property.Type.Name)
			assert.Empty(t, node.Children)
			assert.EqualValues(t, "2392 S Wayside D", node.Parent.Name)
		case "C001":
			require.NotZero(t, len(node.Properties))
			property := node.Properties[0]
			assert.EqualValues(t, "Room Owner", property.Type.Name)
			assert.EqualValues(t, "Elaine", property.Value)
			assert.EqualValues(t, "F2", node.Parent.Name)
		default:
			casesFound--
		}
	}
	assert.Equal(t, 3, casesFound)
}

func TestImportLocationsEdit(t *testing.T) {
	organization := uuid.New().String()
	c := newClient(t, organization, "user@test.com")

	c.log.Debug("adding location types")
	addLocationTypes(t, c)

	c.log.Debug("importing locations[1]")
	importLocations(t, organization, "ExampleLocation.csv")
	importLocations(t, organization, "EditLocation.csv")

	c.log.Debug("loading locations")
	locations, err := c.QueryLocations()
	require.NoError(t, err)
	c.log.Debug("loaded locations",
		zap.Int("count", len(locations.Edges)),
	)

	var casesFound int
	for i := range locations.Edges {
		node := locations.Edges[i].Node
		casesFound++
		c.log.Debug("inspecting node", zap.String("name", string(node.Name)))
		switch node.Name {
		case "2391 S Wayside D":
			assert.EqualValues(t, "id1", node.ExternalID)
			assert.EqualValues(t, 34, node.Latitude)
			assert.EqualValues(t, 35, node.Longitude)
		case "F1":
			require.NotZero(t, len(node.Properties))
			property := node.Properties[0]
			assert.EqualValues(t, "300 sq ft", property.Value)
			assert.EqualValues(t, "id2", node.ExternalID)
			assert.EqualValues(t, 66, node.Latitude)
			assert.EqualValues(t, 67, node.Longitude)
			assert.Empty(t, node.Children)
		default:
			casesFound--
		}
	}
	assert.Equal(t, 2, casesFound)
}

func importLocations(t *testing.T, organization, filename string) {
	var buf bytes.Buffer
	bw := multipart.NewWriter(&buf)

	file, err := os.Open("../../graph/importer/testdata/" + filename)
	require.Nil(t, err)

	fileWriter, err := bw.CreateFormFile("file_0", file.Name())
	require.Nil(t, err)

	_, err = io.Copy(fileWriter, file)
	require.Nil(t, err)

	contentType := bw.FormDataContentType()
	require.NoError(t, bw.Close())

	req, err := http.NewRequest(http.MethodPost, "http://graph/import/location", &buf)
	require.NoError(t, err)

	req.Header.Set("x-auth-organization", organization)
	req.Header.Set("Content-Type", contentType)

	rsp, err := http.DefaultClient.Do(req)
	require.NoError(t, err)
	rsp.Body.Close()
}

func addLocationTypes(t *testing.T, c *client) {
	types := []struct {
		name, property string
	}{
		{name: "Country"},
		{name: "State"},
		{name: "City"},
		{name: "Building"},
		{name: "Floor", property: "Floor Size"},
		{name: "Room", property: "Room Owner"},
	}
	for _, typ := range types {
		if typ.property == "" {
			_, err := c.addLocationType(typ.name)
			assert.NoError(t, err)
		} else {
			_, err := c.addLocationType(typ.name, &models.PropertyTypeInput{
				Name: typ.property,
				Type: models.PropertyKindString,
			})
			assert.NoError(t, err)
		}
	}
}
