/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {EditUser} from '@fbcnms/ui/components/auth/EditUserDialog';
import type {WithStyles} from '@material-ui/core';
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import axios from 'axios';
import React from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Typography from '@material-ui/core/Typography';
import Button from '@material-ui/core/Button';
import EditUserDialog from '@fbcnms/ui/components/auth/EditUserDialog';
import DeleteIcon from '@material-ui/icons/Delete';
import EditIcon from '@material-ui/icons/Edit';
import IconButton from '@material-ui/core/IconButton';
import {UserRoles} from '@fbcnms/auth/types';

import renderList from '@fbcnms/util/renderList';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';
import {withStyles} from '@material-ui/core/styles';

const styles = {
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  paper: {
    margin: '10px',
  },
};

type Props = WithStyles &
  WithAlert & {
    allNetworkIDs: Array<string>,
  };

type State = {
  editingUser: ?EditUser,
  users: Array<EditUser>,
  showDialog: boolean,
};

class UsersSettings extends React.Component<Props, State> {
  state = {
    editingUser: null,
    showDialog: false,
    users: [],
  };

  componentDidMount() {
    axios
      .get('/nms/user/async/')
      .then(response => this.setState({users: response.data.users}));
  }

  render() {
    const rows = this.state.users.map(row => (
      <TableRow key={row.id}>
        <TableCell>{row.email}</TableCell>
        <TableCell>
          {row.role == UserRoles.USER ? 'User' : 'Super User'}
        </TableCell>
        <TableCell>{renderList(row.networkIDs || [])}</TableCell>
        <TableCell>
          <IconButton onClick={this.deleteUser.bind(this, row)}>
            <DeleteIcon />
          </IconButton>
          <IconButton onClick={this.showEditDialog.bind(this, row)}>
            <EditIcon />
          </IconButton>
        </TableCell>
      </TableRow>
    ));

    return (
      <div className={this.props.classes.paper}>
        <div className={this.props.classes.header}>
          <Typography variant="h5">Users</Typography>
          <Button variant="contained" color="primary" onClick={this.showDialog}>
            Add User
          </Button>
        </div>
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
        <EditUserDialog
          key={this.state.editingUser ? this.state.editingUser.id : 'new_user'}
          editingUser={this.state.editingUser}
          open={this.state.showDialog}
          onClose={this.hideDialog}
          onEditUser={this.editUser}
          onCreateUser={this.createUser}
          allNetworkIDs={this.props.allNetworkIDs}
        />
      </div>
    );
  }

  showDialog = () => this.setState({showDialog: true});
  hideDialog = () => this.setState({editingUser: null, showDialog: false});
  showEditDialog = user => this.setState({editingUser: user, showDialog: true});

  deleteUser = user => {
    this.props
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
          axios.delete('/nms/user/async/' + user.id).then(_resp =>
            this.setState({
              users: this.state.users.filter(u => u.id != user.id),
            }),
          );
        }
      });
  };

  createUser = payload => {
    axios
      .post('/nms/user/async/', payload)
      .then(response => this.updateUserState(response.data.user))
      .catch(error => this.props.alert(error.response?.data?.error || error));
  };

  editUser = (userId, payload) => {
    axios
      .put('/nms/user/async/' + userId, payload)
      .then(response => this.updateUserState(response.data.user))
      .catch(error => this.props.alert(error.response?.data?.error || error));
  };

  updateUserState = user =>
    this.setState(state => {
      const users = state.users.slice(0);
      if (this.state.editingUser) {
        const index = users.indexOf(state.editingUser);
        users[index] = user;
      } else {
        users.push(user);
      }
      return {editingUser: null, showDialog: false, users};
    });
}

export default withStyles(styles)(withAlert(UsersSettings));
