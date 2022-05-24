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

import type {network_epc_configs} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import InputAdornment from '@material-ui/core/InputAdornment';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import {AltFormField, AltFormFieldSubheading} from '../../components/FormField';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
// $FlowFixMe migrated to typescript
import type {UpdateNetworkContextProps} from '../../components/context/LteNetworkContext';

import {
  BITRATE_MULTIPLIER,
  DATA_PLAN_UNLIMITED_RATES,
} from '../../components/network/DataPlanConst';
import {useParams} from 'react-router-dom';

type CellularNetworkProfile = $Values<
  $NonMaybeType<$PropertyType<network_epc_configs, 'sub_profiles'>>,
>;

/**
 * Calculates the bitrate value in bps based on the default
 * value input in the form, and the previous (default) bitrate.
 *
 * @param {?CellularNetworkProfile} dataPlan
 *    - network epc configs including dataplans
 * @param {?editedValue} editedValue
 *    - the new specified bit rate in Mbps
 * @param {'max_ul_bit_rate' | 'max_dl_bit_rate'}
 *    - specifies whether we are editing upload of download bitrate
 */
function _getBitRateValue(props: {
  dataPlan: ?CellularNetworkProfile,
  editedValue: ?string,
  field: 'max_ul_bit_rate' | 'max_dl_bit_rate',
}): number {
  const {dataPlan, editedValue, field} = props;
  if (editedValue !== null && editedValue === '') {
    return DATA_PLAN_UNLIMITED_RATES[field];
  } else if (editedValue !== null) {
    return Math.max(
      parseFloat(0),
      parseFloat(editedValue) * BITRATE_MULTIPLIER,
    );
  }
  return dataPlan ? dataPlan[field] : DATA_PLAN_UNLIMITED_RATES[field];
}

/**
 * A prop passed to DataPlanEditDialog
 *
 * @property {boolean} open - Whether the dialog is visible
 * @property {() => void} onClose - Callback after closing dialog
 * @property {?string} dataPlanId
 *    - Supplied if editing a data plan.
 *      Not supplied if creating a new data plan.
 */
type DialogProps = {
  open: boolean,
  onClose: () => void,
  dataPlanId: ?string,
};

/**
 * Modal dialog for adding/editing a single data plan.
 * Displays conditionally depending on props.
 *
 * @param {DialogProps} props
 */
export default function DataPlanEditDialog(props: DialogProps) {
  const isAdd = props.dataPlanId ? false : true;
  const onClose = () => {
    props.onClose();
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="sm">
      <DialogTitle
        label={isAdd ? 'Add New Data Plan' : 'Edit Data Plan'}
        onClose={onClose}
      />
      <DataPlanEdit
        onSave={() => {
          onClose();
        }}
        onClose={onClose}
        dataPlanId={props.dataPlanId || ''}
      />
    </Dialog>
  );
}

/**
 * A prop passed to DataPlanEdit
 *
 * @property {() => void} onSave
 *    - Callback after data plan has been saved
 * @property {() => onClose} onClose
 *    - Callback after dialog has been closed
 * @property {string} dataPlanId
 */
type Props = {
  onSave: () => void,
  onClose: () => void,
  dataPlanId: string,
};

/**
 * Modal dialog for adding/editing a single data plan.
 * Always displays.
 *
 * @param {DialogProps} props
 */
export function DataPlanEdit(props: Props) {
  const params = useParams();
  const networkID = nullthrows(params.networkId);
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(LteNetworkContext);
  const epcConfig = ctx.state.cellular.epc;
  const [editedName, setEditedName] = useState(props.dataPlanId || '');
  const dataPlan =
    (props.dataPlanId && epcConfig.sub_profiles?.[props.dataPlanId]) || null;

  // These bit rates are in Mbps
  const [editedMaxDlBitRate, setEditedMaxDlBitRate] = useState(
    dataPlan ? dataPlan.max_dl_bit_rate / BITRATE_MULTIPLIER : '',
  );
  const [editedMaxUlBitRate, setEditedMaxUlBitRate] = useState(
    dataPlan ? dataPlan.max_ul_bit_rate / BITRATE_MULTIPLIER : '',
  );

  const onSave = async () => {
    if (!props.dataPlanId && (epcConfig.sub_profiles || {})[editedName]) {
      setError('Data plan name is already used. Please use a different name');
      return;
    }

    const subProfiles = epcConfig.sub_profiles
      ? {...epcConfig.sub_profiles}
      : {};

    subProfiles[editedName] = {
      max_dl_bit_rate: _getBitRateValue({
        dataPlan,
        field: 'max_dl_bit_rate',
        editedValue: editedMaxDlBitRate.toString(),
      }),
      max_ul_bit_rate: _getBitRateValue({
        dataPlan,
        field: 'max_ul_bit_rate',
        editedValue: editedMaxUlBitRate.toString(),
      }),
    };

    const newConfig = {
      ...epcConfig,
      sub_profiles: subProfiles,
    };

    try {
      const updateNetworkProps: UpdateNetworkContextProps = {
        networkId: networkID,
        epcConfigs: newConfig,
      };
      await ctx.updateNetworks(updateNetworkProps);
      props.onSave();
      enqueueSnackbar('Data plan saved successfully', {
        variant: 'success',
      });
    } catch (error) {
      setError(error.response?.data?.message || error);
    }
  };
  return (
    <>
      <DialogContent data-testid="dataPlanEditDialog">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <div>
            <ListItem dense disableGutters />
            <AltFormField label={'Data Plan ID'}>
              <OutlinedInput
                data-testid="dataPlanID"
                placeholder="Plan 1"
                fullWidth={true}
                value={editedName}
                onChange={({target}) => setEditedName(target.value)}
              />
            </AltFormField>
            <AltFormField label={'Max Bit Rate'}>
              <AltFormFieldSubheading label={'Download'}>
                <OutlinedInput
                  data-testid="dataPlanMaxDlBitRate"
                  placeholder="Unlimited"
                  type="number"
                  min={0}
                  value={editedMaxDlBitRate}
                  onChange={({target}) => setEditedMaxDlBitRate(target.value)}
                  endAdornment={
                    <InputAdornment position="end">
                      <Text variant="subtitle3">Mbps</Text>
                    </InputAdornment>
                  }
                />
              </AltFormFieldSubheading>
              <AltFormFieldSubheading label={'Upload'}>
                <OutlinedInput
                  data-testid="dataPlanMaxUlBitRate"
                  placeholder="Unlimited"
                  type="number"
                  min={0}
                  value={editedMaxUlBitRate}
                  onChange={({target}) => setEditedMaxUlBitRate(target.value)}
                  endAdornment={
                    <InputAdornment position="end">
                      <Text variant="subtitle3">Mbps</Text>
                    </InputAdornment>
                  }
                />
              </AltFormFieldSubheading>
            </AltFormField>
          </div>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
