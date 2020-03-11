// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"go.opencensus.io/stats/view"
	"go.opencensus.io/trace"
)

func init() {
	MustRegisterTraceExporter("nop", func(TraceExporterOptions) (trace.Exporter, error) {
		return NopExporter{}, nil
	})
	MustRegisterViewExporter("nop", func(ViewExporterOptions) (view.Exporter, error) {
		return NopExporter{}, nil
	})
}

type NopExporter struct{}

func (NopExporter) ExportSpan(*trace.SpanData) {}
func (NopExporter) ExportView(*view.Data)      {}
