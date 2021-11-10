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

package servicers

import (
	"fmt"
	"net/url"

	"github.com/golang/glog"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/utils"
	"magma/gateway/mconfig"
)

const (
	N7ProxyServiceName = "n7_proxy"

	PcfApiRoot      = "PCF_API_ROOT"
	PcfTokenUrl     = "PCF_TOKEN_URL"
	PcfClientId     = "PCF_CLIENT_ID"
	PcfClientSecret = "PCF_CLIENT_SECRET"

	DefaultPcfApiRoot   = "https://localhost"
	DefaultPcfTokenUrl  = "https://localhost/token"
	DefaultClientId     = "magma_client_id"
	DefaultClientSecret = "magma_client_secret"
)

type PCFConfig struct {
	ApiRoot      string
	TokenUrl     string
	ClientId     string
	ClientSecret string
}

type N7ProxyConfig struct {
	DisableN7 bool
	Servers   []*PCFConfig
}

func (cfg *PCFConfig) Validate() error {
	if cfg == nil {
		return fmt.Errorf("no client config")
	}
	_, err := url.ParseRequestURI(cfg.ApiRoot)
	if err != nil {
		return fmt.Errorf("invalid ApiRoot")
	}
	if len(cfg.TokenUrl) == 0 {
		_, err = url.ParseRequestURI(cfg.TokenUrl)
		if err != nil {
			return fmt.Errorf("invalid TokenUrl")
		}
	}
	return nil
}

func GetN7ProxyConfig() *N7ProxyConfig {
	configPtr := &mcfgprotos.N7Config{}
	conf := &N7ProxyConfig{}

	err := mconfig.GetServiceConfigs(N7ProxyServiceName, configPtr)
	if err != nil {
		glog.V(2).Infof("Managed Configs Load Error: %v Using EnvVars", err)
		pcfCfg := &PCFConfig{
			ApiRoot:      utils.GetValueOrEnv("", PcfApiRoot, DefaultPcfApiRoot),
			TokenUrl:     utils.GetValueOrEnv("", PcfTokenUrl, DefaultPcfTokenUrl),
			ClientId:     utils.GetValueOrEnv("", PcfClientId, DefaultClientId),
			ClientSecret: utils.GetValueOrEnv("", PcfClientSecret, DefaultClientSecret),
		}
		conf.Servers = append(conf.Servers, pcfCfg)
		conf.DisableN7 = false
	} else {
		conf.DisableN7 = configPtr.DisableN7
		for _, httpCfg := range configPtr.Servers {
			pcfCfg := &PCFConfig{
				ApiRoot:      httpCfg.GetApiRoot(),
				TokenUrl:     httpCfg.GetTokenUrl(),
				ClientId:     httpCfg.GetClientId(),
				ClientSecret: httpCfg.GetClientSecret(),
			}
			err := pcfCfg.Validate()
			if err != nil {
				glog.Errorf("Managed Config Load Error, HTTP Server config failed: %v", err)
				return nil
			}
			conf.Servers = append(conf.Servers, pcfCfg)
		}
	}

	return conf
}
