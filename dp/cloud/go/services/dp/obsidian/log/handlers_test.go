/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package dp_log

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-openapi/strfmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"magma/dp/cloud/go/dp"
	dp_service "magma/dp/cloud/go/services/dp"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/cloud/go/services/obsidian/tests"
	"magma/orc8r/cloud/go/test_utils"
)

const (
	someTimeStamp                 = "2022-01-14T10:23:49.000Z"
	successfulElasticSearchResult = `{
   "took":982,
   "timed_out":false,
   "_shards":{
      "total":13,
      "successful":13,
      "skipped":0,
      "failed":0
   },
   "hits":{
      "total":{
         "value":10000,
         "relation":"gte"
      },
      "max_score":1.0,
      "hits":[
         {
            "_index":"dp-1234",
            "_type":"_doc",
            "_id":"2ds34f6w-43f5-2344-dsf4-kf9ekw9fke9w",
            "_score":1.0,
            "_source":{
               "event_timestamp": 1642155829,
               "log_message":"some message1",
               "fcc_id":"some_fcc_id",
               "log_from":"SAS",
               "cbsd_serial_number":"some_serial_number",
               "@timestamp":"2022-01-14T10:23:49.871Z",
               "log_to":"DP",
               "log_name":"grantResponse"
            }
         },
         {
            "_index":"dp-1234",
            "_type":"_doc",
            "_id":"2ds34f6w-43f5-2344-dsf4-kf9ekw9fke9w",
            "_score":1.0,
            "_source":{
               "event_timestamp": 1642155829,
               "log_message":"some message2",
               "fcc_id":"some_fcc_id",
               "log_from":"SAS",
               "cbsd_serial_number":"some_serial_number",
               "@timestamp":"2022-01-14T10:23:49.871Z",
               "log_to":"DP",
               "log_name":"grantResponse"
            }
         }
      ]
   }
}`
)

type stubElasticSearchSearver struct {
	expectedPayload string
	response        string
	t               *testing.T
}

func (s *stubElasticSearchSearver) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p, _ := ioutil.ReadAll(r.Body)
	assert.JSONEq(s.t, s.expectedPayload, string(p))
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(s.response))
	assert.Equal(s.t, "/"+wildcardedIndex+"/_search", r.URL.Path)
}

func TestHandlers(t *testing.T) {
	suite.Run(t, &HandlersTestSuite{})
}

type HandlersTestSuite struct {
	suite.Suite
}

func (s *HandlersTestSuite) SetupTest() {
	srv, lis, _ := test_utils.NewTestService(s.T(), dp.ModuleName, dp_service.ServiceName)
	go srv.RunTest(lis, nil)
}

