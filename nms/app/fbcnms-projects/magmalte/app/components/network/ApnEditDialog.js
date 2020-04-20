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

import type {aggregated_maximum_bitrate, apn} from '@fbcnms/magma-api';

import Button from '@fbcnms/ui/components/design-system/Button';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormGroup,
  InputAdornment,
  TextField,
} from '@material-ui/core';

import nullthrows from '@fbcnms/util/nullthrows';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './ApnConst';

// Calculates the bitrate value in bps based on the default
// value input in the form, and the previous (default) bitrate
function _getBitRateValue(props: {
  apnAmbr: ?aggregated_maximum_bitrate,
  editedValue: ?string,
  field: 'max_bandwidth_ul' | 'max_bandwidth_dl',
}): number {
  const {apnAmbr, editedValue, field} = props;
  if (editedValue !== null && editedValue === '') {
    return DATA_PLAN_UNLIMITED_RATES[field];
  } else if (editedValue !== null) {
    return parseFloat(editedValue) * BITRATE_MULTIPLIER;
  }
  return apnAmbr ? apnAmbr[field] : DATA_PLAN_UNLIMITED_RATES[field];
}

type MegabyteFieldProps = {
  label: string,
  field: 'max_bandwidth_ul' | 'max_bandwidth_dl',
  value: ?string,
  apnAmbr: ?aggregated_maximum_bitrate,
  onChange: (event: SyntheticInputEvent<*>) => void,
};

function MegabyteTextField(props: MegabyteFieldProps) {
  const {apnAmbr, field, label, onChange, value} = props;
  const defaultMaxBitRate = apnAmbr ? apnAmbr[field] : null;
  let fieldValue;
  if (value !== null) {
    fieldValue = value || '';
  } else if (
    defaultMaxBitRate === null ||
    defaultMaxBitRate === DATA_PLAN_UNLIMITED_RATES[field]
  ) {
    fieldValue = '';
  } else {
    fieldValue = defaultMaxBitRate / BITRATE_MULTIPLIER;
  }

  return (
    <TextField
      required
      type="number"
      label={label}
      margin="normal"
      placeholder="Unlimited"
      value={fieldValue}
      InputProps={{
        endAdornment: <InputAdornment position="end">Mbps</InputAdornment>,
      }}
      InputLabelProps={{
        shrink: true,
      }}
      onChange={onChange}
    />
  );
}

type Props = {
  onSave: () => void,
  onCancel: () => void,
  apnConfig: ?apn,
  apnName: ?string,
};

export default function ApnEditDialog(props: Props) {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [editedName, setEditedName] = useState(props.apnName || '');
  const [editedMaxDlBitRate, setEditedMaxDlBitRate] = useState(null);
  const [editedMaxUlBitRate, setEditedMaxUlBitRate] = useState(null);
  const [editedClassID, setEditedClassID] = useState(9);
  const [editedPriorityLevel, setEditedPriorityLevel] = useState(15);
  const [editedPreemptionCapability, setEditedPreemptionCapability] = useState(
    0,
  );
  const [
    editedPreemptionVulnerability,
    setEditedPreemptionVulnerability,
  ] = useState(0);

  const {apnConfig} = props;

  const apnConfiguration =
    (props.apnName && apnConfig?.apn_configuration) || null;

  const onSave = () => {
    if (!props.apnName && apnConfig?.apn_name == editedName) {
      enqueueSnackbar('APN is already used. Please use a different name', {
        variant: 'error',
      });
      return;
    }

    const apnAmbr = {
      max_bandwidth_dl: _getBitRateValue({
        apnAmbr: apnConfiguration?.ambr,
        field: 'max_bandwidth_dl',
        editedValue: editedMaxDlBitRate,
      }),
      max_bandwidth_ul: _getBitRateValue({
        apnAmbr: apnConfiguration?.ambr,
        field: 'max_bandwidth_ul',
        editedValue: editedMaxUlBitRate,
      }),
    };

    const apnQosProfile = {
      class_id: parseInt(editedClassID),
      preemption_capability: !!editedPreemptionCapability,
      preemption_vulnerability: !!editedPreemptionVulnerability,
      priority_level: parseInt(editedPriorityLevel),
    };

    const newApnConfig = {
      ambr: apnAmbr,
      qos_profile: apnQosProfile,
    };

    const newApn = {
      apn_name: editedName,
      apn_configuration: newApnConfig,
    };

    if (props.apnName == editedName) {
      MagmaV1API.putLteByNetworkIdApnsByApnName({
        networkId: nullthrows(match.params.networkId),
        apnName: editedName,
        apn: newApn,
      })
        .then(_resp => props.onSave())
        .catch(error =>
          enqueueSnackbar(error.response?.data?.message || error, {
            variant: 'error',
          }),
        );
    } else {
      MagmaV1API.postLteByNetworkIdApns({
        networkId: nullthrows(match.params.networkId),
        apn: newApn,
      })
        .then(_resp => props.onSave())
        .catch(error =>
          enqueueSnackbar(error.response?.data?.message || error, {
            variant: 'error',
          }),
        );
    }
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>
        {props.apnName ? 'Edit' : 'Add'} APN Configuration
      </DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            label="Name"
            margin="normal"
            disabled={!!props.apnName}
            value={editedName}
            onChange={({target}) => setEditedName(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <MegabyteTextField
            label="AMBR DL"
            apnAmbr={apnConfig?.apn_configuration?.ambr}
            field="max_bandwidth_dl"
            value={editedMaxDlBitRate}
            onChange={({target}) => setEditedMaxDlBitRate(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <MegabyteTextField
            label="AMBR UL"
            apnAmbr={apnConfig?.apn_configuration?.ambr}
            field="max_bandwidth_ul"
            value={editedMaxUlBitRate}
            onChange={({target}) => setEditedMaxUlBitRate(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            type="number"
            label="QoS Class ID"
            margin="normal"
            min="0"
            max="255"
            value={editedClassID}
            onChange={({target}) => setEditedClassID(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            type="number"
            label="Priority Level"
            margin="normal"
            min="0"
            max="15"
            value={editedPriorityLevel}
            onChange={({target}) => setEditedPriorityLevel(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            type="number"
            label="Preemption Cap."
            margin="normal"
            min="0"
            max="1"
            value={editedPreemptionCapability}
            onChange={({target}) => setEditedPreemptionCapability(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <TextField
            required
            type="number"
            label="Preemption Vul."
            margin="normal"
            min="0"
            max="1"
            value={editedPreemptionVulnerability}
            onChange={({target}) =>
              setEditedPreemptionVulnerability(target.value)
            }
          />
        </FormGroup>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onCancel} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave}>Save</Button>
      </DialogActions>
    </Dialog>
  );
}
