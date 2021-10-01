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
	"testing"

	"github.com/benbjohnson/clock"
	"github.com/stretchr/testify/assert"
)

func TestGaugeBasic(t *testing.T) {
	// Check that values are appended correctly
	gauge := NewGauge(clock.New(), "metric1", "test metric", []string{})
	gauge.Set(1, map[string]string{"label1": "asdf"})
	gauge.Set(2, map[string]string{})
	values := gauge.GetSamples()
	assert.Equal(t, 2, len(values))
	sample1 := values[0]
	sample2 := values[1]
	assert.Equal(t, float64(1), sample1.Value)
	assert.Equal(t, float64(2), sample2.Value)

	// Check that GetSamples empties out what's stored in the Gauge
	values = gauge.GetSamples()
	assert.Equal(t, 0, len(values))
}
