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
	"fmt"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
)

const (
	// CtracedTableBlobstore is the table where blobstore stores ctraced's data.
	CtracedTableBlobstore = "ctraced_blobstore"

	// CtracedBlobType is the blobstore type field for call traces
	CtracedBlobType = "call_trace"
)

// NewCtracedBlobstore returns a ctraced storage implementation
// backed by the provided blobstore factory.
func NewCtracedBlobstore(factory blobstore.BlobStorageFactory) CtracedStorage {
	return &ctracedBlobStore{factory: factory}
}

type ctracedBlobStore struct {
	factory blobstore.BlobStorageFactory
}

// StoreCallTrace
func (c *ctracedBlobStore) StoreCallTrace(networkID string, callTraceID string, data []byte) error {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: false})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	err = store.CreateOrUpdate(
		networkID,
		blobstore.Blobs{
			{Type: CtracedBlobType, Key: callTraceID, Value: data, Version: 0},
		},
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to store call trace %s", callTraceID))
	}

	return store.Commit()
}

// GetCallTrace
func (c *ctracedBlobStore) GetCallTrace(networkID string, callTraceID string) ([]byte, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(
		networkID,
		storage.TypeAndKey{Type: CtracedBlobType, Key: callTraceID},
	)
	if err == merrors.ErrNotFound {
		return nil, err
	}
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("failed to get call trace %s", callTraceID))
	}

	return blob.Value, store.Commit()
}

// DeleteCallTrace
func (c *ctracedBlobStore) DeleteCallTrace(networkID string, callTraceID string) error {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	err = store.Delete(
		networkID,
		[]storage.TypeAndKey{
			{Type: CtracedBlobType, Key: callTraceID},
		},
	)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to delete call trace %s", callTraceID))
	}

	return store.Commit()
}
