/*
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
import type {KPIRows} from '../../components/KPIGrid';
import type {network_epc_configs} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {AltFormField} from '../../components/FormField';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

type Props = {
  epcConfigs: network_epc_configs,
};

export default function NetworkEpc(props: Props) {
  const kpiData: KPIRows[] = [
    [
      {
        category: 'Policy Enforcement Enabled',
        value: props.epcConfigs.relay_enabled ? 'Enabled' : 'Disabled',
      },
    ],
    [
      {
        category: 'LTE Auth AMF',
        value: props.epcConfigs.lte_auth_amf,
        obscure: true,
      },
    ],
    [
      {
        category: 'MCC',
        value: props.epcConfigs.mcc,
      },
    ],
    [
      {
        category: 'MNC',
        value: props.epcConfigs.mnc,
      },
    ],
    [
      {
        category: 'TAC',
        value: props.epcConfigs.tac,
      },
    ],
  ];

  return (
    <Paper elevation={0} data-testid="epc">
      <KPIGrid data={kpiData} />
    </Paper>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkId: string,
  epcConfigs: ?network_epc_configs,
  onClose: () => void,
  onSave: network_epc_configs => void,
};

export function NetworkEpcEdit(props: EditProps) {
  const [showPassword, setShowPassword] = React.useState(false);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [epcConfigs, setEpcConfigs] = useState<network_epc_configs>(
    props.epcConfigs || {
      cloud_subscriberdb_enabled: false,
      default_rule_id: 'default_rule_1',
      lte_auth_amf: 'gAA=',
      lte_auth_op: 'EREREREREREREREREREREQ==',
      mcc: '001',
      mnc: '01',
      network_services: ['policy_enforcement'],
      relay_enabled: false,
      sub_profiles: {},
      tac: 1,
    },
  );

  const onSave = async () => {
    try {
      MagmaV1API.putLteByNetworkIdCellularEpc({
        networkId: props.networkId,
        config: epcConfigs,
      });
      enqueueSnackbar('EPC configs saved successfully', {variant: 'success'});
      props.onSave(epcConfigs);
    } catch (e) {
      setError(e.data?.message ?? e.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="networkEpcEdit">
        {error !== '' && (
          <AltFormField label={''}>
            <FormLabel error>{error}</FormLabel>
          </AltFormField>
        )}
        <List>
          <AltFormField label={'Policy Enforcement Enabled'}>
            <Select
              variant={'outlined'}
              value={epcConfigs.relay_enabled ? 1 : 0}
              onChange={({target}) => {
                setEpcConfigs({
                  ...epcConfigs,
                  relay_enabled: target.value === 1,
                });
              }}
              input={<OutlinedInput id="relayEnabled" />}>
              <MenuItem value={0}>
                <ListItemText primary={'Disabled'} />
              </MenuItem>
              <MenuItem value={1}>
                <ListItemText primary={'Enabled'} />
              </MenuItem>
            </Select>
          </AltFormField>
          <AltFormField label={'LTE Auth AMF'}>
            <OutlinedInput
              data-testid="password"
              type={showPassword ? 'text' : 'password'}
              fullWidth={true}
              value={epcConfigs.lte_auth_amf}
              onChange={({target}) => {
                setEpcConfigs({...epcConfigs, lte_auth_amf: target.value});
              }}
              endAdornment={
                <InputAdornment position="end">
                  <IconButton
                    aria-label="toggle password visibility"
                    onClick={() => setShowPassword(!showPassword)}
                    onMouseDown={event => event.preventDefault()}>
                    {showPassword ? <Visibility /> : <VisibilityOff />}
                  </IconButton>
                </InputAdornment>
              }
            />
          </AltFormField>
          <AltFormField label={'MCC'}>
            <OutlinedInput
              data-testid="mcc"
              fullWidth={true}
              value={epcConfigs.mcc}
              onChange={({target}) =>
                setEpcConfigs({...epcConfigs, mcc: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'MNC'}>
            <OutlinedInput
              data-testid="mnc"
              fullWidth={true}
              value={epcConfigs.mnc}
              onChange={({target}) =>
                setEpcConfigs({...epcConfigs, mnc: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'TAC'}>
            <OutlinedInput
              data-testid="tac"
              type="number"
              fullWidth={true}
              value={epcConfigs.tac}
              onChange={({target}) =>
                setEpcConfigs({...epcConfigs, tac: parseInt(target.value)})
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose} skin="regular">
          Cancel
        </Button>
        <Button onClick={onSave} variant="contained" color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}
