/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"errors"
	"fmt"
	"log"
	"reflect"

	"magma/lte/cloud/go/protos/mconfig"
	cellular_protos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
	dnsd_config "magma/orc8r/cloud/go/services/dnsd/config"
	dsnd_protos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/proto"
)

type CellularBuilder struct{}

type NonEPSServiceMconfigFields struct {
	csfbRat              mconfig.EnodebD_CSFBRat
	arfcn_2g             []int32
	nonEpsServiceControl mconfig.MME_NonEPSServiceControl
	csfbMcc              string
	csfbMnc              string
	lac                  int32
}

func (builder *CellularBuilder) Build(networkId string, gatewayId string) (map[string]proto.Message, error) {
	emptyRet := map[string]proto.Message{}
	cellularNwConfig, err := getCellularNetworkConfig(networkId)
	if err != nil {
		return nil, err
	}
	if cellularNwConfig == nil {
		return emptyRet, nil
	}
	cellularGwConfig, err := getCellularGatewayConfig(networkId, gatewayId)
	if err != nil {
		return nil, err
	}
	if cellularGwConfig == nil {
		return emptyRet, nil
	}

	// Add if DNS caching should be enabled
	// allows for MME to be restarted when this is changed
	networkDNSconfig, err := dnsd_config.GetNetworkDNSConfig(networkId)
	if err != nil {
		return nil, err
	}
	enableDNSCaching := shouldEnableDNSCaching(networkDNSconfig)

	if err := validateConfigs(cellularNwConfig, cellularGwConfig); err != nil {
		return nil, err
	}

	enbSerialArr := cellularGwConfig.GetAttachedEnodebSerials()
	enodebConfigsBySerial, err := getEnodebConfigsBySerial(networkId, enbSerialArr)
	if err != nil {
		enodebConfigsBySerial = map[string]*mconfig.EnodebD_EnodebConfig{}
	}

	// All guaranteed non-nil by the above check, except gwNonEpsService
	gwRan := cellularGwConfig.GetRan()
	gwEpc := cellularGwConfig.GetEpc()
	gwNonEpsService := cellularGwConfig.GetNonEpsService()
	nwRan := cellularNwConfig.GetRan()
	nwEpc := cellularNwConfig.GetEpc()

	nonEPSServiceMconfig := getNonEPSServiceMconfigFields(gwNonEpsService)

	pipelineDServices, err := cellular_protos.GetPipelineDServicesConfig(nwEpc.GetNetworkServices())
	if err != nil {
		return nil, err
	}

	return map[string]proto.Message{
		"enodebd": &mconfig.EnodebD{
			LogLevel:               protos.LogLevel_INFO,
			Pci:                    gwRan.GetPci(),
			Earfcndl:               nwRan.GetEarfcndl(),
			FddConfig:              getFddConfig(nwRan.GetFddConfig()),
			TddConfig:              getTddConfig(nwRan.GetTddConfig()),
			SubframeAssignment:     nwRan.GetSubframeAssignment(),
			SpecialSubframePattern: nwRan.GetSpecialSubframePattern(),
			BandwidthMhz:           nwRan.GetBandwidthMhz(),
			AllowEnodebTransmit:    gwRan.GetTransmitEnabled(),
			Tac:                    nwEpc.GetTac(),
			PlmnidList:             fmt.Sprintf("%s%s", nwEpc.GetMcc(), nwEpc.GetMnc()),
			CsfbRat:                nonEPSServiceMconfig.csfbRat,
			Arfcn_2G:               nonEPSServiceMconfig.arfcn_2g,
			EnbConfigsBySerial:     enodebConfigsBySerial,
		},
		"mobilityd": &mconfig.MobilityD{
			LogLevel: protos.LogLevel_INFO,
			IpBlock:  gwEpc.GetIpBlock(),
		},
		"mme": &mconfig.MME{
			LogLevel:                 protos.LogLevel_INFO,
			Mcc:                      nwEpc.GetMcc(),
			Mnc:                      nwEpc.GetMnc(),
			Tac:                      nwEpc.GetTac(),
			MmeCode:                  1,
			MmeGid:                   1,
			EnableDnsCaching:         enableDNSCaching,
			NonEpsServiceControl:     nonEPSServiceMconfig.nonEpsServiceControl,
			CsfbMcc:                  nonEPSServiceMconfig.csfbMcc,
			CsfbMnc:                  nonEPSServiceMconfig.csfbMnc,
			Lac:                      nonEPSServiceMconfig.lac,
			RelayEnabled:             nwEpc.GetRelayEnabled(),
			CloudSubscriberdbEnabled: nwEpc.GetCloudSubscriberdbEnabled(),
		},
		"pipelined": &mconfig.PipelineD{
			LogLevel:      protos.LogLevel_INFO,
			UeIpBlock:     gwEpc.GetIpBlock(),
			NatEnabled:    gwEpc.GetNatEnabled(),
			DefaultRuleId: nwEpc.GetDefaultRuleId(),
			RelayEnabled:  nwEpc.GetRelayEnabled(),
			Services:      pipelineDServices,
		},
		"subscriberdb": &mconfig.SubscriberDB{
			LogLevel:     protos.LogLevel_INFO,
			LteAuthOp:    nwEpc.GetLteAuthOp(),
			LteAuthAmf:   nwEpc.GetLteAuthAmf(),
			SubProfiles:  getSubProfiles(nwEpc),
			RelayEnabled: nwEpc.GetRelayEnabled(),
		},
		"policydb": &mconfig.PolicyDB{
			LogLevel: protos.LogLevel_INFO,
		},
		"sessiond": &mconfig.SessionD{
			LogLevel:     protos.LogLevel_INFO,
			RelayEnabled: nwEpc.GetRelayEnabled(),
		},
	}, nil
}

