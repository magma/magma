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

package state

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/storage"
	"magma/orc8r/lib/go/protos"
)

var (
	ErrMissingGateway       = status.Error(codes.PermissionDenied, "missing gateway identity")
	ErrGatewayNotRegistered = status.Error(codes.PermissionDenied, "gateway not registered")
)

func idToTK(id *protos.StateID) storage.TK {
	return storage.TK{Type: id.GetType(), Key: id.GetDeviceID()}
}

func IdsToTKs(ids []*protos.StateID) storage.TKs {
	var tks storage.TKs
	for _, id := range ids {
		tks = append(tks, idToTK(id))
	}
	return tks
}

func IdAndVersionsToTKs(IDs []*protos.IDAndVersion) storage.TKs {
	var ids storage.TKs
	for _, idAndVersion := range IDs {
		ids = append(ids, idToTK(idAndVersion.Id))
	}
	return ids
}

func BlobsToStates(blobs blobstore.Blobs) []*protos.State {
	var states []*protos.State
	for _, b := range blobs {
		st := &protos.State{
			Type:     b.Type,
			DeviceID: b.Key,
			Value:    b.Value,
			Version:  b.Version,
		}
		states = append(states, st)
	}
	return states
}

func StateToBlob(state *protos.State) blobstore.Blob {
	return blobstore.Blob{
		Type:    state.GetType(),
		Key:     state.GetDeviceID(),
		Value:   state.GetValue(),
		Version: state.GetVersion(),
	}
}
