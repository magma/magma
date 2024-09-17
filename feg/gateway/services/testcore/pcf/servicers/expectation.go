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

package servicers

import (
	"encoding/json"
	"fmt"
	"reflect"

	"magma/feg/cloud/go/protos"
	n7_server "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
)

type N7Expectation struct {
	*protos.N7Expectation
}

func (e N7Expectation) GetAnswer() interface{} {
	return e.Answer
}

func (e N7Expectation) DoesMatch(message interface{}) error {
	switch e.RequestType {
	case protos.N7Expectation_CREATE:
		return matchPolicyContext(e.ExpectedRequest, message.(*n7_server.SmPolicyContextData))
	case protos.N7Expectation_UPDATE:
		return matchUpdatePolicyContext(e.ExpectedRequest, message.(*n7_server.SmPolicyUpdateContextData))
	case protos.N7Expectation_TERMINATE:
		return matchDeletePolicyContext(e.ExpectedRequest, message.(*n7_server.SmPolicyDeleteData))
	}
	return fmt.Errorf("invalid request type when matching N7Expectation")
}

func matchPolicyContext(expected string, actualPolicyContext *n7_server.SmPolicyContextData) error {
	var expPolicyContext n7_server.SmPolicyContextData
	err := json.Unmarshal([]byte(expected), &expPolicyContext)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected policy context: %s", err)
	}

	if !reflect.DeepEqual(&expPolicyContext, actualPolicyContext) {
		return fmt.Errorf("expected=%v, actual=%v", expPolicyContext, actualPolicyContext)
	}
	return nil
}

func matchUpdatePolicyContext(expected string, actualPolicyContext *n7_server.SmPolicyUpdateContextData) error {
	var expPolicyContext n7_server.SmPolicyUpdateContextData
	err := json.Unmarshal([]byte(expected), &expPolicyContext)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected update policy context: %s", err)
	}
	if !reflect.DeepEqual(&expPolicyContext, actualPolicyContext) {
		return fmt.Errorf("expected=%v, actual=%v", expPolicyContext, actualPolicyContext)
	}
	return nil
}

func matchDeletePolicyContext(expected string, actualPolicyContext *n7_server.SmPolicyDeleteData) error {
	var expPolicyContext n7_server.SmPolicyDeleteData
	err := json.Unmarshal([]byte(expected), &expPolicyContext)
	if err != nil {
		return fmt.Errorf("failed to unmarshal expected update policy context: %s", err)
	}
	if !reflect.DeepEqual(&expPolicyContext, actualPolicyContext) {
		return fmt.Errorf("expected=%v, actual=%v", expPolicyContext, actualPolicyContext)
	}
	return nil
}
