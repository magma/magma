package analytics

import (
	"fmt"
	"testing"
	"time"

	"magma/cwf/cloud/go/services/analytics/mocks"

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

// Initalize mocked Prometheus clients
func init() {
	// Query returns error
	errClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("query error"))
	errClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("query error"))

	// Query returns matrix datatype
	matrixReturnClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Matrix{}, nil)
	matrixReturnClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(model.Matrix{}, nil)

	// Query returns vector datatype
	vectorReturnClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Vector{}, nil)
	vectorReturnClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(model.Vector{}, nil)
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
	client          PrometheusAPI
	calculation     Calculation
	expectedError   string
	expectedResults []Result
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
	CalculationParams: CalculationParams{
		Days:            7,
		RegisteredGauge: testXAPGauge,
		Labels:          basicLabels,
		Name:            testMetricName,
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
	var vec model.Vector
	vec = []*model.Sample{&sample1}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil)

	expectedSuccessResult := Result{
		value:      1,
		metricName: exampleXAPCalculation.Name,
		labels:     map[string]string{"days": "7", "networkID": "testNetwork"},
	}

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "User Consumption query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        matrixReturnClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "User Consumption query error: unexpected ValueType: matrix",
		},
		{
			name:          "No query data",
			client:        vectorReturnClient,
			calculation:   &exampleXAPCalculation,
			expectedError: "User Consumption query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleXAPCalculation,
			expectedResults: []Result{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleAPThroughputCalculation = APThroughputCalculation{
	CalculationParams: CalculationParams{
		Days:            7,
		RegisteredGauge: testAPThroughputGauge,
		Labels:          basicLabels,
		Name:            testMetricName,
	},
	QueryStepSize: time.Second,
	Direction:     ConsumptionIn,
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

	expectedSuccessResult := Result{
		value:      2,
		metricName: exampleAPThroughputCalculation.Name,
		labels:     map[string]string{"apn": "apn1", "networkID": "network1", "days": "7", "direction": string(ConsumptionIn)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(matrix, nil)

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
			expectedResults: []Result{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleUserThroughputCalculation = UserThroughputCalculation{
	CalculationParams: CalculationParams{
		Days:            7,
		RegisteredGauge: testUserThroughputGauge,
		Labels:          basicLabels,
		Name:            testMetricName,
	},
	QueryStepSize: time.Second,
	Direction:     ConsumptionIn,
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

	expectedSuccessResult := Result{
		value:      2,
		metricName: exampleAPThroughputCalculation.Name,
		labels:     map[string]string{"networkID": "network1", "days": "7", "direction": string(exampleUserThroughputCalculation.Direction)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("QueryRange", mock.Anything, mock.Anything, mock.Anything).Return(matrix, nil)

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "User Throughput query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        vectorReturnClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "User Throughput query error: unexpected ValueType: vector",
		},
		{
			name:          "No query data",
			client:        matrixReturnClient,
			calculation:   &exampleUserThroughputCalculation,
			expectedError: "User Throughput query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleUserThroughputCalculation,
			expectedResults: []Result{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}

var exampleUserConsumptionCalculation = UserConsumptionCalculation{
	CalculationParams: CalculationParams{
		Days:            7,
		RegisteredGauge: testUserConsumptionGauge,
		Labels:          basicLabels,
		Name:            testMetricName,
	},
	Direction: ConsumptionIn,
}

func TestUserConsumptionCalculation(t *testing.T) {
	metric1 := model.Metric{}
	metric1["networkID"] = "network1"

	vec := model.Vector{{
		Metric: metric1,
		Value:  2,
	}}

	expectedSuccessResult := Result{
		value:      2,
		metricName: exampleUserConsumptionCalculation.Name,
		labels:     map[string]string{"networkID": "network1", "days": "7", "direction": string(exampleUserConsumptionCalculation.Direction)},
	}

	successClient := &mocks.PrometheusAPI{}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil)

	testCases := []calculationTestCase{
		{
			name:          "Client Error",
			client:        errClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "User Consumption query error: query error",
		},
		{
			name:          "Unexpected query data",
			client:        matrixReturnClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "User Consumption query error: unexpected ValueType: matrix",
		},
		{
			name:          "No query data",
			client:        vectorReturnClient,
			calculation:   &exampleUserConsumptionCalculation,
			expectedError: "User Consumption query error: no data returned from query",
		},
		{
			name:            "Successful query",
			client:          successClient,
			calculation:     &exampleUserConsumptionCalculation,
			expectedResults: []Result{expectedSuccessResult},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}
