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
	"encoding/json"
	"fmt"
	"net/http"

	"magma/cwf/cloud/go/cwf"
	"magma/cwf/cloud/go/serdes"
	cwfModels "magma/cwf/cloud/go/services/cwf/obsidian/models"
	fegModels "magma/feg/cloud/go/services/feg/obsidian/models"
	lteHandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	directorydTypes "magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/labstack/echo"
	"github.com/pkg/errors"
)

const (
	CwfNetworks                    = "cwf"
	ListNetworksPath               = obsidian.V1Root + CwfNetworks
	ManageNetworkPath              = ListNetworksPath + "/:network_id"
	ManageNetworkNamePath          = ManageNetworkPath + obsidian.UrlSep + "name"
	ManageNetworkDescriptionPath   = ManageNetworkPath + obsidian.UrlSep + "description"
	ManageNetworkFeaturesPath      = ManageNetworkPath + obsidian.UrlSep + "features"
	ManageNetworkDNSPath           = ManageNetworkPath + obsidian.UrlSep + "dns"
	ManageNetworkCarrierWifiPath   = ManageNetworkPath + obsidian.UrlSep + "carrier_wifi"
	ManageNetworkFederationPath    = ManageNetworkPath + obsidian.UrlSep + "federation"
	ListNetworkHAPairsPath         = ManageNetworkPath + obsidian.UrlSep + "ha_pairs"
	ManageNetworkHAPairsPath       = ListNetworkHAPairsPath + obsidian.UrlSep + ":ha_pair_id"
	ManageNetworkHAPairsStatusPath = ManageNetworkHAPairsPath + obsidian.UrlSep + "status"
	ManageNetworkSubscriberPath    = ManageNetworkPath + obsidian.UrlSep + "subscriber_config"
	ManageNetworkBaseNamesPath     = ManageNetworkSubscriberPath + obsidian.UrlSep + "base_names"
	ManageNetworkRuleNamesPath     = ManageNetworkSubscriberPath + obsidian.UrlSep + "rule_names"
	ManageNetworkBaseNamePath      = ManageNetworkBaseNamesPath + obsidian.UrlSep + ":base_name"
	ManageNetworkRuleNamePath      = ManageNetworkRuleNamesPath + obsidian.UrlSep + ":rule_id"
	ManageNetworkLiUesPath         = ManageNetworkPath + obsidian.UrlSep + ":li_ues"

	Gateways                      = "gateways"
	ListGatewaysPath              = ManageNetworkPath + obsidian.UrlSep + Gateways
	ManageGatewayPath             = ListGatewaysPath + obsidian.UrlSep + ":gateway_id"
	ManageGatewayNamePath         = ManageGatewayPath + obsidian.UrlSep + "name"
	ManageGatewayDescriptionPath  = ManageGatewayPath + obsidian.UrlSep + "description"
	ManageGatewayConfigPath       = ManageGatewayPath + obsidian.UrlSep + "magmad"
	ManageGatewayDevicePath       = ManageGatewayPath + obsidian.UrlSep + "device"
	ManageGatewayStatePath        = ManageGatewayPath + obsidian.UrlSep + "status"
	ManageGatewayTierPath         = ManageGatewayPath + obsidian.UrlSep + "tier"
	ManageGatewayCarrierWifiPath  = ManageGatewayPath + obsidian.UrlSep + "carrier_wifi"
	ManageGatewayHealthStatusPath = ManageGatewayPath + obsidian.UrlSep + "health_status"

	Subscribers                   = "subscribers"
	BaseSubscriberPath            = ManageNetworkPath + obsidian.UrlSep + Subscribers + obsidian.UrlSep + ":subscriber_id"
	SubscriberDirectoryRecordPath = BaseSubscriberPath + obsidian.UrlSep + "directory_record"
)

