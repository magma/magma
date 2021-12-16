/*
 Copyright 2020 The Magma Authors.

 This source code is licensed under the BSD-style license found in the
 LICENSE file in the root directory of this source tree.

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package servicers

import (
	"context"
	"fmt"
	"sort"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
	"github.com/golang/glog"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"

	"magma/feg/cloud/go/feg"
	feg_serdes "magma/feg/cloud/go/serdes"
	feg_models "magma/feg/cloud/go/services/feg/obsidian/models"
	"magma/lte/cloud/go/lte"
	lte_mconfig "magma/lte/cloud/go/protos/mconfig"
	"magma/lte/cloud/go/serdes"
	lte_service "magma/lte/cloud/go/services/lte"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	nprobe_models "magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/orc8r/math"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/configurator/mconfig"
	builder_protos "magma/orc8r/cloud/go/services/configurator/mconfig/protos"
	"magma/orc8r/cloud/go/services/orchestrator/obsidian/models"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
)

type builderServicer struct {
	defaultSubscriberdbSyncInterval uint32
}

func NewBuilderServicer(config lte_service.Config) builder_protos.MconfigBuilderServer {
	return &builderServicer{
		defaultSubscriberdbSyncInterval: config.DefaultSubscriberdbSyncInterval,
	}
}

func (s *builderServicer) Build(ctx context.Context, request *builder_protos.BuildRequest) (*builder_protos.BuildResponse, error) {
	ret := &builder_protos.BuildResponse{ConfigsByKey: map[string][]byte{}}

	network, err := (configurator.Network{}).FromProto(request.Network, serdes.Network)
	if err != nil {
		glog.V(4).Infof("LTE mconfig not build (conversion error Network cast failed) '%v' for gateway: %s", err, request.GatewayId)
		return nil, err
	}
	graph, err := (configurator.EntityGraph{}).FromProto(request.Graph, serdes.Entity)
	if err != nil {
		glog.V(4).Infof("LTE mconfig not build (conversion error EntityGraph cast failed) '%v' for gateway: %s", err, request.GatewayId)
		return nil, err
	}
	// Only build mconfig if cellular network and gateway configs exist
	inwConfig, found := network.Configs[lte.CellularNetworkConfigType]
	if !found || inwConfig == nil {
		glog.V(4).Infof("LTE mconfig not build for %s: CellularNetworkConfigType not found", request.GatewayId)
		return ret, nil
	}
	cellularNwConfig := inwConfig.(*lte_models.NetworkCellularConfigs)

	cellGW, err := graph.GetEntity(lte.CellularGatewayEntityType, request.GatewayId)
	if err == merrors.ErrNotFound {
		glog.V(4).Infof("LTE mconfig not build for %s: CellularGatewayEntityType not found in graph", request.GatewayId)
		return ret, nil
	}
	if err != nil {
		return nil, err
	}

	if cellGW.Config == nil {
		glog.V(4).Infof("LTE mconfig not build for %s: CellularGatewayEntityType.Config is nill", request.GatewayId)
		return ret, nil
	}
	cellularGwConfig := cellGW.Config.(*lte_models.GatewayCellularConfigs)

	if err := validateConfigs(cellularNwConfig, cellularGwConfig); err != nil {
		return nil, err
	}

	federatedNetworkConfigs, err := getFederatedNetworkConfigs(network.Type, cellularNwConfig.FegNetworkID, request)
	if err != nil {
		glog.Errorf("Failed to retrieve LTE_federated network config while building lte mconfig for gateway %s", request.GatewayId)
		return nil, err
	}

	enodebs, err := graph.GetAllChildrenOfType(cellGW, lte.CellularEnodebEntityType)
	if err != nil {
		return nil, err
	}

	gwRan := cellularGwConfig.Ran
	gwEpc := cellularGwConfig.Epc
	gwNgc := getGwConfigNgc(cellularGwConfig)
	gwNonEpsService := cellularGwConfig.NonEpsService
	nwRan := cellularNwConfig.Ran
	nwEpc := cellularNwConfig.Epc
	nonEPSServiceMconfig := getNonEPSServiceMconfigFields(gwNonEpsService)

	pipelineDServices, err := getPipelineDServicesConfig(nwEpc.NetworkServices)
	if err != nil {
		return nil, err
	}

	enbConfigsBySerial := getEnodebConfigsBySerial(cellularNwConfig, cellularGwConfig, enodebs)
	heConfig := getHEConfig(cellularGwConfig.HeConfig)
	npTasks, liUes := getNetworkProbeConfig(ctx, network.ID)

	mmePoolRecord, mmeGroupID, err := getMMEPoolConfigs(network.ID, cellularGwConfig.Pooling, cellGW, graph)
	if err != nil {
		return nil, err
	}
	congestionControlEnabled := nwEpc.CongestionControlEnabled
	if gwEpc.CongestionControlEnabled != nil {
		congestionControlEnabled = gwEpc.CongestionControlEnabled
	}

	vals := map[string]proto.Message{
		"enodebd": &lte_mconfig.EnodebD{
			LogLevel:            protos.LogLevel_INFO,
			Pci:                 int32(gwRan.Pci),
			FddConfig:           getFddConfig(nwRan.FddConfig),
			TddConfig:           getTddConfig(nwRan.TddConfig),
			BandwidthMhz:        int32(nwRan.BandwidthMhz),
			AllowEnodebTransmit: swag.BoolValue(gwRan.TransmitEnabled),
			Tac:                 int32(nwEpc.Tac),
			PlmnidList:          fmt.Sprintf("%s%s", nwEpc.Mcc, nwEpc.Mnc),
			CsfbRat:             nonEPSServiceMconfig.csfbRat,
			Arfcn_2G:            nonEPSServiceMconfig.arfcn_2g,
			EnbConfigsBySerial:  enbConfigsBySerial,
		},
		"mobilityd": &lte_mconfig.MobilityD{
			LogLevel:                 protos.LogLevel_INFO,
			IpBlock:                  gwEpc.IPBlock,
			IpAllocatorType:          getMobilityDIPAllocator(nwEpc),
			Ipv6Block:                gwEpc.IPV6Block,
			Ipv6PrefixAllocationType: gwEpc.IPV6PrefixAllocationMode,
			StaticIpEnabled:          getMobilityDStaticIPAllocation(nwEpc),
			MultiApnIpAlloc:          getMobilityDMultuAPNIPAlloc(nwEpc),
		},
		"mme": &lte_mconfig.MME{
			LogLevel:                      protos.LogLevel_INFO,
			Mcc:                           nwEpc.Mcc,
			Mnc:                           nwEpc.Mnc,
			Tac:                           int32(nwEpc.Tac),
			MmeCode:                       int32(mmePoolRecord.MmeCode),
			MmeGid:                        int32(mmeGroupID),
			MmeRelativeCapacity:           int32(mmePoolRecord.MmeRelativeCapacity),
			EnableDnsCaching:              shouldEnableDNSCaching(cellularGwConfig.DNS),
			NonEpsServiceControl:          nonEPSServiceMconfig.nonEpsServiceControl,
			CsfbMcc:                       nonEPSServiceMconfig.csfbMcc,
			CsfbMnc:                       nonEPSServiceMconfig.csfbMnc,
			Lac:                           nonEPSServiceMconfig.lac,
			HssRelayEnabled:               swag.BoolValue(nwEpc.HssRelayEnabled),
			CloudSubscriberdbEnabled:      nwEpc.CloudSubscriberdbEnabled,
			AttachedEnodebTacs:            getEnodebTacs(enbConfigsBySerial),
			DnsPrimary:                    gwEpc.DNSPrimary,
			DnsSecondary:                  gwEpc.DNSSecondary,
			Ipv4PCscfAddress:              string(gwEpc.IPV4pCscfAddr),
			Ipv6DnsAddress:                string(gwEpc.IPV6DNSAddr),
			Ipv6PCscfAddress:              string(gwEpc.IPV6pCscfAddr),
			NatEnabled:                    swag.BoolValue(gwEpc.NatEnabled),
			Ipv4SgwS1UAddr:                gwEpc.IPV4SgwS1uAddr,
			RestrictedPlmns:               getRestrictedPlmns(nwEpc.RestrictedPlmns),
			RestrictedImeis:               getRestrictedImeis(nwEpc.RestrictedImeis),
			ServiceAreaMaps:               getServiceAreaMaps(nwEpc.ServiceAreaMaps),
			FederatedModeMap:              getFederatedModeMap(federatedNetworkConfigs),
			CongestionControlEnabled:      swag.BoolValue(congestionControlEnabled),
			SentryConfig:                  getNetworkSentryConfig(&network),
			Enable5GFeatures:              swag.BoolValue(nwEpc.Enable5gFeatures),
			AmfName:                       gwNgc.AmfName,
			AmfSetId:                      gwNgc.AmfSetID,
			AmfRegionId:                   gwNgc.AmfRegionID,
			AmfPointer:                    gwNgc.AmfPointer,
			AmfDefaultSliceServiceType:    gwNgc.AmfDefaultSst,
			AmfDefaultSliceDifferentiator: gwNgc.AmfDefaultSd,
		},
		"pipelined": &lte_mconfig.PipelineD{
			LogLevel:                   protos.LogLevel_INFO,
			UeIpBlock:                  gwEpc.IPBlock,
			NatEnabled:                 swag.BoolValue(gwEpc.NatEnabled),
			DefaultRuleId:              nwEpc.DefaultRuleID,
			Services:                   pipelineDServices,
			SgiManagementIfaceVlan:     gwEpc.SgiManagementIfaceVlan,
			SgiManagementIfaceIpAddr:   gwEpc.SgiManagementIfaceStaticIP,
			SgiManagementIfaceGw:       gwEpc.SgiManagementIfaceGw,
			SgiManagementIfaceIpv6Addr: gwEpc.SgiManagementIfaceIPV6Addr,
			SgiManagementIfaceIpv6Gw:   string(gwEpc.SgiManagementIfaceIPV6Gw),
			HeConfig:                   heConfig,
			LiUes:                      liUes,
			Enable5GFeatures:           swag.BoolValue(nwEpc.Enable5gFeatures),
			UpfNodeIdentifier:          string(nwEpc.NodeIdentifier),
		},
		"subscriberdb": &lte_mconfig.SubscriberDB{
			LogLevel:         protos.LogLevel_INFO,
			LteAuthOp:        nwEpc.LteAuthOp,
			LteAuthAmf:       nwEpc.LteAuthAmf,
			SubProfiles:      getSubProfiles(nwEpc),
			HssRelayEnabled:  swag.BoolValue(nwEpc.HssRelayEnabled),
			SyncInterval:     s.getRandomizedSyncInterval(cellGW.Key, nwEpc, gwEpc),
			Enable5GFeatures: swag.BoolValue(nwEpc.Enable5gFeatures),
		},
		"policydb": &lte_mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &lte_mconfig.SessionD{
			LogLevel:         protos.LogLevel_INFO,
			GxGyRelayEnabled: swag.BoolValue(nwEpc.GxGyRelayEnabled),
			WalletExhaustDetection: &lte_mconfig.WalletExhaustDetection{
				TerminateOnExhaust: false,
			},
			SentryConfig:     getNetworkSentryConfig(&network),
			Enable5GFeatures: swag.BoolValue(nwEpc.Enable5gFeatures),
		},
		"dnsd": getGatewayCellularDNSMConfig(cellularGwConfig.DNS),
		"liagentd": &lte_mconfig.LIAgentD{
			LogLevel:    protos.LogLevel_INFO,
			NprobeTasks: npTasks,
		},
	}

	ret.ConfigsByKey, err = mconfig.MarshalConfigs(vals)
	if err != nil {
		return nil, err
	}

	return ret, nil
}

func validateConfigs(nwConfig *lte_models.NetworkCellularConfigs, gwConfig *lte_models.GatewayCellularConfigs) error {
	if nwConfig == nil {
		return errors.New("Cellular network config is nil")
	}
	if gwConfig == nil {
		return errors.New("Cellular gateway config is nil")
	}

	if gwConfig.Ran == nil {
		return errors.New("Gateway RAN config is nil")
	}
	if gwConfig.Epc == nil {
		return errors.New("Gateway EPC config is nil")
	}
	if nwConfig.Ran == nil {
		return errors.New("Network RAN config is nil")
	}
	if nwConfig.Epc == nil {
		return errors.New("Network EPC config is nil")
	}
	return nil
}

type nonEPSServiceMconfigFields struct {
	csfbRat              lte_mconfig.EnodebD_CSFBRat
	arfcn_2g             []int32
	nonEpsServiceControl lte_mconfig.MME_NonEPSServiceControl
	csfbMcc              string
	csfbMnc              string
	lac                  int32
}

func getNonEPSServiceMconfigFields(gwNonEpsService *lte_models.GatewayNonEpsConfigs) nonEPSServiceMconfigFields {
	if gwNonEpsService == nil {
		return nonEPSServiceMconfigFields{
			csfbRat:              lte_mconfig.EnodebD_CSFBRAT_2G,
			arfcn_2g:             []int32{},
			nonEpsServiceControl: lte_mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			csfbMcc:              "",
			csfbMnc:              "",
			lac:                  1,
		}
	} else {
		arfcn2g := make([]int32, 0, len(gwNonEpsService.Arfcn2g))
		for _, i := range gwNonEpsService.Arfcn2g {
			arfcn2g = append(arfcn2g, int32(i))
		}

		return nonEPSServiceMconfigFields{
			csfbRat:              lte_mconfig.EnodebD_CSFBRat(swag.Uint32Value(gwNonEpsService.CsfbRat)),
			arfcn_2g:             arfcn2g,
			nonEpsServiceControl: lte_mconfig.MME_NonEPSServiceControl(swag.Uint32Value(gwNonEpsService.NonEpsServiceControl)),
			csfbMcc:              gwNonEpsService.CsfbMcc,
			csfbMnc:              gwNonEpsService.CsfbMnc,
			lac:                  int32(swag.Uint32Value(gwNonEpsService.Lac)),
		}
	}
}

var networkServicesByName = map[string]lte_mconfig.PipelineD_NetworkServices{
	"metering":           lte_mconfig.PipelineD_METERING,
	"dpi":                lte_mconfig.PipelineD_DPI,
	"policy_enforcement": lte_mconfig.PipelineD_ENFORCEMENT,
}

// move this out of this package eventually
func getPipelineDServicesConfig(networkServices []string) ([]lte_mconfig.PipelineD_NetworkServices, error) {
	if len(networkServices) == 0 {
		return []lte_mconfig.PipelineD_NetworkServices{
			lte_mconfig.PipelineD_ENFORCEMENT,
		}, nil
	}
	apps := make([]lte_mconfig.PipelineD_NetworkServices, 0, len(networkServices))
	for _, service := range networkServices {
		mc, found := networkServicesByName[service]
		if !found {
			return nil, errors.Errorf("unknown network service name %s", service)
		}
		apps = append(apps, mc)
	}
	return apps, nil
}

// getGwConfigNgc returns the NGC part of a cellular gateway config, or a
// default value if none exists.
func getGwConfigNgc(configs *lte_models.GatewayCellularConfigs) *lte_models.GatewayNgcConfigs {
	ngc := configs.Ngc
	if ngc == nil {
		ngc = &lte_models.GatewayNgcConfigs{}
	}
	return ngc
}

func getFddConfig(fddConfig *lte_models.NetworkRanConfigsFddConfig) *lte_mconfig.EnodebD_FDDConfig {
	if fddConfig == nil {
		return nil
	}
	return &lte_mconfig.EnodebD_FDDConfig{
		Earfcndl: int32(fddConfig.Earfcndl),
		Earfcnul: int32(fddConfig.Earfcnul),
	}
}

func getTddConfig(tddConfig *lte_models.NetworkRanConfigsTddConfig) *lte_mconfig.EnodebD_TDDConfig {
	if tddConfig == nil {
		return nil
	}

	return &lte_mconfig.EnodebD_TDDConfig{
		Earfcndl:               int32(tddConfig.Earfcndl),
		SubframeAssignment:     int32(tddConfig.SubframeAssignment),
		SpecialSubframePattern: int32(tddConfig.SpecialSubframePattern),
	}
}

func getHEConfig(gwConfig *lte_models.GatewayHeConfig) *lte_mconfig.PipelineD_HEConfig {
	if gwConfig == nil {
		return &lte_mconfig.PipelineD_HEConfig{}
	}

	return &lte_mconfig.PipelineD_HEConfig{
		EnableHeaderEnrichment: swag.BoolValue(gwConfig.EnableHeaderEnrichment),
		EnableEncryption:       swag.BoolValue(gwConfig.EnableEncryption),
		EncryptionAlgorithm:    lte_mconfig.PipelineD_HEConfig_EncryptionAlgorithm(lte_mconfig.PipelineD_HEConfig_EncryptionAlgorithm_value[gwConfig.HeEncryptionAlgorithm]),
		HashFunction:           lte_mconfig.PipelineD_HEConfig_HashFunction(lte_mconfig.PipelineD_HEConfig_HashFunction_value[gwConfig.HeHashFunction]),
		EncodingType:           lte_mconfig.PipelineD_HEConfig_EncodingType(lte_mconfig.PipelineD_HEConfig_EncodingType_value[gwConfig.HeEncodingType]),
		EncryptionKey:          gwConfig.EncryptionKey,
		HmacKey:                gwConfig.HmacKey,
	}
}

// getMMEPoolConfigs returns the gateway pool record and a uint32 specifying
// the MME group ID for a given gateway. If a gateway does not exist in a pool,
// default values are returned.
func getMMEPoolConfigs(networkID string, poolingConfig lte_models.CellularGatewayPoolRecords, cellGateway configurator.NetworkEntity, graph configurator.EntityGraph) (*lte_models.CellularGatewayPoolRecord, uint32, error) {
	// Currently, having multiple (mme group ID, mme code, mme relative
	// capacity) tuples is unsupported. As such, use the first pool record
	// to set all of these values.
	if len(poolingConfig) == 0 {
		return &lte_models.CellularGatewayPoolRecord{
			MmeCode:             1,
			MmeRelativeCapacity: 10,
		}, 1, nil
	}
	pool, err := graph.GetFirstAncestorOfType(cellGateway, lte.CellularGatewayPoolEntityType)
	if err != nil {
		return nil, 0, err
	}
	poolRecord := poolingConfig[0]
	cfg, ok := pool.Config.(*lte_models.CellularGatewayPoolConfigs)
	if !ok {
		err := fmt.Errorf("unable to convert gateway pool config for pool '%s'; pool has invalid config", pool.Key)
		return nil, 0, err
	}
	return poolRecord, cfg.MmeGroupID, nil
}

func getEnodebConfigsBySerial(nwConfig *lte_models.NetworkCellularConfigs, gwConfig *lte_models.GatewayCellularConfigs, enodebs []configurator.NetworkEntity) map[string]*lte_mconfig.EnodebD_EnodebConfig {
	ret := make(map[string]*lte_mconfig.EnodebD_EnodebConfig, len(enodebs))
	for _, ent := range enodebs {
		serial := ent.Key
		ienbConfig := ent.Config
		if ienbConfig == nil {
			glog.Errorf("enb with serial %s is missing config", serial)
		}

		enodebConfig := ienbConfig.(*lte_models.EnodebConfig)
		enbMconfig := &lte_mconfig.EnodebD_EnodebConfig{}

		if enodebConfig.ConfigType == "MANAGED" {
			cellularEnbConfig := enodebConfig.ManagedConfig
			enbMconfig.Earfcndl = int32(cellularEnbConfig.Earfcndl)
			enbMconfig.SubframeAssignment = int32(cellularEnbConfig.SubframeAssignment)
			enbMconfig.SpecialSubframePattern = int32(cellularEnbConfig.SpecialSubframePattern)
			enbMconfig.Pci = int32(cellularEnbConfig.Pci)
			enbMconfig.TransmitEnabled = swag.BoolValue(cellularEnbConfig.TransmitEnabled)
			enbMconfig.DeviceClass = cellularEnbConfig.DeviceClass
			enbMconfig.BandwidthMhz = int32(cellularEnbConfig.BandwidthMhz)
			enbMconfig.Tac = int32(cellularEnbConfig.Tac)
			enbMconfig.CellId = int32(swag.Uint32Value(cellularEnbConfig.CellID))

			// override zero values with network/gateway configs
			if enbMconfig.Earfcndl == 0 {
				enbMconfig.Earfcndl = int32(nwConfig.GetEarfcndl())
			}
			if enbMconfig.SubframeAssignment == 0 {
				if nwConfig.Ran.TddConfig != nil {
					enbMconfig.SubframeAssignment = int32(nwConfig.Ran.TddConfig.SubframeAssignment)
				}
			}
			if enbMconfig.SpecialSubframePattern == 0 {
				if nwConfig.Ran.TddConfig != nil {
					enbMconfig.SpecialSubframePattern = int32(nwConfig.Ran.TddConfig.SpecialSubframePattern)
				}
			}
			if enbMconfig.Pci == 0 {
				enbMconfig.Pci = int32(gwConfig.Ran.Pci)
			}
			if enbMconfig.BandwidthMhz == 0 {
				enbMconfig.BandwidthMhz = int32(nwConfig.Ran.BandwidthMhz)
			}
			if enbMconfig.Tac == 0 {
				enbMconfig.Tac = int32(nwConfig.Epc.Tac)
			}

		} else if enodebConfig.ConfigType == "UNMANAGED" {
			cellularEnbConfig := enodebConfig.UnmanagedConfig
			enbMconfig.CellId = int32(swag.Uint32Value(cellularEnbConfig.CellID))
			enbMconfig.Tac = int32(swag.Uint32Value(cellularEnbConfig.Tac))
			enbMconfig.IpAddress = string(*cellularEnbConfig.IPAddress)

			if enbMconfig.Tac == 0 {
				enbMconfig.Tac = int32(nwConfig.Epc.Tac)
			}
		}

		ret[serial] = enbMconfig
	}
	return ret
}

func getEnodebTacs(enbConfigsBySerial map[string]*lte_mconfig.EnodebD_EnodebConfig) []int32 {
	ret := make([]int32, 0, len(enbConfigsBySerial))
	for _, enbConfig := range enbConfigsBySerial {
		ret = append(ret, enbConfig.Tac)
	}
	sort.Slice(ret, func(i, j int) bool { return ret[i] < ret[j] })
	return ret
}

func getSubProfiles(epc *lte_models.NetworkEpcConfigs) map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile {
	if epc.SubProfiles == nil {
		return map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile{}
	}

	ret := map[string]*lte_mconfig.SubscriberDB_SubscriptionProfile{}
	for name, profile := range epc.SubProfiles {
		ret[name] = &lte_mconfig.SubscriberDB_SubscriptionProfile{
			MaxUlBitRate: profile.MaxUlBitRate,
			MaxDlBitRate: profile.MaxDlBitRate,
		}
	}
	return ret
}

func getMobilityDIPAllocator(epc *lte_models.NetworkEpcConfigs) lte_mconfig.MobilityD_IpAllocatorType {
	if epc.Mobility == nil {
		return lte_mconfig.MobilityD_IP_POOL
	}
	if epc.Mobility.IPAllocationMode == lte_models.DHCPBroadcastAllocationMode {
		return lte_mconfig.MobilityD_DHCP
	}
	// For other modes set IP pool allocator
	return lte_mconfig.MobilityD_IP_POOL
}

func getMobilityDStaticIPAllocation(epc *lte_models.NetworkEpcConfigs) bool {
	if epc.Mobility == nil {
		return false
	}
	return epc.Mobility.EnableStaticIPAssignments
}

func getMobilityDMultuAPNIPAlloc(epc *lte_models.NetworkEpcConfigs) bool {
	if epc.Mobility == nil {
		return false
	}
	return epc.Mobility.EnableMultiApnIPAllocation
}

func getGatewayCellularDNSMConfig(gwDns *lte_models.GatewayDNSConfigs) *lte_mconfig.DnsD {
	if gwDns == nil {
		return &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: true,
			EnableCaching:     false,
			LocalTTL:          0,
			Records:           []*lte_mconfig.GatewayDNSConfigRecordsItems{},
		}
	} else {
		return &lte_mconfig.DnsD{
			LogLevel:          protos.LogLevel_INFO,
			DhcpServerEnabled: swag.BoolValue(gwDns.DhcpServerEnabled),
			EnableCaching:     shouldEnableDNSCaching(gwDns),
			LocalTTL:          *gwDns.LocalTTL,
			Records:           getGatewayDnsRecords(gwDns),
		}
	}
}

func getGatewayDnsRecords(dns *lte_models.GatewayDNSConfigs) []*lte_mconfig.GatewayDNSConfigRecordsItems {
	if dns.Records == nil {
		return []*lte_mconfig.GatewayDNSConfigRecordsItems{}
	}

	ret := make([]*lte_mconfig.GatewayDNSConfigRecordsItems, 0, len(dns.Records))
	for _, record := range dns.Records {
		recordProto := &lte_mconfig.GatewayDNSConfigRecordsItems{}
		recordProto.Domain = record.Domain
		recordProto.ARecord = funk.Map(record.ARecord, func(a strfmt.IPv4) string { return string(a) }).([]string)
		recordProto.AaaaRecord = funk.Map(record.AaaaRecord, func(a strfmt.IPv6) string { return string(a) }).([]string)
		recordProto.CnameRecord = make([]string, 0, len(record.CnameRecord))
		recordProto.CnameRecord = append(recordProto.CnameRecord, record.CnameRecord...)
		ret = append(ret, recordProto)
	}
	return ret
}

func shouldEnableDNSCaching(dns *lte_models.GatewayDNSConfigs) bool {
	if dns == nil {
		return false
	}
	return swag.BoolValue(dns.EnableCaching)
}

func getRestrictedPlmns(plmns []*lte_models.PlmnConfig) []*lte_mconfig.MME_PlmnConfig {
	ret := make([]*lte_mconfig.MME_PlmnConfig, len(plmns))
	for idx, plmn := range plmns {
		ret[idx] = &lte_mconfig.MME_PlmnConfig{Mcc: plmn.Mcc, Mnc: plmn.Mnc}
	}
	return ret
}

func getServiceAreaMaps(serviceAreaMaps map[string]lte_models.TacList) map[string]*lte_mconfig.MME_TacList {
	ret := make(map[string]*lte_mconfig.MME_TacList)
	for k, v := range serviceAreaMaps {
		tacList := &lte_mconfig.MME_TacList{}
		for _, tac := range v {
			tacList.Tac = append(tacList.Tac, uint32(tac))
		}
		ret[k] = tacList
	}
	return ret
}

// getFederatedNetworkConfigs in case this is a federated LTE networkm this function will try to parse out
// feg_models.FederatedNetworkConfigs out of it
func getFederatedNetworkConfigs(networkType string, fegId lte_models.FegNetworkID, request *builder_protos.BuildRequest) (*feg_models.FederatedNetworkConfigs, error) {
	if networkType != feg.FederatedLteNetworkType {
		// this is a non federated network, return nothing
		return nil, nil
	}
	if fegId == "" {
		glog.Warning("federated_id is empty. Ignoring Federated LTE Network config and movign on")
		return nil, nil
	}
	network, err := (configurator.Network{}).FromProto(request.Network, feg_serdes.Network)
	if err != nil {
		return nil, err
	}
	inwConfig, found := network.Configs[feg.FederatedNetworkType]
	if !found || inwConfig == nil {
		return nil, err
	}
	return inwConfig.(*feg_models.FederatedNetworkConfigs), nil
}

// getFederatedModeMap extracts the mapping configuration in case of being a federated network
func getFederatedModeMap(fedNetworkConfigs *feg_models.FederatedNetworkConfigs) *lte_mconfig.FederatedModeMap {
	if fedNetworkConfigs == nil {
		return nil
	}
	return feg_models.ToFederatedModesMap(fedNetworkConfigs.FederatedModesMapping)
}

func getRestrictedImeis(imeis []*lte_models.Imei) []*lte_mconfig.MME_ImeiConfig {
	ret := make([]*lte_mconfig.MME_ImeiConfig, len(imeis))
	for idx, imei := range imeis {
		ret[idx] = &lte_mconfig.MME_ImeiConfig{Tac: imei.Tac, Snr: imei.Snr}
	}
	return ret
}

func getNetworkProbeConfig(ctx context.Context, networkID string) ([]*lte_mconfig.NProbeTask, *lte_mconfig.PipelineD_LiUes) {
	liUes := &lte_mconfig.PipelineD_LiUes{}
	npTasks := []*lte_mconfig.NProbeTask{}
	ents, _, err := configurator.LoadAllEntitiesOfType(
		ctx,
		networkID,
		lte.NetworkProbeTaskEntityType,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err != nil {
		glog.Errorf("Failed to load nprobe task entities %v", err)
		return npTasks, liUes
	}

	for _, ent := range ents {
		task := (&nprobe_models.NetworkProbeTask{}).FromBackendModels(ent)
		if task.TaskDetails.DeliveryType == nprobe_models.NetworkProbeTaskDetailsDeliveryTypeEventsOnly {
			// data plane is not requested.
			continue
		}

		npTasks = append(npTasks, nprobe_models.ToMConfigNProbeTask(task))
		switch task.TaskDetails.TargetType {
		case nprobe_models.NetworkProbeTaskDetailsTargetTypeImsi:
			liUes.Imsis = append(liUes.Imsis, task.TaskDetails.TargetID)
		case nprobe_models.NetworkProbeTaskDetailsTargetTypeImei:
			liUes.Imeis = append(liUes.Imeis, task.TaskDetails.TargetID)
		case nprobe_models.NetworkProbeTaskDetailsTargetTypeMsisdn:
			liUes.Msisdns = append(liUes.Msisdns, task.TaskDetails.TargetID)
		}
	}
	return npTasks, liUes
}

func getNetworkSentryConfig(network *configurator.Network) *lte_mconfig.SentryConfig {
	iSentryConfig, found := network.Configs[orc8r.NetworkSentryConfig]
	if !found || iSentryConfig == nil {
		return nil
	}
	sentryConfig, ok := iSentryConfig.(*models.NetworkSentryConfig)
	if !ok {
		return nil
	}
	return &lte_mconfig.SentryConfig{
		SampleRate:   swag.Float32Value(sentryConfig.SampleRate),
		UploadMmeLog: sentryConfig.UploadMmeLog,
		DsnNative:    string(sentryConfig.URLNative),
		DsnPython:    string(sentryConfig.URLPython),
	}
}

// getSyncInterval takes network-wide subscriberdb sync interval in seconds and overrides it if also set for gateway.
// If sync interval is unset for both network and gateway, a default is read from lte/cloud/configs/lte.yml
func (s *builderServicer) getSyncInterval(nwEpc *lte_models.NetworkEpcConfigs, gwEpc *lte_models.GatewayEpcConfigs) uint32 {
	// minSyncInterval enforces a minimum sync interval to prevent too many
	// sync requests if operator sets the default in lte.yml to lower than 60
	const minSyncInterval = 60
	gwSyncInterval := uint32(gwEpc.SubscriberdbSyncInterval)
	nwSyncInterval := uint32(nwEpc.SubscriberdbSyncInterval)

	if gwSyncInterval >= minSyncInterval {
		return gwSyncInterval
	}
	if nwSyncInterval >= minSyncInterval {
		return nwSyncInterval
	}
	if s.defaultSubscriberdbSyncInterval >= minSyncInterval {
		return s.defaultSubscriberdbSyncInterval
	}
	return minSyncInterval
}

// getRandomizedSyncInterval returns the interval received from getSyncInterval
// as seconds and increases it by a random jitter in the range of
// [0, 0.2 * getSyncInterval()]. Increased sync interval ameliorates the thundering
// herd effect at the Orc8r.
func (s *builderServicer) getRandomizedSyncInterval(gwKey string, nwEpc *lte_models.NetworkEpcConfigs, gwEpc *lte_models.GatewayEpcConfigs) uint32 {
	syncInterval := s.getSyncInterval(nwEpc, gwEpc)
	return math.JitterUint32(syncInterval, gwKey, 0.2)
}
