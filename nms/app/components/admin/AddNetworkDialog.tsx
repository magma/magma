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

import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControl from '@mui/material/FormControl';
import FormLabel from '@mui/material/FormLabel';
import ListItemText from '@mui/material/ListItemText';
import MenuItem from '@mui/material/MenuItem';
import React, {useContext, useState} from 'react';
import Select from '@mui/material/Select';
import axios from 'axios';

import AppContext from '../../context/AppContext';
import nullthrows from '../../../shared/util/nullthrows';
import {
  AllNetworkTypes,
  CWF,
  FEG,
  FEG_LTE,
  NetworkId,
} from '../../../shared/types/network';
import {AltFormField} from '../FormField';
import {OutlinedInput} from '@mui/material';
import {getErrorMessage} from '../../util/ErrorUtils';

type Props = {
  onClose: () => void;
  onSave: (value: NetworkId) => void;
};

type CreateResponse =
  | {
      success: true;
      apiResponse: 'Success';
    }
  | {
      success: false;
      apiResponse: object;
      message: string;
    };

export default function NetworkDialog(props: Props) {
  const [name, setName] = useState('');
  const [description, setDescription] = useState('');
  const [networkID, setNetworkId] = useState('');
  const [networkType, setNetworkType] = useState('');
  const [fegNetworkID, setFegNetworkID] = useState('');
  const [servedNetworkIDs, setServedNetworkIDs] = useState('');
  const [error, setError] = useState<string | null>(null);
  const appContext = useContext(AppContext);

  const onSave = () => {
    const payload = {
      networkID,
      data: {
        name,
        description,
        networkType,
        fegNetworkID,
        servedNetworkIDs,
      },
    };
    axios
      .post<CreateResponse>('/nms/network/create', payload)
      .then(response => {
        if (response.data.success) {
          props.onSave(nullthrows(networkID));
          appContext.addNetworkId(networkID);
        } else {
          setError(response.data.message);
        }
      })
      .catch(error => setError(getErrorMessage(error)));
  };

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Add Network</DialogTitle>
      <DialogContent>
        {error && <FormLabel error>{error}</FormLabel>}
        <AltFormField label="Network ID">
          <OutlinedInput
            fullWidth={true}
            name="networkId"
            value={networkID}
            onChange={({target}) => setNetworkId(target.value)}
          />
        </AltFormField>
        <AltFormField label="Name">
          <OutlinedInput
            fullWidth={true}
            name="name"
            value={name}
            onChange={({target}) => setName(target.value)}
          />
        </AltFormField>
        <AltFormField label="Description">
          <OutlinedInput
            fullWidth={true}
            name="description"
            value={description}
            onChange={({target}) => setDescription(target.value)}
          />
        </AltFormField>
        <AltFormField label="Network Type">
          <FormControl fullWidth>
            <Select
              value={networkType}
              onChange={({target}) => setNetworkType(target.value)}
              input={<OutlinedInput id="types" />}>
              {AllNetworkTypes.map(type => (
                <MenuItem key={type} value={type}>
                  <ListItemText primary={type} />
                </MenuItem>
              ))}
            </Select>
          </FormControl>
        </AltFormField>
        {(networkType === CWF || networkType === FEG_LTE) && (
          <AltFormField label="Federation Network ID">
            <OutlinedInput
              fullWidth={true}
              name="fegNetworkID"
              value={fegNetworkID}
              onChange={({target}) => setFegNetworkID(target.value)}
            />
          </AltFormField>
        )}
        {networkType === FEG && (
          <AltFormField label="Served Network IDs">
            <OutlinedInput
              placeholder="network1,network2"
              value={servedNetworkIDs}
              onChange={({target}) => setServedNetworkIDs(target.value)}
            />
          </AltFormField>
        )}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button onClick={onSave} variant="contained" color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