func getFddConfig(fddConfig *cellular_protos.NetworkRANConfig_FDDConfig) *mconfig.EnodebD_FDDConfig {
	if fddConfig != nil {
		return &mconfig.EnodebD_FDDConfig{
			Earfcndl: fddConfig.GetEarfcndl(),
			Earfcnul: fddConfig.GetEarfcnul(),
		}
	}
	return nil
}

func getTddConfig(tddConfig *cellular_protos.NetworkRANConfig_TDDConfig) *mconfig.EnodebD_TDDConfig {
	if tddConfig != nil {
		return &mconfig.EnodebD_TDDConfig{
			Earfcndl:               tddConfig.GetEarfcndl(),
			SubframeAssignment:     tddConfig.GetSubframeAssignment(),
			SpecialSubframePattern: tddConfig.GetSpecialSubframePattern(),
		}
	}
	return nil
}

func getEnodebConfigsBySerial(
	networkID string,
	enbSerialArr []string,
) (map[string]*mconfig.EnodebD_EnodebConfig, error) {
	nEnb := len(enbSerialArr)
	enbConfigMap := make(map[string]*mconfig.EnodebD_EnodebConfig, nEnb)
	for i := 0; i < nEnb; i++ {
		enbSerial := enbSerialArr[i]
		enbConfig, err := getEnodebConfig(networkID, enbSerial)
		if err != nil {
			// Just exclude the config if we cannot fetch it, probably because it is missing
			log.Printf("Missing config for eNB serial %s", enbSerial)
		} else {
			enbConfigMap[enbSerial] = enbConfig
		}
	}
	return enbConfigMap, nil
}

func getEnodebConfig(
	networkID string,
	enbSerialID string,
) (*mconfig.EnodebD_EnodebConfig, error) {
	cellularEnbConfigStruct, err := config.GetConfig(networkID, CellularEnodebType, enbSerialID)
	if err != nil {
		return nil, err
	}
	if cellularEnbConfigStruct == nil {
		return nil, fmt.Errorf("Missing config for network %s, serial %s", networkID, enbSerialID)
	}
	cellularEnbConfig := cellularEnbConfigStruct.(*cellular_protos.CellularEnodebConfig)
	return &mconfig.EnodebD_EnodebConfig{
		Earfcndl:               cellularEnbConfig.GetEarfcndl(),
		SubframeAssignment:     cellularEnbConfig.GetSubframeAssignment(),
		SpecialSubframePattern: cellularEnbConfig.GetSpecialSubframePattern(),
		Pci:                    cellularEnbConfig.GetPci(),
		TransmitEnabled:        cellularEnbConfig.GetTransmitEnabled(),
		DeviceClass:            cellularEnbConfig.GetDeviceClass(),
		BandwidthMhz:           cellularEnbConfig.GetBandwidthMhz(),
		Tac:                    cellularEnbConfig.GetTac(),
		CellId:                 cellularEnbConfig.GetCellId(),
	}, nil
}

