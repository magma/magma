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
	"context"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"github.com/labstack/echo"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/lte/cloud/go/lte"
	"magma/lte/cloud/go/serdes"
	ltehandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	ltemodels "magma/lte/cloud/go/services/lte/obsidian/models"
	policydbmodels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	subscribermodels "magma/lte/cloud/go/services/subscriberdb/obsidian/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serde"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

const (
	Subscribers               = "subscribers"
	SubscriberState           = "subscriber_state"
	ListSubscribersPath       = ltehandlers.ManageNetworkPath + obsidian.UrlSep + Subscribers
	ManageSubscriberPath      = ListSubscribersPath + obsidian.UrlSep + ":subscriber_id"
	ListSubscriberStatePath   = ltehandlers.ManageNetworkPath + obsidian.UrlSep + SubscriberState
	ManageSubscriberStatePath = ListSubscriberStatePath + obsidian.UrlSep + ":subscriber_id"
	ActivateSubscriberPath    = ManageSubscriberPath + obsidian.UrlSep + "activate"
	DeactivateSubscriberPath  = ManageSubscriberPath + obsidian.UrlSep + "deactivate"
	SubscriberProfilePath     = ManageSubscriberPath + obsidian.UrlSep + "lte" + obsidian.UrlSep + "sub_profile"

	listMSISDNsPath   = ltehandlers.ManageNetworkPath + obsidian.UrlSep + "msisdns"
	manageMSISDNsPath = listMSISDNsPath + obsidian.UrlSep + ":msisdn"

	ParamMSISDN    = "msisdn"
	ParamIP        = "ip"
	ParamPageSize  = "page_size"
	ParamPageToken = "page_token"
	ParamVerbose   = "verbose"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		{Path: ListSubscribersPath, Methods: obsidian.GET, HandlerFunc: listSubscribersHandler},
		{Path: ListSubscribersPath, Methods: obsidian.POST, HandlerFunc: createSubscribersHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.GET, HandlerFunc: getSubscriberHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.PUT, HandlerFunc: updateSubscriberHandler},
		{Path: ManageSubscriberPath, Methods: obsidian.DELETE, HandlerFunc: deleteSubscriberHandler},

		{Path: ListSubscriberStatePath, Methods: obsidian.GET, HandlerFunc: listSubscriberStateHandler},
		{Path: ManageSubscriberStatePath, Methods: obsidian.GET, HandlerFunc: getSubscriberStateHandler},

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

	apnPolicyProfileLoadCriteria = configurator.EntityLoadCriteria{LoadAssocsFromThis: true, LoadAssocsToThis: true}
)

// The following slices comprise the various state types that make up
// subscriber state. Here the state types are separated to allow for more
// efficient state lookup for paginated subscriber requests.

// subscriberStateTypesKeyedByIMSI is a slice of subscriber state types whose
// deviceID is an IMSI.
var subscriberStateTypesKeyedByIMSI = []string{
	lte.ICMPStateType,
	lte.MMEStateType,
	lte.S1APStateType,
	lte.SPGWStateType,
	lte.SubscriberStateType,
	orc8r.DirectoryRecordType,
}

// subscriberStateTypesKeyedByCompositeKey is a slice of subscriber state
// types whose deviceID is a composite key with format <IMSI>.<APN>.
var subscriberStateTypesKeyedByCompositeKey = []string{
	lte.MobilitydStateType,
}

// allSubscriberStateTypes is a composite of all subscriber state types.
var allSubscriberStateTypes = append(subscriberStateTypesKeyedByIMSI, subscriberStateTypesKeyedByCompositeKey...)

type subscriberFilter func(sub *subscribermodels.Subscriber) bool

func acceptAll(*subscribermodels.Subscriber) bool { return true }

