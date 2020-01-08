/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useCallback, useState} from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import {UserRoles} from '@fbcnms/auth/types';

import renderList from '@fbcnms/util/renderList';
import {makeStyles} from '@material-ui/styles';

export type EditUser = {
  id: string,
  email: string,
  role: number,
  networkIDs?: string[],
  organization?: string,
  readOnly: boolean,
};

type SaveUserData = {
  email: string,
  password?: string,
  superUser: boolean,
  networkIds?: string[],
  readOnly: boolean,
};

type Props = {
  editingUser: ?EditUser,
  open: boolean,
  onClose: () => void,
  ssoEnabled: boolean,
  allNetworkIDs?: Array<string>,
  onEditUser: (userId: string, payload: SaveUserData) => void,
  onCreateUser: (payload: SaveUserData) => void,
};

const useStyles = makeStyles({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
});

function getInitialNetworkIDs(userNetworkIds, allNetworkIDs): Set<string> {
  return new Set(allNetworkIDs && userNetworkIds ? userNetworkIds : []);
}

export default function EditUserDialog(props: Props) {
  const {allNetworkIDs} = props;
  const classes = useStyles();

  const [error, setError] = useState<string>('');
  const [email, setEmail] = useState<string>(props.editingUser?.email || '');
  const [password, setPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [isSuperUser, setIsSuperUser] = useState<boolean>(
    props.editingUser?.role === UserRoles.SUPERUSER,
  );
  const [isReadOnly, setIsReadOnly] = useState<boolean>(
    props.editingUser?.readOnly ?? false,
  );
  const [networkIds, setNetworkIds] = useState<Set<string>>(
    getInitialNetworkIDs(props.editingUser?.networkIDs, allNetworkIDs),
  );

  const onSave = useCallback(() => {
    if (password !== confirmPassword) {
      setError('Passwords must match');
      return;
    }

    if (!props.ssoEnabled && !props.editingUser && !password) {
      setError('Password cannot be empty');
      return;
    }

    if (!email) {
      setError('Email cannot be empty');
      return;
    }

    const payload: SaveUserData = {
      email,
      password,
      superUser: isSuperUser,
      networkIDs: isSuperUser ? [] : Array.from(networkIds),
      readOnly: isReadOnly,
    };

    // remove the password field if we are editing a user and the password isn't
    // being updated
    if ((props.editingUser || props.ssoEnabled) && !password) {
      delete payload.password;
    }

    if (props.editingUser) {
      props.onEditUser(props.editingUser.id, payload);
    } else {
      props.onCreateUser(payload);
    }
  }, [
    password,
    confirmPassword,
    props,
    email,
    isSuperUser,
    networkIds,
    isReadOnly,
  ]);

  return (
    <Dialog open={props.open} onClose={props.onClose}>
      <DialogTitle>{props.editingUser ? 'Edit User' : 'Add User'}</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        <TextField
          name="email"
          label="Email"
          className={classes.input}
          disabled={!!props.editingUser}
          value={email}
          onChange={({target}) => setEmail(target.value)}
        />
        {!props.ssoEnabled && (
          <>
            <TextField
              name="password"
              label="Password"
              type="password"
              value={password}
              onChange={({target}) => setPassword(target.value)}
              className={classes.input}
            />
            <TextField
              name="confirm_password"
              label="Confirm Password"
              type="password"
              value={confirmPassword}
              onChange={({target}) => setConfirmPassword(target.value)}
              className={classes.input}
            />
          </>
        )}
        <FormControlLabel
          control={
            <Checkbox
              checked={isSuperUser}
              onChange={({target}) => setIsSuperUser(target.checked)}
              color="primary"
            />
          }
          label="Super User"
        />
        <FormControlLabel
          control={
            <Checkbox
              checked={isReadOnly}
              onChange={({target}) => setIsReadOnly(target.checked)}
              color="primary"
            />
          }
          label="Read Only"
        />
        {allNetworkIDs && (
          <FormControl className={classes.input}>
            <InputLabel htmlFor="network_ids">Accessible Networks</InputLabel>
            <Select
              multiple
              disabled={isSuperUser}
              value={Array.from(networkIds)}
              onChange={({target}) => setNetworkIds(new Set(target.value))}
              renderValue={renderList}
              input={<Input id="network_ids" />}>
              {allNetworkIDs.map(network => (
                <MenuItem key={network} value={network}>
                  <Checkbox checked={networkIds.has(network)} />
                  <ListItemText primary={network} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} color="primary">
          Cancel
        </Button>
        <Button onClick={onSave} color="primary" variant="contained">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
