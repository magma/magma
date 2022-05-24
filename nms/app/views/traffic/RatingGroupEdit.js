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

import type {rating_group} from '../../../generated/MagmaAPIBindings';

import Button from '@material-ui/core/Button';
import Dialog from '@material-ui/core/Dialog';
import DialogActions from '@material-ui/core/DialogActions';
import DialogContent from '@material-ui/core/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import FormLabel from '@material-ui/core/FormLabel';
import List from '@material-ui/core/List';
import ListItem from '@material-ui/core/ListItem';
import ListItemText from '@material-ui/core/ListItemText';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import PolicyContext from '../../components/context/PolicyContext';
import React from 'react';
import Select from '@material-ui/core/Select';

import {AltFormField} from '../../components/FormField';
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext, useEffect, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

const useStyles = makeStyles(() => ({
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
  input: {
    display: 'inline-flex',
    margin: '5px 0',
    width: '50%',
    fullWidth: true,
  },
}));

type Props = {
  open: boolean,
  onClose: () => void,
  ratingGroup?: rating_group,
};

export default function RatingGroupEditDialog(props: Props) {
  const classes = useStyles();
  const ctx = useContext(PolicyContext);
  const ratingGroups = ctx.ratingGroups;
  const enqueueSnackbar = useEnqueueSnackbar();
  const [error, setError] = useState('');

  const [ratingGroup, setRatingGroup] = useState<rating_group>(
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
  const handleRatingGroupChange = (key: string, val) => {
    setRatingGroup({...ratingGroup, [key]: val});
  };
  const onSave = async () => {
    try {
      if (isAdd) {
        if (ratingGroup.id === '') {
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
    } catch (e) {
      setError(e.response?.data?.message ?? e.message);
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
          <ListItem dense disableGutters />
          <AltFormField label={'Rating Group ID'} disableGutters>
            <OutlinedInput
              className={classes.input}
              fullWidth={true}
              type="number"
              data-testid="ratingGroupID"
              placeholder="Eg. 1"
              value={ratingGroup.id}
              onChange={({target}) =>
                handleRatingGroupChange('id', parseInt(target.value) || '')
              }
            />
          </AltFormField>
          <AltFormField disableGutters label={'Limit Type'}>
            <Select
              fullWidth={true}
              variant={'outlined'}
              // eslint-disable-next-line no-warning-comments
              // $FlowFixMe redirect_info type needed
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
        <Button variant="contained" color="primary" onClick={onSave}>
          Save
        </Button>
      </DialogActions>
    </Dialog>
  );
}
