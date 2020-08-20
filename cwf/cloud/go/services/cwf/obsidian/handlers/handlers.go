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
	cwfModels "magma/cwf/cloud/go/services/cwf/obsidian/models"
	fegModels "magma/feg/cloud/go/services/feg/obsidian/models"
	lteHandlers "magma/lte/cloud/go/services/lte/obsidian/handlers"
	policyModels "magma/lte/cloud/go/services/policydb/obsidian/models"
	"magma/orc8r/cloud/go/models"
	"magma/orc8r/cloud/go/obsidian"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/directoryd"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/handlers"
	orc8rModels "magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/go-openapi/swag"
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
	ManageNetworkSubscriberPath    = ManageNetworkPath + obsidian.UrlSep + "subscriber_config"
	ManageNetworkBaseNamesPath     = ManageNetworkSubscriberPath + obsidian.UrlSep + "base_names"
	ManageNetworkRuleNamesPath     = ManageNetworkSubscriberPath + obsidian.UrlSep + "rule_names"
	ManageNetworkBaseNamePath      = ManageNetworkBaseNamesPath + obsidian.UrlSep + ":base_name"
	ManageNetworkRuleNamePath      = ManageNetworkRuleNamesPath + obsidian.UrlSep + ":rule_id"
	ManageNetworkClusterStatusPath = ManageNetworkPath + obsidian.UrlSep + "cluster_status"
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
		handlers.GetListGatewaysHandler(ListGatewaysPath, &cwfModels.MutableCwfGateway{}, makeCwfGateways),
		{Path: ListGatewaysPath, Methods: obsidian.POST, HandlerFunc: createGateway},
		{Path: ManageGatewayPath, Methods: obsidian.GET, HandlerFunc: getGateway},
		{Path: ManageGatewayPath, Methods: obsidian.PUT, HandlerFunc: updateGateway},
		{Path: ManageGatewayPath, Methods: obsidian.DELETE, HandlerFunc: deleteGateway},

		{Path: ManageGatewayStatePath, Methods: obsidian.GET, HandlerFunc: handlers.GetStateHandler},
		{Path: ManageNetworkClusterStatusPath, Methods: obsidian.GET, HandlerFunc: getClusterStatusHandler},
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

	ret = append(ret, handlers.GetTypedNetworkCRUDHandlers(ListNetworksPath, ManageNetworkPath, cwf.CwfNetworkType, &cwfModels.CwfNetwork{})...)

	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkNamePath, new(models.NetworkName), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDescriptionPath, new(models.NetworkDescription), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFeaturesPath, &orc8rModels.NetworkFeatures{}, orc8r.NetworkFeaturesConfig)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkDNSPath, &orc8rModels.NetworkDNSConfig{}, orc8r.DnsdNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkCarrierWifiPath, &cwfModels.NetworkCarrierWifiConfigs{}, cwf.CwfNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkFederationPath, &fegModels.FederatedNetworkConfigs{}, cwf.CwfNetworkType)...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkSubscriberPath, &policyModels.NetworkSubscriberConfig{}, "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkRuleNamesPath, new(policyModels.RuleNames), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkBaseNamesPath, new(policyModels.BaseNames), "")...)
	ret = append(ret, handlers.GetPartialNetworkHandlers(ManageNetworkLiUesPath, new(cwfModels.LiUes), "")...)

	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayNamePath, new(models.GatewayName))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayDescriptionPath, new(models.GatewayDescription))...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayConfigPath, &orc8rModels.MagmadGatewayConfigs{})...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayTierPath, new(orc8rModels.TierID))...)
	ret = append(ret, handlers.GetGatewayDeviceHandlers(ManageGatewayDevicePath)...)
	ret = append(ret, handlers.GetPartialGatewayHandlers(ManageGatewayCarrierWifiPath, &cwfModels.GatewayCwfConfigs{})...)

	return ret
}

