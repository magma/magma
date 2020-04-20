package models

func NewDefaultNetworkCarrierWifiConfigs() *NetworkCarrierWifiConfigs {
	defaultRuleID := "default_rule_1"
	return &NetworkCarrierWifiConfigs{
		AaaServer: &AaaServer{
			AccountingEnabled:    false,
			CreateSessionOnAuth:  false,
			IDLESessionTimeoutMs: 21600000,
		},
		DefaultRuleID: &defaultRuleID,
		EapAka: &EapAka{
			PlmnIds: []string{"123456"},
			Timeout: &EapAkaTimeout{
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
