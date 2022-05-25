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
	"errors"

	"github.com/thoas/go-funk"

	"magma/orc8r/lib/go/protos"
)

func ValidateGetStatesRequest(req *protos.GetStatesRequest) error {
	if !funk.IsEmpty(req.Ids) && funk.IsEmpty(req.NetworkID) {
		return errors.New("network ID must be non-empty for non-empty state IDs")
	}
	return nil
}

func ValidateReportStatesRequest(req *protos.ReportStatesRequest) error {
	if req.GetStates() == nil || len(req.GetStates()) == 0 {
		return errors.New("states value must be specified and non-empty")
	}
	return nil
}

func ValidateDeleteStatesRequest(req *protos.DeleteStatesRequest) error {
	if err := enforceNetworkID(req.NetworkID); err != nil {
		return err
	}
	if funk.IsEmpty(req.Ids) {
		return errors.New("states value must be specified and non-empty")
	}
	return nil
}

func ValidateSyncStatesRequest(req *protos.SyncStatesRequest) error {
	if req.GetStates() == nil || len(req.GetStates()) == 0 {
		return errors.New("states value must be specified and non-empty")
	}
	return nil
}

func enforceNetworkID(networkID string) error {
	if len(networkID) == 0 {
		return errors.New("network ID must be specified")
	}
	return nil
}
