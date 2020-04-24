// Copyright (c) 2004-present Facebook All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package telemetry_test

import (
	"sort"
	"testing"

	"github.com/facebookincubator/symphony/pkg/telemetry"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.opencensus.io/stats/view"
)

type viewIniter struct {
	mock.Mock
}

func (vi *viewIniter) Init(opts telemetry.ViewExporterOptions) (view.Exporter, error) {
	args := vi.Called(opts)
	exporter, _ := args.Get(0).(view.Exporter)
	return exporter, args.Error(1)
}

func TestGetViewExporter(t *testing.T) {
	_, err := telemetry.GetViewExporter("noexist",
		telemetry.ViewExporterOptions{},
	)
	assert.EqualError(t, err, `view exporter "noexist" not found`)
	var vi viewIniter
	vi.On("Init", mock.Anything).
		Return(nil, nil).Once()
	defer vi.AssertExpectations(t)
	assert.NotPanics(t, func() {
		telemetry.MustRegisterViewExporter(t.Name(), vi.Init)
	})
	defer telemetry.UnregisterViewExporter(t.Name())
	_, err = telemetry.GetViewExporter(t.Name(),
		telemetry.ViewExporterOptions{},
	)
	assert.NoError(t, err)
}

func TestAvailableViewExporters(t *testing.T) {
	var vi viewIniter
	defer vi.AssertExpectations(t)
	suffixes := []string{"foo", "bar", "baz"}
	for _, suffix := range suffixes {
		err := telemetry.RegisterViewExporter(t.Name()+suffix, vi.Init)
		assert.NoError(t, err)
	}
	defer func() {
		for _, suffix := range suffixes {
			telemetry.UnregisterViewExporter(t.Name() + suffix)
		}
	}()
	assert.Panics(t, func() {
		telemetry.MustRegisterViewExporter(t.Name()+suffixes[0], vi.Init)
	})
	exporters := telemetry.AvailableViewExporters()
	assert.True(t, sort.IsSorted(sort.StringSlice(exporters)))
	for _, suffix := range suffixes {
		assert.Contains(t, exporters, t.Name()+suffix)
	}
}
