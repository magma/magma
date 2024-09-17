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

import type {RatingGroup} from '../../../generated';

import Button from '@mui/material/Button';
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

import {AltFormField} from '../../components/FormField';
import {getErrorMessage} from '../../util/ErrorUtils';
import {useContext, useEffect, useState} from 'react';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

type Props = {
  open: boolean;
  onClose: () => void;
  ratingGroup?: RatingGroup;
};

export default function RatingGroupEditDialog(props: Props) {
  const ctx = useContext(PolicyContext);
  const ratingGroups = ctx.ratingGroups;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');

  const [ratingGroup, setRatingGroup] = useState<RatingGroup>(
    props.ratingGroup || {
      id: 0,
      limit_type: 'FINITE',
    },
  );

  useEffect(() => {
    setRatingGroup(
      props.ratingGroup || {
        id: 0,
        limit_type: 'FINITE',
      },
    );
    setError('');
  }, [props.open, props.ratingGroup]);

  const isAdd = props.ratingGroup ? false : true;
  const handleRatingGroupChange = (key: string, val: string | number) => {
    setRatingGroup({...ratingGroup, [key]: val});
  };
  const onSave = async () => {
    try {
      if (isAdd) {
        if (isNaN(ratingGroup.id)) {
          setError('empty Rating Group id');
          return;
        }
        if (ratingGroup.id in ratingGroups) {
          setError(`Rating Group ${ratingGroup.id} already exists`);
          return;
        }
      }
      await ctx.setRatingGroups(ratingGroup.id.toString(), ratingGroup);
      enqueueSnackbar('Rating Group saved successfully', {
        variant: 'success',
      });

      props.onClose();
    } catch (error) {
      setError(getErrorMessage(error));
    }
  };

  const onClose = () => {
    props.onClose();
  };

  return (
    <Dialog
      data-testid="editDialog"
      open={props.open}
      scroll="body"
      fullWidth={true}
      maxWidth={'md'}>
      <DialogTitle
        onClose={onClose}
        label={props.ratingGroup ? 'Edit Rating Group' : 'Add New Rating Group'}
      />
      <DialogContent>
        <List>
          {error !== '' && (
            <AltFormField disableGutters label={''}>
              <FormLabel data-testid="configEditError" error>
                {error}
              </FormLabel>
            </AltFormField>
          )}
          <AltFormField label={'Rating Group ID'} disableGutters>
            <OutlinedInput
              fullWidth={true}
              type="number"
              data-testid="ratingGroupID"
              placeholder="Eg. 1"
              value={ratingGroup.id}
              onChange={({target}) =>
                handleRatingGroupChange('id', parseInt(target.value) || NaN)
              }
            />
          </AltFormField>
          <AltFormField disableGutters label={'Limit Type'}>
            <Select
              fullWidth={true}
              variant={'outlined'}
              value={ratingGroup.limit_type || 'FINITE'}
              onChange={({target}) => {
                handleRatingGroupChange('limit_type', target.value);
              }}
              input={<OutlinedInput />}>
              <MenuItem value={'FINITE'}>
                <ListItemText primary={'FINITE'} />
              </MenuItem>
              <MenuItem value={'INFINITE_UNMETERED'}>
                <ListItemText primary={'INFINITE_UNMETERED'} />
              </MenuItem>
              <MenuItem value={'INFINITE_METERED'}>
                <ListItemText primary={'INFINITE_METERED'} />
              </MenuItem>
            </Select>
          </AltFormField>
        </List>
      </DialogContent>
      <DialogActions>
        <Button onClick={props.onClose}>Close</Button>
        <Button
          variant="contained"
          color="primary"
          onClick={() => void onSave()}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
