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

package handlers

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"sort"
	"testing"

	v1 "github.com/prometheus/client_golang/api/prometheus/v1"
	"github.com/stretchr/testify/mock"

	"magma/orc8r/cloud/go/services/metricsd/prometheus/handlers/mocks"
	"magma/orc8r/cloud/go/services/metricsd/prometheus/restrictor"
	"magma/orc8r/cloud/go/services/tenants"
	tenants_test_init "magma/orc8r/cloud/go/services/tenants/test_init"
	"magma/orc8r/lib/go/protos"

	"github.com/imdario/mergo"
	"github.com/labstack/echo"
	"github.com/prometheus/common/model"
	"github.com/stretchr/testify/assert"
)

type prometheusAPITestCase struct {
	Name              string
	RequestURL        string
	ParamNames        []string
	ParamValues       []string
	ApiMethod         string
	ApiExpectedParams []interface{}
	ApiExpectedReturn []interface{}
	FuncToTest        func(PrometheusAPI) func(echo.Context) error
	ExpectedError     string
}

func (tc *prometheusAPITestCase) RunTest(t *testing.T) {
	mockAPI := &mocks.PrometheusAPI{}
	mockAPI.On(tc.ApiMethod, tc.ApiExpectedParams...).Return(tc.ApiExpectedReturn...)
	c := echo.New().NewContext(httptest.NewRequest(http.MethodGet, tc.RequestURL, nil), httptest.NewRecorder())
	c.SetParamNames(tc.ParamNames...)
	c.SetParamValues(tc.ParamValues...)
	err := tc.FuncToTest(mockAPI)(c)
	if tc.ExpectedError != "" {
		assert.EqualError(t, err, tc.ExpectedError)
	} else {
		assert.NoError(t, err)
	}
}

