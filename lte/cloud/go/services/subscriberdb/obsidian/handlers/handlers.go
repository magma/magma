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
	"regexp"

	"magma/lte/cloud/go/lte"
	ltehandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/golang/glog"
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

const (
	mobilitydStateExpectedMatchCount = 2
)

var (
	// mobilitydStateKeyRe captures the IMSI portion of mobilityd state keys.
	// Mobilityd states are keyed as <IMSI>.<APN>.
	mobilitydStateKeyRe = regexp.MustCompile(`^(?P<imsi>IMSI\d+)\..+$`)

	subscriberLoadCriteria = configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true}
)

var subscriberStateTypes = []string{
	lte.ICMPStateType,
	lte.S1APStateType,
	lte.MMEStateType,
	lte.SPGWStateType,
	lte.MobilitydStateType,
	orc8r.DirectoryRecordType,
}

func listSubscribers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	subs, err := loadAllSubscribers(networkID)
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
	subs, err := loadSubscriber(networkID, subscriberID)
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

func loadSubscriber(networkID, key string) (*subscribermodels.Subscriber, error) {
	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, key, subscriberLoadCriteria)
	if err != nil {
		return nil, err
	}
	mutableSub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent)
	if err != nil {
		return nil, err
	}
	sub := mutableSub.ToSubscriber()

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, &key)
	if err != nil {
		return nil, err
	}
	err = sub.FillAugmentedFields(states)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func loadAllSubscribers(networkID string) (map[string]*subscribermodels.Subscriber, error) {
	mutableSubs, err := loadAllMutableSubscribers(networkID)
	if err != nil {
		return nil, err
	}

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil)
	if err != nil {
		return nil, err
	}
	// Each entry in this map contains all the states that the SID cares about.
	// The DeviceID fields of the state IDs in the nested maps do not have to
	// match the SID, as in the case of mobilityd state for example.
	statesBySid := map[string]state_types.StatesByID{}
	for stateID, st := range states {
		sidKey := stateID.DeviceID
		if stateID.Type == lte.MobilitydStateType {
			matches := mobilitydStateKeyRe.FindStringSubmatch(stateID.DeviceID)
			if len(matches) != mobilitydStateExpectedMatchCount {
				glog.Errorf("mobilityd state composite ID %s did not match regex", sidKey)
				continue
			}
			sidKey = matches[1]
		}

		if _, exists := statesBySid[sidKey]; !exists {
			statesBySid[sidKey] = state_types.StatesByID{}
		}
		statesBySid[sidKey][stateID] = st
	}

	subs := map[string]*subscribermodels.Subscriber{}
	for _, mutableSub := range mutableSubs {
		sub := mutableSub.ToSubscriber()
		err = sub.FillAugmentedFields(statesBySid[string(sub.ID)])
		if err != nil {
			return nil, err
		}
		subs[string(sub.ID)] = sub
	}

	return subs, nil
}

func loadAllMutableSubscribers(networkID string) (map[string]*subscribermodels.MutableSubscriber, error) {
	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, subscriberLoadCriteria)
	if err != nil {
		return nil, err
	}
	subs := map[string]*subscribermodels.MutableSubscriber{}
	for _, ent := range ents {
		sub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent)
		if err != nil {
			return nil, err
		}
		subs[ent.Key] = sub
	}
	return subs, nil
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return obsidian.HttpError(err, http.StatusInternalServerError)
}
