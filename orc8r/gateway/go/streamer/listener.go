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

// Package streamer provides streamer client Go implementation for golang based gateways
package streamer

import (
	"magma/orc8r/lib/go/protos"

	"github.com/golang/protobuf/ptypes/any"
)

// Listener interface defines Stream Listener which will become
// the receiver of streaming updates for a registered stream
// Each received update will be unmarshalled into the Listener's update data type determined by
// the actual type returned by Listener's New() receiver method
type Listener interface {
	// GetName() returns name of the stream, the listener is getting updates on
	GetName() string
	// ReportError is going to be called by the streamer on every error.
	// If ReportError() will return nil, streamer will try to continue streaming
	// If ReportError() will return error != nil - streaming on the stream will be terminated
	ReportError(e error) error
	// Update will be called for every new update received from the stream
	// u is guaranteed to be of a type returned by New(), so - myUpdate := u.(MyDataType) should never panic
	// Update() returns bool indicating whether to continue streaming:
	//   true - continue streaming; false - stop streaming
	//   If Update() returns false -> ReportError() will be called with io.EOF,
	//   in this case, if ReportError() returns nil, streaming will continue with the new connection & stream
	Update(u *protos.DataUpdateBatch) bool
	// GetExtraArgs will be called prior to each stream request and its returned value will be used to initialize
	// ExtraArgs field in GetUpdates request payload. Most listeners may just return nil
	GetExtraArgs() *any.Any
}
