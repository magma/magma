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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	n7_server "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
	"magma/feg/gateway/services/testcore/mock_driver"
	lteprotos "magma/lte/cloud/go/protos"
	orcprotos "magma/orc8r/lib/go/protos"
)

const (
	TERM_NOTIF_JSON = `{
		"resourceUri": %s,
		"cause": %s
	}`
)

func (srv *MockPCFServer) SetPCFConfigs(_ context.Context, pcfConfig *protos.PCFConfigs) (*orcprotos.Void, error) {
	srv.serviceConfig = pcfConfig
	return &orcprotos.Void{}, nil
}

func (srv *MockPCFServer) CreateAccount(_ context.Context, subscriberID *lteprotos.SubscriberID) (*orcprotos.Void, error) {
	srv.subscribers[subscriberID.Id] = &subscriberAccount{
		policyDecisions: []policyDecision{},
	}
	glog.V(2).Infof("New account %s added", subscriberID.Id)
	return &orcprotos.Void{}, nil
}

func (srv *MockPCFServer) SetAccountRules(_ context.Context, decision *protos.PolicyDecision) (*orcprotos.Void, error) {
	account, found := srv.subscribers[decision.Imsi]
	if !found {
		err := fmt.Errorf("SetAccountRules subsriber %s not found", decision.Imsi)
		glog.Errorf("%s", err)
		return nil, err
	}
	var sbiPolicyDecision n7_server.SmPolicyDecision
	err := json.Unmarshal([]byte(decision.PolicyDecisionJson), &sbiPolicyDecision)
	if err != nil {
		err = fmt.Errorf("SetAccountRules unable to unmarshal policyDecisionjson: %s", err)
		glog.Errorf("%s", err)
		return nil, err
	}
	srv.addPolicyDecision(account, decision, &sbiPolicyDecision)
	glog.V(2).Infof("Subscriber rules set for %s", decision.Imsi)
	return &orcprotos.Void{}, nil
}

func (srv *MockPCFServer) ClearSubscribers(_ context.Context, void *orcprotos.Void) (*orcprotos.Void, error) {
	srv.policyIdsByImsi = map[string][]string{}
	srv.policySessions = map[string]*policySession{}
	srv.subscribers = map[string]*subscriberAccount{}
	glog.V(2).Info("All accounts deleted.")
	return &orcprotos.Void{}, nil
}

func (srv *MockPCFServer) SmPolicyUpdateNotify(_ context.Context, policyDecision *protos.PolicyDecision) (*protos.UpdateNotificationAnswer, error) {
	sess, err := srv.getPolicySessionByImsiAndPduSession(policyDecision.Imsi, policyDecision.PduSessionId)
	if err != nil {
		err = fmt.Errorf("SmPolicyUpdateNotify session fetch error: %s", err)
		glog.Error(err)
		return nil, err
	}

	resp, err := http.Post(
		sess.notificationUrl,
		"application/json",
		bytes.NewBuffer([]byte(policyDecision.PolicyDecisionJson)),
	)
	if err != nil {
		err = fmt.Errorf("SmPolicyUpdateNotify error on POST: %s", err)
		glog.Error(err)
		return nil, err
	}
	var body []byte
	if resp.Body != nil {
		defer resp.Body.Close()
		body, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			err = fmt.Errorf("SmPolicyUpdateNotify error reading POST response: %s", err)
			glog.Error(err)
			return nil, err
		}
	}
	return &protos.UpdateNotificationAnswer{
		StatusCode:               uint32(resp.StatusCode),
		PartialSuccessReportJson: string(body),
	}, nil
}

func (srv *MockPCFServer) SmPolicyTerminate(_ context.Context, termNotif *protos.TerminateNotification) (*protos.TerminateNotificationAnswer, error) {
	sess, err := srv.getPolicySessionByImsiAndPduSession(termNotif.Imsi, termNotif.PduSessionId)
	if err != nil {
		err = fmt.Errorf("SmPolicyUpdateNotify session fetch error: %s", err)
		glog.Error(err)
		return nil, err
	}

	term_json := fmt.Sprintf(TERM_NOTIF_JSON, srv.getSmPolicyUrl(sess.policyId), termNotif.ReleaseCause)
	resp, err := http.Post(
		sess.notificationUrl,
		"application/json",
		bytes.NewBuffer([]byte(term_json)),
	)
	if err != nil {
		err = fmt.Errorf("SmPolicyUpdateNotify error on POST: %s", err)
		glog.Error(err)
		return nil, err
	}

	return &protos.TerminateNotificationAnswer{
		StatusCode: uint32(resp.StatusCode),
	}, nil
}

func (srv *MockPCFServer) SetExpectations(_ context.Context, n7Expectations *protos.N7Expectations) (*orcprotos.Void, error) {
	es := []mock_driver.Expectation{}
	for _, e := range n7Expectations.Expectations {
		es = append(es, mock_driver.Expectation(N7Expectation{e}))
	}
	srv.mockDriver = mock_driver.NewMockDriver(es, n7Expectations.UnexpectedRequestBehavior, n7Expectations.DefaultAnswer)
	return &orcprotos.Void{}, nil
}

func (srv *MockPCFServer) AssertExpectations(_ context.Context, _ *orcprotos.Void) (*protos.N7ExpectationResult, error) {
	srv.mockDriver.Lock()
	defer srv.mockDriver.Unlock()

	results, errs := srv.mockDriver.AggregateResults()
	return &protos.N7ExpectationResult{Results: results, Errors: errs}, nil
}
