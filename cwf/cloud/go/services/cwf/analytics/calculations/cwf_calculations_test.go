/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package calculations

import (
	"fmt"
	"testing"
	"time"

	"magma/orc8r/cloud/go/services/analytics/calculations"
	"magma/orc8r/cloud/go/services/analytics/protos"
	"magma/orc8r/cloud/go/services/analytics/query_api"
	"magma/orc8r/cloud/go/services/analytics/query_api/mocks"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// mocked prometheus clients that are used for multiple calculation types
var (
	errClient          = &mocks.PrometheusAPI{}
	matrixReturnClient = &mocks.PrometheusAPI{}
	vectorReturnClient = &mocks.PrometheusAPI{}
)

// Initialize mocked Prometheus clients
func init() {
	// Query returns error
	errClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, fmt.Errorf("query error"))
	errClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(nil, nil, fmt.Errorf("query error"))

	// Query returns matrix datatype
	matrixReturnClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Matrix{}, nil, nil)
	matrixReturnClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(model.Matrix{}, nil, nil)

	// Query returns vector datatype
	vectorReturnClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Vector{}, nil, nil)
	vectorReturnClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(model.Vector{}, nil, nil)
}

var (
	basicLabels = prometheus.Labels{"days": "7"}

	testXAPGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test",
	}, []string{"networkID", "days"})

	testAPThroughputGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test",
	}, []string{"networkID", "direction", "apn", "days"})

	testUserThroughputGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test",
	}, []string{"networkID", "direction", "days"})

	testUserConsumptionGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test",
	}, []string{"networkID", "direction", "days"})
)

const (
	testMetricName = "testMetric"
)

type calculationTestCase struct {
	client          query_api.PrometheusAPI
	calculation     calculations.Calculation
	expectedError   string
	expectedResults []*protos.CalculationResult
	name            string
}

func (tc calculationTestCase) RunTest(t *testing.T) {
	results, err := tc.calculation.Calculate(tc.client)
	if tc.expectedError != "" {
		assert.EqualError(t, err, tc.expectedError)
	} else {
		assert.NoError(t, err)
	}
	assert.Equal(t, results, tc.expectedResults)
}

var exampleXAPCalculation = XAPCalculation{
	BaseCalculation: calculations.BaseCalculation{
		CalculationParams: calculations.CalculationParams{
			Days:            7,
			RegisteredGauge: testXAPGauge,
			Labels:          basicLabels,
			Name:            testMetricName,
		},
	},
	ThresholdBytes: 100,
}

