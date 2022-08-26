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

package cbsd_test

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"magma/dp/cloud/go/dp"
	"magma/dp/cloud/go/protos"
	dp_service "magma/dp/cloud/go/services/dp"
	b "magma/dp/cloud/go/services/dp/builders"
	"magma/dp/cloud/go/services/dp/obsidian/cbsd"
	"magma/dp/cloud/go/services/dp/obsidian/models"
	"magma/dp/cloud/go/services/dp/obsidian/to_pointer"
	"magma/dp/cloud/go/services/dp/storage/db"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/cloud/go/services/obsidian/tests"
	"magma/orc8r/cloud/go/test_utils"
)

func TestHandlers(t *testing.T) {
	suite.Run(t, &HandlersTestSuite{})
}

type HandlersTestSuite struct {
	suite.Suite
	cbsdServer *stubCbsdServer
}

func (s *HandlersTestSuite) SetupTest() {
	s.cbsdServer = &stubCbsdServer{}
	s.cbsdServer.t = s.T()
	srv, lis, _ := test_utils.NewTestService(s.T(), dp.ModuleName, dp_service.ServiceName)
	protos.RegisterCbsdManagementServer(srv.GrpcServer, s.cbsdServer)
	go srv.RunTest(lis, nil)
}

type stubCbsdServer struct {
	protos.UnimplementedCbsdManagementServer
	expectedListRequest       *protos.ListCbsdRequest
	listResponse              *protos.ListCbsdResponse
	expectedFetchRequest      *protos.FetchCbsdRequest
	fetchResponse             *protos.FetchCbsdResponse
	expectedCreateRequest     *protos.CreateCbsdRequest
	createResponse            *protos.CreateCbsdResponse
	expectedUpdateRequest     *protos.UpdateCbsdRequest
	updateResponse            *protos.UpdateCbsdResponse
	expectedDeleteRequest     *protos.DeleteCbsdRequest
	deleteResponse            *protos.DeleteCbsdResponse
	expectedDeregisterRequest *protos.DeregisterCbsdRequest
	expectedRelinquishRequest *protos.RelinquishCbsdRequest
	err                       error
	t                         *testing.T
}

