package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/imdario/mergo"
	"github.com/prometheus/alertmanager/api/v2/client/silence"
	"github.com/prometheus/alertmanager/api/v2/models"
	"github.com/stretchr/testify/mock"
)

var (
	postSilenceOK         = silence.PostSilencesOKBody{}
	postSilenceOKBytes, _ = json.Marshal(postSilenceOK)
	postSilenceResponse   = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(postSilenceOKBytes))}

	activeStatus  = models.SilenceStatusStateActive
	pendingStatus = models.SilenceStatusStatePending
	expiredStatus = models.SilenceStatusStateExpired
	getSilence    = []models.GettableSilence{
		{Status: &models.SilenceStatus{State: &activeStatus}},
		{Status: &models.SilenceStatus{State: &pendingStatus}},
		{Status: &models.SilenceStatus{State: &expiredStatus}},
	}
	getSilenceBytes, _ = json.Marshal(getSilence)
	getSilenceResponse = &http.Response{StatusCode: 200, Body: ioutil.NopCloser(bytes.NewBuffer(getSilenceBytes))}
)

func TestGetPostSilencerHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Post",
		ClientExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything},
		ClientExpectedReturn: []interface{}{postSilenceResponse, nil},
		FuncToTest:           GetPostSilencerHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid payload",
			Payload:       []int{},
			ExpectedError: "code=400, message=json: cannot unmarshal array into Go value of type models.Silence",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "code=500, message=server error",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error posting silencer: ",
		},
		{
			Name:                 "invalid server response",
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

func TestGetGetSilencersHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com?filter={a=b,c!=d,e=~f,g!~h}",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Get",
		ClientExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything},
		ClientExpectedReturn: []interface{}{getSilenceResponse, nil},
		FuncToTest:           GetGetSilencersHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:          "invalid filters",
			RequestURL:    "http://url.com?filter={!test}",
			ExpectedError: "code=400, message=bad matcher format: !test",
		},
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "code=500, message=error getting silences: server error",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error getting silences: ",
		},
		{
			Name:                 "invalid server response",
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

func TestGetDeleteSilencersHandler(t *testing.T) {
	baseTest := alertConfigAPITestCase{
		Name:                 "successful request",
		RequestURL:           "http://url.com",
		ParamNames:           []string{"network_id"},
		ParamValues:          []string{"0"},
		ClientMethod:         "Do",
		ClientExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything},
		ClientExpectedReturn: []interface{}{getSilenceResponse, nil},
		FuncToTest:           GetDeleteSilencerHandler,
		ExpectedStatus:       200,
	}
	tests := []alertConfigAPITestCase{
		baseTest,
		{
			Name:                 "server error",
			ClientExpectedReturn: []interface{}{empty500Response, errors.New("server error")},
			ExpectedError:        "code=500, message=error deleting silence: server error",
		},
		{
			Name:                 "server non-200 response",
			ClientExpectedReturn: []interface{}{empty500Response, nil},
			ExpectedError:        "code=500, message=error deleting silence: ",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}
