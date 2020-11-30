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
	"errors"
	"fmt"

	"github.com/golang/glog"
	"golang.org/x/net/context"

	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	merrors "magma/orc8r/lib/go/errors"
	"magma/orc8r/lib/go/protos"
	"magma/orc8r/lib/go/registry"
)

const ServiceName = "ACCESSD"

// getAccessbClient is a utility function to get a RPC connection to the
// accessd service
func getAccessdClient() (accessprotos.AccessControlManagerClient, error) {
	conn, err := registry.GetConnection(ServiceName)
	if err != nil {
		initErr := merrors.NewInitError(err, ServiceName)
		glog.Error(initErr)
		return nil, initErr
	}
	return accessprotos.NewAccessControlManagerClient(conn), err
}

// SetOperator overwrites Permissions to operator Identity to manage/monitor
// entities
func SetOperator(operator *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}

	_, err = client.SetOperator(context.Background(), &accessprotos.AccessControl_ListRequest{Operator: operator, Entities: entities})
	if err != nil {
		errMsg := fmt.Sprintf("Set Permissions for Operator %s error: %s", operator.HashString(), err)
		glog.Error(errMsg)
		return errors.New(errMsg)
	}
	return nil
}

// UpdateOperator adds Permissions to operator Identity to manage/monitor
// entities
func UpdateOperator(operator *protos.Identity, entities []*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.UpdateOperator(
		context.Background(),
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
func DeleteOperator(operator *protos.Identity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.DeleteOperator(context.Background(), operator)
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
	operator *protos.Identity,
) (map[string]*accessprotos.AccessControl_Entity, error) {
	client, err := getAccessdClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetOperatorACL(context.Background(), operator)
	if err != nil {
		errMsg := fmt.Sprintf("Get Permissions for Operator %s error: %s",
			operator.HashString(), err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resp.Entities, nil
}

// GetOperatorsACLs returns the operators' Identities permission lists
func GetOperatorsACLs(operators []*protos.Identity) ([]*accessprotos.AccessControl_List, error) {
	if len(operators) == 0 {
		return nil, nil
	}
	client, err := getAccessdClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.GetOperatorsACLs(context.Background(), &protos.Identity_List{List: operators})
	if err != nil || resp == nil {
		errMsg := fmt.Sprintf("Get Permissions for Operators %v error: %s", operators, err)
		glog.Error(errMsg)
		return nil, errors.New(errMsg)
	}
	return resp.Acls, nil
}

// Returns the operator's permission bitmask for given entity
func GetPermissions(
	operator *protos.Identity,
	entity *protos.Identity,
) (accessprotos.AccessControl_Permission, error) {
	client, err := getAccessdClient()
	if err != nil {
		return accessprotos.AccessControl_NONE, err
	}
	resp, err := client.GetPermissions(
		context.Background(),
		&accessprotos.AccessControl_PermissionsRequest{Operator: operator, Entity: entity})
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
func CheckReadPermission(operator *protos.Identity, ents ...*protos.Identity) error {
	entsPerm := make([]*accessprotos.AccessControl_Entity, len(ents))
	for i, e := range ents {
		entsPerm[i] = &accessprotos.AccessControl_Entity{Id: e, Permissions: accessprotos.AccessControl_READ}
	}
	return CheckPermissions(operator, entsPerm...)
}

// Verifies operator's write permission for given entity and returns error if
// either request fails or the permissions are not granted
func CheckWritePermission(operator *protos.Identity, ents ...*protos.Identity) error {
	entsPerm := make([]*accessprotos.AccessControl_Entity, len(ents))
	for i, e := range ents {
		entsPerm[i] = &accessprotos.AccessControl_Entity{Id: e, Permissions: accessprotos.AccessControl_WRITE}
	}
	return CheckPermissions(operator, entsPerm...)
}

func CheckPermissions(operator *protos.Identity, ents ...*accessprotos.AccessControl_Entity) error {
	client, err := getAccessdClient()
	if err != nil {
		return err
	}
	_, err = client.CheckPermissions(
		context.Background(), &accessprotos.AccessControl_ListRequest{Operator: operator, Entities: ents})
	return err
}

// List all Operator Identities in accessd database
func ListOperators() ([]*protos.Identity, error) {
	client, err := getAccessdClient()
	if err != nil {
		return []*protos.Identity{}, err
	}
	opslist, err := client.ListOperators(context.Background(), &protos.Void{})
	if err != nil || opslist == nil {
		return []*protos.Identity{}, err
	}
	return opslist.List, nil
}
