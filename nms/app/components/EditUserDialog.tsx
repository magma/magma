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

import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';
import ListItemText from '@mui/material/ListItemText';
import MenuItem from '@mui/material/MenuItem';
import React, {useCallback, useMemo, useState} from 'react';
import Select from '@mui/material/Select';
import renderList from '../util/renderList';
import {AltFormField} from './FormField';
import {NetworkId} from '../../shared/types/network';
import {OutlinedInput} from '@mui/material';
import {UserRoles} from '../../shared/roles';

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

function getInitialNetworkIDs(
  userNetworkIds: Array<NetworkId> | undefined,
  allNetworkIDs: Array<NetworkId>,
): Set<string> {
  return new Set(allNetworkIDs && userNetworkIds ? userNetworkIds : []);
}

export default function EditUserDialog(props: Props) {
  const {allNetworkIDs} = props;

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
        <AltFormField label="Email">
          <OutlinedInput
            name="email"
            fullWidth
            disabled={!!props.editingUser}
            value={email}
            onChange={({target}) => setEmail(target.value)}
          />
        </AltFormField>
        {!props.ssoEnabled && (
          <>
            <AltFormField label="Password">
              <OutlinedInput
                autoComplete="off"
                name="password"
                type="password"
                fullWidth
                value={password}
                onChange={({target}) => setPassword(target.value)}
              />
            </AltFormField>
            <AltFormField label="Confirm Password">
              <OutlinedInput
                autoComplete="off"
                name="confirm_password"
                type="password"
                fullWidth
                value={confirmPassword}
                onChange={({target}) => setConfirmPassword(target.value)}
              />
            </AltFormField>
          </>
        )}
        <AltFormField label="Role">
          <FormControl fullWidth>
            <Select
              labelId="role-select-label"
              id="role-select"
              value={role}
              onChange={({target}) =>
                setRole(parseInt(target.value as string))
              }>
              <MenuItem value={UserRoles.USER}>User</MenuItem>
              <MenuItem value={UserRoles.READ_ONLY_USER}>
                Read Only User
              </MenuItem>
              <MenuItem value={UserRoles.SUPERUSER}>Super User</MenuItem>
            </Select>
          </FormControl>
        </AltFormField>
        {allNetworkIDs && (
          <AltFormField label="Accessible Networks">
            <FormControl fullWidth>
              <Select
                multiple
                disabled={isSuperUser}
                value={Array.from(networkIds)}
                onChange={({target}) =>
                  setNetworkIds(new Set(target.value as Array<string>))
                }
                renderValue={networkIds => renderList(networkIds)}
                input={<OutlinedInput id="network_ids" />}>
                {allNetworkIDs.map(network => (
                  <MenuItem key={network} value={network}>
                    <Checkbox checked={networkIds.has(network)} />
                    <ListItemText primary={network} />
                  </MenuItem>
                ))}
              </Select>
            </FormControl>
          </AltFormField>
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
