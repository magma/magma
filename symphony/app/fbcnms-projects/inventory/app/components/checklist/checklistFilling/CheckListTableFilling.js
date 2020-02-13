/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CheckCircle from '@material-ui/icons/CheckCircle';
import CheckListItem from '../CheckListItem';
import RadioButtonUnchecked from '@material-ui/icons/RadioButtonUnchecked';
import React, {useMemo} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import fbt from 'fbt';
import symphony from '@fbcnms/ui/theme/symphony';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import type {CheckListItem_item} from '../__generated__/CheckListItem_item.graphql';
import type {CheckListTableFilling_list} from './__generated__/CheckListTableFilling_list.graphql';

type Props = {
  list: CheckListTableFilling_list,
  onChecklistChanged?: (updatedList: CheckListTableFilling_list) => void,
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

const CheckListTableFilling = (props: Props) => {
  const {list, onChecklistChanged} = props;
  const classes = useStyles();

  const _updateList = (updatedList: CheckListTableFilling_list) => {
    if (!onChecklistChanged) {
      return;
    }

    onChecklistChanged(updatedList);
  };

  const _editItem = itemIndex => (updatedChecklistItem: CheckListItem_item) => {
    if (itemIndex < 0 || itemIndex >= list.length) {
      return;
    }

    const newItem: CheckListItem_item = {
      ...updatedChecklistItem,
      stringValue: updatedChecklistItem.stringValue || '',
      checked: updatedChecklistItem.checked || false,
    };

    const newList: CheckListTableFilling_list = [
      ...Array.prototype.slice.call(list, 0, itemIndex),
      newItem,
      ...Array.prototype.slice.call(list, itemIndex + 1),
    ];

    _updateList(newList);
  };

  const checklistItemsCount = list.length;
  const fullfilledItemsCount = list.reduce((fufilledSoFar, currentItem) => {
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

  const checklistItems = list.map((checkListItem, i) => (
    <TableRow id={checkListItem.id} index={i} key={checkListItem.id}>
      <TableCell className={classes.iconCell} size="small" component="div">
        {checkListItem.checked ? (
          <CheckCircle className={classes.checkIcon} />
        ) : (
          <RadioButtonUnchecked className={classes.checkIcon} />
        )}
      </TableCell>
      <TableCell component="div">
        <CheckListItem item={checkListItem || null} onChange={_editItem(i)} />
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

export default createFragmentContainer(CheckListTableFilling, {
  list: graphql`
    fragment CheckListTableFilling_list on CheckListItem @relay(plural: true) {
      id
      checked
      ...CheckListItem_item
    }
  `,
});
