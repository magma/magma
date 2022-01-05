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

package mock_ocs

import (
	"fmt"
	"reflect"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/golang/glog"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/services/session_proxy/credit_control"
)

type ccrMessage struct {
	SessionID          datatype.UTF8String       `avp:"Session-Id"`
	OriginHost         datatype.DiameterIdentity `avp:"Origin-Host"`
	OriginRealm        datatype.DiameterIdentity `avp:"Origin-Realm"`
	DestinationRealm   datatype.DiameterIdentity `avp:"Destination-Realm"`
	DestinationHost    datatype.DiameterIdentity `avp:"Destination-Host"`
	RequestType        datatype.Enumerated       `avp:"CC-Request-Type"`
	RequestNumber      datatype.Unsigned32       `avp:"CC-Request-Number"`
	MSCC               []*ccrCredit              `avp:"Multiple-Services-Credit-Control"`
	SubscriptionIDs    []*subscriptionID         `avp:"Subscription-Id"`
	ServiceInformation []*serviceInformation     `avp:"Service-Information"`
}

type serviceInformation struct {
	PsInformation []*psInformation `avp:"PS-Information"`
}

type psInformation struct {
	CalledStationId string `avp:"Called-Station-Id"`
}

type subscriptionID struct {
	IDType credit_control.SubscriptionIDType `avp:"Subscription-Id-Type"`
	IDData string                            `avp:"Subscription-Id-Data"`
}

type ccrCredit struct {
	RatingGroup          uint32                `avp:"Rating-Group"`
	UsedServiceUnit      *usedServiceUnit      `avp:"Used-Service-Unit"`
	ReportingReason      uint32                `avp:"Reporting-Reason"`
	RequestedServiceUnit *RequestedServiceUnit `avp:"Requested-Service-Unit"`
}

type usedServiceUnit struct {
	InputOctets     uint64 `avp:"CC-Input-Octets"`
	OutputOctets    uint64 `avp:"CC-Output-Octets"`
	TotalOctets     uint64 `avp:"CC-Total-Octets"`
	ReportingReason uint32 `avp:"Reporting-Reason"`
}

type RequestedServiceUnit struct {
	InputOctets  uint64 `avp:"CC-Input-Octets"`
	OutputOctets uint64 `avp:"CC-Output-Octets"`
	TotalOctets  uint64 `avp:"CC-Total-Octets"`
}

