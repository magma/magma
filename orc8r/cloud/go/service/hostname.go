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

package service

import (
	"os"

	"magma/orc8r/cloud/go/orc8r"

	"github.com/golang/glog"
)

// MustGetHostname gets the hostname of the calling service.
// The hostname is determined by checking an environment variable that's
// required to be set, so fatal on error.
func MustGetHostname() string {
	hostname, exist := os.LookupEnv(orc8r.ServiceHostnameEnvVar)
	if !exist {
		glog.Fatalf("Environment variable %s must be set (in dev: set to localhost) (in prod: set to the public IP of this pod)", orc8r.ServiceHostnameEnvVar)
	}
	return hostname
}