// mapSubscribersForVerbosity filters a subscribers list for the specified verbosity.
// If verbose mode is used, the list is returned as-is (a map of MSISDN strings
// to Subscribers).
// If non-verbose mode is used, an array of MSISDN strings is returned instead.
func mapSubscribersForVerbosity(subs map[string]*subscribermodels.Subscriber, verbose bool) interface{} {
	if verbose {
		return subs
	}

	subsIds := make([]string, 0, len(subs))
	for k := range subs {
		subsIds = append(subsIds, k)
	}

	return subsIds
}

// listSubscribersHandler handles the subscriber endpoint.
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
//
// The returned subscribers can be paginated using the following parameters
//  - page_size
//  - page_token
//
// The page size parameter specifies the maximum number of subscribers to
// return.
//
// The page token parameter is an opaque token used to fetch the next page of
// subscribers. Each API response will contain a page token that can be used
// to fetch the next page.
func listSubscribersHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	var pageSize uint64 = 0
	var err error
	if pageSizeParam := c.QueryParam(ParamPageSize); pageSizeParam != "" {
		pageSize, err = strconv.ParseUint(pageSizeParam, 10, 32)
		if err != nil {
			err := fmt.Errorf("invalid page size parameter: %s", err)
			return echo.NewHTTPError(http.StatusBadRequest, err)
		}
	}
	pageToken := c.QueryParam(ParamPageToken)
	reqCtx := c.Request().Context()

	verbose, err := strconv.ParseBool(c.QueryParam(ParamVerbose))
	if err != nil {
		verbose = true
	}

	// First check for query params to filter by
	if msisdn := c.QueryParam(ParamMSISDN); msisdn != "" {
		queryIMSI, err := subscriberdb.GetIMSIForMSISDN(reqCtx, networkID, msisdn)
		if err != nil {
			return makeErr(err)
		}
		subs, err := loadSubscribers(reqCtx, networkID, acceptAll, queryIMSI)
		if err != nil {
			return makeErr(err)
		}
		return c.JSON(http.StatusOK, mapSubscribersForVerbosity(subs, verbose))
	}
	if ip := c.QueryParam(ParamIP); ip != "" {
		queryIMSIs, err := subscriberdb.GetIMSIsForIP(reqCtx, networkID, ip)
		if err != nil {
			return makeErr(err)
		}
		filter := func(sub *subscribermodels.Subscriber) bool { return sub.IsAssignedIP(ip) }
		subs, err := loadSubscribers(reqCtx, networkID, filter, queryIMSIs...)
		if err != nil {
			return makeErr(err)
		}
		return c.JSON(http.StatusOK, mapSubscribersForVerbosity(subs, verbose))
	}

	// List subscribers for a given page. If no page is specified, the max
	// size will be returned.
	subs, nextPageToken, err := loadSubscriberPage(reqCtx, networkID, uint32(pageSize), pageToken)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	// get total number of subscribers
	count, err := configurator.CountEntitiesOfType(reqCtx, networkID, lte.SubscriberEntityType)
	if err != nil {
		return c.JSON(http.StatusOK, nil)
	}
	var paginatedSubs interface{}
	if verbose {
		paginatedSubs = subscribermodels.PaginatedSubscribers{
			TotalCount:    int64(count),
			NextPageToken: subscribermodels.PageToken(nextPageToken),
			Subscribers:   mapSubscribersForVerbosity(subs, verbose).(map[string]*subscribermodels.Subscriber),
		}
	} else {
		paginatedSubs = subscribermodels.PaginatedSubscriberIds{
			TotalCount:    int64(count),
			NextPageToken: subscribermodels.PageToken(nextPageToken),
			Subscribers:   mapSubscribersForVerbosity(subs, verbose).([]string),
		}
	}
	return c.JSON(http.StatusOK, paginatedSubs)
}

func createSubscribersHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	payload := subscribermodels.MutableSubscribers{}
	if err := c.Bind(&payload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	reqCtx := c.Request().Context()
	if err := payload.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if nerr := validateSubscriberProfiles(reqCtx, networkID, getSubProfiles(payload)...); nerr != nil {
		return nerr
	}

	nerr = createSubscribers(reqCtx, networkID, payload...)
	if nerr != nil {
		return nerr
	}

	return c.NoContent(http.StatusCreated)
}

