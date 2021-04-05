package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/facebookincubator/prometheus-configmanager/alertmanager/config"
	"github.com/imdario/mergo"
	"github.com/labstack/echo"
	amconfig "github.com/prometheus/alertmanager/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/handlers/mocks"
)

const (
	webhookURL = "http://test.com"
	slackURL   = "http://slack.com"
)

var (
	testWebhookConfigBody = fmt.Sprintf(`{
      "name": "test",
      "webhook_configs": [
      {
		  "send_resolved": true,
		  "url": "%s"
      }
      ],
      "slack_configs": [
      {
         "api_url": "%s"
      }
      ]
    }`, webhookURL, slackURL)

	testWebhookURL, _ = url.Parse(webhookURL)
	testWebhookConfig = config.WebhookConfig{
		NotifierConfig: amconfig.NotifierConfig{
			VSendResolved: true,
		},
		URL: &amconfig.URL{
			URL: testWebhookURL,
		},
	}
	testSlackConfig = config.SlackConfig{
		APIURL: slackURL,
	}
	testReceiver = config.Receiver{
		Name: "emptyReceiver",
	}
	testReceiverBytes, _ = json.Marshal([]config.Receiver{testReceiver})

	testRoute = config.Route{
		Receiver: "null",
		Routes:   []*config.Route{},
	}
	testRouteBytes, _ = json.Marshal(&testRoute)

	emptyHTTPResponse     = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}
	empty500Response      = &http.Response{StatusCode: 500, Body: ioutil.NopCloser(bytes.NewBuffer([]byte{}))}
	getReceiverResponse   = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(testReceiverBytes))}
	getAlertRouteResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(testRouteBytes))}
)

type alertConfigAPITestCase struct {
	Name                 string
	RequestURL           string
	ParamNames           []string
	ParamValues          []string
	Payload              interface{}
	ClientMethod         string
	ClientExpectedParams []interface{}
	ClientExpectedReturn []interface{}
	FuncToTest           func(string, HttpClient) func(echo.Context) error
	ExpectedError        string
	ExpectedStatus       int
}

func (tc *alertConfigAPITestCase) RunTest(t *testing.T) {
	mockClient := &mocks.HttpClient{}
	mockClient.On(tc.ClientMethod, tc.ClientExpectedParams...).Return(tc.ClientExpectedReturn...)
	payloadBytes, err := json.Marshal(tc.Payload)
	assert.NoError(t, err)
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, tc.RequestURL, bytes.NewBuffer(payloadBytes)), httptest.NewRecorder())
	c.SetParamNames(tc.ParamNames...)
	c.SetParamValues(tc.ParamValues...)
	err = tc.FuncToTest(tc.RequestURL, mockClient)(c)
	if tc.ExpectedError != "" {
		assert.EqualError(t, err, tc.ExpectedError)
		return
	}
	assert.NoError(t, err)
	assert.Equal(t, tc.ExpectedStatus, c.Response().Status)
}

func TestGetConfigureAlertReceiverHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              testReceiver,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetConfigureAlertReceiverHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid payload",
			Payload:       []int{},
			FuncToTest:    GetConfigureAlertReceiverHandler,
			ExpectedError: "code=400, message=json: cannot unmarshal array into Go value of type config.Receiver",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			FuncToTest:    GetConfigureAlertReceiverHandler,
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			FuncToTest:           GetConfigureAlertReceiverHandler,
			ExpectedError:        "code=500, message=error writing config: <nil>",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetRetrieveAlertReceiverHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              nil,
		ClientMethod:         "Get",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{getReceiverResponse, nil},
		FuncToTest:           GetRetrieveAlertReceiverHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:                 "invalid response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error reading receivers: <nil>",
		},
		{
			Name:                 "request error",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, errors.New("client error")},
			ExpectedError:        "client error",
		},
		{
			Name:                 "invalid response payload",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
			ExpectedStatus:       500,
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetUpdateAlertReceiverHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id", "receiver"},
		ParamValues:          []string{"0", "emptyReceiver"},
		Payload:              testReceiver,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetUpdateAlertReceiverHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid payload",
			Payload:       []int{},
			ExpectedError: "code=400, message=json: cannot unmarshal array into Go value of type config.Receiver",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "receiver name not provided",
			ParamNames:    []string{"network_id"},
			ExpectedError: "code=400, message=receiver name not provided",
		},
		{
			Name:          "receiver name not equal",
			ParamValues:   []string{"0", "test"},
			ExpectedError: "code=400, message=new receiver configuration must have same name",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=code=500, message=error writing config: <nil>",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetDeleteAlertReceiverHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com?receiver=emptyReceiver",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              nil,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetDeleteAlertReceiverHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no receiver name",
			RequestURL:    "http://url.com",
			ExpectedError: "code=400, message=receiver name not provided",
		},
		{
			Name:                 "request non-200",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error deleting receiver: <nil>",
		},
		{
			Name:                 "request error",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("request error")},
			ExpectedError:        "code=500, message=request error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetRetrieveAlertRouteHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              nil,
		ClientMethod:         "Get",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{getAlertRouteResponse, nil},
		FuncToTest:           GetRetrieveAlertRouteHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:                 "server non-500 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error reading alerting route: <nil>",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, errors.New("server error")},
			ExpectedError:        "server error",
		},
		{
			Name:                 "invalid response payload",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
			ExpectedError:        "code=500, message=error decoding server response EOF",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetUpdateAlertRouteHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              testRoute,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetUpdateAlertRouteHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid payload",
			Payload:       []int{},
			ExpectedError: "code=400, message=invalid route specification: json: cannot unmarshal array into Go value of type config.Route",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error updating alert route: error writing config: <nil>",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestBuildReceiverFromContext(t *testing.T) {
	req := httptest.NewRequest(http.MethodPut, "/", strings.NewReader(testWebhookConfigBody))
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)

	receiver, err := buildReceiverFromContext(c)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(receiver.SlackConfigs))
	assert.Equal(t, 1, len(receiver.WebhookConfigs))
	assert.Equal(t, testWebhookConfig, *receiver.WebhookConfigs[0])
	assert.Equal(t, testSlackConfig, *receiver.SlackConfigs[0])
}
