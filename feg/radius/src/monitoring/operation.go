/*
Copyright 2020 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package monitoring

import (
	"context"
	"fmt"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

// Operation counterset to track an operation
type Operation struct {
	name           string
	startTime      int64
	ctx            context.Context
	tagMutators    []tag.Mutator
	start          *stats.Int64Measure
	success        *stats.Int64Measure
	successLatency *stats.Int64Measure
	failed         *stats.Int64Measure
	failedLatency  *stats.Int64Measure
}

// NewOperation creates a new operation counter set
func NewOperation(name string, tags ...tag.Mutator) Operation {
	operation := Operation{
		name:        name,
		tagMutators: tags,
		startTime:   0,
		ctx:         context.Background(),
		start: stats.Int64(
			fmt.Sprintf("%s/start", name),
			fmt.Sprintf("Operation '%s' started", name),
			stats.UnitDimensionless,
		),
		success: stats.Int64(
			fmt.Sprintf("%s/success", name),
			fmt.Sprintf("Operation '%s' succeeded", name),
			stats.UnitDimensionless,
		),
		successLatency: stats.Int64(
			fmt.Sprintf("%s/success_latency", name),
			fmt.Sprintf("Operation '%s' success latency", name),
			stats.UnitMilliseconds,
		),
		failed: stats.Int64(
			fmt.Sprintf("%s/failed", name),
			fmt.Sprintf("Operation '%s' failed", name),
			stats.UnitDimensionless,
		),
		failedLatency: stats.Int64(
			fmt.Sprintf("%s/failed_latency", name),
			fmt.Sprintf("Operation '%s' failure latency", name),
			stats.UnitMilliseconds,
		),
	}
	views := []*view.View{
		{
			Name:        fmt.Sprintf("%s/start", name),
			Measure:     operation.start,
			Description: fmt.Sprintf("The number of time '%s' was started", name),
			Aggregation: view.Count(),
			TagKeys:     AllTagKeys(),
		},
		{
			Name:        fmt.Sprintf("%s/failure", name),
			Measure:     operation.failed,
			Description: fmt.Sprintf("The number of time '%s' has failed", name),
			Aggregation: view.Count(),
			TagKeys:     AllTagKeys(),
		},
		{
			Name:        fmt.Sprintf("%s/failure/latency", name),
			Measure:     operation.failedLatency,
			Description: fmt.Sprintf("The latency of failed '%s' operations", name),
			Aggregation: view.Distribution(),
			TagKeys:     AllTagKeys(),
		},
		{
			Name:        fmt.Sprintf("%s/success/latency", name),
			Measure:     operation.successLatency,
			Description: fmt.Sprintf("The latency of successful '%s' operations", name),
			Aggregation: view.Distribution(),
			TagKeys:     AllTagKeys(),
		},
		{
			Name:        fmt.Sprintf("%s/success", name),
			Measure:     operation.success,
			Description: fmt.Sprintf("The number of time '%s' has succeeded", name),
			Aggregation: view.Count(),
			TagKeys:     AllTagKeys(),
		},
	}

	view.Register(views...)

	return operation
}

// Start indicates the operation has started
func (o Operation) Start(instanceTags ...tag.Mutator) Operation {
	newOp := o
	newOp.tagMutators = append(o.tagMutators, instanceTags...)
	newOp.startTime = time.Now().UnixNano() / int64(time.Millisecond)
	stats.RecordWithTags(
		newOp.ctx,
		newOp.tagMutators,
		newOp.start.M(1),
	)
	return newOp
}

// Success indicates the operation has completed successfully
func (o Operation) Success(tags ...tag.Mutator) {
	n := time.Now().UnixNano() / int64(time.Millisecond)
	stats.RecordWithTags(
		o.ctx,
		append(
			o.tagMutators,
			tags...,
		),
		o.success.M(1),
		o.successLatency.M(n-o.startTime),
	)
	o.startTime = 0
}

// Failure indicates the operation has completed successfully
func (o Operation) Failure(errorCode string, tags ...tag.Mutator) {
	n := time.Now().UnixNano() / int64(time.Millisecond)
	stats.RecordWithTags(
		o.ctx,
		append(
			o.tagMutators,
			append(
				tags,
				tag.Upsert(ErrorCodeTag, errorCode),
			)...,
		),
		o.failed.M(1),
		o.failedLatency.M(n-o.startTime),
	)
	o.startTime = 0
}