func GetHandlers() []obsidian.Handler {
	ret := []obsidian.Handler{
		handlers.GetListGatewaysHandler(ListGatewaysPath, &cwfModels.MutableCwfGateway{}, makeCwfGateways, serdes.Entity, serdes.Device),
		{Path: ListGatewaysPath, Methods: obsidian.POST, HandlerFunc: createGateway},
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageGatewayPath, Methods: obsidian.PUT, HandlerFunc: updateGateway},
		{Path: ManageGatewayPath, Methods: obsidian.DELETE, HandlerFunc: deleteGateway},

		{Path: ManageGatewayStatePath, Methods: obsidian.GET, HandlerFunc: handlers.GetStateHandler},
		{Path: ManageNetworkHAPairsStatusPath, Methods: obsidian.GET, HandlerFunc: getHAPairStatusHandler},
		{Path: ManageGatewayHealthStatusPath, Methods: obsidian.GET, HandlerFunc: getHealthStatusHandler},
		{Path: SubscriberDirectoryRecordPath, Methods: obsidian.GET, HandlerFunc: getSubscriberDirectoryHandler},

		{Path: ListNetworkHAPairsPath, Methods: obsidian.GET, HandlerFunc: listHAPairsHandler},
		{Path: ListNetworkHAPairsPath, Methods: obsidian.POST, HandlerFunc: createHAPairHandler},
		{Path: ManageNetworkHAPairsPath, Methods: obsidian.GET, HandlerFunc: getHAPairHandler},
		{Path: ManageNetworkHAPairsPath, Methods: obsidian.PUT, HandlerFunc: updateHAPairHandler},
		{Path: ManageNetworkHAPairsPath, Methods: obsidian.DELETE, HandlerFunc: deleteHAPairHandler},

		{Path: ManageNetworkBaseNamePath, Methods: obsidian.POST, HandlerFunc: lteHandlers.AddNetworkWideSubscriberBaseName},
		{Path: ManageNetworkRuleNamePath, Methods: obsidian.POST, HandlerFunc: lteHandlers.AddNetworkWideSubscriberRuleName},
		{Path: ManageNetworkBaseNamePath, Methods: obsidian.DELETE, HandlerFunc: lteHandlers.RemoveNetworkWideSubscriberBaseName},
		{Path: ManageNetworkRuleNamePath, Methods: obsidian.DELETE, HandlerFunc: lteHandlers.RemoveNetworkWideSubscriberRuleName},
	}

	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(ListNetworksPath, ManageNetworkPath, cwf.CwfNetworkType, &cwfModels.CwfNetwork{}, serdes.Network)...)

	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "", serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "", serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rModels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig, serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSPath, &orc8rModels.NetworkDNSConfig{}, orc8r.DnsdNetworkType, serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCarrierWifiPath, &cwfModels.NetworkCarrierWifiConfigs{}, cwf.CwfNetworkType, serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFederationPath, &fegModels.FederatedNetworkConfigs{}, cwf.CwfNetworkType, serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkSubscriberPath, &policyModels.NetworkSubscriberConfig{}, "", serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkRuleNamesPath, new(policyModels.RuleNames), "", serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkBaseNamesPath, new(policyModels.BaseNames), "", serdes.Network)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkLiUesPath, new(cwfModels.LiUes), "", serdes.Network)...)

	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayNamePath, new(models.GatewayName), serdes.Entity)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayDescriptionPath, new(models.GatewayDescription), serdes.Entity)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayConfigPath, &orc8rModels.MagmadGatewayConfigs{}, serdes.Entity)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayTierPath, new(orc8rModels.TierID), serdes.Entity)...)
	ret = append(ret, handlers.GetGatewayDeviceHandlers(ManageGatewayDevicePath, serdes.Device)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCarrierWifiPath, &cwfModels.GatewayCwfConfigs{}, serdes.Entity)...)

	return ret
}

func createGateway(c echo.Context) error {
	if nerr := handlers.CreateGateway(c, &cwfModels.MutableCwfGateway{}, serdes.Entity, serdes.Device); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusCreated)
}

func getGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}

	magmadModel, nerr := handlers.LoadMagmadGateway(nid, gid)
	if nerr != nil {
		return nerr
	}

	ent, err := configurator.LoadEntity(
		nid, cwf.CwfGatewayType, gid,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: false},
		serdes.Entity,
	)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "failed to load cwf gateway"), http.StatusInternalServerError)
	}

	ret := &cwfModels.CwfGateway{
		ID:          magmadModel.ID,
		Name:        magmadModel.Name,
		Description: magmadModel.Description,
		Device:      magmadModel.Device,
		Status:      magmadModel.Status,
		Tier:        magmadModel.Tier,
		Magmad:      magmadModel.Magmad,
	}
	if ent.Config != nil {
		ret.CarrierWifi = ent.Config.(*cwfModels.GatewayCwfConfigs)
	}

	return c.JSON(http.StatusOK, ret)
}

func updateGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	if nerr = handlers.UpdateGateway(c, nid, gid, &cwfModels.MutableCwfGateway{}, serdes.Entity, serdes.Device); nerr != nil {
		return nerr
	}
	return c.NoContent(http.StatusNoContent)
}

func deleteGateway(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	err := handlers.DeleteMagmadGateway(nid, gid, storage.TKs{{Type: cwf.CwfGatewayType, Key: gid}})
	if err != nil {
		return makeErr(err)
	}
	return c.NoContent(http.StatusNoContent)
}

type cwfAndMagmadGateway struct {
	magmadGateway, cwfGateway configurator.NetworkEntity
}

