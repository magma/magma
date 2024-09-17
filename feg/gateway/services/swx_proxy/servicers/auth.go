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

// Package servicers implements Swx GRPC proxy service which sends MAR/SAR messages over
// diameter connection, waits (blocks) for diameter's MAA/SAAs and returns their RPC representation
package servicers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/fiorix/go-diameter/v4/diam"
	"github.com/fiorix/go-diameter/v4/diam/avp"
	"github.com/fiorix/go-diameter/v4/diam/datatype"
	"github.com/fiorix/go-diameter/v4/diam/dict"
	"github.com/golang/glog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"magma/feg/cloud/go/protos"
	"magma/feg/gateway/diameter"
	"magma/feg/gateway/services/swx_proxy/metrics"
)

const (
	MinRequestedVectors uint32 = 5
	MaxReturnedVectors  int    = 5
)

// AuthenticateImpl sends MAR over diameter connection,
// waits (blocks) for MAA & returns its RPC representation
func (s *swxProxy) AuthenticateImpl(req *protos.AuthenticationRequest) (*protos.AuthenticationAnswer, error) {
	var (
		res       = &protos.AuthenticationAnswer{}
		err       = validateAuthRequest(req)
		cachedRes *protos.AuthenticationAnswer
	)
	if err != nil {
		return res, status.Errorf(codes.InvalidArgument, err.Error())
	}
	res.UserName = req.GetUserName()
	shouldSendSar := s.config.VerifyAuthorization || req.GetRetrieveUserProfile()
	requestedVectors := int(req.SipNumAuthVectors)
	if requestedVectors < 1 {
		requestedVectors = 1
	} else if requestedVectors > MaxReturnedVectors {
		requestedVectors = MaxReturnedVectors
	}
	if s.cache != nil {
		// Check if we still have valid vectors for the user in the cache
		if len(req.GetResyncInfo()) == 0 { // Only try to get cached vectors if it's not resync request
			cachedRes = s.cache.Get(res.UserName, requestedVectors)
			if cachedRes != nil && ((!shouldSendSar) || cachedRes.GetUserProfile() != nil) &&
				len(cachedRes.SipAuthVectors) >= requestedVectors {
				return cachedRes, nil // We have a valid result in the cache, return it
			}
		}

		if req.SipNumAuthVectors < MinRequestedVectors {
			req.SipNumAuthVectors = MinRequestedVectors // Get Max allowed # of vectors for caching
		}
	}
	sid := s.genSID(req.GetUserName())
	maa, err := s.sendMAR(req, sid)
	if err != nil {
		if protos.SwxErrorCode(status.Code(err)) != protos.SwxErrorCode_IDENTITY_ALREADY_REGISTERED {
			if cachedRes != nil {
				// if we have a cached result with insufficient # of vectors, log an error & return it
				glog.Errorf("SWx send MAR Error: %v, returning %d cached vectors", err, len(cachedRes.SipAuthVectors))
				return cachedRes, nil
			}
			return res, err
		}
		aaaHost := string(maa.AAAServerName)
		if len(aaaHost) > 0 {
			originRalm := s.config.ClientCfg.Realm
			if s.config.DeriveUnregisterRealm {
				ha := strings.SplitN(aaaHost, ".", 2)
				if len(ha) == 2 && len(ha[1]) > 0 {
					originRalm = ha[1]
				} else {
					glog.Errorf("Cannot derive Origin-Realm from AAA-Server-Name: %s", aaaHost)
				}
			}
			// deregister
			s.sendSARExt(
				req.GetUserName(),
				ServerAssignnmentType_USER_DEREGISTRATION,
				aaaHost,
				originRalm, "")
			// repeat MAR after deregistration
			sid = s.genSID(req.GetUserName())
			maa, err = s.sendMAR(req, sid)
		}
		if err != nil {
			return res, err
		}
	}
	res.SessionId = sid
	if shouldSendSar {
		profile, authorized, err := s.retrieveUserProfile(req.GetUserName(), sid)
		if err != nil {
			glog.Error(err)
		}
		// If user is unauthorized, don't send back auth vectors
		if !authorized {
			return res, err
		}
		res.UserProfile = profile
	} else if s.config.RegisterOnAuth {
		err := s.registerUser(req.GetUserName(), sid)
		if err != nil {
			glog.Error(err)
		}
	}
	res.SipAuthVectors = getSIPAuthenticationVectors(maa.SIPAuthDataItems)
	// The only point when we cache vectors
	if s.cache != nil {
		if cachedRes != nil {
			cacheVectors := len(cachedRes.SipAuthVectors)
			requestedVectors -= cacheVectors
			res = s.cache.Put(res, requestedVectors)
			res.SipAuthVectors = append(cachedRes.SipAuthVectors, res.SipAuthVectors...)
		} else {
			res = s.cache.Put(res, requestedVectors)
		}
	}
	return res, err
}

