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

package subscriberdb

import (
	"magma/lte/cloud/go/lte"
	"magma/orc8r/lib/go/service/config"

	"github.com/golang/glog"
)

type Config struct {
	// FlatDigestEnabled is a feature flag for the flat digest functionality.
	FlatDigestEnabled bool `yaml:"flatDigestEnabled"`
	// ChangesetSizeTheshold specifies the max size of the cloud-agw changeset
	// past which a resync signal will be sent back to the agw.
	ChangesetSizeTheshold int `yaml:"changesetSizeTheshold"`
	SyncInterval      uint32 `yaml:"syncInterval"`
	// DefaultSyncInterval is the the default interval in seconds between
	// gateway requests to sync its subscriberdb with the cloud.
	DefaultSyncInterval uint32 `yaml:"defaultSyncInterval"`
}

func MustGetServiceConfig() Config {
	var serviceConfig Config
	_, _, err := config.GetStructuredServiceConfig(lte.ModuleName, ServiceName, &serviceConfig)
	if err != nil {
		glog.Fatalf("Failed parsing the subscriberdb config file: %+v", err)
	}

	return serviceConfig
}