func getSubscriberHandler(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}
	subs, err := loadSubscriber(c.Request().Context(), networkID, subscriberID)
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	reqCtx := c.Request().Context()
	if err := payload.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if string(payload.ID) != subscriberID {
		err := fmt.Errorf("subscriber ID from parameters (%s) and payload (%s) must match", subscriberID, payload.ID)
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	if nerr := validateSubscriberProfiles(reqCtx, networkID, string(payload.Lte.SubProfile)); nerr != nil {
		return nerr
	}

	err := updateSubscriber(reqCtx, networkID, payload)
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
	err := deleteSubscriber(c.Request().Context(), networkID, subscriberID)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNoContent)
	}
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}
	return c.NoContent(http.StatusNoContent)
}

func listSubscriberStateHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	statesBySID, err := loadAllStatesForIMSIs(c.Request().Context(), networkID, []string{})
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	modelsBySID := map[string]*subscribermodels.SubscriberState{}
	for sid, states := range statesBySID {
		modelsBySID[sid] = makeSubscriberState(sid, states)
	}

	return c.JSON(http.StatusOK, modelsBySID)
}

func getSubscriberStateHandler(c echo.Context) error {
	networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
	if nerr != nil {
		return nerr
	}

	states, err := getStatesForIMSIs(c.Request().Context(), networkID, allSubscriberStateTypes, subscriberID, serdes.State)
	if err != nil {
		return makeErr(err)
	}

	subState := makeSubscriberState(subscriberID, states)
	return c.JSON(http.StatusOK, subState)
}

func listMSISDNsHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}

	msisdns, err := subscriberdb.ListMSISDNs(c.Request().Context(), networkID)
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}
	if err := payload.ValidateModel(context.Background()); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	err := subscriberdb.SetIMSIForMSISDN(c.Request().Context(), networkID, string(payload.Msisdn), string(payload.ID))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return c.NoContent(http.StatusCreated)
}

func getMSISDNHandler(c echo.Context) error {
	networkID, msisdn, nerr := getNetworkAndMSISDN(c)
	if nerr != nil {
		return nerr
	}
	imsi, err := subscriberdb.GetIMSIForMSISDN(c.Request().Context(), networkID, msisdn)
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

	err := subscriberdb.DeleteMSISDN(c.Request().Context(), networkID, msisdn)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
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
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	reqCtx := c.Request().Context()
	if err := payload.ValidateModel(reqCtx); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	currentCfg, err := configurator.LoadEntityConfig(reqCtx, networkID, lte.SubscriberEntityType, subscriberID, serdes.Entity)
	if err != nil {
		return makeErr(err)
	}

	desiredCfg := currentCfg.(*subscribermodels.SubscriberConfig)
	desiredCfg.Lte.SubProfile = *payload
	if nerr := validateSubscriberProfiles(reqCtx, networkID, string(desiredCfg.Lte.SubProfile)); nerr != nil {
		return nerr
	}

	_, err = configurator.UpdateEntity(
		reqCtx,
		networkID,
		configurator.EntityUpdateCriteria{Type: lte.SubscriberEntityType, Key: subscriberID, NewConfig: desiredCfg},
		serdes.Entity,
	)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, errors.Wrap(err, "failed to update profile"))
	}
	return c.NoContent(http.StatusNoContent)
}

func makeSubscriberStateHandler(desiredState string) echo.HandlerFunc {
	return func(c echo.Context) error {
		networkID, subscriberID, nerr := getNetworkAndSubIDs(c)
		if nerr != nil {
			return nerr
		}
		reqCtx := c.Request().Context()

		cfg, err := configurator.LoadEntityConfig(reqCtx, networkID, lte.SubscriberEntityType, subscriberID, serdes.Entity)
		if err != nil {
			return makeErr(err)
		}

		newConfig := cfg.(*subscribermodels.SubscriberConfig)
		newConfig.Lte.State = desiredState
		err = configurator.CreateOrUpdateEntityConfig(reqCtx, networkID, lte.SubscriberEntityType, subscriberID, newConfig, serdes.Entity)
		if err != nil {
			return echo.NewHTTPError(http.StatusInternalServerError, err)
		}
		return c.NoContent(http.StatusOK)
	}
}

