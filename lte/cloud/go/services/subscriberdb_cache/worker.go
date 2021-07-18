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
		return nil, nil, errors.Wrapf(err, "get networks to update")
	}

	err = digestStore.DeleteDigests(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "remove flat digests of invalid networks")
	}
	err = perSubDigestStore.DeleteDigests(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "remove per sub digests of invalid networks")
	}
	err = subStore.DeleteSubscribersForNetworks(networksToRemove)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "remove sub protos of invalid networks")
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
	digest, prevDigest, err := updateDigest(network, digestStore)
	if err != nil {
		multierror.Append(errs, err)
	}

	var perSubDigests []*lte_protos.SubscriberDigestWithID
	if errs.ErrorOrNil() == nil {
		perSubDigests, err = updatePerSubDigests(network, perSubDigestStore)
		if err != nil {
			multierror.Append(errs, err)
		}
	}

	// If all digest-related operations succeeded, and the generated digest is the same as
	// the previous digest, then no need to update the subscriber proto cache
	if errs.ErrorOrNil() == nil && prevDigest == digest {
		return digest, perSubDigests, nil
	}

	err = updateSubscribers(network, subStore)
	if err != nil {
		multierror.Append(errs, errors.Wrapf(err, "update subscriber protos for network %+v", network))
	}
	// TODO(wangyyt1013): add logs for updated sub protos
	return digest, perSubDigests, errs.ErrorOrNil()
}

func updatePerSubDigests(network string, perSubDigestStore *storage.PerSubDigestStore) ([]*lte_protos.SubscriberDigestWithID, error) {
	// The per-sub digests in store are updated en masse (collectively serialized into one blob per network);
	// this update takes place along with every flat digest update for consistency.
	// If an error occurs during this step, the overall last_updated_at timestamp for the network will
	// not update, and will indicate outdated-ness instead, forcing a redo in the next loop.
	perSubDigests, err := subscriberdb.GetPerSubscriberDigests(network)
	if err != nil {
		return nil, errors.Wrapf(err, "get per sub dgests to update")
	}
	err = perSubDigestStore.SetDigest(network, perSubDigests)
	if err != nil {
		return nil, errors.Wrapf(err, "set per sub digest")
	}
	return perSubDigests, nil
}

func updateDigest(network string, digestStore storage.DigestStore) (string, string, error) {
	digest, err := subscriberdb.GetDigest(network)
	if err != nil {
		return "", "", errors.Wrapf(err, "generate flat digest")
	}
	prevDigest, err := storage.GetDigest(digestStore, network)
	if err != nil {
		return "", "", errors.Wrapf(err, "get previous flat digest")
	}
	err = digestStore.SetDigest(network, digest)
	if err != nil {
		return "", "", errors.Wrapf(err, "set flat digest")
	}
	return digest, prevDigest, nil
}

func updateSubscribers(network string, subStore *storage.SubStore) error {
	apnsByName, err := subscriberdb.LoadApnsByName(network)
	if err != nil {
		return err
	}
	// Ensure the temporary table is empty for use
	err = subStore.InitiateUpdate()
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

		err = subStore.InsertMany(network, subProtos)
		if err != nil {
			return err
		}
		foundEmptyToken = nextToken == ""
		token = nextToken
	}
	return subStore.ApplyUpdate(network)
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
