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

package subscriberdb_cache

import (
	"magma/lte/cloud/go/lte"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

type Config struct {
	// SleepIntervalSecs is the time interval between each service worker loop.
	SleepIntervalSecs int `yaml:"sleepIntervalSecs"`
	// UpdateIntervalSecs is the target time interval to update each digest.
	UpdateIntervalSecs int `yaml:"updateIntervalSecs"`
}

func MustGetServiceConfig() Config {
	var serviceConfig Config
	_, _, err := config.GetStructuredServiceConfig(lte.ModuleName, ServiceName, &serviceConfig)
	if err != nil {
		glog.Fatalf("Failed parsing the subscriberdb_cache config file: %+v", err)
	}

	return serviceConfig
}
