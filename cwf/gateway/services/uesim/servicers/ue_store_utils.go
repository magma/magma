/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package servicers

import (
	cwfprotos "magma/cwf/cloud/go/protos"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/lib/go/protos"

	"github.com/pkg/errors"
)

func addUeToStore(srvstore blobstore.BlobStorageFactory, ue *cwfprotos.UEConfig) {
	blob, err := ueToBlob(ue)
	store, err := srvstore.StartTransaction(nil)
	if err != nil {
		err = errors.Wrap(err, "Error while starting transaction")
		err = ConvertStorageErrorToGrpcStatus(err)
		return
	}
	defer func() {
		switch err {
		case nil:
			if commitErr := store.Commit(); commitErr != nil {
				err = errors.Wrap(err, "Error while committing transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		default:
			if rollbackErr := store.Rollback(); rollbackErr != nil {
				err = errors.Wrap(err, "Error while rolling back transaction")
				err = ConvertStorageErrorToGrpcStatus(err)
			}
		}
	}()

	err = store.CreateOrUpdate(networkIDPlaceholder, blobstore.Blobs{blob})

}

// Converts UE data to a blob for storage.
func ueToBlob(ue *cwfprotos.UEConfig) (blobstore.Blob, error) {
	marshaledUE, err := protos.Marshal(ue)
	if err != nil {
		return blobstore.Blob{}, err
	}
	return blobstore.Blob{
		Type:  blobTypePlaceholder,
		Key:   ue.GetImsi(),
		Value: marshaledUE,
	}, nil
}
