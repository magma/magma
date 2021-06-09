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
	"sort"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	merrors "magma/orc8r/lib/go/errors"

	"github.com/pkg/errors"
)

const (
	// DirectorydTableBlobstore is the table where blobstore stores directoryd's data.
	DirectorydTableBlobstore = "directoryd_blobstore"

	// DirectorydTypeHWIDToHostname is the blobstore type field for the hardware ID to hostname mapping.
	DirectorydTypeHWIDToHostname = "hwid_to_hostname"

	// DirectorydTypeSessionIDToIMSI is the blobstore type field for the session ID to IMSI mapping.
	DirectorydTypeSessionIDToIMSI = "sessionid_to_imsi"

	// DirectorydTypeSessionIDToIMSI is the blobstore type field for the session ID to IMSI mapping.
	DirectorydTypeSgwCteidToHwid = "sgwCteid_to_hwid"

	// Blobstore needs a network ID, so for network-agnostic types we use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

// NewDirectorydBlobstore returns a directoryd storage implementation
// backed by the provided blobstore factory.
func NewDirectorydBlobstore(factory blobstore.BlobStorageFactory) DirectorydStorage {
	return &directorydBlobstore{factory: factory}
}

type directorydBlobstore struct {
	factory blobstore.BlobStorageFactory
}

func (d *directorydBlobstore) GetHostnameForHWID(hwid string) (string, error) {
	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return "", errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(
		placeholderNetworkID,
		storage.TypeAndKey{Type: DirectorydTypeHWIDToHostname, Key: hwid},
	)
	if err == merrors.ErrNotFound {
		return "", err
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to get hostname")
	}

	hostname := string(blob.Value)
	return hostname, store.Commit()
}

func (d *directorydBlobstore) MapHWIDsToHostnames(hwidToHostname map[string]string) error {
	store, err := d.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blobs := convertKVToBlobs(DirectorydTypeHWIDToHostname, hwidToHostname)
	err = store.CreateOrUpdate(placeholderNetworkID, blobs)
	if err != nil {
		return errors.Wrap(err, "failed to create or update HWID to hostname mapping")
	}
	return store.Commit()
}

func (d *directorydBlobstore) GetIMSIForSessionID(networkID, sessionID string) (string, error) {
	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return "", errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blob, err := store.Get(
		networkID,
		storage.TypeAndKey{Type: DirectorydTypeSessionIDToIMSI, Key: sessionID},
	)
	if err == merrors.ErrNotFound {
		return "", err
	}
	if err != nil {
		return "", errors.Wrap(err, "failed to get IMSI")
	}

	imsi := string(blob.Value)
	return imsi, store.Commit()
}

func (d *directorydBlobstore) MapSessionIDsToIMSIs(networkID string, sessionIDToIMSI map[string]string) error {
	store, err := d.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction")
	}
	defer store.Rollback()

	blobs := convertKVToBlobs(DirectorydTypeSessionIDToIMSI, sessionIDToIMSI)
	err = store.CreateOrUpdate(networkID, blobs)
	if err != nil {
		return errors.Wrap(err, "failed to create or update session ID to IMSI mapping")
	}
	return store.Commit()
}

func (d *directorydBlobstore) GetHWIDForSgwCTeid(networkID, teid string) (string, error) {
	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to start transaction to get HwID from teid %s", teid))
	}
	defer store.Rollback()

	blob, err := store.Get(
		networkID,
		storage.TypeAndKey{Type: DirectorydTypeSgwCteidToHwid, Key: teid},
	)
	if err == merrors.ErrNotFound {
		return "", err
	}
	if err != nil {
		return "", errors.Wrap(err, fmt.Sprintf("failed to get HwID from teid %s", teid))
	}

	hwid := string(blob.Value)
	return hwid, store.Commit()
}

func (d *directorydBlobstore) MapSgwCTeidToHWID(networkID string, sgwCTeidToHwid map[string]string) error {
	store, err := d.factory.StartTransaction(&storage.TxOptions{})
	if err != nil {
		return errors.Wrap(err, "failed to start transaction to Map sgwCTeidToHwid")
	}
	defer store.Rollback()

	blobs := convertKVToBlobs(DirectorydTypeSgwCteidToHwid, sgwCTeidToHwid)
	err = store.CreateOrUpdate(networkID, blobs)
	if err != nil {
		return errors.Wrap(err, "failed to create or update sgwCTeidToHwid map")
	}
	return store.Commit()
}

// convertKVToBlobs deterministically converts a string-string map to blobstore blobs.
func convertKVToBlobs(typ string, kv map[string]string) blobstore.Blobs {
	var blobs blobstore.Blobs
	for k, v := range kv {
		blobs = append(blobs, blobstore.Blob{Type: typ, Key: k, Value: []byte(v)})
	}

	// Sort by key for deterministic behavior in tests
	sort.Slice(blobs, func(i, j int) bool { return blobs[i].Key < blobs[j].Key })

	return blobs
}
