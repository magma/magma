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

	"magma/lte/cloud/go/lte"
	ltehandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	Subscribers              = "subscribers"
	ListSubscribersPath      = ltehandlers.ManageNetworkPath + obsidian.UrlSep + Subscribers
	ManageSubscriberPath     = ListSubscribersPath + obsidian.UrlSep + ":subscriber_id"
	ActivateSubscriberPath   = ManageSubscriberPath + obsidian.UrlSep + "activate"
	DeactivateSubscriberPath = ManageSubscriberPath + obsidian.UrlSep + "deactivate"
	SubscriberProfilePath    = ManageSubscriberPath + obsidian.UrlSep + "lte" + obsidian.UrlSep + "sub_profile"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListSubscribersPath, Methods: obsidian.GET, HandlerFunc: listSubscribers},
		{Path: ListSubscribersPath, Methods: obsidian.POST, HandlerFunc: createSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.GET, HandlerFunc: getSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriber},
		{Path: ManageSubscriberPath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriber},
		{Path: ActivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateACTIVE)},
		{Path: DeactivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateINACTIVE)},
		{Path: SubscriberProfilePath, Methods: obsidian.PUT, HandlerFunc: updateSubscriberProfile},
	}
	return ret
}

func listSubscribers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subs, err := (&subscribermodels.Subscriber{}).LoadAll(networkID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, subs)
}

func createSubscriber(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &subscribermodels.MutableSubscriber{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if nerr := validateSubscriberProfile(networkID, payload.Lte); nerr != nil {
		return nerr
	}

	err := payload.Create(networkID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}
	subs, err := (&subscribermodels.Subscriber{}).Load(networkID, subscriberID)
	if err != nil {
		return makeErr(err)
	}
	return c.JSON(http.StatusOK, subs)
}

func updateSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	payload := &subscribermodels.MutableSubscriber{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if string(payload.ID) != subscriberID {
		err := fmt.Errorf("subscriber ID from parameters (%s) and payload (%s) must match", subscriberID, payload.ID)
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	if nerr := validateSubscriberProfile(networkID, payload.Lte); nerr != nil {
		return nerr
	}

	err := payload.Update(networkID)
	if err != nil {
		return makeErr(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func deleteSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}
	err := (&subscribermodels.MutableSubscriber{}).Delete(networkID, subscriberID)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNoContent)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func updateSubscriberProfile(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	var payload = new(subscribermodels.SubProfile)
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	currentCfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID)
	if err != nil {
		return makeErr(err)
	}

	desiredCfg := currentCfg.(*subscribermodels.SubscriberConfig)
	desiredCfg.Lte.SubProfile = *payload
	if nerr := validateSubscriberProfile(networkID, desiredCfg.Lte); nerr != nil {
		return nerr
	}

	_, err = configurator.UpdateEntity(networkID, configurator.EntityUpdateCriteria{Type: lte.SubscriberEntityType, Key: subscriberID, NewConfig: desiredCfg})
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to update profile"), http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func makeSubscriberStateHandler(desiredState string) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
		if nerr != nil {
			return nerr
		}

		cfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID)
		if err != nil {
			return makeErr(err)
		}

		newConfig := cfg.(*subscribermodels.SubscriberConfig)
		newConfig.Lte.State = desiredState
		err = configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, newConfig)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		return c.NoContent(http.StatusOK)
	}
}

func getNetworkAndSubIDs(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "subscriber_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func validateSubscriberProfile(networkID string, sub *subscribermodels.LteSubscription) *echo.HTTPError {
	// Check the sub profiles available on the network if sub profile is not
	// default (which is always available)
	if sub.SubProfile != "default" {
		netConf, err := configurator.LoadNetworkConfig(networkID, lte.CellularNetworkConfigType)
		switch {
		case err == merrors.ErrNotFound:
			return obsidian.HttpError(errors.New("no cellular config found for network"), http.StatusInternalServerError)
		case err != nil:
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}

		cellNetConf := netConf.(*ltemodels.NetworkCellularConfigs)
		profName := string(sub.SubProfile)
		if _, profileExists := cellNetConf.Epc.SubProfiles[profName]; !profileExists {
			return obsidian.HttpError(errors.Errorf("subscriber profile %s does not exist for the network", profName), http.StatusBadRequest)
		}
	}
	return nil
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return obsidian.HttpError(err, http.StatusInternalServerError)
}
