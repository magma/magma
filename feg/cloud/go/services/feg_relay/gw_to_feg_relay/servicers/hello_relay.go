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

package servicers

import (
	"context"
	"regexp"
	"strings"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
)

// should match "... @NH-FEG-FOR: IMSI123456", "... @NH-FeG-FOR IMSI123456", @nh-Feg-for Imsi 123456", etc.
var nHTargetSuffixRe = regexp.MustCompile(`(?i)(?:\s|^)@NH-FEG-FOR:?\s*(?:IMSI:?\s*)?(\d{5,15})\s*$`)

// FeG Hello implementation
//
// SayHello sends HelloRequest to default FeG if the greeting doesn not end with '@NH-FeG-FOR <IMSI>' (see RegEx above)
// If the greeting ends with '@NH-FeG-FOR <IMSI>', SayHello will try to route the request to Neutral Host's FeG network
// matching the given IMSI
func (s *RelayRouter) SayHello(ctx context.Context, req *protos.HelloRequest) (*protos.HelloReply, error) {
	var imsi string
	match := nHTargetSuffixRe.FindStringSubmatch(req.GetGreeting())
	if len(match) > 1 {
		imsi = match[1]
		glog.V(1).Infof("SayHello with NH IMSI '%s': %s", imsi, req.GetGreeting())
		req.Greeting = strings.TrimSuffix(req.Greeting, match[0])
	}
	conn, ctx, cancel, err := s.GetFegServiceConnection(ctx, imsi, FegHello)
	if err != nil {
		glog.Errorf("SayHello error for NH IMSI '%s', greeting '%s': %v", imsi, req.GetGreeting(), err)
		return nil, err
	}
	defer cancel()
	return protos.NewHelloClient(conn).SayHello(ctx, req)
}
