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

	"magma/dp/cloud/go/active_mode_controller/protos/requests"
)

func NewSasMessage(data string) *sasMessage {
	return &sasMessage{data: data}
}

type sasMessage struct {
	data string
}

func (s *sasMessage) Send(ctx context.Context, provider ClientProvider) error {
	payload := &requests.RequestPayload{Payload: s.data}
	client := provider.GetRequestsClient()
	_, err := client.UploadRequests(ctx, payload)
	return err
}

func (s *sasMessage) String() string {
	return fmt.Sprintf("request: %s", s.data)
}