func getStatesForIMSIs(ctx context.Context, networkID string, typeFilter []string, keyPrefix string, serdes serde.Registry) (state_types.StatesByID, error) {
	states, err := state.SearchStates(ctx, networkID, typeFilter, nil, &keyPrefix, serdes)
	if err != nil {
		return nil, err
	}
	// Returned states contain matches by prefix, so filter out non-exact matches
	for stateID := range states {
		imsi := stateID.DeviceID
		if stateID.Type == lte.MobilitydStateType {
			matches := mobilitydStateKeyRe.FindStringSubmatch(stateID.DeviceID)
			if len(matches) != mobilitydStateExpectedMatchCount {
				glog.Infof("state device ID '%s' with type '%s' did not match IMSI-prefixed regex", stateID.DeviceID, stateID.Type)
				continue
			}
			imsi = matches[1]
		}

		if imsi != keyPrefix {
			delete(states, stateID)
		}
	}

	return states, nil
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

func getSubProfiles(subs subscribermodels.MutableSubscribers) []string {
	profiles := map[string]struct{}{}
	for _, sub := range subs {
		profiles[string(sub.Lte.SubProfile)] = struct{}{}
	}
	return funk.Keys(profiles).([]string)
}

func validateSubscriberProfiles(ctx context.Context, networkID string, profiles ...string) *echo.HTTPError {
	nonDefaultProfiles := funk.FilterString(profiles, func(s string) bool { return s != "default" })

	if len(nonDefaultProfiles) == 0 {
		return nil
	}

	networkConfig, err := configurator.LoadNetworkConfig(ctx, networkID, lte.CellularNetworkConfigType, serdes.Network)
	if err == merrors.ErrNotFound {
		return echo.NewHTTPError(http.StatusBadRequest, errors.New("no cellular config found for network"))
	}
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	networkProfiles := networkConfig.(*ltemodels.NetworkCellularConfigs).Epc.SubProfiles
	errs := &multierror.Error{}
	for _, p := range nonDefaultProfiles {
		if _, ok := networkProfiles[p]; !ok {
			errs = multierror.Append(errs, errors.Errorf("subscriber profile '%s' does not exist for the network", p))
		}
	}
	err = errs.ErrorOrNil()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err)
	}

	return nil
}

