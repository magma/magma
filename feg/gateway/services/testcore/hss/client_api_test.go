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

package hss_test

import (
	"context"
	"testing"

	"magma/feg/gateway/services/testcore/hss"
	"magma/feg/gateway/services/testcore/hss/storage"
	"magma/feg/gateway/services/testcore/hss/test_init"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/status"
)

func TestHSSClient(t *testing.T) {
	srv, err := test_init.StartTestService(t)
	require.NoError(t, err, "Error starting test service: %v", err)

	expectedID := "0123456789"
	expectedType1 := lteprotos.SubscriberID_IDType(1)
	sub := &lteprotos.SubscriberData{
		Sid: &lteprotos.SubscriberID{
			Id:   expectedID,
			Type: expectedType1,
		},
	}

	// Add Subscriber Request
	err = hss.AddSubscriber(sub)
	assert.NoError(t, err, "AddSubscriberError: %v", err)

	// Get Subscriber Data Request
	subRes, err := hss.GetSubscriberData(expectedID)
	assert.NoError(t, err, "GetSubscriberData Error: %v", err)
	assert.Equal(t, expectedID, subRes.GetSid().GetId())
	assert.Equal(t, expectedType1, subRes.GetSid().GetType())

	// Update Subscriber Data Request
	expectedType2 := lteprotos.SubscriberID_IDType(2)
	sub.GetSid().Type = expectedType2
	err = hss.UpdateSubscriber(sub)
	assert.NoError(t, err, "UpdateSubscriber Error: %v", err)

	// Verify that data was updated
	subRes, err = hss.GetSubscriberData(expectedID)
	assert.NoError(t, err)
	assert.Equal(t, expectedID, subRes.GetSid().GetId())
	assert.Equal(t, expectedType2, subRes.GetSid().GetType())

	// Delete Subscriber Data Request
	err = hss.DeleteSubscriber(expectedID)
	assert.NoError(t, err, "DeleteSubscriber Error: %v", err)

	// Verify that subscriber was deleted
	subRes, err = hss.GetSubscriberData(expectedID)
	assert.Nil(t, subRes)
	// check the error matches with a GRPC NewUnknownSubscriberError
	expectedGrpcErr := storage.ConvertStorageErrorToGrpcStatus(storage.NewUnknownSubscriberError(expectedID))
	expectedErr, ok := status.FromError(expectedGrpcErr)
	assert.True(t, ok)
	actualErr, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, expectedErr.Code(), actualErr.Code())
	assert.Equal(t, expectedErr.Message(), actualErr.Message())

	_, err = srv.StopService(context.Background(), &orcprotos.Void{})
	assert.NoError(t, err)
}
