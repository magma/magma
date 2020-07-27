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

// identity_helper provides Identity setter methods, missing in protobuf 3
// while protoc generates oneof type getters, setters re missing
package protos

import (
	"fmt"
	"reflect"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Identity type names table. Every Identity type should add a unique type name
// into identityTypeNameTable below. The names must be unique and should only
// include alphanumeric ASCII characters
var identityTypeNameTable = map[reflect.Type]string{
	reflect.TypeOf(&Identity_Gateway_{}):  "Gateway",
	reflect.TypeOf(&Identity_Operator{}):  "Operator",
	reflect.TypeOf(&Identity_Network{}):   "Network",
	reflect.TypeOf(&Identity_Wildcard_{}): "Wildcard",
}

// SetGateway sets Identity to be given gateway identity (gwid) and returns the
// receiver (Identity pointer)
func (id *Identity) SetGateway(gwId *Identity_Gateway) *Identity {
	if id != nil {
		id.Value = &Identity_Gateway_{gwId}
	}
	return id
}

// SetOperator sets Identity to be given operator identity (opid) and returns the
// receiver (Identity pointer)
func (id *Identity) SetOperator(opId string) *Identity {
	if id != nil {
		id.Value = &Identity_Operator{opId}
	}
	return id
}

// SetNetwork sets Identity to be given network identity (networkId) and returns
// the receiver (Identity pointer)
func (id *Identity) SetNetwork(networkId string) *Identity {
	if id != nil {
		id.Value = &Identity_Network{networkId}
	}
	return id
}

// NewGatewayIdentity returns Gateway identity corresponding to given hardware network &
// logical gateway IDs
func NewGatewayIdentity(hwId, networkId, logicalId string) *Identity {
	return new(Identity).SetGateway(
		&Identity_Gateway{
			HardwareId: hwId, NetworkId: networkId, LogicalId: logicalId})
}

// NewOperatorIdentity returns Operator identity corresponding to given opId
func NewOperatorIdentity(opId string) *Identity {
	return (&Identity{}).SetOperator(opId)
}

// NewNetworkIdentity returns Network identity corresponding to given networkId
func NewNetworkIdentity(networkId string) *Identity {
	return (&Identity{}).SetNetwork(networkId)
}

// NewGatewayWildcardIdentity returns Gateway wildcard Identity
func NewGatewayWildcardIdentity() *Identity {
	return &Identity{
		Value: &Identity_Wildcard_{Wildcard: &Identity_Wildcard{Type: Identity_Wildcard_Gateway}}}
}

// NewOperatorWildcardIdentity returns Operator wildcard Identity
func NewOperatorWildcardIdentity() *Identity {
	return &Identity{
		Value: &Identity_Wildcard_{Wildcard: &Identity_Wildcard{Type: Identity_Wildcard_Operator}}}
}

// NewNetworkWildcardIdentity returns Network Wildcard Identity
func NewNetworkWildcardIdentity() *Identity {
	return &Identity{
		Value: &Identity_Wildcard_{Wildcard: &Identity_Wildcard{Type: Identity_Wildcard_Network}}}
}

func (id *Identity) TypeName() *string {
	if tn, ok := identityTypeNameTable[reflect.TypeOf(id.Value)]; ok {
		return &tn
	}
	return nil
}

// HashString returns a unique string suitable to be used as a map/hash key
// uniquely identifying an Identity
// For instance, a network Identity's hash string would look like:
// "Id_Network_MyNetworkId"
// where 'MyNetworkId' is a registered network ID
func (id *Identity) HashString() string {
	var cn, tn string // common name & type name
	if id != nil {
		tnPtr := id.TypeName()
		if tnPtr != nil {
			tn = *tnPtr
		} else {
			// We don't want to fatal or panic in case of missing Identity type
			// name since Identity is a basic building block of Magma cloud
			// infrastructure.
			// For the same reasons - we avoid using external logging package
			// such as glog to not to introduce outside dependencies in Magma
			// core. log Print will never be disabled/filtered & should always
			// appear in std logs of a service/app alone with the offending
			// Identity type name.
			tn = "<UNDEFINED>"
			typ := reflect.TypeOf(id.Value)
			if typ.Kind() == reflect.Ptr {
				typ = typ.Elem() // dereference if ptr
			}
			glog.Errorf(
				"Magma Identity ERROR: type '%s' is missing Hashable Type Name",
				typ.Name())
		}
		cnp := id.ToCommonName()
		if cnp != nil {
			cn = *cnp
		} else {
			cn = "<nil>"
		}
	} else {
		tn = "<nil>"
	}
	return fmt.Sprintf("Id_%s_%s", tn, cn)
}

// GetHashableIdentitiesNumber returns a new table of type names, it's to be
// used only by unit tests and should not expose the original table
func GetHashableIdentitiesTable() map[string]string {
	res := map[string]string{}
	for k, v := range identityTypeNameTable {
		res[k.String()] = v
	}
	return res
}

// An embedded test type to be used by identity_test.go to make sure the protoc
// Oneof auto-generation assumptions used by the helper are still valid
// the test type cannot conflict with anything generated by protoc due to its
// type & implementation capitalization
type testIdentityWrapper struct {
	myTestIdentityImpl int
}

func (*testIdentityWrapper) isIdentity_Value() {}
func CreateTestIdentityImplValue() isIdentity_Value {
	return new(testIdentityWrapper)
}

// context key to get caller identity
type clientIdentityKey struct{}

// context key to set/get caller's Client Certificate Expiration Time
type clientCertificateExpirationKey struct{}

// NewContextWithCertExpiration returns a new Context that carries the given certificate expiration time
func NewContextWithCertExpiration(ctx context.Context, certExp int64) context.Context {
	if ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, clientCertificateExpirationKey{}, certExp)
}

