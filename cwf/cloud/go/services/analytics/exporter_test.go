package analytics

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"

	"magma/cwf/cloud/go/services/analytics/mocks"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	metricsPrefix   = "prefix"
	appSecret       = "abc"
	appID           = "123"
	metricExportURL = "export.com"
	categoryName    = "category"
)

const (
	wwwDatapointJSONStringTemplate = `[{"entity":"%s","key":"%s","value":"%f"}]`
	exportURLTemplate              = `%s?access_token=%s|%s`
)

var (
	sampleResult = Result{
		value:      1,
		metricName: "testMetric",
		labels:     prometheus.Labels{"networkID": "testNetwork", "label1": "labelValue"},
	}

	noNetworkResult = Result{
		value:      1,
		metricName: "testMetric",
		labels:     prometheus.Labels{},
	}

	testExporter = &wwwExporter{
		metricsPrefix:   metricsPrefix,
		appSecret:       appSecret,
		appID:           appID,
		metricExportURL: metricExportURL,
		categoryName:    categoryName,
	}
)

type exportTestCase struct {
	client             HttpClient
	exporter           Exporter
	exportResult       Result
	expectedError      string
	assertExpectations func(t *testing.T)
	name               string
}

func (tc exportTestCase) RunTest(t *testing.T) {
	err := tc.exporter.Export(tc.exportResult, tc.client)
	if tc.expectedError != "" {
		assert.EqualError(t, err, tc.expectedError)
	} else {
		assert.NoError(t, err)
	}
	tc.assertExpectations(t)
}

func TestWwwExporter_Export(t *testing.T) {
	errClient := &mocks.HttpClient{}
	errClient.On("PostForm", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("error making post"))

	badStatusClient := &mocks.HttpClient{}
	badStatusResponse := &http.Response{
		Body:       ioutil.NopCloser(bytes.NewBufferString("bad status")),
		StatusCode: 404,
	}
	badStatusClient.On("PostForm", mock.Anything, mock.Anything).Return(badStatusResponse, nil)

	successClient := &mocks.HttpClient{}
	successResponse := &http.Response{
		StatusCode: 200,
	}
	successClient.On("PostForm", mock.Anything, mock.Anything).Return(successResponse, nil)

	testCases := []exportTestCase{
		{
			name:          "Post form error",
			client:        errClient,
			exporter:      testExporter,
			exportResult:  sampleResult,
			expectedError: "error making post",
			assertExpectations: func(t *testing.T) {
				errClient.AssertCalled(t, "PostForm", mock.Anything, mock.Anything)
			},
		},
		{
			name:          "Bad client status",
			client:        badStatusClient,
			exporter:      testExporter,
			exportResult:  sampleResult,
			expectedError: "bad status",
			assertExpectations: func(t *testing.T) {
				badStatusClient.AssertCalled(t, "PostForm", mock.Anything, mock.Anything)
			},
		},
		{
			name:          "No networkID in result",
			client:        successClient,
			exporter:      testExporter,
			exportResult:  noNetworkResult,
			expectedError: "no networkID for exported metric",
			assertExpectations: func(t *testing.T) {
				successClient.AssertNotCalled(t, "PostForm", mock.Anything, mock.Anything)
			},
		},
		{
			name:          "Successful export",
			client:        successClient,
			exporter:      testExporter,
			exportResult:  sampleResult,
			expectedError: "",
			assertExpectations: func(t *testing.T) {
				expectedURL := fmt.Sprintf(exportURLTemplate, metricExportURL, appID, appSecret)
				expectedDatapointJSON := fmt.Sprintf(wwwDatapointJSONStringTemplate, testExporter.FormatEntity(sampleResult, "testNetwork"), testExporter.FormatKey(sampleResult), sampleResult.value)
				expectedPostData := url.Values{"datapoints": {string(expectedDatapointJSON)}, "category": {categoryName}}
				successClient.AssertCalled(t, "PostForm", expectedURL, expectedPostData)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, test.RunTest)
	}
}