// getCCRHandler returns a handler to be called when the server receives a CCR
func getCCRHandler(srv *OCSDiamServer) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received CCR from %s\n", c.RemoteAddr())
		glog.V(2).Infof("Received Gy CCR message\n%s\n", m)
		srv.lastDiamMessageReceived = m
		var ccr ccrMessage
		if err := m.Unmarshal(&ccr); err != nil {
			glog.Errorf("Failed to unmarshal CCR %s", err)
			return
		}
		imsi := ccr.GetIMSI()
		if len(imsi) == 0 {
			glog.Errorf("Could not find IMSI in CCR")
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		requestType := credit_control.CreditRequestType(ccr.RequestType)
		account, found := srv.accounts[imsi]
		if !found {
			glog.Errorf("Account not found!")
			sendAnswer(ccr, c, m, diam.AuthenticationRejected)
			return
		}
		account.CurrentState = &SubscriberSessionState{
			Connection: c,
			SessionID:  string(ccr.SessionID),
		}
		if srv.ocsConfig.UseMockDriver {
			srv.mockDriver.Lock()
			iAnswer := srv.mockDriver.GetAnswerFromExpectations(ccr)
			srv.mockDriver.Unlock()
			if iAnswer == nil {
				sendAnswer(ccr, c, m, diam.UnableToComply)
				return
			}
			if iAnswer.(GyAnswer).GetLinkFailure() {
				return
			}
			avps, resultCode := iAnswer.(GyAnswer).toAVPs()
			sendAnswer(ccr, c, m, resultCode, avps...)
			return
		}

		if requestType == credit_control.CRTTerminate {
			sendAnswer(ccr, c, m, diam.Success)
			return
		}

		creditAnswers := make([]*diam.AVP, 0, len(ccr.MSCC))

		for _, mscc := range ccr.MSCC {
			// Only check usage for CCR-U and CCR-T,
			if requestType != credit_control.CRTInit {
				if mscc.UsedServiceUnit != nil {
					glog.V(2).Infof("Received credit usage from %s:%d, balance will be decremented by Total:%d Tx:%d Rx:%d",
						imsi, mscc.RatingGroup,
						mscc.UsedServiceUnit.TotalOctets,
						mscc.UsedServiceUnit.OutputOctets,
						mscc.UsedServiceUnit.InputOctets,
					)
					decrementUsedCredit(
						account.ChargingCredit[mscc.RatingGroup],
						mscc.UsedServiceUnit,
					)
					glog.V(2).Infof("Current balance for %s:%d is Total:%d Tx:%d Rx:%d",
						imsi, mscc.RatingGroup,
						account.ChargingCredit[mscc.RatingGroup].Volume.TotalOctets,
						account.ChargingCredit[mscc.RatingGroup].Volume.OutputOctets,
						account.ChargingCredit[mscc.RatingGroup].Volume.InputOctets)
				}
			}
			// Only return credit for CCR-I and CCR-U,
			if requestType != credit_control.CRTTerminate {
				// Requested-Service-Unit AVP must always exist
				if mscc.RequestedServiceUnit == nil {
					sendAnswer(ccr, c, m, diam.UnableToComply)
					return
				}

				returnOctets, final, creditLimitReached :=
					getQuotaGrant(srv, account.ChargingCredit[mscc.RatingGroup])
				if creditLimitReached {
					sendAnswer(ccr, c, m, DiameterCreditLimitReached)
					return
				}
				creditAnswers = append(
					creditAnswers,
					toGrantedUnitsAVP(
						diam.Success,
						srv.ocsConfig.ValidityTime,
						returnOctets,
						final,
						mscc.RatingGroup,
						srv.ocsConfig.FinalUnitIndication.FinalUnitAction,
						srv.ocsConfig.FinalUnitIndication.RedirectAddress,
						srv.ocsConfig.FinalUnitIndication.RestrictRules,
					))
			}
		}
		sendAnswer(ccr, c, m, diam.Success, creditAnswers...)
	}
}

func decrementUsedCredit(credit *CreditBucket, usage *usedServiceUnit) {
	credit.Volume.TotalOctets = decrementOrZero(credit.Volume.GetTotalOctets(), usage.TotalOctets)
	credit.Volume.OutputOctets = decrementOrZero(credit.Volume.GetOutputOctets(), usage.OutputOctets)
	credit.Volume.InputOctets = decrementOrZero(credit.Volume.GetInputOctets(), usage.InputOctets)
}

func decrementOrZero(first, second uint64) uint64 {
	if second >= first {
		// subtraction between uints is never negative!!
		return 0
	}
	result := first - second
	return result
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
	a.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID)
	for _, avp := range additionalAVPs {
		a.InsertAVP(avp)
	}
	// SessionID must be the first AVP
	a.InsertAVP(diam.NewAVP(avp.SessionID, avp.Mbit, 0, ccr.SessionID))

	glog.V(2).Infof("Sending Gy CCA message\n%s\n", a)
	_, err := a.WriteTo(conn)
	if err != nil {
		glog.Errorf("Failed to write message to %s: %s\n%s\n",
			conn.RemoteAddr(), err, a)
		return
	}
	glog.V(2).Infof("Sent CCA to %s:\n", conn.RemoteAddr())
}

// getIMSI finds the account IMSI in a CCR message
func (message ccrMessage) GetIMSI() string {
	for _, subID := range message.SubscriptionIDs {
		if subID.IDType == credit_control.EndUserIMSI {
			return subID.IDData
		}
	}
	return ""
}

