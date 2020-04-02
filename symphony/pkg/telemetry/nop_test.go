// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNopExporter(t *testing.T) {
	_, err := GetTraceExporter("nop", TraceExporterOptions{})
	assert.NoError(t, err)
	_, err = GetViewExporter("nop", ViewExporterOptions{})
	assert.NoError(t, err)
}
