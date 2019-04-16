/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"testing"

	"magma/orc8r/cloud/go/services/metricsd/graphite/exporters"

	"github.com/stretchr/testify/assert"
)

func TestValidateQuery(t *testing.T) {
	goodQueries := []string{"metric_name", "metric.name", "metric_name,tag1=val1,tag2=val2"}
	for _, query := range goodQueries {
		assert.True(t, validateQuery(query))
	}
	badQueries := []string{"bad$name", "function(metric_name)", "metric;badtags", "metric;delimiter=bad"}
	for _, query := range badQueries {
		assert.False(t, validateQuery(query))
	}
}

func TestParseTagsFromQuery(t *testing.T) {
	query := "metric,tag1=val1"
	tags := parseTagsFromQuery(query)
	assert.Equal(t, 1, len(tags))
	assert.Equal(t, tags.String(), ";tag1=val1")

	query = "metric,tag1=val1,tag2=val2"
	tags = parseTagsFromQuery(query)
	assert.Equal(t, 2, len(tags))
	assert.Equal(t, tags.String(), ";tag1=val1;tag2=val2")

	query = "metric_name"
	tags = parseTagsFromQuery(query)
	assert.Equal(t, exporters.TagSet{}, tags)
}

func TestBuildSeriesByTagQuery(t *testing.T) {
	tags := make(exporters.TagSet)

	name := "metric_name"
	query := buildTaggedQuery(name, tags)
	assert.Equal(t, "seriesByTag('name=~^metric_name$')", query)

	tags.Insert("tag1", "val1")
	query = buildTaggedQuery(name, tags)
	assert.Equal(t, "seriesByTag('name=~^metric_name$','tag1=val1')", query)

	tags.Insert("tag2", "val2")
	query = buildTaggedQuery(name, tags)
	assert.Equal(t, "seriesByTag('name=~^metric_name$','tag1=val1','tag2=val2')", query)
}
