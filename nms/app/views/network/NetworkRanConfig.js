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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {DataRows} from '../../components/DataGrid';
import type {network_ran_configs} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FddConfig from './NetworkRanFddConfig';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
// $FlowFixMe migrated to typescript
import LteNetworkContext from '../../components/context/LteNetworkContext';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';
import TddConfig from './NetworkRanTddConfig';

// $FlowFixMe migrated to typescript
import {AltFormField, FormDivider} from '../../components/FormField';
import {useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

type Props = {
  lteRanConfigs: network_ran_configs,
};

export default function NetworkRan(props: Props) {
  const tdd: DataRows[] = [
    [
      {
        category: 'EARFCNDL',
        value: props.lteRanConfigs?.tdd_config?.earfcndl || '-',
      },
    ],
    [
      {
        category: 'Special Subframe Pattern',
        value: props.lteRanConfigs?.tdd_config?.special_subframe_pattern || '-',
      },
    ],
    [
      {
        category: 'Subframe Assignment',
        value: props.lteRanConfigs?.tdd_config?.subframe_assignment || '-',
      },
    ],
  ];

  const fdd: DataRows[] = [
    [
      {
        category: 'EARFCNDL',
        value: props.lteRanConfigs?.fdd_config?.earfcndl || '-',
      },
    ],
    [
      {
        category: 'EARFCNUL',
        value: props.lteRanConfigs?.fdd_config?.earfcnul || '-',
      },
    ],
  ];

  const ran: DataRows[] = [
    [
      {
        category: 'Bandwidth',
        value: props.lteRanConfigs?.bandwidth_mhz || '-',
      },
    ],
    [
      {
        category: 'RAN Config',
        value: props.lteRanConfigs?.fdd_config
          ? 'FDD'
          : props.lteRanConfigs?.tdd_config
          ? 'TDD'
          : '-',
        collapse: props.lteRanConfigs?.fdd_config ? (
          <DataGrid data={fdd} />
        ) : props.lteRanConfigs?.tdd_config ? (
          <DataGrid data={tdd} />
        ) : (
          false
        ),
      },
    ],
  ];

  return <DataGrid data={ran} testID="ran" />;
}

type EditProps = {
  saveButtonTitle: string,
  networkId: string,
  lteRanConfigs: ?network_ran_configs,
  onClose: () => void,
  onSave: network_ran_configs => void,
};
type BandType = 'tdd' | 'fdd';
const ValidBandwidths = [3, 5, 10, 15, 20];

export function NetworkRanEdit(props: EditProps) {
  const ctx = useContext(LteNetworkContext);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');
  const [bandType, setBandType] = useState<BandType>('tdd');
  const defaultTddConfig = {
    earfcndl: 44590,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
  };
  const defaulFddConfig = {
    earfcndl: 0,
    earfcnul: 0,
  };

  const [lteRanConfigs, setLteRanConfigs] = useState(
    props.lteRanConfigs == null || Object.keys(props.lteRanConfigs).length === 0
      ? {
          bandwidth_mhz: 20,
          fdd_config: undefined,
          tdd_config: defaultTddConfig,
        }
      : props.lteRanConfigs,
  );

  const onSave = async () => {
    const config: network_ran_configs = {
      ...lteRanConfigs,
    };
    if (bandType === 'tdd') {
      config.fdd_config = undefined;
    } else {
      config.tdd_config = undefined;
    }

    try {
      await ctx.updateNetworks({
        networkId: props.networkId,
        lteRanConfigs: config,
      });
      enqueueSnackbar('RAN configs saved successfully', {variant: 'success'});
      props.onSave(lteRanConfigs);
    } catch (e) {
      setError(e.response?.data?.message ?? e?.message);
    }
  };

  return (
    <>
      <DialogContent data-testid="networkRanEdit">
        {error !== '' && (
          <AltFormField label={''}>
            <FormLabel error>{error}</FormLabel>
          </AltFormField>
        )}
        <List>
          <AltFormField label={'Bandwidth'}>
            <Select
              variant={'outlined'}
              fullWidth={true}
              value={lteRanConfigs.bandwidth_mhz}
              onChange={({target}) => {
                if (
                  target.value === 3 ||
                  target.value === 5 ||
                  target.value === 10 ||
                  target.value === 15 ||
                  target.value === 20
                ) {
                  setLteRanConfigs({
                    ...lteRanConfigs,
                    bandwidth_mhz: target.value,
                  });
                }
              }}
              input={<OutlinedInput fullWidth={true} id="bandwidth" />}>
              {ValidBandwidths.map((k: number, idx: number) => (
                <MenuItem key={idx} value={k}>
                  <ListItemText primary={k} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>
          <AltFormField label={'Band Type'}>
            <Select
              variant={'outlined'}
              fullWidth={true}
              value={bandType}
              onChange={({target}) => {
                if (target.value === 'fdd') {
                  setLteRanConfigs({
                    fdd_config: defaulFddConfig,
                    ...lteRanConfigs,
                  });
                  setBandType('fdd');
                } else {
                  setLteRanConfigs({
                    tdd_config: defaultTddConfig,
                    ...lteRanConfigs,
                  });
                  setBandType(target.value === 'tdd' ? 'tdd' : 'fdd');
                }
              }}
              input={<OutlinedInput fullWidth={true} id="bandType" />}>
              <MenuItem value={'tdd'}>
                <ListItemText primary={'TDD'} />
              </MenuItem>
              <MenuItem value={'fdd'}>
                <ListItemText primary={'FDD'} />
              </MenuItem>
            </Select>
          </AltFormField>
          <FormDivider />
          {bandType === 'tdd' && (
            <TddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
          {bandType === 'fdd' && (
            <FddConfig
              lteRanConfigs={lteRanConfigs}
              setLteRanConfigs={setLteRanConfigs}
            />
          )}
        </List>
      </DialogContent>
      <DialogActions>
        <Button data-testid="ranCancelButton" onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          data-testid="ranSaveButton"
          onClick={onSave}
          variant="contained"
          color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}
