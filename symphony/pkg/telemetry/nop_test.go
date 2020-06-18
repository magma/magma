// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"testing"

	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/stretchr/testify/assert"
)

func TestNopExporter(t *testing.T) {
	_, err := telemetry.GetTraceExporter("nop",
		telemetry.TraceExporterOptions{},
	)
	assert.NoError(t, err)
	_, err = telemetry.GetViewExporter("nop",
		telemetry.ViewExporterOptions{},
	)
	assert.NoError(t, err)
}
