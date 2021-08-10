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
	"encoding/binary"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
)

const (
	lastResyncTimeBlobstoreType = "gateway_last_resync_time"
)

type LastResyncTimeStore struct {
	fact blobstore.BlobStorageFactory
}

func NewLastResyncTimeStore(fact blobstore.BlobStorageFactory) *LastResyncTimeStore {
	return &LastResyncTimeStore{fact: fact}
}

func (l *LastResyncTimeStore) Get(network string, gateway string) (uint32, error) {
	store, err := l.fact.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return uint32(0), errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(network, storage.TypeAndKey{Type: lastResyncTimeBlobstoreType, Key: gateway})
	if err == merrors.ErrNotFound {
		// If this store has never been resynced, return 0 to enforce first resync
		return uint32(0), nil
	}
	if err != nil {
		return uint32(0), errors.Wrapf(err, "get last resync time of network %+v, gateway %+v from blobstore", network, gateway)
	}

	lastResyncTime := binary.LittleEndian.Uint32(blob.Value)
	return lastResyncTime, store.Commit()
}

func (l *LastResyncTimeStore) Set(network string, gateway string, unixTime uint32) error {
	store, err := l.fact.StartTransaction(nil)
	if err != nil {
		return errors.Wrapf(err, "error starting transaction")
	}
	defer store.Rollback()

	lastResyncTimeBytes := make([]byte, 8)
	binary.LittleEndian.PutUint32(lastResyncTimeBytes, unixTime)
	err = store.CreateOrUpdate(network, blobstore.Blobs{{
		Type:  lastResyncTimeBlobstoreType,
		Key:   gateway,
		Value: lastResyncTimeBytes,
	}})
	if err != nil {
		return errors.Wrapf(err, "set last resync time of network %+v, gateway %+v in blobstore", network, gateway)
	}

	return store.Commit()
}
