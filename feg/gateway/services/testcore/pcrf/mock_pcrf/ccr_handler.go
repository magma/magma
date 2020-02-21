/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 */

package mock_pcrf

import (
	"errors"
	"fmt"

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
func getCCRHandler(srv *PCRFDiamServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received CCR from %s\n", c.RemoteAddr())
		var ccr ccrMessage
		if err := m.Unmarshal(&ccr); err != nil {
			glog.Errorf("Failed to unmarshal CCR %s", err)
			return
		}
		imsi, err := getIMSI(ccr)
		if err != nil {
			glog.Errorf("Could not parse CCR: %s", err.Error())
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account, found := srv.subscribers[imsi]
		if !found {
			glog.Error("IMSI not found in subscribers")
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}

		avps := []*diam.AVP{}
		if credit_control.CreditRequestType(ccr.RequestType) == credit_control.CRTInit {
			// Install all rules attached to the subscriber for the initial answer
			ruleInstalls := toRuleInstallAVPs(account.RuleNames, account.RuleBaseNames, account.RuleDefinitions)
			// Install all monitors attached to the subscriber for the initial answer
			usageMonitors := toUsageMonitoringInfoAVPs(account.UsageMonitors)
			avps = append(ruleInstalls, usageMonitors...)
		} else {
			// Update the subscriber state with the usage updates in CCR-U/T
			creditByMkey, err := updateSubscriberAccountWithUsageUpdates(account, ccr.UsageMonitors)
			if err != nil {
				glog.Errorf("Failed to update quota: %v", err)
			}
			avps = toUsageMonitoringInfoAVPs(creditByMkey)
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

func getIMSI(message ccrMessage) (string, error) {
	for _, subID := range message.SubscriptionIDs {
		if subID.IDType == credit_control.EndUserIMSI {
			return subID.IDData, nil
		}
	}
	return "", errors.New("Could not obtain IMSI from CCR message")
}

func updateSubscriberAccountWithUsageUpdates(account *subscriberAccount, monitorUpdates []*usageMonitorRequestAVP) (creditByMkey, error) {
	credits := make(creditByMkey, len(monitorUpdates))
	for _, update := range monitorUpdates {
		monitorCredit, ok := account.UsageMonitors[update.MonitoringKey]
		if !ok {
			return credits, fmt.Errorf("Unknown monitoring key %s", update.MonitoringKey)
		}
		monitorCredit.Volume -= update.UsedServiceUnit.TotalOctets
		if monitorCredit.Volume < 0 {
			monitorCredit.Volume = 0
		}
		credits[update.MonitoringKey] = monitorCredit
	}
	return credits, nil
}

func getQuotaGrant(monitorCredit *protos.UsageMonitorCredit) uint64 {
	// get the min of return bytes or the total volume
	if monitorCredit.ReturnBytes > monitorCredit.Volume {
		return monitorCredit.Volume
	}
	return monitorCredit.ReturnBytes
}
