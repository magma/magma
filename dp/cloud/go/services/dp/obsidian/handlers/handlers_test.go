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

package handlers_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/dp"
	"magma/dp/cloud/go/protos"
	dp_service "magma/dp/cloud/go/services/dp"
	"magma/dp/cloud/go/services/dp/obsidian/handlers"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/obsidian/tests"
	"magma/orc8r/cloud/go/test_utils"
)

func TestHandlers(t *testing.T) {
	suite.Run(t, &HandlersTestSuite{})
}

type HandlersTestSuite struct {
	suite.Suite
	cbsdServer *stubCbsdServer
	logsServer *stubLogsServer
}

func (s *HandlersTestSuite) SetupTest() {
	s.cbsdServer = &stubCbsdServer{}
	s.cbsdServer.t = s.T()
	s.logsServer = &stubLogsServer{}
	s.logsServer.t = s.T()
	srv, lis := test_utils.NewTestService(s.T(), dp.ModuleName, dp_service.ServiceName)
	protos.RegisterCbsdManagementServer(srv.GrpcServer, s.cbsdServer)
	protos.RegisterLogFetcherServer(srv.GrpcServer, s.logsServer)
	go srv.RunTest(lis)
}

type stubCbsdServer struct {
	protos.UnimplementedCbsdManagementServer
	expectedListRequest   *protos.ListCbsdRequest
	listResponse          *protos.ListCbsdResponse
	expectedFetchRequest  *protos.FetchCbsdRequest
	fetchResponse         *protos.FetchCbsdResponse
	expectedCreateRequest *protos.CreateCbsdRequest
	createResponse        *protos.CreateCbsdResponse
	expectedUpdateRequest *protos.UpdateCbsdRequest
	updateResponse        *protos.UpdateCbsdResponse
	expectedDeleteRequest *protos.DeleteCbsdRequest
	deleteResponse        *protos.DeleteCbsdResponse
	err                   error
	t                     *testing.T
}

type stubLogsServer struct {
	protos.UnimplementedLogFetcherServer
	expectedListRequest *protos.ListLogsRequest
	listResponse        *protos.ListLogsResponse
	err                 error
	t                   *testing.T
}

