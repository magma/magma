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
// Package providers encapsulates supported EAP Authenticator Providers
//
//go:generate protoc -I ../protos -I . --go_out=plugins=grpc,paths=source_relative:. protos/eap_provider.proto
//
package providers

import (
	"fmt"

	"magma/feg/gateway/services/aaa/protos"
)

// Method is the Interface for Eap Provider
type Method interface {
	// Stringer -> String() string with Provider Name/description
	fmt.Stringer
	// EAPType should return a valid EAP Type
	EAPType() uint8
	// Handle - handles EAP Resp message (protos.EapRequest)
	Handle(*protos.Eap) (*protos.Eap, error)
	// WillHandleIdentity returns true if the provider 1) recognizes the given Identity and 2) can handle authentication
	// for this type of identity.
	// Note: a negative (false) result doesn't necessary mean that the provider cannot handle the auth for the client,
	//       it may also mean that the client did not pass enough information for the provider to recognize it
	WillHandleIdentity(identityData []byte) bool
}
