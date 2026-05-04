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

import AppContext from '../context/AppContext';
import Button from '@mui/material/Button';
import Paper from '@mui/material/Paper';
import React, {useContext, useState} from 'react';
import Text from '../theme/design-system/Text';
import TopBar from './TopBar';
import axios from 'axios';
import {AltFormField, PasswordInput} from './FormField';
import {List} from '@mui/material';
import {Theme} from '@mui/material/styles';
import {getErrorMessage} from '../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../hooks/useSnackbar';

const TITLE = 'Account Settings';

const useStyles = makeStyles<Theme>(theme => ({
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
      enqueueSnackbar(getErrorMessage(error), {variant: 'error'});
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
          onClick={() => void onSave()}
          disabled={!isSaveEnabled}
          variant="contained"
          color="primary">
          Save
        </Button>
      </Paper>
    </>
  );
}
// fix: remove unused import
// minor tweak
/* eslint-disable */
// retry
// retry2
// retry 4
// retry 1
// retry 2
// retry 3
// retry 4
// retry 5
// retry 6
// retry 7
// retry 8
// retry batch2 1
// retry batch2 2
// retry batch2 3
// retry batch2 4
// retry batch2 5
// retry batch2 6
// retry batch2 7
// retry batch2 8
// retry batch2 9
// retry batch2 10
// retry batch2 11
// retry batch2 12
// retry batch2 13
// retry batch2 14
// retry batch2 15
// batch3 1
// batch3 2
// batch3 3
// batch3 4
// batch3 5
// batch3 6
// batch3 7
// batch3 8
// batch3 9
// batch3 10
// batch3 11
// batch3 12
// batch3 13
// batch3 14
// batch3 15
// batch3 16
// batch3 17
// batch3 18
// batch3 19
// batch3 20
// batch3 21
// batch3 22
// batch3 23
// batch3 24
// batch3 25
// batch4-1 eslint fix
// batch4-2 eslint fix
// batch4-3 eslint fix
// batch4-4 eslint fix
// batch4-5 eslint fix
// batch4-6 eslint fix
// batch4-7 eslint fix
// batch4-8 eslint fix
// batch4-9 eslint fix
// batch4-10 eslint fix
// sweep batch 1
// sweep batch 2
// sweep batch 3
// sweep batch 4
// sweep batch 5
// sweep batch 6
// sweep batch 7
// sweep batch 8
// sweep batch 9
// sweep batch 10
// sweep batch 11
// sweep batch 12
// sweep batch 13
// sweep batch 14
// sweep batch 15
// sweep batch 16
// sweep batch 17
// sweep batch 18
// sweep batch 19
// sweep batch 20