func loadSubscriber(ctx context.Context, networkID, key string) (*subscribermodels.Subscriber, error) {
	loadCriteria := getSubscriberLoadCriteria(0, "")
	ent, err := configurator.LoadEntity(ctx, networkID, lte.SubscriberEntityType, key, loadCriteria, serdes.Entity)
	if err != nil {
		return nil, err
	}

	// Configurator doesn't currently support loading a specified subgraph,
	// so we have to load the subscriber and its apn_policy_profile ents in
	// separate calls.
	var policyProfileEnts configurator.NetworkEntities
	if ppAssocs := ent.Associations.Filter(lte.APNPolicyProfileEntityType); len(ppAssocs) != 0 {
		policyProfileEnts, _, err = configurator.LoadEntities(
			ctx,
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

	states, err := getStatesForIMSIs(ctx, networkID, allSubscriberStateTypes, key, serdes.State)
	if err != nil {
		return nil, err
	}

	sub := mutableSub.ToSubscriber()
	sub.FillAugmentedFields(states)
	return sub, nil
}

func loadSubscribers(ctx context.Context, networkID string, includeSub subscriberFilter, keys ...string) (map[string]*subscribermodels.Subscriber, error) {
	subs := map[string]*subscribermodels.Subscriber{}
	for _, key := range keys {
		sub, err := loadSubscriber(ctx, networkID, key)
		if err != nil {
			return nil, errors.Wrapf(err, "error loading subscriber %s", key)
		}
		if includeSub(sub) {
			subs[string(sub.ID)] = sub
		}
	}
	return subs, nil
}

func loadSubscriberPage(ctx context.Context, networkID string, pageSize uint32, pageToken string) (map[string]*subscribermodels.Subscriber, string, error) {
	mutableSubs, nextPageToken, err := loadMutableSubscriberPage(ctx, networkID, pageSize, pageToken)
	if err != nil {
		return nil, "", err
	}
	imsis := make([]string, 0, len(mutableSubs))
	for imsi := range mutableSubs {
		imsis = append(imsis, imsi)
	}
	states, err := loadAllStatesForIMSIs(ctx, networkID, imsis)
	if err != nil {
		return nil, "", err
	}

	subs := map[string]*subscribermodels.Subscriber{}
	for _, mutableSub := range mutableSubs {
		sub := mutableSub.ToSubscriber()
		sub.FillAugmentedFields(states[string(sub.ID)])
		subs[string(sub.ID)] = sub
	}

	return subs, nextPageToken, nil
}

func loadMutableSubscriberPage(ctx context.Context, networkID string, pageSize uint32, pageToken string) (map[string]*subscribermodels.MutableSubscriber, string, error) {
	loadCriteria := getSubscriberLoadCriteria(pageSize, pageToken)
	ents, nextPageToken, err := configurator.LoadAllEntitiesOfType(ctx, networkID, lte.SubscriberEntityType, loadCriteria, serdes.Entity)
	if err != nil {
		return nil, "", err
	}
	profileEnts, _, err := configurator.LoadAllEntitiesOfType(
		ctx,
		networkID, lte.APNPolicyProfileEntityType,
		apnPolicyProfileLoadCriteria,
		serdes.Entity,
	)
	if err != nil {
		return nil, "", err
	}
	profileEntsBySub := profileEnts.MakeByParentTK()

	subs := map[string]*subscribermodels.MutableSubscriber{}
	for _, ent := range ents {
		sub, err := (&subscribermodels.MutableSubscriber{}).FromEnt(ent, profileEntsBySub[ent.GetTK()])
		if err != nil {
			return nil, "", err
		}
		subs[ent.Key] = sub
	}
	return subs, nextPageToken, nil
}

func createSubscribers(ctx context.Context, networkID string, subs ...*subscribermodels.MutableSubscriber) *echo.HTTPError {
	var ents configurator.NetworkEntities
	var ids []string
	uniqueIDs := map[string]int{}

	for _, s := range subs {
		ents = append(ents, getCreateSubscriberEnts(s)...)

		id := string(s.ID)
		ids = append(ids, id)
		uniqueIDs[id] = uniqueIDs[id] + 1
	}

	if len(uniqueIDs) != len(ids) {
		duplicates := funk.FilterString(ids, func(s string) bool { return uniqueIDs[s] > 1 })
		return echo.NewHTTPError(http.StatusBadRequest, errors.Errorf("found multiple subscriber models for IDs: %+v", duplicates))
	}

	// TODO(hcgatewood) iterate over this to remove "too many placeholders" error
	tks := storage.MakeTKs(lte.SubscriberEntityType, ids)
	found, _, err := configurator.LoadSerializedEntities(ctx, networkID, nil, nil, nil, tks, configurator.EntityLoadCriteria{})
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}
	if len(found) != 0 {
		return echo.NewHTTPError(http.StatusBadRequest, errors.Errorf("found %v existing subscribers which would have been overwritten: %+v", len(found), found.TKs()))
	}

	_, err = configurator.CreateEntities(ctx, networkID, ents, serdes.Entity)
	if err != nil {
		return obsidian.MakeHTTPError(err, http.StatusInternalServerError)
	}

	return nil
}

func getCreateSubscriberEnts(sub *subscribermodels.MutableSubscriber) configurator.NetworkEntities {
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
			Lte:                   sub.Lte,
			StaticIps:             sub.StaticIps,
			ForbiddenNetworkTypes: sub.ForbiddenNetworkTypes,
		},
		Associations: sub.GetAssocs(),
	}

	var ents []configurator.NetworkEntity
	ents = append(ents, sub.ActivePoliciesByApn.ToEntities(subEnt.Key)...)
	ents = append(ents, subEnt)

	return ents
}

