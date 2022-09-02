/*
Copyright 2022 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package message

import (
	"context"
	"fmt"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func NewRelinquishMessage(id int64) *relinquishMessage {
	return &relinquishMessage{id: id}
}

type relinquishMessage struct {
	id    int64
	delta int64
}

func (u *relinquishMessage) Send(ctx context.Context, client active_mode.ActiveModeControllerClient) error {
	req := &active_mode.AcknowledgeCbsdRelinquishRequest{Id: u.id}
	_, err := client.AcknowledgeCbsdRelinquish(ctx, req)
	return err
}

func (u *relinquishMessage) String() string {
	return fmt.Sprintf("relinquish: %d", u.id)
}