func TestGetPrometheusTargetsMetadata(t *testing.T) {
	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com",
		ParamNames:        []string{"network_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "TargetsMetadata",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{[]v1.MetricMetadata{}, nil},
		FuncToTest:        GetPrometheusTargetsMetadata,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{[]v1.MetricMetadata{}, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetPrometheusQueryHandler(t *testing.T) {
	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com?query=up",
		ParamNames:        []string{"network_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "Query",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{model.Vector{}, nil, nil},
		FuncToTest:        GetPrometheusQueryHandler,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "invalid query",
			RequestURL:    "http://url.com?query=!test",
			ExpectedError: "code=500, message=error parsing query: parse error at char 1: unexpected character after '!': 't'",
		},
		{
			Name:          "invalid time",
			RequestURL:    "http://url.com?query=up&time=abc",
			ExpectedError: `code=400, message=unable to parse time parameter: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{model.Vector{}, nil, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetTenantQueryHandler(t *testing.T) {
	tenants_test_init.StartTestService(t)
	tenants.CreateTenant(0, &protos.Tenant{
		Name:     "0",
		Networks: []string{"test"},
	})

	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com?query=up",
		ParamNames:        []string{"tenant_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "Query",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{model.Vector{}, nil, nil},
		FuncToTest:        GetTenantQueryHandler,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no tenant",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Tenant ID",
		},
		{
			Name:          "invalid tenant",
			ParamValues:   []string{"99"},
			ExpectedError: "Not found",
		},
		{
			Name:          "invalid query",
			RequestURL:    "http://url.com?query=!test",
			ExpectedError: "code=500, message=error parsing query: parse error at char 1: unexpected character after '!': 't'",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetPrometheusQueryRangeHandler(t *testing.T) {
	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com?query=up&start=0",
		ParamNames:        []string{"network_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "QueryRange",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{model.Matrix{}, nil, nil},
		FuncToTest:        GetPrometheusQueryRangeHandler,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no network",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Network ID",
		},
		{
			Name:          "invalid query",
			RequestURL:    "http://url.com?query=!test&start=0",
			ExpectedError: "code=500, message=error parsing query: parse error at char 1: unexpected character after '!': 't'",
		},
		{
			Name:          "invalid start",
			RequestURL:    "http://url.com?query=up&start=abc",
			ExpectedError: `code=400, message=unable to parse start parameter: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:          "invalid end",
			RequestURL:    "http://url.com?query=up&start=0&end=abc",
			ExpectedError: `code=400, message=unable to parse end parameter: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:          "invalid step",
			RequestURL:    "http://url.com?query=up&start=0&step=abc",
			ExpectedError: `code=400, message=unable to parse step parameter: time: invalid duration abc`,
		},
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{model.Matrix{}, nil, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetTenantQueryRangeHandler(t *testing.T) {
	tenants_test_init.StartTestService(t)
	tenants.CreateTenant(0, &protos.Tenant{
		Name:     "0",
		Networks: []string{"test"},
	})

	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com?query=up&start=0",
		ParamNames:        []string{"tenant_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "QueryRange",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{model.Matrix{}, nil, nil},
		FuncToTest:        GetTenantPromQueryRangeHandler,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no tenant",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=Missing Tenant ID",
		},
		{
			Name:          "invalid tenant",
			ParamValues:   []string{"99"},
			ExpectedError: "code=500, message=Not found",
		},
		{
			Name:          "invalid query",
			RequestURL:    "http://url.com?query=!test&start=0",
			ExpectedError: "code=500, message=error parsing query: parse error at char 1: unexpected character after '!': 't'",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetPrometheusSeriesHandler(t *testing.T) {
	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com",
		ParamNames:        []string{"network_id"},
		ParamValues:       []string{"0"},
		ApiMethod:         "Series",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, nil},
		FuncToTest:        GetPrometheusSeriesHandler,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no networkID",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=code=400, message=Missing Network ID",
		},
		{
			Name:          "invalid start",
			RequestURL:    "http://url.com?start=abc",
			ExpectedError: `code=400, message=unable to parse start parameter: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:          "invalid end",
			RequestURL:    "http://url.com?end=abc",
			ExpectedError: `code=400, message=unable to parse end parameter: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:          "invalid matchers",
			RequestURL:    "http://url.com?match={!invalid}",
			ExpectedError: "code=400, message=Error parsing series matchers: code=500, message=unable to secure match parameter: error parsing query: parse error at char 2: unexpected character after '!' inside braces: 'i'",
		},
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetTenantPromSeriesHandler(t *testing.T) {
	tenants_test_init.StartTestService(t)
	tenants.CreateTenant(0, &protos.Tenant{
		Name:     "0",
		Networks: []string{"test"},
	})

	mockAPI := &mocks.PrometheusAPI{}
	mockAPI.On("Series", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.LabelSet{}, nil, nil)
	req := httptest.NewRequest(http.MethodGet, "http://url.com", nil)
	rec := httptest.NewRecorder()

	// Successful Request
	c := echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err := GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.NoError(t, err)

	// No tenant ID
	c = echo.New().NewContext(req, rec)
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, "code=400, message=code=400, message=Missing Tenant ID")

	// Invalid Tenant ID
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("1")
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, "code=500, message=Not found")

	// Invalid start
	req = httptest.NewRequest(http.MethodGet, "http://url.com?start=abc", nil)
	rec = httptest.NewRecorder()
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, `code=400, message=parse start time: abc: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`)

	// Invalid end
	req = httptest.NewRequest(http.MethodGet, "http://url.com?end=abc", nil)
	rec = httptest.NewRecorder()
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, `code=400, message=parse end time: abc: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`)

	// Invalid Query
	req = httptest.NewRequest(http.MethodGet, "http://url.com?match[]={!invalidQuery}", nil)
	rec = httptest.NewRecorder()
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, "code=400, message=Error parsing series matchers: code=500, message=unable to secure match parameter: error parsing query: parse error at char 2: unexpected character after '!' inside braces: 'i'")

	// Prometheus error
	req = httptest.NewRequest(http.MethodGet, "http://url.com", nil)
	rec = httptest.NewRecorder()
	mockAPI = &mocks.PrometheusAPI{}
	mockAPI.On("Series", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.LabelSet{}, nil, errors.New("prometheus error"))
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err = GetTenantPromSeriesHandler(mockAPI, false)(c)
	assert.EqualError(t, err, "code=500, message=prometheus error")

	// Use cache
	mockAPI = &mocks.PrometheusAPI{}
	mockAPI.On("Series", mock.Anything, mock.Anything, mock.Anything, mock.Anything).Return([]model.LabelSet{}, nil, nil)
	c = echo.New().NewContext(req, rec)
	c.SetParamNames("tenant_id")
	c.SetParamValues("0")
	err = GetTenantPromSeriesHandler(mockAPI, true)(c)
	assert.NoError(t, err)
}

func TestGetTenantPromValuesHandler(t *testing.T) {
	tenants_test_init.StartTestService(t)
	tenants.CreateTenant(0, &protos.Tenant{
		Name:     "0",
		Networks: []string{"test"},
	})

	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com",
		ParamNames:        []string{"tenant_id", "label_name"},
		ParamValues:       []string{"0", "test"},
		ApiMethod:         "Series",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, nil},
		FuncToTest:        GetTenantPromValuesHandler,
	}

	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "no tenant",
			ParamNames:    []string{"label_name"},
			ExpectedError: "code=400, message=code=400, message=Missing Tenant ID",
		},
		{
			Name:          "no label name",
			ParamNames:    []string{"tenant_id"},
			ExpectedError: "code=400, message=label_name is required",
		},
		{
			Name:          "invalid tenant",
			ParamValues:   []string{"99", "test"},
			ExpectedError: "code=500, message=Not found",
		},
		{
			Name:          "invalid label",
			ParamValues:   []string{"0", "!test"},
			ExpectedError: "code=500, message=error parsing query: parse error at char 2: unexpected character after '!' inside braces: 't'",
		},
		{
			Name:          "invalid start",
			RequestURL:    "http://url.com?start=abc",
			ExpectedError: `code=400, message=parse start time: abc: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:          "invalid end",
			RequestURL:    "http://url.com?end=abc",
			ExpectedError: `code=400, message=parse end time: abc: parsing time "abc" as "2006-01-02T15:04:05Z07:00": cannot parse "abc" as "2006"`,
		},
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

type seriesMatchesTestCase struct {
	name            string
	inputURL        string
	start           string
	end             string
	restrictor      restrictor.QueryRestrictor
	expectedStrings []string
	expectedError   string
}

func (tc *seriesMatchesTestCase) RunTest(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, tc.inputURL, nil)
	rec := httptest.NewRecorder()

	c := echo.New().NewContext(req, rec)
	params, err := getSeriesMatches(c, "match", tc.restrictor)
	if tc.expectedError == "" {
		assert.NoError(t, err)
		assert.Equal(t, tc.expectedStrings, params)
	} else {
		assert.EqualError(t, err, tc.expectedError)
	}
}

func TestGetSeriesMatches(t *testing.T) {
	testCases := []seriesMatchesTestCase{
		{
			name:            "single match",
			inputURL:        "/?match=up",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`},
		},
		{
			name:            "two match",
			inputURL:        "/?match=up%20down",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`},
		},
		{
			name:            "complicated match",
			inputURL:        "/?match=up%20down%20{gatewayID=\"gw1\"}",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`up{networkID="test"}`, `down{networkID="test"}`, `{gatewayID="gw1",networkID="test"}`},
		},
		{
			name:            "no match",
			inputURL:        "/",
			restrictor:      networkQueryRestrictorProvider("test"),
			expectedStrings: []string{`{networkID="test"}`},
		},
		{
			name:            "tenant match",
			inputURL:        "/",
			restrictor:      *restrictor.NewQueryRestrictor(restrictor.Opts{ReplaceExistingLabel: false}).AddMatcher("networkID", "net1", "net2"),
			expectedStrings: []string{`{networkID=~"net1|net2"}`},
		},
		{
			name:            "tenant two match",
			inputURL:        "/?match=up%20down",
			restrictor:      *restrictor.NewQueryRestrictor(restrictor.Opts{ReplaceExistingLabel: false}).AddMatcher("networkID", "net1", "net2"),
			expectedStrings: []string{`up{networkID=~"net1|net2"}`, `down{networkID=~"net1|net2"}`},
		},
		{
			name:          "invalid query",
			inputURL:      "/?match={!invalidQuery}",
			restrictor:    networkQueryRestrictorProvider("test"),
			expectedError: "code=500, message=unable to secure match parameter: error parsing query: parse error at char 2: unexpected character after '!' inside braces: 'i'",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, tc.RunTest)
	}
}

func TestTenantSeriesHandlerProvider(t *testing.T) {
	tenants_test_init.StartTestService(t)
	tenants.CreateTenant(0, &protos.Tenant{
		Name:     "0",
		Networks: []string{"test"},
	})

	baseTest := prometheusAPITestCase{
		Name:              "successful request",
		RequestURL:        "http://url.com",
		ParamNames:        []string{"tenant_id", "label_name"},
		ParamValues:       []string{"0", "test"},
		ApiMethod:         "Series",
		ApiExpectedParams: []interface{}{mock.Anything, mock.Anything, mock.Anything, mock.Anything},
		ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, nil},
		FuncToTest:        TenantSeriesHandlerProvider,
	}
	tests := []prometheusAPITestCase{
		baseTest,
		{
			Name:          "invalid tenant",
			ParamValues:   []string{"99"},
			ExpectedError: "code=500, message=Not found",
		},
		{
			Name:          "no tenant",
			ParamNames:    []string{"nil"},
			ExpectedError: "code=400, message=code=400, message=Missing Tenant ID",
		},
		{
			Name:              "prometheus error",
			ApiExpectedReturn: []interface{}{[]model.LabelSet{}, nil, errors.New("prometheus error")},
			ExpectedError:     "code=500, message=prometheus error",
		},
		{
			Name:          "invalid match",
			RequestURL:    "http://url.com?match={!invalidMatch}",
			ExpectedError: "code=400, message=Error parsing series matchers: code=500, message=unable to secure match parameter: error parsing query: parse error at char 2: unexpected character after '!' inside braces: 'i'",
		},
	}
	for i := range tests {
		mergo.Merge(&tests[i], baseTest)
	}
	for _, tc := range tests {
		t.Run(tc.Name, tc.RunTest)
	}
}

func TestGetSetOfValuesFromLabel(t *testing.T) {
	seriesList := []model.LabelSet{{"__name__": "test", "label1": "val1"}, {"__name__": "test2", "label1": "val2"}, {"__name__": "test"}}
	vals := getSetOfValuesFromLabel(seriesList, "__name__")

	sort.Strings(vals)
	assert.Equal(t, []string{"test", "test2"}, vals)
}