func makeCwfGateways(
	entsByTK configurator.NetworkEntitiesByTK,
	devicesByID map[string]interface{},
	statusesByID map[string]*orc8rModels.GatewayStatus,
) map[string]handlers.GatewayModel {
	gatewayEntsByKey := map[string]*cwfAndMagmadGateway{}
	for tk, ent := range entsByTK.MultiFilter(orc8r.MagmadGatewayType, cwf.CwfGatewayType) {
		existing, found := gatewayEntsByKey[tk.Key]
		if !found {
			existing = &cwfAndMagmadGateway{}
			gatewayEntsByKey[tk.Key] = existing
		}

		switch ent.Type {
		case orc8r.MagmadGatewayType:
			existing.magmadGateway = ent
		case cwf.CwfGatewayType:
			existing.cwfGateway = ent
		}
	}

	ret := make(map[string]handlers.GatewayModel, len(gatewayEntsByKey))
	for key, ents := range gatewayEntsByKey {
		hwID := ents.magmadGateway.PhysicalID
		var devCasted *orc8rModels.GatewayDevice
		if devicesByID[hwID] != nil {
			devCasted = devicesByID[hwID].(*orc8rModels.GatewayDevice)
		}
		ret[key] = (&cwfModels.CwfGateway{}).FromBackendModels(ents.magmadGateway, ents.cwfGateway, devCasted, statusesByID[hwID])
	}
	return ret
}

func getSubscriberDirectoryHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	configuratorNetwork, err := configurator.LoadNetwork(networkID, false, false, serdes.Network)
	if err != nil {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if configuratorNetwork.Type != cwf.CwfNetworkType {
		return obsidian.HttpError(fmt.Errorf("NetworkID %s is not a CWF network", networkID), http.StatusBadRequest)
	}
	subscriberID := c.Param("subscriber_id")
	if subscriberID == "" {
		return obsidian.HttpError(fmt.Errorf("SubscriberID cannot be empty"), http.StatusBadRequest)
	}
	directoryState, err := state.GetState(networkID, orc8r.DirectoryRecordType, subscriberID, serdes.State)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	cwfRecord, err := convertDirectoryRecordToSubscriberRecord(directoryState.ReportedState)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cwfRecord)
}

func convertDirectoryRecordToSubscriberRecord(iRecord interface{}) (*cwfModels.CwfSubscriberDirectoryRecord, error) {
	record, ok := iRecord.(*directorydTypes.DirectoryRecord)
	if !ok {
		return nil, fmt.Errorf("Could not convert retrieved state to DirectoryRecord")
	}
	b, err := json.Marshal(record.Identifiers)
	if err != nil {
		return nil, err
	}
	c := &cwfModels.CwfSubscriberDirectoryRecord{}
	err = json.Unmarshal(b, c)
	if err != nil {
		return nil, fmt.Errorf("Error converting DirectoryRecord to CWF Record: %s, %v", err, *record)
	}
	c.LocationHistory = record.LocationHistory
	return c, nil
}

func getHAPairStatusHandler(c echo.Context) error {
	nid, haPairID, nerr := getNetworkIDAndHaPairID(c)
	if nerr != nil {
		return nerr
	}
	network, err := configurator.LoadNetwork(nid, true, true, serdes.Network)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if network.Type != cwf.CwfNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a <%s> network", nid, cwf.CwfNetworkType))
	}
	haPairStatus, err := getCwfHaPairStatus(nid, haPairID)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, haPairStatus)
}

func getHealthStatusHandler(c echo.Context) error {
	nid, gid, nerr := obsidian.GetNetworkAndGatewayIDs(c)
	if nerr != nil {
		return nerr
	}
	pid, err := configurator.GetPhysicalIDOfEntity(nid, orc8r.MagmadGatewayType, gid)
	if err == merrors.ErrNotFound || len(pid) == 0 {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	healthState, err := getCwfGatewayHealth(nid, gid)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, healthState)
}

func listHAPairsHandler(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	haPairEnts, _, err := configurator.LoadAllEntitiesOfType(nid, cwf.CwfHAPairType, configurator.FullEntityLoadCriteria(), serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	ret := make(map[string]*cwfModels.CwfHaPair, len(haPairEnts))
	for _, haPairEnt := range haPairEnts {
		cwfHaPair := &cwfModels.CwfHaPair{}
		err = cwfHaPair.FromBackendModels(haPairEnt)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
		cwfHaPair.State = getHaPairState(nid, cwfHaPair)
		ret[haPairEnt.Key] = cwfHaPair
	}
	return c.JSON(http.StatusOK, ret)
}

func createHAPairHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	haPair := new(cwfModels.MutableCwfHaPair)
	if err := c.Bind(haPair); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := haPair.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	_, err := configurator.CreateEntity(networkID, haPair.ToEntity(), serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, haPair.HaPairID)
}

