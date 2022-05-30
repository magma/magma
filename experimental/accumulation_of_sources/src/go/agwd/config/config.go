// Copyright 2021 The Magma Authors.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"sync"

	"github.com/magma/magma/src/go/protos/magma/config"
	"google.golang.org/grpc/resolver"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	"github.com/magma/magma/src/go/agwd/config/internal/grpcutil"
	"github.com/magma/magma/src/go/log"
)

//go:generate go run github.com/golang/mock/mockgen -package config -destination mock_config/mock_config.go . Configer

// LogLevel translates protobuf defined config.AgwD_LogLevel to log.Level.
func LogLevel(l config.AgwD_LogLevel) log.Level {
	switch l {
	case config.AgwD_DEBUG:
		return log.DebugLevel
	case config.AgwD_INFO:
		return log.InfoLevel
	case config.AgwD_WARN:
		return log.WarnLevel
	case config.AgwD_ERROR:
		return log.ErrorLevel
	}
	return log.InfoLevel
}

const (
	ipv4Scheme = "ipv4"
	tcp4Scheme = "tcp4"
	ipv6Scheme = "ipv6"
	tcp6Scheme = "tcp6"
	TCP        = "tcp"
)

// ParseTarget takes a target in string form and returns a resolved Target.
// Extends functionality in grpc/internal/grpcutil.ParseTarget to support ipv4
// and ipv6 schemes. Return tcp4Scheme or tcp6Scheme so the Scheme can be passed
// directly to net.Listen.
func ParseTarget(target string) resolver.Target {
	if strings.HasPrefix(target, ipv4Scheme+":") {
		return resolver.Target{
			Scheme:   tcp4Scheme,
			Endpoint: target[len(ipv4Scheme)+1:],
		}
	}
	if strings.HasPrefix(target, ipv6Scheme+":") {
		return resolver.Target{
			Scheme:   tcp6Scheme,
			Endpoint: target[len(ipv6Scheme)+1:],
		}
	}
	if strings.HasPrefix(target, tcp4Scheme+":") {
		return resolver.Target{
			Scheme:   tcp4Scheme,
			Endpoint: target[len(tcp4Scheme)+1:],
		}
	}
	return grpcutil.ParseTarget(target, false)
}

// GetVagrantTarget concatenates the vagrant private ip with the port posses it
// into ParseTarget and returns the result.
func GetVagrantTarget(vagrantIP, port string) resolver.Target {
	return ParseTarget(fmt.Sprintf("tcp4:%s:%s", vagrantIP, port))
}

// Configer returns a parsed config.
type Configer interface {
	Config() *config.AgwD
	UpdateConfig(*config.AgwD) error
	ReplaceConfig(d *config.AgwD) error
}

// ConfigManager implements Configer via a loaded config.
type ConfigManager struct {
	config *config.AgwD

	sync.RWMutex
}

func newDefaultConfig() *config.AgwD {
	return &config.AgwD{
		LogLevel:                        config.AgwD_INFO,
		SctpdDownstreamServiceTarget:    "unix:///tmp/sctpd_downstream.sock",
		SctpdUpstreamServiceTarget:      "unix:///tmp/sctpd_upstream.sock",
		MmeSctpdDownstreamServiceTarget: "unix:///tmp/mme_sctpd_downstream.sock",
		MmeSctpdUpstreamServiceTarget:   "unix:///tmp/mme_sctpd_upstream.sock",
		PipelinedServiceTarget:          "tcp4:0.0.0.0:12345",
		// Sentry is disabled if DSN is not set.
		SentryDsn:               "",
		VagrantPrivateNetworkIp: "192.168.60.142",
		ConfigServicePort:       "6000",
		CaptureServicePort:      "6001",
		CaptureConfig:           &config.CaptureConfig{MatchSpecs: []*config.CaptureConfig_MatchSpec{}},
	}
}

const (
	cStyleCommentStart = "/*"
	cStyleCommentEnd   = "*/"
)

func filterCStyleComments(in string) string {
	var filtered string
	for {
		idx := strings.Index(in, cStyleCommentStart)
		if idx == -1 {
			filtered += in
			break
		}
		filtered += in[:idx]
		in = in[idx:]
		idx = strings.Index(in, cStyleCommentEnd)
		if idx == -1 {
			break
		}
		in = in[idx+2:]
	}
	return filtered
}

func filterJSONComments(in string) string {
	var filtered []string
	for _, line := range strings.Split(filterCStyleComments(in), "\n") {
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[0:idx]
		}
		line = strings.TrimSpace(line)
		if line != "" {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, "\n")
}

// NewConfigManager constructs a *ConfigManager with default config values.
func NewConfigManager() *ConfigManager {
	return &ConfigManager{config: newDefaultConfig()}
}

// Config returns the current config.
func (c *ConfigManager) Config() *config.AgwD {
	c.RLock()
	defer c.RUnlock()

	return c.config
}

// Merge updates the managed config.
func (c *ConfigManager) Merge(update *config.AgwD) {
	c.Lock()
	defer c.Unlock()

	// clone to prevent data race on proto fields
	config, ok := proto.Clone(c.config).(*config.AgwD)
	if !ok {
		panic("clone of config not *config.AgwD")
	}
	proto.Merge(config, update)
	c.config = config
}

// Replace overwrites the managed config.
func (c *ConfigManager) Replace(replace *config.AgwD) {
	c.Lock()
	defer c.Unlock()
	config, ok := proto.Clone(replace).(*config.AgwD)
	if !ok {
		panic("clone of config not *config.AgwD")
	}
	c.config = config
}

func loadConfigFile(
	osStat func(string) (os.FileInfo, error),
	readFile func(string) ([]byte, error),
	unmarshalProto func([]byte, proto.Message) error,
	path string,
) (*config.AgwD, error) {
	if _, err := osStat(path); err != nil {
		return nil, fmt.Errorf("path=%s: %w", path, err)
	}

	bytes, err := readFile(path)
	if err != nil {
		return nil, fmt.Errorf("path=%s: %w", path, err)
	}
	filtered := []byte(filterJSONComments(string(bytes)))
	config := &config.AgwD{}
	if err := unmarshalProto(filtered, config); err != nil {
		return nil, fmt.Errorf(
			"path=%s filtered=%s: %w",
			path,
			string(filtered),
			err)
	}
	return config, nil
}

// LoadConfigFile updates ConfigManager with a config file if it can be read
// successfully.
func LoadConfigFile(cm *ConfigManager, path string) error {
	loaded, err := loadConfigFile(
		os.Stat, ioutil.ReadFile, protojson.Unmarshal, path)
	if err != nil {
		return err
	}

	cm.Merge(loaded)
	return nil
}

// UpdateConfig updates CongfigManager with a config file.
// TODO: Add validation and error checks on fields being updated.
func (cm *ConfigManager) UpdateConfig(config *config.AgwD) error {
	cm.Merge(config)
	return nil
}

// ReplaceConfig overwites the CongfigManager with a config file.
// TODO: Add validation and error checks.
// Note: this can lead to a race condition.
func (cm *ConfigManager) ReplaceConfig(config *config.AgwD) error {
	cm.Replace(config)
	return nil
}
