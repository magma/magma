/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

// package service implements the core of configurator
package service

import (
	"log"

	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/service/config"
)

// MagmadCfg represents magmad.yml based configuration
type MagmadCfg struct {
	LogLevel                         string   `yaml:"log_level"`
	MagmaServices                    []string `yaml:"magma_services"`
	NonService303Services            []string `yaml:"non_service303_services"`
	RegisteredDynamicServices        []string `yaml:"registered_dynamic_services"`
	SkipCheckinIfMissingMetaServices []string `yaml:"skip_checkin_if_missing_meta_services"`
	InitSystem                       string   `yaml:"init_system"`
	BootstrapConfig                  struct {
		ChallengeKey string `yaml:"challenge_key"`
	} `yaml:"bootstrap_config"`

	EnableConfigStreamer           bool                 `yaml:"enable_config_streamer"`
	EnableUpgradeMamager           bool                 `yaml:"enable_upgrade_manager"`
	EnableNetworkMonitor           bool                 `yaml:"enable_network_monitor"`
	EnableSystemdTailer            bool                 `yaml:"enable_systemd_tailer"`
	EnableSyncRpc                  bool                 `yaml:"enable_sync_rpc"`
	EnableKernelVersionChecking    bool                 `yaml:"enable_kernel_version_checking"`
	SystemdTailerPollInterval      int                  `yaml:"systemd_tailer_poll_interval"`
	NetworkMonitorConfig           NetworkMonitorConfig `yaml:"network_monitor_config"`
	UpgraderFactory                UpgraderFactory      `yaml:"upgrader_factory"`
	MconfigModules                 []string             `yaml:"mconfig_modules"`
	Metricsd                       Metricsd             `yaml:"metricsd"`
	GenericCommandConfig           GenericCommandConfig `yaml:"generic_command_config"`
	ConfigStreamErrorRetryInterval int                  `yaml:"config_stream_error_retry_interval"`
}

// NetworkMonitorConfig is network_monitor_config configuration block from magmad.yml
type NetworkMonitorConfig struct {
	SamplingPeriod int `yaml:"sampling_period"`
	PingConfig     struct {
		Hosts       []string `yaml:"hosts"`
		NumPackets  int      `yaml:"num_packets"`
		TimeoutSecs int      `yaml:"timeout_secs"`
	} `yaml:"ping_config"`
}

// NetworkMonitorConfig is upgrader_factory configuration block from magmad.yml
type UpgraderFactory struct {
	Module string `yaml:"module"`
	Class  string `yaml:"class"`
}

// UpgraderFactory is metricsd configuration block from magmad.yml
type Metricsd struct {
	LogLevel        string   `yaml:"log_level"`
	CollectInterval int      `yaml:"collect_interval"`
	SyncInterval    int      `yaml:"sync_interval"`
	GrpcTimeout     int      `yaml:"grpc_timeout"`
	QueueLength     int      `yaml:"queue_length"`
	Services        []string `yaml:"services"`
}

// GenericCommandConfig is generic_command_config configuration block from magmad.yml
type GenericCommandConfig struct {
	Module        string         `yaml:"module"`
	Class         string         `yaml:"class"`
	ShellCommands []ShellCommand `yaml:"shell_commands"`
}

// ShellCommand magmad shell command definition
type ShellCommand struct {
	Name        string
	Command     string
	AllowParams string `yaml:"allow_params"`
}

// NewDefaultMgmadCfg returns new default magmad configs
func NewDefaultMgmadCfg() *MagmadCfg {
	return &MagmadCfg{
		LogLevel:                         "INFO",
		MagmaServices:                    []string{},
		NonService303Services:            []string{},
		RegisteredDynamicServices:        []string{},
		SkipCheckinIfMissingMetaServices: []string{},
		InitSystem:                       "",
		BootstrapConfig: struct {
			ChallengeKey string `yaml:"challenge_key"`
		}{
			ChallengeKey: "/var/opt/magma/certs/gw_challenge.key"},
		EnableConfigStreamer:        true,
		EnableUpgradeMamager:        false,
		EnableNetworkMonitor:        false,
		EnableSystemdTailer:         false,
		EnableSyncRpc:               true,
		EnableKernelVersionChecking: false,
		SystemdTailerPollInterval:   30,
		NetworkMonitorConfig: NetworkMonitorConfig{
			SamplingPeriod: 60,
			PingConfig: struct {
				Hosts       []string `yaml:"hosts"`
				NumPackets  int      `yaml:"num_packets"`
				TimeoutSecs int      `yaml:"timeout_secs"`
			}{
				Hosts:       []string{"8.8.8.8"},
				NumPackets:  1,
				TimeoutSecs: 20,
			},
		},
		UpgraderFactory: UpgraderFactory{},
		MconfigModules:  []string{},
		Metricsd: Metricsd{
			LogLevel:        "INFO",
			CollectInterval: 60,
			SyncInterval:    60,
			GrpcTimeout:     30,
			QueueLength:     1000,
			Services:        []string{},
		},
		GenericCommandConfig:           GenericCommandConfig{},
		ConfigStreamErrorRetryInterval: 60,
	}
}

func (mdc *MagmadCfg) updateFromMagmadCfg() *MagmadCfg {
	newCfg := *mdc // copy current configs
	err := config.GetStructuredServiceConfig("", definitions.MagmadServiceName, &newCfg)
	if err != nil {
		log.Printf("Error Getting Magmad Configs: %v,\n\tcontinue using old configs: %+v", err, mdc)
	} else {
		// success, copy over the new configs
		*mdc = newCfg
	}
	return mdc
}
