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

// Package accessd provides a thin client for access management service.
package accessd

import (
	"context"
	"errors"
	"fmt"

	"github.com/golang/glog"

	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	"magma/orc8r/lib/go/merrors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const ServiceName = "accessd"

// getAccessbClient is a utility function to get a RPC connection to the
// accessd service
func getAccessdClient() (accessprotos.AccessControlManagerClient, error) {
	conn, err := registry.GetConnection(ServiceName, protos.ServiceType_PROTECTED)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return accessprotos.NewAccessControlManagerClient(conn), err
}

// SetOperator overwrites Permissions to operator Identity to manage/monitor
// entities
func SetOperator(ctx context.Context, operator *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}

	_, err = client.SetOperator(ctx, &accessprotos.AccessControl_ListRequest{Operator: operator, Entities: entities})
	if err != nil {
		errMsg := fmt.Sprintf("Set Permissions for Operator %s error: %s", operator.HashString(), err)
		glog.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// UpdateOperator adds Permissions to operator Identity to manage/monitor
// entities
func UpdateOperator(ctx context.Context, operator *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.UpdateOperator(
		ctx,
		&accessprotos.AccessControl_ListRequest{Operator: operator, Entities: entities})
	if err != nil {
		errMsg := fmt.Sprintf("Add Permissions for Operator %s error: %s",
			operator.HashString(), err)
		glog.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// Removes all operator's permissions (the entire operator's ACL)
func DeleteOperator(ctx context.Context, operator *protos.Identity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteOperator(ctx, operator)
	if err != nil {
		errMsg := fmt.Sprintf("Revoke Permissions for Operator %s error: %s",
			operator.HashString(), err)
		glog.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// GetOperatorACL returns the operator's Identity's permission list
func GetOperatorACL(
	ctx context.Context,
	operator *protos.Identity,
) (map[string]*accessprotos.AccessControl_Entity, error) {
	client, err := getAccessdClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetOperatorACL(ctx, operator)
	if err != nil {
		errMsg := fmt.Sprintf("Get Permissions for Operator %s error: %s",
			operator.HashString(), err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resp.Entities, nil
}

// GetOperatorsACLs returns the operators' Identities permission lists
func GetOperatorsACLs(ctx context.Context, operators []*protos.Identity) ([]*accessprotos.AccessControl_List, error) {
	if len(operators) == 0 {
		return nil, nil
	}
	client, err := getAccessdClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetOperatorsACLs(ctx, &protos.Identity_List{List: operators})
	if err != nil || resp == nil {
		errMsg := fmt.Sprintf("Get Permissions for Operators %v error: %s", operators, err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resp.Acls, nil
}

// Returns the operator's permission bitmask for given entity
func GetPermissions(
	ctx context.Context,
	operator *protos.Identity,
	entity *protos.Identity,
) (accessprotos.AccessControl_Permission, error) {
	client, err := getAccessdClient()
	if err != nil {
		return accessprotos.AccessControl_NONE, err
	}
	resp, err := client.GetPermissions(ctx, &accessprotos.AccessControl_PermissionsRequest{Operator: operator, Entity: entity})
	if err != nil {
		errMsg := fmt.Sprintf("Operator %s Permissions Check for %s error: %s",
			operator.HashString(), entity.HashString(), err)
		glog.Error(errMsg)
		return accessprotos.AccessControl_NONE, errors.New(errMsg)
	}
	return resp.Permissions, nil
}

// Verifies operator's read permission for given entity and returns error if
// either request fails or the permissions are not granted
func CheckReadPermission(ctx context.Context, operator *protos.Identity, ents ...*protos.Identity) error {
	entsPerm := make([]*accessprotos.AccessControl_Entity, len(ents))
	for i, e := range ents {
		entsPerm[i] = &accessprotos.AccessControl_Entity{Id: e, Permissions: accessprotos.AccessControl_READ}
	}
	return CheckPermissions(ctx, operator, entsPerm...)
}

// Verifies operator's write permission for given entity and returns error if
// either request fails or the permissions are not granted
func CheckWritePermission(ctx context.Context, operator *protos.Identity, ents ...*protos.Identity) error {
	entsPerm := make([]*accessprotos.AccessControl_Entity, len(ents))
	for i, e := range ents {
		entsPerm[i] = &accessprotos.AccessControl_Entity{Id: e, Permissions: accessprotos.AccessControl_WRITE}
	}
	return CheckPermissions(ctx, operator, entsPerm...)
}

func CheckPermissions(ctx context.Context, operator *protos.Identity, ents ...*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.CheckPermissions(ctx, &accessprotos.AccessControl_ListRequest{Operator: operator, Entities: ents})
	return err
}

// List all Operator Identities in accessd database
func ListOperators(ctx context.Context) ([]*protos.Identity, error) {
	client, err := getAccessdClient()
	if err != nil {
		return []*protos.Identity{}, err
	}
	opslist, err := client.ListOperators(ctx, &protos.Void{})
	if err != nil || opslist == nil {
		return []*protos.Identity{}, err
	}
	return opslist.List, nil
}
