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
	"os"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/gateway/mconfig"

	"github.com/golang/glog"
)

const (
	S8ProxyServiceName = "s8_proxy"

	ClientAddrEnv = "CLIENT_ADDRESS"
	ServerAddrEnv = "SERVER_ADDRESS"
)

func GetS8ProxyConfig() *S8ProxyConfig {
	configPtr := &mcfgprotos.S8Config{}
	err := mconfig.GetServiceConfigs(S8ProxyServiceName, configPtr)
	if err != nil {
		glog.V(2).Infof("%s Managed Configs Load Error: %v\nUsing EnvVars", S8ProxyServiceName, err)
		return &S8ProxyConfig{
			ClientAddr: os.Getenv(ClientAddrEnv),
			ServerAddr: os.Getenv(ServerAddrEnv),
		}
	}

	glog.V(2).Infof("Loaded %s configs: %+v", S8ProxyConfig{}, *configPtr)

	return &S8ProxyConfig{
		ClientAddr: configPtr.LocalAddress,
		ServerAddr: configPtr.PgwAddress,
	}
}
