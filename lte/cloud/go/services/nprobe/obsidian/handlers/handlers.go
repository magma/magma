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
	"math/rand"
	"net/http"
	"time"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/lte/obsidian/handlers"
	"magma/lte/cloud/go/services/nprobe/obsidian/models"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	NetworkProbeTasks                        = "network_probe_tasks"
	ManageNetworkNetworkProbeTasksPath       = handlers.ManageNetworkPath + obsidian.UrlSep + NetworkProbeTasks
	ManageNetworkNetworkProbeTaskDetailsPath = ManageNetworkNetworkProbeTasksPath + obsidian.UrlSep + ":task_id"

	NetworkProbeDestinations                        = "network_probe_destinations"
	ManageNetworkNetworkProbeDestinationsPath       = handlers.ManageNetworkPath + obsidian.UrlSep + NetworkProbeDestinations
	ManageNetworkNetworkProbeDestinationDetailsPath = ManageNetworkNetworkProbeDestinationsPath + obsidian.UrlSep + ":destination_id"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ManageNetworkNetworkProbeTasksPath, Methods: obsidian.GET, HandlerFunc: listNetworkProbeTasks},
		{Path: ManageNetworkNetworkProbeTasksPath, Methods: obsidian.POST, HandlerFunc: createNetworkProbeTask},
		{Path: ManageNetworkNetworkProbeTaskDetailsPath, Methods: obsidian.GET, HandlerFunc: getNetworkProbeTask},
		{Path: ManageNetworkNetworkProbeTaskDetailsPath, Methods: obsidian.PUT, HandlerFunc: updateNetworkProbeTask},
		{Path: ManageNetworkNetworkProbeTaskDetailsPath, Methods: obsidian.DELETE, HandlerFunc: deleteNetworkProbeTask},

		{Path: ManageNetworkNetworkProbeDestinationsPath, Methods: obsidian.GET, HandlerFunc: listNetworkProbeDestinations},
		{Path: ManageNetworkNetworkProbeDestinationsPath, Methods: obsidian.POST, HandlerFunc: createNetworkProbeDestination},
		{Path: ManageNetworkNetworkProbeDestinationDetailsPath, Methods: obsidian.GET, HandlerFunc: getNetworkProbeDestination},
		{Path: ManageNetworkNetworkProbeDestinationDetailsPath, Methods: obsidian.PUT, HandlerFunc: updateNetworkProbeDestination},
		{Path: ManageNetworkNetworkProbeDestinationDetailsPath, Methods: obsidian.DELETE, HandlerFunc: deleteNetworkProbeDestination},
	}
	return ret
}

func getParamValues(c echo.Context, paramNames []string) ([]string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, paramNames...)
	if err != nil {
		return []string{}, err
	}
	return vals, nil
}

func getNetworkProbeEntityAndIDs(c echo.Context, entityType string) (configurator.NetworkEntity, string, string, error) {
	ent := configurator.NetworkEntity{}
	paramNames := []string{"network_id"}
	switch {
	case entityType == lte.NetworkProbeTaskEntityType:
		paramNames = append(paramNames, "task_id")
	case entityType == lte.NetworkProbeDestinationEntityType:
		paramNames = append(paramNames, "destination_id")
	}

	values, err1 := getParamValues(c, paramNames)
	if err1 != nil {
		return ent, "", "", err1
	}
	networkID, taskID := values[0], values[1]
	ent, err2 := configurator.LoadEntity(networkID, entityType, taskID, configurator.EntityLoadCriteria{LoadConfig: true}, serdes.Entity)
	return ent, networkID, taskID, err2
}

func getNetworkProbeTask(c echo.Context) error {
	ent, _, _, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeTaskEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&models.NetworkProbeTask{}).FromBackendModels(ent)
	return c.JSON(http.StatusOK, ret)
}

func deleteNetworkProbeTask(c echo.Context) error {
	ent, networkID, _, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeTaskEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	var deletes []storage.TypeAndKey
	deletes = append(deletes, ent.GetTypeAndKey())
	err = configurator.DeleteEntities(networkID, deletes)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getNetworkProbeDestination(c echo.Context) error {
	ent, _, _, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeDestinationEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&models.NetworkProbeDestination{}).FromBackendModels(ent)
	return c.JSON(http.StatusOK, ret)
}

