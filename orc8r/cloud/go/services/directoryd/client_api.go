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

	"github.com/golang/glog"

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/directoryd/protos"
	"magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state"
	"magma/orc8r/lib/go/merrors"
	lib_protos "magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const ServiceName = "directoryd"

//-------------------------------
// Directoryd service client APIs
//-------------------------------

// getDirectorydClient returns an RPC connection to the directoryd service.
func getDirectorydClient() (protos.DirectoryLookupClient, error) {
	conn, err := registry.GetConnection(ServiceName, lib_protos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewDirectoryLookupClient(conn), err
}

// GetHostnameForHWID returns the hostname mapped to by hardware ID.
// Derived state, stored in directoryd service.
func GetHostnameForHWID(ctx context.Context, hwid string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}

	res, err := client.GetHostnameForHWID(ctx, &protos.GetHostnameForHWIDRequest{Hwid: hwid})
	if err != nil {
		return "", fmt.Errorf("failed to get hostname for hwid %s: %s", hwid, err)
	}

	return res.Hostname, nil
}

// MapHWIDToHostname maps a single hwid to a hostname.
// Derived state, stored in directoryd service.
func MapHWIDToHostname(ctx context.Context, hwid, hostname string) error {
	return MapHWIDsToHostnames(ctx, map[string]string{hwid: hostname})
}

// MapHWIDsToHostnames maps {hwid -> hostname}.
// Derived state, stored in directoryd service.
func MapHWIDsToHostnames(ctx context.Context, hwidToHostname map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.MapHWIDsToHostnames(ctx, &protos.MapHWIDToHostnameRequest{HwidToHostname: hwidToHostname})
	if err != nil {
		return fmt.Errorf("failed to map hwids to hostnames %v: %s", hwidToHostname, err)
	}

	return nil
}

// UnmapHWIDsToHostnames removes the {hwid -> hostname} map
// Derived state, stored in directoryd service.
func UnmapHWIDsToHostnames(ctx context.Context, hwids []string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.UnmapHWIDsToHostnames(ctx, &protos.UnmapHWIDToHostnameRequest{Hwids: hwids})
	if err != nil {
		return fmt.Errorf("failed to ummap hwids to hostnames %v: %s", hwids, err)
	}
	return nil
}

