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

package directoryd

import (
	"context"
	"fmt"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const ServiceName = "DIRECTORYD"

//-------------------------------
// Directoryd service client APIs
//-------------------------------

// getDirectorydClient returns an RPC connection to the directoryd service.
func getDirectorydClient() (protos.DirectoryLookupClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDirectoryLookupClient(conn), err
}

// GetHostnameForHWID returns the hostname mapped to by hardware ID.
// Derived state, stored in directoryd service.
func GetHostnameForHWID(hwid string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", errors.Wrap(err, "failed to get directoryd client")
	}

	res, err := client.GetHostnameForHWID(context.Background(), &protos.GetHostnameForHWIDRequest{Hwid: hwid})
	if err != nil {
		return "", fmt.Errorf("failed to get hostname for hwid %s: %s", hwid, err)
	}

	return res.Hostname, nil
}

// MapHWIDToHostname maps a single hwid to a hostname.
// Derived state, stored in directoryd service.
func MapHWIDToHostname(hwid, hostname string) error {
	return MapHWIDsToHostnames(map[string]string{hwid: hostname})
}

// MapHWIDsToHostnames maps {hwid -> hostname}.
// Derived state, stored in directoryd service.
func MapHWIDsToHostnames(hwidToHostname map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directoryd client")
	}

	_, err = client.MapHWIDsToHostnames(context.Background(), &protos.MapHWIDToHostnameRequest{HwidToHostname: hwidToHostname})
	if err != nil {
		return fmt.Errorf("failed to map hwids to hostnames %v: %s", hwidToHostname, err)
	}

	return nil
}

// GetIMSIForSessionID returns the IMSI mapped to by session ID.
// Derived state, stored in directoryd service.
// NOTE: this mapping is provided on a best-effort basis, meaning
//	- a {session ID -> IMSI} mapping may be missing even though the IMSI has a session ID record
//	- a {session ID -> IMSI} mapping may be stale
func GetIMSIForSessionID(networkID, sessionID string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", errors.Wrap(err, "failed to get directoryd client")
	}

	res, err := client.GetIMSIForSessionID(context.Background(), &protos.GetIMSIForSessionIDRequest{
		NetworkID: networkID,
		SessionID: sessionID,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get IMSI for session ID %s under network ID %s: %s", sessionID, networkID, err)
	}

	return res.Imsi, nil
}

// MapSessionIDsToIMSIs maps {session ID -> IMSI}.
// Derived state, stored in directoryd service.
func MapSessionIDsToIMSIs(networkID string, sessionIDToIMSI map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directoryd client")
	}

	_, err = client.MapSessionIDsToIMSIs(context.Background(), &protos.MapSessionIDToIMSIRequest{
		NetworkID:       networkID,
		SessionIDToIMSI: sessionIDToIMSI,
	})
	if err != nil {
		return fmt.Errorf("failed to map session IDs to IMSIs %v under network ID %s: %s", sessionIDToIMSI, networkID, err)
	}

	return nil
}

// GetHWIDForSgwCTeid returns the HwID mapped to by teid
// Derived state, stored in directoryd service.
// NOTE: this mapping is provided on a best-effort basis, meaning
//	- a {teid -> HwId} mapping may be missing even though the IMSI has a session ID record
//	- a {teid -> HwId} mapping may be stale
func GetHWIDForSgwCTeid(networkID, teid string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", errors.Wrap(err, "failed to get directoryd client")
	}

	res, err := client.GetHWIDForSgwCTeid(context.Background(), &protos.GetHWIDForSgwCTeidRequest{
		NetworkID: networkID,
		Teid:      teid,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get HwID for sgw c teid %s under network ID %s: %s", teid, networkID, err)
	}

	return res.GetHwid(), nil
}

// MapSgwCTeidToHWID maps {Teid -> HwId}
func MapSgwCTeidToHWID(networkID string, teidToHWID map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return errors.Wrap(err, "failed to get directoryd client")
	}

	_, err = client.MapSgwCTeidToHWID(context.Background(), &protos.MapSgwCTeidToHWIDRequest{
		NetworkID:  networkID,
		TeidToHwid: teidToHWID,
	})
	if err != nil {
		return fmt.Errorf("failed to map sgw c teid to HwId %v under network ID %s: %s", teidToHWID, networkID, err)
	}

	return nil
}

//--------------------------
// State service client APIs
//--------------------------

// GetHWIDForIMSI returns the HWID mapped to by the IMSI.
// Primary state, stored in state service.
func GetHWIDForIMSI(networkID, imsi string) (string, error) {
	st, err := state.GetState(networkID, orc8r.DirectoryRecordType, imsi, serdes.State)
	if err != nil {
		return "", err
	}
	record, ok := st.ReportedState.(*types.DirectoryRecord)
	if !ok || len(record.LocationHistory) == 0 {
		return "", fmt.Errorf("failed to convert reported state to DirectoryRecord for device id: %s", st.ReporterID)
	}
	return record.LocationHistory[0], nil
}

// GetSessionIDForIMSI returns the session ID mapped to by the IMSI.
// Primary state, stored in state service.
func GetSessionIDForIMSI(networkID, imsi string) (string, error) {
	st, err := state.GetState(networkID, orc8r.DirectoryRecordType, imsi, serdes.State)
	if err != nil {
		return "", err
	}

	record, ok := st.ReportedState.(*types.DirectoryRecord)
	if !ok {
		return "", fmt.Errorf("failed to convert reported state to DirectoryRecord for device id: %s", st.ReporterID)
	}

	return record.GetSessionID()
}