// GetClientCertExpiration retreives from cxt & returns Unix time in seconds of client certificate
// expiration if present, 0 if not present
func GetClientCertExpiration(ctx context.Context) int64 {
	expTime, ok := ctx.Value(clientCertificateExpirationKey{}).(int64)
	if ok {
		return expTime
	}
	return 0
}

// GetGatewayIdentity returns the identity of the Gateway caller.
// Returns an error if the gateway is not registered.
func GetGatewayIdentity(ctx context.Context) (*Identity_Gateway, error) {
	id := GetClientGateway(ctx)
	if id == nil || !id.Registered() {
		return nil, status.Errorf(codes.PermissionDenied, "Gateway not registered")
	}
	return id, nil
}

// GetClientIdentity returns Identity of the RPC caller retrieved from GRPC/HTTP
// Context (if present) where it's set by the middleware or Obsidian
//
// NOTE: nil Identity is equivalent to missing Identity for all intents and
// purposes
func GetClientIdentity(ctx context.Context) *Identity {
	id, ok := ctx.Value(clientIdentityKey{}).(*Identity)
	if ok {
		return id
	}
	return nil
}

// GetClientGateway returns Identity of the Gateway caller retrieved from GRPC/HTTP
// Context (if present) where it's set by the middleware
// For use by all Gateway facing cloud services.
//   Example:
//      func (srv *MeteringdRecordsServer) UpdateFlows(
//			ctx context.Context,
//			tbl *protos.FlowTable) (*protos.Void, error) {
//
//          gw := protos.GetClientGateway(ctx)
//			if gw == nil || !gw.Registered() {
//              return &protos.Void{}, status.Errorf(
//                  codes.PermissionDenied, "Gateway not registered")
//          }
//
//		    ...
func GetClientGateway(ctx context.Context) *Identity_Gateway {
	return GetClientIdentity(ctx).GetGateway()
}

