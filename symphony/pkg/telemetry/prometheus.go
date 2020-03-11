// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"contrib.go.opencensus.io/exporter/prometheus"
	"go.opencensus.io/stats/view"
	"go.uber.org/zap"
)

func init() {
	MustRegisterViewExporter("prometheus", NewPrometheusExporter)
}

// NewJaegerExporter creates a new opencensus view exporter.
func NewPrometheusExporter(opts ViewExporterOptions) (view.Exporter, error) {
	o := prometheus.Options{
		ConstLabels: opts.Labels,
		OnError: func(err error) {
			zap.L().Error("cannot export view data to prometheus", zap.Error(err))
		},
	}
	exporter, err := prometheus.NewExporter(o)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
