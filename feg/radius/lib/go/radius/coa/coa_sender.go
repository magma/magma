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

package coa

import (
	"context"
	"fmt"
	"time"

	"fbc/lib/go/radius/dictionaries/ruckus"
	"fbc/lib/go/radius/dictionaries/xwf"

	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"
	"fbc/lib/go/radius/rfc2869"
)

// AuthorizeTrafficClasses key value map represent the bytes left for each class
type AuthorizeTrafficClasses struct {
	AuthorizeClassName string
	AuthorizeBytesLeft uint64
}

// Client object that handle packet send
type Client radius.Client

// Packet - encoded bytes data to send forward
type Packet *radius.Packet

// Vendor represents a vendor
type Vendor int

const (
	// VendorXWFCertified Cambium, Mojo, SoMA, or other XWF certified
	VendorXWFCertified Vendor = iota
	// VendorRuckus Ruckus
	VendorRuckus
)

// Params params needed in order to create CoA request
type Params struct {
	NASIdentifier       string
	AcctInterimInterval uint32
	AcctSessionID       string
	CallingStationID    string
	CaptivePortalToken  string
	TrafficClasses      []AuthorizeTrafficClasses
	VendorName          Vendor
}

// DisconnectParams params needed in order to create disconnect request
type DisconnectParams struct {
	NASIdentifier    string
	AcctSessionID    string
	CallingStationID string
}

// Code - return by send function to inform the user about the command status
type Code int

const (
	// CodeUnknown - request failed or if the response is unknown
	CodeUnknown Code = iota
	// CodeACK - request succeeded
	CodeACK
	// CodeNAK - request failed (probably session is disconnected)
	CodeNAK
)

func trafficClassesToXWFCertified(trafficClasses []AuthorizeTrafficClasses) []xwf.XWFAuthorizeTrafficClasses {
	res := make([]xwf.XWFAuthorizeTrafficClasses, len(trafficClasses))
	for i, val := range trafficClasses {
		res[i] = xwf.XWFAuthorizeTrafficClasses{
			XWFAuthorizeClassName: val.AuthorizeClassName,
			XWFAuthorizeBytesLeft: val.AuthorizeBytesLeft,
		}
	}
	return res
}

func trafficClassesToRuckus(trafficClasses []AuthorizeTrafficClasses) []ruckus.RuckusTCAttrIdsWithQuota {
	res := make([]ruckus.RuckusTCAttrIdsWithQuota, len(trafficClasses))
	for i, val := range trafficClasses {
		res[i] = ruckus.RuckusTCAttrIdsWithQuota{
			RuckusTCNameQuota: val.AuthorizeClassName,
			RuckusTCQuota:     val.AuthorizeBytesLeft,
		}
	}
	return res
}

func updateCoARequestXWFCertified(params Params, p *radius.Packet) (Packet, error) {
	err := xwf.XWFAuthorizeTrafficClasses_Set(p, trafficClassesToXWFCertified(params.TrafficClasses))
	if err != nil {
		return nil, err
	}
	err = xwf.XWFCaptivePortalToken_SetString(p, params.CaptivePortalToken)
	if err != nil {
		return nil, err
	}

	return Packet(p), nil
}

func updateCoARequestRuckus(params Params, p *radius.Packet) (Packet, error) {
	err := ruckus.RuckusTCAttrIdsWithQuota_Set(p, trafficClassesToRuckus(params.TrafficClasses))
	if err != nil {
		return nil, err
	}
	err = ruckus.RuckusCPToken_SetString(p, params.CaptivePortalToken)
	if err != nil {
		return nil, err
	}

	return Packet(p), nil
}

func updateCoARequestByVendor(params Params, p *radius.Packet) (Packet, error) {
	switch params.VendorName {
	case VendorRuckus:
		return updateCoARequestRuckus(params, p)
	case VendorXWFCertified:
		return updateCoARequestXWFCertified(params, p)
	}
	return nil, fmt.Errorf("unknown vendor %d", params.VendorName)
}

// CreateCoARequest - create the request and return encoded packet
func CreateCoARequest(params Params, secret []byte) (Packet, error) {
	p := radius.New(radius.CodeCoARequest, secret)

	err := rfc2865.NASIdentifier_SetString(p, params.NASIdentifier)
	if err != nil {
		return nil, err
	}
	err = rfc2869.AcctInterimInterval_Set(p, rfc2869.AcctInterimInterval(params.AcctInterimInterval))
	if err != nil {
		return nil, err
	}
	err = rfc2866.AcctSessionID_SetString(p, params.AcctSessionID)
	if err != nil {
		return nil, err
	}
	err = rfc2865.CallingStationID_SetString(p, params.CallingStationID)
	if err != nil {
		return nil, err
	}

	return updateCoARequestByVendor(params, p)
}

// CreateCoADisconnect create the request and return encoded packet
func CreateCoADisconnect(params DisconnectParams, secret []byte) (Packet, error) {
	p := radius.New(radius.CodeDisconnectRequest, secret)

	err := rfc2865.NASIdentifier_SetString(p, params.NASIdentifier)
	if err != nil {
		return nil, err
	}
	err = rfc2866.AcctSessionID_SetString(p, params.AcctSessionID)
	if err != nil {
		return nil, err
	}
	err = rfc2865.CallingStationID_SetString(p, params.CallingStationID)
	if err != nil {
		return nil, err
	}

	return Packet(p), nil
}

// CreateClient - return client object that handle packet send
//
// retriesInterval - interval on which to resend packet (zero or negative value
// means no retry).
//
// maxPacketErrors controls how many packet parsing and validation errors
// the client will ignore before returning the error from Exchange (zero
// means drop all packet parsing errors).
func CreateClient(retriesInterval time.Duration, maxPacketErrors int) Client {
	return Client(radius.Client{
		Retry:           retriesInterval,
		MaxPacketErrors: maxPacketErrors,
	})
}

// Send send the packet to the given target
func (client Client) Send(ctx context.Context, p Packet, addr string) (Code, error) {
	radiusClient := radius.Client(client)
	response, err := radiusClient.Exchange(ctx, p, addr)
	if err != nil {
		return CodeUnknown, err
	}

	switch response.Code {
	case radius.CodeCoAACK, radius.CodeDisconnectACK:
		return CodeACK, nil
	case radius.CodeCoANAK, radius.CodeDisconnectNAK:
		return CodeNAK, nil
	default:
		return CodeUnknown, nil
	}
}
