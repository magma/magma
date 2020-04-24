/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package status implements magmad status amd metrics collectors & reporters
package status

import (
	"sync"

	prometheus "github.com/prometheus/client_model/go"
)

// MetricsQueue - definition of metrics queue "object"
// MetricsQueue is not a traditional queue, it only implements:
//   append, collect/reset & prepend with a limit functionality
type MetricsQueue struct {
	sync.Mutex
	items []*prometheus.MetricFamily
}

// Append adds elements to the end of the queue
func (q *MetricsQueue) Append(elems ...*prometheus.MetricFamily) *MetricsQueue {
	if q != nil {
		q.Lock()
		q.items = append(q.items, elems...)
		q.Unlock()
	}
	return q
}

// Collect returns the que items and empties the queue
func (q *MetricsQueue) Collect() []*prometheus.MetricFamily {
	var result = []*prometheus.MetricFamily{}
	if q != nil {
		q.Lock()
		result, q.items = q.items, result
		q.Unlock()
	}
	return result
}

// Prepend inserts up to len(elems) elements at the start of the queue while maintaining given max len of the queue
// it returns the actual number of elements inserted
// Note: if len(elems) > # elements inserted, inserted elements are taken from the end elems slice
func (q *MetricsQueue) Prepend(elems []*prometheus.MetricFamily, maxQueueLen int) int {
	if q == nil {
		return 0
	}
	el := len(elems)
	q.Lock()
	startIdx := el + len(q.items) - maxQueueLen
	if startIdx < el {
		if startIdx <= 0 {
			q.items = append(elems, q.items...)
		} else {
			q.items = append(elems[startIdx:], q.items...)
			el -= startIdx
		}
		return el
	}
	return 0
}

// Reset clears the queue
func (q *MetricsQueue) Reset() {
	if q != nil {
		q.Lock()
		q.items = []*prometheus.MetricFamily{}
		q.Unlock()
	}
}
