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

package akamagma

import (
	"context"
	"encoding/json"
	"errors"

	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/eap/methods"
	"fbc/cwf/radius/modules/eap/methods/common"
	"fbc/cwf/radius/modules/eap/packet"
	aaa "fbc/cwf/radius/modules/protos"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// EapAkaMagmaMethod Implementation ofthe EAP-AKA method impl with Magma binding
type EapAkaMagmaMethod struct {
	config    Config
	akaClient aaa.AuthenticatorClient
}

// Config the aka-magma configuration
type Config struct {
	FegEndpoint string
}

// Create ...
func Create(config methods.MethodConfig) (methods.EapMethod, error) {
	// Parse config
	var akaConfig Config
	err := mapstructure.Decode(config, &akaConfig)
	if err != nil {
		return nil, errors.New("failed to parse AKAMAGMA configuration")
	}

	// Get EAP Authenticator GRPC client
	conn, err := grpc.Dial(akaConfig.FegEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return &EapAkaMagmaMethod{
		config:    akaConfig,
		akaClient: aaa.NewAuthenticatorClient(conn),
	}, nil
}

// Handle ...
func (m EapAkaMagmaMethod) Handle(
	c *modules.RequestContext,
	p *packet.Packet,
	s string,
	r *radius.Request,
) (*methods.HandlerResponse, error) {
	radiusReqAuthenticator := r.Authenticator[:]
	radiusSecret := r.Secret
	sessionID := c.SessionID
	eapPacket := p
	state := s
	eapLogger := c.Logger

	bytes, err := eapPacket.Bytes()
	if err != nil {
		return nil, err
	}

	// Get the client MAC address
	var clientMac, apn string
	attr, err := rfc2865.CallingStationID_Lookup(r.Packet)
	if err == nil {
		clientMac = string(attr)
	}
	attr, err = rfc2865.CalledStationID_Lookup(r.Packet)
	if err == nil {
		apn = string(attr)
	} else {
		eapLogger.Warn("Error Getting Called-Station_ID ", zap.Error(err))
	}

	UnmarshalProtocolState.Start()
	eapContext := aaa.Context{}
	if err := json.Unmarshal([]byte(state), &eapContext); err != nil {
		// This is not an invalid flow, but rather might happen when context has
		// not yet been registered (i.e: on first handshake message)
		UnmarshalProtocolState.Failure(err.Error())
		eapContext = aaa.Context{
			SessionId: sessionID,
			MacAddr:   clientMac,
			Apn:       apn,
		}
		eapLogger.Debug("EAP state not found, created a new state", zap.Any("state", eapContext))
	} else {
		eapContext.SessionId = sessionID // Always get the session id from RADIUS
		eapLogger.Debug("EAP state unmarshaled successfully", zap.Any("state", eapContext))
		UnmarshalProtocolState.Success()

		// Verify & warn if MAC address was already set on session but now changed
		if eapContext.MacAddr != "" && eapContext.MacAddr != clientMac {
			eapLogger.Warn(
				"Found incompatible MAC address on session",
				zap.String("previous", eapContext.MacAddr),
				zap.String("current", clientMac),
			)
		}
		if len(eapContext.Apn) == 0 {
			eapLogger.Warn("Empty Context APN,", zap.String("setting to Called-Station-Id", apn))
			eapContext.Apn = apn
		}
	}

	c.SessionStorage.Set(session.State{
		MACAddress:      clientMac,
		MSISDN:          eapContext.GetMsisdn(),
		CalledStationID: apn,
	})

	var eapResponse *aaa.Eap
	if eapPacket.EAPType == packet.EAPTypeIDENTITY {
		c.Logger.Debug("Handling EAP-Identity request")
		eapResponse, err = m.akaClient.HandleIdentity(
			context.Background(),
			&aaa.EapIdentity{
				Payload: bytes,
				Ctx:     &eapContext,
				Method:  0, // pass undefined method & let EAP router to find an EAP provider to handle the request
			},
		)
	} else {
		c.Logger.Debug("Handling EAP-non-Identity request")
		eapResponse, err = m.akaClient.Handle(
			context.Background(),
			&aaa.Eap{
				Payload: bytes,
				Ctx:     &eapContext,
			},
		)
	}
	if err != nil {
		eapLogger.Error("Failed handling EAP message", zap.Error(err))
		return nil, err
	}

	// Marshal protocol new state
	MarshalProtocolState.Start()
	postHandlerContext := eapResponse.GetCtx()
	newProtocolState, err := json.Marshal(postHandlerContext)
	if err != nil {
		// We mark this as error, but this is allowed (for example: in case of new auth session)
		MarshalProtocolState.Failure(err.Error())
		newProtocolState = []byte("{}")
	} else {
		eapLogger.Debug("EAP state marshaled successfully", zap.Any("state", postHandlerContext))
		MarshalProtocolState.Success()
	}

	// Build the returned packet
	eapResponsePacket, err := packet.NewPacketFromRaw(eapResponse.GetPayload())
	if err != nil {
		return nil, err
	}

	result := &methods.HandlerResponse{
		Packet:           eapResponsePacket,
		RadiusCode:       methods.ToRadiusCode(eapResponsePacket.Code),
		NewProtocolState: string(newProtocolState),
	}

	// Add key material for Access-Accept/EAP-Success message
	if eapResponsePacket.Code == packet.CodeSUCCESS {
		result.ExtraAttributes = radius.Attributes{}

		// Add MPPE keys
		keyingMaterialAttrs, err := common.GetKeyingAttributes(
			postHandlerContext.GetMsk(),
			radiusSecret,
			radiusReqAuthenticator,
		)
		if err != nil {
			return nil, err
		}
		result.ExtraAttributes[rfc2865.VendorSpecific_Type] = keyingMaterialAttrs

		// Add User-Name attribute, which is mandatory
		result.ExtraAttributes[rfc2865.UserName_Type] =
			[]radius.Attribute{
				radius.Attribute(postHandlerContext.Identity),
			}
	}
	return result, nil
}
