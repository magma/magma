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
import AppContext from './context/AppContext';
import Button from '@material-ui/core/Button';
import Paper from '@material-ui/core/Paper';
import React, {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../theme/design-system/Text';
import TopBar from './TopBar';
import axios from 'axios';
// $FlowFixMe migrated to typescript
import {AltFormField, PasswordInput} from './FormField';
import {List} from '@material-ui/core';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../app/hooks/useSnackbar';

const TITLE = 'Account Settings';

const useStyles = makeStyles(theme => ({
  title: {
    fontSize: '18px',
  },
  input: {
    width: '100%',
    maxWidth: '400px',
  },
  formContainer: {
    paddingBottom: theme.spacing(2),
  },
  paper: {
    margin: theme.spacing(4),
    padding: theme.spacing(3),
    paddingBottom: theme.spacing(6),
  },
}));

export default function AccountSettings() {
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [currentPassword, setCurrentPassword] = useState('');
  const [newPassword, setNewPassword] = useState('');
  const [confirmPassword, setConfirmPassword] = useState('');
  const {isOrganizations} = useContext(AppContext);

  const isSaveEnabled = currentPassword && newPassword && confirmPassword;

  const onSave = async () => {
    if (newPassword !== confirmPassword) {
      enqueueSnackbar('Passwords do not match', {variant: 'error'});
      return;
    }

    try {
      await axios.post('/user/change_password', {
        currentPassword: currentPassword,
        newPassword: newPassword,
      });

      enqueueSnackbar('Success', {variant: 'success'});
      setCurrentPassword('');
      setNewPassword('');
      setConfirmPassword('');
    } catch (error) {
      enqueueSnackbar(error.response.data.error, {variant: 'error'});
    }
  };

  return (
    <>
      {!isOrganizations && <TopBar header={TITLE} tabs={[]} />}
      <Paper className={classes.paper}>
        <Text data-testid="change-password-title" variant="body1">
          Change Password
        </Text>
        <List className={classes.formContainer}>
          <AltFormField label="Current Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              placeholder="Enter Current Password"
              value={currentPassword}
              onChange={setCurrentPassword}
            />
          </AltFormField>

          <AltFormField label="New Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              autoComplete="off"
              placeholder="Enter New Password"
              value={newPassword}
              onChange={setNewPassword}
            />
          </AltFormField>

          <AltFormField label="Confirm New Password" disableGutters={true}>
            <PasswordInput
              className={classes.input}
              required
              autoComplete="off"
              placeholder="Confirm New Password"
              value={confirmPassword}
              onChange={setConfirmPassword}
            />
          </AltFormField>
        </List>
        <Button
          onClick={onSave}
          disabled={!isSaveEnabled}
          variant="contained"
          color="primary">
          Save
        </Button>
      </Paper>
    </>
  );
}
