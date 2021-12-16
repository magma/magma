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

	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/certifier/constants"
	"magma/orc8r/cloud/go/services/certifier/protos"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"
)

const (
	// CertifierTableBlobstore is the service-wide blobstore table for certifier data
	CertifierTableBlobstore = "certificate_info_blobstore"

	// Blobstore needs a network ID, but certifier is network-agnostic so we
	// will use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

type certifierBlobstore struct {
	factory blobstore.StoreFactory
}

// NewCertifierBlobstore returns an initialized instance of certifierBlobstore as CertifierStorage.
func NewCertifierBlobstore(factory blobstore.StoreFactory) CertifierStorage {
	return &certifierBlobstore{factory: factory}
}

func (c *certifierBlobstore) ListSerialNumbers() ([]string, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	serialNumbers, err := blobstore.ListKeys(store, placeholderNetworkID, constants.CertInfoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list keys")
	}

	return serialNumbers, store.Commit()
}

func (c *certifierBlobstore) GetCertInfo(serialNumber string) (*protos.CertificateInfo, error) {
	infos, err := c.GetManyCertInfo([]string{serialNumber})
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		return info, nil
	}
	return nil, merrors.ErrNotFound
}

func (c *certifierBlobstore) GetManyCertInfo(serialNumbers []string) (map[string]*protos.CertificateInfo, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	tks := storage.MakeTKs(constants.CertInfoType, serialNumbers)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get many certificate info")
	}

	ret := make(map[string]*protos.CertificateInfo)
	for _, blob := range blobs {
		info := &protos.CertificateInfo{}
		err = proto.Unmarshal(blob.Value, info)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cert info")
		}
		ret[blob.Key] = info
	}

	return ret, store.Commit()
}

func (c *certifierBlobstore) GetAllCertInfo() (map[string]*protos.CertificateInfo, error) {
	infos := map[string]*protos.CertificateInfo{}

	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	serialNumbers, err := blobstore.ListKeys(store, placeholderNetworkID, constants.CertInfoType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list keys")
	}

	if len(serialNumbers) == 0 {
		return infos, store.Commit()
	}

	tks := storage.MakeTKs(constants.CertInfoType, serialNumbers)
	blobs, err := store.GetMany(placeholderNetworkID, tks)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get many certificate info")
	}

	for _, blob := range blobs {
		info := &protos.CertificateInfo{}
		err = proto.Unmarshal(blob.Value, info)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal cert info")
		}
		infos[blob.Key] = info
	}

	return infos, store.Commit()
}

func (c *certifierBlobstore) PutCertInfo(serialNumber string, certInfo *protos.CertificateInfo) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	marshaledCertInfo, err := proto.Marshal(certInfo)
	if err != nil {
		return errors.Wrap(err, "failed to marshal cert info")
	}

	blob := blobstore.Blob{Type: constants.CertInfoType, Key: serialNumber, Value: marshaledCertInfo}
	err = store.Write(placeholderNetworkID, blobstore.Blobs{blob})
	if err != nil {
		return errors.Wrap(err, "failed to put certificate info")
	}

	return store.Commit()
}

func (c *certifierBlobstore) DeleteCertInfo(serialNumber string) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	tk := storage.TK{Type: constants.CertInfoType, Key: serialNumber}
	err = store.Delete(placeholderNetworkID, storage.TKs{tk})
	if err != nil {
		return errors.Wrap(err, "failed to delete certificate info")
	}

	return store.Commit()
}

func (c *certifierBlobstore) ListUser() ([]string, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	users, err := blobstore.ListKeys(store, placeholderNetworkID, constants.UserType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list keys")
	}

	return users, store.Commit()
}

func (c *certifierBlobstore) GetUser(username string) (*protos.User, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()
	userBlob, err := store.Get(placeholderNetworkID, storage.TK{Type: constants.UserType, Key: username})
	if err != nil {
		return nil, err
	}
	user, err := protos.UserFromBlob(userBlob)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (c *certifierBlobstore) PutUser(username string, user *protos.User) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	userBlob, err := user.UserToBlob(username)
	if err != nil {
		return err
	}

	err = store.Write(placeholderNetworkID, blobstore.Blobs{userBlob})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to update password for user %s", username))
	}

	return store.Commit()
}

func (c *certifierBlobstore) DeleteUser(username string) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	tk := storage.TK{Type: constants.UserType, Key: username}
	err = store.Delete(placeholderNetworkID, storage.TKs{tk})

	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete users: %s", err)
	}

	return store.Commit()
}

func (c *certifierBlobstore) GetPolicy(token string) (*protos.Policy, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	policyBlob, err := store.Get(placeholderNetworkID, storage.TK{Type: constants.PolicyType, Key: token})
	if err != nil {
		return nil, err
	}
	policy, err := protos.PolicyFromBlob(policyBlob)
	if err != nil {
		return nil, err
	}
	return &policy, nil

}

func (c *certifierBlobstore) PutPolicy(token string, policy *protos.Policy) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	policyBlob, err := policy.PolicyToBlob(token)
	if err != nil {
		return err
	}

	err = store.Write(placeholderNetworkID, blobstore.Blobs{policyBlob})
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("failed to create or update policy for token %s", token))
	}

	return store.Commit()
}

// ListUsers lists the usernames of all the current users
func (c *certifierBlobstore) ListUsers() ([]*protos.User, error) {
	store, err := c.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return nil, status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()
	blobs, err := blobstore.GetAllOfType(store, placeholderNetworkID, constants.UserType)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all users")
	}
	users := make([]*protos.User, len(blobs))
	for i, blob := range blobs {
		user := &protos.User{}
		err = proto.Unmarshal(blob.Value, user)
		if err != nil {
			return nil, errors.Wrap(err, "failed to unmarshal user")
		}
		users[i] = user
	}
	return users, store.Commit()
}

func (c *certifierBlobstore) DeletePolicy(token string) error {
	store, err := c.factory.StartTransaction(nil)
	if err != nil {
		return status.Errorf(codes.Unavailable, "failed to start transaction: %s", err)
	}
	defer store.Rollback()

	tk := storage.TK{Type: constants.PolicyType, Key: token}
	err = store.Delete(placeholderNetworkID, storage.TKs{tk})
	if err != nil {
		return status.Errorf(codes.Internal, "failed to delete policy: %s", err)
	}

	return store.Commit()
}
