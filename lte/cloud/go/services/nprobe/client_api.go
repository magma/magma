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

package nprobe

import (
	"context"

	"magma/lte/cloud/go/services/nprobe/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/registry"

	"github.com/golang/glog"
)

// GetNProbeState retrieves an nprobe state for a target id.
func GetNProbeState(networkID, taskID string) (int64, uint64, error) {
	client, err := getClient()
	if err != nil {
		return 0, 0, err
	}

	resp, err := client.GetNProbeState(
		context.Background(),
		&protos.GetNProbeStateRequest{
			NetworkId: networkID,
			TaskId:    taskID,
		},
	)
	if err != nil {
		return 0, 0, err
	}
	return resp.LastExported, resp.SequenceNumber, nil
}

// SetNProbeState update an nprobe state for a target id.
func SetNProbeState(networkID, taskID, targetID string, timestamp int64, seqNbr uint64) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	state := &protos.NProbeState{
		LastExported:   timestamp,
		SequenceNumber: seqNbr,
	}

	_, err = client.SetNProbeState(
		context.Background(),
		&protos.SetNProbeStateRequest{
			NetworkId:   networkID,
			TaskId:      taskID,
			TargetId:    targetID,
			NprobeState: state,
		},
	)
	return err
}

// DeleteNProbeState deletes an nprobe state for a target id.
func DeleteNProbeState(networkID, taskID string) error {
	client, err := getClient()
	if err != nil {
		return err
	}

	_, err = client.DeleteNProbeState(
		context.Background(),
		&protos.DeleteNProbeStateRequest{
			NetworkId: networkID,
			TaskId:    taskID,
		},
	)
	return err
}

func getClient() (protos.NProbeStateServiceClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return protos.NewNProbeStateServiceClient(conn), nil
}