func TestXAPCalculation(t *testing.T) {
	// setup successful expected result
	metric1 := model.Metric{}
	metric1["networkID"] = "testNetwork"
	sample1 := model.Sample{
		Metric:    metric1,
		Value:     1,
		Timestamp: 123,
	}
	var vec model.Vector = []*model.Sample{&sample1}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil, nil)

	expectedSuccessResult := &protos.CalculationResult{
		Value:      1,
		MetricName: exampleXAPCalculation.Name,
		Labels:     map[string]string{"days": "7", "networkID": "testNetwork"},
	}

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "user Consumption query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        matrixReturnClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "user Consumption query error: unexpected ValueType: matrix",
		},
		{
			name:          "No query data",
			client:        vectorReturnClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "user Consumption query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleXAPCalculation,
			expectedResults: []*protos.CalculationResult{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleAPThroughputCalculation = APNThroughputCalculation{
	BaseCalculation: calculations.BaseCalculation{
		CalculationParams: calculations.CalculationParams{
			Days:            7,
			RegisteredGauge: testAPThroughputGauge,
			Labels:          basicLabels,
			Name:            testMetricName,
		},
	},
	QueryStepSize: time.Second,
	Direction:     calculations.ConsumptionIn,
}

func TestAPThroughputCalculation(t *testing.T) {
	metric1 := model.Metric{}
	metric1["apn"] = "apn1"
	metric1["networkID"] = "network1"

	// Values: 1, 2, 3. average is (1+2+3)/3 = 2
	values := []model.SamplePair{
		{
			Value: 1,
		}, {
			Value: 2,
		}, {
			Value: 3,
		},
	}

	matrix := model.Matrix{{
		Metric: metric1,
		Values: values,
	}}

	expectedSuccessResult := &protos.CalculationResult{
		Value:      2,
		MetricName: exampleAPThroughputCalculation.Name,
		Labels:     map[string]string{"apn": "apn1", "networkID": "network1", "days": "7", "direction": string(calculations.ConsumptionIn)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(matrix, nil, nil)

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleAPThroughputCalculation,
			expectedError: "AP Throughput query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        vectorReturnClient,
			calculation:   &exampleAPThroughputCalculation,
			expectedError: "AP Throughput query error: unexpected ValueType: vector",
		},
		{
			name:          "No query data",
			client:        matrixReturnClient,
			calculation:   &exampleAPThroughputCalculation,
			expectedError: "AP Throughput query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleAPThroughputCalculation,
			expectedResults: []*protos.CalculationResult{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleUserThroughputCalculation = UserThroughputCalculation{
	BaseCalculation: calculations.BaseCalculation{
		CalculationParams: calculations.CalculationParams{
			Days:            7,
			RegisteredGauge: testUserThroughputGauge,
			Labels:          basicLabels,
			Name:            testMetricName,
		},
	},
	QueryStepSize: time.Second,
	Direction:     calculations.ConsumptionIn,
}

func TestUserThroughputCalculation(t *testing.T) {
	metric1 := model.Metric{}
	metric1["networkID"] = "network1"

	// Values: 1, 2, 3. average is (1+2+3)/3 = 2
	values := []model.SamplePair{
		{
			Value: 1,
		}, {
			Value: 2,
		}, {
			Value: 3,
		},
	}

	matrix := model.Matrix{{
		Metric: metric1,
		Values: values,
	}}

	expectedSuccessResult := &protos.CalculationResult{
		Value:      2,
		MetricName: exampleAPThroughputCalculation.Name,
		Labels:     map[string]string{"networkID": "network1", "days": "7", "direction": string(exampleUserThroughputCalculation.Direction)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(matrix, nil, nil)

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "user Throughput query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        vectorReturnClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "user Throughput query error: unexpected ValueType: vector",
		},
		{
			name:          "No query data",
			client:        matrixReturnClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "user Throughput query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleUserThroughputCalculation,
			expectedResults: []*protos.CalculationResult{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleUserConsumptionCalculation = UserConsumptionCalculation{
	BaseCalculation: calculations.BaseCalculation{
		CalculationParams: calculations.CalculationParams{
			Days:            7,
			RegisteredGauge: testUserConsumptionGauge,
			Labels:          basicLabels,
			Name:            testMetricName,
		},
	},
	Direction: calculations.ConsumptionIn,
}

func TestUserConsumptionCalculation(t *testing.T) {
	metric1 := model.Metric{}
	metric1["networkID"] = "network1"

	vec := model.Vector{{
		Metric: metric1,
		Value:  2,
	}}

	expectedSuccessResult := &protos.CalculationResult{
		Value:      2,
		MetricName: exampleUserConsumptionCalculation.Name,
		Labels:     map[string]string{"networkID": "network1", "days": "7", "direction": string(exampleUserConsumptionCalculation.Direction)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil, nil)

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "user Consumption query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        matrixReturnClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "user Consumption query error: unexpected ValueType: matrix",
		},
		{
			name:          "No query data",
			client:        vectorReturnClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "user Consumption query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleUserConsumptionCalculation,
			expectedResults: []*protos.CalculationResult{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

func TestCheckLabelsMatch(t *testing.T) {
	expectedLabels := []string{"label1", "label2"}
	assert.True(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"label1": "val", "label2": "val"}))
	assert.True(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"label2": "val", "label1": "val"}))

	assert.False(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"label1": "val"}))
	assert.False(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"label2": "val"}))
	assert.False(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"newLabel": "val"}))
	assert.False(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{}))
	assert.False(t, calculations.CheckLabelsMatch(expectedLabels, prometheus.Labels{"label2": "val", "label1": "val", "newLabel": "val"}))
}
