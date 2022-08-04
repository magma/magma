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

func NewUpdateMessage(id int64) *updateMessage {
	return &updateMessage{id: id}
}

type updateMessage struct {
	id    int64
	delta int64
}

func (u *updateMessage) Send(ctx context.Context, provider ClientProvider) error {
	req := &active_mode.AcknowledgeCbsdUpdateRequest{Id: u.id}
	client := provider.GetActiveModeClient()
	_, err := client.AcknowledgeCbsdUpdate(ctx, req)
	return err
}

func (u *updateMessage) String() string {
	return fmt.Sprintf("update: %d", u.id)
}
