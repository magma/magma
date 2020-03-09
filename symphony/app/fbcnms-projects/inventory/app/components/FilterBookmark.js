/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {
  EntityType,
  FiltersQuery,
} from './comparison_view/ComparisonViewTypes';

import * as React from 'react';

import BookmarksIcon from '@material-ui/icons/Bookmarks';
import BookmarksOutlinedIcon from '@material-ui/icons/BookmarksOutlined';
import Button from '@fbcnms/ui/components/design-system/Button';
import DialogActions from '@material-ui/core/DialogActions';
import Popover from '@material-ui/core/Popover';
import Strings from '../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import symphony from '../../../../fbcnms-packages/fbcnms-ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(() => ({
  filledBookmarkButton: {
    cursor: 'pointer',
    color: symphony.palette.B600,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    margin: '0px 4px',
    '&:hover:not($disabled)': {
      color: symphony.palette.B700,
    },
  },
  bookmarkButton: {
    color: symphony.palette.D500,
    margin: '0px 4px',
    '&:hover:not($disabled)': {
      color: symphony.palette.D900,
    },
  },
  dialogActions: {
    padding: '8px 0px 0px 0px',
  },
  popup: {
    backgroundColor: 'white',
    width: '350px',
    padding: '16px 16px',
  },
  text: {
    margin: '8px 2px',
  },
}));
type Props = {
  isBookmark: boolean,
  filters: FiltersQuery,
  entity: EntityType,
};

const FilterBookmark = (props: Props) => {
  const {isBookmark} = props;
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [name, setName] = React.useState('');

  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setName('');
    setAnchorEl(null);
  };
  const open = Boolean(anchorEl);

  const saveFilter = () => {
    handleClose();
  };

  return (
    <>
      <Button variant="text" skin="gray">
        {isBookmark ? (
          <BookmarksIcon
            className={classes.filledBookmarkButton}
            color="inherit"
            onClick={handleClick}
          />
        ) : (
          <BookmarksOutlinedIcon
            className={classes.bookmarkButton}
            color="inherit"
            onClick={handleClick}
          />
        )}
      </Button>
      <Popover
        open={open}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'center',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}>
        <div className={classes.popup}>
          <div className={classes.text}>
            <Text variant="body2" color="regular">
              <fbt desc="">SAVE SEARCH</fbt>
            </Text>
          </div>
          <TextInput
            type="string"
            placeholder={`${fbt('Bookmark name', '')}`}
            onChange={({target}) => setName(target.value)}
            value={name}
          />
          <DialogActions classes={{root: classes.dialogActions}}>
            <Button onClick={handleClose} skin="regular">
              {Strings.common.cancelButton}
            </Button>
            <Button disabled={name == ''} onClick={saveFilter}>
              {Strings.common.saveButton}
            </Button>
          </DialogActions>
        </div>
      </Popover>
    </>
  );
};

export default FilterBookmark;
