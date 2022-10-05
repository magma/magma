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

package exporters

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.opencensus.io/stats/view"
	"go.opencensus.io/tag"
)

var (
	category = "123"
	token    = "456|789"
	entity   = "my_component"
)

// For testing purpose im omitting the time of the datapoint
type DatapointTest struct {
	Entity string   `json:"entity"`
	Key    string   `json:"key"`
	Value  string   `json:"value"`
	Tags   []string `json:"tags"`
}

// Used to sort datapoints for comparison
type Interface interface {
	// Len is the number of elements in the collection.
	Len() int
	// Less reports whether the element with
	// index i should sort before the element with index j.
	Less(i, j int) bool
	// Swap swaps the elements with indexes i and j.
	Swap(i, j int)
}

type ByKey []DatapointTest

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

type mockConfigProvider struct {
	mock.Mock
}

func (m *mockConfigProvider) getConfig() Config {
	return Config{
		Category:    category,
		Token:       token,
		Entity:      entity,
		DisablePost: true,
		GraphURL:    "https://graph.facebook.com/ods_metrics",
	}
}

func (m *mockConfigProvider) setConfig(cfg Config) Config {
	return m.getConfig()
}

func TestPostODS(t *testing.T) {

	tests := []struct {
		testName    string
		metricsData map[string]string
		err         error
		output      []byte
		resp        func(*assert.Assertions, http.ResponseWriter)
	}{
		{
			testName:    "no_datapoints_should_return_error",
			metricsData: make(map[string]string),
			err:         errors.New("empty datapoints: no valid datapoints found not posting to ODS"),
			output:      nil,
			resp:        nil,
		},
		{
			testName: "actuall_datapoints_should_pass",
			metricsData: map[string]string{
				"success": "1",
				"failed":  "3",
			},
			err:    nil,
			output: []byte(`[{"entity":"dummy.my_component","key":"failed","value":"3","tags":null},{"entity":"dummy.my_component","key":"success","value":"1","tags":null}]`),
			resp: func(assert *assert.Assertions, w http.ResponseWriter) {
				w.WriteHeader(http.StatusOK)
			},
		},
	}

	// Iteratting over tests while mocking ODS graph endpoint each time.
	for _, test := range tests {
		test := test
		t.Run(test.testName, func(t *testing.T) {
			assert := assert.New(t)
			mockGraph := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				test.resp(assert, w)
				v, err := ioutil.ReadAll(r.Body)
				require.NoError(t, err)
				actual, err := getActualURI(v, t)
				expected := "access_token=" + token + "&category_id=" + category + "&datapoints=" + string(test.output)
				require.NoError(t, err)
				assert.Equal(expected, actual)
			}))
			defer mockGraph.Close()
			u, err := url.Parse(mockGraph.URL)
			require.NoError(t, err)
			var odsCfg = Config{
				Category: category,
				Token:    token,
				Entity:   entity,
				GraphURL: u.String(),
			}
			postErr := PostToODS(test.metricsData, odsCfg)
			assert.Equal(test.err, postErr)
		})
	}
}

func getActualURI(v []byte, t *testing.T) (string, error) {
	actual, err := url.QueryUnescape(string(v))
	require.NoError(t, err)

	// Extracting the time from datapoints so comparison will work and rebuild the string
	urlValues, err := url.ParseQuery(string(v))
	require.NoError(t, err)
	var dataPoints []DatapointTest
	err = json.Unmarshal([]byte(urlValues.Get("datapoints")), &dataPoints)
	require.NoError(t, err)
	sort.Sort(ByKey(dataPoints))
	dps, err := json.Marshal(dataPoints)
	actual = actual[:strings.Index(actual, "datapoints")] + "datapoints=" + string(dps)
	require.NoError(t, err)

	return actual, err
}

var ome odsMetricsExporter

func TestExportView(t *testing.T) {

	tests := []struct {
		testName string
		vd       view.Data
		row      *view.Row
		output   map[string]string
	}{
		{
			testName: "no_rows_should_return_error",
			vd:       view.Data{View: &view.View{}, Rows: []*view.Row{}},
			output:   map[string]string{},
		},
		{
			testName: "has_single_count_row_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Success",
					Aggregation: &view.Aggregation{
						Type: view.AggTypeCount,
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.CountData{Value: 1},
					},
				},
			},
			output: map[string]string{"Success.count": "1"},
		},
		{
			testName: "has_count_row_and_tags_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Success",
					Aggregation: &view.Aggregation{
						Type: view.AggTypeCount,
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.CountData{Value: 1},
						Tags: []tag.Tag{{Value: "stam"}},
					},
				},
			},
			output: map[string]string{"stam.Success.count": "1"},
		},
		{
			testName: "has_count_rows_and_tags_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Success",
					Aggregation: &view.Aggregation{
						Type: view.AggTypeCount,
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.CountData{Value: 1},
						Tags: []tag.Tag{{Value: "stam"}},
					},
					{
						Data: &view.CountData{Value: 3},
						Tags: []tag.Tag{{Value: "stam3"}},
					},
				},
			},
			output: map[string]string{"stam.Success.count": "1", "stam3.Success.count": "3"},
		},
		{
			testName: "has_sum_rows_and_tags_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Failed",
					Aggregation: &view.Aggregation{
						Type: view.AggTypeSum,
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.SumData{Value: 2},
						Tags: []tag.Tag{{Value: "stam"}},
					},
				},
			},
			output: map[string]string{"stam.Failed.sum": "2.000000"},
		},
		{
			testName: "has_distribution_rows_and_tags_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Latency",
					Aggregation: &view.Aggregation{
						Type:    view.AggTypeDistribution,
						Buckets: []float64{0, 25, 50, 75, 100, 200, 400},
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.DistributionData{CountPerBucket: []int64{12}, Min: 2, Max: 200, Count: 11, Mean: 100},
						Tags: []tag.Tag{{Value: "some_Request"}},
					},
				},
			},
			output: map[string]string{
				"some_Request.Latency.avg":   "100.000000",
				"some_Request.Latency.count": "11",
				"some_Request.Latency.max":   "200.000000",
				"some_Request.Latency.min":   "2.000000",
				"some_Request.Latency.sum":   "1100.000000",
			},
		},
		{
			testName: "has_single_lastvalue_row_should_pass",
			vd: view.Data{
				View: &view.View{
					Name: "Latency",
					Aggregation: &view.Aggregation{
						Type: view.AggTypeLastValue,
					},
				},
				Rows: []*view.Row{
					{
						Data: &view.LastValueData{Value: 103},
					},
				},
			},
			output: map[string]string{"Latency.gauge": "103.000000"},
		},
	}
	// Iteratting over tests while mocking ODS graph endpoint each time.
	for _, test := range tests {
		test := test
		t.Run(test.testName, func(t *testing.T) {
			assert := assert.New(t)
			// creating the mock
			ome.configProvider = new(mockConfigProvider)
			ome.configProvider.setConfig(Config{})

			ome.ExportView(&test.vd)
			actual := ome.metricsData
			expected := test.output
			assert.True(reflect.DeepEqual(actual, expected), fmt.Sprintf("Actual %v, Expected %v", actual, expected))

		})
	}
}
