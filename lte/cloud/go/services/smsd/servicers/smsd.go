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

package servicers

import (
	"context"
	"time"

	lteProtos "magma/lte/cloud/go/protos"
	"magma/lte/cloud/go/services/smsd/storage"
	"magma/lte/cloud/go/sms_ll"
	"magma/orc8r/cloud/go/identity"

	"github.com/golang/protobuf/ptypes"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const defaultTimeout = 6 * time.Minute

type smsdServicer struct {
	store storage.SMSStorage
	serde sms_ll.SMSSerde
}

func NewSMSDServicer(store storage.SMSStorage, serde sms_ll.SMSSerde) lteProtos.SmsDServer {
	return &smsdServicer{store: store, serde: serde}
}

func (s *smsdServicer) GetMessages(ctx context.Context, request *lteProtos.GetMessagesRequest) (*lteProtos.GetMessagesResponse, error) {
	networkID, err := identity.GetClientNetworkID(ctx)
	if err != nil {
		return &lteProtos.GetMessagesResponse{}, err
	}

	messages, err := s.store.GetSMSsToDeliver(networkID, request.Imsis, defaultTimeout)
	if err != nil {
		return &lteProtos.GetMessagesResponse{}, status.Error(codes.Internal, err.Error())
	}

	ret := &lteProtos.GetMessagesResponse{}
	for _, storedSMS := range messages {
		// Maybe we should aggregate encoding errors and also return the
		// messages that successfully encoded instead.
		// For now it seems a reasonable ask to delete any malformed SMS's
		// using the API.
		createdTime, err := ptypes.Timestamp(storedSMS.CreatedTime)
		if err != nil {
			return &lteProtos.GetMessagesResponse{}, status.Errorf(codes.Internal, "could not encode message timestamp %s: %s", storedSMS.Pk, err)
		}
		encodedMessages, err := s.serde.EncodeMessage(storedSMS.Message, storedSMS.SourceMsisdn, createdTime, storedSMS.RefNums)
		if err != nil {
			return &lteProtos.GetMessagesResponse{}, status.Errorf(codes.Internal, "could not encode message %s: %s", storedSMS.Pk, err)
		}

		for _, encoded := range encodedMessages {
			ret.Messages = append(ret.Messages, &lteProtos.SMODownlinkUnitdata{
				Imsi:                storedSMS.Imsi,
				NasMessageContainer: encoded,
			})
		}
	}
	return ret, nil
}

func (s *smsdServicer) ReportDelivery(ctx context.Context, request *lteProtos.ReportDeliveryRequest) (*lteProtos.ReportDeliveryResponse, error) {
	ret := &lteProtos.ReportDeliveryResponse{}
	if request.Report == nil {
		return ret, nil
	}
	networkID, err := identity.GetClientNetworkID(ctx)
	if err != nil {
		return ret, err
	}

	decoded, err := s.serde.DecodeDelivery(request.Report.NasMessageContainer)
	if err != nil {
		return ret, errors.Wrap(err, "failed to decode report")
	}

	delivered, failed := map[string][]storage.SMSRef{}, map[string][]storage.SMSFailureReport{}
	if decoded.IsSuccessful {
		delivered[request.Report.Imsi] = append(delivered[request.Report.Imsi], decoded.Reference)
	} else {
		failed[request.Report.Imsi] = append(failed[request.Report.Imsi], storage.SMSFailureReport{
			Ref:          decoded.Reference,
			ErrorMessage: decoded.ErrorMessage,
		})
	}

	err = s.store.ReportDelivery(networkID, delivered, failed)
	if err != nil {
		return ret, errors.Wrap(err, "failed to report delivery")
	}
	return ret, nil
}
