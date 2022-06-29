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

import (
	"fmt"
	"net/http"
)

type logDecorator func(string, ...interface{}) string

// getDecorator returns a logDecorator that appends remote address, URI, and
// certificate CN as available from the passed context.
func getDecorator(req *http.Request) logDecorator {
	return func(fmtStr string, args ...interface{}) string {
		if req != nil {
			fmtStr += "; remote: %s, URI: %s"
			args = append(args, req.RemoteAddr, req.RequestURI)
			ccn := req.Header.Get(CLIENT_CERT_CN_KEY)
			if len(ccn) > 0 {
				fmtStr += ", cert CN: %s"
				args = append(args, ccn)
			}
		}
		return fmt.Sprintf(fmtStr, args...)
	}
}