func getCellularNetworkConfig(networkId string) (*cellular_protos.CellularNetworkConfig, error) {
	iCellularNwConfigs, err := config.GetConfig(networkId, CellularNetworkType, networkId)
	if err != nil || iCellularNwConfigs == nil {
		return nil, err
	}
	cellularNwConfigs, ok := iCellularNwConfigs.(*cellular_protos.CellularNetworkConfig)
	if !ok {
		return nil, fmt.Errorf(
			"Received unexpected type for network record. "+
				"Expected *CellularNetworkConfig but got %s",
			reflect.TypeOf(iCellularNwConfigs),
		)
	}
	return cellularNwConfigs, nil
}

func getCellularGatewayConfig(networkId string, gatewayId string) (*cellular_protos.CellularGatewayConfig, error) {
	iGatewayConfigs, err := config.GetConfig(networkId, CellularGatewayType, gatewayId)
	if err != nil || iGatewayConfigs == nil {
		return nil, err
	}
	gatewayConfigs, ok := iGatewayConfigs.(*cellular_protos.CellularGatewayConfig)
	if !ok {
		return nil, fmt.Errorf(
			"received unexpected type for gateway record. "+
				"Expected *CellularGatewayConfig but got %s",
			reflect.TypeOf(iGatewayConfigs),
		)
	}
	return gatewayConfigs, nil
}

func getNonEPSServiceMconfigFields(gwNonEpsService *cellular_protos.GatewayNonEPSConfig) NonEPSServiceMconfigFields {
	if gwNonEpsService == nil {
		return NonEPSServiceMconfigFields{
			csfbRat:              mconfig.EnodebD_CSFBRAT_2G,
			arfcn_2g:             []int32{},
			nonEpsServiceControl: mconfig.MME_NON_EPS_SERVICE_CONTROL_OFF,
			csfbMcc:              "",
			csfbMnc:              "",
			lac:                  1,
		}
	} else {
		return NonEPSServiceMconfigFields{
			csfbRat:              mconfig.EnodebD_CSFBRat(gwNonEpsService.GetCsfbRat()),
			arfcn_2g:             gwNonEpsService.GetArfcn_2G(),
			nonEpsServiceControl: mconfig.MME_NonEPSServiceControl(gwNonEpsService.GetNonEpsServiceControl()),
			csfbMcc:              gwNonEpsService.GetCsfbMcc(),
			csfbMnc:              gwNonEpsService.GetCsfbMnc(),
			lac:                  gwNonEpsService.GetLac(),
		}
	}
}

func getSubProfiles(nwEpc *cellular_protos.NetworkEPCConfig) map[string]*mconfig.SubscriberDB_SubscriptionProfile {
	subProfiles := make(map[string]*mconfig.SubscriberDB_SubscriptionProfile)
	if nwEpc.GetSubProfiles() != nil {
		for name, profile := range nwEpc.GetSubProfiles() {
			subProfiles[name] = &mconfig.SubscriberDB_SubscriptionProfile{
				MaxUlBitRate: profile.MaxUlBitRate,
				MaxDlBitRate: profile.MaxDlBitRate,
			}
		}
	}
	return subProfiles
}

func shouldEnableDNSCaching(dnsConfig *dsnd_protos.NetworkDNSConfig) bool {
	if dnsConfig == nil {
		return false
	} else {
		return dnsConfig.GetEnableCaching()
	}
}

func validateConfigs(nwConfig *cellular_protos.CellularNetworkConfig, gwConfig *cellular_protos.CellularGatewayConfig) error {
	if nwConfig == nil {
		return errors.New("Cellular network config is nil")
	}
	if gwConfig == nil {
		return errors.New("Cellular gateway config is nil")
	}

	if gwConfig.GetRan() == nil {
		return errors.New("Gateway RAN config is nil")
	}
	if gwConfig.GetEpc() == nil {
		return errors.New("Gateway EPC config is nil")
	}
	if nwConfig.GetRan() == nil {
		return errors.New("Network RAN config is nil")
	}
	if nwConfig.GetEpc() == nil {
		return errors.New("Network EPC config is nil")
	}
	return nil
}
