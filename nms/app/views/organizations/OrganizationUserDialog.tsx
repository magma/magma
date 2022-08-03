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

import type {DialogProps} from './OrganizationDialog';

import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';

import {AltFormField, PasswordInput} from '../../components/FormField';
import {UserRoles} from '../../../shared/roles';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles({
  addButton: {
    minWidth: '150px',
  },
  selectItem: {
    fontSize: '12px',
    fontFamily: '"Inter", sans-serif',
    fontWeight: 600,
  },
});

/**
 * Create User Tab
 * This component displays a form used to create a user that belongs to a new organization
 */
export default function OrganizationUserDialog(props: DialogProps) {
  const {user} = props;
  const classes = useStyles();

  return (
    <List>
      {props.error && (
        <AltFormField label={''}>
          <FormLabel error>{props.error}</FormLabel>
        </AltFormField>
      )}
      <AltFormField label={'Email'}>
        <OutlinedInput
          data-testid="email"
          placeholder="Email"
          fullWidth={true}
          disabled={user.id !== undefined}
          value={user.email || ''}
          onChange={({target}) => {
            props.onUserChange({...user, email: target.value});
          }}
        />
      </AltFormField>
      <AltFormField label={'Password'}>
        <PasswordInput
          data-testid="password"
          placeholder="Enter Password"
          value={user.password || ''}
          onChange={target => {
            props.onUserChange({...user, password: target});
          }}
        />
      </AltFormField>
      <AltFormField label={'Confirm Password'}>
        <PasswordInput
          data-testid="passwordConfirmation"
          placeholder="Enter Password Confirmation"
          value={user?.passwordConfirmation || ''}
          onChange={target => {
            props.onUserChange({...user, passwordConfirmation: target});
          }}
        />
      </AltFormField>
      <AltFormField
        label={'Role'}
        subLabel={
          'The role decides permissions that the user has to areas and features '
        }>
        <Select
          fullWidth={true}
          variant={'outlined'}
          value={user.role ?? 0}
          onChange={({target}) => {
            props.onUserChange({...user, role: target.value as number});
          }}
          input={<OutlinedInput id="direction" />}>
          <MenuItem key={UserRoles.USER} value={UserRoles.USER}>
            <ListItemText
              classes={{primary: classes.selectItem}}
              primary={'User'}
            />
          </MenuItem>
          <MenuItem
            key={UserRoles.READ_ONLY_USER}
            value={UserRoles.READ_ONLY_USER}>
            <ListItemText
              classes={{primary: classes.selectItem}}
              primary={'Read Only User'}
            />
          </MenuItem>
          <MenuItem key={UserRoles.SUPERUSER} value={UserRoles.SUPERUSER}>
            <ListItemText
              classes={{primary: classes.selectItem}}
              primary={'SuperUser'}
            />
          </MenuItem>
        </Select>
      </AltFormField>
    </List>
  );
}
