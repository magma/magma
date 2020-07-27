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

package helpers

import (
	"context"
	"runtime"
	"time"

	"go.opencensus.io/stats"
	"go.opencensus.io/tag"
)

func elapsed(startTime time.Time) int64 {
	// Converting duration to ms
	return time.Since(startTime).Nanoseconds() / 1000 / 1000
}

// AddComponentTag should be called once when component start to tag all events
// by that name
func AddComponentTag(ctx context.Context, component string) context.Context {
	newCtx, _ := tag.New(ctx, tag.Upsert(KeyComponent, component))
	return newCtx
}

// AddPartnetTag should be called per partner identification in MT systems
func AddPartnetTag(ctx context.Context, partnerShortname string) context.Context {
	newCtx, _ := tag.New(ctx, tag.Upsert(KeyComponent, partnerShortname))
	return newCtx
}

// Operation counterset to track an operation
type Operation struct {
	name      string
	startTime time.Time
	// tags that will be added to the measurement once it is recorded
	mutatorTags []tag.Mutator
}

// NewOperation creates a new operation counter set
func NewOperation(name string) Operation {
	operation := Operation{
		name:        name,
		startTime:   time.Now(),
		mutatorTags: []tag.Mutator{tag.Upsert(KeyOperation, name)},
	}
	return operation
}

// AddTag attaches a tag to the operation
func (o *Operation) AddTag(mutatorTag tag.Mutator) {
	o.mutatorTags = append(o.mutatorTags, mutatorTag)
}

// AddTags attaches a set of tag to the operation
func (o *Operation) AddTags(mutatorTags []tag.Mutator) {
	o.mutatorTags = append(o.mutatorTags, mutatorTags...)
}

// MarkAsSuccess indicates the operation has completed successfully
func (o *Operation) MarkAsSuccess(ctx context.Context) {
	stats.RecordWithTags(ctx, o.mutatorTags, MLatencyMs.M(elapsed(o.startTime)), MSuccess.M(1))
}

// MarkAsError indicates the operation has completed with error
func (o *Operation) MarkAsError(ctx context.Context, e error) {
	keyTags := append(o.mutatorTags, tag.Upsert(KeyError, e.Error()))
	stats.RecordWithTags(ctx, keyTags, MLatencyMs.M(elapsed(o.startTime)), MErrors.M(1))
}

// MarkAsFailed indicates the operation has completed yet failed to do the required action
func (o *Operation) MarkAsFailed(ctx context.Context) {
	stats.RecordWithTags(ctx, o.mutatorTags, MLatencyMs.M(elapsed(o.startTime)), MErrors.M(1))
}

// MarkEvents indicate the event occurred 'count' time - count it
func (o *Operation) MarkEvents(ctx context.Context, count int64) {
	stats.RecordWithTags(ctx, o.mutatorTags, MEvents.M(count))
}

// MarkEvent indicate the event occurred once - count it
func (o *Operation) MarkEvent(ctx context.Context) {
	o.MarkEvents(ctx, 1)
}

// Add64Inc add given value to an In64 measurement
func (o *Operation) Int64Add(ctx context.Context, counter *stats.Int64Measure, add int64) {
	stats.RecordWithTags(ctx, o.mutatorTags, counter.M(add))
}

// Int64Inc increment an In64 measurement
func (o *Operation) Int64Inc(ctx context.Context, counter *stats.Int64Measure) {
	stats.RecordWithTags(ctx, o.mutatorTags, counter.M(1))
}

// CollectRuntimeMetrics should be used to get runtime metrics, it will run every
// minute and will record measurements to opencensus
func CollectRuntimeMetrics(component string) {
	ticker := time.NewTicker(60 * time.Second)
	ctx := context.Background()
	ms := &runtime.MemStats{}
	runtime.ReadMemStats(ms)
	go func(component string) {
		for {
			select {
			case <-ticker.C:
				// Goroutines
				stats.RecordWithTags(
					ctx,
					[]tag.Mutator{tag.Upsert(KeyComponent, component)},
					mGoroutines.M(int64(runtime.NumGoroutine())),
				)
				// Memory only
				stats.RecordWithTags(
					ctx,
					[]tag.Mutator{tag.Upsert(KeyComponent, component)},
					mHeapAllocs.M(int64(ms.HeapAlloc)),
					mFrees.M(int64(ms.Frees)),
					mPtrLookups.M(int64(ms.Lookups)),
					mStackSys.M(int64(ms.StackSys)),
					mHeapObjects.M(int64(ms.HeapObjects)),
					mHeapReleased.M(int64(ms.HeapReleased)),
				)
			}
		}
	}(component)
}
