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
	"magma/lte/cloud/go/serdes"
	ltehandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
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

	listMSISDNsPath   = ltehandlers.ManageNetworkPath + obsidian.UrlSep + "msisdns"
	manageMSISDNsPath = listMSISDNsPath + obsidian.UrlSep + ":msisdn"

	ParamMSISDN = "msisdn"
	ParamIP     = "ip"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListSubscribersPath, Methods: obsidian.GET, HandlerFunc: listSubscribersHandler},
		{Path: ListSubscribersPath, Methods: obsidian.POST, HandlerFunc: createSubscriberHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.GET, HandlerFunc: getSubscriberHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriberHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriberHandler},

		{Path: ActivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateACTIVE)},
		{Path: DeactivateSubscriberPath, Methods: obsidian.POST, HandlerFunc: makeSubscriberStateHandler(subscribermodels.LteSubscriptionStateINACTIVE)},
		{Path: SubscriberProfilePath, Methods: obsidian.PUT, HandlerFunc: updateSubscriberProfile},

		{Path: listMSISDNsPath, Methods: obsidian.GET, HandlerFunc: listMSISDNsHandler},
		{Path: listMSISDNsPath, Methods: obsidian.POST, HandlerFunc: createMSISDNsHandler},
		{Path: manageMSISDNsPath, Methods: obsidian.GET, HandlerFunc: getMSISDNHandler},
		{Path: manageMSISDNsPath, Methods: obsidian.DELETE, HandlerFunc: deleteMSISDNHandler},
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

	subscriberLoadCriteria       = configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true}
	apnPolicyProfileLoadCriteria = configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true}
)

var subscriberStateTypes = []string{
	lte.ICMPStateType,
	lte.S1APStateType,
	lte.MMEStateType,
	lte.SPGWStateType,
	lte.MobilitydStateType,
	orc8r.DirectoryRecordType,
}

type subscriberFilter func(sub *subscribermodels.Subscriber) bool

func acceptAll(sub *subscribermodels.Subscriber) bool { return true }

// listSubscribersHandler handles the base subscriber endpoint.
// The returned subscribers can be filtered using the following query
// parameters
//	- msisdn
//	- ip
//
// The MSISDN parameter is config-based, and is enforced to be a unique
// identifier.
//
// The IP parameter is state-based, and not guaranteed to be unique. The
// IP->IMSI mapping is cached as the output of a mobilityd state indexer, then
// each reported subscriber is checked to ensure it actually is assigned the
// requested IP.
func listSubscribersHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	// First check for query params to filter by
	if msisdn := c.QueryParam(ParamMSISDN); msisdn != "" {
		queryIMSI, err := subscriberdb.GetIMSIForMSISDN(networkID, msisdn)
		if err != nil {
			return makeErr(err)
		}
		subs, err := loadSubscribers(networkID, acceptAll, queryIMSI)
		if err != nil {
			return makeErr(err)
		}
		return c.JSON(http.StatusOK, subs)
	}
	if ip := c.QueryParam(ParamIP); ip != "" {
		queryIMSIs, err := subscriberdb.GetIMSIsForIP(networkID, ip)
		if err != nil {
			return makeErr(err)
		}
		filter := func(sub *subscribermodels.Subscriber) bool { return sub.IsAssignedIP(ip) }
		subs, err := loadSubscribers(networkID, filter, queryIMSIs...)
		if err != nil {
			return makeErr(err)
		}
		return c.JSON(http.StatusOK, subs)
	}

	// Default to listing all subscribers
	subs, err := loadAllSubscribers(networkID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, subs)
}

func createSubscriberHandler(c echo.Context) error {
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

	err := createSubscriber(networkID, payload)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getSubscriberHandler(c echo.Context) error {
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

func updateSubscriberHandler(c echo.Context) error {
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

	err := updateSubscriber(networkID, payload)
	if err != nil {
		return makeErr(err)
	}

	return c.NoContent(http.StatusNoContent)
}

func deleteSubscriberHandler(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}
	err := deleteSubscriber(networkID, subscriberID)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNoContent)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func listMSISDNsHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	msisdns, err := subscriberdb.ListMSISDNs(networkID)
	if err != nil {
		return makeErr(err)
	}
	// Normalize for JSON output
	if msisdns == nil {
		msisdns = map[string]string{}
	}

	return c.JSON(http.StatusOK, msisdns)
}

func createMSISDNsHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := &subscribermodels.MsisdnAssignment{}
	if err := c.Bind(payload); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := payload.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}

	err := subscriberdb.SetIMSIForMSISDN(networkID, string(payload.Msisdn), string(payload.ID))
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return c.NoContent(http.StatusCreated)
}

func getMSISDNHandler(c echo.Context) error {
	networkID, msisdn, nerr := getNetworkAndMSISDN(c)
	if nerr != nil {
		return nerr
	}
	imsi, err := subscriberdb.GetIMSIForMSISDN(networkID, msisdn)
	if err != nil {
		return makeErr(err)
	}
	return c.JSON(http.StatusOK, imsi)
}

func deleteMSISDNHandler(c echo.Context) error {
	networkID, msisdn, nerr := getNetworkAndMSISDN(c)
	if nerr != nil {
		return nerr
	}

	err := subscriberdb.DeleteMSISDN(networkID, msisdn)
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

	currentCfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, serdes.Entity)
	if err != nil {
		return makeErr(err)
	}

	desiredCfg := currentCfg.(*subscribermodels.SubscriberConfig)
	desiredCfg.Lte.SubProfile = *payload
	if nerr := validateSubscriberProfile(networkID, desiredCfg.Lte); nerr != nil {
		return nerr
	}

	_, err = configurator.UpdateEntity(
		networkID,
		configurator.EntityUpdateCriteria{Type: lte.SubscriberEntityType, Key: subscriberID, NewConfig: desiredCfg},
		serdes.Entity,
	)
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

		cfg, err := configurator.LoadEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, serdes.Entity)
		if err != nil {
			return makeErr(err)
		}

		newConfig := cfg.(*subscribermodels.SubscriberConfig)
		newConfig.Lte.State = desiredState
		err = configurator.CreateOrUpdateEntityConfig(networkID, lte.SubscriberEntityType, subscriberID, newConfig, serdes.Entity)
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

func getNetworkAndMSISDN(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "msisdn")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func validateSubscriberProfile(networkID string, sub *subscribermodels.LteSubscription) *echo.HTTPError {
	// Check the sub profiles available on the network if sub profile is not
	// default (which is always available)
	if sub.SubProfile != "default" {
		netConf, err := configurator.LoadNetworkConfig(networkID, lte.CellularNetworkConfigType, serdes.Network)
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
	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, key, subscriberLoadCriteria, serdes.Entity)
	if err != nil {
		return nil, err
	}

	// Configurator doesn't currently support loading a specified subgraph,
	// so we have to load the subscriber and its apn_policy_profile ents in
	// separate calls.
	var policyProfileEnts configurator.NetworkEntities
	if ppAssocs := ent.Associations.Filter(lte.APNPolicyProfileEntityType); len(ppAssocs) != 0 {
		policyProfileEnts, _, err = configurator.LoadEntities(
			ent.NetworkID, nil, nil, nil,
			ppAssocs,
			apnPolicyProfileLoadCriteria,
			serdes.Entity,
		)
		if err != nil {
			return nil, err
		}
	}

	mutableSub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent, policyProfileEnts)
	if err != nil {
		return nil, err
	}

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, &key, serdes.State)
	if err != nil {
		return nil, err
	}

	sub := mutableSub.ToSubscriber()
	err = sub.FillAugmentedFields(states)
	if err != nil {
		return nil, err
	}

	return sub, nil
}

func loadSubscribers(networkID string, includeSub subscriberFilter, keys ...string) (map[string]*subscribermodels.Subscriber, error) {
	subs := map[string]*subscribermodels.Subscriber{}
	for _, key := range keys {
		sub, err := loadSubscriber(networkID, key)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading subscriber %s", key)
		}
		if includeSub(sub) {
			subs[string(sub.ID)] = sub
		}
	}
	return subs, nil
}

