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

package providers

import (
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

// StreamProvider provides a streamer policy. Given a gateway hardware ID,
// return a serialized data bundle of updates to stream back to the gateway.
type StreamProvider interface {
	// GetUpdates returns updates to stream updates back to a gateway given its hardware ID
	// If GetUpdates returns error, the stream will be closed without sending any updates
	// If GetUpdates returns error == nil, updates will be sent & the stream will be closed after that
	// If GetUpdates returns error == io.EAGAIN - the returned updates will be sent & GetUpdates will be called again
	// on the same stream
	GetUpdates(gatewayId string, extraArgs *any.Any) ([]*protos.DataUpdate, error)
}
