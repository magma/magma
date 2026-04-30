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
	"encoding/json"
	"fmt"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/orc8r"
	"magma/orc8r/cloud/go/services/directoryd/types"
	"magma/orc8r/cloud/go/services/state"
	state_types "magma/orc8r/cloud/go/services/state/types"
	"magma/orc8r/lib/go/protos"
)

// Helper to fetch network ID and state client
func getNetworkAndStateClient(ctx context.Context, cloud bool) (string, error, interface{}) {
	networkId, err := identity.GetClientNetworkID(ctx)
	if err != nil {
		return "", err, nil
	}
	if cloud {
		client, err := state.GetCloudStateClient()
		return networkId, err, client
	}
	client, err := state.GetStateClient()
	return networkId, err, client
}

// Helper to unmarshal DirectoryRecord from state
func unmarshalDirectoryRecord(value []byte) (*types.DirectoryRecord, error) {
	serialized := &state_types.SerializedState{}
	err := json.Unmarshal(value, serialized)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal json-encoded state proto value")
	}
	dr := &types.DirectoryRecord{}
	err = dr.UnmarshalBinary([]byte(serialized.SerializedReportedState))
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal DirectoryRecord: %v", err)
	}
	return dr, nil
}

type directoryUpdateServicer struct {
}

// NewDirectoryUpdateServicer creates & returns GatewayDirectoryServiceServer interface implementation
func NewDirectoryUpdateServicer() protos.GatewayDirectoryServiceServer {
	return &directoryUpdateServicer{}
}

// UpdateRecord creates or overwrites an existing directory_record state in the state service DB
// Current implementation will overwrite Locations instead of extending it, extension is TBD
func (d *directoryUpdateServicer) UpdateRecord(c context.Context, r *protos.UpdateRecordRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if r == nil || len(r.GetId()) == 0 {
		return ret, nil
	}

	if len(r.GetLocation()) == 0 || r.GetLocation() == "hwid" {
		gw := protos.GetClientGateway(c)
		r.Location = gw.GetHardwareId()
	}
	_, err, client := getNetworkAndStateClient(c, false)
	if err != nil {
		return ret, err
	}
	stateClient := client.(protos.StateServiceClient)
	dr := &types.DirectoryRecord{LocationHistory: []string{r.GetLocation()}, Identifiers: map[string]interface{}{}}
	for k, v := range r.GetFields() {
		dr.Identifiers[k] = v
	}
	serialized, _ := dr.MarshalBinary()
	st := &protos.State{
		Type:     orc8r.DirectoryRecordType,
		DeviceID: r.Id,
		Value:    serialized,
	}
	res, err := stateClient.ReportStates(makeOutgoingCtx(c), &protos.ReportStatesRequest{States: []*protos.State{st}})
	if err != nil {
		return ret, err
	}
	fmt.Printf("UpdateRecord: reported state for id=%s location=%s\n", r.Id, r.Location)
	nid, _ := identity.GetClientNetworkID(c)
	fmt.Printf("UpdateRecord: networkId from context=%s\n", nid)
	if len(res.GetUnreportedStates()) > 0 {
		return ret, fmt.Errorf(res.GetUnreportedStates()[0].Error)
	}
	return ret, nil
}

// DeleteRecord deletes directory record of an object from the directory service
func (d *directoryUpdateServicer) DeleteRecord(c context.Context, r *protos.DeleteRecordRequest) (*protos.Void, error) {
	ret := &protos.Void{}
	if r == nil || len(r.GetId()) == 0 {
		return ret, nil
	}
	_, err, client := getNetworkAndStateClient(c, false)
	if err != nil {
		return ret, err
	}
	stateClient := client.(protos.StateServiceClient)
	_, err = stateClient.DeleteStates(
		makeOutgoingCtx(c),
		&protos.DeleteStatesRequest{Ids: []*protos.StateID{{Type: orc8r.DirectoryRecordType, DeviceID: r.Id}}},
	)
	return ret, err
}

// GetDirectoryField returns directory field for a given id and key
func (d *directoryUpdateServicer) GetDirectoryField(
	c context.Context, r *protos.GetDirectoryFieldRequest) (*protos.DirectoryField, error) {

	ret := &protos.DirectoryField{Key: r.GetFieldKey()}
	networkId, err, client := getNetworkAndStateClient(c, true)
	if err != nil {
		return ret, err
	}
	cloudClient := client.(protos.CloudStateServiceClient)
	res, err := cloudClient.GetStates(
		makeOutgoingCtx(c),
		&protos.GetStatesRequest{
			NetworkID: networkId,
			Ids:       []*protos.StateID{{Type: orc8r.DirectoryRecordType, DeviceID: r.GetId()}},
		},
	)
	if err != nil {
		return ret, err
	}
	if len(res.GetStates()) != 1 {
		return ret, status.Errorf(codes.NotFound, "directory record for ID: %s is not found", r.GetId())
	}
	dr, err := unmarshalDirectoryRecord(res.States[0].Value)
	if err != nil {
		return ret, status.Errorf(codes.Internal, "%v", err)
	}
	if dr.Identifiers != nil {
		iVal, found := dr.Identifiers[r.GetFieldKey()]
		if found {
			if val, ok := iVal.(string); ok {
				ret.Value = val
				return ret, nil
			}
			return ret, status.Errorf(codes.NotFound, "record identifier '%s' is not a string", r.GetFieldKey())
		}
	}
	return ret, status.Errorf(codes.NotFound, "record identifier key '%s' does not exist", r.GetFieldKey())
}

// GetAllDirectoryRecords returns all directory records
func (d *directoryUpdateServicer) GetAllDirectoryRecords(
	c context.Context, r *protos.Void) (*protos.AllDirectoryRecords, error) {

	ret := &protos.AllDirectoryRecords{}
	networkId, err, client := getNetworkAndStateClient(c, true)
	fmt.Printf("GetAllDirectoryRecords: %s, %v\n", networkId, err)
	if err != nil {
		return ret, err
	}
	cloudClient := client.(protos.CloudStateServiceClient)
	res, err := cloudClient.GetStates(
		context.Background(),
		&protos.GetStatesRequest{
			NetworkID:  networkId,
			TypeFilter: []string{orc8r.DirectoryRecordType},
			LoadValues: true})
	if err != nil {
		return ret, err
	}
	ret.Records = make([]*protos.DirectoryRecord, 0, len(res.GetStates()))
	for _, st := range res.GetStates() {
		if st == nil {
			continue
		}
		dr, err := unmarshalDirectoryRecord(st.Value)
		if err != nil {
			continue
		}
		pdr := &protos.DirectoryRecord{
			Id:              st.DeviceID,
			LocationHistory: dr.LocationHistory,
			Fields:          map[string]string{},
		}
		fmt.Printf("this is hte record that we have received %s", pdr.String())
		for k, iV := range dr.Identifiers {
			if v, ok := iV.(string); ok {
				pdr.Fields[k] = v
			}
		}
		ret.Records = append(ret.Records, pdr)
	}
	return ret, nil
}

func makeOutgoingCtx(incomingCtx context.Context) context.Context {
	md, _ := metadata.FromIncomingContext(incomingCtx)
	return metadata.NewOutgoingContext(incomingCtx, md)
}
