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
	"fmt"
	"net/http"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	"magma/orc8r/cloud/go/services/ctraced/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	tracing = "tracing"
	// v1/networks/:network_id/tracing
	tracingRootPath = obsidian.V1Root + obsidian.MagmaNetworksUrlPart + obsidian.UrlSep + ":" + pathParamNetworkID + obsidian.UrlSep + tracing
	// v1/networks/:network_id/tracing/:trace_id
	tracingPath = tracingRootPath + obsidian.UrlSep + ":" + pathParamTraceID
	// v1/networks/:network_id/tracing/:trace_id/download
	tracingDownloadPath = tracingPath + obsidian.UrlSep + "download"

	pathParamTraceID   = "trace_id"
	pathParamNetworkID = "network_id"
)

func GetObsidianHandlers(client GwCtracedClient, storage storage.CtracedStorage) []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: tracingRootPath, Methods: obsidian.GET, HandlerFunc: listCallTraces},
		{Path: tracingRootPath, Methods: obsidian.POST, HandlerFunc: getCreateCallTraceHandlerFunc(client)},
		{Path: tracingPath, Methods: obsidian.GET, HandlerFunc: getCallTrace},
		{Path: tracingPath, Methods: obsidian.PUT, HandlerFunc: getUpdateCallTraceHandlerFunc(client, storage)},
		{Path: tracingPath, Methods: obsidian.DELETE, HandlerFunc: getDeleteCallTraceHandlerFunc(client, storage)},
		{Path: tracingDownloadPath, Methods: obsidian.GET, HandlerFunc: getDownloadCallTraceHandlerFunc(storage)},
	}

	return ret
}

func listCallTraces(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	callTraces, _, err := configurator.LoadAllEntitiesOfType(
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

func getCreateCallTraceHandlerFunc(client GwCtracedClient) echo.HandlerFunc {
	return func(c echo.Context) error {
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
		exists, err := configurator.DoesEntityExist(networkID, orc8r.CallTraceEntityType, cfg.TraceID)
		if exists {
			return obsidian.HttpError(errors.New(fmt.Sprintf("Call trace id: %s already exists", cfg.TraceID)))
		}
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		if err := ctr.ValidateModel(); err != nil {
			return obsidian.HttpError(err, http.StatusBadRequest)
		}

		req, err := buildStartTraceRequest(cfg)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to build call trace request"), http.StatusInternalServerError)
		}

		resp, err := client.StartCallTrace(networkID, cfg.GatewayID, req)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to start call trace"), http.StatusInternalServerError)
		}
		if !resp.Success {
			return obsidian.HttpError(errors.New("failed to start call trace"), http.StatusInternalServerError)
		}

		createdEntity := ctr.ToEntity()
		_, err = configurator.CreateEntity(networkID, createdEntity, serdes.Entity)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to create call trace"), http.StatusInternalServerError)
		}
		return c.JSON(http.StatusCreated, cfg.TraceID)
	}
}

func getCallTrace(c echo.Context) error {
	callTrace, err := getCallTraceModel(c)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, callTrace)
}

func getUpdateCallTraceHandlerFunc(client GwCtracedClient, storage storage.CtracedStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
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

		callTrace, err := getCallTraceModel(c)
		if err != nil {
			return err
		}
		if !shouldEndTraceBeTriggered(callTrace, mutableCallTrace) {
			return obsidian.HttpError(errors.New("Error: call trace end already triggered earlier"), http.StatusBadRequest)
		}

		req := &protos.EndTraceRequest{
			TraceId: callTraceID,
		}
		resp, err := client.EndCallTrace(networkID, callTrace.Config.GatewayID, req)
		if err != nil {
			return err
		}
		if !resp.Success {
			return obsidian.HttpError(errors.New("Failed to end call trace"), http.StatusInternalServerError)
		}

		err = storage.StoreCallTrace(networkID, callTraceID, resp.TraceContent)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, fmt.Sprintf("failed to save call trace data, network-id: %s, gateway-id: %s, calltrace-id: %s", networkID, callTrace.Config.GatewayID, callTraceID)), http.StatusInternalServerError)
		}

		_, err = configurator.UpdateEntity(networkID, mutableCallTrace.ToEntityUpdateCriteria(callTraceID, *callTrace), serdes.Entity)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func getDeleteCallTraceHandlerFunc(client GwCtracedClient, storage storage.CtracedStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkID, callTraceID, nerr := getNetworkIDAndCallTraceID(c)
		if nerr != nil {
			return nerr
		}

		err := storage.DeleteCallTrace(networkID, callTraceID)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to delete call trace data"), http.StatusInternalServerError)
		}

		err = configurator.DeleteEntity(networkID, orc8r.CallTraceEntityType, callTraceID)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusNoContent)
	}
}

func getDownloadCallTraceHandlerFunc(storage storage.CtracedStorage) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkID, callTraceID, nerr := getNetworkIDAndCallTraceID(c)
		if nerr != nil {
			return nerr
		}

		callTrace, err := storage.GetCallTrace(networkID, callTraceID)
		if err != nil {
			return obsidian.HttpError(errors.Wrap(err, "failed to retrieve call trace data"), http.StatusInternalServerError)
		}

		res := c.Response()
		header := res.Header()
		header.Set(echo.HeaderContentType, "application/pcapng")
		header.Set(echo.HeaderContentDisposition, "attachment; filename="+fmt.Sprintf("%s.pcapng", callTraceID))
		res.WriteHeader(http.StatusOK)
		_, err = res.Write(callTrace)
		return err
	}
}

func getCallTraceModel(c echo.Context) (*models.CallTrace, error) {
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

func buildStartTraceRequest(cfg *models.CallTraceConfig) (*protos.StartTraceRequest, error) {
	req := &protos.StartTraceRequest{
		TraceId:        cfg.TraceID,
		TraceType:      protos.StartTraceRequest_ALL,
		Timeout:        cfg.Timeout,
		CaptureFilters: cfg.CaptureFilters,
		DisplayFilters: cfg.DisplayFilters,
	}
	return req, nil
}

func shouldEndTraceBeTriggered(callTrace *models.CallTrace, mutable *models.MutableCallTrace) bool {
	if callTrace.State.CallTraceEnding {
		return false
	}
	return *mutable.RequestedEnd
}
