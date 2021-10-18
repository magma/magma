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

	"github.com/pkg/errors"

	"magma/orc8r/cloud/go/services/directoryd/storage"
	"magma/orc8r/lib/go/protos"
)

type directoryLookupServicer struct {
	store storage.DirectorydStorage
}

func NewDirectoryLookupServicer(store storage.DirectorydStorage) (protos.DirectoryLookupServer, error) {
	srv := &directoryLookupServicer{store: store}
	return srv, nil
}

func (d *directoryLookupServicer) GetHostnameForHWID(
	ctx context.Context, req *protos.GetHostnameForHWIDRequest,
) (*protos.GetHostnameForHWIDResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate GetHostnameForHWIDRequest")
	}

	hostname, err := d.store.GetHostnameForHWID(req.Hwid)
	res := &protos.GetHostnameForHWIDResponse{Hostname: hostname}

	return res, err
}

func (d *directoryLookupServicer) MapHWIDsToHostnames(ctx context.Context, req *protos.MapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate MapHWIDToHostnameRequest")
	}

	err = d.store.MapHWIDsToHostnames(req.HwidToHostname)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) DeMapHWIDsToHostnames(ctx context.Context, req *protos.DeMapHWIDToHostnameRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate DeMapHWIDToHostnameRequest")
	}

	err = d.store.DeMapHWIDsToHostnames(req.Hwid)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetIMSIForSessionID(
	ctx context.Context, req *protos.GetIMSIForSessionIDRequest,
) (*protos.GetIMSIForSessionIDResponse, error) {
	err := req.Validate()

	if err != nil {
		return nil, errors.Wrap(err, "failed to validate GetIMSIForSessionIDRequest")
	}

	imsi, err := d.store.GetIMSIForSessionID(req.NetworkID, req.SessionID)
	res := &protos.GetIMSIForSessionIDResponse{Imsi: imsi}

	return res, err
}

func (d *directoryLookupServicer) MapSessionIDsToIMSIs(ctx context.Context, req *protos.MapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate MapSessionIDToIMSIRequest")
	}

	err = d.store.MapSessionIDsToIMSIs(req.NetworkID, req.SessionIDToIMSI)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) DeMapSessionIDsToIMSIs(ctx context.Context, req *protos.DeMapSessionIDToIMSIRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate DeMapSessionIDToIMSIRequest")
	}

	err = d.store.DeMapSessionIDsToIMSIs(req.NetworkID, req.SessionID)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) MapSgwCTeidToHWID(ctx context.Context, req *protos.MapSgwCTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate MapSgwCTeidToHWIDRequest")
	}

	err = d.store.MapSgwCTeidToHWID(req.NetworkID, req.TeidToHwid)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) DeMapSgwCTeidToHWID(ctx context.Context, req *protos.DeMapSgwCTeidToHWIDRequest) (*protos.Void, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate DeMapSgwCTeidToHWIDRequest")
	}

	err = d.store.DeMapSgwCTeidToHWID(req.NetworkID, req.Teid)

	return &protos.Void{}, err
}

func (d *directoryLookupServicer) GetHWIDForSgwCTeid(
	ctx context.Context, req *protos.GetHWIDForSgwCTeidRequest,
) (*protos.GetHWIDForSgwCTeidResponse, error) {
	err := req.Validate()
	if err != nil {
		return nil, errors.Wrap(err, "failed to validate GetHWIDForSgwCTeidRequest")
	}

	hwid, err := d.store.GetHWIDForSgwCTeid(req.NetworkID, req.Teid)
	res := &protos.GetHWIDForSgwCTeidResponse{Hwid: hwid}

	return res, err
}
