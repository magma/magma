/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"magma/fbinternal/cloud/go/serdes"
	"magma/fbinternal/cloud/go/services/testcontroller"
	"magma/fbinternal/cloud/go/services/testcontroller/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/labstack/echo"
)

const (
	E2eTestsElement           = "tests" + obsidian.UrlSep + "e2e"
	EnodebdElement            = "enodebd"
	TestPkArg                 = ":test_pk"
	E2ETestsRootPath          = obsidian.V1Root + E2eTestsElement
	E2ETestsEnodebdRootPath   = E2ETestsRootPath + obsidian.UrlSep + EnodebdElement
	E2ETestsEnodebdDetailPath = E2ETestsEnodebdRootPath + obsidian.UrlSep + TestPkArg
)

func GetObsidianHandlers() []obsidian.Handler {
	return []obsidian.Handler{
		{Path: CINodesRootPath, Methods: obsidian.GET, HandlerFunc: listCINodes},
		{Path: CINodesRootPath, Methods: obsidian.POST, HandlerFunc: createCINode},
		{Path: CINodesGetPath, Methods: obsidian.GET, HandlerFunc: getCINode},
		{Path: CINodesGetPath, Methods: obsidian.PUT, HandlerFunc: updateCINode},
		{Path: CINodesGetPath, Methods: obsidian.DELETE, HandlerFunc: deleteCINode},
		{Path: CINodesReservePath, Methods: obsidian.POST, HandlerFunc: leaseCINode},
		{Path: CINodesManuallyReservePath, Methods: obsidian.POST, HandlerFunc: reserveCINode},
		{Path: CINodesManuallyReleasePath, Methods: obsidian.POST, HandlerFunc: returnManuallyReservedCINode},
		{Path: CINodesReleasePath, Methods: obsidian.POST, HandlerFunc: releaseCINode},

		{Path: E2ETestsRootPath, Methods: obsidian.GET, HandlerFunc: listTestCases},
		{Path: E2ETestsEnodebdRootPath, Methods: obsidian.GET, HandlerFunc: listEnodebdTestCase},
		{Path: E2ETestsEnodebdRootPath, Methods: obsidian.POST, HandlerFunc: createEnodebdTestCase},
		{Path: E2ETestsEnodebdDetailPath, Methods: obsidian.GET, HandlerFunc: getEnodebdTestCase},
		{Path: E2ETestsEnodebdDetailPath, Methods: obsidian.PUT, HandlerFunc: updateEnodebdTestCase},
		{Path: E2ETestsEnodebdDetailPath, Methods: obsidian.DELETE, HandlerFunc: deleteEnodebdTestCase},
	}
}

func listTestCases(c echo.Context) error {
	tcs, err := testcontroller.GetTestCases(nil, serdes.TestController)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := make([]*models.E2eTestCase, 0, len(tcs))
	for _, tc := range tcs {
		ret = append(ret, unmarshalledTestCaseToModel(tc))
	}
	sort.Slice(ret, func(i, j int) bool { return *ret[i].Pk < *ret[j].Pk })
	return c.JSON(http.StatusOK, ret)
}

func listEnodebdTestCase(c echo.Context) error {
	// We should add some filter criteria to the RPC method but this will be
	// low-scale enough that it won't matter
	tcs, err := testcontroller.GetTestCases(nil, serdes.TestController)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := []*models.EnodebdE2eTest{}
	for _, tc := range tcs {
		if tc.TestCaseType != testcontroller.EnodedTestCaseType && tc.TestCaseType != testcontroller.EnodedTestExcludeTraffic {
			continue
		}
		ret = append(ret, unmarshalledTestCaseToEnodebdTestCase(tc))
	}
	sort.Slice(ret, func(i, j int) bool { return *ret[i].Pk < *ret[j].Pk })
	return c.JSON(http.StatusOK, ret)
}

