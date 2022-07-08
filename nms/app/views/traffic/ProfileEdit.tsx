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

import type {PolicyQosProfile} from '../../../generated-ts';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import PolicyContext from '../../components/context/PolicyContext';
import React from 'react';

import {AltFormField} from '../../components/FormField';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

const useStyles = makeStyles({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
});

type Props = {
  open: boolean;
  onClose: () => void;
  profile?: PolicyQosProfile;
};

export default function ProfileEditDialog(props: Props) {
  const classes = useStyles();
  const ctx = useContext(PolicyContext);
  const qosProfiles = ctx.qosProfiles;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');

  const [profile, setProfile] = useState<PolicyQosProfile>(
    props.profile || {
      id: '',
      class_id: 0,
      max_req_bw_dl: 9,
      max_req_bw_ul: 9,
    },
  );

  useEffect(() => {
    setProfile(
      props.profile || {
        id: '',
        class_id: 0,
        max_req_bw_dl: 9,
        max_req_bw_ul: 9,
      },
    );
    setError('');
  }, [props.open, props.profile]);

  const isAdd = props.profile ? false : true;
  const handleProfileChange = (key: string, val: string | number) => {
    setProfile({...profile, [key]: val});
  };
  const onSave = async () => {
    try {
      if (isAdd) {
        if (profile.id === '') {
          setError('empty profile id');
          return;
        }
        if (profile.id in qosProfiles) {
          setError(`Profile ${profile.id} already exists`);
          return;
        }
      }
      await ctx.setQosProfiles(profile.id, profile);
      enqueueSnackbar('Profile saved successfully', {
        variant: 'success',
      });

      props.onClose();
    } catch (error) {
      setError(getErrorMessage(error));
    }
  };

  const onClose = () => {
    props.onClose();
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      scroll="body"
      fullWidth={true}
      maxWidth={'md'}>
      <DialogTitle
        onClose={onClose}
        label={props.profile ? 'Edit Profile' : 'Add New Profile'}
      />
      <DialogContent>
        <List>
          {error !== '' && (
            <AltFormField disableGutters label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <ListItem dense disableGutters />
          <AltFormField label={'Profile ID'} disableGutters>
            <OutlinedInput
              className={classes.input}
              fullWidth={true}
              data-testid="profileID"
              placeholder="test_profile"
              value={profile.id}
              onChange={({target}) => handleProfileChange('id', target.value)}
            />
          </AltFormField>
          <AltFormField label={'Class ID'} disableGutters>
            <OutlinedInput
              className={classes.input}
              type="number"
              fullWidth={true}
              data-testid="profileClassID"
              placeholder="9"
              value={profile.class_id}
              onChange={({target}) =>
                handleProfileChange('class_id', parseInt(target.value) || '')
              }
            />
          </AltFormField>
          <AltFormField label={'Max Bandwidth Downlink(bps)'} disableGutters>
            <OutlinedInput
              className={classes.input}
              type="number"
              fullWidth={true}
              inputProps={{min: 0}}
              data-testid="maxReqBwDl"
              placeholder="1000"
              value={profile.max_req_bw_dl}
              onChange={({target}) =>
                handleProfileChange(
                  'max_req_bw_dl',
                  parseInt(target.value) || '',
                )
              }
            />
          </AltFormField>
          <AltFormField label={'Max Bandwidth Uplink(bps)'} disableGutters>
            <OutlinedInput
              className={classes.input}
              type="number"
              fullWidth={true}
              data-testid="maxReqBwUl"
              placeholder="1000"
              value={profile.max_req_bw_ul}
              onChange={({target}) =>
                handleProfileChange(
                  'max_req_bw_ul',
                  parseInt(target.value) || '',
                )
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Close</Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => void onSave()}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
