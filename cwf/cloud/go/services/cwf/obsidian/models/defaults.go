package models

import fegmodels "magma/feg/cloud/go/services/feg/obsidian/models"

func NewDefaultNetworkCarrierWifiConfigs() *NetworkCarrierWifiConfigs {
	defaultRuleID := "default_rule_1"
	return &NetworkCarrierWifiConfigs{
		AaaServer: &fegmodels.AaaServer{
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
			IdleSessionTimeoutMs: 21600000,
		},
		DefaultRuleID: &defaultRuleID,
		EapAka: &fegmodels.EapAka{
			PlmnIds: []string{"123456"},
			Timeout: &fegmodels.EapAkaTimeouts{
				ChallengeMs:            20000,
				ErrorNotificationMs:    10000,
				SessionAuthenticatedMs: 5000,
				SessionMs:              43200000,
			},
		},
		NetworkServices: []string{"dpi", "policy_enforcement"},
	}
}

func NewDefaultModifiedNetworkCarrierWifiConfigs() *NetworkCarrierWifiConfigs {
	configs := NewDefaultNetworkCarrierWifiConfigs()
	configs.AaaServer.AccountingEnabled = true
	return configs
}
