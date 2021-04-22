package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/facebookincubator/prometheus-configmanager/prometheus/alert"
	"github.com/imdario/mergo"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/metricsd/obsidian/models"
	"magma/orc8r/cloud/go/services/metricsd/test_common"
)

var (
	testRule = alert.RuleJSONWrapper{
		Alert: "test",
		Expr:  "up",
	}
	testRuleBytes, _ = json.Marshal([]alert.RuleJSONWrapper{testRule})
	badExprRule      = alert.RuleJSONWrapper{
		Alert: "test",
		Expr:  "!up",
	}
	badLabelsRule = alert.RuleJSONWrapper{
		Alert:  "test",
		Expr:   "up",
		Labels: map[string]string{"!labelName": "val"},
	}

	testBulkRules    = []alert.RuleJSONWrapper{{Alert: "test1", Expr: "up"}, {Alert: "test2", Expr: "up"}}
	testRuleResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(testRuleBytes))}

	firingAlerts = []models.GettableAlert{{
		Name: test_common.MakeStrPtr("test"),
	}}
	firingAlertsBytes, _ = json.Marshal(firingAlerts)
	firingAlertsResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(firingAlertsBytes))}
)

func TestGetConfigurePrometheusAlertHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              testRule,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetConfigurePrometheusAlertHandler,
		ExpectedStatus:       201,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid rule expression",
			Payload:       badExprRule,
			ExpectedError: "code=400, message=error parsing query: parse error at char 1: unexpected character after '!': 'u'",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "invalid rule",
			Payload:       badLabelsRule,
			ExpectedError: "code=400, message=invalid rule: [invalid label name: !labelName]\n",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=code=500, message=error writing config: <nil>",
		},
		{
			Name:          "bad payload",
			Payload:       []int{},
			ExpectedError: "code=400, message=misconfigured rule: json: cannot unmarshal array into Go value of type alert.RuleJSONWrapper",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetRetrieveAlertRuleHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com?alert_name=test",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Get",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{testRuleResponse, nil},
		FuncToTest:           GetRetrieveAlertRuleHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error reading rules: <nil>",
		},
		{
			Name:                 "server error response",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "server error",
		},
		{
			Name:                 "bad response body",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
			ExpectedError:        "code=500, message=error decoding server response: EOF",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetDeleteAlertRuleHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com?alert_name=test",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetDeleteAlertRuleHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "no alert name",
			RequestURL:    "http://url.com",
			ExpectedError: "code=400, message=alert name not provided",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error deleting rule: <nil>",
		},
		{
			Name:                 "server error response",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "server error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetUpdateAlertRuleHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id", "alert_name"},
		ParamValues:          []string{"0", "test"},
		Payload:              testRule,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetUpdateAlertRuleHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "no alert name",
			ParamNames:    []string{"network_id"},
			ExpectedError: "code=400, message=alert name not provided",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=code=500, message=error writing config: <nil>",
		},
		{
			Name:          "bad payload",
			Payload:       []int{},
			ExpectedError: "code=400, message=misconfigured rule: json: cannot unmarshal array into Go value of type alert.RuleJSONWrapper",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetBulkUpdateAlertHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		Payload:              testBulkRules,
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
		FuncToTest:           GetBulkUpdateAlertHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error writing config: <nil>",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "code=500, message=make PUT request: server error",
		},
		{
			Name:          "bad payload",
			Payload:       "string",
			ExpectedError: "code=400, message=error parsing rule payload: json: cannot unmarshal string into Go value of type []alert.RuleJSONWrapper",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetViewFiringAlertsHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Get",
		ClientExpectedParams: []interface{}{mock.Anything},
		ClientExpectedReturn: []interface{}{firingAlertsResponse, nil},
		FuncToTest:           GetViewFiringAlertHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=alertmanager error: <nil>",
		},
		{
			Name:                 "bad response payload",
			ClientExpectedReturn: []interface{}{emptyHTTPResponse, nil},
			ExpectedError:        "code=500, message=error decoding alertmanager response: EOF",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}
