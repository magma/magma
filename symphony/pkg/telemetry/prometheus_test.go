// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewPrometheusExporter(t *testing.T) {
	exporter, err := NewPrometheusExporter(ViewExporterOptions{})
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}
