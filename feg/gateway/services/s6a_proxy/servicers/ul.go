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

// Package service implements S6a GRPC proxy service which sends AIR, ULR messages over diameter connection,
// waits (blocks) for diameter's AIAs, ULAs & returns their RPC representation
package servicers

import (
	"log"
	"strconv"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/s6a_proxy/metrics"
)

// sendULR - sends ULR with given Session ID (sid)
func (s *s6aProxy) sendULR(sid string, req *protos.UpdateLocationRequest, retryCount uint) error {
	c, err := s.connMan.GetConnection(s.smClient, s.config.ServerCfg)
	if err != nil {
		return err
	}
	m := diameter.NewProxiableRequest(diam.UpdateLocation, diam.TGPP_S6A_APP_ID, dict.Default)
	m.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	s.addDiamOriginAVPs(m)
	m.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.UserName))
	m.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	m.NewAVP(avp.RATType, avp.Mbit, diameter.Vendor3GPP, datatype.Enumerated(ULR_RAT_TYPE))
	m.NewAVP(avp.ULRFlags, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.Unsigned32(createULR_Flags(req)))
	m.NewAVP(avp.VisitedPLMNID, avp.Vbit|avp.Mbit, diameter.Vendor3GPP, datatype.OctetString(req.VisitedPlmn))

	// Supported feature-List-id-1
	if req.FeatureListId_1 != nil {
		m.NewAVP(avp.SupportedFeatures, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
				diam.NewAVP(avp.FeatureListID, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(1)),
				diam.NewAVP(avp.FeatureList, avp.Vbit, diameter.Vendor3GPP,
					datatype.Unsigned32(createFeatureListID1(req.FeatureListId_1))),
			},
		})
	}

	// Supported feature-List-id-2
	if req.FeatureListId_2 != nil {
		m.NewAVP(avp.SupportedFeatures, avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{
			AVP: []*diam.AVP{
				diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(10415)),
				diam.NewAVP(avp.FeatureListID, avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(2)),
				diam.NewAVP(avp.FeatureList, avp.Vbit, diameter.Vendor3GPP,
					datatype.Unsigned32(createFeatureListID2(req.FeatureListId_2))),
			},
		})
	}

	glog.V(2).Infof("Sending S6a ULR message\n%s\n", m)
	err = c.SendRequest(m, retryCount)
	if err != nil {
		err = Error(codes.DataLoss, err)
	}
	return err
}

// createULR_Flags creates ULR Flags based on TS 29.272
func createULR_Flags(req *protos.UpdateLocationRequest) int {
	// ULR_FLAGS contains defaults flags for ULR-Flags
	ulrFlags := ULR_FLAGS
	switch {
	case req.DualRegistration_5GIndicator:
		// add TS 29.272 Dual-Registration5G-Indicator on bit 8
		ulrFlags = ulrFlags | FlagBit8
	}
	return int(ulrFlags)
}

// createFeatureListID1 creates the Feature-List-ID 1 based on TS 29.272
func createFeatureListID1(featureList *protos.FeatureListId1) int {
	if featureList == nil {
		return 0
	}
	res := EmptyFlagBit
	switch {
	case featureList.RegionalSubscription:
		// set bit 9 (TS 29.272 Regional Subscription)
		res = res | FlagBit9
	}
	return int(res)
}

// createFeatureListID2 creates the Feature-List-ID 2 based on TS 29.272
func createFeatureListID2(featureList *protos.FeatureListId2) int {
	if featureList == nil {
		return 0
	}
	res := EmptyFlagBit
	switch {
	case featureList.NrAsSecondaryRat:
		// set bit 27 (TS 29.272 Nr As Secondary Rat)
		res = res | FlagBit27
	}
	return int(res)
}

