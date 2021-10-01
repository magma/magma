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

// Collector is used as the central point on the Magma access gateway for
// metrics collection. All metrics should be registered with the Collector,
// and periodically, it will send all new metric samples to the MetricSender
type Collector struct {
	clock     clock.Clock
	period    time.Duration
	sender    MetricSender
	gauges    []Gauge
	startOnce sync.Once
	// gates access to gauges
	sync.Mutex
}

func NewCollector(c clock.Clock, period time.Duration, sender MetricSender) *Collector {
	return &Collector{clock: c, period: period, sender: sender, startOnce: sync.Once{}}
}

// RegisterGauge is required for the Collector to be aware of a Gauge-type
// metric. Without it, the Collector will not collect the metric samples, nor
// pass them to the MetricSender
func (c *Collector) RegisterGauge(gauge Gauge) {
	c.Lock()
	defer c.Unlock()
	c.gauges = append(c.gauges, gauge)
}

// Starts periodic collection of registered metric samples, and passes them to
// the MetricSender
func (c *Collector) Start() {
	c.startOnce.Do(func() {
		go func() {
			ticker := c.clock.Ticker(c.period)
			for {
				c.collectMetrics()
				<-ticker.C
			}
		}()
	})
}

func (c *Collector) collectMetrics() {
	c.Lock()
	defer c.Unlock()
	c.sender.Send(c.gauges)
}
