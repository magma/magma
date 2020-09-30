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

package eap

import (
	"errors"
	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/eap/authstate"
	"fbc/cwf/radius/modules/eap/methods"
	"fbc/cwf/radius/modules/eap/methods/akamagma"
	"fbc/cwf/radius/modules/eap/methods/akatataipx"
	"fbc/cwf/radius/modules/eap/packet"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2869"
	"fmt"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
)

// Method A definition for an EAP method with its config
type Method struct {
	Name   string               `json:"name"`
	Config methods.MethodConfig `json:"config"`
}

// Config configuration structure for the EAP module
type Config struct {
	Methods []Method
}

// stateManager a state manage instance
type ModuleCtx struct {
	stateManager authstate.Manager
	method       methods.EapMethod
}

// Init module interface implementation
//nolint:deadcode
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var mCtx ModuleCtx

	// Parse config
	var eapConfig Config
	err := mapstructure.Decode(config, &eapConfig)
	if err != nil {
		return nil, err
	}

	// Initialize State Manager singleton
	// TODO: sync object
	if mCtx.stateManager == nil {
		mCtx.stateManager = authstate.NewMemoryManager()
	}

	// TODO: handle multiple methods (currently assuming only one)
	mCtx.method, err = getMethod(eapConfig.Methods[0])
	if err != nil {
		return nil, err
	}

	// We're done without any error!
	return mCtx, nil
}

// GetMethod factory method, instatiates and initializes an EAP method
func getMethod(method Method) (methods.EapMethod, error) {
	switch method.Name {
	case "akamagma":
		return akamagma.Create(method.Config)
	case "akatataipx":
		return akatataipx.Create(method.Config)
	default:
		return nil, errors.New("unsupported eap method '%s' ('akamagma', 'akatataipx' are supported")
	}
}

// Handle module interface implementation
//nolint:deadcode
func Handle(m modules.Context, c *modules.RequestContext, r *radius.Request, next modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)
	c.Logger.Debug("Starting to handle radius request")

	// Extract EAP packet
	ExtractEapPacket.Start()
	eapPacket, err := packet.NewPacketFromRadius(r.Packet)
	if err != nil {
		c.Logger.Error("Failed to extract EAP packet", zap.Error(err))
		ExtractEapPacket.Failure("missing_or_invalid_eap_packet")
		return next(c, r)
	}
	ExtractEapPacket.Success()

	// Build EAP logger for the current request
	eapLogger := c.Logger.
		With(zap.Int64("correlation_id", c.RequestID)).
		With(zap.Int("eap_type", int(eapPacket.EAPType))).
		With(zap.Int("eap_code", int(eapPacket.Code)))

	// Restore EAP authentication state (reset if we got Identity Response)
	RestoreProtocolState.Start()
	eapAuthState := &authstate.Container{}
	if eapPacket.EAPType == packet.EAPTypeIDENTITY {
		err := mCtx.stateManager.Reset(r.Packet, packet.EAPTypeIDENTITY)
		if err != nil {
			c.Logger.Error("Failed to load EAP state", zap.Error(err))
			RestoreProtocolState.Failure("reset_on_eap_identity_failed")
			return next(c, r)
		}
		if err := mCtx.stateManager.Set(r.Packet, packet.EAPTypeIDENTITY, authstate.Container{}); err != nil {
			c.Logger.Error("Failed to load EAP state", zap.Error(err))
			RestoreProtocolState.Failure("set_empty_on_eap_identity")
			return next(c, r)
		}
	} else {
		eapAuthState, err = mCtx.stateManager.Get(r.Packet, eapPacket.EAPType)
		if err != nil {
			c.Logger.Error("Missing or invalid EAP auth state", zap.Error(err))
			RestoreProtocolState.Failure("missing_or_invalid_auth_state")
			return next(c, r)
		}
	}
	RestoreProtocolState.Success()

	HandleEapPacket.Start()
	// Check if EAP method is supported
	if eapPacket.EAPType != packet.EAPTypeAKA && eapPacket.EAPType != packet.EAPTypeSIM &&
		eapPacket.EAPType != packet.EAPTypeIDENTITY {
		c.Logger.Error("Unsupported EAP method requested", zap.Int("eap_method", int(eapPacket.EAPType)))
		HandleEapPacket.Failure(fmt.Sprintf("unsupported_eap_type_%d", int(eapPacket.EAPType)))
	}

	// Handle the EAP-method state machine
	logger := c.Logger
	c.Logger = eapLogger
	eapResponse, err := mCtx.method.Handle(c, eapPacket, eapAuthState.ProtocolState, r)
	if err != nil {
		c.Logger.Error("Failed handling EAP packet", zap.Error(err))
		HandleEapPacket.Failure("unknown")
		return nil, err
	}
	c.Logger = logger
	HandleEapPacket.Success()

	// Persist state
	PersistProtocolState.Start()
	eapAuthState.ProtocolState = eapResponse.NewProtocolState
	err = mCtx.stateManager.Set(r.Packet, eapPacket.EAPType, *eapAuthState)
	if err != nil {
		PersistProtocolState.Failure("unknown")
		c.Logger.Error("Failed to persist state", zap.Error(err))
		return nil, err
	}
	PersistProtocolState.Success()

	// Add EAP Packet to the EAP-Message AVP
	radiusResponse := r.Response(eapResponse.RadiusCode)
	if eapResponse.Packet != nil {
		eapBytes, err := eapResponse.Packet.Bytes()
		if err != nil {
			c.Logger.Error("Failed serializing EAP response", zap.Error(err))
			return nil, err
		}
		radiusResponse.Add(rfc2869.EAPMessage_Type, eapBytes)
	}

	// Add the extra attributes to the radius packet
	if eapResponse.ExtraAttributes != nil {
		for t, attrs := range eapResponse.ExtraAttributes {
			for _, attr := range attrs {
				radiusResponse.Add(t, attr)
			}
		}
	}

	return &modules.Response{
		Code:       eapResponse.RadiusCode,
		Attributes: radiusResponse.Attributes,
	}, nil
}
