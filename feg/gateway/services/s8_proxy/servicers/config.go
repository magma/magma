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
	"net"
	"os"
	"strconv"
	"strings"

	mcfgprotos "magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/utils"
	"magma/gateway/mconfig"

	"github.com/golang/glog"
)

const (
	S8ProxyServiceName = "s8_proxy"

	ClientAddrEnv     = "S8_CLIENT_ADDRESS"
	ServerAddrEnv     = "S8_SERVER_ADDRESS"
	ApnOperatorSuffix = "S8_APN_OPERATOR_SUFFIX"
)

func GetS8ProxyConfig() *S8ProxyConfig {
	configPtr := &mcfgprotos.S8Config{}
	conf := &S8ProxyConfig{}
	err := mconfig.GetServiceConfigs(S8ProxyServiceName, configPtr)
	if err != nil {
		glog.V(2).Infof("%s Managed Configs Load Error: %v Using EnvVars", S8ProxyServiceName, err)
		conf.ClientAddr = os.Getenv(ClientAddrEnv)
		conf.ServerAddr = ParseAddress(os.Getenv(ServerAddrEnv))
		conf.ApnOperatorSuffix = os.Getenv(ApnOperatorSuffix)
	} else {
		conf.ClientAddr = utils.GetValueOrEnv("", ClientAddrEnv, configPtr.LocalAddress)
		conf.ServerAddr = ParseAddress(utils.GetValueOrEnv("", ServerAddrEnv, configPtr.PgwAddress))
		conf.ApnOperatorSuffix = utils.GetValueOrEnv("", ApnOperatorSuffix, configPtr.ApnOperatorSuffix)
	}
	glog.V(2).Infof("Loaded configs: %+v", conf)
	return conf
}

//parseAddress will parse an ip:port address. If parse fails it will just return nil
func ParseAddress(ipAndPort string) *net.UDPAddr {
	if ipAndPort == "" {
		return nil
	}
	splitted := strings.Split(ipAndPort, ":")
	if len(splitted) != 2 {
		glog.Warningf("Malformed address. It must be formatted as IP:Port, but %s was received", ipAndPort)
		return nil
	}
	ip := splitted[0]
	if ip == "" {
		glog.Warningf("Empty IP during parsing address on config file: %s", ipAndPort)
		return nil
	}
	port, err := strconv.Atoi(splitted[1])
	if err != nil {
		glog.Warningf("Malformed PORT during parsing address on config file: %s", ipAndPort)
		return nil
	}
	ipAddr := net.ParseIP(ip)
	if ipAddr == nil {
		glog.Warningf("Malformed IP during parsing address on config file: %s", ipAndPort)
		return nil
	}
	return &net.UDPAddr{IP: ipAddr, Port: port, Zone: ""}
}
