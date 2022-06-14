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

import type {EditUser} from '../../EditUserDialog';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {WithAlert} from '../../Alert/withAlert';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AppContext from '../../../../app/components/context/AppContext';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import EditUserDialog from '../../EditUserDialog';
import IconButton from '@material-ui/core/IconButton';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../../theme/design-system/Text';
import axios from 'axios';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {UserRoles} from '../../../../shared/roles';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import renderList from '../../../util/renderList';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import withAlert from '../../Alert/withAlert';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useAxios} from '../../../hooks';
import {useCallback, useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../../app/hooks/useSnackbar';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

type Props = {...WithAlert};

function UsersSettings(props: Props) {
  const classes = useStyles();
  const [editingUser, setEditingUser] = useState<?EditUser>(null);
  const [users, setUsers] = useState<Array<EditUser>>([]);
  const [showDialog, setShowDialog] = useState<boolean>(false);
  const {networkIds, ssoEnabled} = useContext(AppContext);
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, error} = useAxios({
    url: '/user/async/',
    onResponse: useCallback(res => setUsers(res.data.users), []),
  });

  if (isLoading || error) {
    return <LoadingFiller />;
  }

  const handleError = error =>
    enqueueSnackbar(error.response?.data?.error || error, {variant: 'error'});

  const deleteUser = user => {
    props
      .confirm({
        message: (
          <span>
            Are you sure you want to delete the user{' '}
            <strong>{user.email}</strong>?
          </span>
        ),
        confirmLabel: 'Delete',
      })
      .then(confirmed => {
        if (confirmed) {
          axios
            .delete('/user/async/' + user.id)
            .then(_resp => setUsers(users.filter(u => u.id != user.id)))
            .catch(handleError);
        }
      });
  };

  const updateUserState = user => {
    const newUsers = users.slice(0);
    if (editingUser) {
      const index = users.indexOf(editingUser);
      newUsers[index] = user;
    } else {
      newUsers.push(user);
    }

    setShowDialog(false);
    setEditingUser(null);
    setUsers(newUsers);
  };

  const rows = users.map(row => (
    <TableRow key={row.id}>
      <TableCell>{row.email}</TableCell>
      <TableCell>
        {row.role == UserRoles.USER
          ? 'User'
          : row.role === UserRoles.READ_ONLY_USER
          ? 'Read Only User'
          : 'Super User'}
      </TableCell>
      <TableCell>{renderList(row.networkIDs || [])}</TableCell>
      <TableCell>
        <IconButton onClick={() => deleteUser(row)}>
          <DeleteIcon />
        </IconButton>
        <IconButton
          onClick={() => {
            setShowDialog(true);
            setEditingUser(row);
          }}>
          <EditIcon />
        </IconButton>
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.paper}>
      <div className={classes.header}>
        <Text variant="h5">Users</Text>
        <Button
          variant="contained"
          color="primary"
          onClick={() => setShowDialog(true)}>
          Add User
        </Button>
      </div>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Email</TableCell>
              <TableCell>Role</TableCell>
              <TableCell>Accessible Networks</TableCell>
              <TableCell />
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </Paper>
      {showDialog && (
        <EditUserDialog
          editingUser={editingUser}
          open={true}
          onClose={() => {
            setShowDialog(false);
            setEditingUser(null);
          }}
          ssoEnabled={ssoEnabled}
          allNetworkIDs={networkIds}
          onEditUser={(userId, payload) => {
            axios
              .put('/user/async/' + userId, payload)
              .then(response => updateUserState(response.data.user))
              .catch(handleError);
          }}
          onCreateUser={payload => {
            axios
              .post('/user/async/', payload)
              .then(response => updateUserState(response.data.user))
              .catch(handleError);
          }}
        />
      )}
    </div>
  );
}

export default withAlert(UsersSettings);
