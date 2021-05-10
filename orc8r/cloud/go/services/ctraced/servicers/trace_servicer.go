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

	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/serdes"
	"magma/orc8r/cloud/go/services/configurator"
	"magma/orc8r/cloud/go/services/ctraced/obsidian/models"
	"magma/orc8r/cloud/go/services/ctraced/storage"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type callTraceServicer struct {
	storage storage.CtracedStorage
}

func NewCallTraceServicer(storage storage.CtracedStorage) protos.CallTraceControllerServer {
	return &callTraceServicer{storage: storage}
}

func (srv *callTraceServicer) ReportEndedCallTrace(ctx context.Context, req *protos.ReportEndedTraceRequest) (*protos.ReportEndedTraceResponse, error) {
	networkID, err := getNetworkID(ctx)
	if err != nil {
		return nil, err
	}
	callTrace, err := getCallTraceModel(networkID, req.TraceId)
	if err != nil {
		return nil, err
	}

	err = srv.storage.StoreCallTrace(networkID, req.TraceId, req.TraceContent)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, fmt.Sprintf("failed to save call trace data, network-id: %s, gateway-id: %s, calltrace-id: %s", networkID, callTrace.Config.GatewayID, req.TraceId))
	}

	callTrace.State.CallTraceEnding = req.Success
	callTrace.State.CallTraceAvailable = req.Success

	update := configurator.EntityUpdateCriteria{
		Type:      orc8r.CallTraceEntityType,
		Key:       req.TraceId,
		NewConfig: callTrace,
	}

	_, err = configurator.UpdateEntity(networkID, update, serdes.Entity)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, fmt.Sprintf("failed to update call trace, network-id: %s, gateway-id: %s, calltrace-id: %s", networkID, callTrace.Config.GatewayID, req.TraceId))
	}
	return &protos.ReportEndedTraceResponse{}, nil
}

func getNetworkID(ctx context.Context) (string, error) {
	id, err := protos.GetGatewayIdentity(ctx)
	if err != nil {
		return "", err
	}
	return id.GetNetworkId(), nil
}

func getCallTraceModel(networkID string, callTraceID string) (*models.CallTrace, error) {
	ent, err := configurator.LoadEntity(
		networkID, orc8r.CallTraceEntityType, callTraceID,
		configurator.EntityLoadCriteria{LoadConfig: true},
		serdes.Entity,
	)
	if err == merrors.ErrNotFound {
		return nil, status.Errorf(codes.InvalidArgument, "Call trace not found")
	}
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "Failed to load call trace")
	}
	callTrace := &models.CallTrace{}
	err = callTrace.FromBackendModels(ent)
	if err != nil {
		return nil, status.Errorf(codes.Aborted, "Failed to load call trace")
	}
	return callTrace, nil
}
