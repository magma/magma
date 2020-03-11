// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opencensus.io/stats/view"
)

type viewIniter struct {
	mock.Mock
}

func (vi *viewIniter) Init(opts ViewExporterOptions) (view.Exporter, error) {
	args := vi.Called(opts)
	exporter, _ := args.Get(0).(view.Exporter)
	return exporter, args.Error(1)
}

func TestGetViewExporter(t *testing.T) {
	_, err := GetViewExporter("noexist", ViewExporterOptions{})
	assert.EqualError(t, err, `view exporter "noexist" not found`)
	var vi viewIniter
	vi.On("Init", mock.Anything).Return(nil, nil).Once()
	defer vi.AssertExpectations(t)
	assert.NotPanics(t, func() { MustRegisterViewExporter(t.Name(), vi.Init) })
	defer traceExporters.Delete(t.Name())
	_, err = GetViewExporter(t.Name(), ViewExporterOptions{})
	assert.NoError(t, err)
}

func TestAvailableViewExporters(t *testing.T) {
	var vi viewIniter
	defer vi.AssertExpectations(t)
	suffixes := []string{"foo", "bar", "baz"}
	for _, suffix := range suffixes {
		err := RegisterViewExporter(t.Name()+suffix, vi.Init)
		assert.NoError(t, err)
	}
	defer func() {
		for _, suffix := range suffixes {
			viewExporters.Delete(t.Name() + suffix)
		}
	}()
	assert.Panics(t, func() {
		MustRegisterViewExporter(t.Name()+suffixes[0], vi.Init)
	})
	exporters := AvailableViewExporters()
	assert.True(t, sort.IsSorted(sort.StringSlice(exporters)))
	for _, suffix := range suffixes {
		assert.Contains(t, exporters, t.Name()+suffix)
	}
}