func deleteNetworkProbeDestination(c echo.Context) error {
	ent, networkID, _, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeDestinationEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	var deletes []storage.TypeAndKey
	deletes = append(deletes, ent.GetTypeAndKey())
	err = configurator.DeleteEntities(networkID, deletes)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateNetworkProbeTask(c echo.Context) error {
	_, networkID, taskID, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeTaskEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load existing NetworkProbeTask"), http.StatusInternalServerError)
	}

	payload := &models.NetworkProbeTask{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	err = configurator.CreateOrUpdateEntityConfig(networkID, lte.NetworkProbeTaskEntityType, taskID, payload.TaskDetails, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateNetworkProbeDestination(c echo.Context) error {
	_, networkID, taskID, err := getNetworkProbeEntityAndIDs(c, lte.NetworkProbeDestinationEntityType)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load existing NetworkProbeDestination"), http.StatusInternalServerError)
	}

	payload := &models.NetworkProbeDestination{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	err = configurator.CreateOrUpdateEntityConfig(networkID, lte.NetworkProbeDestinationEntityType, taskID, payload.DestinationDetails, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listNetworkProbeTasks(c echo.Context) error {
	networkID, err1 := obsidian.GetNetworkId(c)
	if err1 != nil {
		return err1
	}

	ents, err2 := configurator.LoadAllEntitiesOfType(
		networkID, lte.NetworkProbeTaskEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	switch {
	case err2 == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err2 != nil:
		return obsidian.HttpError(errors.Wrap(err2, "failed to load existing NetworkProbeTasks"), http.StatusInternalServerError)
	}

	ret := make(map[string]*models.NetworkProbeTask, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&models.NetworkProbeTask{}).FromBackendModels(ent)
	}
	return c.JSON(http.StatusOK, ret)
}

func listNetworkProbeDestinations(c echo.Context) error {
	networkID, err1 := obsidian.GetNetworkId(c)
	if err1 != nil {
		return err1
	}

	ents, err2 := configurator.LoadAllEntitiesOfType(
		networkID, lte.NetworkProbeDestinationEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err2 != nil {
		return obsidian.HttpError(err2, http.StatusInternalServerError)
	}

	ret := make(map[string]*models.NetworkProbeDestination, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&models.NetworkProbeDestination{}).FromBackendModels(ent)
	}
	return c.JSON(http.StatusOK, ret)
}

func createNetworkProbeTask(c echo.Context) error {
	networkID, err1 := obsidian.GetNetworkId(c)
	if err1 != nil {
		return err1
	}

	payload := &models.NetworkProbeTask{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	// generate random correlation ID if not provided
	if payload.TaskDetails.CorrelationID == 0 {
		payload.TaskDetails.CorrelationID = rand.Uint64()
	}

	if time.Time(payload.TaskDetails.Timestamp).IsZero() {
		payload.TaskDetails.Timestamp = strfmt.DateTime(time.Now().UTC())
	}

	_, err2 := configurator.CreateEntity(
		networkID,
		configurator.NetworkEntity{
			Type:   lte.NetworkProbeTaskEntityType,
			Key:    string(payload.TaskID),
			Config: payload.TaskDetails,
		},
		serdes.Entity,
	)
	if err2 != nil {
		return obsidian.HttpError(err2, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func createNetworkProbeDestination(c echo.Context) error {
	networkID, err1 := obsidian.GetNetworkId(c)
	if err1 != nil {
		return err1
	}

	payload := &models.NetworkProbeDestination{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	_, err2 := configurator.CreateEntity(
		networkID,
		configurator.NetworkEntity{
			Type:   lte.NetworkProbeDestinationEntityType,
			Key:    string(payload.DestinationID),
			Config: payload.DestinationDetails,
		},
		serdes.Entity,
	)
	if err2 != nil {
		return obsidian.HttpError(err2, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}
