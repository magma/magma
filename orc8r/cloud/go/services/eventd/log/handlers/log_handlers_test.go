package handlers

import (
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo"
	"github.com/olivere/elastic/v7"
	"github.com/stretchr/testify/assert"
)

var (
	elasticCases = []elasticTestCase{
		{
			name: "full query params",
			params: logQueryParams{
				SimpleQuery: "magmad",
				Fields:      []string{"message"},
				Filters:     map[string]string{"test": "result"},
				StartTime:   "2019-09-27T10:59:56.248Z",
				EndTime:     "2019-09-27T11:14:56.248Z",
			},
		},
		{
			name: "no time",
			params: logQueryParams{
				SimpleQuery: "magmad",
				Fields:      []string{"message"},
			},
		},
		{
			name: "no query",
			params: logQueryParams{
				StartTime: "2019-09-27T10:59:56.248Z",
				EndTime:   "2019-09-27T11:14:56.248Z",
			},
		},
		{
			name:   "no params",
			params: logQueryParams{},
		},
	}

	queryParamsTestCases = []queryParamTestCase{
		{
			name:      "no params",
			urlString: "",
			expectedParams: logQueryParams{
				Fields: []string{},
				From:   0,
				Size:   defaultSearchSize,
			},
		},
		{
			name:      "only string params",
			urlString: "?size=123&start=12345&end=23456",
			expectedParams: logQueryParams{
				Size:      123,
				From:      0,
				StartTime: "12345",
				EndTime:   "23456",
				Fields:    []string{},
			},
		},
		{
			name:      "string, list and map params",
			urlString: "?from=10&size=123&start=12345&end=23456&fields=message&filters=test:result",
			expectedParams: logQueryParams{
				Size:      123,
				From:      10,
				StartTime: "12345",
				EndTime:   "23456",
				Fields:    []string{"message"},
				Filters:   map[string]string{"test": "result"},
			},
		},
		{
			name:      "multiple values",
			urlString: "?fields=message,log,test&filters=test:result,key:value,foo:baz",
			expectedParams: logQueryParams{
				From:    0,
				Size:    defaultSearchSize,
				Fields:  []string{"message", "log", "test"},
				Filters: map[string]string{"test": "result", "key": "value", "foo": "baz"},
			},
		},
	}

	countParamsTestCases = []countParamTestCase{
		{
			name:      "no params",
			urlString: "",
			expectedParams: logQueryParams{
				Fields: []string{},
			},
		},
		{
			name:      "string, list and map params",
			urlString: "?start=12345&end=23456&fields=message&filters=test:result",
			expectedParams: logQueryParams{
				StartTime: "12345",
				EndTime:   "23456",
				Fields:    []string{"message"},
				Filters:   map[string]string{"test": "result"},
			},
		},
	}
)

// TestToElasticBoolQuery tests that various combinations of query parameters
// results in a valid and expected ElasticSearch Query
func TestToElasticBoolQuery(t *testing.T) {
	for _, test := range elasticCases {
		t.Run(test.name, func(t *testing.T) {
			runToElasticBoolQueryTestCase(t, test)
		})
	}
}

// TestSecureQuery tests that all queries include filters for network_id and
// gateway_id to enable multi-tenancy
func TestSecureQuery(t *testing.T) {
	for _, test := range elasticCases {
		t.Run(test.name, func(t *testing.T) {
			runSecureElasticQueryTestCase(t, test)
		})
	}
}

// TestGetQueryParams tests that parameters in the url are parsed correctly
func TestGetQueryParams(t *testing.T) {
	for _, test := range queryParamsTestCases {
		t.Run(test.name, func(t *testing.T) {
			runQueryParamTestCase(t, test)
		})
	}
}

// TestGetCountParams tests that parameters in the url are parsed correctly
func TestGetCountParams(t *testing.T) {
	for _, test := range countParamsTestCases {
		t.Run(test.name, func(t *testing.T) {
			runCountParamTestCase(t, test)
		})
	}
}

type elasticTestCase struct {
	name     string
	params   logQueryParams
	expected elastic.BoolQuery
}

