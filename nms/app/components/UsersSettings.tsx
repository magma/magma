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

import AppContext from './context/AppContext';
import Button from '@material-ui/core/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import EditUserDialog from './EditUserDialog';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from './LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '../theme/design-system/Text';
import axios, {AxiosResponse} from 'axios';
import {UserRoles} from '../../shared/roles';
import type {EditUser} from './EditUserDialog';
import type {WithAlert} from './Alert/withAlert';

import renderList from '../util/renderList';
import withAlert from './Alert/withAlert';
import {Theme} from '@material-ui/core/styles';
import {getErrorMessage} from '../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '../hooks';
import {useCallback, useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

const useStyles = makeStyles<Theme>(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

function UsersSettings(props: WithAlert) {
  const classes = useStyles();
  const [editingUser, setEditingUser] = useState<EditUser | null>(null);
  const [users, setUsers] = useState<Array<EditUser>>([]);
  const [showDialog, setShowDialog] = useState<boolean>(false);
  const {networkIds, ssoEnabled} = useContext(AppContext);
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, error} = useAxios<{users: Array<EditUser>}>({
    url: '/user/async/',
    onResponse: useCallback(
      (res: AxiosResponse<{users: Array<EditUser>}>) =>
        setUsers(res.data.users),
      [],
    ),
  });

  if (isLoading || error) {
    return <LoadingFiller />;
  }

  const handleError = (error: unknown) =>
    enqueueSnackbar(getErrorMessage(error), {variant: 'error'});

  const deleteUser = (user: EditUser) => {
    void props
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
            .then(() => setUsers(users.filter(u => u.id != user.id)))
            .catch(handleError);
        }
      });
  };

  const updateUserState = (user: EditUser) => {
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
              .put<{user: EditUser}>('/user/async/' + userId, payload)
              .then(response => updateUserState(response.data.user))
              .catch(handleError);
          }}
          onCreateUser={payload => {
            axios
              .post<{user: EditUser}>('/user/async/', payload)
              .then(response => updateUserState(response.data.user))
              .catch(handleError);
          }}
        />
      )}
    </div>
  );
}

export default withAlert(UsersSettings);
