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
import type {DataRows} from '../../components/DataGrid';
import type {network_epc_configs} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import DataGrid from '../../components/DataGrid';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import IconButton from '@material-ui/core/IconButton';
import InputAdornment from '@material-ui/core/InputAdornment';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import Switch from '@material-ui/core/Switch';
import Visibility from '@material-ui/icons/Visibility';
import VisibilityOff from '@material-ui/icons/VisibilityOff';

import {AltFormField} from '../../components/FormField';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

type Props = {
  epcConfigs: network_epc_configs,
};

export default function NetworkEpc(props: Props) {
  const kpiData: DataRows[] = [
    [
      {
        category: 'Policy Enforcement Enabled',
        value: props.epcConfigs?.hss_relay_enabled ? 'Enabled' : 'Disabled',
      },
    ],
    [
      {
        category: 'LTE Auth AMF',
        value: props.epcConfigs?.lte_auth_amf,
        obscure: true,
      },
    ],
    [
      {
        category: 'MCC',
        value: props.epcConfigs?.mcc,
      },
    ],
    [
      {
        category: 'MNC',
        value: props.epcConfigs?.mnc,
      },
    ],
    [
      {
        category: 'TAC',
        value: props.epcConfigs?.tac,
      },
    ],
  ];

  return <DataGrid data={kpiData} testID="epc" />;
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
  const ctx = useContext(LteNetworkContext);
  const IPallocationMode = ['NAT', 'DHCP_BROADCAST'];
  const [epcConfigs, setEpcConfigs] = useState<network_epc_configs>(
    props.epcConfigs == null || Object.keys(props.epcConfigs).length === 0
      ? {
          cloud_subscriberdb_enabled: false,
          default_rule_id: 'default_rule_1',
          lte_auth_amf: 'gAA=',
          lte_auth_op: 'EREREREREREREREREREREQ==',
          mcc: '001',
          mnc: '01',
          network_services: ['policy_enforcement'],
          hss_relay_enabled: false,
          gx_gy_relay_enabled: false,
          sub_profiles: {},
          tac: 1,
        }
      : props.epcConfigs,
  );
  const [epcMobility, setEpcMobility] = useState(
    props.epcConfigs?.mobility || {
      ip_allocation_mode: 'NAT',
      enable_static_ip_assignments: false,
      enable_multi_apn_ip_allocation: false,
    },
  );
  const handleMobilityChange = (key: string, val) =>
    setEpcMobility({...epcMobility, [key]: val});
  const onSave = async () => {
    try {
      await ctx.updateNetworks({
        networkId: props.networkId,
        epcConfigs: {...epcConfigs, mobility: epcMobility},
      });
      props.onSave({...epcConfigs, mobility: epcMobility});
      enqueueSnackbar('EPC configs saved successfully', {variant: 'success'});
    } catch (e) {
      setError(e.response?.data?.message ?? e?.message);
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
          <AltFormField label={'IP Allocation Mode'}>
            <Select
              variant={'outlined'}
              displayEmpty={true}
              value={epcMobility.ip_allocation_mode}
              onChange={({target}) =>
                handleMobilityChange('ip_allocation_mode', target.value)
              }
              data-testid="IpAllocationMode"
              input={<OutlinedInput />}>
              {IPallocationMode.map(mode => (
                <MenuItem key={mode} value={mode}>
                  <ListItemText primary={mode} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>
          <AltFormField label={'Static IP Assignments'}>
            <Switch
              onChange={() => {
                handleMobilityChange(
                  'enable_static_ip_assignments',
                  !epcMobility.enable_static_ip_assignments,
                );
              }}
              checked={epcMobility.enable_static_ip_assignments}
            />
          </AltFormField>
          <AltFormField label={'Multi APN IP Allocation'}>
            <Switch
              onChange={() => {
                handleMobilityChange(
                  'enable_multi_apn_ip_allocation',
                  !epcMobility.enable_multi_apn_ip_allocation,
                );
              }}
              checked={epcMobility.enable_multi_apn_ip_allocation}
              disabled={!(epcMobility.ip_allocation_mode === 'DHCP_BROADCAST')}
            />
          </AltFormField>
          <AltFormField label={'Policy Enforcement Enabled'}>
            <Select
              variant={'outlined'}
              value={epcConfigs.hss_relay_enabled ? 1 : 0}
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
              placeholder="Enter Auth AMF"
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
              placeholder="Enter MCC"
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
              placeholder="Enter MNC"
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
              placeholder="Enter TAC"
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
        <Button data-testid="epcCancelButton" onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          data-testid="epcSaveButton"
          onClick={onSave}
          variant="contained"
          color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}
