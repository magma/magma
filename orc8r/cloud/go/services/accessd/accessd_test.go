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

package accessd_test

import (
	"testing"

	"magma/orc8r/cloud/go/identity"
	"magma/orc8r/cloud/go/services/accessd"
	accessprotos "magma/orc8r/cloud/go/services/accessd/protos"
	accessd_test_service "magma/orc8r/cloud/go/services/accessd/test_init"
	"magma/orc8r/lib/go/protos"

	"github.com/stretchr/testify/assert"
)

func TestAccessManager(t *testing.T) {
	accessd_test_service.StartTestService(t)

	op1 := identity.NewOperator("operator1")
	assert.NotEmpty(t, op1.ToCommonName())
	op2 := identity.NewOperator("operator2")
	assert.Equal(t, op2.HashString(), "Id_Operator_operator2")

	net1 := identity.NewNetwork("network1")
	assert.Equal(t, *net1.ToCommonName(), "network1")
	net2 := identity.NewNetwork("network2")
	assert.Equal(t, *net2.ToCommonName(), "network2")

	entities := []*accessprotos.AccessControl_Entity{
		{Id: net1, Permissions: accessprotos.AccessControl_READ},
		{Id: net2, Permissions: accessprotos.AccessControl_WRITE},
	}
	err := accessd.UpdateOperator(op1, entities)
	assert.Error(t, err)

	err = accessd.SetOperator(op1, entities)
	assert.NoError(t, err)

	acl, err := accessd.GetOperatorACL(op1)
	assert.NoError(t, err)
	assert.NotEmpty(t, acl)

	assert.Equal(t, len(acl), 2)
	ent, ok := acl[net1.HashString()]
	assert.True(t, ok)
	assert.Equal(t, ent.Permissions, accessprotos.AccessControl_READ)

	ent, ok = acl[net2.HashString()]
	assert.True(t, ok)
	assert.Equal(t, ent.Permissions, accessprotos.AccessControl_WRITE)

	assert.NoError(t, accessd.CheckReadPermission(op1, net1))
	assert.Error(t, accessd.CheckWritePermission(op1, net1))

	assert.Error(t, accessd.CheckReadPermission(op1, net2))
	assert.NoError(t, accessd.CheckWritePermission(op1, net2))

	net3 := (&protos.Identity{}).SetNetwork("network3")

	entitiesToAdd := []*accessprotos.AccessControl_Entity{
		{Id: net3, Permissions: accessprotos.AccessControl_READ},
	}
	err = accessd.UpdateOperator(op1, entitiesToAdd)
	assert.NoError(t, err)

	assert.NoError(t, accessd.CheckReadPermission(op1, net3))
	assert.Error(t, accessd.CheckWritePermission(op1, net3))

	err = accessd.CheckPermissions(
		op1,
		&accessprotos.AccessControl_Entity{Id: net1, Permissions: accessprotos.AccessControl_READ},
		&accessprotos.AccessControl_Entity{Id: net2, Permissions: accessprotos.AccessControl_WRITE},
		&accessprotos.AccessControl_Entity{Id: net3, Permissions: accessprotos.AccessControl_READ})

	assert.NoError(t, err)

	err = accessd.CheckPermissions(
		op1,
		&accessprotos.AccessControl_Entity{Id: net2, Permissions: accessprotos.AccessControl_READ},
		&accessprotos.AccessControl_Entity{Id: net3, Permissions: accessprotos.AccessControl_READ})

	assert.Error(t, err)

	removeNet13Ents := []*accessprotos.AccessControl_Entity{
		{Id: net2, Permissions: accessprotos.AccessControl_WRITE},
	}
	err = accessd.SetOperator(op1, removeNet13Ents)
	assert.NoError(t, err)

	assert.Error(t, accessd.CheckReadPermission(op1, net1))
	assert.Error(t, accessd.CheckWritePermission(op1, net1))

	assert.Error(t, accessd.CheckReadPermission(op1, net3))
	assert.NoError(t, accessd.CheckWritePermission(op1, net2))

	perm, err := accessd.GetPermissions(op1, net1)
	assert.NoError(t, err)
	assert.Equal(t, perm, accessprotos.AccessControl_NONE)

	_, err = accessd.GetPermissions(op1, net3)
	assert.NoError(t, err)

	perm, err = accessd.GetPermissions(op1, net2)
	assert.NoError(t, err)
	assert.Equal(t, perm, accessprotos.AccessControl_WRITE)

	err = accessd.SetOperator(op2, entities)
	assert.NoError(t, err)

	assert.NoError(t, accessd.CheckReadPermission(op2, net1))
	assert.Error(t, accessd.CheckWritePermission(op2, net1))

	assert.Error(t, accessd.CheckReadPermission(op2, net2))
	assert.NoError(t, accessd.CheckWritePermission(op2, net2))

	opers, err := accessd.ListOperators()
	assert.NoError(t, err)
	assert.Len(t, opers, 2)
	if len(opers) >= 2 {
		expected := []string{"Id_Operator_operator1", "Id_Operator_operator2"}
		assert.Contains(t, expected, opers[0].HashString())
		assert.Contains(t, expected, opers[1].HashString())
	}

	err = accessd.DeleteOperator(op1)
	assert.NoError(t, err)
	assert.Error(t, accessd.CheckReadPermission(op1, net1))
	assert.Error(t, accessd.CheckWritePermission(op1, net2))
	_, err = accessd.GetPermissions(op1, net2)
	assert.Error(t, err)

	entitiesToAdd = []*accessprotos.AccessControl_Entity{ // WRITE perm for all Networks
		{Id: identity.NewNetworkWildcard(), Permissions: accessprotos.AccessControl_WRITE},
	}
	err = accessd.UpdateOperator(op2, entitiesToAdd)
	assert.NoError(t, err)
	assert.NoError(t, accessd.CheckReadPermission(op2, net1))
	assert.NoError(t, accessd.CheckWritePermission(op2, net1))
	assert.Error(t, accessd.CheckReadPermission(op2, net2))
	assert.NoError(t, accessd.CheckWritePermission(op2, net3))
	assert.NoError(t, accessd.CheckWritePermission(
		op2, identity.NewNetwork("some_network")))
	assert.Error(t, accessd.CheckReadPermission(
		op2, identity.NewNetwork("some_network2")))

	err = accessd.CheckPermissions(
		op2,
		&accessprotos.AccessControl_Entity{Id: net1, Permissions: accessprotos.AccessControl_READ},
		&accessprotos.AccessControl_Entity{Id: net2, Permissions: accessprotos.AccessControl_WRITE},
		&accessprotos.AccessControl_Entity{Id: net3, Permissions: accessprotos.AccessControl_WRITE})

	assert.NoError(t, err)

	err = accessd.CheckPermissions(
		op2,
		&accessprotos.AccessControl_Entity{Id: net1, Permissions: accessprotos.AccessControl_READ},
		&accessprotos.AccessControl_Entity{Id: net2, Permissions: accessprotos.AccessControl_READ})

	assert.Error(t, err)

	opers, err = accessd.ListOperators()
	assert.NoError(t, err)
	assert.Len(t, opers, 1)
	if len(opers) > 0 {
		assert.Equal(t, "Id_Operator_operator2", opers[0].HashString())
	}
}