func loadAllSubscribers(networkID string) (map[string]*subscribermodels.Subscriber, error) {
	mutableSubs, err := loadAllMutableSubscribers(networkID)
	if err != nil {
		return nil, err
	}

	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil, serdes.State)
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
	ents, err := configurator.LoadAllEntitiesOfType(networkID, lte.SubscriberEntityType, subscriberLoadCriteria, serdes.Entity)
	if err != nil {
		return nil, err
	}
	profileEnts, err := configurator.LoadAllEntitiesOfType(
		networkID, lte.APNPolicyProfileEntityType,
		apnPolicyProfileLoadCriteria,
		serdes.Entity,
	)
	if err != nil {
		return nil, err
	}
	profileEntsBySub := profileEnts.MakeByParentTK()

	subs := map[string]*subscribermodels.MutableSubscriber{}
	for _, ent := range ents {
		sub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent, profileEntsBySub[ent.GetTypeAndKey()])
		if err != nil {
			return nil, err
		}
		subs[ent.Key] = sub
	}
	return subs, nil
}

func createSubscriber(networkID string, sub *subscribermodels.MutableSubscriber) error {
	// New ents
	//	- active_policies_by_apn
	//		- Assocs: policy_rule..., apn
	//	- subscriber
	//		- Assocs: active_policies_by_apn

	subEnt := configurator.NetworkEntity{
		Type: lte.SubscriberEntityType,
		Key:  string(sub.ID),
		Name: sub.Name,
		Config: &subscribermodels.SubscriberConfig{
			Lte:       sub.Lte,
			StaticIps: sub.StaticIps,
		},
		Associations: sub.GetAssocs(),
	}

	var ents []configurator.NetworkEntity
	ents = append(ents, sub.ActivePoliciesByApn.ToEntities(subEnt.Key)...)
	ents = append(ents, subEnt)

	_, err := configurator.CreateEntities(networkID, ents, serdes.Entity)
	if err != nil {
		return err
	}

	return nil
}

func updateSubscriber(networkID string, sub *subscribermodels.MutableSubscriber) error {
	var writes []configurator.EntityWriteOperation

	existingSub, err := configurator.LoadEntity(
		networkID, lte.SubscriberEntityType, string(sub.ID),
		configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err != nil {
		return err
	}

	// For simplicity, delete all of subscriber's existing
	// apn_policy_profile, then add new
	policyMapTKs := existingSub.Associations.Filter(lte.APNPolicyProfileEntityType)
	for _, tk := range policyMapTKs {
		writes = append(writes, configurator.EntityUpdateCriteria{Type: tk.Type, Key: tk.Key, DeleteEntity: true})
	}
	for _, e := range sub.ActivePoliciesByApn.ToEntities(string(sub.ID)) {
		writes = append(writes, e)
	}

	subUpdate := configurator.EntityUpdateCriteria{
		Key:     string(sub.ID),
		Type:    lte.SubscriberEntityType,
		NewName: swag.String(sub.Name),
		NewConfig: &subscribermodels.SubscriberConfig{
			Lte:       sub.Lte,
			StaticIps: sub.StaticIps,
		},
		AssociationsToSet: sub.GetAssocs(),
	}
	writes = append(writes, subUpdate)

	err = configurator.WriteEntities(networkID, writes, serdes.Entity)
	if err != nil {
		return err
	}

	return nil
}

func deleteSubscriber(networkID, key string) error {
	ent, err := configurator.LoadEntity(
		networkID, lte.SubscriberEntityType, key,
		configurator.EntityLoadCriteria{LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err != nil {
		return err
	}
	// Configurator doesn't currently support loading a specified subgraph,
	// so we have to load the subscriber and its apn_policy_profile ents in
	// separate calls.
	var policyProfileEnts configurator.NetworkEntities
	if ppAssocs := ent.Associations.Filter(lte.APNPolicyProfileEntityType); len(ppAssocs) != 0 {
		policyProfileEnts, _, err = configurator.LoadEntities(
			ent.NetworkID, nil, nil, nil,
			ppAssocs,
			apnPolicyProfileLoadCriteria,
			serdes.Entity,
		)
		if err != nil {
			return err
		}
	}

	sub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent, policyProfileEnts)
	if err != nil {
		return err
	}

	var deletes []storage.TypeAndKey
	deletes = append(deletes, sub.ToTK())
	deletes = append(deletes, sub.ActivePoliciesByApn.ToTKs(string(sub.ID))...)

	err = configurator.DeleteEntities(networkID, deletes)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	return nil
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return obsidian.HttpError(err, http.StatusInternalServerError)
}
