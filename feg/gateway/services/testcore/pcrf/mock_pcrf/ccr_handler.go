/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mock_pcrf

import (
	"errors"
	"fmt"
	"reflect"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/session_proxy/credit_control"
	"magma/feg/gateway/services/session_proxy/credit_control/gx"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"
)

type ccrMessage struct {
	SessionID        datatype.UTF8String       `avp:"Session-Id"`
	OriginHost       datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm      datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost  datatype.DiameterIdentity `avp:"Destination-Host"`
	RequestType      datatype.Enumerated       `avp:"CC-Request-Type"`
	RequestNumber    datatype.Unsigned32       `avp:"CC-Request-Number"`
	SubscriptionIDs  []*subscriptionIDDiam     `avp:"Subscription-Id"`
	IPAddr           datatype.OctetString      `avp:"Framed-IP-Address"`
	UsageMonitors    []*usageMonitorRequestAVP `avp:"Usage-Monitoring-Information"`
	CalledStationId  string                    `avp:"Called-Station-Id"`
	EventTrigger     datatype.Enumerated       `avp:"Event-Trigger"`
}

type subscriptionIDDiam struct {
	IDType credit_control.SubscriptionIDType `avp:"Subscription-Id-Type"`
	IDData string                            `avp:"Subscription-Id-Data"`
}

type usageMonitorRequestAVP struct {
	MonitoringKey   string             `avp:"Monitoring-Key"`
	UsedServiceUnit usedServiceUnitAVP `avp:"Used-Service-Unit"`
	Level           gx.MonitoringLevel `avp:"Usage-Monitoring-Level"`
}

type usedServiceUnitAVP struct {
	InputOctets  uint64 `avp:"CC-Input-Octets"`
	OutputOctets uint64 `avp:"CC-Output-Octets"`
	TotalOctets  uint64 `avp:"CC-Total-Octets"`
}

// getCCRHandler returns a handler to be called when the server receives a CCR
func getCCRHandler(srv *PCRFServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received CCR from %s\n", c.RemoteAddr())
		srv.lastDiamMessageReceived = m
		var ccr ccrMessage
		if err := m.Unmarshal(&ccr); err != nil {
			glog.Errorf("Failed to unmarshal CCR %s", err)
			return
		}

		imsi, err := ccr.GetIMSI()
		if err != nil {
			glog.Errorf("Could not parse CCR: %s", err.Error())
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}

		// save connection for RaRs
		account, found := srv.subscribers[imsi]
		if !found {
			glog.Errorf("IMSI %v not found in subscribers", imsi)
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account.CurrentState = &SubscriberSessionState{
			Connection: c,
			SessionID:  string(ccr.SessionID),
		}

		if srv.serviceConfig.UseMockDriver {
			srv.mockDriver.Lock()
			iAnswer := srv.mockDriver.GetAnswerFromExpectations(ccr)
			srv.mockDriver.Unlock()
			if iAnswer == nil {
				sendAnswer(ccr, c, m, diam.UnableToComply)
				return
			}
			avps, resultCode := iAnswer.(GxAnswer).toAVPs()
			sendAnswer(ccr, c, m, resultCode, avps...)
			return
		}

		var avps []*diam.AVP
		ccrType := credit_control.CreditRequestType(ccr.RequestType)
		if ccrType == credit_control.CRTInit {
			glog.V(2).Infof("\tGot xCRT-Init from %s. Rules will be sent: %v - %v - %v\n", imsi, account.RuleNames, account.RuleBaseNames, account.UsageMonitors)
			// Install all rules attached to the subscriber for the initial answer
			ruleInstalls := toRuleInstallAVPs(account.RuleNames, account.RuleBaseNames, account.RuleDefinitions, nil, nil)
			// Install all monitors attached to the subscriber for the initial answer
			usageMonitors := toUsageMonitorAVPs(account.UsageMonitors)
			avps = append(ruleInstalls, usageMonitors...)
		} else {
			// Update the subscriber state with the usage updates in CCR-U/T
			glog.V(2).Infof("\tGot xCCR type \"%d\" from IMSI %s. Quota be updated\n", ccrType, imsi)
			creditByMkey, err := updateSubscriberAccountWithUsageUpdates(account, ccr.UsageMonitors)
			if err != nil {
				glog.Errorf("Failed to update quota: %v", err)
			}
			avps = toUsageMonitorAVPs(creditByMkey)
		}
		sendAnswer(ccr, c, m, diam.Success, avps...)
	}
}

