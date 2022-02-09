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
	"fmt"

	"github.com/google/uuid"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/sbi"
	sbi_NpcfSMPolicyControlServer "magma/feg/gateway/sbi/specs/TS29512NpcfSMPolicyControlServer"
	"magma/feg/gateway/services/testcore/mock_driver"
)

var BASE_PATH = "/sm-policy-control/v1"

type MockPCFServer struct {
	// HTTP SBI server
	*sbi.NotifierServer
	// N7 server config
	serverConfig *sbi.NotifierConfig
	// subsribers stores the policy rules with imsi as the key
	subscribers map[string]*subscriberAccount
	// policySessions stores the SM Policy sessions with policyId as the key
	policySessions map[string]*policySession
	// policyIdsByImsi is a map of SmPolicyIds by IMSI
	policyIdsByImsi map[string][]string
	//
	mockDriver *mock_driver.MockDriver
	// Configuration set for the GRPC service
	serviceConfig *protos.PCFConfigs
}

type subscriberAccount struct {
	policyDecisions []policyDecision
}

type policyDecision struct {
	pduSessionId uint32
	decision     *sbi_NpcfSMPolicyControlServer.SmPolicyDecision
}

type policySession struct {
	policyId        string
	imsi            string
	pduSessionId    uint32
	notificationUrl string
	policy          *sbi_NpcfSMPolicyControlServer.SmPolicyControl
}

func NewMockPCFServer(serverConfig *sbi.NotifierConfig) (*MockPCFServer, error) {
	pcfServer := &MockPCFServer{
		serverConfig:    serverConfig,
		subscribers:     map[string]*subscriberAccount{},
		policySessions:  map[string]*policySession{},
		policyIdsByImsi: map[string][]string{},
	}
	pcfServer.NotifierServer = sbi.NewNotifierServer(*serverConfig)
	sbi_NpcfSMPolicyControlServer.RegisterHandlersWithBaseURL(pcfServer.NotifierServer.Server, pcfServer, BASE_PATH)

	err := pcfServer.NotifierServer.Start()
	if err != nil {
		return nil, err
	}
	return pcfServer, nil
}

func (srv *MockPCFServer) addPolicyDecision(
	account *subscriberAccount,
	decision *protos.PolicyDecision,
	sbiPolicyDecision *sbi_NpcfSMPolicyControlServer.SmPolicyDecision,
) {
	found := false
	for i, item := range account.policyDecisions {
		if decision.PduSessionId == item.pduSessionId {
			// replace
			account.policyDecisions[i] = policyDecision{
				pduSessionId: decision.PduSessionId,
				decision:     sbiPolicyDecision,
			}
			found = true
			break
		}
	}
	if !found {
		// add new
		account.policyDecisions = append(account.policyDecisions, policyDecision{
			pduSessionId: decision.PduSessionId,
			decision:     sbiPolicyDecision,
		})
	}
}

func (srv *MockPCFServer) fetchPolicyDecision(imsi string, pduSessionId uint32) (decision *sbi_NpcfSMPolicyControlServer.SmPolicyDecision, err error) {
	account, found := srv.subscribers[imsi]
	if !found {
		err = fmt.Errorf("subsriber account not found for %s", imsi)
		return
	}
	for _, item := range account.policyDecisions {
		if pduSessionId == 0 || pduSessionId == item.pduSessionId {
			decision = item.decision
			return
		}
	}
	err = fmt.Errorf("subsriber account not found for imsi %s and pduSession %d", imsi, pduSessionId)
	return
}

func (srv *MockPCFServer) createPolicySession(
	ctxData *sbi_NpcfSMPolicyControlServer.SmPolicyContextData,
	policy *sbi_NpcfSMPolicyControlServer.SmPolicyControl,
) string {
	policyId := uuid.New().String()
	srv.policySessions[policyId] = &policySession{
		policyId:        policyId,
		imsi:            string(ctxData.Supi),
		pduSessionId:    uint32(ctxData.PduSessionId),
		notificationUrl: string(ctxData.NotificationUri),
		policy:          policy,
	}
	policyIdList, found := srv.policyIdsByImsi[string(ctxData.Supi)]
	if !found {
		policyIdList = []string{}
	}
	srv.policyIdsByImsi[string(ctxData.Supi)] = append(policyIdList, policyId)
	return policyId
}

func (srv *MockPCFServer) deletePolicySession(policyId string) error {
	sess, found := srv.policySessions[policyId]
	if !found {
		return fmt.Errorf("policy session for %s not found", policyId)
	}
	delete(srv.policySessions, policyId)
	policyIdList, found := srv.policyIdsByImsi[sess.imsi]
	if !found {
		return fmt.Errorf("policy id %s for imsi %s not found", policyId, sess.imsi)
	}

	for i, item := range policyIdList {
		if item == policyId {
			// remove the policyId from the list
			policyIdList[i] = policyIdList[len(policyIdList)-1]
			srv.policyIdsByImsi[sess.imsi] = policyIdList[:len(policyIdList)-1]
			break
		}
	}

	return nil
}

func (srv *MockPCFServer) getPolicySessionByImsiAndPduSession(imsi string, pduSessionId uint32) (sess *policySession, err error) {
	policyIdList, found := srv.policyIdsByImsi[imsi]
	if !found {
		err = fmt.Errorf("policy id imsi %s not found", imsi)
		return
	}
	for _, item := range policyIdList {
		sess, found = srv.policySessions[item]
		if !found {
			err = fmt.Errorf("session not found for imsi %s and policy id %s", imsi, item)
			return
		}
		if pduSessionId == 0 || pduSessionId == sess.pduSessionId {
			return
		}
	}
	err = fmt.Errorf("policy session not found for imsi %s", imsi)
	return
}

func (srv *MockPCFServer) getSmPolicyUrl(policyId string) string {
	return fmt.Sprintf("%s/sm-policies/%s", srv.serverConfig.NotifyApiRoot, policyId)
}