func getHAPairHandler(c echo.Context) error {
	networkID, haPairID, nerr := getNetworkIDAndHaPairID(c)
	if nerr != nil {
		return nerr
	}
	ent, err := configurator.LoadEntity(
		networkID, cwf.CwfHAPairType, haPairID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	cwfHaPair := &cwfModels.CwfHaPair{}
	err = cwfHaPair.FromBackendModels(ent)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	cwfHaPair.State = getHaPairState(networkID, cwfHaPair)
	return c.JSON(http.StatusOK, cwfHaPair)
}

func updateHAPairHandler(c echo.Context) error {
	networkID, haPairID, nerr := getNetworkIDAndHaPairID(c)
	if nerr != nil {
		return nerr
	}
	mutableHaPair := new(cwfModels.MutableCwfHaPair)
	if err := c.Bind(mutableHaPair); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := mutableHaPair.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if mutableHaPair.HaPairID != haPairID {
		err := fmt.Errorf("ha pair ID from parameters (%s) and payload (%s) must match", haPairID, mutableHaPair.HaPairID)
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	// 404 if pair doesn't exist
	exists, err := configurator.DoesEntityExist(networkID, cwf.CwfHAPairType, haPairID)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if ha pair exists"), http.StatusInternalServerError)
	}
	if !exists {
		return echo.ErrNotFound
	}
	_, err = configurator.UpdateEntity(networkID, mutableHaPair.ToEntityUpdateCriteria(haPairID), serdes.Entity)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func deleteHAPairHandler(c echo.Context) error {
	networkID, haPairID, nerr := getNetworkIDAndHaPairID(c)
	if nerr != nil {
		return nerr
	}
	err := configurator.DeleteEntity(networkID, cwf.CwfHAPairType, haPairID)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusNoContent)
}

func getHaPairState(networkID string, haPair *cwfModels.CwfHaPair) *cwfModels.CarrierWifiHaPairState {
	ret := &cwfModels.CarrierWifiHaPairState{}
	gateway1Health, err := getCwfGatewayHealth(networkID, haPair.GatewayID1)
	if err == nil {
		ret.Gateway1Health = gateway1Health
	}
	gateway2Health, err := getCwfGatewayHealth(networkID, haPair.GatewayID2)
	if err == nil {
		ret.Gateway2Health = gateway2Health
	}
	status, err := getCwfHaPairStatus(networkID, haPair.HaPairID)
	if err == nil {
		ret.HaPairStatus = status
	}
	return ret
}

func getCwfGatewayHealth(networkID string, gatewayID string) (*cwfModels.CarrierWifiGatewayHealthStatus, error) {
	reportedGatewayState, err := state.GetState(networkID, cwf.CwfGatewayHealthType, gatewayID, serdes.State)
	if err == merrors.ErrNotFound {
		return nil, obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	healthState, ok := reportedGatewayState.ReportedState.(*cwfModels.CarrierWifiGatewayHealthStatus)
	if !ok {
		return nil, obsidian.HttpError(
			fmt.Errorf("could not convert retrieved type %T to CarrierWifiGatewayHealthStatus", reportedGatewayState.ReportedState),
			http.StatusInternalServerError,
		)
	}
	return healthState, nil
}

func getCwfHaPairStatus(networkID string, haPairID string) (*cwfModels.CarrierWifiHaPairStatus, error) {
	reportedHaPairStatus, err := state.GetState(networkID, cwf.CwfHAPairStatusType, haPairID, serdes.State)
	if err == merrors.ErrNotFound {
		return nil, obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return nil, obsidian.HttpError(err, http.StatusInternalServerError)
	}
	haPairStatus, ok := reportedHaPairStatus.ReportedState.(*cwfModels.CarrierWifiHaPairStatus)
	if !ok {
		return nil, obsidian.HttpError(
			fmt.Errorf("could not convert retrieved type %T to CarrierWifiHaPairStatus", reportedHaPairStatus.ReportedState),
			http.StatusInternalServerError,
		)
	}
	return haPairStatus, nil
}

func getNetworkIDAndHaPairID(c echo.Context) (string, string, *echo.HTTPError) {
	vals, err := obsidian.GetParamValues(c, "network_id", "ha_pair_id")
	if err != nil {
		return "", "", err
	}
	return vals[0], vals[1], nil
}

func makeErr(err error) *echo.HTTPError {
	if err == merrors.ErrNotFound {
		return echo.ErrNotFound
	}
	return obsidian.HttpError(err, http.StatusInternalServerError)
}
