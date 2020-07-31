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
import type {network, network_dns_config} from '@fbcnms/magma-api';

import Button from '@material-ui/core/Button';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import KPIGrid from '../../components/KPIGrid';
import List from '@material-ui/core/List';
import ListItemText from '@material-ui/core/ListItemText';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Select from '@material-ui/core/Select';
import axios from 'axios';

import {AllNetworkTypes} from '@fbcnms/types/network';
import {AltFormField} from '../../components/FormField';
import {CWF, FEG, LTE} from '@fbcnms/types/network';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

type Props = {
  networkInfo: network,
};

export default function NetworkInfo(props: Props) {
  const kpiData: KPIRows[] = [
    [
      {
        category: 'ID',
        value: props.networkInfo.id,
      },
    ],
    [
      {
        category: 'Name',
        value: props.networkInfo.name,
      },
    ],
    [
      {
        category: 'Network Type',
        value:
          typeof props.networkInfo.type != 'undefined'
            ? props.networkInfo.type
            : '-',
      },
    ],
    [
      {
        category: 'Description',
        value: props.networkInfo.description,
      },
    ],
  ];
  return (
    <Paper elevation={0} data-testid="info">
      <KPIGrid data={kpiData} />
    </Paper>
  );
}

type EditProps = {
  saveButtonTitle: string,
  networkInfo: ?network,
  onClose: () => void,
  onSave: network => void,
};

const DEFAULT_DNS_CONFIG: network_dns_config = {
  enable_caching: false,
  local_ttl: 0,
  records: [],
};

export function NetworkInfoEdit(props: EditProps) {
  const [error, setError] = useState('');
  const [fegNetworkID, setFegNetworkID] = useState('');
  const [servedNetworkIDs, setServedNetworkIDs] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const [networkType, setNetworkType] = useState(
    props.networkInfo?.type || LTE,
  );
  const [networkInfo, setNetworkInfo] = useState<network>(
    props.networkInfo || {
      name: '',
      id: '',
      description: '',
      dns: DEFAULT_DNS_CONFIG,
    },
  );

  const onSave = async () => {
    const payload = {
      networkID: networkInfo.id,
      data: {
        name: networkInfo.name,
        description: networkInfo.description,
        networkType,
        fegNetworkID,
        servedNetworkIDs,
      },
    };
    if (props.networkInfo) {
      // edit
      try {
        await MagmaV1API.putNetworksByNetworkId({
          networkId: networkInfo.id,
          network: {
            ...networkInfo,
            type: networkType,
          },
        });
        enqueueSnackbar('Network configs saved successfully', {
          variant: 'success',
        });
        props.onSave(networkInfo);
      } catch (e) {
        setError(e.data?.message ?? e.message);
      }
    } else {
      try {
        const response = await axios.post('/nms/network/create', payload);
        if (response.data.success) {
          enqueueSnackbar(`Network $networkInfo.name} successfully created`, {
            variant: 'success',
          });
          props.onSave(networkInfo);
        } else {
          setError(response.data.message);
        }
      } catch (e) {
        setError(e.data?.message ?? e.message);
      }
    }
  };

  return (
    <>
      <DialogContent data-testid="networkInfoEdit">
        {error !== '' && (
          <AltFormField label={''}>
            <FormLabel error>{error}</FormLabel>
          </AltFormField>
        )}
        <List>
          <AltFormField label={'Network ID'}>
            <OutlinedInput
              data-testid="networkID"
              fullWidth={true}
              value={networkInfo.id}
              onChange={({target}) =>
                setNetworkInfo({...networkInfo, id: target.value})
              }
              disabled={props.networkInfo ? true : false}
            />
          </AltFormField>
          <AltFormField label={'Network Name'}>
            <OutlinedInput
              data-testid="networkName"
              fullWidth={true}
              value={networkInfo.name}
              onChange={({target}) =>
                setNetworkInfo({...networkInfo, name: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'Network Type'}>
            <Select
              variant={'outlined'}
              fullWidth={true}
              value={networkType}
              onChange={({target}) => {
                setNetworkType(target.value);
              }}
              data-testid="networkType"
              input={<OutlinedInput fullWidth={true} id="networkType" />}>
              {AllNetworkTypes.map(type => (
                <MenuItem key={type} value={type}>
                  <ListItemText primary={type} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>

          {networkType === CWF && (
            <AltFormField label={'Federation Network ID'}>
              <OutlinedInput
                fullWidth={true}
                value={fegNetworkID}
                onChange={({target}) => setFegNetworkID(target.value)}
              />
            </AltFormField>
          )}
          {networkType === FEG && (
            <AltFormField label={'Served Network IDs'}>
              <OutlinedInput
                placeholder="network1,network2"
                fullWidth={true}
                value={servedNetworkIDs}
                onChange={({target}) => setServedNetworkIDs(target.value)}
              />
            </AltFormField>
          )}
          <AltFormField label={'Add Description'}>
            <OutlinedInput
              data-testid="networkDescription"
              fullWidth={true}
              multiline
              rows={4}
              value={networkInfo.description}
              onChange={({target}) =>
                setNetworkInfo({...networkInfo, description: target.value})
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
