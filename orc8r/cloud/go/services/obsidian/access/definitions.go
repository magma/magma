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

package access

// getOperator relies on x-magma-client-cert-serial HTTP request header,
// the header string is redefined here to avoid sharing it with magma GRPC
// Identity middleware & to comply with specific to Go's net/http header
// capitalization: https://golang.org/pkg/net/http/#Request
const (
	// Client Certificate CN Header (for logging only)
	CLIENT_CERT_CN_KEY = "X-Magma-Client-Cert-Cn"
	// Client Certificate Serial Number Header
	CLIENT_CERT_SN_KEY = "X-Magma-Client-Cert-Serial"
)
