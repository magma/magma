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
 * @flow strict-local
 * @format
 */

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import FormLabel from '@material-ui/core/FormLabel';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';
import axios from 'axios';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import {
  AllNetworkTypes,
  CWF,
  FEG,
  FEG_LTE,
  XWFM,
  // $FlowFixMe migrated to typescript
} from '../../../shared/types/network';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {triggerAlertSync} from '../../state/SyncAlerts';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

type Props = {
  onClose: () => void,
  onSave: string => void,
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
  const [error, setError] = useState(null);

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
      .post('/nms/network/create', payload)
      .then(response => {
        if (response.data.success) {
          props.onSave(nullthrows(networkID));
          if (payload.data.networkType === XWFM) {
            triggerAlertSync(networkID, enqueueSnackbar);
          }
        } else {
          setError(response.data.message);
        }
      })
      .catch(error => setError(error.response?.data?.error || error));
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