func (s *HandlersTestSuite) TestListCbsds() {
	testCases := []struct {
		testName            string
		paramNames          []string
		ParamValues         []string
		model               db.Model
		expectedStatus      int
		expectedResult      *models.PaginatedCbsds
		expectedError       string
		queryParamsString   string
		expectedListRequest *protos.ListCbsdRequest
	}{
		{
			testName:          "test list cbsds without query params",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    b.GetPaginatedCbsds(b.NewCbsdModelPayloadBuilder().WithGrant()),
			expectedError:     "",
			queryParamsString: "",
			expectedListRequest: &protos.ListCbsdRequest{
				NetworkId:  "n1",
				Pagination: &protos.Pagination{},
			},
		},
		{
			testName:          "test list cbsds with limit and offset",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    b.GetPaginatedCbsds(b.NewCbsdModelPayloadBuilder().WithGrant()),
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
			testName:          "test list cbsds with all params",
			paramNames:        []string{"network_id"},
			ParamValues:       []string{"n1"},
			expectedStatus:    http.StatusOK,
			expectedResult:    b.GetPaginatedCbsds(b.NewCbsdModelPayloadBuilder().WithGrant()),
			expectedError:     "",
			queryParamsString: "?limit=4&offset=3&serial_number=foo123",
			expectedListRequest: &protos.ListCbsdRequest{
				NetworkId: "n1",
				Pagination: &protos.Pagination{
					Limit:  wrapperspb.Int64(4),
					Offset: wrapperspb.Int64(3),
				},
				Filter: &protos.CbsdFilter{
					SerialNumber: "foo123",
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
	obsidianHandlers := cbsd.GetHandlers()
	s.cbsdServer.listResponse = &protos.ListCbsdResponse{
		Details: []*protos.CbsdDetails{
			b.NewDetailedProtoCbsdBuilder(
				b.NewCbsdProtoPayloadBuilder()).
				WithGrant().
				WithCbsdId("some_cbsd_id").
				Details},
		TotalCount: 1,
	}
	listCbsds := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdsPath, obsidian.GET).HandlerFunc
	for _, t := range testCases {
		s.Run(t.testName, func() {
			s.cbsdServer.expectedListRequest = t.expectedListRequest
			tc := tests.Test{
				Method:         http.MethodGet,
				URL:            cbsd.ManageCbsdsPath + t.queryParamsString,
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
	obsidianHandlers := cbsd.GetHandlers()
	s.cbsdServer.fetchResponse = &protos.FetchCbsdResponse{Details: b.NewDetailedProtoCbsdBuilder(
		b.NewCbsdProtoPayloadBuilder()).
		WithGrant().
		WithCbsdId("some_cbsd_id").
		Details}
	s.cbsdServer.expectedFetchRequest = &protos.FetchCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	fetchCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.GET).HandlerFunc
	expectedResult := b.NewCbsdModelPayloadBuilder().WithGrant().Payload
	tc := tests.Test{
		Method:         http.MethodGet,
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
	obsidianHandlers := cbsd.GetHandlers()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedFetchRequest = &protos.FetchCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	fetchCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.GET).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodGet,
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
	testCases := []struct {
		name            string
		inputPayload    *models.MutableCbsd
		expectedPayload *protos.CbsdData
		expectedStatus  int
		expectedError   string
	}{{
		name:            "test create without installation param",
		inputPayload:    b.NewMutableCbsdModelPayloadBuilder().Payload,
		expectedPayload: b.NewCbsdProtoPayloadBuilder().Payload,
		expectedStatus:  http.StatusCreated,
		expectedError:   "",
	}, {
		name: "test create with antenna gain",
		inputPayload: b.NewMutableCbsdModelPayloadBuilder().
			WithAntennaGain(10.5).Payload,
		expectedPayload: b.NewCbsdProtoPayloadBuilder().
			WithEmptyInstallationParam().
			WithAntennaGain(10.5).Payload,
		expectedStatus: http.StatusCreated,
		expectedError:  "",
	}, {
		name: "test create with carrier aggregation is enabled and grant_redundancy is false",
		inputPayload: b.NewMutableCbsdModelPayloadBuilder().
			WithGrantRedundancy(to_pointer.Bool(false)).
			WithCarrierAggregationEnabled(to_pointer.Bool(true)).
			Payload,
		expectedStatus: http.StatusBadRequest,
		expectedError:  "grant_redundancy cannot be set to false when carrier_aggregation_enabled is enabled",
	}, {
		name: "test create with max_ibw_mhz lesser than bandwidth_mhz",
		inputPayload: b.NewMutableCbsdModelPayloadBuilder().
			WithBandwidth(10).
			WithMaxIbwMhz(5).
			Payload,
		expectedStatus: http.StatusBadRequest,
		expectedError:  "max_ibw_mhz cannot be less than bandwidth_mhz",
	}, {
		name: "test failed model validation raises 400",
		inputPayload: b.NewMutableCbsdModelPayloadBuilder().
			WithSingleStepEnabled(nil).
			Payload,
		expectedStatus: http.StatusBadRequest,
		expectedError:  "single_step_enabled in body is required",
	}}
	for _, tc := range testCases {
		s.Run(tc.name, func() {
			e := echo.New()
			obsidianHandlers := cbsd.GetHandlers()
			payload := tc.inputPayload
			s.cbsdServer.createResponse = &protos.CreateCbsdResponse{}
			s.cbsdServer.expectedCreateRequest = &protos.CreateCbsdRequest{
				NetworkId: "n1",
				Data:      tc.expectedPayload,
			}
			createCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdsPath, obsidian.POST).HandlerFunc
			tc := tests.Test{
				Method:         http.MethodPost,
				URL:            "/magma/v1/dp/n1/cbsds",
				Payload:        payload,
				ParamNames:     []string{"network_id"},
				ParamValues:    []string{"n1"},
				Handler:        createCbsd,
				ExpectedStatus: tc.expectedStatus,
				ExpectedError:  tc.expectedError,
			}
			tests.RunUnitTest(s.T(), e, tc)
		})
	}
}

func (s *HandlersTestSuite) TestCreateWithDuplicateUniqueFieldsReturnsConflict() {
	const errMsg = "some error"
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	payload := b.NewMutableCbsdModelPayloadBuilder().Payload
	data, _ := models.CbsdToBackend(payload)
	s.cbsdServer.err = status.Error(codes.AlreadyExists, errMsg)
	s.cbsdServer.expectedCreateRequest = &protos.CreateCbsdRequest{
		NetworkId: "n1",
		Data:      data,
	}
	createCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdsPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPost,
		URL:                    "/magma/v1/dp/n1/cbsds",
		Payload:                payload,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                createCbsd,
		ExpectedStatus:         http.StatusConflict,
		ExpectedErrorSubstring: errMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestCreateCbsdWithoutAllRequiredParams() {
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	payload := b.NewMutableCbsdModelPayloadBuilder().Empty().WithSerialNumber("someSerialNumber").Payload
	createCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdsPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPost,
		URL:                    "/magma/v1/dp/n1/cbsds",
		Payload:                payload,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                createCbsd,
		ExpectedStatus:         http.StatusBadRequest,
		ExpectedErrorSubstring: "validation failure list",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestDeleteCbsd() {
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	s.cbsdServer.deleteResponse = &protos.DeleteCbsdResponse{}
	deleteCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.DELETE).HandlerFunc
	s.cbsdServer.expectedDeleteRequest = &protos.DeleteCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	tc := tests.Test{
		Method:         http.MethodDelete,
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
	const errorMsg = "some msg"
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedDeleteRequest = &protos.DeleteCbsdRequest{
		NetworkId: "n1",
		Id:        1,
	}
	deleteCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.DELETE).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodDelete,
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
	obsidianHandlers := cbsd.GetHandlers()
	s.cbsdServer.updateResponse = &protos.UpdateCbsdResponse{}
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.PUT).HandlerFunc
	payload := b.NewMutableCbsdModelPayloadBuilder().Payload
	s.cbsdServer.expectedUpdateRequest = &protos.UpdateCbsdRequest{
		NetworkId: "n1",
		Data:      b.NewCbsdProtoPayloadBuilder().Payload,
	}
	tc := tests.Test{
		Method:         http.MethodPut,
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

func (s *HandlersTestSuite) TestUpdateCbsdWithoutAllRequiredParams() {
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	payload := b.NewMutableCbsdModelPayloadBuilder().Empty().WithSerialNumber("someSerialNumber").Payload
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.PUT).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPut,
		URL:                    "/magma/v1/dp/n1/cbsds",
		Payload:                payload,
		ParamNames:             []string{"network_id"},
		ParamValues:            []string{"n1"},
		Handler:                updateCbsd,
		ExpectedStatus:         http.StatusBadRequest,
		ExpectedErrorSubstring: "validation failure list",
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestUpdateNonexistentCbsd() {
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	payload := b.NewMutableCbsdModelPayloadBuilder().Payload
	s.cbsdServer.expectedUpdateRequest = &protos.UpdateCbsdRequest{
		NetworkId: "n1",
		Data:      b.NewCbsdProtoPayloadBuilder().Payload,
	}
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.PUT).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPut,
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

func (s *HandlersTestSuite) TestUpdateCbsdWithDuplicateUniqueFieldsReturnsConflict() {
	e := echo.New()
	obsidianHandlers := cbsd.GetHandlers()
	payload := b.NewMutableCbsdModelPayloadBuilder().Payload
	const errMsg = "some error"
	data, _ := models.CbsdToBackend(payload)
	s.cbsdServer.err = status.Error(codes.AlreadyExists, errMsg)
	s.cbsdServer.expectedUpdateRequest = &protos.UpdateCbsdRequest{
		NetworkId: "n1",
		Id:        1,
		Data:      data,
	}
	updateCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.ManageCbsdPath, obsidian.PUT).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPost,
		URL:                    "/magma/v1/dp/n1/cbsds/1",
		Payload:                payload,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "1"},
		Handler:                updateCbsd,
		ExpectedStatus:         http.StatusConflict,
		ExpectedErrorSubstring: errMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestDeregisterCbsd() {
	e := echo.New()
	s.cbsdServer.expectedDeregisterRequest = &protos.DeregisterCbsdRequest{
		NetworkId: "n1",
		Id:        0,
	}
	obsidianHandlers := cbsd.GetHandlers()
	deregisterCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.DeregisterCbsdPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:         http.MethodPut,
		URL:            "/magma/v1/n1/cbsds/0",
		Handler:        deregisterCbsd,
		ParamNames:     []string{"network_id", "cbsd_id"},
		ParamValues:    []string{"n1", "0"},
		ExpectedStatus: http.StatusNoContent,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestDeregisterNonexistentCbsd() {
	e := echo.New()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedDeregisterRequest = &protos.DeregisterCbsdRequest{
		NetworkId: "n1",
		Id:        0,
	}
	obsidianHandlers := cbsd.GetHandlers()
	deregisterCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.DeregisterCbsdPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPut,
		URL:                    "/magma/v1/dp/n1/cbsds/0",
		Handler:                deregisterCbsd,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "0"},
		ExpectedStatus:         http.StatusNotFound,
		ExpectedErrorSubstring: errorMsg,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestRelinquishCbsd() {
	e := echo.New()
	s.cbsdServer.expectedRelinquishRequest = &protos.RelinquishCbsdRequest{
		NetworkId: "n1",
		Id:        0,
	}
	obsidianHandlers := cbsd.GetHandlers()
	relinquishCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.RelinquishCbsdPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:         http.MethodPut,
		URL:            "/magma/v1/n1/cbsds/0",
		Handler:        relinquishCbsd,
		ParamNames:     []string{"network_id", "cbsd_id"},
		ParamValues:    []string{"n1", "0"},
		ExpectedStatus: http.StatusNoContent,
	}
	tests.RunUnitTest(s.T(), e, tc)
}

func (s *HandlersTestSuite) TestRelinquishNonexistentCbsd() {
	e := echo.New()
	const errorMsg = "some msg"
	s.cbsdServer.err = status.Error(codes.NotFound, errorMsg)
	s.cbsdServer.expectedRelinquishRequest = &protos.RelinquishCbsdRequest{
		NetworkId: "n1",
		Id:        0,
	}
	obsidianHandlers := cbsd.GetHandlers()
	relinquishCbsd := tests.GetHandlerByPathAndMethod(s.T(), obsidianHandlers, cbsd.RelinquishCbsdPath, obsidian.POST).HandlerFunc
	tc := tests.Test{
		Method:                 http.MethodPut,
		URL:                    "/magma/v1/dp/n1/cbsds/0",
		Handler:                relinquishCbsd,
		ParamNames:             []string{"network_id", "cbsd_id"},
		ParamValues:            []string{"n1", "0"},
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
			pag, _ := cbsd.GetPagination(c)
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
	pag, err := cbsd.GetPagination(c)
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
		pag, err := cbsd.GetPagination(c)
		assert.Nil(s.T(), pag)
		assert.EqualError(s.T(), err, "code=400, message="+tc.expectedError)
	}
}

func (s *HandlersTestSuite) TestGetCbsdFilter() {
	testCases := []struct {
		testName             string
		URL                  string
		expectedSerialNumber string
		expectedError        string
	}{
		{
			testName:             "test filter with serial_number",
			URL:                  "/some/url?serial_number=some_serial_number",
			expectedSerialNumber: "some_serial_number",
		},
		{
			testName:             "test filter without params",
			URL:                  "/some/url",
			expectedSerialNumber: "",
		},
	}
	for _, tc := range testCases {
		s.Run(tc.testName, func() {
			t := tests.Test{}
			req := *httptest.NewRequest(t.Method, tc.URL, bytes.NewReader(nil))
			c := echo.New().NewContext(&req, httptest.NewRecorder())
			f := cbsd.GetCbsdFilter(c)
			s.Equal(tc.expectedSerialNumber, f.SerialNumber)
		})
	}
}

func (s *stubCbsdServer) CreateCbsd(_ context.Context, request *protos.CreateCbsdRequest) (*protos.CreateCbsdResponse, error) {
	assert.Equal(s.t, s.expectedCreateRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedCreateRequest.Data, request.Data)
	return s.createResponse, s.err
}

func (s *stubCbsdServer) UserUpdateCbsd(_ context.Context, request *protos.UpdateCbsdRequest) (*protos.UpdateCbsdResponse, error) {
	assert.Equal(s.t, s.expectedUpdateRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedUpdateRequest.Id, request.Id)
	assert.Equal(s.t, s.expectedUpdateRequest.Data, request.Data)
	return s.updateResponse, s.err
}

func (s *stubCbsdServer) DeleteCbsd(_ context.Context, request *protos.DeleteCbsdRequest) (*protos.DeleteCbsdResponse, error) {
	assert.Equal(s.t, s.expectedDeleteRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedDeleteRequest.Id, request.Id)
	return s.deleteResponse, s.err
}

func (s *stubCbsdServer) FetchCbsd(_ context.Context, request *protos.FetchCbsdRequest) (*protos.FetchCbsdResponse, error) {
	assert.Equal(s.t, s.expectedFetchRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedFetchRequest.Id, request.Id)
	return s.fetchResponse, s.err
}

func (s *stubCbsdServer) ListCbsds(_ context.Context, request *protos.ListCbsdRequest) (*protos.ListCbsdResponse, error) {
	assert.Equal(s.t, s.expectedListRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Limit, request.Pagination.Limit)
	assert.Equal(s.t, s.expectedListRequest.Pagination.Offset, request.Pagination.Offset)
	return s.listResponse, s.err
}

func (s *stubCbsdServer) DeregisterCbsd(_ context.Context, request *protos.DeregisterCbsdRequest) (*protos.DeregisterCbsdResponse, error) {
	assert.Equal(s.t, s.expectedDeregisterRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedDeregisterRequest.Id, request.Id)
	return &protos.DeregisterCbsdResponse{}, s.err
}

func (s *stubCbsdServer) RelinquishCbsd(_ context.Context, request *protos.RelinquishCbsdRequest) (*protos.RelinquishCbsdResponse, error) {
	assert.Equal(s.t, s.expectedRelinquishRequest.NetworkId, request.NetworkId)
	assert.Equal(s.t, s.expectedRelinquishRequest.Id, request.Id)
	return &protos.RelinquishCbsdResponse{}, s.err
}
