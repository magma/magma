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

	"github.com/go-openapi/swag"
	"github.com/labstack/echo/v4"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	"magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/obsidian"
	"magma/orc8r/lib/go/merrors"
)

const (
	ratingGroupIDParam = "rating_group_id"
)

func ListRatingGroups(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, _, err := configurator.LoadAllEntitiesOfType(
		c.Request().Context(),
		networkID, lte.RatingGroupEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	groupsByID := map[models.RatingGroupID]*models.RatingGroup{}
	for _, ent := range ents {
		r := (&models.RatingGroup{}).FromEntity(ent)
		groupsByID[r.ID] = r
	}
	return c.JSON(http.StatusOK, groupsByID)
}

func CreateRatingGroup(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	reqCtx := c.Request().Context()

	group := new(models.RatingGroup)
	if err := c.Bind(group); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := group.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	_, err := configurator.CreateEntity(reqCtx, networkID, group.ToEntity(), serdes.Entity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusCreated)
}

func GetRatingGroup(c echo.Context) error {
	networkID, ratingGroupID, nerr := getNetworkAndGroupIDs(c)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		c.Request().Context(),
		networkID, lte.RatingGroupEntityType, ratingGroupID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, (&models.RatingGroup{}).FromEntity(ent))
}

func UpdateRatingGroup(c echo.Context) error {
	networkID, ratingGroupID, nerr := getNetworkAndGroupIDs(c)
	if nerr != nil {
		return nerr
	}
	reqCtx := c.Request().Context()

	ratingGroup := new(models.MutableRatingGroup)
	if err := c.Bind(ratingGroup); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	if err := ratingGroup.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	groupID, err := swag.ConvertUint32(ratingGroupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	// 404 if rating group doesn't exist
	exists, err := configurator.DoesEntityExist(reqCtx, networkID, lte.RatingGroupEntityType, ratingGroupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to check if rating group exists: %v", err))
	}
	if !exists {
		return echo.ErrNotFound
	}

	_, err = configurator.UpdateEntity(reqCtx, networkID, ratingGroup.ToEntityUpdateCriteria(groupID), serdes.Entity)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func DeleteRatingGroup(c echo.Context) error {
	networkID, ratingGroupID, nerr := getNetworkAndGroupIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(c.Request().Context(), networkID, lte.RatingGroupEntityType, ratingGroupID)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.NoContent(http.StatusNoContent)
}

func getNetworkAndGroupIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", ratingGroupIDParam)
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
