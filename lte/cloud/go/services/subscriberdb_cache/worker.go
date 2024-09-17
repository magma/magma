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
	"context"
	"fmt"
	"time"

	"github.com/golang/glog"
	"github.com/hashicorp/go-multierror"
	"github.com/thoas/go-funk"

	lte_models "magma/lte/cloud/go/services/lte/obsidian/models"
	"magma/lte/cloud/go/services/subscriberdb"
	"magma/orc8r/cloud/go/clock"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/syncstore"
	"magma/orc8r/lib/go/protos"
)

func MonitorDigests(config Config, store syncstore.SyncStore) {
	for {
		rootDigests, leafDigests, err := RenewDigests(config, store)
		if err != nil {
			glog.Errorf("Error monitoring digests: %+v", err)
		}
		if len(rootDigests) > 0 {
			glog.Infof("Generated root digests per network: %+v", rootDigests)
			glog.V(2).Infof("Generated leaf digests per network: %+v", leafDigests)
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
func RenewDigests(config Config, store syncstore.SyncStore) (map[string]string, map[string][]*protos.LeafDigest, error) {
	tracked, err := configurator.ListNetworkIDs(context.Background())
	if err != nil {
		return nil, nil, fmt.Errorf("Load current networks for subscriberdb cache: %w", err)
	}
	// Garbage collection needs to happen here to ensure that untracked networks are removed from store
	// before the next step
	store.CollectGarbage(tracked)
	toUpdate, err := getNetworksToUpdate(store, tracked, config.UpdateIntervalSecs)
	if err != nil {
		return nil, nil, fmt.Errorf("get networks to update: %w", err)
	}

	errs := &multierror.Error{}
	rootDigestsByNetwork := map[string]string{}
	leafDigestsByNetwork := map[string][]*protos.LeafDigest{}
	for _, network := range toUpdate {
		rootDigest, leaveDigests, err := renewDigestsForNetwork(network, store)
		errs = multierror.Append(errs, err)
		rootDigestsByNetwork[network] = rootDigest
		leafDigestsByNetwork[network] = leaveDigests
	}

	return rootDigestsByNetwork, leafDigestsByNetwork, errs.ErrorOrNil()
}

// renewDigestsForNetwork updates the digest stores and subscriber proto cache for a given network.
func renewDigestsForNetwork(
	network string,
	store syncstore.SyncStore,
) (string, []*protos.LeafDigest, error) {
	prevDigestTree, err := syncstore.GetDigestTree(store, network)
	if err != nil {
		return "", nil, fmt.Errorf("get previous root digest for network %+v: %w", network, err)
	}
	rootDigest, leafDigests, err := updateDigestTree(network, store)
	// If the digest-related operations succeeded, and the generated root digest is the same as
	// the previous root digest, then no need to update the subscriber cache
	if err == nil && prevDigestTree.RootDigest.GetMd5Base64Digest() == rootDigest {
		return rootDigest, leafDigests, nil
	}

	err = updateSubscribers(network, store)
	if err != nil {
		return "", nil, fmt.Errorf("update subscribers cache for network %+v: %w", network, err)
	}
	// TODO(wangyyt1013): add logs for updated sub protos

	return rootDigest, leafDigests, nil
}

func updateDigestTree(network string, store syncstore.SyncStore) (string, []*protos.LeafDigest, error) {
	rootDigest, err := subscriberdb.GetDigest(network)
	if err != nil {
		return "", nil, fmt.Errorf("generate root digest: %w", err)
	}
	// The leaf digests in store are updated en masse (collectively serialized into one blob per network);
	// this update takes place along with every root digest update for consistency.
	// If an error occurs during this step, the overall last_updated_at timestamp for the network will
	// not update, and will indicate outdated-ness instead, forcing a redo in the next loop.
	leafDigests, err := subscriberdb.GetPerSubscriberDigests(network)
	if err != nil {
		return "", nil, fmt.Errorf("get per-subscriber digests to update: %w", err)
	}
	digestTree := &protos.DigestTree{
		RootDigest:  &protos.Digest{Md5Base64Digest: rootDigest},
		LeafDigests: leafDigests,
	}
	err = store.SetDigest(network, digestTree)
	if err != nil {
		return "", nil, fmt.Errorf("set digest for network %+v: %w", network, err)
	}
	return rootDigest, leafDigests, nil
}

func updateSubscribers(network string, store syncstore.SyncStore) error {
	apnsByName, err := subscriberdb.LoadApnsByName(network)
	if err != nil {
		return err
	}
	writer, err := store.UpdateCache(network)
	if err != nil {
		return fmt.Errorf("get new cache writer for network %+v: %w", network, err)
	}

	token := ""
	foundEmptyToken := false
	for !foundEmptyToken {
		subProtos, nextToken, err := subscriberdb.LoadSubProtosPage(0, token, network, apnsByName, lte_models.ApnResources{})
		if err != nil {
			return err
		}
		subProtosSerialized, err := subscriberdb.SerializeSubscribers(subProtos)
		if err != nil {
			return err
		}
		err = writer.InsertMany(subProtosSerialized)
		if err != nil {
			return err
		}
		foundEmptyToken = nextToken == ""
		token = nextToken
	}
	return writer.Apply()
}

// getNetworksToUpdate returns the networks that need to be updated, given that
// all networks in store are tracked (but not all tracked networks are in store).
func getNetworksToUpdate(store syncstore.SyncStore, tracked []string, updateIntervalSecs int) ([]string, error) {
	storedDigests, err := store.GetDigests([]string{}, clock.Now().Unix(), false)
	if err != nil {
		return nil, fmt.Errorf("Load digests in store for subscriberdb cache: %w", err)
	}
	outdatedDigests, err := store.GetDigests([]string{}, clock.Now().Unix()-int64(updateIntervalSecs), false)
	if err != nil {
		return nil, fmt.Errorf("Load outdated digests in store for subscriberdb cache: %w", err)
	}

	newlyCreated, _ := funk.DifferenceString(tracked, storedDigests.Networks())
	return append(newlyCreated, outdatedDigests.Networks()...), nil
}
