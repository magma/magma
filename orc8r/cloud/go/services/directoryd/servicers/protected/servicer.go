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

package servicers

import (
	"context"
	"fmt"

	"github.com/golang/glog"

	directoryd_protos "magma/orc8r/cloud/go/services/directoryd/protos"
	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/lib/go/protos"
)

type directoryLookupServicer struct {
	store       storage.DirectorydStorage
	sgwCteidGen *IdGenerator
	sgwUteidGen *IdGenerator
}

func NewDirectoryLookupServicer(store storage.DirectorydStorage) (directoryd_protos.DirectoryLookupServer, error) {
	srv := &directoryLookupServicer{
		store:       store,
		sgwCteidGen: NewIdGenerator(),
		sgwUteidGen: NewIdGenerator(),
	}
	return srv, nil
}

func (d *directoryLookupServicer) GetHostnameForHWID(
	ctx context.Context, req *directoryd_protos.GetHostnameForHWIDRequest,
) (*directoryd_protos.GetHostnameForHWIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate GetHostnameForHWIDRequest: %w", err)
	}

	hostname, err := d.store.GetHostnameForHWID(req.Hwid)
	res := &directoryd_protos.GetHostnameForHWIDResponse{Hostname: hostname}

	return res, err
}

func (d *directoryLookupServicer) MapHWIDsToHostnames(ctx context.Context, req *directoryd_protos.MapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate MapHWIDToHostnameRequest: %w", err)
	}

	err = d.store.MapHWIDsToHostnames(req.HwidToHostname)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) UnmapHWIDsToHostnames(ctx context.Context, req *directoryd_protos.UnmapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate UnmapHWIDToHostnameRequest: %w", err)
	}

	err = d.store.UnmapHWIDsToHostnames(req.Hwids)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetIMSIForSessionID(
	ctx context.Context, req *directoryd_protos.GetIMSIForSessionIDRequest,
) (*directoryd_protos.GetIMSIForSessionIDResponse, error) {
	err := req.Validate()

	if err != nil {
		return nil, fmt.Errorf("failed to validate GetIMSIForSessionIDRequest: %w", err)
	}

	imsi, err := d.store.GetIMSIForSessionID(req.NetworkID, req.SessionID)
	res := &directoryd_protos.GetIMSIForSessionIDResponse{Imsi: imsi}

	return res, err
}

func (d *directoryLookupServicer) MapSessionIDsToIMSIs(ctx context.Context, req *directoryd_protos.MapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate MapSessionIDToIMSIRequest: %w", err)
	}

	err = d.store.MapSessionIDsToIMSIs(req.NetworkID, req.SessionIDToIMSI)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) UnmapSessionIDsToIMSIs(ctx context.Context, req *directoryd_protos.UnmapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate UnmapSessionIDToIMSIRequest: %w", err)
	}

	err = d.store.UnmapSessionIDsToIMSIs(req.NetworkID, req.SessionIDs)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) MapSgwCTeidToHWID(ctx context.Context, req *directoryd_protos.MapSgwCTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate MapSgwCTeidToHWIDRequest: %w", err)
	}

	err = d.store.MapSgwCTeidToHWID(req.NetworkID, req.TeidToHwid)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) UnmapSgwCTeidToHWID(ctx context.Context, req *directoryd_protos.UnmapSgwCTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate UnmapSgwCTeidToHWIDRequest: %w", err)
	}

	err = d.store.UnmapSgwCTeidToHWID(req.NetworkID, req.Teids)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetHWIDForSgwCTeid(
	ctx context.Context, req *directoryd_protos.GetHWIDForSgwCTeidRequest,
) (*directoryd_protos.GetHWIDForSgwCTeidResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate GetHWIDForSgwCTeidRequest: %w", err)
	}

	hwid, err := d.store.GetHWIDForSgwCTeid(req.NetworkID, req.Teid)
	res := &directoryd_protos.GetHWIDForSgwCTeidResponse{Hwid: hwid}

	return res, err
}

func (d *directoryLookupServicer) GetNewSgwCTeid(ctx context.Context, req *directoryd_protos.GetNewSgwCTeidRequest) (*directoryd_protos.GetNewSgwCTeidResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate GetNewSgwCTeid: %w", err)
	}

	sgwCTeid, err := d.sgwCteidGen.GetUniqueId(req.NetworkID, d.store.GetHWIDForSgwCTeid)
	if err != nil {
		err = fmt.Errorf("GetNewSgwCTeid could not get unique TEID: %s", err)
		glog.Error(err)
		return nil, err
	}
	return &directoryd_protos.GetNewSgwCTeidResponse{Teid: fmt.Sprint(sgwCTeid)}, nil
}

func (d *directoryLookupServicer) MapSgwUTeidToHWID(ctx context.Context, req *directoryd_protos.MapSgwUTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate MapSgwUTeidToHWIDRequest: %w", err)
	}

	err = d.store.MapSgwUTeidToHWID(req.NetworkID, req.TeidToHwid)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) UnmapSgwUTeidToHWID(ctx context.Context, req *directoryd_protos.UnmapSgwUTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate UnmapSgwUTeidToHWIDRequest: %w", err)
	}

	err = d.store.UnmapSgwUTeidToHWID(req.NetworkID, req.Teids)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetHWIDForSgwUTeid(
	ctx context.Context, req *directoryd_protos.GetHWIDForSgwUTeidRequest,
) (*directoryd_protos.GetHWIDForSgwUTeidResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate GetHWIDForSgwUTeidRequest: %w", err)
	}

	hwid, err := d.store.GetHWIDForSgwUTeid(req.NetworkID, req.Teid)
	res := &directoryd_protos.GetHWIDForSgwUTeidResponse{Hwid: hwid}

	return res, err
}

func (d *directoryLookupServicer) GetNewSgwUTeid(ctx context.Context, req *directoryd_protos.GetNewSgwUTeidRequest) (*directoryd_protos.GetNewSgwUTeidResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, fmt.Errorf("failed to validate GetNewSgwUTeid: %w", err)
	}

	SgwUTeid, err := d.sgwUteidGen.GetUniqueId(req.NetworkID, d.store.GetHWIDForSgwUTeid)
	if err != nil {
		err = fmt.Errorf("GetNewSgwUTeid could not get unique TEID: %s", err)
		glog.Error(err)
		return nil, err
	}
	return &directoryd_protos.GetNewSgwUTeidResponse{Teid: fmt.Sprint(SgwUTeid)}, nil
}
