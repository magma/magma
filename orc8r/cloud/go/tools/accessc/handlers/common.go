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

// Package handlers implements individual accessc commands as well as common
// across multiple commands functionality
package handlers

import (
	"fmt"

	"magma/orc8r/cloud/go/identity"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/cloud/go/tools/commands"
	"magma/orc8r/lib/go/protos"
)

var (
	CommandRegistry = new(commands.Map)
	makeAdmin       bool // flag for add/modify commands
)

// BuildACLForEntities builds ACL from networks, operators & gateways command
// line flags
func BuildACLForEntities(networks, operators, gateways Entities) []*accessprotos.AccessControl_Entity {
	// Assemble ACL for the Operator
	// NOTE: Gateway granularity access control is not implemented yet
	var acl []*accessprotos.AccessControl_Entity
	for _, e := range networks {
		if e.id == "*" {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id:          identity.NewNetworkWildcard(),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		} else {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id:          identity.NewNetwork(e.id),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		}
	}
	for _, e := range operators {
		if e.id == "*" {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id:          identity.NewOperatorWildcard(),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		} else {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id:          identity.NewOperator(e.id),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		}
	}
	for _, e := range gateways {
		if e.id == "*" {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id:          identity.NewGatewayWildcard(),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		} else {
			acl = append(
				acl,
				&accessprotos.AccessControl_Entity{
					Id: new(protos.Identity).SetGateway(
						&protos.Identity_Gateway{HardwareId: e.id}),
					Permissions: accessprotos.AccessControl_Permission(e.perm)})
		}
	}
	return acl
}

// CreateAdminACL Constructs "can do everything" ACL for administrators
func CreateAdminACL() []*accessprotos.AccessControl_Entity {
	perm := accessprotos.ACCESS_CONTROL_ALL_PERMISSIONS
	return []*accessprotos.AccessControl_Entity{
		{Id: identity.NewNetworkWildcard(), Permissions: perm},
		{Id: identity.NewOperatorWildcard(), Permissions: perm},
		{Id: identity.NewGatewayWildcard(), Permissions: perm}}
}

// PrintACL - prints operator ID, its certificate SNs & ACL
func PrintACL(acl *accessprotos.AccessControl_List, certSNs []string) {
	op := acl.GetOperator()
	opname := op.HashString()
	cn := "<nil>"
	cnPtr := op.ToCommonName()
	if cnPtr != nil {
		cn = *cnPtr
	}
	fmt.Printf("\t%s (%s): Certificates:%s;\n\t\tACL:\n", cn, opname, certSNs)
	for entname, ent := range acl.Entities {
		fmt.Printf(
			"\t\t  %s: %s (%d)\n",
			entname,
			ent.Permissions.ToString(),
			ent.Permissions)
	}
	fmt.Println()
}
