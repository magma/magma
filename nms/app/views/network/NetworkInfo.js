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
import type {
  feg_lte_network,
  lte_network,
} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import DataGrid from '../../components/DataGrid';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import LteNetworkContext from '../../components/context/LteNetworkContext';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../components/context/NetworkContext';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import axios from 'axios';

import {AltFormField} from '../../components/FormField';
// $FlowFixMe migrated to typescript
import {FEG_LTE, LTE} from '../../../shared/types/network';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

type Props = {
  lteNetwork: $Shape<lte_network & feg_lte_network>,
};

export default function NetworkInfo(props: Props) {
  const networkCtx = useContext(NetworkContext);

  const kpiData: DataRows[] = [
    [
      {
        category: 'ID',
        value: props.lteNetwork.id,
      },
    ],
    [
      {
        category: 'Name',
        value: props.lteNetwork.name,
      },
    ],
    [
      {
        category: 'Description',
        value: props.lteNetwork.description || '-',
      },
    ],
  ];
  if (networkCtx.networkType === FEG_LTE) {
    kpiData.push(
      [
        {
          category: 'Federation',
          value: props.lteNetwork?.federation?.feg_network_id || '-',
        },
      ],
      [
        {
          category: 'Federated Mapping Mode',
          value:
            props.lteNetwork?.federation?.federated_modes_mapping?.enabled ===
            true
              ? 'On'
              : 'Off',
        },
      ],
    );
  }
  return <DataGrid data={kpiData} testID="info" />;
}

type EditProps = {
  saveButtonTitle: string,
  lteNetwork: $Shape<lte_network & feg_lte_network>,
  onClose: () => void,
  onSave: ($Shape<lte_network & feg_lte_network>) => void,
};

export function NetworkInfoEdit(props: EditProps) {
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();
  const ctx = useContext(LteNetworkContext);
  const [lteNetwork, setLteNetwork] = useState<
    $Shape<lte_network & feg_lte_network>,
  >(props.lteNetwork);

  const onSave = async () => {
    if (props.lteNetwork?.id) {
      // edit
      try {
        await ctx.updateNetworks({networkId: lteNetwork.id, lteNetwork});
        enqueueSnackbar('Network configs saved successfully', {
          variant: 'success',
        });
        props.onSave(lteNetwork);
      } catch (e) {
        setError(e.response?.data?.message ?? e?.message);
      }
    } else {
      // network creation a special case. We have to update the organization
      // information in db, so we hijack the request and update the org info
      // with the networkID
      try {
        const payload = {
          networkID: lteNetwork.id,
          data: {
            name: lteNetwork.name,
            description: lteNetwork.description,
            networkType: LTE,
          },
        };
        const response = await axios.post('/nms/network/create', payload);
        if (response.data.success) {
          enqueueSnackbar(`Network ${lteNetwork.name} successfully created`, {
            variant: 'success',
          });
          props.onSave(lteNetwork);
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
              placeholder="Enter ID"
              fullWidth={true}
              value={lteNetwork.id}
              onChange={({target}) =>
                setLteNetwork({...lteNetwork, id: target.value})
              }
              disabled={props.lteNetwork?.id ? true : false}
            />
          </AltFormField>
          <AltFormField label={'Network Name'}>
            <OutlinedInput
              data-testid="networkName"
              placeholder="Enter Name"
              fullWidth={true}
              value={lteNetwork.name}
              onChange={({target}) =>
                setLteNetwork({...lteNetwork, name: target.value})
              }
            />
          </AltFormField>
          <AltFormField label={'Add Description'}>
            <OutlinedInput
              data-testid="networkDescription"
              placeholder="Enter Description"
              fullWidth={true}
              multiline
              rows={4}
              value={lteNetwork.description}
              onChange={({target}) =>
                setLteNetwork({...lteNetwork, description: target.value})
              }
            />
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button data-testid="cancelButton" onClick={props.onClose}>
          Cancel
        </Button>
        <Button
          data-testid="saveButton"
          onClick={onSave}
          variant="contained"
          color="primary">
          {props.saveButtonTitle}
        </Button>
      </DialogActions>
    </>
  );
}