func createEnodebdTestCase(e echo.Context) error {
	tc := &models.MutableEnodebdE2eTest{}
	if err := e.Bind(tc); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := tc.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	runTraffic := *(tc.Config).RunTrafficTests
	enodedTestType := testcontroller.EnodedTestCaseType
	if !runTraffic {
		enodedTestType = testcontroller.EnodedTestExcludeTraffic
	}

	err := testcontroller.CreateOrUpdateTestCase(*tc.Pk, enodedTestType, tc.Config, serdes.TestController)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return e.NoContent(http.StatusCreated)
}

func getEnodebdTestCase(e echo.Context) error {
	pk, nerr := getTestPk(e)
	if nerr != nil {
		return nerr
	}

	res, err := testcontroller.GetTestCases([]int64{pk}, serdes.TestController)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	ret, ok := res[pk]
	if !ok || (ret.TestCaseType != testcontroller.EnodedTestCaseType && ret.TestCaseType != testcontroller.EnodedTestExcludeTraffic) {
		return echo.ErrNotFound
	}
	return e.JSON(http.StatusOK, unmarshalledTestCaseToEnodebdTestCase(ret))
}

func updateEnodebdTestCase(e echo.Context) error {
	pk, nerr := getTestPk(e)
	if nerr != nil {
		return nerr
	}
	cfg := &models.EnodebdTestConfig{}
	if err := e.Bind(cfg); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := cfg.Validate(strfmt.Default); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	runTraffic := *cfg.RunTrafficTests
	enodedTestType := testcontroller.EnodedTestCaseType
	if !runTraffic {
		enodedTestType = testcontroller.EnodedTestExcludeTraffic
	}
	err := testcontroller.CreateOrUpdateTestCase(pk, enodedTestType, cfg, serdes.TestController)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return e.NoContent(http.StatusNoContent)
}

func deleteEnodebdTestCase(e echo.Context) error {
	pk, nerr := getTestPk(e)
	if nerr != nil {
		return nerr
	}

	err := testcontroller.DeleteTestCase(pk)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return e.NoContent(http.StatusNoContent)
}

func unmarshalledTestCaseToModel(tc *testcontroller.UnmarshalledTestCase) *models.E2eTestCase {
	ret := &models.E2eTestCase{
		Config: tc.UnmarshaledConfig,
		Pk:     swag.Int64(tc.Pk),
		State: &models.E2eTestCaseState{
			CurrentState: tc.State,
			Error:        tc.Error,
			IsExecuting:  swag.Bool(tc.IsCurrentlyExecuting),
		},
		TestType: swag.String(tc.TestCaseType),
	}

	lastTime := timestampToStrfmt(tc.LastExecutionTime)
	if lastTime != nil {
		ret.State.LastExecutionTime = *lastTime
	}
	nextTime := timestampToStrfmt(tc.NextScheduledTime)
	if nextTime != nil {
		ret.State.NextScheduledTime = *nextTime
	}
	return ret
}

func unmarshalledTestCaseToEnodebdTestCase(tc *testcontroller.UnmarshalledTestCase) *models.EnodebdE2eTest {
	genericTc := unmarshalledTestCaseToModel(tc)
	return &models.EnodebdE2eTest{
		Config: genericTc.Config.(*models.EnodebdTestConfig),
		Pk:     genericTc.Pk,
		State:  genericTc.State,
	}
}

func timestampToStrfmt(ts *timestamp.Timestamp) *strfmt.DateTime {
	if ts == nil {
		return nil
	}

	tt, err := ptypes.Timestamp(ts)
	if err != nil {
		glog.Errorf("timestamp failed validation: %s", err)
		return nil
	}
	ret := strfmt.DateTime(tt)
	return &ret
}

func getTestPk(e echo.Context) (int64, *echo.HTTPError) {
	params, nerr := obsidian.GetParamValues(e, "test_pk")
	if nerr != nil {
		return 0, nerr
	}
	i, err := strconv.ParseInt(params[0], 10, 64)
	if err != nil {
		return 0, obsidian.HttpError(err, http.StatusBadRequest)
	}
	return i, nil
}
