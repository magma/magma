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
	mErrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
)

const networkWildcard = "*"

type Store interface {
	SetTokenInfo(oldNonce string, tokenInfo protos.TokenInfo) error
	GetTokenInfoFromLogicalID(networkID string, logicalID string) (*protos.TokenInfo, error)
	GetTokenInfoFromNonce(nonce string) (*protos.TokenInfo, error)
	IsNonceUnique(nonce string) (bool, error)
}

type blobstoreStore struct {
	factory blobstore.StoreFactory
}

func NewBlobstoreStore(factory blobstore.StoreFactory) Store {
	return &blobstoreStore{factory}
}

func (b *blobstoreStore) SetTokenInfo(oldNonce string, tokenInfo protos.TokenInfo) error {
	networkID := tokenInfo.GatewayPreregisterInfo.NetworkId
	logicalID := tokenInfo.GatewayPreregisterInfo.LogicalId

	store, err := b.factory.StartTransaction(nil)
	if err != nil {
		return err
	}
	defer store.Rollback()

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
	err = store.Write(networkWildcard, blobstore.Blobs{nonceBlob})
	if err != nil {
		return err
	}

	if oldNonce != "" {
		oldNonceTK := storage.TKs{{
			Type: string(bootstrapper.NonceTokenToInfoMap),
			Key:  oldNonce,
		}}
		err = store.Delete(networkWildcard, oldNonceTK)
		if err != nil {
			if err != mErrors.ErrNotFound{
				return err
			}
		}
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

	return &tokenInfo, store.Commit()
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
	nonceBlob, err := store.Get(networkWildcard, nonceTK)
	if err != nil {
		return nil, err
	}

	tokenInfo, err := tokenInfoFromBlob(nonceBlob)
	if err != nil {
		return nil, err
	}

	return &tokenInfo, store.Commit()
}

func (b *blobstoreStore) IsNonceUnique(nonce string) (bool, error) {
	tokenInfo, err := b.GetTokenInfoFromNonce(nonce)
	if err != nil {
		if err == mErrors.ErrNotFound {
			return true, nil
		}
		return false, err
	}
	return tokenInfo == nil, nil
}

// tokenInfoToBlob turns the input tokenInfo into a blob
// For blobType bootstrapper.LogicalIDToTokenInfo, key should be a LogicalID
// For blobType bootstrapper.NonceTokenToInfoMap, key should be a Nonce
func tokenInfoToBlob(blobType bootstrapper.DBBlobType, key string, tokenInfo protos.TokenInfo) (blobstore.Blob, error) {
	marshaledTokenInfo, err := protos.Marshal(&tokenInfo)
	if err != nil {
		return blobstore.Blob{}, errors.Wrap(err, "Error marshaling protobuf")
	}
	return blobstore.Blob{
		Type:  string(blobType),
		Key:   key,
		Value: marshaledTokenInfo,
	}, nil
}

func tokenInfoFromBlob(blob blobstore.Blob) (protos.TokenInfo, error) {
	tokenInfo := protos.TokenInfo{}
	err := protos.Unmarshal(blob.Value, &tokenInfo)
	if err != nil {
		return protos.TokenInfo{}, errors.Wrap(err, "Error unmarshaling protobuf")
	}
	return tokenInfo, nil
}
