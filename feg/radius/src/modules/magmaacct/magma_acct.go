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

package magmaacct

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"fbc/cwf/radius/modules"
	"fbc/cwf/radius/modules/protos"
	"fbc/cwf/radius/session"
	"fbc/lib/go/radius"
	"fbc/lib/go/radius/rfc2865"
	"fbc/lib/go/radius/rfc2866"

	"github.com/mitchellh/mapstructure"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// Config configuration structure for proxy module
type Config struct {
	FegEndpoint string
}

// ModuleCtx ...
type ModuleCtx struct {
	client protos.AccountingClient
}

// Init module interface implementation
func Init(logger *zap.Logger, config modules.ModuleConfig) (modules.Context, error) {
	var acctConfig Config
	err := mapstructure.Decode(config, &acctConfig)
	if err != nil {
		return nil, err
	}

	if acctConfig.FegEndpoint == "" {
		return nil, errors.New("magma acct module cannot be initialize with empty FegEndpoint value")
	}

	// Initialize the client
	conn, err := grpc.Dial(acctConfig.FegEndpoint, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}

	return ModuleCtx{client: protos.NewAccountingClient(conn)}, nil
}

// Handle module interface implementation
func Handle(m modules.Context, ctx *modules.RequestContext, r *radius.Request, _ modules.Middleware) (*modules.Response, error) {
	mCtx := m.(ModuleCtx)

	// Load the state
	state, err := ctx.SessionStorage.Get()
	if err != nil {
		state = &session.State{}
		attr, err := rfc2865.CallingStationID_Lookup(r.Packet)
		if err == nil {
			state.MACAddress = string(attr)
		}
		attr, err = rfc2865.CalledStationID_Lookup(r.Packet)
		if err == nil {
			state.CalledStationID = string(attr)
		}
	}

	// TODO: this should be moved to the server itself (and not in a module)
	// Update Acct-Session-ID attribute in state
	acctSessionID := r.Get(rfc2866.AcctSessionID_Type)
	if acctSessionID != nil && state.AcctSessionID != string(acctSessionID) {
		state.AcctSessionID = string(acctSessionID)
		ctx.Logger.Debug(
			"updating Acct-Session-Id in session state",
			zap.String("acct_session_id", state.AcctSessionID),
		)
		_ = ctx.SessionStorage.Set(*state)
	}

	// Get accounting type
	acctTypeAttr, exists := r.Lookup(rfc2866.AcctStatusType_Type)
	if !exists {
		return nil, errors.New("invalid RADIUS request without Acct-Status-Type attribute")
	}
	acctType := rfc2866.AcctStatusType(binary.BigEndian.Uint32(acctTypeAttr))

	// Restore Context
	c := &protos.Context{
		SessionId: ctx.SessionID,
		Msisdn:    state.MSISDN,
		MacAddr:   state.MACAddress,
		IpAddr:    strings.Split(r.RemoteAddr.String(), ":")[0],
		Apn:       state.CalledStationID,
	}

	// Call magma client
	switch acctType {
	case rfc2866.AcctStatusType_Value_AccountingOn:
	case rfc2866.AcctStatusType_Value_Start:
		_, err = mCtx.client.Start(context.Background(), c)
		if err != nil {
			return nil, err
		}
		ctx.Logger.Debug("MagmaAccounting.Start succeeded", zap.Any("context", c))
		break
	case rfc2866.AcctStatusType_Value_AccountingOff:
	case rfc2866.AcctStatusType_Value_Stop:
		stopRequest := &protos.StopRequest{
			Cause:     protos.StopRequest_NAS_REQUEST,
			Ctx:       c,
			OctetsIn:  getValue(r, rfc2866.AcctInputOctets_Type),
			OctetsOut: getValue(r, rfc2866.AcctOutputOctets_Type),
		}
		_, err = mCtx.client.Stop(context.Background(), stopRequest)
		if err != nil {
			return nil, err
		}
		ctx.Logger.Debug("MagmaAccounting.Stop succeeded", zap.Any("context", c))
		break
	case rfc2866.AcctStatusType_Value_InterimUpdate:
		updateRequest := &protos.UpdateRequest{
			OctetsIn:   getValue(r, rfc2866.AcctInputOctets_Type),
			OctetsOut:  getValue(r, rfc2866.AcctOutputOctets_Type),
			PacketsIn:  getValue(r, rfc2866.AcctInputPackets_Type),
			PacketsOut: getValue(r, rfc2866.AcctOutputPackets_Type),
			Ctx:        c,
		}
		_, err = mCtx.client.InterimUpdate(context.Background(), updateRequest)
		if err != nil {
			return nil, err
		}
		ctx.Logger.Debug("MagmaAccounting.InterimUpdate succeeded", zap.Any("context", c))
		break
	default:
		return nil, fmt.Errorf("unknown Acct-Status-Type received: %d", acctType)
	}

	// Build response
	result := &modules.Response{
		Code: radius.CodeAccountingResponse,
		Attributes: radius.Attributes{
			rfc2866.AcctSessionID_Type: []radius.Attribute{radius.Attribute(c.SessionId)},
		},
	}

	ctx.Logger.Debug(
		"successfully handled Accounting Request",
		zap.Any("context", c),
		zap.Any("result", result),
	)
	return result, nil
}

func getValue(r *radius.Request, t radius.Type) uint32 {
	valueAttr, exists := r.Lookup(t)
	var value uint32
	if exists {
		value = binary.BigEndian.Uint32(valueAttr)
	}
	return value
}
