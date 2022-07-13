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

import Button from '@material-ui/core/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useCallback, useMemo, useState} from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import renderList from '../util/renderList';
import {NetworkId} from '../../shared/types/network';
import {UserRoles} from '../../shared/roles';
import {makeStyles} from '@material-ui/styles';

export type EditUser = {
  id: string;
  email: string;
  role: number;
  networkIDs?: Array<string>;
  organization?: string;
};

export type SaveUserData = {
  email: string;
  password?: string;
  role: number;
  networkIDs?: Array<string>;
};

type Props = {
  editingUser: EditUser | null | undefined;
  open: boolean;
  onClose: () => void;
  ssoEnabled: boolean;
  allNetworkIDs: Array<string>;
  onEditUser: (userId: string, payload: SaveUserData) => void;
  onCreateUser: (payload: SaveUserData) => void;
};

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
  select: {
    marginTop: '16px',
  },
}));

function getInitialNetworkIDs(
  userNetworkIds: Array<NetworkId> | undefined,
  allNetworkIDs: Array<NetworkId>,
): Set<string> {
  return new Set(allNetworkIDs && userNetworkIds ? userNetworkIds : []);
}

export default function EditUserDialog(props: Props) {
  const {allNetworkIDs} = props;
  const classes = useStyles();

  const [error, setError] = useState<string>('');
  const [email, setEmail] = useState<string>(props.editingUser?.email || '');
  const [password, setPassword] = useState<string>('');
  const [confirmPassword, setConfirmPassword] = useState<string>('');
  const [role, setRole] = useState<typeof UserRoles[keyof typeof UserRoles]>(
    props.editingUser?.role ?? UserRoles.USER,
  );
  const [networkIds, setNetworkIds] = useState<Set<string>>(
    getInitialNetworkIDs(props.editingUser?.networkIDs, allNetworkIDs),
  );
  const isSuperUser = useMemo(() => role === UserRoles.SUPERUSER, [role]);

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
      role,
      networkIDs: isSuperUser ? [] : Array.from(networkIds),
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
  }, [password, confirmPassword, props, email, role, isSuperUser, networkIds]);

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
              autoComplete="off"
              name="password"
              label="Password"
              type="password"
              value={password}
              onChange={({target}) => setPassword(target.value)}
              className={classes.input}
            />
            <TextField
              autoComplete="off"
              name="confirm_password"
              label="Confirm Password"
              type="password"
              value={confirmPassword}
              onChange={({target}) => setConfirmPassword(target.value)}
              className={classes.input}
            />
          </>
        )}
        <FormControl className={classes.input}>
          <InputLabel id="role-select-label">Role</InputLabel>
          <Select
            labelId="role-select-label"
            id="role-select"
            value={role}
            onChange={({target}) => setRole(parseInt(target.value as string))}>
            <MenuItem value={UserRoles.USER}>User</MenuItem>
            <MenuItem value={UserRoles.READ_ONLY_USER}>Read Only User</MenuItem>
            <MenuItem value={UserRoles.SUPERUSER}>Super User</MenuItem>
          </Select>
        </FormControl>
        {allNetworkIDs && (
          <FormControl className={classes.input}>
            <InputLabel htmlFor="network_ids">Accessible Networks</InputLabel>
            <Select
              multiple
              disabled={isSuperUser}
              value={Array.from(networkIds)}
              onChange={({target}) =>
                setNetworkIds(new Set(target.value as Array<string>))
              }
              renderValue={networkIds =>
                renderList(networkIds as Array<string>)
              }
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
