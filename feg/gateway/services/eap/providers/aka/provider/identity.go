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

// Package aka implements EAP-AKA provider
package provider

import "regexp"

var akaRe = regexp.MustCompile(`^0\d{6,15}@\w(?:\w|\.|-)*\w$`)

// WillHandleIdentity returns true if the provider 1) recognizes the given Identity and 2) can hendle authentication
// for this type of identity.
// Note: a negative (false) result doesn't necessary mean that the provider cannot handle the auth for the client,
//       it may also mean that the client did not pass enough information for the provider to recognize it
func (p *providerImpl) WillHandleIdentity(identityData []byte) bool {
	return len(identityData) > 10 && akaRe.Match(identityData)
}
