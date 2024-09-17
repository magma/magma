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

// package servcers implements WiFi AAA GRPC services
package servicers

import (
	"fmt"
	"time"

	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos/mconfig"
	"magma/feg/gateway/services/aaa"
)

// GetIdleSessionTimeout returns Idle Session Timeout Duration if set in mconfigs or DefaultSessionTimeout otherwise
func GetIdleSessionTimeout(cfg *mconfig.AAAConfig) time.Duration {
	if cfg != nil {
		if tout := time.Millisecond * time.Duration(cfg.GetIdleSessionTimeoutMs()); tout > 0 {
			return tout
		}
	}
	return aaa.DefaultSessionTimeout
}

func Errorf(code codes.Code, format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	glog.Errorf("%s; [RPC: %s]", msg, code.String())
	return status.Errorf(code, msg)
}

func Error(code codes.Code, err error) error {
	if err != nil {
		if se, ok := err.(interface{ GRPCStatus() *status.Status }); ok {
			code = se.GRPCStatus().Code()
		}
		glog.Errorf("%v; [RPC: %s]", err, code.String())
		return status.Error(code, err.Error())
	}
	return nil
}