// S6a ULA
func handleULA(s *s6aProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		glog.V(2).Infof("Received S6a ULA message:\n%s\n", m)
		var ula ULA
		err := m.Unmarshal(&ula)
		if err != nil {
			log.Printf("ULA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(ula.SessionID)
		if ch != nil {
			ch <- &ula
		} else {
			log.Printf("ULA SessionID %s not found. Message: %s, Remote: %s", ula.SessionID, m, c.RemoteAddr())
		}
	}
}

// UpdateLocationImpl sends ULR (Code 316) over diameter connection,
// waits (blocks) for ULA & returns its RPC representation
func (s *s6aProxy) UpdateLocationImpl(req *protos.UpdateLocationRequest) (*protos.UpdateLocationAnswer, error,
) {
	res := &protos.UpdateLocationAnswer{}
	if req == nil {
		return res, Errorf(codes.InvalidArgument, "Nil UL Request")
	}
	sid := s.genSID()
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	defer s.requestTracker.DeregisterRequest(sid)

	var (
		err     error
		retries uint = MAX_DIAM_RETRIES
	)

	err = s.sendULR(sid, req, retries)

	if err != nil {
		metrics.ULRSendFailures.Inc()
		log.Printf("Error sending ULR with SID %s: %v", sid, err)
	}

	if err == nil {
		metrics.ULRRequests.Inc()
		select {
		case resp, open := <-ch:
			if open {
				ula, ok := resp.(*ULA)
				if ok {
					metrics.S6aResultCodes.WithLabelValues(strconv.FormatUint(uint64(ula.ResultCode), 10)).Inc()
					err = diameter.TranslateDiamResultCode(ula.ResultCode)
					res.ErrorCode = protos.ErrorCode(ula.ExperimentalResult.ExperimentalResultCode)
					res.Msisdn = ula.SubscriptionData.MSISDN.Serialize()
					res.DefaultContextId = ula.SubscriptionData.APNConfigurationProfile.ContextIdentifier
					res.TotalAmbr = ula.SubscriptionData.AMBR.getProtoAmbr()
					res.DefaultChargingCharacteristics = ula.SubscriptionData.TgppChargingCharacteristics
					res.AllApnsIncluded =
						ula.SubscriptionData.APNConfigurationProfile.AllAPNConfigurationsIncludedIndicator == 0
					res.NetworkAccessMode = protos.UpdateLocationAnswer_NetworkAccessMode(ula.SubscriptionData.NetworkAccessMode)
					res.RegionalSubscriptionZoneCode = make([][]byte, len(ula.SubscriptionData.RegionalSubscriptionZoneCode))
					res.FeatureListId_1 = getFeatureListID1(ula.SupportedFeatures)
					res.FeatureListId_2 = getFeatureListID2(ula.SupportedFeatures)
					for i, code := range ula.SubscriptionData.RegionalSubscriptionZoneCode {
						res.RegionalSubscriptionZoneCode[i] = code.Serialize()
					}
					for _, apnCfg := range ula.SubscriptionData.APNConfigurationProfile.APNConfigs {
						res.Apn = append(
							res.Apn,
							&protos.UpdateLocationAnswer_APNConfiguration{
								ContextId:        apnCfg.ContextIdentifier,
								Pdn:              protos.UpdateLocationAnswer_APNConfiguration_PDNType(apnCfg.PDNType),
								ServiceSelection: apnCfg.ServiceSelection,
								QosProfile: &protos.UpdateLocationAnswer_APNConfiguration_QoSProfile{
									ClassId:                 apnCfg.EPSSubscribedQoSProfile.QoSClassIdentifier,
									PriorityLevel:           apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PriorityLevel,
									PreemptionCapability:    apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PreemptionCapability == 0,
									PreemptionVulnerability: apnCfg.EPSSubscribedQoSProfile.AllocationRetentionPriority.PreemptionVulnerability == 0,
								},
								Ambr:                    apnCfg.AMBR.getProtoAmbr(),
								ChargingCharacteristics: apnCfg.TgppChargingCharacteristics,
							})
					}
					return res, err
				} else {
					err = Errorf(codes.Internal, "Invalid Response Type: %T, ULA expected.", resp)
					metrics.S6aUnparseableMsg.Inc()
				}
			} else {
				err = Errorf(codes.Aborted, "ULR for Session ID: %s is canceled", sid)
			}
		case <-time.After(time.Second * TIMEOUT_SECONDS):
			err = Errorf(codes.DeadlineExceeded, "ULR Timed Out for Session ID: %s", sid)
			metrics.S6aTimeouts.Inc()
		}
	}
	return res, err
}

// getFeatureListID1 creates the Feature-List-ID 1 based on TS 29.272
func getFeatureListID1(supportedFeatures []SupportedFeatures) *protos.FeatureListId1 {
	if len(supportedFeatures) == 0 {
		return nil
	}
	protoFeatureList := &protos.FeatureListId1{}
	for _, features := range supportedFeatures {
		if features.FeatureListID == 1 {
			// get bit 27 (TS 29.272 Nr As Secondary Rat)
			if features.FeatureList&(1<<9) != 0 {
				protoFeatureList.RegionalSubscription = true
			}
		}
	}
	return protoFeatureList
}

// getFeatureListID2 creates the Feature-List-ID 2 based on TS 29.272
func getFeatureListID2(supportedFeatures []SupportedFeatures) *protos.FeatureListId2 {
	if len(supportedFeatures) == 0 {
		return nil
	}
	protoFeatureList := &protos.FeatureListId2{}
	for _, features := range supportedFeatures {
		if features.FeatureListID == 2 {
			// get bit 27 (TS 29.272 Nr As Secondary Rat)
			if features.FeatureList&(1<<27) != 0 {
				protoFeatureList.NrAsSecondaryRat = true
			}
		}
	}
	return protoFeatureList
}

func (ambr *AMBR) getProtoAmbr() *protos.UpdateLocationAnswer_AggregatedMaximumBitrate {
	if ambr.ExtendMaxRequestedBwDL != 0 && ambr.ExtendMaxRequestedBwUL != 0 {
		return &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
			MaxBandwidthUl: ambr.ExtendMaxRequestedBwUL,
			MaxBandwidthDl: ambr.ExtendMaxRequestedBwDL,
			Unit:           protos.UpdateLocationAnswer_AggregatedMaximumBitrate_KBPS,
		}
	}
	return &protos.UpdateLocationAnswer_AggregatedMaximumBitrate{
		MaxBandwidthUl: ambr.MaxRequestedBandwidthUL,
		MaxBandwidthDl: ambr.MaxRequestedBandwidthDL,
		Unit:           protos.UpdateLocationAnswer_AggregatedMaximumBitrate_BPS,
	}
}
