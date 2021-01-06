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

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	Tracing = "tracing"
	// v1/networks/:network_id/tracing
	TracingRootPath = obsidian.V1Root + obsidian.MagmaNetworksUrlPart + obsidian.UrlSep + ":" + pathParamNetworkID + obsidian.UrlSep + Tracing
	// v1/networks/:network_id/tracing/:trace_id
	TracingPath = TracingRootPath + obsidian.UrlSep + ":" + pathParamTraceID

	pathParamTraceID   = "trace_id"
	pathParamNetworkID = "network_id"
)

func GetObsidianHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: TracingRootPath, Methods: obsidian.GET, HandlerFunc: ListCallTraces},
		{Path: TracingRootPath, Methods: obsidian.POST, HandlerFunc: CreateCallTrace},
		{Path: TracingPath, Methods: obsidian.GET, HandlerFunc: GetCallTrace},
		{Path: TracingPath, Methods: obsidian.PUT, HandlerFunc: UpdateCallTrace},
		{Path: TracingPath, Methods: obsidian.DELETE, HandlerFunc: DeleteCallTrace},
	}

	return ret
}

func ListCallTraces(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	callTraces, err := configurator.LoadAllEntitiesOfType(
		networkID, orc8r.CallTraceEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := map[string]*models.CallTrace{}
	for _, ctEnt := range callTraces {
		ret[ctEnt.Key] = (&models.CallTrace{}).FromEntity(ctEnt)
	}
	return c.JSON(http.StatusOK, ret)
}

func CreateCallTrace(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	cfg := &models.CallTraceConfig{}
	if err := c.Bind(cfg); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	ctr := &models.CallTrace{
		Config: cfg,
		State: &models.CallTraceState{
			CallTraceAvailable: false,
			CallTraceEnding:    false,
		},
	}

	if err := ctr.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	createdEntity := ctr.ToEntity()
	_, err := configurator.CreateEntity(networkID, createdEntity, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to create call trace"), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, string(cfg.TraceID))
}

func GetCallTrace(c echo.Context) error {
	callTrace, err := getCallTrace(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, callTrace)
}

func UpdateCallTrace(c echo.Context) error {
	networkID, callTraceID, nerr := getNetworkIDAndCallTraceID(c)
	if nerr != nil {
		return nerr
	}
	mutableCallTrace := &models.MutableCallTrace{}
	if err := c.Bind(mutableCallTrace); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := mutableCallTrace.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	callTrace, err := getCallTrace(c)
	if err != nil {
		return err
	}

	_, err = configurator.UpdateEntity(networkID, mutableCallTrace.ToEntityUpdateCriteria(callTraceID, *callTrace), serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteCallTrace(c echo.Context) error {
	networkID, callTraceID, nerr := getNetworkIDAndCallTraceID(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, orc8r.CallTraceEntityType, callTraceID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getCallTrace(c echo.Context) (*models.CallTrace, error) {
	networkID, callTraceID, nerr := getNetworkIDAndCallTraceID(c)
	if nerr != nil {
		return nil, nerr
	}
	ent, err := configurator.LoadEntity(
		networkID, orc8r.CallTraceEntityType, callTraceID,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return nil, obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	callTrace := &models.CallTrace{}
	err = callTrace.FromBackendModels(ent)
	if err != nil {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return callTrace, nil
}

func getNetworkIDAndCallTraceID(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, pathParamNetworkID, pathParamTraceID)
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
