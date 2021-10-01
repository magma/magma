// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package metrics

import (
	"sync"
	"time"

	"github.com/benbjohnson/clock"
)

// Gauge represents an integer metric value.
// This is intended to be similar to:
// 	- the Prometheus Gauge
//  - Google Cloud DOUBLE gauge metric
// Each Gauge should be defined with a minimal set of label names. For each
// published sample value, there should be a corresponding value for each label
// name.
//	ie.
//	For a Gauge representing "cpu_usage_percentage", we may want a label named
//	"service_name". We could use this in the following way:
//	```
//	name          := "cpu_usage_percentage"
//	description   := "Percentage of CPU resources used by AGW service"
//	labelNames    := []string{"service_name"}
//  cpuUsageGauge := NewGauge(name, description, labelNames)
//	// At this instant, 4.4% of CPU resources used by service agwd
//  cpuUsageGauge.Set(4.4, map[string]string{"service_name": "agwd"})
//  // And 3.5% of CPU resources are used by service enodebd
//  cpuUsageGauge.Set(3.5, map[string]string{"service_name": "enodebd"})
//	```
//
//  Each Gauge can thus represent a separate time series for each combination
//  of label values.
//
//	NOTE: Due to current limitations with Prometheus, it is not recommended to
//        use labels to represent dimensions with high cardinality (many
//        different label values), such as IMSI, email addresses, or other
//        unbounded sets of values.
type Gauge interface {
	// Set sets the Gauge to an arbitrary value.
	Set(value float64, labels map[string]string)

	// GetSamples returns all values set since it was last called.
	GetSamples() []GaugeSample

	Name() string

	Description() string
}

type GaugeSample struct {
	Value       float64
	LabelValues map[string]string
	TimestampMs int64
}

type gaugeImpl struct {
	name        string
	description string
	labelNames  []string
	clock       clock.Clock
	samples     []GaugeSample
	sync.Mutex
}

func NewGauge(clock clock.Clock, name string, description string, labelNames []string) Gauge {
	return &gaugeImpl{name: name, description: description, labelNames: labelNames, clock: clock}
}

// Set sets the Gauge to an arbitrary value.
// NOTE: labels is a pointer to a map, so the caller should not change it
func (g *gaugeImpl) Set(value float64, labels map[string]string) {
	g.Lock()
	defer g.Unlock()

	g.samples = append(g.samples, GaugeSample{value, labels, int64(time.Nanosecond) * g.clock.Now().UnixNano() / int64(time.Millisecond)})
}

// GetSamples returns all values set since it was last called.
func (g *gaugeImpl) GetSamples() []GaugeSample {
	g.Lock()
	defer g.Unlock()
	samples := g.samples
	g.samples = nil
	return samples
}

func (g *gaugeImpl) Name() string {
	return g.name
}

func (g *gaugeImpl) Description() string {
	return g.description
}
