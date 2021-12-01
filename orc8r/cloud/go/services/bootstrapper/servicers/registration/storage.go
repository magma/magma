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

package registration

import (
	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/bootstrapper"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
)

const placeholderNetworkID = "placeholder_network"

type Store interface {
	SetTokenInfo(tokenInfo *protos.TokenInfo) error
	GetTokenInfoFromLogicalID(networkID string, logicalID string) (*protos.TokenInfo, error)
	GetTokenInfoFromNonce(nonce string) (*protos.TokenInfo, error)
}

type blobstoreStore struct {
	factory blobstore.StoreFactory
}

func NewBlobstoreStore(factory blobstore.StoreFactory) Store {
	return &blobstoreStore{factory: factory}
}

func (b *blobstoreStore) SetTokenInfo(tokenInfo *protos.TokenInfo) error {
	networkID := tokenInfo.GatewayDeviceInfo.NetworkId
	logicalID := tokenInfo.GatewayDeviceInfo.LogicalId

	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	defer store.Rollback()

	isUnique, err := b.isNonceUnique(store, tokenInfo.Nonce)
	if err != nil {
		return err
	}
	if !isUnique {
		return errors.Errorf("token is not unique")
	}

	// Write to 2 blobstores so that we can key for tokenInfo with both the logicalID and the nonce

	logicalIDBlob, err := tokenInfoToBlob(bootstrapper.LogicalIDToTokenInfo, logicalID, tokenInfo)
	if err != nil {
		return err
	}
	err = store.Write(networkID, blobstore.Blobs{logicalIDBlob})
	if err != nil {
		return err
	}

	nonceBlob, err := tokenInfoToBlob(bootstrapper.NonceTokenToInfoMap, tokenInfo.Nonce, tokenInfo)
	if err != nil {
		return err
	}
	err = store.Write(placeholderNetworkID, blobstore.Blobs{nonceBlob})
	if err != nil {
		return err
	}

	return store.Commit()
}

func (b *blobstoreStore) GetTokenInfoFromLogicalID(networkID string, logicalID string) (*protos.TokenInfo, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	logicalIDTK := storage.TK{
		Type: string(bootstrapper.LogicalIDToTokenInfo),
		Key:  logicalID,
	}
	logicalIDBlob, err := store.Get(networkID, logicalIDTK)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := tokenInfoFromBlob(logicalIDBlob)
	if err != nil {
		return nil, err
	}

	return tokenInfo, store.Commit()
}

func (b *blobstoreStore) GetTokenInfoFromNonce(nonce string) (*protos.TokenInfo, error) {
	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	nonceTK := storage.TK{
		Type: string(bootstrapper.NonceTokenToInfoMap),
		Key:  nonce,
	}
	nonceBlob, err := store.Get(placeholderNetworkID, nonceTK)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := tokenInfoFromBlob(nonceBlob)
	if err != nil {
		return nil, err
	}

	return tokenInfo, store.Commit()
}

// isNonceUnique must be run within a transaction
func (b *blobstoreStore) isNonceUnique(store blobstore.Store, nonce string) (bool, error) {
	nonceTK := storage.TK{
		Type: string(bootstrapper.NonceTokenToInfoMap),
		Key:  nonce,
	}
	_, err := store.Get(placeholderNetworkID, nonceTK)
	if err == merrors.ErrNotFound {
		return true, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}

// tokenInfoToBlob turns the input tokenInfo into a blob
// For blobType bootstrapper.LogicalIDToTokenInfo, key should be a LogicalID
// For blobType bootstrapper.NonceTokenToInfoMap, key should be a Nonce
func tokenInfoToBlob(blobType bootstrapper.BlobType, key string, tokenInfo *protos.TokenInfo) (blobstore.Blob, error) {
	marshaledTokenInfo, err := protos.Marshal(tokenInfo)
	if err != nil {
		return blobstore.Blob{}, errors.Wrap(err, "Error marshaling protobuf")
	}
	blob := blobstore.Blob{
		Type:  string(blobType),
		Key:   key,
		Value: marshaledTokenInfo,
	}
	return blob, nil
}

func tokenInfoFromBlob(blob blobstore.Blob) (*protos.TokenInfo, error) {
	tokenInfo := protos.TokenInfo{}
	err := protos.Unmarshal(blob.Value, &tokenInfo)
	if err != nil {
		return &protos.TokenInfo{}, errors.Wrap(err, "Error unmarshaling protobuf")
	}
	return &tokenInfo, nil
}
