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

// Package client (eap.client) provides interface to supported & registered EAP Authenticator Providers
//
package client

import (
	"errors"
	"fmt"
	"strings"

	"github.com/golang/glog"

	"magma/feg/gateway/services/aaa/protos"
	"magma/feg/gateway/services/eap"
	"magma/feg/gateway/services/eap/providers"
	"magma/feg/gateway/services/eap/providers/registry"
)

var supportedEapTypes []uint8 // a read only copy of supported types for internal use

func init() {
	supportedEapTypes = SupportedTypes()
}

// HandleIdentityResponse passes Identity EAP payload to corresponding method provider & returns corresponding
// EAP result
// NOTE: Identity Request is handled by APs & does not involve EAP Authenticator's support
func HandleIdentityResponse(providerType uint8, msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Nil EAP Request")
	}
	err := verifyEapPayload(msg.Payload)
	if err != nil {
		return newFailureMsg(msg), err
	}
	if msg.Payload[eap.EapMsgMethodType] != eap.MethodIdentity {
		return newFailureMsg(msg), fmt.Errorf(
			"Invalid EAP Method Type for Identity Response: %d. Expecting EAP Identity (%d)",
			msg.Payload[eap.EapMsgMethodType], eap.MethodIdentity)
	}
	td := eap.Packet(msg.Payload).TypeData()
	if msg.Ctx != nil {
		if len(td) > 0 {
			msg.Ctx.Identity = string(td)
		}
	}
	var p providers.Method
	if providerType == 0 { // no EAP type specified, find a provider willing to handle the Identity
		if len(td) > 0 {
			for _, typ := range supportedEapTypes {
				pt := registry.GetProvider(typ)
				if pt != nil && pt.WillHandleIdentity(td) {
					p = pt
					break
				}
			}
		}
	} else {
		p = registry.GetProvider(providerType)
	}
	if p == nil && len(supportedEapTypes) > 0 {
		// we still have no viable provider, try the first one available (AKA)
		p = registry.GetProvider(supportedEapTypes[0])
		if p != nil {
			glog.Warningf("No EAP Provider found for Identity: '%s', trying %s", string(td), p.String())
		}
	}
	if p == nil {
		return newFailureMsg(msg), unsupportedProviderError(providerType)
	}
	glog.V(2).Infof("Handling %s (%d)", p.String(), providerType)
	return p.Handle(msg)
}

// SupportedTypes returns sorted list (ascending, by type) of registered EAP Providers
// SupportedTypes makes copy of an internally maintained supported types list, so callers
// are advised to save the result locally and re-use it if needed
func SupportedTypes() []uint8 {
	return registry.SupportedTypes()
}

// Handle handles passed EAP payload & returns corresponding EAP result
func Handle(msg *protos.Eap) (*protos.Eap, error) {
	if msg == nil {
		return nil, errors.New("Nil EAP Message")
	}
	err := verifyEapPayload(msg.Payload)
	if err != nil {
		return newFailureMsg(msg), err
	}
	method := msg.Payload[eap.EapMsgMethodType]
	var p providers.Method
	// Legacy Nak Based Auth Method Discovery
	// If the method is Nak, try to find a replacement handle (if specified by the Nak)
	if method == eap.MethodNak {
		// Get Nak's desired auth types array. If the peer did not provide desired auth types
		// or none of the types is supported - return EAP Failure
		td := eap.Packet(msg.Payload).TypeDataUnsafe()
		for _, method = range td {
			// Find first supported desired auth type
			p = registry.GetProvider(method)
			if p != nil {
				// a matching handler is found, call it with a simulated EAP Identity (1) Request,
				// use previously saved identity (if any) to create the request
				identity := []byte{eap.MethodIdentity}
				if msg.Ctx != nil {
					identity = append(identity, []byte(msg.Ctx.Identity)...)
				}
				// Create new EAP request with the simulated EAP Identity packet
				msg = &protos.Eap{
					Payload: eap.NewPacket(eap.ResponseCode, msg.Payload[eap.EapMsgIdentifier], identity),
					Ctx:     msg.Ctx,
				}
				break // first found provider
			}
		}
	} else {
		p = registry.GetProvider(method)
	}
	if p == nil {
		feap := newFailureMsg(msg)
		return feap, unsupportedProviderError(method)
	}
	return p.Handle(msg)
}

// newFailureMsg returns a new *protos.Eap with Payload set to EAP Failure packet
func newFailureMsg(msg *protos.Eap) *protos.Eap {
	var (
		ctx     *protos.Context
		payload eap.Packet
	)
	if msg != nil {
		ctx, payload = msg.Ctx, eap.Packet(msg.Payload)
	}
	return &protos.Eap{Payload: payload.Failure(), Ctx: ctx}
}

// verifyEapPayload checks validity of EAP message & it's length
func verifyEapPayload(payload []byte) error {
	el := len(payload)
	if el < eap.EapMsgData {
		return fmt.Errorf("EAP Message is too short: %d bytes", el)
	}
	mLen := uint16(payload[eap.EapMsgLenHigh])<<8 + uint16(payload[eap.EapMsgLenLow])
	if el < int(mLen) {
		return fmt.Errorf("Invalid EAP Message: bytes received %d are below specified length %d", el, mLen)
	}
	if payload[eap.EapMsgCode] != eap.CodeResponse {
		return fmt.Errorf(
			"Unsupported EAP Code: %d. Expecting EAP-Response (%d)",
			payload[eap.EapMsgCode], eap.CodeResponse)
	}
	return nil
}

func unsupportedProviderError(methodType uint8) error {
	return fmt.Errorf("Unsupported EAP Provider for Method Type: %d", methodType)
}

// BytesToStr returns Go compatible byte slice string
func BytesToStr(b []byte) string {
	return strings.Trim(strings.Replace(fmt.Sprint(b), " ", ", ", -1), "[]")
}
