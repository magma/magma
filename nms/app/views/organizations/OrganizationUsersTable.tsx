/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
import ActionTable, {TableRef} from '../../components/ActionTable';
import React from 'react';
import Text from '../../theme/design-system/Text';
import axios from 'axios';
import withAlert from '../../components/Alert/withAlert';
import {OrganizationUser} from './types';
import {QueryResult} from '@material-table/core';
import {UserRoles} from '../../../shared/roles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import {useState} from 'react';
import type {WithAlert} from '../../components/Alert/withAlert';

type OrganizationUsersTableProps = WithAlert & {
  editUser: (user: OrganizationUser | null) => void;
  tableRef: TableRef;
};

/**
 * Table of users that belong to a specific organization
 */
function OrganizationUsersTable(props: OrganizationUsersTableProps) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [users, setUsers] = React.useState<Array<OrganizationUser>>([]);
  const [currRow, setCurrRow] = useState<OrganizationUser>(
    {} as OrganizationUser,
  );
  const params = useParams();

  const onDeleteUser = (user: OrganizationUser) => {
    void props
      .confirm({
        message: (
          <span>
            {'Are you sure you want to delete the user '}
            <strong>{user.email}</strong>?
          </span>
        ),
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (confirmed) {
          axios
            .delete(`/user/async/${user.id}`)
            .then(() => {
              props.tableRef.current?.onQueryChange();
            })
            .catch(() => {
              enqueueSnackbar(`Unable to delete user: ${user.id}`, {
                variant: 'error',
              });
            });
        }
      });
  };

  const menuItems = [
    {
      name: 'Edit',
      handleFunc: () => {
        const user = users.find(user => user.id === currRow.id);
        props.editUser?.(user!);
        props.tableRef.current?.onQueryChange();
      },
    },
    {
      name: 'Remove',
      handleFunc: () => {
        onDeleteUser(currRow);
      },
    },
  ];
  const columnStruct = [
    {
      title: '',
      field: '',
      width: '40px',
      render: (rowData: OrganizationUser) => (
        <Text variant="subtitle3">
          {((rowData as unknown) as {tableData: {id: number}}).tableData?.id +
            1}
        </Text>
      ),
    },
    {
      title: 'Email',
      field: 'email',
    },
    {
      title: 'Role',
      field: 'role',
      render: (rowData: OrganizationUser) => {
        const userRole = (Object.keys(UserRoles) as Array<
          keyof typeof UserRoles
        >).find(role => UserRoles[role] === rowData.role);
        return <>{userRole}</>;
      },
    },
  ];

  return (
    <>
      <ActionTable
        tableRef={props.tableRef}
        data={() =>
          axios
            .get<Array<OrganizationUser>>(
              `/host/organization/async/${params.name!}/users`,
            )
            .then(result => {
              const users = result.data.map(user => {
                return {
                  email: user.email,
                  role: user.role,
                  id: user.id,
                  networkIDs: user.networkIDs,
                  organization: user.organization,
                };
              });
              setUsers(users);
              return {
                data: users,
              } as QueryResult<OrganizationUser>;
            })
        }
        columns={columnStruct}
        handleCurrRow={(row: OrganizationUser) => {
          setCurrRow(row);
        }}
        menuItems={menuItems}
        localization={{
          // hide 'Actions' in table header
          header: {actions: ''},
        }}
        options={{
          actionsColumnIndex: -1,
          sorting: true,
          // hide table title and toolbar
          toolbar: false,
          paging: false,
          pageSizeOptions: [100, 200],
        }}
      />
    </>
  );
}

export default withAlert(OrganizationUsersTable);
