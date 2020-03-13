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
  AddReportFilterMutationResponse,
  AddReportFilterMutationVariables,
  FilterEntity,
  FilterOperator,
} from '../mutations/__generated__/AddReportFilterMutation.graphql';
import type {FiltersQuery} from './comparison_view/ComparisonViewTypes';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithSnackbarProps} from 'notistack';

import * as React from 'react';

import AddReportFilterMutation from '../mutations/AddReportFilterMutation';
import BookmarksIcon from '@material-ui/icons/Bookmarks';
import BookmarksOutlinedIcon from '@material-ui/icons/BookmarksOutlined';
import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import DialogActions from '@material-ui/core/DialogActions';
import Popover from '@material-ui/core/Popover';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '../../../../fbcnms-packages/fbcnms-ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {withSnackbar} from 'notistack';

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
type Props = WithSnackbarProps & {
  isBookmark: boolean,
  filters: FiltersQuery,
  entity: FilterEntity,
};

const FilterBookmark = (props: Props) => {
  const {isBookmark} = props;
  const classes = useStyles();
  const [anchorEl, setAnchorEl] = React.useState(null);
  const [name, setName] = React.useState('');
  const [saving, setSaving] = React.useState(false);
  const [bookmarked, setBookmarked] = React.useState(isBookmark);

  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setName('');
    setAnchorEl(null);
  };
  const open = Boolean(anchorEl);

  const saveFilter = () => {
    saveFilterReport();
    handleClose();
  };

  const toOperator = (op: string): FilterOperator => {
    switch (op) {
      case 'is':
        return 'IS';
      case 'contains':
        return 'CONTAINS';
      case 'date_greater_than':
        return 'DATE_GREATER_THAN';
      case 'date_less_than':
        return 'DATE_LESS_THAN';
      case 'is_not_one_of':
        return 'IS_NOT_ONE_OF';
      case 'is_one_of':
        return 'IS_ONE_OF';
    }
    throw new Error(`Operator ${op} is not supported`);
  };

  const saveFilterReport = () => {
    setSaving(true);
    const filterInput = props.filters.map(f => {
      if (
        f.propertyValue &&
        (!f.propertyValue?.name || !f.propertyValue?.type)
      ) {
        throw new Error(`Property is not supported`);
      }
      return {
        filterType: f.name.toUpperCase(),
        operator: toOperator(f.operator),
        stringValue: f.stringValue,
        idSet: f.idSet,
        stringSet: f.stringSet,
        boolValue: f.boolValue,
        propertyValue: f.propertyValue
          ? {
              ...f.propertyValue,
              name: nullthrows(f.propertyValue?.name),
              type: nullthrows(f.propertyValue?.type),
            }
          : null,
      };
    });
    const variables: AddReportFilterMutationVariables = {
      input: {
        name: name,
        entity: props.entity,
        filters: filterInput,
      },
    };
    const callbacks: MutationCallbacks<AddReportFilterMutationResponse> = {
      onCompleted: (response, errors) => {
        setSaving(false);
        setBookmarked(true);
        if (errors && errors[0]) {
          props.enqueueSnackbar(errors[0].message, {
            children: key => (
              <SnackbarItem
                id={key}
                message={errors[0].message}
                variant="error"
              />
            ),
          });
        }
      },
      onError: (error: Error) => {
        setSaving(false);
        props.enqueueSnackbar(error.message, {
          children: key => (
            <SnackbarItem id={key} message={error.message} variant="error" />
          ),
        });
      },
    };
    AddReportFilterMutation(variables, callbacks);
  };

  return (
    <>
      <Button variant="text" skin="gray">
        {bookmarked ? (
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
        <>
          {saving ? (
            <CircularProgress size={24} />
          ) : (
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
          )}
        </>
      </Popover>
    </>
  );
};

export default withSnackbar(FilterBookmark);