// NewContextWithIdentity returns a new Context that carries the given Identity
// Since nil Identity is equivalent to no Identity, no new CTX will be created
// if the passed id is nil and the passed ctx will be returned unmodified
func (id *Identity) NewContextWithIdentity(ctx context.Context) context.Context {
	if id == nil || ctx == nil {
		return ctx
	}
	return context.WithValue(ctx, clientIdentityKey{}, id)
}

// TBD (not currently implemented/enabled):
//
// Identity (such as a registered Gateway) may exist in a scope of another
// Identity (such as network).
// Such Identities may implement scopedIdentity interface and their GetScope()
// receiver will return an Identity of the scope that they exists within.
//
// The global scope is assumed for an Identity which returns nil from GetScope()
//
// The scoped Identity interface will be defined as specified below:
//
// type scopedIdentity interface {
//	 GetScope() *Identity
// }

type identityWithCommonName interface {
	ToCommonName() *string
}

// ToCommonName receiver assumes that each Identity implementation would
// provide its own ToCommonName and return a string to be used in cert Subject
// Common Name
// ToCommonName returns string pointer, nil will be returned if underlying
// Identity provider is missing (nil)
func (id *Identity) ToCommonName() *string {
	var v interface{} = id.GetValue()
	if v != nil {
		if icn, ok := v.(identityWithCommonName); ok {
			return icn.ToCommonName()
		}
	}
	return nil
}

// ToCommonName receiver for Gateway Identity
func (gw *Identity_Gateway_) ToCommonName() *string {
	if gw != nil && gw.Gateway != nil {
		return &gw.Gateway.HardwareId
	}
	return nil
}

// ToCommonName receiver for Operator Identity
func (oper *Identity_Operator) ToCommonName() *string {
	if oper != nil {
		return &oper.Operator
	}
	return nil
}

// ToCommonName receiver for Network Identity
func (nid *Identity_Network) ToCommonName() *string {
	if nid != nil {
		return &nid.Network
	}
	return nil
}

// ToCommonName receiver for Wildcard Identity
func (wc *Identity_Wildcard_) ToCommonName() *string {
	if wc != nil && wc.Wildcard != nil {
		if tn, ok := Identity_Wildcard_Type_name[int32(wc.Wildcard.Type)]; ok {
			return &tn // return Type Name (Network, Operator, etc.)
		}
	}
	return nil
}

// Registered check for Gateway Identity
func (gw *Identity_Gateway) Registered() bool {
	if gw != nil {
		return len(gw.GetLogicalId()) > 0 && len(gw.GetNetworkId()) > 0
	}
	return false
}

// A helper to verify if a given entity matches the wildcard pattern
// For now only match all (*) is supported
func (wildcard *Identity) Match(entity *Identity) bool {
	if wildcard != nil && entity != nil {
		wc, ok := wildcard.Value.(*Identity_Wildcard_)
		if ok && wc != nil && wc.Wildcard != nil {
			switch entity.Value.(type) {
			case *Identity_Gateway_:
				return wc.Wildcard.Type == Identity_Wildcard_Gateway
			case *Identity_Operator:
				return wc.Wildcard.Type == Identity_Wildcard_Operator
			case *Identity_Network:
				return wc.Wildcard.Type == Identity_Wildcard_Network
			}
		}
	}
	return false
}

func (id *Identity) GetWildcardForIdentity() *Identity {
	if id != nil {
		switch id.Value.(type) {
		case *Identity_Gateway_:
			return NewGatewayWildcardIdentity()
		case *Identity_Operator:
			return NewOperatorWildcardIdentity()
		case *Identity_Network:
			return NewNetworkWildcardIdentity()
		}
	}
	return nil
}

// GetHashToIdentity converts the passed slice to a map, whose keys are the hash strings of each Identity proto.
func GetHashToIdentity(ids []*Identity) map[string]*Identity {
	ret := make(map[string]*Identity)
	for _, id := range ids {
		ret[id.HashString()] = id
	}
	return ret
}
