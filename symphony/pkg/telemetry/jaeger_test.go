// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"os"
	"testing"

	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewJaegerExporter(t *testing.T) {
	err := os.Setenv("JAEGER_AGENT_ENDPOINT", "localhost:6831")
	require.NoError(t, err)
	defer os.Unsetenv("JAEGER_AGENT_ENDPOINT")
	exporter, err := telemetry.GetTraceExporter("jaeger",
		telemetry.TraceExporterOptions{ServiceName: t.Name()},
	)
	assert.NoError(t, err)
	assert.NotNil(t, exporter)
}
