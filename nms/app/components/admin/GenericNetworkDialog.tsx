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
import Button from '@mui/material/Button';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '@mui/material/DialogTitle';
import FormControl from '@mui/material/FormControl';
import ListItemText from '@mui/material/ListItemText';
import MenuItem from '@mui/material/MenuItem';
import Select from '@mui/material/Select';
import {useState} from 'react';

import {AltFormField} from '../FormField';
import {OutlinedInput} from '@mui/material';

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
        <AltFormField label="Name">
          <OutlinedInput
            name="name"
            fullWidth
            value={networkConfig.name}
            onChange={({target}) => updateNetwork('name', target.value)}
          />
        </AltFormField>
        <AltFormField label="Description">
          <OutlinedInput
            name="description"
            fullWidth
            value={networkConfig.description}
            onChange={({target}) => updateNetwork('description', target.value)}
          />
        </AltFormField>
        <AltFormField label="Network Type">
          <FormControl fullWidth>
            <Select
              value={networkConfig.type}
              onChange={({target}) => updateNetwork('type', target.value)}
              input={<OutlinedInput id="types" />}
              disabled={true}>
              <MenuItem key={networkConfig.type} value={networkConfig.type}>
                <ListItemText primary={networkConfig.type} />
              </MenuItem>
            </Select>
          </FormControl>
        </AltFormField>
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
