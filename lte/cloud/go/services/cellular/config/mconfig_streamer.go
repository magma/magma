/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package config

import (
	"fmt"

	"magma/lte/cloud/go/protos/mconfig"
	cellularprotos "magma/lte/cloud/go/services/cellular/protos"
	"magma/orc8r/cloud/go/protos"
	"magma/orc8r/cloud/go/services/config"
	"magma/orc8r/cloud/go/services/config/streaming"
	dnsd_config "magma/orc8r/cloud/go/services/dnsd/config"
	dnsd_protos "magma/orc8r/cloud/go/services/dnsd/protos"

	"github.com/golang/protobuf/ptypes"
)

// Subset of mconfig fields that this streamer manages
var managedFields = []string{
	"enodebd",
	"mobilityd",
	"mme",
	"pipelined",
	"subscriberdb",
	"directoryd",
	"policydb",
	"sessiond",
}

type CellularStreamer struct{}

func (*CellularStreamer) GetSubscribedConfigTypes() []string {
	return []string{CellularNetworkType, CellularGatewayType, dnsd_config.DnsdNetworkType}
}

func (*CellularStreamer) SeedNewGatewayMconfig(networkId string, gatewayId string, mconfigOut *protos.GatewayConfigs) error {
	// Seed with network config
	nwCfg, err := config.GetConfig(networkId, CellularNetworkType, networkId)
	if err != nil {
		return err
	}
	if nwCfg == nil {
		return nil
	}
	return applyNwConfigUpdate(streaming.CreateOperation, nwCfg.(*cellularprotos.CellularNetworkConfig), mconfigOut)
}

func (*CellularStreamer) ApplyMconfigUpdate(
	update *streaming.ConfigUpdate,
	oldMconfigsByGatewayId map[string]*protos.GatewayConfigs,
) (map[string]*protos.GatewayConfigs, error) {
	switch update.ConfigType {
	case CellularNetworkType:
		newValueCasted := castConfigValueToCellularNetwork(update.NewValue)
		for _, mconfigValue := range oldMconfigsByGatewayId {
			err := applyNwConfigUpdate(update.Operation, newValueCasted, mconfigValue)
			if err != nil {
				return oldMconfigsByGatewayId, err
			}
		}

		return oldMconfigsByGatewayId, nil
	case CellularGatewayType:
		newValueCasted := castConfigValueToCellularGateway(update.NewValue)
		for _, mconfigValue := range oldMconfigsByGatewayId {
			err := applyGwConfigUpdate(update.Operation, newValueCasted, mconfigValue)
			if err != nil {
				return oldMconfigsByGatewayId, err
			}
		}

		return oldMconfigsByGatewayId, nil
	case dnsd_config.DnsdNetworkType:
		newValueCasted := castConfigValueToDnsdNetwork(update.NewValue)
		for _, mconfigValue := range oldMconfigsByGatewayId {
			err := applyDnsdConfigUpdate(update.Operation, newValueCasted, mconfigValue)
			if err != nil {
				return oldMconfigsByGatewayId, err
			}
		}

		return oldMconfigsByGatewayId, nil
	default:
		return oldMconfigsByGatewayId, fmt.Errorf(
			"Cellular streamer encountered uncrecognized config type: %s",
			update.ConfigType,
		)
	}
}

func castConfigValueToCellularNetwork(v interface{}) *cellularprotos.CellularNetworkConfig {
	if v == nil {
		return nil
	}
	return v.(*cellularprotos.CellularNetworkConfig)
}

func castConfigValueToCellularGateway(v interface{}) *cellularprotos.CellularGatewayConfig {
	if v == nil {
		return nil
	}
	return v.(*cellularprotos.CellularGatewayConfig)
}

func castConfigValueToDnsdNetwork(v interface{}) *dnsd_protos.NetworkDNSConfig {
	if v == nil {
		return nil
	}
	return v.(*dnsd_protos.NetworkDNSConfig)
}

// Struct to hold the subset of mconfig fields that this streamer manages
type cellularPartialMconfig struct {
	// Export the fields for easier access via reflect
	Enodebd      *mconfig.EnodebD
	Mobilityd    *mconfig.MobilityD
	Mme          *mconfig.MME
	Pipelined    *mconfig.PipelineD
	Subscriberdb *mconfig.SubscriberDB
	Policydb     *mconfig.PolicyDB
	Sessiond     *mconfig.SessionD
}

