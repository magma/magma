/*
Copyright 2021 The Magma Authors.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sbi

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"magma/orc8r/lib/go/util"
)

// SbiLogger logs to glog when verbose level is set to 2
type SbiLogger struct{}

// LogRequest logs http request relateed info before making the client request
func (logger *SbiLogger) LogRequest(method string, url *url.URL, reqBody []byte, header http.Header) {
	var extraReqInfo string
	if len(reqBody) != 0 {
		extraReqInfo = fmt.Sprintf("\nBody = %s", string(reqBody))
	}
	util.SafeLog("Request %s %v %s\n", method, url.Path, extraReqInfo)
}

// LogResponse logs http response related info after receiving the response from the server
func (logger *SbiLogger) LogResponse(url *url.URL, status string, resBody []byte, header http.Header, latency time.Duration) {
	var extraResInfo string
	location, found := header["Location"]
	if found {
		extraResInfo = fmt.Sprintf("\nLocation = %s", location[0])
	}
	if len(resBody) != 0 {
		extraResInfo = fmt.Sprintf("%s\nBody = %s", extraResInfo, string(resBody))
	}
	util.SafeLog("Response %v for %v took %dms %s\n", status, url.Path, latency.Milliseconds(), extraResInfo)
}
