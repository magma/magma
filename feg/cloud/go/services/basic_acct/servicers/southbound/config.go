/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicers

import (
	"flag"
	"strconv"

	"github.com/golang/glog"
)

func (cfg *Config) updateFromFlags() {
	if cfg == nil {
		return
	}
	if sv, set := getFlagValue(remoteAddrFlagName); set {
		cfg.RemoteAddr = sv
	}
	if sv, set := getFlagValue(caCertFlagName); set {
		cfg.RootCaCert = sv
	}
	if sv, set := getFlagValue(clientCertFlagName); set {
		cfg.ClientCrt = sv
	}
	if sv, set := getFlagValue(clientKeyFlagName); set {
		cfg.ClientCrtKey = sv
	}
	if bv, set := getBoolFlagValue(notlsFlagName); set {
		cfg.NoTls = bv
	}
	if bv, set := getBoolFlagValue(insecureFlagName); set {
		cfg.Insecure = bv
	}
}

// getFlagValue returns the value of the flagValue & True if it exists and was set
func getFlagValue(flagName string) (string, bool) {
	var (
		res   string
		isSet bool
	)
	if len(flagName) > 0 {
		flag.Visit(func(f *flag.Flag) {
			if f.Name == flagName {
				res = f.Value.String()
				isSet = true
				glog.V(1).Infof("Using runtime flag: %s => %s", flagName, res)
			}
		})
	}
	return res, isSet
}

// getBoolFlagValue returns the value of the bool flag & True if it exists and was set
func getBoolFlagValue(flagName string) (bool, bool) {
	if val, isSet := getFlagValue(flagName); isSet {
		if boolVal, err := strconv.ParseBool(val); err == nil {
			return boolVal, isSet
		}
	}
	return false, false
}