func updateSubscriber(ctx context.Context, networkID string, sub *subscribermodels.MutableSubscriber) error {
	var writes []configurator.EntityWriteOperation

	existingSub, err := configurator.LoadEntity(
		ctx,
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
			Lte:                   sub.Lte,
			StaticIps:             sub.StaticIps,
			ForbiddenNetworkTypes: sub.ForbiddenNetworkTypes,
		},
		AssociationsToSet: sub.GetAssocs(),
	}
	writes = append(writes, subUpdate)

	err = configurator.WriteEntities(ctx, networkID, writes, serdes.Entity)
	if err != nil {
		return err
	}

	return nil
}

func deleteSubscriber(ctx context.Context, networkID, key string) error {
	ent, err := configurator.LoadEntity(
		ctx,
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
			ctx,
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

	var deletes storage.TKs
	deletes = append(deletes, sub.ToTK())
	deletes = append(deletes, sub.ActivePoliciesByApn.ToTKs(string(sub.ID))...)

	err = configurator.DeleteEntities(ctx, networkID, deletes)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err)
	}

	return nil
}

// loadAllStatesForIMSIs loads all states whose IMSI prefix is contained in the
// IMSI array passed in as argument. If passed IMSIs is nil,
// loads states for all IMSIs in the network.
func loadAllStatesForIMSIs(ctx context.Context, networkID string, imsis []string) (map[string]state_types.StatesByID, error) {
	requestedIMSIs := map[string]struct{}{}
	for _, v := range imsis {
		requestedIMSIs[v] = struct{}{}
	}

	shouldLoadState := func(imsi string) bool {
		if len(requestedIMSIs) == 0 {
			// load all states regardless of their IMSIs if requested IMSIs is nil,
			return true
		}
		_, ok := requestedIMSIs[imsi]
		return ok
	}

	imsiKeyStates, err := state.SearchStates(ctx, networkID, subscriberStateTypesKeyedByIMSI, imsis, nil, serdes.State)
	if err != nil {
		return nil, err
	}
	imsiCompositeKeyStates, err := state.SearchStates(ctx, networkID, subscriberStateTypesKeyedByCompositeKey, nil, nil, serdes.State)
	if err != nil {
		return nil, err
	}
	states := mergeStates(imsiKeyStates, imsiCompositeKeyStates)
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
			if !shouldLoadState(matches[1]) {
				continue
			}
			sidKey = matches[1]
		}

		if _, exists := statesBySid[sidKey]; !exists {
			statesBySid[sidKey] = state_types.StatesByID{}
		}
		statesBySid[sidKey][stateID] = st
	}

	return statesBySid, nil
}

func makeSubscriberState(subscriberID string, states state_types.StatesByID) *subscribermodels.SubscriberState {
	// Create anonymous subscriber (may or may not have a backing configurator
	// entity), then extract its formatted state
	sub := &subscribermodels.Subscriber{ID: policydbmodels.SubscriberID(subscriberID)}
	sub.FillAugmentedFields(states)
	if sub.State == nil {
		return &subscribermodels.SubscriberState{}
	}
	return sub.State
}

func mergeStates(s1 state_types.StatesByID, s2 state_types.StatesByID) state_types.StatesByID {
	for id, state := range s2 {
		s1[id] = state
	}
	return s1
}

func getSubscriberLoadCriteria(pageSize uint32, pageToken string) configurator.EntityLoadCriteria {
	loadCriteria := configurator.EntityLoadCriteria{
		LoadMetadata:       true,
		LoadConfig:         true,
		LoadAssocsFromThis: true,
		PageSize:           pageSize,
		PageToken:          pageToken,
	}
	return loadCriteria
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return echo.NewHTTPError(http.StatusInternalServerError, err)
}
