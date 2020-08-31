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

	"magma/lte/cloud/go/lte"
	policydb_streamer "magma/lte/cloud/go/services/policydb/streamer"
	subscriber_streamer "magma/lte/cloud/go/services/subscriberdb/streamer"
	streamer_protos "magma/orc8r/cloud/go/services/streamer/protos"
	"magma/orc8r/cloud/go/services/streamer/providers"
	"magma/orc8r/lib/go/protos"
)

type providerServicer struct{}

func NewProviderServicer() streamer_protos.StreamProviderServer {
	return &providerServicer{}
}

func (s *providerServicer) GetUpdates(ctx context.Context, req *protos.StreamRequest) (*protos.DataUpdateBatch, error) {
	var streamer providers.StreamProvider
	switch req.GetStreamName() {
	case lte.SubscriberStreamName:
		streamer = &subscriber_streamer.SubscribersProvider{}
	case lte.PolicyStreamName:
		streamer = &policydb_streamer.PoliciesProvider{}
	case lte.ApnRuleMappingsStreamName:
		streamer = &policydb_streamer.ApnRuleMappingsProvider{}
	case lte.BaseNameStreamName:
		streamer = &policydb_streamer.BaseNamesProvider{}
	case lte.NetworkWideRulesStreamName:
		streamer = &policydb_streamer.NetworkWideRulesProvider{}
	case lte.RatingGroupStreamName:
		streamer = &policydb_streamer.RatingGroupsProvider{}
	default:
		return nil, fmt.Errorf("GetUpdates failed: unknown stream name provided: %s", req.GetStreamName())
	}

	updates, err := streamer.GetUpdates(req.GetGatewayId(), req.GetExtraArgs())
	if err != nil {
		return &protos.DataUpdateBatch{}, err
	}
	return &protos.DataUpdateBatch{Updates: updates}, nil
}
