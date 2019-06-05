/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package cache

import (
	"testing"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

var (
	testName   = "testName"
	testValue  = "testValue"
	testLabels = []*dto.LabelPair{{Name: &testName, Value: &testValue}}
)

func TestMetricQueue(t *testing.T) {
	gaugeA := prometheus.NewGauge(prometheus.GaugeOpts{Name: "A"})
	gaugeB := prometheus.NewGauge(prometheus.GaugeOpts{Name: "B"})
	gaugeC := prometheus.NewGauge(prometheus.GaugeOpts{Name: "C"})
	gaugeD := prometheus.NewGauge(prometheus.GaugeOpts{Name: "D"})

	metricA := &dto.Metric{}
	metricB := &dto.Metric{}
	metricC := &dto.Metric{}
	metricD := &dto.Metric{}
	gaugeA.Write(metricA)
	gaugeB.Write(metricB)
	gaugeC.Write(metricC)
	gaugeD.Write(metricD)

	queue := NewMetricQueue(3)

	queue.Push(metricA)
	queue.Push(metricB)
	queue.Push(metricC)
	// [A B C]
	assert.Equal(t, metricA, queue.Pop())
	assert.Equal(t, metricB, queue.Pop())
	assert.Equal(t, metricC, queue.Pop())

	queue = NewMetricQueue(3)
	queue.Push(metricA)
	queue.Push(metricB)
	queue.Push(metricC)
	queue.Push(metricD)
	// [B C D] represented by [D B C]
	assert.Equal(t, metricB, queue.Pop())
	assert.Equal(t, metricC, queue.Pop())
	assert.Equal(t, metricD, queue.Pop())

	queue = NewMetricQueue(3)
	queue.Push(metricA)
	// [A _ _]
	assert.Equal(t, metricA, queue.Pop())
	assert.Nil(t, queue.Pop())
	assert.Nil(t, queue.Pop())

	queue = NewMetricQueue(3)
	queue.Push(metricA)
	queue.Push(metricB)
	// [A B _ ]
	assert.Equal(t, metricA, queue.Pop())
	assert.Equal(t, metricB, queue.Pop())

	queue.Push(metricA)
	queue.Push(metricB)
	queue.Push(metricC)
	// [A B C]
	assert.Equal(t, metricA, queue.Pop())
	assert.Equal(t, metricB, queue.Pop())
	assert.Equal(t, metricC, queue.Pop())
}
