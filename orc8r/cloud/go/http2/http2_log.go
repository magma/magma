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

package http2

import (
	"net/http"

	"github.com/golang/glog"
)

//LogRequestWithVerbosity prints out request when the service binary is run
// with log_verbosity verbosity
func LogRequestWithVerbosity(req *http.Request, verbosity glog.Level) {
	glog.V(verbosity).Infof("Printing request metadata: \nHeader: %v\n"+
		"Host: %v\nURL: %v\nTrailer: %v\nProto: %v\nRequestURI: %v\n"+
		"RemoteAddr: %v\nMethod: %v\n", req.Header, req.Host, req.URL,
		req.Trailer, req.Proto, req.RequestURI, req.RemoteAddr, req.Method)
}
