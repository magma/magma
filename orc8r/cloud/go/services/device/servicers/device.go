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

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/services/device/protos"
	commonProtos "magma/orc8r/lib/go/protos"
)

type deviceServicer struct {
	factory blobstore.BlobStorageFactory
}

func NewDeviceServicer(factory blobstore.BlobStorageFactory) (protos.DeviceServer, error) {
	if factory == nil {
		return nil, fmt.Errorf("storage cannot be nil")
	}
	return &deviceServicer{factory: factory}, nil
}

func (srv *deviceServicer) RegisterDevices(ctx context.Context, req *protos.RegisterOrUpdateDevicesRequest) (*commonProtos.Void, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	blobs := protos.EntitiesToBlobs(req.GetEntities())
	existingKeys, err := store.GetExistingKeys(blobs.Keys(), blobstore.SearchFilter{})
	if err != nil {
		return nil, err
	}
	if len(existingKeys) > 0 {
		return nil, fmt.Errorf("the following keys: %v are already registered", existingKeys)
	}

	err = store.CreateOrUpdate(req.NetworkID, blobs)
	if err != nil {
		return nil, err
	}

	return &commonProtos.Void{}, store.Commit()
}

func (srv *deviceServicer) UpdateDevices(ctx context.Context, req *protos.RegisterOrUpdateDevicesRequest) (*commonProtos.Void, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	blobs := protos.EntitiesToBlobs(req.GetEntities())
	err = store.CreateOrUpdate(req.NetworkID, blobs)
	if err != nil {
		return nil, err
	}

	return &commonProtos.Void{}, store.Commit()
}

func (srv *deviceServicer) GetDeviceInfo(ctx context.Context, req *protos.GetDeviceInfoRequest) (*protos.GetDeviceInfoResponse, error) {
	res := &protos.GetDeviceInfoResponse{}
	if err := req.Validate(); err != nil {
		return nil, err
	}

	ids := protos.DeviceIDsToTypeAndKey(req.DeviceIDs)
	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	blobs, err := store.GetMany(req.NetworkID, ids)
	if err != nil {
		return res, err
	}

	res.DeviceMap = protos.BlobsToEntityByDeviceID(blobs)
	return res, store.Commit()
}

func (srv *deviceServicer) DeleteDevices(ctx context.Context, req *protos.DeleteDevicesRequest) (*commonProtos.Void, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	store, err := srv.factory.StartTransaction(nil)
	if err != nil {
		return nil, err
	}
	defer store.Rollback()

	ids := protos.DeviceIDsToTypeAndKey(req.DeviceIDs)
	err = store.Delete(req.NetworkID, ids)
	if err != nil {
		return nil, err
	}

	return &commonProtos.Void{}, store.Commit()
}
