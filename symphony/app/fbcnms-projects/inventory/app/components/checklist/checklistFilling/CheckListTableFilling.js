/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ChecklistItemsDialogStateType} from '../checkListCategory/ChecklistItemsDialogMutateState';

import CheckCircle from '@material-ui/icons/CheckCircle';
import CheckListItem from '../CheckListItem';
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import RadioButtonUnchecked from '@material-ui/icons/RadioButtonUnchecked';
import React, {useContext, useMemo} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {makeStyles} from '@material-ui/styles';

type Props = {
  items: ChecklistItemsDialogStateType,
};

const useStyles = makeStyles(() => ({
  container: {
    maxWidth: '1366px',
    overflowX: 'auto',
  },
  root: {
    marginBottom: '12px',
    maxWidth: '100%',
  },
  cell: {
    paddingLeft: '0px',
    width: 'unset',
  },
  checkIcon: {
    fill: symphony.palette.D300,
    marginTop: '5px',
  },
  iconCell: {
    width: '20px',
  },
}));

const CheckListTableFilling = ({items}: Props) => {
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const classes = useStyles();

  const checklistItemsCount = items.length;
  const fullfilledItemsCount = items.reduce((fufilledSoFar, currentItem) => {
    if (currentItem.checked) {
      return fufilledSoFar + 1;
    }
    return fufilledSoFar;
  }, 0);

  // Could not use the FBT formation for plural here
  const tableHeader = useMemo(
    () => (
      <TableRow component="div">
        <TableCell padding="none" component="div" />
        <TableCell component="div">
          {checklistItemsCount > 0
            ? `${fbt(
                'Items',
                'Checklist items table header',
              )} (${fullfilledItemsCount}/${checklistItemsCount})`
            : fbt(
                'No Items',
                'Checklist items table header when there are no items in list',
              )}
        </TableCell>
      </TableRow>
    ),
    [checklistItemsCount, fullfilledItemsCount],
  );

  const checklistItems = items.map((checkListItem, i) => (
    <TableRow id={checkListItem.id} index={i} key={checkListItem.id}>
      <TableCell className={classes.iconCell} size="small" component="div">
        {checkListItem.checked ? (
          <CheckCircle className={classes.checkIcon} />
        ) : (
          <RadioButtonUnchecked className={classes.checkIcon} />
        )}
      </TableCell>
      <TableCell component="div">
        <CheckListItem
          item={checkListItem || null}
          onChange={updatedItem =>
            dispatch({
              type: 'EDIT_ITEM',
              value: updatedItem,
            })
          }
        />
      </TableCell>
    </TableRow>
  ));

  return (
    <div className={classes.container}>
      <Table component="div" className={classes.root}>
        <TableHead component="div">{tableHeader}</TableHead>
        <TableBody>{checklistItems}</TableBody>
      </Table>
    </div>
  );
};

export default CheckListTableFilling;
