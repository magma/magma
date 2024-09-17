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

import type {BaseNameRecord} from '../../../generated';

import Button from '@mui/material/Button';
import Checkbox from '@mui/material/Checkbox';
import Dialog from '@mui/material/Dialog';
import DialogActions from '@mui/material/DialogActions';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@mui/material/FormLabel';
import List from '@mui/material/List';
import ListItemText from '@mui/material/ListItemText';
import MenuItem from '@mui/material/MenuItem';
import OutlinedInput from '@mui/material/OutlinedInput';
import PolicyContext from '../../context/PolicyContext';
import React from 'react';
import Select from '@mui/material/Select';
import SubscriberContext from '../../context/SubscriberContext';
import {AltFormField} from '../../components/FormField';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useContext, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

/**
 * A prop passed to DataPlanEditDialog
 *
 * @property {boolean} open - Whether the dialog is visible
 * @property {() => void} onClose - Callback after closing dialog
 * @property {?string} baseNameId
 *    - Supplied if editing a base name.
 *      Not supplied if creating a new base name.
 */
type DialogProps = {
  open: boolean;
  onClose: () => void;
  baseNameId?: string;
};

/**
 * Modal dialog for adding/editing a single base name.
 * Displays conditionally depending on props.
 *
 * @param {DialogProps} props
 */
export default function BaseNameEditDialog(props: DialogProps) {
  const isAdd: boolean = props.baseNameId ? false : true;
  const onClose = () => {
    props.onClose();
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      fullWidth={true}
      maxWidth="sm">
      <DialogTitle
        label={isAdd ? 'Add New Base Name' : 'Edit Base Name'}
        onClose={onClose}
      />
      <BaseNameEdit
        onSave={() => {
          onClose();
        }}
        onClose={onClose}
        baseNameId={props.baseNameId || ''}
      />
    </Dialog>
  );
}

/**
 * A prop passed to BaseNameEdit
 *
 * @property {() => void} onSave
 *    - Callback after data plan has been saved
 * @property {() => onClose} onClose
 *    - Callback after dialog has been closed
 * @property {string} baseNameId
 */
type Props = {
  onSave: () => void;
  onClose: () => void;
  baseNameId: string;
};

/**
 * Modal dialog for adding/editing a single base name.
 * Always displays.
 *
 * @param {DialogProps} props
 */
export function BaseNameEdit(props: Props) {
  // Basic necessities
  const [error, setError] = useState('');
  const enqueueSnackbar = useEnqueueSnackbar();

  // Context
  const subscriberCtx = useContext(SubscriberContext);
  const subscriberIMSIs: Array<string> = Object.keys(subscriberCtx.state);
  const ctx = useContext(PolicyContext);
  const baseName =
    (props.baseNameId && ctx.baseNames[props.baseNameId]) || null;

  // Component state
  const [editedName, setEditedName] = useState(props.baseNameId || '');
  const [editedRuleNames, setEditedRuleNames] = useState(
    baseName?.rule_names || [],
  );
  const [editedAssignedSubscribers, setEditedAssignedSubscribers] = useState(
    baseName?.assigned_subscribers || [],
  );

  // Called when saving a base name, either new or edited
  const onSave = async () => {
    if (!props.baseNameId && ctx.baseNames[editedName]) {
      setError('Base name ID is already used. Please use a different name');
      return;
    }

    const savingRecord: BaseNameRecord = {
      name: editedName,
      rule_names: editedRuleNames,
      assigned_subscribers: editedAssignedSubscribers,
    };

    try {
      await ctx.setBaseNames(editedName, savingRecord);
      props.onSave();
      enqueueSnackbar('Base name saved successfully', {
        variant: 'success',
      });
    } catch (error) {
      setError(getErrorMessage(error));
    }
  };

  return (
    <>
      <DialogContent data-testid="baseNameEditDialog">
        <List>
          {error !== '' && (
            <AltFormField label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <div>
            <AltFormField label={'Base Name ID'}>
              <OutlinedInput
                data-testid="baseNameID"
                placeholder="Base Name 1"
                fullWidth={true}
                value={editedName}
                onChange={({target}) => setEditedName(target.value)}
              />
            </AltFormField>
            <AltFormField label={'Included Rule Names'}>
              <Select
                fullWidth={true}
                multiple
                variant={'outlined'}
                value={editedRuleNames}
                onChange={({target}) => {
                  // On autofill we get a stringified value.
                  setEditedRuleNames(
                    typeof target.value === 'string'
                      ? target.value.split(',')
                      : target.value,
                  );
                }}
                renderValue={selected => `${selected.length} rules`}
                input={<OutlinedInput id="ruleNames" />}>
                {Object.keys(ctx.state).map(profileID => (
                  <MenuItem key={profileID} value={profileID}>
                    <Checkbox checked={editedRuleNames.includes(profileID)} />
                    <ListItemText primary={profileID} />
                  </MenuItem>
                ))}
              </Select>
            </AltFormField>
            <AltFormField label={'Assigned Subscribers'}>
              <Select
                fullWidth={true}
                multiple
                variant={'outlined'}
                value={editedAssignedSubscribers}
                onChange={({target}) => {
                  // On autofill we get a stringified value.
                  setEditedAssignedSubscribers(
                    typeof target.value === 'string'
                      ? target.value.split(',')
                      : target.value,
                  );
                }}
                renderValue={selected => `${selected.length} subscribers`}
                input={<OutlinedInput id="assignedSubscribers" />}>
                {subscriberIMSIs.map(imsi => (
                  <MenuItem key={imsi} value={imsi}>
                    <Checkbox
                      checked={editedAssignedSubscribers.includes(imsi)}
                    />
                    <ListItemText primary={imsi} />
                  </MenuItem>
                ))}
              </Select>
            </AltFormField>
          </div>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Cancel</Button>
        <Button
          onClick={() => void onSave()}
          variant="contained"
          color="primary">
          Save
        </Button>
      </DialogActions>
    </>
  );
}
