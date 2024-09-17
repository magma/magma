/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package n7

import (
	"fmt"
	"net"
	"net/url"

	"magma/feg/gateway/sbi"

	"github.com/golang/glog"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/utils"
	"magma/gateway/mconfig"
)

const (
	N7N40ProxyServiceName = "n7_n40_proxy"

	PcfApiRoot            = "PCF_API_ROOT"
	PcfTokenUrl           = "PCF_TOKEN_URL"
	PcfClientId           = "PCF_CLIENT_ID"
	PcfClientSecret       = "PCF_CLIENT_SECRET"
	N7ClientLocalAddr     = "N7_CONSUMER_LOCAL_ADDR"
	N7ClientNotifyApiRoot = "N&_CONSUMER_NOTIFY_API_ROOT"

	DefaultPcfApiRoot      = "https://localhost"
	DefaultPcfTokenUrl     = "https://localhost/token"
	DefaultClientId        = "magma_client_id"
	DefaultClientSecret    = "magma_client_secret"
	DefaultN7ClientAddr    = "localhost:0"
	DefaultN7ClientApiRoot = "https://localhost/npcf-smpolicycontrol/v1"
)

type N7Config struct {
	DisableN7    bool
	ServerConfig sbi.RemoteConfig
	ClientConfig sbi.NotifierConfig
}

func GetN7Config() (*N7Config, error) {
	configPtr := &mcfgprotos.N7N40ProxyConfig{}
	conf := &N7Config{}

	err := mconfig.GetServiceConfigs(N7N40ProxyServiceName, configPtr)
	if err != nil || !validManagedConfig(configPtr) {
		glog.V(2).Infof("Managed Configs Load Error: %v Using EnvVars", err)
		apiRoot, err := url.ParseRequestURI(utils.GetValueOrEnv("", PcfApiRoot, DefaultPcfApiRoot))
		if err != nil {
			return nil, fmt.Errorf("invalid NotifierServer ApiRoot - %s", err)
		}
		conf.ServerConfig = sbi.RemoteConfig{
			ApiRoot:      *apiRoot,
			TokenUrl:     utils.GetValueOrEnv("", PcfTokenUrl, DefaultPcfTokenUrl),
			ClientId:     utils.GetValueOrEnv("", PcfClientId, DefaultClientId),
			ClientSecret: utils.GetValueOrEnv("", PcfClientSecret, DefaultClientSecret),
		}
		conf.DisableN7 = false
		conf.ClientConfig = sbi.NotifierConfig{
			LocalAddr:     utils.GetValueOrEnv("", N7ClientLocalAddr, DefaultN7ClientAddr),
			NotifyApiRoot: utils.GetValueOrEnv("", N7ClientNotifyApiRoot, DefaultN7ClientApiRoot),
		}
	} else {
		n7configPtr := configPtr.N7Config
		conf.DisableN7 = n7configPtr.DisableN7
		apiRoot, err := url.ParseRequestURI(utils.GetValueOrEnv("", PcfApiRoot, n7configPtr.Server.GetApiRoot()))
		if err != nil {
			return nil, fmt.Errorf("invalid NotifierServer ApiRoot - %s", err)
		}
		conf.ServerConfig = sbi.RemoteConfig{
			ApiRoot:      *apiRoot,
			TokenUrl:     utils.GetValueOrEnv("", PcfTokenUrl, n7configPtr.Server.GetTokenUrl()),
			ClientId:     utils.GetValueOrEnv("", PcfClientId, n7configPtr.Server.GetClientId()),
			ClientSecret: utils.GetValueOrEnv("", PcfClientSecret, n7configPtr.Server.GetClientSecret()),
		}
		conf.ClientConfig = sbi.NotifierConfig{
			LocalAddr:     utils.GetValueOrEnv("", N7ClientLocalAddr, n7configPtr.Client.LocalAddr),
			NotifyApiRoot: utils.GetValueOrEnv("", N7ClientNotifyApiRoot, n7configPtr.Client.NotifyApiRoot),
		}
	}
	err = validateN7Config(conf)
	if err != nil {
		return nil, err
	}
	return conf, nil
}

func validManagedConfig(config *mcfgprotos.N7N40ProxyConfig) bool {
	if config.N7Config == nil || config.N7Config.Server == nil || config.N7Config.Client == nil {
		return false
	}
	return true
}

func validateN7Config(config *N7Config) error {
	_, err := url.ParseRequestURI(config.ServerConfig.TokenUrl)
	if err != nil {
		return fmt.Errorf("invalid NotifierServer TokenUrl - %s", err)
	}
	_, err = url.ParseRequestURI(config.ClientConfig.NotifyApiRoot)
	if err != nil {
		return fmt.Errorf("invalid BaseClientWithNotifier NotifyApiRoot - %s", err)
	}
	_, err = net.ResolveTCPAddr("tcp", config.ClientConfig.LocalAddr)
	if err != nil {
		return fmt.Errorf("invalid BaseClientWithNotifier LocalAddr - %s", err)
	}
	return nil
}
