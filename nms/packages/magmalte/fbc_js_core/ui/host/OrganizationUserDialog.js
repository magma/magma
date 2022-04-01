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

import type {DialogProps} from './OrganizationDialog';

import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';

import {AltFormField} from '../../../fbc_js_core/ui/components/design-system/FormField/FormField';
import {UserRoles} from '../../../fbc_js_core/auth/types';

/**
 * Create User Tab
 * This component displays a form used to create a user that belongs to a new organization
 */
export default function OrganizationUserDialog(props: DialogProps) {
  const {user} = props;

  return (
    <DialogContent>
      <List>
        {props.error && (
          <AltFormField label={''}>
            <FormLabel error>{props.error}</FormLabel>
          </AltFormField>
        )}
        <AltFormField disableGutters label={'Email'}>
          <OutlinedInput
            data-testid="name"
            placeholder="Email"
            fullWidth={true}
            value={user.email}
            onChange={({target}) => {
              props.onUserChange({...user, email: target.value});
            }}
          />
        </AltFormField>
        <AltFormField disableGutters label={'Password'}>
          <OutlinedInput
            data-testid="name"
            placeholder="Enter Password"
            fullWidth={true}
            value={user.password}
            onChange={({target}) => {
              props.onUserChange({...user, password: target.value});
            }}
          />
        </AltFormField>
        <AltFormField disableGutters label={'Confirm Password'}>
          <OutlinedInput
            data-testid="name"
            placeholder="Enter Password Confirmation"
            fullWidth={true}
            value={user.passwordConfirmation}
            onChange={({target}) => {
              props.onUserChange({...user, passwordConfirmation: target.value});
            }}
          />
        </AltFormField>
        <AltFormField
          disableGutters
          label={'Role'}
          subLabel={
            'The role decides permissions that the user has to areas and features '
          }>
          <Select
            fullWidth={true}
            variant={'outlined'}
            value={user.role ?? 0}
            onChange={({target}) => {
              props.onUserChange({...user, role: target.value});
            }}
            input={<OutlinedInput id="direction" />}>
            <MenuItem key={UserRoles.USER} value={UserRoles.USER}>
              <ListItemText primary={'User'} />
            </MenuItem>
            <MenuItem
              key={UserRoles.READ_ONLY_USER}
              value={UserRoles.READ_ONLY_USER}>
              <ListItemText primary={'Read Only User'} />
            </MenuItem>
            <MenuItem key={UserRoles.SUPERUSER} value={UserRoles.SUPERUSER}>
              <ListItemText primary={'SuperUser'} />
            </MenuItem>
          </Select>
        </AltFormField>
      </List>
    </DialogContent>
  );
}