// TODO: Remove this when not needed anymore (use findAVP from diam library)
// Searches on ccr message for an specific AVP message based on the avp tag on ccr type (ie "Session-Id")
// It returns on the first match it finds.
func GetAVP(message *ccrMessage, AVPToFind string) (interface{}, error) {
	elem := reflect.ValueOf(message)
	avpFound, err := findAVP(elem, "avp", AVPToFind)
	if err != nil {
		glog.Errorf("Failed to find  %s: %s\n", AVPToFind, err)
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

// getQuotaGrant gets how much credit to return in a CCA-update, which is the
// minimum between the max usage and how much credit is in the account
// Returns credits to return, true if these are the final bytes, true if we have exceeded the quota
// Depending on OCS configuration grantTypeProcedure it will use TOTAl bytes or TX bytes for calculations
func getQuotaGrant(srv *OCSDiamServer, bucket *CreditBucket) (*protos.Octets, bool, bool) {
	switch srv.ocsConfig.grantTypeProcedure {
	case protos.OCSConfig_TotalOnly:
		return getQuotaGrantOnlyTotal(srv, bucket)
	case protos.OCSConfig_TxOnly:
		return getQuotaGrantOnlyTX(srv, bucket)
	default:
		panic("getQuotaGrant type not implemented")
	}
}

// getQuotaGrantOnlyTotal gets how much credit to return in a CCA-update, which is the
// minimum between the max usage and how much credit is in the account
// Returns credits to return, true if these are the final bytes, true if we have exceeded the quota
func getQuotaGrantOnlyTotal(srv *OCSDiamServer, bucket *CreditBucket) (*protos.Octets, bool, bool) {
	var grant *protos.Octets
	var selectedMaxGrant uint64

	switch bucket.Unit {
	case protos.CreditInfo_Bytes:
		maxGrantedServiceUnits := srv.ocsConfig.MaxUsageOctets
		selectedMaxGrant = maxGrantedServiceUnits.GetTotalOctets()
		perRequest := bucket.Volume
		grant = &protos.Octets{
			TotalOctets:  getMin(maxGrantedServiceUnits.GetTotalOctets(), perRequest.GetTotalOctets()),
			InputOctets:  getMin(maxGrantedServiceUnits.GetInputOctets(), perRequest.GetInputOctets()),
			OutputOctets: getMin(maxGrantedServiceUnits.GetOutputOctets(), perRequest.GetOutputOctets())}

	case protos.CreditInfo_Time:
		selectedMaxGrant = uint64(srv.ocsConfig.MaxUsageTime)
		grant = &protos.Octets{TotalOctets: getMin(uint64(srv.ocsConfig.MaxUsageTime), bucket.Volume.GetTotalOctets())}
	}
	if grant.GetTotalOctets() <= selectedMaxGrant {
		return grant, true, false
	}
	if grant.GetTotalOctets() <= 0 {
		return grant, true, true
	}
	return grant, false, false
}

// getQuotaGrantOnlyTX does the same getQuotaGrantOnlyTotal but only check TX bytes (output Octets)
func getQuotaGrantOnlyTX(srv *OCSDiamServer, bucket *CreditBucket) (*protos.Octets, bool, bool) {
	var grant *protos.Octets
	var selectedMaxGrant uint64

	switch bucket.Unit {
	case protos.CreditInfo_Bytes:
		maxGrantedServiceUnits := srv.ocsConfig.MaxUsageOctets
		selectedMaxGrant = maxGrantedServiceUnits.GetOutputOctets()
		perRequest := bucket.Volume
		grant = &protos.Octets{
			TotalOctets:  getMin(maxGrantedServiceUnits.GetTotalOctets(), perRequest.GetTotalOctets()),
			InputOctets:  getMin(maxGrantedServiceUnits.GetInputOctets(), perRequest.GetInputOctets()),
			OutputOctets: getMin(maxGrantedServiceUnits.GetOutputOctets(), perRequest.GetOutputOctets())}

	case protos.CreditInfo_Time:
		selectedMaxGrant = uint64(srv.ocsConfig.MaxUsageTime)
		grant = &protos.Octets{TotalOctets: getMin(uint64(srv.ocsConfig.MaxUsageTime), bucket.Volume.GetOutputOctets())}
	}
	if grant.GetOutputOctets() <= selectedMaxGrant {
		return grant, true, false
	}
	if grant.GetOutputOctets() <= 0 {
		return grant, true, true
	}
	return grant, false, false
}

func getMin(first, second uint64) uint64 {
	if first > second {
		return second
	}
	return first
}

func toFinalUnitActionAVP(finalUnitAction protos.FinalUnitAction, redirectAddress string, restrict_rules []string) []*diam.AVP {
	fuaAVPs := []*diam.AVP{
		diam.NewAVP(avp.FinalUnitAction, avp.Mbit, 0, datatype.Enumerated(finalUnitAction)),
	}

	if finalUnitAction == protos.FinalUnitAction_Restrict {
		if len(restrict_rules) == 0 {
			glog.Errorf("RestrictRules must be provided when final unit action is set to restrict\n")
			return fuaAVPs
		}
		for _, rule := range restrict_rules {
			fuaAVPs = append(fuaAVPs,
				diam.NewAVP(avp.FilterID, avp.Mbit, 0, datatype.UTF8String(rule)),
			)
		}
	}

	if finalUnitAction == protos.FinalUnitAction_Redirect {
		fuaAVPs = append(
			fuaAVPs,
			diam.NewAVP(avp.RedirectServer, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: []*diam.AVP{
					diam.NewAVP(avp.RedirectServerAddress, avp.Mbit, 0, datatype.UTF8String(redirectAddress)),
					diam.NewAVP(avp.RedirectAddressType, avp.Mbit, 0, datatype.Unsigned32(0)),
				},
			}),
		)
	}
	return fuaAVPs
}

func toGrantedUnitsAVP(resultCode uint32, validityTime uint32, quotaGrant *protos.Octets, isFinalUnit bool, ratingGroup uint32, fuAction protos.FinalUnitAction, redirectAddr string, restrict_rules []string) *diam.AVP {
	creditGroup := &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.GrantedServiceUnit, avp.Mbit, 0, &diam.GroupedAVP{
				AVP: toGrantedServiceUnitAVP(quotaGrant),
			}),
			diam.NewAVP(avp.ValidityTime, avp.Mbit, 0, datatype.Unsigned32(validityTime)),
			diam.NewAVP(avp.RatingGroup, avp.Mbit, 0, datatype.Unsigned32(ratingGroup)),
			diam.NewAVP(avp.ResultCode, avp.Mbit, 0, datatype.Unsigned32(resultCode)),
		},
	}
	if isFinalUnit {
		creditGroup.AddAVP(
			diam.NewAVP(avp.FinalUnitIndication, avp.Mbit, 0, &diam.GroupedAVP{AVP: toFinalUnitActionAVP(fuAction, redirectAddr, restrict_rules)}),
		)
	}
	return diam.NewAVP(avp.MultipleServicesCreditControl, avp.Mbit, 0, creditGroup)
}

func toGrantedServiceUnitAVP(quotaGrant *protos.Octets) []*diam.AVP {
	res := []*diam.AVP{}
	if quotaGrant.GetTotalOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCTotalOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetTotalOctets())))
	}
	if quotaGrant.GetInputOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCInputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetInputOctets())))
	}
	if quotaGrant.GetOutputOctets() != 0 {
		res = append(res, diam.NewAVP(avp.CCOutputOctets, avp.Mbit, 0, datatype.Unsigned64(quotaGrant.GetOutputOctets())))
	}
	return res
}
