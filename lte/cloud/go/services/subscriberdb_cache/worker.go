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

package subscriberdb_cache

import (
	"time"

	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func MonitorDigests(store storage.DigestLookup, config Config) {
	for {
		generatedDigestsByNetworks, err := RenewDigests(store, config)
		if err != nil {
			glog.Errorf("Error monitoring digests: %+v", err)
		}
		if len(generatedDigestsByNetworks) > 0 {
			glog.Infof("Generated digests per network: %+v", generatedDigestsByNetworks)
		}

		time.Sleep(time.Duration(config.SleepIntervalSecs) * time.Second)
	}
}

// RenewDigests continuously monitors changes in network configs and keeps a
// list of up-to-date digests per network, updated at a configurable interval,
// in the db store.
//
// Note: RenewDigests renews digests only a single time. Prefer MonitorDigests
// for continuously updating the digests.
func RenewDigests(store storage.DigestLookup, config Config) (map[string]string, error) {
	networksToRenew, networksToRemove, err := getNetworksToUpdate(store, config.UpdateIntervalSecs)
	if err != nil {
		return nil, errors.Wrapf(err, "get networks to update")
	}
	err = store.DeleteDigests(networksToRemove)
	if err != nil {
		return nil, errors.Wrapf(err, "remove invalid networks")
	}

	errs := &multierror.Error{}
	digestsByNetwork := map[string]string{}
	for _, network := range networksToRenew {
		digest, err := subscriberdb.GetDigest(network)
		if err != nil {
			multierror.Append(errors.Wrapf(err, "generate digest"))
			continue
		}

		err = store.SetDigest(network, digest)
		if err != nil {
			multierror.Append(errors.Wrapf(err, "set digest"))
			continue
		}
		digestsByNetwork[network] = digest
	}
	return digestsByNetwork, errs.ErrorOrNil()
}

// getNetworksToUpdate returns networks to renew or delete in the store.
func getNetworksToUpdate(store storage.DigestLookup, updateIntervalSecs int) ([]string, []string, error) {
	all, err := configurator.ListNetworkIDs()
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load current networks for subscriberdb cache")
	}
	tracked, err := storage.GetAllNetworks(store)
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load networks in store for subscriberdb cache")
	}

	newlyCreated, deleted := funk.DifferenceString(all, tracked)
	outdated, err := storage.GetOutdatedNetworks(store, clock.Now().Unix()-int64(updateIntervalSecs))
	if err != nil {
		return nil, nil, errors.Wrap(err, "Load outdated networks for subscriberdb cache")
	}
	trackedToRenew, _ := funk.DifferenceString(outdated, deleted)
	toRenew := append(newlyCreated, trackedToRenew...)

	return toRenew, deleted, nil
}