func applyNwConfigUpdate(
	operation streaming.ChangeOperation,
	newConfig *cellularprotos.CellularNetworkConfig,
	mconfigOut *protos.GatewayConfigs, // output param
) error {
	switch operation {
	case streaming.DeleteOperation:
		for _, field := range managedFields {
			delete(mconfigOut.ConfigsByKey, field)
		}
		return nil
	case streaming.ReadOperation, streaming.UpdateOperation, streaming.CreateOperation:
		partialFields := &cellularPartialMconfig{}
		err := streaming.GetPartialMconfig(mconfigOut, partialFields)
		if err != nil {
			return err
		}
		setDefaultFields(partialFields)

		epc := getNwEpc(newConfig)
		ran := getNwRan(newConfig)

		partialFields.Enodebd.Earfcndl = ran.Earfcndl
		partialFields.Enodebd.SubframeAssignment = ran.SubframeAssignment
		partialFields.Enodebd.SpecialSubframePattern = ran.SpecialSubframePattern

		if ran.TddConfig != nil {
			partialFields.Enodebd.TddConfig = &mconfig.EnodebD_TDDConfig{
				Earfcndl:               ran.TddConfig.Earfcndl,
				SpecialSubframePattern: ran.TddConfig.SpecialSubframePattern,
				SubframeAssignment:     ran.TddConfig.SubframeAssignment,
			}
		}
		if ran.FddConfig != nil {
			partialFields.Enodebd.FddConfig = &mconfig.EnodebD_FDDConfig{
				Earfcndl: ran.FddConfig.Earfcndl,
				Earfcnul: ran.FddConfig.Earfcnul,
			}
		}
		partialFields.Enodebd.BandwidthMhz = ran.BandwidthMhz
		partialFields.Enodebd.PlmnidList = fmt.Sprintf("%s%s", epc.Mcc, epc.Mnc)
		partialFields.Enodebd.Tac = epc.Tac

		partialFields.Mme.Mcc = epc.Mcc
		partialFields.Mme.Mnc = epc.Mnc
		partialFields.Mme.Tac = epc.Tac
		partialFields.Mme.RelayEnabled = epc.RelayEnabled

		partialFields.Pipelined.DefaultRuleId = epc.DefaultRuleId
		partialFields.Pipelined.RelayEnabled = epc.RelayEnabled
		pipelineDServices, err := cellularprotos.GetPipelineDServicesConfig(epc.GetNetworkServices())
		if err != nil {
			return err
		}
		partialFields.Pipelined.Services = pipelineDServices

		partialFields.Sessiond.RelayEnabled = epc.RelayEnabled

		partialFields.Subscriberdb.LteAuthOp = epc.LteAuthOp
		partialFields.Subscriberdb.LteAuthAmf = epc.LteAuthAmf
		partialFields.Subscriberdb.SubProfiles = getSubProfiles(epc)
		partialFields.Subscriberdb.RelayEnabled = epc.RelayEnabled

		return streaming.UpdateMconfig(partialFields, mconfigOut)
	default:
		return fmt.Errorf("Unrecognized stream change operation %s", operation)
	}
}

func applyGwConfigUpdate(
	operation streaming.ChangeOperation,
	newConfig *cellularprotos.CellularGatewayConfig,
	mconfigOut *protos.GatewayConfigs, // output param
) error {
	switch operation {
	case streaming.DeleteOperation:
		for _, field := range managedFields {
			delete(mconfigOut.ConfigsByKey, field)
		}
		return nil
	case streaming.ReadOperation, streaming.UpdateOperation, streaming.CreateOperation:
		partialFields := &cellularPartialMconfig{}
		err := streaming.GetPartialMconfig(mconfigOut, partialFields)
		if err != nil {
			return err
		}
		setDefaultFields(partialFields)

		epc := getGwEpc(newConfig)
		ran := getGwRan(newConfig)
		eps := getNonEPSServiceMconfigFields(newConfig.NonEpsService)

		partialFields.Enodebd.Pci = ran.Pci
		partialFields.Enodebd.AllowEnodebTransmit = ran.TransmitEnabled
		partialFields.Enodebd.CsfbRat = eps.csfbRat
		partialFields.Enodebd.Arfcn_2G = eps.arfcn_2g

		partialFields.Mobilityd.IpBlock = epc.IpBlock

		partialFields.Mme.NonEpsServiceControl = eps.nonEpsServiceControl
		partialFields.Mme.CsfbMcc = eps.csfbMcc
		partialFields.Mme.CsfbMnc = eps.csfbMnc
		partialFields.Mme.Lac = eps.lac

		partialFields.Pipelined.UeIpBlock = epc.IpBlock
		partialFields.Pipelined.NatEnabled = epc.NatEnabled

		return streaming.UpdateMconfig(partialFields, mconfigOut)
	default:
		return fmt.Errorf("Unrecognized stream change operation %s", operation)
	}
}

