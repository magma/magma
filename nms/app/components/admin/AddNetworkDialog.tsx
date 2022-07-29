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
import Input from '@mui/material/Input';
import InputLabel from '@mui/material/InputLabel';
import ListItemText from '@mui/material/ListItemText';
import MenuItem from '@mui/material/MenuItem';
import React, {useContext} from 'react';
import Select from '@mui/material/Select';
import TextField from '@mui/material/TextField';
import axios from 'axios';

import AppContext from '../../context/AppContext';
import nullthrows from '../../../shared/util/nullthrows';
import {
  AllNetworkTypes,
  CWF,
  FEG,
  FEG_LTE,
  NetworkId,
  XWFM,
} from '../../../shared/types/network';
import {getErrorMessage} from '../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {triggerAlertSync} from '../../util/SyncAlerts';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

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
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
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
          if (payload.data.networkType === XWFM) {
            void triggerAlertSync(networkID, enqueueSnackbar);
          }
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
        <TextField
          name="networkId"
          label="Network ID"
          className={classes.input}
          value={networkID}
          onChange={({target}) => setNetworkId(target.value)}
        />
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={name}
          onChange={({target}) => setName(target.value)}
        />
        <TextField
          name="description"
          label="Description"
          className={classes.input}
          value={description}
          onChange={({target}) => setDescription(target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Network Type</InputLabel>
          <Select
            value={networkType}
            onChange={({target}) => setNetworkType(target.value)}
            input={<Input id="types" />}>
            {AllNetworkTypes.map(type => (
              <MenuItem key={type} value={type}>
                <ListItemText primary={type} />
              </MenuItem>
            ))}
          </Select>
        </FormControl>
        {(networkType === CWF || networkType === FEG_LTE) && (
          <TextField
            name="fegNetworkID"
            label="Federation Network ID"
            className={classes.input}
            value={fegNetworkID}
            onChange={({target}) => setFegNetworkID(target.value)}
          />
        )}
        {networkType === FEG && (
          <TextField
            placeholder="network1,network2"
            label="Served Network IDs"
            className={classes.input}
            value={servedNetworkIDs}
            onChange={({target}) => setServedNetworkIDs(target.value)}
          />
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
