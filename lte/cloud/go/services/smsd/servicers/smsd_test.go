/*
 *  Copyright 2020 The Magma Authors.
 *
 *  This source code is licensed under the BSD-style license found in the
 *  LICENSE file in the root directory of this source tree.
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

package servicers_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/smsd/servicers"
	"magma/lte/cloud/go/services/smsd/storage"
	"magma/lte/cloud/go/services/smsd/storage/mocks"
	"magma/lte/cloud/go/sms_ll"
	mocks2 "magma/lte/cloud/go/sms_ll/mocks"
	protos2 "magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/stretchr/testify/assert"
)

func TestSMSDServicer_GetMessages(t *testing.T) {
	store := new(mocks.SMSStorage)
	serde := new(mocks2.SMSSerde)
	srv := servicers.NewSMSDServicer(store, serde)
	ctx := getTestContext(context.Background())

	// 0 case
	store.On("GetSMSsToDeliver", "n1", []string{"IMSI1"}, 6*time.Minute).Return([]*storage.SMS{}, nil).Once()
	actual, err := srv.GetMessages(ctx, &protos.GetMessagesRequest{Imsis: []string{"IMSI1"}})
	assert.NoError(t, err)
	assert.Empty(t, actual.Messages)

	// Return some messages
	tsClock := tsProto(t, time.Unix(1000, 0))
	expClock, err := ptypes.Timestamp(tsClock)
	assert.NoError(t, err)

	serde.On("EncodeMessage", "foobar", "123", expClock, []uint8{0x1, 0x2}).
		Return([][]byte{{0x1, 0x2}, {0x3, 0x4}}, nil).
		Once()
	serde.On("EncodeMessage", "barbaz", "456", expClock, []uint8{0x2, 0x3}).
		Return([][]byte{{0x2, 0x3}, {0x4, 0x5}}, nil).
		Once()
	store.On("GetSMSsToDeliver", "n1", []string{"IMSI1"}, 6*time.Minute).
		Return(
			[]*storage.SMS{
				{
					Pk:           "1",
					Status:       storage.MessageStatus_WAITING,
					Imsi:         "IMSI1",
					SourceMsisdn: "123",
					Message:      "foobar",
					RefNums:      []byte{0x1, 0x2},
					CreatedTime:  tsClock,
				},
				{
					Pk:           "2",
					Status:       storage.MessageStatus_WAITING,
					Imsi:         "IMSI1",
					SourceMsisdn: "456",
					Message:      "barbaz",
					RefNums:      []byte{0x2, 0x3},
					CreatedTime:  tsClock,
				},
			},
			nil,
		).
		Times(2)

	actual, err = srv.GetMessages(ctx, &protos.GetMessagesRequest{Imsis: []string{"IMSI1"}})
	assert.NoError(t, err)
	expected := &protos.GetMessagesResponse{
		Messages: []*protos.SMODownlinkUnitdata{
			{
				Imsi:                "IMSI1",
				NasMessageContainer: []byte{0x1, 0x2},
			},
			{
				Imsi:                "IMSI1",
				NasMessageContainer: []byte{0x3, 0x4},
			},
			{
				Imsi:                "IMSI1",
				NasMessageContainer: []byte{0x2, 0x3},
			},
			{
				Imsi:                "IMSI1",
				NasMessageContainer: []byte{0x4, 0x5},
			},
		},
	}
	assert.Equal(t, expected, actual)

	// Error in encoding
	// mock store will return the configured response from above (times 2)
	serde.On("EncodeMessage", "foobar", "123", expClock, []uint8{0x1, 0x2}).
		Return(nil, errors.New("oopsies")).
		Once()
	actual, err = srv.GetMessages(ctx, &protos.GetMessagesRequest{Imsis: []string{"IMSI1"}})
	assert.EqualError(t, err, "rpc error: code = Internal desc = could not encode message 1: oopsies")
	assert.Empty(t, actual.Messages)

	// Error fetching from store
	store.On("GetSMSsToDeliver", "n1", []string{"IMSI2", "IMSI3"}, 6*time.Minute).Return(nil, errors.New("oop")).Once()
	actual, err = srv.GetMessages(ctx, &protos.GetMessagesRequest{Imsis: []string{"IMSI2", "IMSI3"}})
	assert.EqualError(t, err, "rpc error: code = Internal desc = oop")
	assert.Empty(t, actual.Messages)

	serde.AssertExpectations(t)
	store.AssertExpectations(t)
}

func TestSMSDServicer_ReportDelivery(t *testing.T) {
	store := new(mocks.SMSStorage)
	serde := new(mocks2.SMSSerde)
	srv := servicers.NewSMSDServicer(store, serde)
	ctx := getTestContext(context.Background())

	// 0 case
	_, err := srv.ReportDelivery(ctx, &protos.ReportDeliveryRequest{Report: nil})
	assert.NoError(t, err)

	// Happy paths, delivered and failed messages
	expDelivered := map[string][]storage.SMSRef{"IMSI1": {0x1}}
	expFailed := map[string][]storage.SMSFailureReport{}
	expNasContainer := []byte{0x1, 0x2}

	// args are refs so modifications below will update the mock expectation
	store.On("ReportDelivery", "n1", expDelivered, expFailed).
		Return(nil).
		Twice()
	serde.On("DecodeDelivery", expNasContainer).
		Return(sms_ll.SMSDeliveryReport{Reference: 1, IsSuccessful: true}, nil).
		Once()
	_, err = srv.ReportDelivery(ctx, &protos.ReportDeliveryRequest{Report: &protos.SMOUplinkUnitdata{
		Imsi:                "IMSI1",
		NasMessageContainer: expNasContainer,
	}})
	assert.NoError(t, err)

	expFailed["IMSI1"] = []storage.SMSFailureReport{{
		Ref:          1,
		ErrorMessage: "foobar",
	}}
	delete(expDelivered, "IMSI1")
	serde.On("DecodeDelivery", expNasContainer).
		Return(sms_ll.SMSDeliveryReport{Reference: 1, IsSuccessful: false, ErrorMessage: "foobar"}, nil).
		Twice()
	_, err = srv.ReportDelivery(ctx, &protos.ReportDeliveryRequest{Report: &protos.SMOUplinkUnitdata{
		Imsi:                "IMSI1",
		NasMessageContainer: expNasContainer,
	}})
	assert.NoError(t, err)

	// storage error
	store.On("ReportDelivery", "n1", expDelivered, expFailed).Return(errors.New("store")).Once()
	_, err = srv.ReportDelivery(ctx, &protos.ReportDeliveryRequest{Report: &protos.SMOUplinkUnitdata{
		Imsi:                "IMSI1",
		NasMessageContainer: expNasContainer,
	}})
	assert.EqualError(t, err, "failed to report delivery: store")

	// serde error
	serde.On("DecodeDelivery", expNasContainer).
		Return(sms_ll.SMSDeliveryReport{}, errors.New("serde")).
		Once()
	_, err = srv.ReportDelivery(ctx, &protos.ReportDeliveryRequest{Report: &protos.SMOUplinkUnitdata{
		Imsi:                "IMSI1",
		NasMessageContainer: expNasContainer,
	}})
	assert.EqualError(t, err, "failed to decode report: serde")

	serde.AssertExpectations(t)
	store.AssertExpectations(t)
}

func tsProto(t *testing.T, ti time.Time) *timestamp.Timestamp {
	ret, err := ptypes.TimestampProto(ti)
	assert.NoError(t, err)
	return ret
}

func getTestContext(ctx context.Context) context.Context {
	return protos2.NewGatewayIdentity("hw1", "n1", "gw1").NewContextWithIdentity(ctx)
}
