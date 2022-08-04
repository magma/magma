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

import * as React from 'react';
import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '@material-ui/core/DialogTitle';
import FormControl from '@material-ui/core/FormControl';
import Input from '@material-ui/core/Input';
import InputLabel from '@material-ui/core/InputLabel';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import Select from '@material-ui/core/Select';
import TextField from '@material-ui/core/TextField';

import {makeStyles} from '@material-ui/styles';
import {useState} from 'react';

const useStyles = makeStyles(() => ({
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '100%',
  },
}));

export type GenericConfig = {
  id: string;
  name: string;
  description: string;
  type?: string;
};

type Props = {
  onClose: () => void;
  onSave: (config: GenericConfig) => void;
  networkConfig: GenericConfig;
  children?: React.ReactNode;
};

export default function GenericNetworkDialog(props: Props) {
  const classes = useStyles();
  const [networkConfig, setNetworkConfig] = useState(props.networkConfig);

  const updateNetwork = (
    field: 'name' | 'description' | 'type',
    value: string,
  ) =>
    setNetworkConfig({
      ...networkConfig,
      [field]: value,
    });

  return (
    <Dialog open={true} onClose={props.onClose}>
      <DialogTitle>Editing "{networkConfig.id}"</DialogTitle>
      <DialogContent>
        <TextField
          name="name"
          label="Name"
          className={classes.input}
          value={networkConfig.name}
          onChange={({target}) => updateNetwork('name', target.value)}
        />
        <TextField
          name="description"
          label="Description"
          className={classes.input}
          value={networkConfig.description}
          onChange={({target}) => updateNetwork('description', target.value)}
        />
        <FormControl className={classes.input}>
          <InputLabel htmlFor="types">Network Type</InputLabel>
          <Select
            value={networkConfig.type}
            onChange={({target}) =>
              updateNetwork('type', target.value as string)
            }
            input={<Input id="types" />}
            disabled={true}>
            <MenuItem key={networkConfig.type} value={networkConfig.type}>
              <ListItemText primary={networkConfig.type} />
            </MenuItem>
          </Select>
        </FormControl>
        {props.children}
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button
          onClick={() => props.onSave(networkConfig)}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