func applyDnsdConfigUpdate(
	operation streaming.ChangeOperation,
	newConfig *dnsd_protos.NetworkDNSConfig,
	mconfigOut *protos.GatewayConfigs, // output param
) error {
	switch operation {
	case streaming.DeleteOperation:
		// In this case, don't delete the cellular config - just set the
		// dns caching to disabled (default/empty value)
		return setEnableDnsCaching(mconfigOut, false)
	case streaming.ReadOperation, streaming.UpdateOperation, streaming.CreateOperation:
		return setEnableDnsCaching(mconfigOut, newConfig.EnableCaching)
	default:
		return fmt.Errorf("Unrecognized stream change operation %s", operation)
	}
}

// These fields are constant
func setDefaultFields(partialConfigOut *cellularPartialMconfig) {
	partialConfigOut.Enodebd.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Mobilityd.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Mme.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Pipelined.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Subscriberdb.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Policydb.LogLevel = protos.LogLevel_INFO
	partialConfigOut.Sessiond.LogLevel = protos.LogLevel_INFO

	partialConfigOut.Mme.MmeCode = 1
	partialConfigOut.Mme.MmeGid = 1
}

func getNwEpc(nwCellularConfig *cellularprotos.CellularNetworkConfig) *cellularprotos.NetworkEPCConfig {
	if nwCellularConfig.Epc == nil {
		return &cellularprotos.NetworkEPCConfig{}
	} else {
		return nwCellularConfig.Epc
	}
}

func getNwRan(nwCellularConfig *cellularprotos.CellularNetworkConfig) *cellularprotos.NetworkRANConfig {
	if nwCellularConfig.Ran == nil {
		return &cellularprotos.NetworkRANConfig{}
	} else {
		return nwCellularConfig.Ran
	}
}

func getGwEpc(gwCellularConfig *cellularprotos.CellularGatewayConfig) *cellularprotos.GatewayEPCConfig {
	if gwCellularConfig.Epc == nil {
		return &cellularprotos.GatewayEPCConfig{}
	} else {
		return gwCellularConfig.Epc
	}
}

func getGwRan(gwCellularConfig *cellularprotos.CellularGatewayConfig) *cellularprotos.GatewayRANConfig {
	if gwCellularConfig.Ran == nil {
		return &cellularprotos.GatewayRANConfig{}
	} else {
		return gwCellularConfig.Ran
	}
}

func setEnableDnsCaching(mconfigOut *protos.GatewayConfigs, enableDnsCaching bool) error {
	mmeCfg, err := getMmeMconfigOrDefault(mconfigOut)
	if err != nil {
		return err
	}
	mmeCfg.EnableDnsCaching = enableDnsCaching

	mmeAny, err := ptypes.MarshalAny(mmeCfg)
	if err != nil {
		return err
	}

	mconfigOut.ConfigsByKey["mme"] = mmeAny
	return nil
}

func getMmeMconfigOrDefault(cfg *protos.GatewayConfigs) (*mconfig.MME, error) {
	ret := &mconfig.MME{}
	if cfg == nil || cfg.ConfigsByKey == nil {
		return ret, nil
	}

	existing, exists := cfg.ConfigsByKey["mme"]
	if exists {
		err := ptypes.UnmarshalAny(existing, ret)
		return ret, err
	} else {
		return ret, nil
	}
}
