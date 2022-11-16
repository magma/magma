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

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/merrors"
)

const (
	// DirectorydTableBlobstore is the table where blobstore stores directoryd's data.
	DirectorydTableBlobstore = "directoryd_blobstore"

	// DirectorydTypeHWIDToHostname is the blobstore type field for the hardware ID to hostname mapping.
	DirectorydTypeHWIDToHostname = "hwid_to_hostname"

	// DirectorydTypeSessionIDToIMSI is the blobstore type field for the session ID to IMSI mapping.
	DirectorydTypeSessionIDToIMSI = "sessionid_to_imsi"

	// DirectorydTypeSgwCteidToHwid is the blobstore type field for the C TEID to HWID SI mapping.
	DirectorydTypeSgwCteidToHwid = "sgwCteid_to_hwid"

	// DirectorydTypeSgwUteidToHwid is the blobstore type field for U TEID to HWID mapping.
	DirectorydTypeSgwUteidToHwid = "sgwUteid_to_hwid"

	// Blobstore needs a network ID, so for network-agnostic types we use a placeholder value.
	placeholderNetworkID = "placeholder_network"
)

// NewDirectorydBlobstore returns a directoryd storage implementation
// backed by the provided blobstore factory.
func NewDirectorydBlobstore(factory blobstore.StoreFactory) DirectorydStorage {
	return &directorydBlobstore{factory: factory}
}

type directorydBlobstore struct {
	factory blobstore.StoreFactory
}

func (d *directorydBlobstore) GetHostnameForHWID(hwid string) (string, error) {
	res, err := d.getFromStore(placeholderNetworkID, DirectorydTypeHWIDToHostname, hwid)
	printIfError(err, "Error GetHostnameForHWID: %+v", err)
	return res, err
}

func (d *directorydBlobstore) GetIMSIForSessionID(networkID, sessionID string) (string, error) {
	res, err := d.getFromStore(networkID, DirectorydTypeSessionIDToIMSI, sessionID)
	printIfError(err, "Error GetIMSIForSessionID: %+v", err)
	return res, err
}

func (d *directorydBlobstore) GetHWIDForSgwCTeid(networkID, teid string) (string, error) {
	res, err := d.getFromStore(networkID, DirectorydTypeSgwCteidToHwid, teid)
	printIfError(err, "Error GetHWIDForSgwCTeid: %+v", err)
	return res, err
}

func (d *directorydBlobstore) GetHWIDForSgwUTeid(networkID, teid string) (string, error) {
	res, err := d.getFromStore(networkID, DirectorydTypeSgwUteidToHwid, teid)
	printIfError(err, "Error GetHWIDForSgwUTeid: %+v", err)
	return res, err
}

func (d *directorydBlobstore) MapHWIDsToHostnames(hwidToHostname map[string]string) error {
	err := d.mapToStore(placeholderNetworkID, DirectorydTypeHWIDToHostname, hwidToHostname)
	printIfError(err, "Error MapHWIDsToHostnames: %+v", err)
	return err
}

func (d *directorydBlobstore) MapSessionIDsToIMSIs(networkID string, sessionIDToIMSI map[string]string) error {
	err := d.mapToStore(networkID, DirectorydTypeSessionIDToIMSI, sessionIDToIMSI)
	printIfError(err, "Error MapSessionIDsToIMSIs: %s", err)
	return err
}

func (d *directorydBlobstore) MapSgwCTeidToHWID(networkID string, sgwCTeidToHwid map[string]string) error {
	err := d.mapToStore(networkID, DirectorydTypeSgwCteidToHwid, sgwCTeidToHwid)
	printIfError(err, "Error MapSgwCTeidToHWID: %+v", err)
	return err
}

func (d *directorydBlobstore) MapSgwUTeidToHWID(networkID string, sgwUTeidToHwid map[string]string) error {
	err := d.mapToStore(networkID, DirectorydTypeSgwUteidToHwid, sgwUTeidToHwid)
	printIfError(err, "Error MapSgwUTeidToHWID: %+v", err)
	return err
}

func (d *directorydBlobstore) UnmapHWIDsToHostnames(hwids []string) error {
	err := d.unmapFromStore(placeholderNetworkID, DirectorydTypeHWIDToHostname, hwids)
	printIfError(err, "Error UnmapHWIDsToHostnames: %+v", err)
	return err
}

func (d *directorydBlobstore) UnmapSessionIDsToIMSIs(networkID string, sessionIDs []string) error {
	err := d.unmapFromStore(networkID, DirectorydTypeSessionIDToIMSI, sessionIDs)
	printIfError(err, "Error UnmapSessionIDsToIMSIs: %+v", err)
	return err
}

func (d *directorydBlobstore) UnmapSgwCTeidToHWID(networkID string, teids []string) error {
	err := d.unmapFromStore(networkID, DirectorydTypeSgwCteidToHwid, teids)
	printIfError(err, "Error UnmapSgwCTeidToHWID: %+v", err)
	return err
}

func (d *directorydBlobstore) UnmapSgwUTeidToHWID(networkID string, teids []string) error {
	err := d.unmapFromStore(networkID, DirectorydTypeSgwUteidToHwid, teids)
	printIfError(err, "Error UnmapSgwUTeidToHWID: %+v", err)
	return err
}

func (d *directorydBlobstore) getFromStore(networkID, tkType, key string) (string, error) {
	store, err := d.factory.StartTransaction(&storage.TxOptions{ReadOnly: true})
	if err != nil {
		return "", fmt.Errorf("failed to start transaction to get %s for tkKey %s: %w", key, tkType, err)
	}
	defer store.Rollback()

	blob, err := store.Get(networkID, storage.TK{Type: tkType, Key: key})
	if err == merrors.ErrNotFound {
		return "", err
	}
	if err != nil {
		return "", fmt.Errorf("failed to get %s from %s: %w", key, tkType, err)
	}
	return string(blob.Value), store.Commit()
}

func (d *directorydBlobstore) mapToStore(networkID, tkType string, keyToValueMap map[string]string) error {
	store, err := d.factory.StartTransaction(nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction mapToStore with map %+v for tkKey %s: %w", keyToValueMap, tkType, err)
	}
	defer store.Rollback()

	blobs := convertKVToBlobs(tkType, keyToValueMap)
	err = store.Write(networkID, blobs)
	if err != nil {
		return fmt.Errorf("failed to mapToStore with map %+v for tkKey %s: %w", keyToValueMap, tkType, err)
	}
	return store.Commit()
}

func (d *directorydBlobstore) unmapFromStore(networkID, tkType string, keys []string) error {
	store, err := d.factory.StartTransaction(nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction unmapFromStore with key %s for tkKey %s: %w", keys, tkType, err)
	}
	defer store.Rollback()

	err = store.Delete(networkID, storage.MakeTKs(tkType, keys))
	if err != nil {
		return fmt.Errorf("failed to unmapFromStore with keys %s for tkKey %s: %w", keys, tkType, err)
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

// printIfError prints in case of errors.
// Args:
//   - err -- error to check
//   - msg -- message to print
//   - a   -- Arguments are handled in the manner of fmt.Printf; a newline is appended if missing.
func printIfError(err error, msg string, a ...interface{}) {
	if err != nil {
		glog.Errorf(msg, a...)
	}
}
