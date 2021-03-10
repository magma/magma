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

package collection_test

import (
	"sort"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/metricsd/collection"

	prometheus_proto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

type TestMetricCollector struct {
	ret []*prometheus_proto.MetricFamily
}

func (t *TestMetricCollector) GetMetrics() ([]*prometheus_proto.MetricFamily, error) {
	return t.ret, nil
}

func TestMetricsGatherer_Gather(t *testing.T) {
	expected1 := collection.MakeSingleGaugeFamily("name1", "help1", nil, 12.34)
	expected2 := collection.MakeSingleGaugeFamily("name2", "help2", &collection.MetricLabel{Name: "lname", Value: "lvalue"}, 56.78)
	expected3 := collection.MakeSingleGaugeFamily("name3", "help3", nil, 0)

	output := make(chan *prometheus_proto.MetricFamily)
	gatherer, err := collection.NewMetricsGatherer(
		[]collection.MetricCollector{
			&TestMetricCollector{ret: []*prometheus_proto.MetricFamily{expected1}},
			&TestMetricCollector{ret: []*prometheus_proto.MetricFamily{expected2, expected3}},
		},
		time.Second,
		output,
	)
	assert.NoError(t, err)

	go gatherer.Run()
	timeout := make(chan struct{}, 1)
	go func() {
		time.Sleep(5 * time.Second)
		timeout <- struct{}{}
	}()

	var actual []*prometheus_proto.MetricFamily
	for i := 0; i < 3; i++ {
		select {
		case recv := <-output:
			actual = append(actual, recv)
		case <-timeout:
			assert.Fail(t, "Did not gather expected metrics within timeout")
		}
	}

	assert.Equal(t, 3, len(actual))
	expected := []*prometheus_proto.MetricFamily{expected1, expected2, expected3}
	sort.Slice(actual, func(i, j int) bool { return *actual[i].Name < *actual[j].Name })
	assert.Equal(t, expected, actual)
}
