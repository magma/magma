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

package storage

import (
	"testing"
	"time"

	"magma/lte/cloud/go/services/nprobe/obsidian/models"
	"magma/orc8r/cloud/go/blobstore"
	"magma/orc8r/cloud/go/blobstore/mocks"
	"magma/orc8r/cloud/go/storage"

	strfmt "github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

const (
	placeholderNetworkID = "placeholder_network"
)

func TestStoreNProbeData(t *testing.T) {
	var blobFactMock *mocks.BlobStorageFactory
	var blobStoreMock *mocks.TransactionalBlobStorage

	taskID := "task_id1"
	nprobeData := models.NetworkProbeData{
		LastExported:   strfmt.DateTime(time.Now()),
		TargetID:       "imsi01",
		SequenceNumber: 0,
	}

	blob, err := nprobeDataToBlob(taskID, nprobeData)
	assert.NoError(t, err)

	// Store nprobe data
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("CreateOrUpdate", placeholderNetworkID, blobstore.Blobs{blob}).
		Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()

	store := NewNProbeBlobstore(blobFactMock)
	err = store.StoreNProbeData(placeholderNetworkID, taskID, nprobeData)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Get nprobe data
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	tk := storage.TypeAndKey{Type: NProbeBlobType, Key: taskID}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Get", placeholderNetworkID, tk).Return(blob, nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()

	store = NewNProbeBlobstore(blobFactMock)
	nprobeReceived, err := store.GetNProbeData(placeholderNetworkID, taskID)
	assert.NoError(t, err)
	assert.Equal(t, "imsi01", nprobeReceived.TargetID)
	assert.Equal(t, uint32(0), nprobeReceived.SequenceNumber)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)

	// Delete nprobe data
	blobFactMock = &mocks.BlobStorageFactory{}
	blobStoreMock = &mocks.TransactionalBlobStorage{}
	tkSet := []storage.TypeAndKey{tk}
	blobFactMock.On("StartTransaction", mock.Anything).Return(blobStoreMock, nil).Once()
	blobStoreMock.On("Rollback").Return(nil).Once()
	blobStoreMock.On("Delete", placeholderNetworkID, tkSet).Return(nil).Once()
	blobStoreMock.On("Commit").Return(nil).Once()
	store = NewNProbeBlobstore(blobFactMock)

	err = store.DeleteNProbeData(placeholderNetworkID, taskID)
	assert.NoError(t, err)
	blobFactMock.AssertExpectations(t)
	blobStoreMock.AssertExpectations(t)
}
