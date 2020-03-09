// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"fmt"
	"sort"
	"sync"

	"go.opencensus.io/trace"
)

// TraceExporterOptions defines a set of options shared between trace exporters.
type TraceExporterOptions struct {
	ServiceName string
	Tags        map[string]string
}

// TraceExporterInitFunc is the function that is called to initialize a trace exporter.
type TraceExporterInitFunc func(TraceExporterOptions) (trace.Exporter, error)

var traceExporters sync.Map

// RegisterTraceExporter registers a trace exporter.
func RegisterTraceExporter(name string, f TraceExporterInitFunc) error {
	if _, loaded := traceExporters.LoadOrStore(name, f); loaded {
		return fmt.Errorf("trace exporter %q already registered", name)
	}
	return nil
}

// MustRegisterTraceExporter registers a trace exporter and panics on error.
func MustRegisterTraceExporter(name string, f TraceExporterInitFunc) {
	if err := RegisterTraceExporter(name, f); err != nil {
		panic(err)
	}
}

// GetTraceExporter gets the specified trace exporter passing in the options to the exporter init function.
func GetTraceExporter(name string, opts TraceExporterOptions) (trace.Exporter, error) {
	f, ok := traceExporters.Load(name)
	if !ok {
		return nil, fmt.Errorf("trace exporter %q not found", name)
	}
	return f.(TraceExporterInitFunc)(opts)
}

func availableExporters(exporters sync.Map) []string {
	var names []string
	exporters.Range(func(key, _ interface{}) bool {
		names = append(names, key.(string))
		return true
	})
	sort.Strings(names)
	return names
}

// AvailableTraceExporters gets the names of registered trace exporters.
func AvailableTraceExporters() []string {
	return availableExporters(traceExporters)
}

// WithoutNameSampler returns a trace sampler filtering out a set of span names.
func WithoutNameSampler(name string, names ...string) trace.Sampler {
	return func(params trace.SamplingParameters) trace.SamplingDecision {
		if params.Name == name {
			return trace.SamplingDecision{Sample: false}
		}
		for _, name := range names {
			if params.Name == name {
				return trace.SamplingDecision{Sample: false}
			}
		}
		return trace.SamplingDecision{Sample: true}
	}
}
