/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Button from '@fbcnms/ui/components/Button';
import CheckListItem, {
  CHECKLIST_ITEM_TYPES,
  GetValidChecklistItemType,
} from '../CheckListItem';
import DeleteIcon from '@material-ui/icons/Delete';
import DraggableTableRow from '../../draggable/DraggableTableRow';
import DroppableTableBody from '../../draggable/DroppableTableBody';
import IconButton from '@material-ui/core/IconButton';
import MenuItem from '@material-ui/core/MenuItem';
import React, {useMemo, useState} from 'react';
import Table from '@material-ui/core/Table';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import TextField from '@material-ui/core/TextField';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import {makeStyles} from '@material-ui/styles';
import {reorder} from '../../draggable/DraggableUtils';
import type {CheckListItemType} from '../../configure/__generated__/AddEditWorkOrderTypeCard_editingWorkOrderType.graphql.js';
import type {CheckListTableDefinition_list} from './__generated__/CheckListTableDefinition_list.graphql';
import type {ChecklistItemInput} from '../../../mutations/__generated__/AddWorkOrderMutation.graphql';

type Props = {
  list: CheckListTableDefinition_list,
  onChecklistChanged?: (updatedList: CheckListTableDefinition_list) => void,
};

const useStyles = makeStyles({
  container: {
    maxWidth: '1366px',
    overflowX: 'auto',
  },
  table: {
    marginBottom: '12px',
    maxWidth: '100%',
  },
  cell: {
    paddingLeft: '0px',
    width: 'unset',
  },
  input: {
    marginTop: '0px',
    marginBottom: '0px',
    width: '100%',
  },
  selectMenu: {
    height: '14px',
  },
  iconCell: {
    width: '20px',
    paddingRight: '0px',
  },
  addButton: {
    marginBottom: '12px',
  },
});

const CheckListTableDefinition = (props: Props) => {
  const {list, onChecklistChanged} = props;
  const classes = useStyles();
  const [nextNewItemTempId, setNextNewItemTempId] = useState(1);

  const _updateList = (updatedList: CheckListTableDefinition_list) => {
    if (!onChecklistChanged) {
      return;
    }

    onChecklistChanged(updatedList);
  };

  const _removeItem = (item, itemIndex) => {
    _updateList([...list.slice(0, itemIndex), ...list.slice(itemIndex + 1)]);
  };

  const _createNewItem: () => ChecklistItemInput = () => {
    const newId = nextNewItemTempId;
    setNextNewItemTempId(newId + 1);
    return {
      id: `@tmp${newId}`,
      title: '',
      type: 'simple',
      index: list.length,
    };
  };

  const _addItem = () => _updateList([...list, _createNewItem()]);

  const _editItem = itemIndex => (updatedChecklistItem: ChecklistItemInput) => {
    if (itemIndex < 0 || itemIndex >= list.length) {
      return;
    }

    const newList: CheckListTableDefinition_list = [
      ...Array.prototype.slice.call(list, 0, itemIndex),
      updatedChecklistItem,
      ...Array.prototype.slice.call(list, itemIndex + 1),
    ];

    _updateList(newList);
  };

  const _changeItemType = (checklistItemIndex: number, newType: string) => {
    const checklistItemType: ?CheckListItemType = GetValidChecklistItemType(
      newType,
    );
    if (!checklistItemType) {
      return;
    }
    const newItem = {
      ...list[checklistItemIndex],
      type: checklistItemType,
    };

    _editItem(checklistItemIndex)(newItem);
  };

  const _changeItemPosition = positionChange => {
    if (!positionChange.destination) {
      return;
    }

    const updatedList = reorder(
      list,
      positionChange.source.index,
      positionChange.destination.index,
    ).map((item, index) => {
      return {
        ...item,
        index,
      };
    });

    _updateList(updatedList);
  };

  const tableHeader = useMemo(
    () => (
      <TableRow component="div">
        <TableCell padding="none" component="div" />
        <TableCell size="small" component="div" className={classes.cell}>
          {fbt(
            'Type',
            'Title of the check list type column at the check list definition table',
          )}
        </TableCell>
        <TableCell component="div">
          {fbt(
            'Definition',
            'Title of the check list definition column at the check list definition table',
          )}
        </TableCell>
        <TableCell component="div" />
      </TableRow>
    ),
    [classes.cell],
  );

  const checklistItems = list.map((checkListItem, i) => (
    <DraggableTableRow id={checkListItem.id} index={i} key={i}>
      <TableCell className={classes.cell} size="small" component="div">
        <TextField
          select
          variant="outlined"
          className={classes.input}
          value={checkListItem.type}
          SelectProps={{
            classes: {selectMenu: classes.selectMenu},
            MenuProps: {
              className: classes.menu,
            },
          }}
          onChange={event => _changeItemType(i, event.target.value)}
          margin="dense">
          {Object.keys(CHECKLIST_ITEM_TYPES).map(type => (
            <MenuItem key={type} value={type}>
              {CHECKLIST_ITEM_TYPES[type].description}
            </MenuItem>
          ))}
        </TextField>
      </TableCell>
      <TableCell component="div">
        <CheckListItem
          item={checkListItem}
          designMode={true}
          onChange={_editItem(i)}
        />
      </TableCell>
      <TableCell className={classes.iconCell} align="right" component="div">
        <IconButton onClick={() => _removeItem(checkListItem, i)}>
          <DeleteIcon />
        </IconButton>
      </TableCell>
    </DraggableTableRow>
  ));

  return (
    <div className={classes.container}>
      <Table component="div" className={classes.table}>
        <TableHead component="div">{tableHeader}</TableHead>
        <DroppableTableBody onDragEnd={_changeItemPosition}>
          {checklistItems}
        </DroppableTableBody>
      </Table>
      <Button
        className={classes.addButton}
        color="primary"
        variant="outlined"
        onClick={_addItem}>
        {fbt(
          'Add Item',
          'Caption of the Add Checklist Item button (under the checklist table)',
        )}
      </Button>
    </div>
  );
};

export default createFragmentContainer(CheckListTableDefinition, {
  list: graphql`
    fragment CheckListTableDefinition_list on CheckListItem
      @relay(plural: true) {
      id
      type
      index
      ...CheckListItem_item
    }
  `,
});