func (s *HandlersTestSuite) TestBoolQuery() {
	testCases := []struct {
		name            string
		params          ListLogsRequest
		expectedPayload json.RawMessage
	}{
		{
			name: "test query with no params",
			params: ListLogsRequest{
				Index:     "dp",
				NetworkId: "someNetworkId",
				Filter:    &LogsFilter{},
			},
			expectedPayload: []byte(`{
   "must":
      {
         "match":{
            "network_id":{
               "query":"someNetworkId"
            }
         }
      }
   }`),
		},
		{
			name: "test query with some params",
			params: ListLogsRequest{
				Index:     "dp",
				NetworkId: "someNetworkId",
				Filter: &LogsFilter{
					LogTo:           "SAS",
					Name:            "someLog",
					SerialNumber:    "someSerialNumber123",
					EndTimestampSec: to_pointer.Int64(1000),
				},
			},
			expectedPayload: []byte(`{
   "must":[
      {
         "range":{
            "event_timestamp":{
               "from":null,
               "include_lower":true,
               "include_upper":true,
               "to":1000
            }
         }
      },
      {
         "match":{
            "network_id":{
               "query":"someNetworkId"
            }
         }
      },
      {
         "match":{
            "log_to":{
               "query":"SAS"
            }
         }
      },
      {
         "match":{
            "log_name":{
               "query":"someLog"
            }
         }
      },
      {
         "match":{
            "cbsd_serial_number":{
               "query":"someSerialNumber123"
            }
         }
      }
   ]
}`),
		},
		{
			name: "test query with all params",
			params: ListLogsRequest{
				Index:     "dp",
				NetworkId: "someNetworkId",
				Filter: &LogsFilter{
					LogFrom:           "DP",
					LogTo:             "SAS",
					Name:              "someLog",
					SerialNumber:      "someSerialNumber123",
					FccId:             "someFccId123",
					ResponseCode:      to_pointer.Int64(0),
					BeginTimestampSec: to_pointer.Int64(100),
					EndTimestampSec:   to_pointer.Int64(1000),
				},
			},
			expectedPayload: []byte(`{
   "must":[
      {
         "range":{
            "event_timestamp":{
               "from":100,
               "include_lower":true,
               "include_upper":true,
               "to":1000
            }
         }
      },
	  {
         "match":{
            "network_id":{
               "query":"someNetworkId"
            }
         }
      },
      {
         "match":{
            "log_from":{
               "query":"DP"
            }
         }
      },
      {
         "match":{
            "log_to":{
               "query":"SAS"
            }
         }
      },
      {
         "match":{
            "log_name":{
               "query":"someLog"
            }
         }
      },
      {
         "match":{
            "cbsd_serial_number":{
               "query":"someSerialNumber123"
            }
         }
      },
      {
         "match":{
            "fcc_id":{
               "query":"someFccId123"
            }
         }
      },
      {
         "match":{
            "response_code":{
               "query":0
            }
         }
      }
   ]
}`),
		},
	}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			t := s.T()
			query := tc.params.toElasticSearchQuery()
			source, err := query.Source()
			assert.NoError(t, err)
			s := source.(map[string]interface{})
			boolQuery, _ := s["bool"].(map[string]interface{})
			expected := new(bytes.Buffer)
			if err := json.Compact(expected, tc.expectedPayload); err != nil {
				panic(err)
			}
			actual, _ := json.Marshal(boolQuery)
			assert.JSONEq(t, expected.String(), string(actual))
		})
	}
}

func (s *HandlersTestSuite) TestGetPagination() {
	testCases := []struct {
		testName       string
		URL            string
		expectedLimit  *int
		expectedOffset *int
	}{
		{
			testName:       "test pagination with limit and offset",
			URL:            "/some/url?limit=1&offset=2",
			expectedLimit:  to_pointer.Int(1),
			expectedOffset: to_pointer.Int(2),
		},
		{
			testName:       "test pagination with limit only",
			URL:            "/some/url?limit=1",
			expectedLimit:  to_pointer.Int(1),
			expectedOffset: nil,
		},
		{
			testName:       "test pagination without limit and offset",
			URL:            "/some/url",
			expectedLimit:  nil,
			expectedOffset: nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			t := tests.Test{}
			req := *httptest.NewRequest(t.Method, tc.URL, bytes.NewReader(nil))
			c := echo.New().NewContext(&req, httptest.NewRecorder())
			pag, _ := getPagination(c)
			if tc.expectedLimit != nil {
				s.Equal(tc.expectedLimit, pag.Size)
			} else {
				s.Nil(pag.Size)
			}
			if tc.expectedOffset != nil {
				s.Equal(tc.expectedOffset, pag.From)
			} else {
				s.Nil(pag.From)
			}
		})
	}
}