func createGateway(c echo.Context) error {
	if nerr := handlers.CreateGateway(c, &cwfModels.MutableCwfGateway{}); nerr != nil {
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
	if nerr = handlers.UpdateGateway(c, nid, gid, &cwfModels.MutableCwfGateway{}); nerr != nil {
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
	configuratorNetwork, err := configurator.LoadNetwork(networkID, false, false)
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
	directoryState, err := state.GetState(networkID, orc8r.DirectoryRecordType, subscriberID)
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
	record, ok := iRecord.(*directoryd.DirectoryRecord)
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

func getClusterStatusHandler(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	network, err := configurator.LoadNetwork(nid, true, true)
	if err == merrors.ErrNotFound {
		return c.NoContent(http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	if network.Type != cwf.CwfNetworkType {
		return echo.NewHTTPError(http.StatusBadRequest, fmt.Sprintf("network %s is not a <%s> network", nid, cwf.CwfNetworkType))
	}
	reportedClusterStatus, err := state.GetState(nid, cwf.CwfClusterHealthType, "cluster")
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	clusterStatus, ok := reportedClusterStatus.ReportedState.(*cwfModels.CarrierWifiNetworkClusterStatus)
	if !ok {
		return obsidian.HttpError(fmt.Errorf("could not convert reported retrieved state to ClusterStatus"), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, clusterStatus)
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
	reportedGatewayState, err := state.GetState(nid, cwf.CwfGatewayHealthType, gid)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	} else if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	healthState, ok := reportedGatewayState.ReportedState.(*cwfModels.CarrierWifiGatewayHealthStatus)
	if !ok {
		return obsidian.HttpError(fmt.Errorf("could not convert reported retrieved state to HealthStatus"), http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, healthState)
}

func listHAPairsHandler(c echo.Context) error {
	nid, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	ids, err := configurator.ListEntityKeys(nid, cwf.CwfHAPairType)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	haPairTKs := make([]storage.TypeAndKey, 0, len(ids))
	for _, id := range ids {
		haPairTKs = append(haPairTKs, storage.TypeAndKey{Type: cwf.CwfHAPairType, Key: id})
	}
	haPairEnts, _, err := configurator.LoadEntities(
		nid,
		swag.String(cwf.CwfHAPairType),
		nil,
		nil,
		haPairTKs,
		configurator.FullEntityLoadCriteria(),
	)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	haPairEntsByTK := haPairEnts.MakeByTK()
	ret := make(map[string]*cwfModels.CwfHaPair, len(haPairEntsByTK))
	for tk, haPair := range haPairEntsByTK {
		ret[tk.Key], err = (&cwfModels.CwfHaPair{}).FromBackendModels(haPair)
		if err != nil {
			return obsidian.HttpError(err, http.StatusInternalServerError)
		}
	}
	return c.JSON(http.StatusOK, ret)
}

func createHAPairHandler(c echo.Context) error {
	networkID, nerr := obsidian.GetNetworkId(c)
	if nerr != nil {
		return nerr
	}
	haPair := new(cwfModels.CwfHaPair)
	if err := c.Bind(haPair); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	_, err := configurator.CreateEntity(networkID, haPair.ToEntity())
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
		networkID,
		cwf.CwfHAPairType,
		haPairID,
		configurator.EntityLoadCriteria{LoadConfig: true, LoadAssocsFromThis: true},
	)
	if err == merrors.ErrNotFound {
		return obsidian.HttpError(err, http.StatusNotFound)
	}
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	cwfHaPair, err := (&cwfModels.CwfHaPair{}).FromBackendModels(ent)
	if err != nil {
		return obsidian.HttpError(err, http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK, cwfHaPair)
}

func updateHAPairHandler(c echo.Context) error {
	networkID, haPairID, nerr := getNetworkIDAndHaPairID(c)
	if nerr != nil {
		return nerr
	}
	haPair := new(cwfModels.CwfHaPair)
	if err := c.Bind(haPair); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if err := haPair.ValidateModel(); err != nil {
		return obsidian.HttpError(err, http.StatusBadRequest)
	}
	if haPairID != haPair.HaPairID {
		return obsidian.HttpError(errors.New("ha pair in body does not match URL param"), http.StatusBadRequest)
	}
	// 404 if pair doesn't exist
	exists, err := configurator.DoesEntityExist(networkID, cwf.CwfHAPairType, haPairID)
	if err != nil {
		return obsidian.HttpError(errors.Wrap(err, "Failed to check if ha pair exists"), http.StatusInternalServerError)
	}
	if !exists {
		return echo.ErrNotFound
	}
	_, err = configurator.UpdateEntity(networkID, haPair.ToEntityUpdateCriteria())
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
