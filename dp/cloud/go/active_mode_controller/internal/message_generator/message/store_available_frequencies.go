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
	"strconv"
	"strings"

	"magma/dp/cloud/go/active_mode_controller/protos/active_mode"
)

func NewStoreAvailableFrequenciesMessage(id int64, frequencies []uint32) *storeAvailableFrequenciesMessage {
	return &storeAvailableFrequenciesMessage{
		id:          id,
		frequencies: frequencies,
	}
}

type storeAvailableFrequenciesMessage struct {
	id          int64
	frequencies []uint32
}

func (s *storeAvailableFrequenciesMessage) Send(ctx context.Context, provider ClientProvider) error {
	req := &active_mode.StoreAvailableFrequenciesRequest{
		Id:                   s.id,
		AvailableFrequencies: s.frequencies,
	}
	client := provider.GetActiveModeClient()
	_, err := client.StoreAvailableFrequencies(ctx, req)
	return err
}

func (s *storeAvailableFrequenciesMessage) String() string {
	b := strings.Builder{}
	_, _ = b.WriteString("store available frequencies: ")
	_, _ = b.WriteString(strconv.FormatInt(s.id, 10))
	_, _ = b.WriteString(" (")
	for i, f := range s.frequencies {
		_, _ = b.WriteString(strconv.FormatUint(uint64(f), 2))
		if i != len(s.frequencies)-1 {
			_, _ = b.WriteString(", ")
		}
	}
	_, _ = b.WriteString(")")
	return b.String()
}