func runToElasticBoolQueryTestCase(t *testing.T, tc elasticTestCase) {
	query := tc.params.ToElasticBoolQuery()
	source, err := query.Source()
	assert.NoError(t, err)

	s := source.(map[string]interface{})

	boolQuery, ok := s["bool"].(map[string]interface{})
	assert.True(t, ok)

	if tc.params.EndTime != "" || tc.params.StartTime != "" {
		must, ok := boolQuery["must"].(map[string]interface{})
		assert.True(t, ok)
		assert.Len(t, must, 1)

		_, ok = must["range"]
		assert.True(t, ok)
	} else {
		_, ok := boolQuery["must"]
		assert.False(t, ok)
	}

	// Check simple_string_query
	if tc.params.SimpleQuery != "" {
		f, ok := boolQuery["filter"]
		assert.True(t, ok)
		simpleQueryExists := false
		if filters, ok := f.([]interface{}); ok {
			for _, filter := range filters {
				foundQuery, isSimple := filter.(map[string]interface{})["simple_query_string"].(map[string]interface{})
				if isSimple && tc.params.SimpleQuery == foundQuery["query"] {
					simpleQueryExists = true
				}
			}
		} else {
			filters := f.(map[string]interface{})
			simple, ok := filters["simple_query_string"].(map[string]interface{})
			assert.True(t, ok)
			if simple["query"] == tc.params.SimpleQuery {
				simpleQueryExists = true
			}
		}
		assert.True(t, simpleQueryExists)
	}
	// Check that filters are applied as expected
	if len(tc.params.Filters) > 0 {
		f, ok := boolQuery["filter"]
		assert.True(t, ok)
		if filters, ok := f.([]interface{}); ok {
			for expectedKey, expectedVal := range tc.params.Filters {
				filterExists := false
				for _, filter := range filters {
					foundTerm, isTerm := filter.(map[string]interface{})["term"].(map[string]interface{})
					if isTerm && foundTerm[expectedKey] == expectedVal {
						filterExists = true
					}
				}
				assert.True(t, filterExists)
			}

		} else {
			filters := f.(map[string]interface{})
			for _, filter := range filters {
				for expectedKey, expectedValue := range tc.params.Filters {
					assert.Equal(t, expectedValue, filter.(map[string]interface{})[expectedKey])
				}
			}
		}
	}
}

func runSecureElasticQueryTestCase(t *testing.T, tc elasticTestCase) {
	networkID := "testNetwork"
	secureQuery := secureElasticQuery(networkID, tc.params)

	source, err := secureQuery.Source()
	assert.NoError(t, err)

	f, ok := source.(map[string]interface{})["bool"].(map[string]interface{})["filter"]
	assert.True(t, ok)

	if filters, ok := f.([]interface{}); ok {
		networkFound := false
		for _, filter := range filters {
			foundTerm, ok := filter.(map[string]interface{})["term"].(map[string]interface{})
			if !ok {
				continue
			}
			if foundTerm[NetworkLogLabel] == networkID {
				networkFound = true
				continue
			}
		}
		assert.True(t, networkFound)
	} else {
		filters := f.(map[string]interface{})
		for _, val := range filters {
			assert.Equal(t, networkID, val.(map[string]interface{})[NetworkLogLabel])
		}
	}
}

type queryParamTestCase struct {
	name           string
	urlString      string
	expectedParams logQueryParams
}

func runQueryParamTestCase(t *testing.T, tc queryParamTestCase) {
	req := httptest.NewRequest(echo.GET, fmt.Sprintf("/%s", tc.urlString), nil)
	c := echo.New().NewContext(req, httptest.NewRecorder())
	params, err := getQueryParameters(c)
	assert.NoError(t, err)
	assert.Equal(t, tc.expectedParams, params)
}

type countParamTestCase struct {
	name           string
	urlString      string
	expectedParams logQueryParams
}

func runCountParamTestCase(t *testing.T, tc countParamTestCase) {
	req := httptest.NewRequest(echo.GET, fmt.Sprintf("/%s", tc.urlString), nil)
	c := echo.New().NewContext(req, httptest.NewRecorder())
	params, err := getCountParameters(c)
	assert.NoError(t, err)
	assert.Equal(t, tc.expectedParams, params)
}
