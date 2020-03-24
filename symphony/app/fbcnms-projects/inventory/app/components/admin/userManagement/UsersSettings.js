/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {EditUser} from '@fbcnms/ui/components/auth/EditUserDialog';
import type {EditUserMutationResponse} from '../../../mutations/__generated__/EditUserMutation.graphql';
import type {MutationCallbacks} from '../../../mutations/MutationCallbacks.js';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@fbcnms/ui/components/design-system/Button';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import EditUserDialog from '@fbcnms/ui/components/auth/EditUserDialog';
import EditUserMutation from '../../../mutations/EditUserMutation';
import IconButton from '@material-ui/core/IconButton';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import RelayEnvironment from '../../../common/RelayEnvironment';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Text from '@fbcnms/ui/components/design-system/Text';
import axios from 'axios';
import {UserRoles} from '@fbcnms/auth/types';

import renderList from '@fbcnms/util/renderList';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {fetchQuery, graphql} from 'relay-runtime';
import {makeStyles} from '@material-ui/styles';
import {useAxios} from '@fbcnms/ui/hooks';
import {useCallback, useContext, useState} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';

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

const userQuery = graphql`
  query UsersSettings_UserQuery($authID: String!) {
    user(authID: $authID) {
      id
    }
  }
`;

type Props = {...WithAlert};

function UsersSettings(props: Props) {
  const classes = useStyles();
  const [editingUser, setEditingUser] = useState<?EditUser>(null);
  const [users, setUsers] = useState<Array<EditUser>>([]);
  const [showDialog, setShowDialog] = useState<boolean>(false);
  const {networkIds, ssoEnabled, isTabEnabled} = useContext(AppContext);
  const enqueueSnackbar = useEnqueueSnackbar();

  const isNmsEnabled = isTabEnabled('nms');

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

  const editUserInfo = (id, role) => {
    return new Promise((resolve, reject) => {
      const callbacks: MutationCallbacks<EditUserMutationResponse> = {
        onCompleted: (response, errors) => {
          if (errors && errors[0]) {
            reject(errors[0].message);
          }
          resolve(response.editUser);
        },
        onError: () => {
          reject('Error saving service');
        },
      };
      EditUserMutation(
        {
          input: {
            id: id,
            role: role === UserRoles.SUPERUSER ? 'OWNER' : 'USER',
          },
        },
        callbacks,
      );
    });
  };

  const fetchAndEditUser = (email, role) => {
    return fetchQuery(RelayEnvironment, userQuery, {
      authID: email,
    }).then(data => editUserInfo(data.user.id, role));
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
      {isNmsEnabled ? (
        <TableCell>{renderList(row.networkIDs || [])}</TableCell>
      ) : null}
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
        <Text variant="h5">Users - a</Text>
        <Button onClick={() => setShowDialog(true)}>Add User</Button>
      </div>
      <Paper elevation={2}>
        <Table>
          <TableHead>
            <TableRow>
              <TableCell>Email</TableCell>
              <TableCell>Role</TableCell>
              {isNmsEnabled ? <TableCell>Accessible Networks</TableCell> : null}
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
              .then(() => fetchAndEditUser(payload.email, payload.role))
              .catch(handleError);
          }}
          onCreateUser={payload => {
            axios
              .post('/user/async/', payload)
              .then(response => updateUserState(response.data.user))
              .then(() => fetchAndEditUser(payload.email, payload.role))
              .catch(handleError);
          }}
        />
      )}
    </div>
  );
}

export default withAlert(UsersSettings);
