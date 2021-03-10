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

package protos

import (
	"github.com/pkg/errors"
)

func (m *GetNProbeStateRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}

func (m *SetNProbeStateRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	if m.TargetId == "" {
		return errors.New("target ID cannot be empty")
	}
	return nil
}

func (m *DeleteNProbeStateRequest) Validate() error {
	if m.NetworkId == "" {
		return errors.New("network ID cannot be empty")
	}
	if m.TaskId == "" {
		return errors.New("task ID cannot be empty")
	}
	return nil
}
