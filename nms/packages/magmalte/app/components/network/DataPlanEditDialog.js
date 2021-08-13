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

import type {network_epc_configs} from '@fbcnms/magma-api';

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

import {BITRATE_MULTIPLIER, DATA_PLAN_UNLIMITED_RATES} from './DataPlanConst';

type CellularNetworkProfile = $Values<
  $NonMaybeType<$PropertyType<network_epc_configs, 'sub_profiles'>>,
>;

// Calculates the bitrate value in bps based on the default
// value input in the form, and the previous (default) bitrate
function _getBitRateValue(props: {
  dataPlan: ?CellularNetworkProfile,
  editedValue: ?string,
  field: 'max_ul_bit_rate' | 'max_dl_bit_rate',
}): number {
  const {dataPlan, editedValue, field} = props;
  if (editedValue !== null && editedValue === '') {
    return DATA_PLAN_UNLIMITED_RATES[field];
  } else if (editedValue !== null) {
    return parseFloat(editedValue) * BITRATE_MULTIPLIER;
  }
  return dataPlan ? dataPlan[field] : DATA_PLAN_UNLIMITED_RATES[field];
}

const MegabyteTextField = (props: {
  label: string,
  field: 'max_ul_bit_rate' | 'max_dl_bit_rate',
  value: ?string,
  dataPlan: ?CellularNetworkProfile,
  onChange: (event: SyntheticInputEvent<*>) => void,
}) => {
  const {dataPlan, field, label, onChange, value} = props;
  const defaultMaxBitRate = dataPlan ? dataPlan[field] : null;
  let fieldValue;
  if (value !== null) {
    fieldValue = value || '';
  } else if (defaultMaxBitRate === null) {
    fieldValue = '';
  } else if (defaultMaxBitRate === DATA_PLAN_UNLIMITED_RATES[field]) {
    fieldValue = '';
  } else {
    fieldValue = (defaultMaxBitRate || 0) / BITRATE_MULTIPLIER + '';
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
};

type Props = {
  onSave: (dataPlanId: string, newEPCConfig: network_epc_configs) => void,
  onCancel: () => void,
  epcConfig: network_epc_configs,
  dataPlanId: ?string,
};

export default function DataPlanEditDialog(props: Props) {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [editedName, setEditedName] = useState(props.dataPlanId || '');
  const [editedMaxDlBitRate, setEditedMaxDlBitRate] = useState(null);
  const [editedMaxUlBitRate, setEditedMaxUlBitRate] = useState(null);

  const {epcConfig} = props;
  const dataPlan =
    (props.dataPlanId && epcConfig.sub_profiles?.[props.dataPlanId]) || null;

  const onSave = () => {
    if (!props.dataPlanId && (epcConfig.sub_profiles || {})[editedName]) {
      enqueueSnackbar(
        'Data plan name is already used. Please use a different name',
        {variant: 'error'},
      );
      return;
    }

    const subProfiles = epcConfig.sub_profiles
      ? {...epcConfig.sub_profiles}
      : {};

    subProfiles[editedName] = {
      max_dl_bit_rate: _getBitRateValue({
        dataPlan,
        field: 'max_dl_bit_rate',
        editedValue: editedMaxDlBitRate,
      }),
      max_ul_bit_rate: _getBitRateValue({
        dataPlan,
        field: 'max_ul_bit_rate',
        editedValue: editedMaxUlBitRate,
      }),
    };

    const newConfig = {
      ...epcConfig,
      sub_profiles: subProfiles,
    };

    MagmaV1API.putLteByNetworkIdCellularEpc({
      networkId: nullthrows(match.params.networkId),
      config: newConfig,
    })
      .then(_resp => props.onSave(editedName, newConfig))
      .catch(error =>
        enqueueSnackbar(error.response?.data?.message || error, {
          variant: 'error',
        }),
      );
  };

  return (
    <Dialog open={true} onClose={props.onCancel} scroll="body">
      <DialogTitle>{props.dataPlanId ? 'Edit' : 'Add'} Data Plan</DialogTitle>
      <DialogContent>
        <FormGroup row>
          <TextField
            required
            label="Name"
            margin="normal"
            disabled={!!props.dataPlanId}
            value={editedName}
            onChange={({target}) => setEditedName(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <MegabyteTextField
            label="Download Limit"
            dataPlan={dataPlan}
            field="max_dl_bit_rate"
            value={editedMaxDlBitRate}
            onChange={({target}) => setEditedMaxDlBitRate(target.value)}
          />
        </FormGroup>
        <FormGroup row>
          <MegabyteTextField
            label="Upload Limit"
            dataPlan={dataPlan}
            field="max_ul_bit_rate"
            value={editedMaxUlBitRate}
            onChange={({target}) => setEditedMaxUlBitRate(target.value)}
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
