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
} from '../mutations/__generated__/AddReportFilterMutation.graphql';
import type {
  DeleteReportFilterMutationResponse,
  DeleteReportFilterMutationVariables,
} from '../mutations/__generated__/DeleteReportFilterMutation.graphql';
import type {
  EditReportFilterMutationResponse,
  EditReportFilterMutationVariables,
} from '../mutations/__generated__/EditReportFilterMutation.graphql';
import type {FiltersQuery} from './comparison_view/ComparisonViewTypes';
import type {MutationCallbacks} from '../mutations/MutationCallbacks.js';
import type {WithSnackbarProps} from 'notistack';

import * as React from 'react';

import AddReportFilterMutation from '../mutations/AddReportFilterMutation';
import BookmarksIcon from '@material-ui/icons/Bookmarks';
import BookmarksOutlinedIcon from '@material-ui/icons/BookmarksOutlined';
import Button from '@fbcnms/ui/components/design-system/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import DeleteOutlineIcon from '@material-ui/icons/DeleteOutline';
import DeleteReportFilterMutation from '../mutations/DeleteReportFilterMutation';
import DialogActions from '@material-ui/core/DialogActions';
import EditReportFilterMutation from '../mutations/EditReportFilterMutation';
import Popover from '@material-ui/core/Popover';
import SnackbarItem from '@fbcnms/ui/components/SnackbarItem';
import Strings from '../common/CommonStrings';
import Text from '@fbcnms/ui/components/design-system/Text';
import TextInput from '@fbcnms/ui/components/design-system/Input/TextInput';
import fbt from 'fbt';
import nullthrows from '@fbcnms/util/nullthrows';
import symphony from '../../../../fbcnms-packages/fbcnms-ui/theme/symphony';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {makeStyles} from '@material-ui/styles';
import {stringToOperator} from './comparison_view/FilterUtils';
import {useEffect} from 'react';
import {usePowerSearch} from './power_search/PowerSearchContext';
import {withSnackbar} from 'notistack';

const useStyles = makeStyles(() => ({
  filledBookmarkButton: {
    cursor: 'pointer',
    color: symphony.palette.B600,
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'center',
    margin: '0px 4px 0px 8px',
    '&:hover:not($disabled)': {
      color: symphony.palette.B700,
    },
  },
  bookmarkButton: {
    color: symphony.palette.D500,
    margin: '0px 4px 0px 8px',
    '&:hover:not($disabled)': {
      color: symphony.palette.D900,
    },
  },
  dialogActions: {
    padding: '8px 0px 0px 0px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  popup: {
    backgroundColor: 'white',
    width: '320px',
    padding: '16px 16px',
  },
  text: {
    margin: '8px 2px',
  },
}));
type Props = WithSnackbarProps & {
  filters: FiltersQuery,
  entity: FilterEntity,
};

