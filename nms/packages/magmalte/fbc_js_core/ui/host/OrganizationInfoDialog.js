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
import type {DialogProps} from './OrganizationDialog';

import ArrowDropDown from '@material-ui/icons/ArrowDropDown';
import Button from '../../../fbc_js_core/ui/components/design-system/Button';
import Checkbox from '@material-ui/core/Checkbox';
import Collapse from '@material-ui/core/Collapse';
import DialogContent from '@material-ui/core/DialogContent';
import FormControlLabel from '@material-ui/core/FormControlLabel';
import FormHelperText from '@material-ui/core/FormHelperText';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React from 'react';
import Select from '@material-ui/core/Select';

import {AltFormField} from '../../../fbc_js_core/ui/components/design-system/FormField/FormField';
import {useState} from 'react';

const ENABLE_ALL_NETWORKS_HELPER =
  'By checking this, the organization will have access to all existing and future networks.';

/**
 * Create Organization Tab
 * This component displays a form used to create an organization
 */
export default function OrganizationInfoDialog(props: DialogProps) {
  const {
    organization,
    allNetworks,
    shouldEnableAllNetworks,
    setShouldEnableAllNetworks,
  } = props;
  const [open, setOpen] = useState(false);

  return (
    <DialogContent>
      <List>
        {props.error && (
          <AltFormField label={''}>
            <FormLabel error>{props.error}</FormLabel>
          </AltFormField>
        )}
        <AltFormField disableGutters label={'Organization Name'}>
          <OutlinedInput
            data-testid="name"
            placeholder="Organization Name"
            fullWidth={true}
            value={organization.name}
            onChange={({target}) => {
              props.onOrganizationChange({...organization, name: target.value});
            }}
          />
        </AltFormField>
        <ListItem disableGutters>
          <Button variant="text" onClick={() => setOpen(!open)}>
            Advanced Settings
          </Button>
          <ArrowDropDown />
        </ListItem>
        <Collapse in={open}>
          <AltFormField
            disableGutters
            label={'Accessible Networks'}
            subLabel={'The networks that the organization have access to'}>
            <Select
              fullWidth={true}
              variant={'outlined'}
              multiple={true}
              renderValue={selected => selected.join(', ')}
              value={organization.networkIDs || []}
              onChange={({target}) => {
                props.onOrganizationChange({
                  ...organization,
                  networkIDs: [...target.value],
                });
              }}
              input={<OutlinedInput id="direction" />}>
              {allNetworks.map(network => (
                <MenuItem key={network} value={network}>
                  <ListItemText primary={network} />
                </MenuItem>
              ))}
            </Select>
          </AltFormField>
          <FormControlLabel
            disableGutters
            label={'Give this organization access to all networks'}
            control={
              <Checkbox
                checked={shouldEnableAllNetworks}
                onChange={() =>
                  setShouldEnableAllNetworks(!shouldEnableAllNetworks)
                }
              />
            }
          />
          <FormHelperText>{ENABLE_ALL_NETWORKS_HELPER}</FormHelperText>
        </Collapse>
      </List>
    </DialogContent>
  );
}
