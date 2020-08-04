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
 * @flow
 * @format
 */

import type {apn_list, subscriber} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import InputLabel from '@material-ui/core/InputLabel';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import TypedSelect from '@fbcnms/ui/components/TypedSelect';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import nullthrows from '@fbcnms/util/nullthrows';
import {base64ToHex, hexToBase64, isValidHex} from '@fbcnms/util/strings';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type EditingSubscriber = {
  imsiID: string,
  lteState: 'ACTIVE' | 'INACTIVE',
  authKey: string,
  authOpc: string,
  subProfile: string,
  apnList: apn_list,
};

type Props = {
  onClose: () => void,
  onSave: (subscriberID: string) => void,
  onSaveError: (reason: string) => void,
  editingSubscriber?: subscriber,
  subProfiles: Array<string>,
  apns: apn_list,
};

function buildEditingSubscriber(
  editingSubscriber: ?subscriber,
): EditingSubscriber {
  if (!editingSubscriber) {
    return {
      imsiID: '',
      lteState: 'ACTIVE',
      authKey: '',
      authOpc: '',
      subProfile: 'default',
      apnList: [],
    };
  }

  const authKey = editingSubscriber.lte.auth_key
    ? base64ToHex(editingSubscriber.lte.auth_key)
    : '';

  const authOpc =
    editingSubscriber.lte.auth_opc != undefined
      ? base64ToHex(editingSubscriber.lte.auth_opc)
      : '';

  return {
    imsiID: editingSubscriber.id,
    lteState: editingSubscriber.lte.state,
    authKey,
    authOpc,
    subProfile: editingSubscriber.lte.sub_profile,
    apnList: editingSubscriber.active_apns || [],
  };
}

export default function AddEditSubscriberDialog(props: Props) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [editingSubscriber, setEditingSubscriber] = useState(
    buildEditingSubscriber(props.editingSubscriber),
  );

  const onSave = () => {
    if (!editingSubscriber.imsiID || !editingSubscriber.authKey) {
      enqueueSnackbar('Please complete all fields', {variant: 'error'});
      return;
    }

    let {imsiID} = editingSubscriber;
    if (!imsiID.startsWith('IMSI')) {
      imsiID = `IMSI${imsiID}`;
    }

    const data = {
      id: imsiID,
      lte: {
        state: editingSubscriber.lteState,
        auth_algo: 'MILENAGE', // default auth algo
        auth_key: editingSubscriber.authKey,
        auth_opc: editingSubscriber.authOpc || undefined,
        sub_profile: editingSubscriber.subProfile,
      },
      active_apns: editingSubscriber.apnList,
    };
    if (data.lte.auth_key && isValidHex(data.lte.auth_key)) {
      data.lte.auth_key = hexToBase64(data.lte.auth_key);
    }
    if (data.lte.auth_opc != undefined && isValidHex(data.lte.auth_opc)) {
      data.lte.auth_opc = hexToBase64(data.lte.auth_opc);
    }
    if (props.editingSubscriber) {
      MagmaV1API.putLteByNetworkIdSubscribersBySubscriberId({
        networkId: nullthrows(match.params.networkId),
        subscriberId: data.id,
        subscriber: data,
      })
        .then(() => props.onSave(data.id))
        .catch(e => props.onSaveError(e.response.data.message));
    } else {
      MagmaV1API.postLteByNetworkIdSubscribers({
        networkId: match.params.networkId || '',
        subscriber: data,
      })
        .then(() => props.onSave(data.id))
        .catch(e => props.onSaveError(e.response.data.message));
    }
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>
        {props.editingSubscriber ? 'Edit Subscriber' : 'Add Subscriber'}
      </DialogTitle>
      <DialogContent>
        <TextField
          label="IMSI"
          className={classes.input}
          disabled={!!props.editingSubscriber}
          value={editingSubscriber.imsiID}
          onChange={({target}) =>
            setEditingSubscriber({
              ...editingSubscriber,
              imsiID: target.value,
            })
          }
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="lteState">LTE Subscription State</InputLabel>
          <TypedSelect
            inputProps={{id: 'lteState'}}
            value={editingSubscriber.lteState}
            items={{
              ACTIVE: 'Active',
              INACTIVE: 'Inactive',
            }}
            onChange={lteState =>
              setEditingSubscriber({...editingSubscriber, lteState})
            }
          />
        </FormControl>
        <TextField
          label="LTE Auth Key"
          className={classes.input}
          value={editingSubscriber.authKey}
          onChange={({target}) =>
            setEditingSubscriber({
              ...editingSubscriber,
              authKey: target.value,
            })
          }
        />
        <TextField
          label="LTE Auth OPc"
          className={classes.input}
          value={editingSubscriber.authOpc}
          onChange={({target}) =>
            setEditingSubscriber({
              ...editingSubscriber,
              authOpc: target.value,
            })
          }
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="subProfile">Data Plan</InputLabel>
          <Select
            inputProps={{id: 'subProfile'}}
            value={editingSubscriber.subProfile}
            onChange={({target}) =>
              setEditingSubscriber({
                ...editingSubscriber,
                subProfile: target.value,
              })
            }>
            {props.subProfiles.map(p => (
              <MenuItem value={p} key={p}>
                {p}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        <FormControl className={classes.input}>
          <InputLabel htmlFor="apnList">Access Point Names</InputLabel>
          <Select
            inputProps={{id: 'apnList'}}
            value={editingSubscriber.apnList}
            multiple={true}
            onChange={({target}) =>
              setEditingSubscriber({
                ...editingSubscriber,
                apnList: ((target.value: any): apn_list),
              })
            }>
            {props.apns.map(apn => (
              <MenuItem value={apn} key={apn}>
                {apn}
              </MenuItem>
            ))}
          </Select>
        </FormControl>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
