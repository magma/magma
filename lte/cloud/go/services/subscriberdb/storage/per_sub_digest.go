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

package storage

import (
	"magma/lte/cloud/go/lte"
	lte_protos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/subscriberdb/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/mproto"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/hashicorp/go-multierror"
	"github.com/pkg/errors"
)

const (
	perSubDigestBlobstoreNetworkKey = "per_sub_digest_network"
)

type PerSubDigestLookup struct {
	fact blobstore.BlobStorageFactory
}

func NewPerSubDigestLookup(fact blobstore.BlobStorageFactory) *PerSubDigestLookup {
	return &PerSubDigestLookup{fact: fact}
}

// GetDigest returns a list of per-subscriber digests of a network, ordered by their subscriber ID.
func (l *PerSubDigestLookup) GetDigest(network string) ([]*lte_protos.SubscriberDigestByID, error) {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(perSubDigestBlobstoreNetworkKey, storage.TypeAndKey{Type: lte.PerSubDigestBlobstoreType, Key: network})
	if err == merrors.ErrNotFound {
		// if per-sub digests for this network is not set yet, return empty list
		return []*lte_protos.SubscriberDigestByID{}, nil
	}
	if err != nil {
		return nil, errors.Wrapf(err, "get per-sub digests of network %+v from blobstore", network)
	}

	perSubDigests := &protos.SubscriberDigestByIDs{}
	err = mproto.UnmarshalDeterministic(blob.Value, perSubDigests)
	if err != nil {
		return nil, errors.Wrapf(err, "deserialize per-sub digests")
	}

	return perSubDigests.Digests, store.Commit()
}

// SetDigest creates or updates the per-subscriber digests of a network.
func (l *PerSubDigestLookup) SetDigest(network string, digests []*lte_protos.SubscriberDigestByID) error {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	// The sorted list of per-sub digests are serialized to be stored as a single blob.
	// This is to preserve the sorted order of the digests AND decrease the frequency
	// of writes to the store.
	blobValue := &protos.SubscriberDigestByIDs{Digests: digests}
	blobValueSerialized, err := mproto.MarshalDeterministic(blobValue)
	if err != nil {
		return errors.Wrapf(err, "serialize per-sub digests")
	}
	err = store.CreateOrUpdate(perSubDigestBlobstoreNetworkKey, blobstore.Blobs{{
		Type:  lte.PerSubDigestBlobstoreType,
		Key:   network,
		Value: blobValueSerialized,
	}})
	if err != nil {
		return errors.Wrapf(err, "set per-sub digests of network %+v in blobstore", network)
	}

	return store.Commit()
}

// DeleteDigests deletes the per-subscriber digests for the networks specified.
func (l *PerSubDigestLookup) DeleteDigests(networks []string) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	errs := &multierror.Error{}
	for _, network := range networks {
		err = store.Delete(perSubDigestBlobstoreNetworkKey, []storage.TypeAndKey{{
			Type: lte.PerSubDigestBlobstoreType,
			Key:  network,
		}})
		if err != nil {
			multierror.Append(errs, err)
		}
	}

	if errs.ErrorOrNil() != nil {
		return errors.Wrapf(errs, "Delete per-sub digests from blobstore")
	}
	return store.Commit()
}