// sendAnswer sends a CCA to the connection given
func sendAnswer(
	ccr ccrMessage,
	conn diam.Conn,
	message *diam.Message,
	statusCode uint32,
	additionalAVPs ...*diam.AVP,
) {
	a := message.Answer(statusCode)
	a.NewAVP(avp.OriginHost, avp.Mbit, 0, ccr.DestinationHost)
	a.NewAVP(avp.OriginRealm, avp.Mbit, 0, ccr.DestinationRealm)
	a.NewAVP(avp.DestinationRealm, avp.Mbit, 0, ccr.OriginRealm)
	a.NewAVP(avp.DestinationHost, avp.Mbit, 0, ccr.OriginHost)
	a.NewAVP(avp.CCRequestType, avp.Mbit, 0, ccr.RequestType)
	a.NewAVP(avp.CCRequestNumber, avp.Mbit, 0, ccr.RequestNumber)
	for _, avp := range additionalAVPs {
		a.InsertAVP(avp)
	}
	// SessionID must be the first AVP
	a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID))

	_, err := a.WriteTo(conn)
	if err != nil {
		glog.V(2).Infof("Failed to write message to %s: %s\n%s\n",
			conn.RemoteAddr(), err, a)
		return
	}
	glog.V(2).Infof("Sent CCA to %s:\n", conn.RemoteAddr())
}

func (m *ccrMessage) GetIMSI() (string, error) {
	for _, subID := range m.SubscriptionIDs {
		if subID.IDType == credit_control.EndUserIMSI {
			return subID.IDData, nil
		}
	}
	return "", errors.New("Could not obtain IMSI from CCR message")
}

// TODO: Remove this when not needed anymore (use findAVP from diam library)
// Searches on ccr message for an specific AVP message based on the avp tag on ccr type (ie "Session-Id")
// It returns on the first match it finds.
func GetAVP(message *ccrMessage, AVPToFind string) (interface{}, error) {
	elem := reflect.ValueOf(message)
	avpFound, err := findAVP(elem, "avp", AVPToFind)
	if err != nil {
		glog.Errorf("Failed to find %s: %s\n", AVPToFind, err)
		return "", err
	}
	return avpFound, nil
}

// Depth Search First of a specific tag:value on a element (accepts structs, pointers, slices)
func findAVP(elem reflect.Value, tag, AVPtoFind string) (interface{}, error) {
	switch elem.Kind() {
	case reflect.Ptr:
		return findAVP(elem.Elem(), tag, AVPtoFind)
	case reflect.Struct:
		for i := 0; i < elem.NumField(); i += 1 {
			fieldT := elem.Type().Field(i)
			if fieldT.Tag.Get(tag) == AVPtoFind {
				fieldV := elem.Field(i)
				return fieldV.Interface(), nil
			}
			result, err := findAVP(elem.Field(i), tag, AVPtoFind)
			if err == nil {
				return result, err
			}
		}
	case reflect.Slice:
		for i := 0; i < elem.Len(); i += 1 {
			result, err := findAVP(elem.Index(i), tag, AVPtoFind)
			if err == nil {
				return result, err
			}
		}
	}
	return "", fmt.Errorf("Could not find AVP %s:%s", tag, AVPtoFind)
}

func updateSubscriberAccountWithUsageUpdates(account *subscriberAccount, monitorUpdates []*usageMonitorRequestAVP) (creditByMkey, error) {
	credits := make(creditByMkey, len(monitorUpdates))
	for _, update := range monitorUpdates {
		monitorCredit, ok := account.UsageMonitors[update.MonitoringKey]
		if !ok {
			return credits, fmt.Errorf("Unknown monitoring key %s", update.MonitoringKey)
		}
		credits[update.MonitoringKey] = updateSubscriberQuota(monitorCredit, update)
	}
	return credits, nil
}

func updateSubscriberQuota(credit *protos.UsageMonitor, usedServiceUnit *usageMonitorRequestAVP) *protos.UsageMonitor {
	total := credit.GetTotalQuota()
	update := usedServiceUnit.UsedServiceUnit
	credit.TotalQuota.TotalOctets = decrementOrZero(total.GetTotalOctets(), update.TotalOctets)
	credit.TotalQuota.InputOctets = decrementOrZero(total.GetInputOctets(), update.InputOctets)
	credit.TotalQuota.OutputOctets = decrementOrZero(total.GetOutputOctets(), update.OutputOctets)
	return credit
}

func getQuotaGrant(monitorCredit *protos.UsageMonitor) *protos.Octets {
	total := monitorCredit.GetTotalQuota()
	perRequest := monitorCredit.GetMonitorInfoPerRequest().GetOctets()
	return &protos.Octets{
		TotalOctets:  getMin(total.GetTotalOctets(), perRequest.TotalOctets),
		InputOctets:  getMin(total.GetInputOctets(), perRequest.InputOctets),
		OutputOctets: getMin(total.GetOutputOctets(), perRequest.OutputOctets),
	}
}

func decrementOrZero(first, second uint64) uint64 {
	if second > first {
		return 0
	}
	return first - second
}

func getMin(first, second uint64) uint64 {
	if first > second {
		return second
	}
	return first
}