func (s *swxProxy) sendMAR(req *protos.AuthenticationRequest, sid string) (*MAA, error) {
	if len(sid) == 0 {
		sid = s.genSID(req.GetUserName())
	}
	ch := make(chan interface{})
	s.requestTracker.RegisterRequest(sid, ch)
	// if request hasn't been removed by end of transaction, remove it
	defer s.requestTracker.DeregisterRequest(sid)

	marMsg, err := s.createMAR(sid, req)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	marStartTime := time.Now()
	err = s.sendDiameterMsg(marMsg, MAX_DIAM_RETRIES)
	if err != nil {
		metrics.MARSendFailures.Inc()
		err = status.Errorf(codes.Internal, "Error while sending MAR with SID %s: %s", sid, err)
		glog.Error(err)
		return nil, err
	}
	metrics.MARRequests.Inc()
	select {
	case resp, open := <-ch:
		metrics.MARLatency.Observe(time.Since(marStartTime).Seconds())
		if !open {
			metrics.SwxInvalidSessions.Inc()
			err = status.Errorf(codes.Aborted, "MAA for Session ID: %s is canceled", sid)
			glog.Error(err)
			return nil, err
		}
		maa, ok := resp.(*MAA)
		if !ok {
			metrics.SwxUnparseableMsg.Inc()
			err = status.Errorf(codes.Internal, "Invalid Response Type: %T, MAA expected.", resp)
			glog.Error(err)
			return nil, err
		}
		err = diameter.TranslateDiamResultCode(maa.ResultCode)
		metrics.SwxResultCodes.WithLabelValues(strconv.FormatUint(uint64(maa.ResultCode), 10)).Inc()
		// If there is no base diameter error, check that there is no experimental error either
		if err == nil {
			err = diameter.TranslateDiamResultCode(maa.ExperimentalResult.ExperimentalResultCode)
			metrics.SwxExperimentalResultCodes.WithLabelValues(strconv.FormatUint(uint64(maa.ExperimentalResult.ExperimentalResultCode), 10)).Inc()
		}
		// According to spec 29.273, SIP-Auth-Data-Item(s) only present on SUCCESS
		return maa, err

	case <-time.After(time.Second * TIMEOUT_SECONDS):
		metrics.MARLatency.Observe(time.Since(marStartTime).Seconds())
		metrics.SwxTimeouts.Inc()
		err = status.Errorf(codes.DeadlineExceeded, "MAA Timed Out for Session ID: %s", sid)
		glog.Error(err)
		return nil, err
	}
}

// retrieveUserProfile sends SARs with ServerAssignmentType AAA_USER_DATA_REQUEST or REGISTRATION, receives back SAA
// and returns the subscribers's Non-3GPP-User-Data profile
func (s *swxProxy) retrieveUserProfile(userName, sid string) (*protos.AuthenticationAnswer_UserProfile, bool, error) {
	var sat uint32 = ServerAssignmentType_AAA_USER_DATA_REQUEST
	if s.config.RegisterOnAuth {
		sat = ServerAssignmentType_REGISTRATION
	}
	saa, err := s.sendSAR(userName, sat, sid)
	if err != nil {
		return nil, true, err
	}
	if s.config.VerifyAuthorization && saa.UserData.Non3GPPIPAccess != datatype.Enumerated(Non3GPPIPAccess_ENABLED) {
		metrics.UnauthorizedAuthAttempts.Inc()
		return nil, false, status.Errorf(codes.PermissionDenied, "User %s is not authorized for Non-3GPP Subscription Access", userName)
	}
	if saa.UserData.SubscriptionId.SubscriptionIdType != END_USER_E164 {
		return nil, true, status.Error(
			codes.Internal,
			"Subscription ID type is not END_USER_E164; Cannot retrieve MSISDN",
		)
	}
	userProfile := &protos.AuthenticationAnswer_UserProfile{
		Msisdn: string(saa.UserData.SubscriptionId.SubscriptionIdData),
	}
	return userProfile, true, nil
}

// registerUser sends SARs with ServerAssignmentType REGISTRATION
func (s *swxProxy) registerUser(userName, sid string) error {
	_, err := s.sendSAR(userName, ServerAssignmentType_REGISTRATION, sid)
	return err
}