func (s *HandlersTestSuite) TestListLogs() {
	testCases := []struct {
		testName            string
		paramNames          []string
		ParamValues         []string
		model               db.Model
		expectedStatus      int
		expectedResult      []*models.Message
		expectedError       string
		queryParamsString   url.Values
		expectedListRequest *protos.ListLogsRequest
	}{
		{
			testName:       "test list logs without query params",
			paramNames:     []string{"network_id"},
			ParamValues:    []string{"n1"},
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Message{getLog()},
			expectedError:  "",
			expectedListRequest: &protos.ListLogsRequest{
				NetworkId:  "n1",
				Filter:     &protos.LogFilter{},
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list logs with limit and offset",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "",
			queryParamsString: url.Values{"limit": {"4"}, "offset": {"3"}},
			expectedListRequest: &protos.ListLogsRequest{
				NetworkId: "n1",
				Filter:    &protos.LogFilter{},
				Pagination: &protos.Pagination{
					Limit:  wrapperspb.Int64(4),
					Offset: wrapperspb.Int64(3),
				},
			},
		},
		{
			testName:       "test list logs with all query params",
			paramNames:     []string{"network_id"},
			ParamValues:    []string{"n1"},
			expectedStatus: http.StatusOK,
			expectedResult: []*models.Message{getLog()},
			expectedError:  "",
			queryParamsString: url.Values{
				"from":          {"CBSD"},
				"to":            {"DP"},
				"type":          {"someType"},
				"serial_number": {"abc123"},
				"fcc_id":        {"some_id"},
				"response_code": {"0"},
				"begin":         {"2022-01-14T10:23:49.871036+00:00"},
				"end":           {"2022-02-14T10:23:49.871036+00:00"},
			},
			expectedListRequest: &protos.ListLogsRequest{
				NetworkId: "n1",
				Filter: &protos.LogFilter{
					From:                "CBSD",
					To:                  "DP",
					Name:                "someType",
					SerialNumber:        "abc123",
					FccId:               "some_id",
					ResponseCode:        wrapperspb.Int64(0),
					BeginTimestampMilli: wrapperspb.Int64(1642155829871),
					EndTimestampMilli:   wrapperspb.Int64(1644834229871),
				},
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list logs begin date without offset",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    []*models.Message{getLog()},
			queryParamsString: url.Values{"begin": {"2022-01-14T10:23:49.871036Z"}},
			expectedListRequest: &protos.ListLogsRequest{
				NetworkId: "n1",
				Filter: &protos.LogFilter{
					BeginTimestampMilli: wrapperspb.Int64(1642155829871),
				},
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list logs end date without offset",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    []*models.Message{getLog()},
			queryParamsString: url.Values{"end": {"2022-01-14T10:23:49.871036Z"}},
			expectedListRequest: &protos.ListLogsRequest{
				NetworkId: "n1",
				Filter: &protos.LogFilter{
					EndTimestampMilli: wrapperspb.Int64(1642155829871),
				},
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list logs with incorrect begin param",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "'2334321344-33321' is not a proper value for begin",
			queryParamsString: url.Values{"begin": {"2334321344-33321"}},
		},
		{
			testName:          "test list logs with incorrect end param",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "'2334321344-33321' is not a proper value for end",
			queryParamsString: url.Values{"end": {"2334321344-33321"}},
		},
		{
			testName:          "test list logs with incorrect begin format",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "'incorrect_date_string' is not a proper value for begin",
			queryParamsString: url.Values{"begin": {"incorrect_date_string"}},
		},
		{
			testName:          "test list logs with incorrect end format",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "'incorrect_date_string' is not a proper value for end",
			queryParamsString: url.Values{"end": {"incorrect_date_string"}},
		},
		{
			testName:          "test list logs with incorrect response_code format",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedResult:    []*models.Message{getLog()},
			expectedError:     "'foo' is not a proper value for response_code",
			queryParamsString: url.Values{"response_code": {"foo"}},
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
	}
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	s.logsServer.listResponse = &protos.ListLogsResponse{Logs: []*protos.Log{getLogDetails()}}
	listLogs := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageLogsPath, obsidian.GET).HandlerFunc
	for _, t := range testCases {
		s.Run(t.testName, func() {
			s.logsServer.expectedListRequest = t.expectedListRequest
			tc := tests.Test{
				Method:         "GET",
				URL:            handlers.ManageLogsPath + "?" + t.queryParamsString.Encode(),
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

func (s *HandlersTestSuite) TestListCbsds() {
	testCases := []struct {
		testName            string
		paramNames          []string
		ParamValues         []string
		model               db.Model
		expectedStatus      int
		expectedResult      []*models.Cbsd
		expectedError       string
		queryParamsString   string
		expectedListRequest *protos.ListCbsdRequest
	}{
		{
			testName:          "test list cbsds without query params",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    []*models.Cbsd{getCbsd()},
			expectedError:     "",
			queryParamsString: "",
			expectedListRequest: &protos.ListCbsdRequest{
				NetworkId:  "n1",
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list cbsds without limit and offset",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    []*models.Cbsd{getCbsd()},
			expectedError:     "",
			queryParamsString: "?limit=4&offset=3",
			expectedListRequest: &protos.ListCbsdRequest{
				NetworkId: "n1",
				Pagination: &protos.Pagination{
					Limit:  wrapperspb.Int64(4),
					Offset: wrapperspb.Int64(3),
				},
			},
		},
		{
			testName:          "test list cbsds with incorrect limit value",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'incorrect_limit_value' is not a proper value for limit",
			queryParamsString: "?limit=incorrect_limit_value",
		},
		{
			testName:          "test list cbsds with incorrect offset value",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusBadRequest,
			expectedError:     "'incorrect_offset_value' is not a proper value for offset",
			queryParamsString: "?offset=incorrect_offset_value",
		},
	}
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	s.cbsdServer.listResponse = &protos.ListCbsdResponse{Details: []*protos.CbsdDetails{getCbsdDetails()}}
	listCbsds := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdsPath, obsidian.GET).HandlerFunc
	for _, t := range testCases {
		s.Run(t.testName, func() {
			s.cbsdServer.expectedListRequest = t.expectedListRequest
			tc := tests.Test{
				Method:         "GET",
				URL:            handlers.ManageCbsdsPath + t.queryParamsString,
				Payload:        nil,
				ParamNames:     t.paramNames,
				ParamValues:    t.ParamValues,
				Handler:        listCbsds,
				ExpectedStatus: t.expectedStatus,
				ExpectedResult: tests.JSONMarshaler(t.expectedResult),
				ExpectedError:  t.expectedError,
			}
			tests.RunUnitTest(s.T(), e, tc)
		})
	}
}

func (s *HandlersTestSuite) TestFetchCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	s.cbsdServer.fetchResponse = &protos.FetchCbsdResponse{Details: getCbsdDetails()}
	s.cbsdServer.expectedFetchRequest = &protos.FetchCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	fetchCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.GET).HandlerFunc
	expectedResult := getCbsd()
	tc := tests.Test{
		Method:         "GET",
		URL:            "/magma/v1/dp/n1/cbsds/1",
		Payload:        nil,
		ParamNames:     []string{"network_id", "cbsd_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        fetchCbsd,
		ExpectedStatus: http.StatusOK,
		ExpectedResult: tests.JSONMarshaler(expectedResult),
		ExpectedError:  "",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestFetchNonexistentCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedFetchRequest = &protos.FetchCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	fetchCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.GET).HandlerFunc
	tc := tests.Test{
		Method:                 "GET",
		URL:                    "/magma/v1/dp/n1/cbsds/1",
		Payload:                nil,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "1"},
		Handler:                fetchCbsd,
		ExpectedStatus:         http.StatusNotFound,
		ExpectedErrorSubstring: errorMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestCreateCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	payload := createOrUpdateCbsdPayload()
	s.cbsdServer.createResponse = &protos.CreateCbsdResponse{}
	s.cbsdServer.expectedCreateRequest = &protos.CreateCbsdRequest{
		NetworkId: "n1",
		Data: &protos.CbsdData{
			UserId:       *payload.UserID,
			FccId:        *payload.FccID,
			SerialNumber: *payload.SerialNumber,
			Capabilities: &protos.Capabilities{
				MinPower:         *payload.Capabilities.MinPower,
				MaxPower:         *payload.Capabilities.MaxPower,
				NumberOfAntennas: *payload.Capabilities.NumberOfAntennas,
				AntennaGain:      *payload.Capabilities.AntennaGain,
			},
		},
	}
	createCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdsPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/dp/n1/cbsds",
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createCbsd,
		ExpectedStatus: http.StatusCreated,
		ExpectedError:  "",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestCreateCbsdWithoutAllRequiredParams() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	payload := &models.Cbsd{
		Capabilities: &models.Capabilities{
			AntennaGain:      to_pointer.Float(1),
			NumberOfAntennas: to_pointer.Int(1),
		},
		SerialNumber: to_pointer.Str("someSerialNumber"),
	}
	s.cbsdServer.createResponse = &protos.CreateCbsdResponse{}
	s.cbsdServer.expectedCreateRequest = &protos.CreateCbsdRequest{
		NetworkId: "n1",
		Data: &protos.CbsdData{
			SerialNumber: *payload.SerialNumber,
			Capabilities: &protos.Capabilities{
				NumberOfAntennas: *payload.Capabilities.NumberOfAntennas,
				AntennaGain:      *payload.Capabilities.AntennaGain,
			},
		},
	}
	createCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdsPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:         "POST",
		URL:            "/magma/v1/dp/n1/cbsds",
		Payload:        payload,
		ParamNames:     []string{"network_id"},
		ParamValues:    []string{"n1"},
		Handler:        createCbsd,
		ExpectedStatus: http.StatusBadGateway,
		ExpectedError:  "Bad Gateway",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestDeleteCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	s.cbsdServer.deleteResponse = &protos.DeleteCbsdResponse{}
	deleteCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.DELETE).HandlerFunc
	s.cbsdServer.expectedDeleteRequest = &protos.DeleteCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	tc := tests.Test{
		Method:         "DELETE",
		URL:            "/magma/v1/dp/n1/cbsds/1",
		ParamNames:     []string{"network_id", "cbsd_id"},
		ParamValues:    []string{"n1", "1"},
		Handler:        deleteCbsd,
		ExpectedStatus: http.StatusNoContent,
		ExpectedError:  "",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestDeleteNonexistentCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedDeleteRequest = &protos.DeleteCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	deleteCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.DELETE).HandlerFunc
	tc := tests.Test{
		Method:                 "DELETE",
		URL:                    "/magma/v1/dp/n1/cbsds/1",
		Payload:                nil,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "1"},
		Handler:                deleteCbsd,
		ExpectedStatus:         http.StatusNotFound,
		ExpectedErrorSubstring: errorMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestUpdateCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	s.cbsdServer.updateResponse = &protos.UpdateCbsdResponse{}
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.PUT).HandlerFunc
	payload := createOrUpdateCbsdPayload()
	s.cbsdServer.expectedUpdateRequest = &protos.UpdateCbsdRequest{
		NetworkId: "n1",
		Data: &protos.CbsdData{
			UserId:       *payload.UserID,
			FccId:        *payload.FccID,
			SerialNumber: *payload.SerialNumber,
			Capabilities: &protos.Capabilities{
				MinPower:         *payload.Capabilities.MinPower,
				MaxPower:         *payload.Capabilities.MaxPower,
				NumberOfAntennas: *payload.Capabilities.NumberOfAntennas,
				AntennaGain:      *payload.Capabilities.AntennaGain,
			},
		},
	}
	tc := tests.Test{
		Method:         "PUT",
		URL:            "/magma/v1/dp/n1/cbsds/0",
		Payload:        payload,
		ParamNames:     []string{"network_id", "cbsd_id"},
		ParamValues:    []string{"n1", "0"},
		Handler:        updateCbsd,
		ExpectedStatus: http.StatusNoContent,
		ExpectedError:  "",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestUpdateNonexistentCbsd() {
	e := echo.New()
	obsidianHandlers := handlers.GetHandlers()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	payload := createOrUpdateCbsdPayload()
	s.cbsdServer.expectedUpdateRequest = &protos.UpdateCbsdRequest{
		NetworkId: "n1",
		Data: &protos.CbsdData{
			UserId:       *payload.UserID,
			FccId:        *payload.FccID,
			SerialNumber: *payload.SerialNumber,
			Capabilities: &protos.Capabilities{
				MinPower:         *payload.Capabilities.MinPower,
				MaxPower:         *payload.Capabilities.MaxPower,
				NumberOfAntennas: *payload.Capabilities.NumberOfAntennas,
				AntennaGain:      *payload.Capabilities.AntennaGain,
			},
		},
	}
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, handlers.ManageCbsdPath, obsidian.PUT).HandlerFunc
	tc := tests.Test{
		Method:                 "PUT",
		URL:                    "/magma/v1/dp/n1/cbsds/0",
		Payload:                payload,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "0"},
		Handler:                updateCbsd,
		ExpectedStatus:         http.StatusNotFound,
		ExpectedErrorSubstring: errorMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestGetPagination() {
	testCases := []struct {
		testName       string
		URL            string
		expectedLimit  *wrapperspb.Int64Value
		expectedOffset *wrapperspb.Int64Value
	}{
		{
			testName:       "test pagination with limit and offset",
			URL:            "/magma/v1/dp/some_network/cbsds?limit=1&offset=2",
			expectedLimit:  wrapperspb.Int64(1),
			expectedOffset: wrapperspb.Int64(2),
		},
		{
			testName:       "test pagination with limit only",
			URL:            "/magma/v1/dp/some_network/cbsds?limit=1",
			expectedLimit:  wrapperspb.Int64(1),
			expectedOffset: nil,
		},
		{
			testName:       "test pagination without limit and offset",
			URL:            "/magma/v1/dp/some_network/cbsds",
			expectedLimit:  nil,
			expectedOffset: nil,
		},
	}
	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			t := tests.Test{}
			req := *httptest.NewRequest(t.Method, tc.URL, bytes.NewReader(nil))
			c := echo.New().NewContext(&req, httptest.NewRecorder())
			pag, _ := handlers.GetPagination(c)
			if tc.expectedLimit != nil {
				assert.Equal(s.T(), tc.expectedLimit, pag.Limit)
			} else {
				assert.Nil(s.T(), pag.Limit)
			}
			if tc.expectedOffset != nil {
				assert.Equal(s.T(), tc.expectedOffset, pag.Offset)
			} else {
				assert.Nil(s.T(), pag.Offset)
			}
		})
	}
}

func (s *HandlersTestSuite) TestGetPaginationWithoutLimit() {
	t := tests.Test{}
	req := *httptest.NewRequest(t.Method, "/magma/v1/dp/some_network/cbsds?offset=2", bytes.NewReader(nil))
	c := echo.New().NewContext(&req, httptest.NewRecorder())
	pag, err := handlers.GetPagination(c)
	assert.Nil(s.T(), pag)
	assert.EqualError(s.T(), err, "code=400, message=offset requires a limit")
}

func (s *HandlersTestSuite) TestGetPaginationWithIncorrectLimitAndOffset() {
	testCases := []struct {
		URL           string
		expectedError string
	}{
		{
			URL:           "/magma/v1/dp/some_network/cbsds?offset=foo",
			expectedError: "'foo' is not a proper value for offset",
		},
		{
			URL:           "/magma/v1/dp/some_network/cbsds?offset=.",
			expectedError: "'.' is not a proper value for offset",
		},
		{
			URL:           "/magma/v1/dp/some_network/cbsds?limit=foo",
			expectedError: "'foo' is not a proper value for limit",
		},
		{
			URL:           "/magma/v1/dp/some_network/cbsds?limit=.",
			expectedError: "'.' is not a proper value for limit",
		},
	}
	for _, tc := range testCases {
		t := tests.Test{}
		req := *httptest.NewRequest(t.Method, tc.URL, bytes.NewReader(nil))
		c := echo.New().NewContext(&req, httptest.NewRecorder())
		pag, err := handlers.GetPagination(c)
		assert.Nil(s.T(), pag)
		assert.EqualError(s.T(), err, "code=400, message="+tc.expectedError)
	}
}

func (s *stubCbsdServer) CreateCbsd(ctx context.Context, request *protos.CreateCbsdRequest) (*protos.CreateCbsdResponse, error) {
	assert.Equal(s.t, s.expectedCreateRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedCreateRequest.Data, request.Data)
	return s.createResponse, s.err
}

func (s *stubCbsdServer) UpdateCbsd(ctx context.Context, request *protos.UpdateCbsdRequest) (*protos.UpdateCbsdResponse, error) {
	assert.Equal(s.t, s.expectedUpdateRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedUpdateRequest.Id, request.Id)
	assert.Equal(s.t, s.expectedUpdateRequest.Data, request.Data)
	return s.updateResponse, s.err
}

func (s *stubCbsdServer) DeleteCbsd(ctx context.Context, request *protos.DeleteCbsdRequest) (*protos.DeleteCbsdResponse, error) {
	assert.Equal(s.t, s.expectedDeleteRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedDeleteRequest.Id, request.Id)
	return s.deleteResponse, s.err
}

func (s *stubCbsdServer) FetchCbsd(ctx context.Context, request *protos.FetchCbsdRequest) (*protos.FetchCbsdResponse, error) {
	assert.Equal(s.t, s.expectedFetchRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedFetchRequest.Id, request.Id)
	return s.fetchResponse, s.err
}

func (s *stubCbsdServer) ListCbsds(ctx context.Context, request *protos.ListCbsdRequest) (*protos.ListCbsdResponse, error) {
	assert.Equal(s.t, s.expectedListRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Limit, request.Pagination.Limit)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Offset, request.Pagination.Offset)
	return s.listResponse, s.err
}

func (s *stubLogsServer) ListLogs(ctx context.Context, request *protos.ListLogsRequest) (*protos.ListLogsResponse, error) {
	assert.Equal(s.t, s.expectedListRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Limit, request.Pagination.Limit)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Offset, request.Pagination.Offset)
	assert.Equal(s.t, s.expectedListRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedListRequest.Filter.From, request.Filter.From)
	assert.Equal(s.t, s.expectedListRequest.Filter.To, request.Filter.To)
	assert.Equal(s.t, s.expectedListRequest.Filter.Name, request.Filter.Name)
	assert.Equal(s.t, s.expectedListRequest.Filter.SerialNumber, request.Filter.SerialNumber)
	if s.expectedListRequest.Filter.BeginTimestampMilli != nil {
		assert.Equal(s.t, s.expectedListRequest.Filter.BeginTimestampMilli, request.Filter.BeginTimestampMilli)
	}
	if s.expectedListRequest.Filter.EndTimestampMilli != nil {
		assert.Equal(s.t, s.expectedListRequest.Filter.EndTimestampMilli, request.Filter.EndTimestampMilli)
	}
	if s.expectedListRequest.Filter.ResponseCode != nil {
		assert.Equal(s.t, s.expectedListRequest.Filter.ResponseCode, request.Filter.ResponseCode)
	}
	assert.Equal(s.t, s.expectedListRequest.Filter.FccId, request.Filter.FccId)
	return s.listResponse, s.err
}

func getCbsd() *models.Cbsd {
	return &models.Cbsd{
		Capabilities: &models.Capabilities{
			AntennaGain:      to_pointer.Float(1),
			MaxPower:         to_pointer.Float(24),
			MinPower:         to_pointer.Float(0),
			NumberOfAntennas: to_pointer.Int(1),
		},
		CbsdID: "someCbsdId",
		FccID:  to_pointer.Str("someFCCId"),
		Grant: &models.Grant{
			BandwidthMhz:       0,
			FrequencyMhz:       0,
			GrantExpireTime:    *to_pointer.TimeToDateTime(0),
			MaxEirp:            to_pointer.Float(0),
			State:              "someState",
			TransmitExpireTime: *to_pointer.TimeToDateTime(0),
		},
		ID:           0,
		SerialNumber: to_pointer.Str("someSerialNumber"),
		State:        "unregistered",
		UserID:       to_pointer.Str("someUserId"),
		IsActive:     false,
	}
}

func createOrUpdateCbsdPayload() *models.Cbsd {
	return &models.Cbsd{
		Capabilities: &models.Capabilities{
			AntennaGain:      to_pointer.Float(1),
			MaxPower:         to_pointer.Float(24),
			MinPower:         to_pointer.Float(0),
			NumberOfAntennas: to_pointer.Int(1),
		},
		FccID:        to_pointer.Str("someFCCId"),
		SerialNumber: to_pointer.Str("someSerialNumber"),
		UserID:       to_pointer.Str("someUserId"),
	}
}

func getCbsdDetails() *protos.CbsdDetails {
	return &protos.CbsdDetails{
		Id:       0,
		Data:     getCbsdData(),
		CbsdId:   "someCbsdId",
		State:    "unregistered",
		IsActive: false,
		Grant: &protos.GrantDetails{
			BandwidthMhz:            0,
			FrequencyMhz:            0,
			MaxEirp:                 0,
			State:                   "someState",
			TransmitExpireTimestamp: 0,
			GrantExpireTimestamp:    0,
		},
	}
}

func getCbsdData() *protos.CbsdData {
	return &protos.CbsdData{
		UserId:       "someUserId",
		FccId:        "someFCCId",
		SerialNumber: "someSerialNumber",
		Capabilities: &protos.Capabilities{
			MinPower:         0,
			MaxPower:         24,
			NumberOfAntennas: 1,
			AntennaGain:      1,
		},
	}
}

func getLog() *models.Message {
	return &models.Message{
		From:         "CBSD",
		To:           "DP",
		Body:         "some message",
		FccID:        "some fcc id",
		SerialNumber: "12345a",
		Time:         *to_pointer.TimeMilliToDate(123),
		Type:         "RegistrationRequest",
	}
}

func getLogDetails() *protos.Log {
	return &protos.Log{
		From:           "CBSD",
		To:             "DP",
		Name:           "RegistrationRequest",
		Message:        "some message",
		SerialNumber:   "12345a",
		FccId:          "some fcc id",
		TimestampMilli: 123,
	}
}
