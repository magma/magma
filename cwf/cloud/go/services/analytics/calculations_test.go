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

var (
	testGauge = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Name: "test",
	}, []string{"networkID", "label1"})

	testMetricLabels = prometheus.Labels{"label1": "value1"}
)

const (
	testMetricName = "testMetric"
)

var testXAPCalculation = XAPCalculation{
	Days:            7,
	ThresholdBytes:  100,
	QueryStepSize:   time.Second * 100,
	RegisteredGauge: testGauge,
	Labels:          testMetricLabels,
	Name:            testMetricName,
}

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

func TestXAPCalculation(t *testing.T) {
	// Query returns error
	errClient := &mocks.PrometheusAPI{}
	errClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(nil, fmt.Errorf("query error"))

	results, err := testXAPCalculation.Calculate(errClient)
	assert.Error(t, err)
	assert.Len(t, results, 0)

	// Query returns unexpected datatype
	nonVecClient := &mocks.PrometheusAPI{}
	nonVecClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Matrix{}, nil)

	results, err = testXAPCalculation.Calculate(nonVecClient)
	assert.Error(t, err)
	assert.Len(t, results, 0)

	// Query returns no data
	noDataClient := &mocks.PrometheusAPI{}
	noDataClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(model.Vector{}, nil)

	results, err = testXAPCalculation.Calculate(noDataClient)
	assert.EqualError(t, err, "no data returned from query")
	assert.Len(t, results, 0)

	// Query returns expected data
	successClient := &mocks.PrometheusAPI{}
	metric1 := model.Metric{}
	metric1["networkID"] = "testNetwork"
	sample1 := model.Sample{
		Metric:    metric1,
		Value:     1,
		Timestamp: 123,
	}
	var vec model.Vector
	vec = []*model.Sample{&sample1}
	successClient.On("Query", mock.Anything, mock.Anything, mock.Anything).Return(vec, nil)

	expectedResult := Result{
		value:      1,
		metricName: testXAPCalculation.Name,
		labels:     map[string]string{"label1": "value1", "networkID": "testNetwork"},
	}

	results, err = testXAPCalculation.Calculate(successClient)
	assert.NoError(t, err)
	assert.Len(t, results, 1)
	assert.Equal(t, expectedResult, results[0])
}