// createMAR creates a Multimedia Authentication Request diameter msg with provided SessionID (sid)
// to be sent to HSS
func (s *swxProxy) createMAR(sid string, req *protos.AuthenticationRequest) (*diam.Message, error) {
	authScheme, err := convertAuthSchemeToString(req.GetAuthenticationScheme())
	if err != nil {
		return nil, err
	}

	msg := diameter.NewProxiableRequest(diam.MultimediaAuthentication, diam.TGPP_SWX_APP_ID, dict.Default)
	msg.NewAVP(avp.SessionID, avp.Mbit, 0, datatype.UTF8String(sid))
	msg.NewAVP(avp.VendorSpecificApplicationID, avp.Mbit, 0, &diam.GroupedAVP{
		AVP: []*diam.AVP{
			diam.NewAVP(avp.AuthApplicationID, avp.Mbit, 0, datatype.Unsigned32(diam.TGPP_SWX_APP_ID)),
			diam.NewAVP(avp.VendorID, avp.Mbit, 0, datatype.Unsigned32(diameter.Vendor3GPP)),
		},
	})
	msg.NewAVP(avp.OriginHost, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Host))
	msg.NewAVP(avp.OriginRealm, avp.Mbit, 0, datatype.DiameterIdentity(s.config.ClientCfg.Realm))
	msg.NewAVP(avp.UserName, avp.Mbit, 0, datatype.UTF8String(req.GetUserName()))
	msg.NewAVP(avp.AuthSessionState, avp.Mbit, 0, datatype.Enumerated(1))
	msg.NewAVP(avp.SIPNumberAuthItems, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Unsigned32(req.GetSipNumAuthVectors()))
	msg.NewAVP(avp.RATType, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.Enumerated(RadioAccessTechnologyType_WLAN))
	authDataAvp := []*diam.AVP{
		diam.NewAVP(avp.SIPAuthenticationScheme, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.UTF8String(authScheme)),
	}
	if len(req.GetResyncInfo()) > 0 {
		authDataAvp = append(
			authDataAvp,
			diam.NewAVP(avp.SIPAuthorization, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, datatype.OctetString(req.GetResyncInfo())),
		)
	}
	msg.NewAVP(avp.SIPAuthDataItem, avp.Mbit|avp.Vbit, diameter.Vendor3GPP, &diam.GroupedAVP{AVP: authDataAvp})
	return msg, nil
}

func handleMAA(s *swxProxy) diam.HandlerFunc {
	return func(c diam.Conn, m *diam.Message) {
		var maa MAA
		err := m.Unmarshal(&maa)
		if err != nil {
			metrics.SwxUnparseableMsg.Inc()
			glog.Errorf("MAA Unmarshal failed for remote %s & message %s: %s", c.RemoteAddr(), m, err)
			return
		}
		ch := s.requestTracker.DeregisterRequest(maa.SessionID)
		if ch != nil {
			ch <- &maa
		} else {
			metrics.SwxInvalidSessions.Inc()
			glog.Errorf("MAA SessionID %s not found. Message: %s, Remote: %s", maa.SessionID, m, c.RemoteAddr())
		}
	}
}

func getSIPAuthenticationVectors(items []SIPAuthDataItem) []*protos.AuthenticationAnswer_SIPAuthVector {
	var authVectors []*protos.AuthenticationAnswer_SIPAuthVector
	for _, item := range items {
		// If the auth scheme is unrecognized, don't include the vector
		authScheme, err := convertStringToAuthScheme(item.AuthScheme)
		if err != nil {
			glog.Error(err)
			continue
		}
		authVectors = append(
			authVectors,
			&protos.AuthenticationAnswer_SIPAuthVector{
				AuthenticationScheme: authScheme,
				RandAutn:             item.Authenticate.Serialize(),
				Xres:                 item.Authorization.Serialize(),
				ConfidentialityKey:   item.ConfidentialityKey.Serialize(),
				IntegrityKey:         item.IntegrityKey.Serialize()})
	}
	return authVectors
}

func validateAuthRequest(req *protos.AuthenticationRequest) error {
	if req == nil {
		return fmt.Errorf("nil authentication request provided")
	}
	if len(req.GetUserName()) == 0 {
		return fmt.Errorf("empty user-name provided in authentication request")
	}
	if req.SipNumAuthVectors == 0 {
		return fmt.Errorf("SIPNumAuthVectors in authentication request must be greater than 0")
	}
	// imsi cannot be greater than 15 digits according to 3GPP Spec 23.003
	if len(req.GetUserName()) > 15 {
		return fmt.Errorf("provided username %s is greater than 15 digits", req.GetUserName())
	}
	return nil
}

func convertStringToAuthScheme(maaScheme string) (protos.AuthenticationScheme, error) {
	switch maaScheme {
	case SipAuthScheme_EAP_AKA:
		return protos.AuthenticationScheme_EAP_AKA, nil
	case SipAuthScheme_EAP_AKA_PRIME:
		return protos.AuthenticationScheme_EAP_AKA_PRIME, nil
	default:
		return protos.AuthenticationScheme_EAP_AKA, fmt.Errorf("unrecognized Authentication Scheme returned: %s", maaScheme)
	}
}

func convertAuthSchemeToString(scheme protos.AuthenticationScheme) (string, error) {
	switch scheme {
	case protos.AuthenticationScheme_EAP_AKA:
		return SipAuthScheme_EAP_AKA, nil
	case protos.AuthenticationScheme_EAP_AKA_PRIME:
		return SipAuthScheme_EAP_AKA_PRIME, nil
	default:
		return "", fmt.Errorf("unrecognized Authentication Scheme returned: %v", scheme)
	}
}
