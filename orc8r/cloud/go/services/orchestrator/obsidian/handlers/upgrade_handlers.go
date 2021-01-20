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
	"errors"
	"fmt"
	"net/http"
	"sort"

	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
	"github.com/labstack/echo"
)

func listChannelsHandler(c echo.Context) error {
	channelNames, err := configurator.ListInternalEntityKeys(orc8r.UpgradeReleaseChannelEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(channelNames)
	return c.JSON(http.StatusOK, channelNames)
}

func createChannelHandler(c echo.Context) error {
	payload, nerr := GetAndValidatePayload(c, &models.ReleaseChannel{})
	if nerr != nil {
		return nerr
	}
	channel := payload.(*models.ReleaseChannel)

	entity := configurator.NetworkEntity{
		Type:   orc8r.UpgradeReleaseChannelEntityType,
		Key:    string(channel.ID),
		Name:   channel.Name,
		Config: channel,
	}
	_, err := configurator.CreateInternalEntity(entity, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func readChannelHandler(c echo.Context) error {
	channelID, nerr := getChannelID(c)
	if nerr != nil {
		return nerr
	}
	entity, err := configurator.LoadInternalEntity(
		orc8r.UpgradeReleaseChannelEntityType, channelID,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, entity.Config)
}

func updateChannelHandler(c echo.Context) error {
	channelID, nerr := getChannelID(c)
	if nerr != nil {
		return nerr
	}

	payload, nerr := GetAndValidatePayload(c, &models.ReleaseChannel{})
	if nerr != nil {
		return nerr
	}
	channel := payload.(*models.ReleaseChannel)

	update := configurator.EntityUpdateCriteria{
		Type:      orc8r.UpgradeReleaseChannelEntityType,
		Key:       channelID,
		NewName:   swag.String(channel.Name),
		NewConfig: channel,
	}
	_, err := configurator.UpdateInternalEntity(update, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteChannelHandler(c echo.Context) error {
	channelID, nerr := getChannelID(c)
	if nerr != nil {
		return nerr
	}
	err := configurator.DeleteInternalEntity(orc8r.UpgradeReleaseChannelEntityType, channelID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listTiersHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	tiers, err := configurator.ListEntityKeys(networkID, orc8r.UpgradeTierEntityType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	sort.Strings(tiers)
	return c.JSON(http.StatusOK, tiers)
}

func createTierHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	payload, nerr := GetAndValidatePayload(c, &models.Tier{})
	if nerr != nil {
		return nerr
	}
	tier := payload.(*models.Tier)
	entity := tier.ToNetworkEntity()
	_, err := configurator.CreateEntity(networkID, entity, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusCreated)
}

func updateTierHandler(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	payload, nerr := GetAndValidatePayload(c, &models.Tier{})
	if nerr != nil {
		return nerr
	}
	tier := payload.(*models.Tier)

	if string(tier.ID) != tierID {
		return obsidian.HttpError(fmt.Errorf("TierID in URL and payload do not match."), http.StatusBadRequest)
	}
	update := tier.ToUpdateCriteria()
	_, err := configurator.UpdateEntity(networkID, update, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func readTierHandler(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	entity, err := configurator.LoadEntity(
		networkID, orc8r.UpgradeTierEntityType, tierID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true, LoadMetadata: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	tier := &models.Tier{}
	return c.JSON(http.StatusOK, tier.FromBackendModel(entity))
}

func deleteTierHandler(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	err := configurator.DeleteEntity(networkID, orc8r.UpgradeTierEntityType, tierID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func createTierImage(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	image, nerr := GetAndValidatePayload(c, &models.TierImage{})
	if nerr != nil {
		return nerr
	}

	updates, err := image.(*models.TierImage).ToUpdateCriteria(networkID, tierID)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	_, err = configurator.UpdateEntities(networkID, updates, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteImage(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	params, nerr := obsidian.GetParamValues(c, "image_name")
	if nerr != nil {
		return nerr
	}
	update, err := (&models.TierImage{}).ToDeleteImageUpdateCriteria(networkID, tierID, params[0])
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	_, err = configurator.UpdateEntity(networkID, update, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func createTierGateway(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	var gatewayID string
	if err := c.Bind(&gatewayID); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	update := (&models.TierGateways{}).ToAddGatewayUpdateCriteria(tierID, gatewayID)
	_, err := configurator.UpdateEntity(networkID, update, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteTierGateway(c echo.Context) error {
	networkID, tierID, nerr := getNetworkAndTierIDs(c)
	if nerr != nil {
		return nerr
	}
	_, gatewayID, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	update := (&models.TierGateways{}).ToDeleteGatewayUpdateCriteria(tierID, gatewayID)
	_, err := configurator.UpdateEntity(networkID, update, serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getChannelID(c echo.Context) (string, *echo.HTTPError) {
	channelID := c.Param("channel_id")
	if channelID == "" {
		return "", obsidian.HttpError(errors.New("Missing release channel ID"), http.StatusBadRequest)
	}
	return channelID, nil
}

func getNetworkAndTierIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "tier_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}
