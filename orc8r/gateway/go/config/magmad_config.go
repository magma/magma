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

package config

import (
	"strings"
	"time"

	"github.com/golang/glog"

	"magma/orc8r/lib/go/definitions"
	"magma/orc8r/lib/go/service/config"
)

const (
	MAGMAD_YML_FRESHNESS_CHECK_INTERVAL = time.Minute * 10

	// Defaults
	DefaultChallengeKeyFile = "/var/opt/magma/certs/gw_challenge.key"
	DefaultStaticConfigDir  = "/etc/magma"
	DefaultDynamicConfigDir = "/var/opt/magma/configs"
)

// BootstrapConfig bootstrapper related configuration - `yaml:"bootstrap_config"`
type BootstrapConfig struct {
	ChallengeKey string `yaml:"challenge_key"`
}

// MagmadCfg represents magmad.yml based configuration
type MagmadCfg struct {
	LogLevel                         string   `yaml:"log_level"`
	MagmaServices                    []string `yaml:"magma_services"`
	NonService303Services            []string `yaml:"non_service303_services"`
	RegisteredDynamicServices        []string `yaml:"registered_dynamic_services"`
	SkipCheckinIfMissingMetaServices []string `yaml:"skip_checkin_if_missing_meta_services"`
	InitSystem                       string   `yaml:"init_system"`
	// When cloud managed configs (gateway.cmconfig) are loaded by a magma service, the service first tries to
	// load them from dynamic (most recent) configs directory - DynamicMconfigDir, if unsuccessful, the service
	// falls back to static configs directory - StaticMconfigDir; this allows services to operate
	StaticMconfigDir  string `yaml:"static_mconfig_dir"`
	DynamicMconfigDir string `yaml:"dynamic_mconfig_dir"`
	// StaticMconfigUpdateIntervalMin specifies interval in minutes dynamic gateway.mconfig from DynamicMconfigDir
	// will be synchronized with (copied to) static gateway.mconfig in StaticMconfigDir
	// if StaticMconfigUpdateIntervalMin <= 0 (default) - static gateway.mconfig in StaticMconfigDir will never
	// be overwritten
	StaticMconfigUpdateIntervalMin int                  `yaml:"static_mconfig_update_interval_minutes"`
	BootstrapConfig                BootstrapConfig      `yaml:"bootstrap_config"`
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

// UpgraderFactory is upgrader_factory configuration block from magmad.yml
type UpgraderFactory struct {
	Module string `yaml:"module"`
	Class  string `yaml:"class"`
}

// Metricsd is metricsd configuration block from magmad.yml
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
	Module        string                  `yaml:"module"`
	Class         string                  `yaml:"class"`
	ShellCommands []ShellCommand          `yaml:"shell_commands"`
	CommandsMap   map[string]ShellCommand `yaml:"-"`
}

// ShellCommand magmad shell command definition
type ShellCommand struct {
	Name        string
	Command     string
	AllowParams bool   `yaml:"allow_params"`
	CommandFmt  string `yaml:"-"`
}

// UpdateShellCmdMap creates a new CommandsMap and populates it from ShellCommands list
// UpdateShellCmdMap also does basic format conversion
func (gcf *GenericCommandConfig) UpdateShellCmdMap() {
	if gcf != nil {
		gcf.CommandsMap = map[string]ShellCommand{}
		for _, cmd := range gcf.ShellCommands {
			cmd.CommandFmt = strings.ReplaceAll(cmd.Command, `{}`, `%v`)
			gcf.CommandsMap[cmd.Name] = cmd
		}
	}
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
		StaticMconfigDir:                 DefaultStaticConfigDir,
		DynamicMconfigDir:                DefaultDynamicConfigDir,
		BootstrapConfig:                  BootstrapConfig{ChallengeKey: DefaultChallengeKeyFile},
		EnableConfigStreamer:             true,
		EnableUpgradeMamager:             false,
		EnableNetworkMonitor:             false,
		EnableSystemdTailer:              false,
		EnableSyncRpc:                    true,
		EnableKernelVersionChecking:      false,
		SystemdTailerPollInterval:        30,
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

// UpdateFromYml of StructuredConfign interface - updates given magmad config struct from corresponding YML file
// returns updated MagmadCfg, main YML CFG file path & overwrite YML CFG file path (if any)
func (mdc *MagmadCfg) UpdateFromYml() (StructuredConfig, string, string) {
	var newCfg *MagmadCfg
	if mdc != nil {
		newCfg = &MagmadCfg{}
		*newCfg = *mdc // copy current configs
	} else {
		newCfg = NewDefaultMgmadCfg()
		mdc = newCfg
	}
	ymlFile, ymlOWFile, err := config.GetStructuredServiceConfig("", definitions.MagmadServiceName, newCfg)
	if err != nil {
		glog.Warningf("Error Getting Magmad Configs: %v,\n\tcontinue using old configs: %+v", err, mdc)
	} else {
		if mdc != newCfg { // success, copy if needed
			*mdc = *newCfg
		}
		mdc.GenericCommandConfig.UpdateShellCmdMap()
	}
	return mdc, ymlFile, ymlOWFile
}

// FreshnessCheckInterval of StructuredConfig interface
func (_ *MagmadCfg) FreshnessCheckInterval() time.Duration {
	return MAGMAD_YML_FRESHNESS_CHECK_INTERVAL
}

var magmadConfigs AtomicStore

func magmadCfgFactory() StructuredConfig {
	return NewDefaultMgmadCfg()
}

// GetMagmadConfigs returns current magmad configuration
func GetMagmadConfigs() *MagmadCfg {
	return magmadConfigs.GetCurrent(magmadCfgFactory).(*MagmadCfg)
}

// OverwriteMagmadConfigs overwrites current magmad configs
func OverwriteMagmadConfigs(cfg *MagmadCfg) {
	magmadConfigs.Overwrite(cfg)
}
