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
 *
 * @flow strict-local
 * @format
 */
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import ActionTable from '../../components/ActionTable';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
import axios from 'axios';
import withAlert from '../../components/Alert/withAlert';
import type {EditUser} from './OrganizationEdit';
import type {WithAlert} from '../../components/Alert/withAlert';

import {UserRoles} from '../../../shared/roles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import {useState} from 'react';

type OrganizationUsersTableProps = WithAlert & {
  editUser: (user: ?EditUser) => void,
  tableRef: {current: null | {onQueryChange(): void}},
};

/**
 * Table of users that belong to a specific organization
 */
function OrganizationUsersTable(props: OrganizationUsersTableProps) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const [users, setUsers] = React.useState<Array<EditUser>>([]);
  const [currRow, setCurrRow] = useState<EditUser>({});
  const params = useParams();

  const onDeleteUser = user => {
    props
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
            .delete('/user/async/' + user.id)
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
        const user: ?EditUser = users.find(user => user.id === currRow.id);
        props.editUser?.(user);
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
      render: rowData => (
        <Text variant="subtitle3">{rowData.tableData?.id + 1}</Text>
      ),
    },
    {
      title: 'Email',
      field: 'email',
    },
    {
      title: 'Role',
      field: 'role',
      render: rowData => {
        const userRole = Object.keys(UserRoles).find(
          role => UserRoles[role] === rowData.role,
        );
        return <>{userRole}</>;
      },
    },
  ];

  return (
    <>
      <ActionTable
        tableRef={props.tableRef}
        data={() =>
          new Promise((resolve, _reject) => {
            axios
              .get(`/host/organization/async/${params.name}/users`)
              .then(result => {
                const users: Array<EditUser> = result.data.map(user => {
                  return {
                    email: user.email,
                    role: user.role,
                    id: user.id,
                    networkIDs: user.networkIDs,
                    organization: user.organization,
                  };
                });
                setUsers(users);
                resolve({
                  data: users,
                });
              });
          })
        }
        columns={columnStruct}
        handleCurrRow={(row: EditUser) => {
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
