/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package handlers

import (
	"net/http"
	"regexp"

	"magma/lte/cloud/go/lte"
	ltehandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
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

var subscriberStateTypes = []string{
	lte.ICMPStateType,
	lte.S1APStateType,
	lte.MMEStateType,
	lte.SPGWStateType,
	lte.MobilitydStateType,
}

// mobilityd states are keyed as <ISMI>.<APN>. This captures just the imsi
// portion in a named match group
var mobilitydStateKeyRe = regexp.MustCompile(`^(?P<imsi>IMSI\d+)\..+$`)

const mobilitydStateExpectedMatchCount = 2

var subscriberLoadCriteria = configurator.EntityLoadCriteria{LoadMetadata: true, LoadConfig: true, LoadAssocsFromThis: true}

func listSubscribers(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	ents, err := configurator.LoadAllEntitiesInNetwork(networkID, lte.SubscriberEntityType, subscriberLoadCriteria)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	subStates, err := state.SearchStates(networkID, subscriberStateTypes, nil, nil)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	// Each entry in this map contains all the states that the SID cares about.
	// The DeviceID fields of the state IDs in the nested maps do not have to
	// match the SID, as in the case of mobilityd state for example.
	statesBySid := map[string]state_types.StatesByID{}
	for stateID, st := range subStates {
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

	ret := make(map[string]*subscribermodels.Subscriber, len(ents))
	for _, ent := range ents {
		ret[ent.Key] = (&subscribermodels.Subscriber{}).FromBackendModels(ent, statesBySid[ent.Key])
	}
	return c.JSON(http.StatusOK, ret)
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

	_, err := configurator.CreateEntity(networkID, configurator.NetworkEntity{
		Type:         lte.SubscriberEntityType,
		Key:          string(payload.ID),
		Name:         payload.Name,
		Config:       payload.Lte,
		Associations: payload.ActiveApns.ToAssocs(),
	})
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

	ent, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, subscriberLoadCriteria)
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	keyPrefix := swag.String(ent.Key)
	states, err := state.SearchStates(networkID, subscriberStateTypes, nil, keyPrefix)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}

	ret := (&subscribermodels.Subscriber{}).FromBackendModels(ent, states)
	return c.JSON(http.StatusOK, ret)
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

	_, err := configurator.LoadEntity(networkID, lte.SubscriberEntityType, subscriberID, configurator.EntityLoadCriteria{LoadAssocsFromThis: true})
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "failed to load existing subscriber"), http.StatusInternalServerError)
	}

	if nerr := validateSubscriberProfile(networkID, payload.Lte); nerr != nil {
		return nerr
	}

	updateCriteria := configurator.EntityUpdateCriteria{
		Key:               subscriberID,
		Type:              lte.SubscriberEntityType,
		NewName:           swag.String(payload.Name),
		NewConfig:         payload.Lte,
		AssociationsToSet: payload.ActiveApns.ToAssocs(),
	}
	_, err = configurator.UpdateEntities(networkID, []configurator.EntityUpdateCriteria{updateCriteria})
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteSubscriber(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	err := configurator.DeleteEntity(networkID, lte.SubscriberEntityType, subscriberID)
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
	switch {
	case err == merrors.ErrNotFound:
		return echo.ErrNotFound
	case err != nil:
		return obsidian.HttpError(errors.Wrap(err, "could not load subscriber"), http.StatusInternalServerError)
	}

	desiredCfg := currentCfg.(*subscribermodels.LteSubscription)
	desiredCfg.SubProfile = *payload
	if nerr := validateSubscriberProfile(networkID, desiredCfg); nerr != nil {
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
		switch {
		case err == merrors.ErrNotFound:
			return echo.ErrNotFound
		case err != nil:
			return obsidian.HttpError(errors.Wrap(err, "failed to load existing subscriber"), http.StatusInternalServerError)
		}

		newConfig := cfg.(*subscribermodels.LteSubscription)
		newConfig.State = desiredState
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
		netConf, err := configurator.LoadNetworkConfig(networkID, lte.CellularNetworkType)
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
