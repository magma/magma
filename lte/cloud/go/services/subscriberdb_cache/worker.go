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

	lte_protos "magma/lte/cloud/go/protos"
	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/lte/cloud/go/services/subscriberdb/storage"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"

	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
	"github.com/thoas/go-funk"
)

func MonitorDigests(config Config, digestStore storage.DigestStore, perSubDigestStore *storage.PerSubDigestStore, subStore *storage.SubStore) {
	for {
		flatDigestsByNetworks, perSubDigestsByNetworks, err := RenewDigests(config, digestStore, perSubDigestStore, subStore)
		if err != nil {
			glog.Errorf("Error monitoring digests: %+v", err)
		}
		if len(flatDigestsByNetworks) > 0 {
			glog.Infof("Generated digests per network: %+v", flatDigestsByNetworks)
			glog.V(2).Infof("Generated per-sub digests per network: %+v", perSubDigestsByNetworks)
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
func RenewDigests(
	config Config,
	digestStore storage.DigestStore,
	perSubDigestStore *storage.PerSubDigestStore,
	subStore *storage.SubStore,
) (map[string]string, map[string][]*lte_protos.SubscriberDigestWithID, error) {
	networksToRenew, networksToRemove, err := getNetworksToUpdate(digestStore, config.UpdateIntervalSecs)
	if err != nil {
		return nil, nil, errors.Wrap(err, "get networks to update")
	}

	err = digestStore.DeleteDigests(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrap(err, "remove flat digests of invalid networks")
	}
	err = perSubDigestStore.DeleteDigests(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrap(err, "remove per sub digests of invalid networks")
	}
	err = subStore.DeleteSubscribersForNetworks(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrap(err, "remove sub protos of invalid networks")
	}

	errs := &multierror.Error{}
	perSubDigestsByNetwork := map[string][]*lte_protos.SubscriberDigestWithID{}
	flatDigestsByNetwork := map[string]string{}
	for _, network := range networksToRenew {
		digest, perSubDigests, err := renewDigestsForNetwork(network, digestStore, perSubDigestStore, subStore)
		if err != nil {
			multierror.Append(errs, err)
		}
		flatDigestsByNetwork[network] = digest
		perSubDigestsByNetwork[network] = perSubDigests
	}

	return flatDigestsByNetwork, perSubDigestsByNetwork, errs.ErrorOrNil()
}

// renewDigestsForNetwork updates the digest stores and subscriber proto cache for a given network.
func renewDigestsForNetwork(
	network string,
	digestStore storage.DigestStore,
	perSubDigestStore *storage.PerSubDigestStore,
	subStore *storage.SubStore,
) (string, []*lte_protos.SubscriberDigestWithID, error) {
	errs := &multierror.Error{}
	digest, prevDigest, updateDigestErr := updateDigest(network, digestStore)
	if updateDigestErr != nil {
		multierror.Append(errs, updateDigestErr)
	}

	perSubDigests, updatePerSubDigestsErr := updatePerSubDigests(network, perSubDigestStore)
	if updatePerSubDigestsErr != nil {
		multierror.Append(errs, updatePerSubDigestsErr)
	}

	// If all digest-related operations succeeded, and the generated digest is the same as
	// the previous digest, then no need to update the subscriber proto cache
	if updateDigestErr == nil && updatePerSubDigestsErr == nil && prevDigest == digest {
		return digest, perSubDigests, nil
	}

	updateSubscribersErr := updateSubscribers(network, subStore)
	if updateSubscribersErr != nil {
		multierror.Append(errs, errors.Wrapf(updateSubscribersErr, "update subscriber protos for network %+v", network))
	}
	// TODO(wangyyt1013): add logs for updated sub protos

	return digest, perSubDigests, errs.ErrorOrNil()
}

func updatePerSubDigests(network string, store *storage.PerSubDigestStore) ([]*lte_protos.SubscriberDigestWithID, error) {
	// The per-sub digests in store are updated en masse (collectively serialized into one blob per network);
	// this update takes place along with every flat digest update for consistency.
	// If an error occurs during this step, the overall last_updated_at timestamp for the network will
	// not update, and will indicate outdated-ness instead, forcing a redo in the next loop.
	perSubDigests, err := subscriberdb.GetPerSubscriberDigests(network)
	if err != nil {
		return nil, errors.Wrap(err, "get per sub dgests to update")
	}
	err = store.SetDigest(network, perSubDigests)
	if err != nil {
		return nil, errors.Wrap(err, "set per sub digest")
	}
	return perSubDigests, nil
}

// updateDigest returns the the current digest (1st return) and previous digest (2nd return) in store.
func updateDigest(network string, store storage.DigestStore) (string, string, error) {
	digest, err := subscriberdb.GetDigest(network)
	if err != nil {
		return "", "", errors.Wrap(err, "generate flat digest")
	}
	prevDigest, err := storage.GetDigest(store, network)
	if err != nil {
		return "", "", errors.Wrap(err, "get previous flat digest")
	}
	err = store.SetDigest(network, digest)
	if err != nil {
		return "", "", errors.Wrap(err, "set flat digest")
	}
	return digest, prevDigest, nil
}

func updateSubscribers(network string, store *storage.SubStore) error {
	apnsByName, err := subscriberdb.LoadApnsByName(network)
	if err != nil {
		return err
	}
	err = store.InitializeUpdate()
	if err != nil {
		return err
	}

	token := ""
	foundEmptyToken := false
	for !foundEmptyToken {
		subProtos, nextToken, err := subscriberdb.LoadSubProtosPage(0, token, network, apnsByName, lte_models.ApnResources{})
		if err != nil {
			return err
		}

		err = store.InsertMany(network, subProtos)
		if err != nil {
			return err
		}
		foundEmptyToken = nextToken == ""
		token = nextToken
	}
	return store.ApplyUpdate(network)
}

// getNetworksToUpdate returns networks to renew or delete in the store.
func getNetworksToUpdate(store storage.DigestStore, updateIntervalSecs int) ([]string, []string, error) {
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
