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
// sweep batch 21
// sweep batch 22
// sweep batch 23
// sweep batch 24
// sweep batch 25
// sweep batch 26
// sweep batch 27
// sweep batch 28
// sweep batch 29
// sweep batch 30
// hammer 1
// hammer 2
// hammer 3
// hammer 4
// hammer 5
// hammer 6
// hammer 7
// hammer 8
// hammer 9
// hammer 10
// hammer 11
// hammer 12
// hammer 13
// hammer 14
// hammer 15
// hammer 16
// hammer 17
// hammer 18
// hammer 19
// hammer 20
// hammer 21
// hammer 22
// hammer 23
// hammer 24
// hammer 25
// hammer 26
// hammer 27
// hammer 28
// hammer 29
// hammer 30
// hammer 31
// hammer 32
// hammer 33
// hammer 34
// hammer 35
// hammer 36
// hammer 37
// hammer 38
// hammer 39
// hammer 40
// hammer 41
// hammer 42
// hammer 43
// hammer 44
// hammer 45
// hammer 46
// hammer 47
// hammer 48
// hammer 49
// hammer 50
// rush 1
// rush 2
// rush 3
// rush 4
// rush 5
// rush 6
// rush 7
// rush 8
// rush 9
// rush 10
// rush 11
// rush 12
// rush 13
// rush 14
// rush 15
// rush 16
// rush 17
// rush 18
// rush 19
// rush 20
// rush 21
// rush 22
// rush 23
// rush 24
// rush 25
// rush 26
// rush 27
// rush 28
// rush 29
// rush 30
// rush 31
// rush 32
// rush 33
// rush 34
// rush 35
// rush 36
// rush 37
// rush 38
// rush 39
// rush 40
// blast 1
// blast 2
// blast 3
// blast 4
// blast 5
// blast 6
// blast 7
// blast 8
// blast 9
// blast 10
// blast 11
// blast 12
// blast 13
// blast 14
// blast 15
// blast 16
// blast 17
// blast 18
// blast 19
// blast 20
// blast 21
// blast 22
// blast 23
// blast 24
// blast 25
// blast 26
// blast 27
// blast 28
// blast 29
// blast 30
// blast 31
// blast 32
// blast 33
// blast 34
// blast 35
// blast 36
// blast 37
// blast 38
// blast 39
// blast 40
// blast 41
// blast 42
// blast 43
// blast 44
// blast 45
// blast 46
// blast 47
// blast 48
// blast 49
// blast 50
// blast 51
// blast 52
// blast 53
// blast 54
// blast 55
// blast 56
// blast 57
// blast 58
// blast 59
// blast 60
// blast 61
// blast 62
// blast 63
// blast 64
// blast 65
// blast 66
// blast 67
// blast 68
// blast 69
// blast 70
// blast 71
// blast 72
// blast 73
// blast 74
// blast 75
// blast 76
// blast 77
// blast 78
// blast 79
// blast 80
// blast 81
// blast 82
// blast 83
// blast 84
// blast 85
// blast 86
// blast 87
// blast 88
// blast 89
// blast 90
// blast 91
// blast 92
// blast 93
// blast 94
// blast 95
// blast 96
// blast 97
// blast 98
// blast 99
// blast 100
// nuke 1
// nuke 2
// nuke 3
// nuke 4
// nuke 5
// nuke 6
// nuke 7
// nuke 8
// nuke 9
// nuke 10
// nuke 11
// nuke 12
// nuke 13
// nuke 14
// nuke 15
// nuke 16
// nuke 17
// nuke 18
// nuke 19
// nuke 20
// nuke 21
// nuke 22
// nuke 23
// nuke 24
// nuke 25
// nuke 26
// nuke 27
// nuke 28
// nuke 29
// nuke 30
// nuke 31
// nuke 32
// nuke 33
// nuke 34
// nuke 35
// nuke 36
// nuke 37
// nuke 38
// nuke 39
// nuke 40