// GetIMSIForSessionID returns the IMSI mapped to by session ID.
// Derived state, stored in directoryd service.
// NOTE: this mapping is provided on a best-effort basis, meaning
//   - a {session ID -> IMSI} mapping may be missing even though the IMSI has a session ID record
//   - a {session ID -> IMSI} mapping may be stale
func GetIMSIForSessionID(ctx context.Context, networkID, sessionID string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}

	res, err := client.GetIMSIForSessionID(ctx, &protos.GetIMSIForSessionIDRequest{
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
func MapSessionIDsToIMSIs(ctx context.Context, networkID string, sessionIDToIMSI map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.MapSessionIDsToIMSIs(ctx, &protos.MapSessionIDToIMSIRequest{
		NetworkID:       networkID,
		SessionIDToIMSI: sessionIDToIMSI,
	})
	if err != nil {
		return fmt.Errorf("failed to map session IDs to IMSIs %v under network ID %s: %s", sessionIDToIMSI, networkID, err)
	}

	return nil
}

// MapSessionIDsToIMSIs removes {session ID -> IMSI} mapping
// Derived state, stored in directoryd service.
func UnmapSessionIDsToIMSIs(ctx context.Context, networkID string, sessionIDs []string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.UnmapSessionIDsToIMSIs(ctx, &protos.UnmapSessionIDToIMSIRequest{
		NetworkID:  networkID,
		SessionIDs: sessionIDs,
	})
	if err != nil {
		return fmt.Errorf("failed to unmap session IDs %v under network ID %s: %s", sessionIDs, networkID, err)
	}

	return nil
}

// GetHWIDForSgwCTeid returns the HwID mapped to by Control teid
// Derived state, stored in directoryd service.
// NOTE: this mapping is provided on a best-effort basis, meaning
//   - a {teid -> HwId} mapping may be missing even though the IMSI has a session ID record
//   - a {teid -> HwId} mapping may be stale
func GetHWIDForSgwCTeid(ctx context.Context, networkID, teid string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}

	res, err := client.GetHWIDForSgwCTeid(ctx, &protos.GetHWIDForSgwCTeidRequest{
		NetworkID: networkID,
		Teid:      teid,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get HwID for sgw c teid %s under network ID %s: %s", teid, networkID, err)
	}

	return res.GetHwid(), nil
}

// MapSgwCTeidToHWID maps {Teid -> HwId}
func MapSgwCTeidToHWID(ctx context.Context, networkID string, teidToHWID map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.MapSgwCTeidToHWID(ctx, &protos.MapSgwCTeidToHWIDRequest{
		NetworkID:  networkID,
		TeidToHwid: teidToHWID,
	})
	if err != nil {
		return fmt.Errorf("failed to map sgw c teid to HwId %v under network ID %s: %s", teidToHWID, networkID, err)
	}

	return nil
}

// UnmapSgwCTeidToHWID removes {Teid -> HwId} mapping
func UnmapSgwCTeidToHWID(ctx context.Context, networkID string, teids []string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.UnmapSgwCTeidToHWID(ctx, &protos.UnmapSgwCTeidToHWIDRequest{
		NetworkID: networkID,
		Teids:     teids,
	})
	if err != nil {
		return fmt.Errorf("failed to ummap sgw c teid %v under network ID %s: %s", teids, networkID, err)
	}

	return nil
}

// GetNewSgwCTeid get an available teid
func GetNewSgwCTeid(ctx context.Context, networkID string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}
	res, err := client.GetNewSgwCTeid(ctx, &protos.GetNewSgwCTeidRequest{NetworkID: networkID})
	if err != nil {
		return "", fmt.Errorf("failed to get new sgw c teid under network ID %s: %s", networkID, err)
	}
	return res.Teid, nil
}

// GetHWIDForSgwUTeid returns the HwID mapped to by User teid
// Derived state, stored in directoryd service.
// NOTE: this mapping is provided on a best-effort basis, meaning
//   - a {teid -> HwId} mapping may be missing even though the IMSI has a session ID record
//   - a {teid -> HwId} mapping may be stale
func GetHWIDForSgwUTeid(ctx context.Context, networkID, teid string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}

	res, err := client.GetHWIDForSgwUTeid(ctx, &protos.GetHWIDForSgwUTeidRequest{
		NetworkID: networkID,
		Teid:      teid,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get HwID for sgw c teid %s under network ID %s: %s", teid, networkID, err)
	}

	return res.GetHwid(), nil
}

// MapSgwUTeidToHWID maps {Teid -> HwId}
func MapSgwUTeidToHWID(ctx context.Context, networkID string, teidToHWID map[string]string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.MapSgwUTeidToHWID(ctx, &protos.MapSgwUTeidToHWIDRequest{
		NetworkID:  networkID,
		TeidToHwid: teidToHWID,
	})
	if err != nil {
		return fmt.Errorf("failed to map sgw c teid to HwId %v under network ID %s: %s", teidToHWID, networkID, err)
	}

	return nil
}

// UnmapSgwUTeidToHWID removes {Teid -> HwId} mapping
func UnmapSgwUTeidToHWID(ctx context.Context, networkID string, teids []string) error {
	client, err := getDirectorydClient()
	if err != nil {
		return fmt.Errorf("failed to get directoryd client: %w", err)
	}

	_, err = client.UnmapSgwUTeidToHWID(ctx, &protos.UnmapSgwUTeidToHWIDRequest{
		NetworkID: networkID,
		Teids:     teids,
	})
	if err != nil {
		return fmt.Errorf("failed to ummap sgw c teid %v under network ID %s: %s", teids, networkID, err)
	}

	return nil
}

// GetNewSgwUTeid get an available teid
func GetNewSgwUTeid(ctx context.Context, networkID string) (string, error) {
	client, err := getDirectorydClient()
	if err != nil {
		return "", fmt.Errorf("failed to get directoryd client: %w", err)
	}
	res, err := client.GetNewSgwUTeid(ctx, &protos.GetNewSgwUTeidRequest{NetworkID: networkID})
	if err != nil {
		return "", fmt.Errorf("failed to get new sgw c teid under network ID %s: %s", networkID, err)
	}
	return res.Teid, nil
}

//--------------------------
// State service client APIs
//--------------------------

// GetHWIDForIMSI returns the HWID mapped to by the IMSI.
// Primary state, stored in state service.
func GetHWIDForIMSI(ctx context.Context, networkID, imsi string) (string, error) {
	st, err := state.GetState(ctx, networkID, orc8r.DirectoryRecordType, imsi, serdes.State)
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
func GetSessionIDForIMSI(ctx context.Context, networkID, imsi string) (string, error) {
	st, err := state.GetState(ctx, networkID, orc8r.DirectoryRecordType, imsi, serdes.State)
	if err != nil {
		return "", err
	}

	record, ok := st.ReportedState.(*types.DirectoryRecord)
	if !ok {
		return "", fmt.Errorf("failed to convert reported state to DirectoryRecord for device id: %s", st.ReporterID)
	}

	return record.GetSessionID()
}