func (s *HandlersTestSuite) TestListLogs() {
	testCases := []struct {
		testName              string
		paramNames            []string
		ParamValues           []string
		model                 db.Model
		expectedStatus        int
		expectedResult        *models.PaginatedLogs
		expectedError         string
		queryParamsString     url.Values
		elasticSearchResponse string
		expectedPayload       string
	}{
		{
			testName:              "test list logs without query params",
			paramNames:            []string{"network_id"},
			ParamValues:           []string{"n1"},
			expectedStatus:        http.StatusOK,
			expectedResult:        getPaginatedLogs(),
			elasticSearchResponse: successfulElasticSearchResult,
			expectedPayload: `{
   "query":{
      "bool":{
         "must":{
		    "match":{
			   "network_id":{
				  "query":"n1"
			   }
		    }
		 }
      }
   },
   "sort":[
      {
         "event_timestamp":{
            "order":"desc"
         }
      }
   ]
}`,
		},
		{
			testName:              "test list logs with limit and offset",
			paramNames:            []string{"network_id"},
			ParamValues:           []string{"n1"},
			expectedStatus:        http.StatusOK,
			expectedResult:        getPaginatedLogs(),
			queryParamsString:     url.Values{"limit": {"4"}, "offset": {"3"}},
			elasticSearchResponse: successfulElasticSearchResult,
			expectedPayload: `{
   "from":3,
   "query":{
      "bool":{
         "must":{
            "match":{
               "network_id":{
                  "query":"n1"
               }
            }
         }
      }
   },
   "size":4,
   "sort":[
      {
         "event_timestamp":{
            "order":"desc"
         }
      }
   ]
}`,
		},
		{
			testName:       "test list logs with all query params",
			paramNames:     []string{"network_id"},
			ParamValues:    []string{"n1"},
			expectedStatus: http.StatusOK,
			expectedError:  "",
			queryParamsString: url.Values{
				"limit":         {"4"},
				"offset":        {"3"},
				"from":          {"SAS"},
				"to":            {"DP"},
				"type":          {"grantResponse"},
				"serial_number": {"some_serial_number"},
				"fcc_id":        {"some_fcc_id"},
				"response_code": {"0"},
				"begin":         {"2022-01-14T10:23:49.871036Z"},
			},
			expectedResult:        getPaginatedLogs(),
			elasticSearchResponse: successfulElasticSearchResult,
			expectedPayload: `{
   "from":3,
   "query":{
      "bool":{
         "must":[
            {
               "range":{
                  "event_timestamp":{
                     "from":1642155829,
                     "include_lower":true,
                     "include_upper":true,
                     "to":null
                  }
               }
            },
            {
               "match":{
                  "network_id":{
                     "query":"n1"
                  }
               }
            },
            {
               "match":{
                  "log_from":{
                     "query":"SAS"
                  }
               }
            },
            {
               "match":{
                  "log_to":{
                     "query":"DP"
                  }
               }
            },
            {
               "match":{
                  "log_name":{
                     "query":"grantResponse"
                  }
               }
            },
            {
               "match":{
                  "cbsd_serial_number":{
                     "query":"some_serial_number"
                  }
               }
            },
            {
               "match":{
                  "fcc_id":{
                     "query":"some_fcc_id"
                  }
               }
            },
            {
               "match":{
                  "response_code":{
                     "query":0
                  }
               }
            }
         ]
      }
   },
   "size":4,
   "sort":[
      {
         "event_timestamp":{
            "order":"desc"
         }
      }
   ]
}`,
		},
		{
			testName:          "test list logs with incorrect limit value",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'incorrect_limit_value' is not a proper value for limit",
			queryParamsString: url.Values{"limit": {"incorrect_limit_value"}},
		},
		{
			testName:          "test list logs with incorrect offset value",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'incorrect_offset_value' is not a proper value for offset",
			queryParamsString: url.Values{"offset": {"incorrect_offset_value"}},
		},
		{
			testName:          "test list logs with offset and no limit",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "offset requires a limit",
			queryParamsString: url.Values{"offset": {"1"}},
		},
		{
			testName:          "test list logs with incorrect begin date",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'2022-01-14T10' is not a proper value for begin",
			queryParamsString: url.Values{"begin": {"2022-01-14T10"}},
		},
		{
			testName:          "test list logs with incorrect end date",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'2022-01-14T10' is not a proper value for end",
			queryParamsString: url.Values{"end": {"2022-01-14T10"}},
		},
		{
			testName:          "test list logs with incorrect response code",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'foo' is not a proper value for response_code",
			queryParamsString: url.Values{"response_code": {"foo"}},
		},
	}
	e := echo.New()

	for _, t := range testCases {
		s.Run(t.testName, func() {
			testServer := httptest.NewServer(&stubElasticSearchSearver{
				expectedPayload: t.expectedPayload,
				response:        t.elasticSearchResponse,
				t:               s.T(),
			})
			defer testServer.Close()
			obsidianHandlers := NewHandlersGetter(GetElasticClient, testServer.URL).GetHandlers()
			listLogs := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, ManageLogsPath, obsidian.GET).HandlerFunc
			tc := tests.Test{
				Method:         http.MethodGet,
				URL:            ManageLogsPath + "?" + t.queryParamsString.Encode(),
				Payload:        nil,
				ParamNames:     t.paramNames,
				ParamValues:    t.ParamValues,
				Handler:        listLogs,
				ExpectedStatus: t.expectedStatus,
				ExpectedResult: tests.JSONMarshaler(t.expectedResult),
				ExpectedError:  t.expectedError,
			}
			tests.RunUnitTest(s.T(), e, tc)
		})
	}
}