const FilterBookmark = (props: Props) => {
  const classes = useStyles();
  const {filters, entity} = props;
  const {bookmark, setBookmark} = usePowerSearch();

  const [anchorEl, setAnchorEl] = React.useState(null);
  const [name, setName] = React.useState('');
  const [saving, setSaving] = React.useState(false);

  useEffect(() => {
    setName(bookmark?.name ?? '');
  }, [bookmark]);

  const handleClick = event => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setSaving(false);
    setAnchorEl(null);
  };

  const saveBookmark = () => {
    saveFilterReport();
    handleClose();
  };

  const toCapitalLetter = (x: string): string => {
    return x[0].toUpperCase() + x.substring(1).toLowerCase();
  };

  const entityToLabel = (entity: FilterEntity): string => {
    let entitySplit = entity.split('_');
    entitySplit = entitySplit.map(w => toCapitalLetter(w));
    return entitySplit.join(' ');
  };

  const filtersQueryToFilterInput = (filterQuery: FiltersQuery) => {
    return filterQuery.map(f => {
      if (
        f.propertyValue &&
        (!f.propertyValue?.name || !f.propertyValue?.type)
      ) {
        throw new Error(`Property is not supported`);
      }
      return {
        filterType: f.name.toUpperCase(),
        key: f.key,
        operator: stringToOperator(f.operator),
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
  };

  const removeBookmark = () => {
    setSaving(true);
    const variables: DeleteReportFilterMutationVariables = {
      id: nullthrows(bookmark?.id),
    };
    const callbacks: MutationCallbacks<DeleteReportFilterMutationResponse> = {
      onCompleted: (response, errors) => {
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
        } else {
          ServerLogger.info(LogEvents.SAVED_SEARCH_DELETED, {
            bookmark_name: bookmark?.name,
          });
        }
        handleClose();
        setBookmark(null);
      },
      onError: (error: Error) => {
        props.enqueueSnackbar(error.message, {
          children: key => (
            <SnackbarItem id={key} message={error.message} variant="error" />
          ),
        });
        handleClose();
        setBookmark(null);
      },
    };
    DeleteReportFilterMutation(variables, callbacks);
  };

  const editBookmark = () => {
    setSaving(true);
    const variables: EditReportFilterMutationVariables = {
      input: {
        name: name,
        id: nullthrows(bookmark?.id),
      },
    };
    const callbacks: MutationCallbacks<EditReportFilterMutationResponse> = {
      onCompleted: (response, errors) => {
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
        } else {
          ServerLogger.info(LogEvents.SAVED_SEARCH_EDITED, {
            bookmark_name: name,
          });
        }
        handleClose();
        setBookmark({
          id: response.editReportFilter.id,
          name: response.editReportFilter.name,
        });
      },
      onError: (error: Error) => {
        handleClose();
        setBookmark(null);
        props.enqueueSnackbar(error.message, {
          children: key => (
            <SnackbarItem id={key} message={error.message} variant="error" />
          ),
        });
      },
    };
    EditReportFilterMutation(variables, callbacks);
  };

  const saveFilterReport = () => {
    setSaving(true);
    const filterInput = filtersQueryToFilterInput(filters);
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
        } else {
          ServerLogger.info(LogEvents.SAVED_SEARCH_CREATED, {
            bookmark_name: name,
          });
        }
        setBookmark({
          id: response.addReportFilter.id,
          name: response.addReportFilter.name,
        });
      },
      onError: (error: Error) => {
        setSaving(false);
        setBookmark(null);
        props.enqueueSnackbar(error.message, {
          children: key => (
            <SnackbarItem id={key} message={error.message} variant="error" />
          ),
        });
      },
    };
    AddReportFilterMutation(variables, callbacks);
  };
  const isBookmark = bookmark != null;
  return (
    <>
      <div onClick={handleClick}>
        <Button variant="text" skin="gray">
          {isBookmark ? (
            <BookmarksIcon
              className={classes.filledBookmarkButton}
              color="inherit"
            />
          ) : (
            <BookmarksOutlinedIcon
              className={classes.bookmarkButton}
              color="inherit"
            />
          )}
        </Button>
      </div>
      <Popover
        open={Boolean(anchorEl)}
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
                <Text variant="body1" color="regular">
                  <fbt desc="">Saved Search</fbt>
                </Text>
                <div>
                  <Text variant="subtitle2" color="gray">
                    <fbt desc="">
                      You can find it under the
                      <fbt:param name="entity name">
                        {"'" + entityToLabel(entity) + "'"}
                      </fbt:param>{' '}
                      filter search bar.
                    </fbt>
                  </Text>
                </div>
              </div>
              <TextInput
                type="string"
                placeholder={
                  isBookmark ? name : `${fbt('Saved search name', '')}`
                }
                onChange={({target}) => setName(target.value)}
                value={name}
              />
              <DialogActions classes={{root: classes.dialogActions}}>
                {isBookmark ? (
                  <Button variant="text" skin="gray" onClick={removeBookmark}>
                    <DeleteOutlineIcon />
                  </Button>
                ) : (
                  <div />
                )}
                <div>
                  <Button onClick={handleClose} skin="regular">
                    {Strings.common.cancelButton}
                  </Button>
                  {isBookmark ? (
                    <Button
                      disabled={name.trim() == '' || filters.length === 0}
                      onClick={editBookmark}>
                      {Strings.common.saveButton}
                    </Button>
                  ) : (
                    <Button
                      disabled={name.trim() == '' || filters.length === 0}
                      onClick={saveBookmark}>
                      {Strings.common.createButton}
                    </Button>
                  )}
                </div>
              </DialogActions>
            </div>
          )}
        </>
      </Popover>
    </>
  );
};

export default withSnackbar(FilterBookmark);