func (s *HandlersTestSuite) TestGetLogsFilter() {
	testCases := []struct {
		testName             string
		URL                  string
		expectedLogFrom      string
		expectedLogTo        string
		expectedName         string
		expectedSerialNumber string
		expectedFccId        string
		expectedResponseCode *int64
		expectedBeginTS      *int64
		expectedEndTS        *int64
		expectedError        string
	}{
		{
			testName:             "test filter with all params",
			URL:                  "/some/url?from=SAS&to=DP&type=grantResponse&serial_number=some_serial_number&fcc_id=some_fcc_id&response_code=0&begin=2022-01-14T10%3A23%3A49.871036Z&end=2022-01-15T10%3A23%3A49.871036Z",
			expectedLogFrom:      "SAS",
			expectedLogTo:        "DP",
			expectedName:         "grantResponse",
			expectedSerialNumber: "some_serial_number",
			expectedFccId:        "some_fcc_id",
			expectedResponseCode: to_pointer.Int64(0),
			expectedBeginTS:      to_pointer.Int64(1642155829),
			expectedEndTS:        to_pointer.Int64(1642242229),
		},
		{
			testName:             "test filter without params",
			URL:                  "/some/url",
			expectedLogFrom:      "",
			expectedLogTo:        "",
			expectedName:         "",
			expectedSerialNumber: "",
			expectedFccId:        "",
			expectedResponseCode: nil,
			expectedBeginTS:      nil,
			expectedEndTS:        nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			t := tests.Test{}
			req := *httptest.NewRequest(t.Method, tc.URL, bytes.NewReader(nil))
			c := echo.New().NewContext(&req, httptest.NewRecorder())
			f, err := getLogsFilter(c)
			if err != nil {
				s.Equal(tc.expectedError, err.Error())
			} else {
				s.Equal(tc.expectedLogFrom, f.LogFrom)
				s.Equal(tc.expectedLogTo, f.LogTo)
				s.Equal(tc.expectedName, f.Name)
				s.Equal(tc.expectedSerialNumber, f.SerialNumber)
				s.Equal(tc.expectedFccId, f.FccId)
				if tc.expectedResponseCode != nil {
					s.Equal(tc.expectedResponseCode, f.ResponseCode)
				} else {
					s.Nil(tc.expectedResponseCode, f.ResponseCode)
				}
				if tc.expectedBeginTS != nil {
					s.Equal(tc.expectedBeginTS, f.BeginTimestampSec)
				} else {
					s.Nil(tc.expectedBeginTS, f.BeginTimestampSec)
				}
				if tc.expectedEndTS != nil {
					s.Equal(tc.expectedEndTS, f.EndTimestampSec)
				} else {
					s.Nil(tc.expectedEndTS, f.EndTimestampSec)
				}
			}
		})
	}
}

func getTestDateTime(s string) strfmt.DateTime {
	dt, _ := time.Parse(time.RFC3339, s)
	return strfmt.DateTime(dt)
}

func getPaginatedLogs() *models.PaginatedLogs {
	return &models.PaginatedLogs{
		Logs: []*models.Log{
			{
				Body:         "some message1",
				FccID:        "some_fcc_id",
				From:         "SAS",
				SerialNumber: "some_serial_number",
				Time:         getTestDateTime(someTimeStamp),
				To:           "DP",
				Type:         "grantResponse",
			},
			{
				Body:         "some message2",
				FccID:        "some_fcc_id",
				From:         "SAS",
				SerialNumber: "some_serial_number",
				Time:         getTestDateTime(someTimeStamp),
				To:           "DP",
				Type:         "grantResponse",
			},
		},
		TotalCount: 10000,
	}
}
