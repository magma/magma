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

import Button from '@fbcnms/ui/components/design-system/Button';
import CheckListItemDefinition, {
  CHECKLIST_ITEM_DEFINITION_TYPES,
} from './CheckListItemDefinition';
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import DeleteIcon from '@material-ui/icons/Delete';
import DraggableTableRow from '../../draggable/DraggableTableRow';
import DroppableTableBody from '../../draggable/DroppableTableBody';
import FormAction from '@fbcnms/ui/components/design-system/Form/FormAction';
import FormField from '@fbcnms/ui/components/design-system/FormField/FormField';
import React, {useContext, useMemo} from 'react';
import Select from '@fbcnms/ui/components/design-system/Select/Select';
import Table from '@material-ui/core/Table';
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
    width: '100%',
  },
  iconCell: {
    width: '20px',
    paddingRight: '0px',
  },
  addButton: {
    height: '32px',
    padding: '4px 18px',
    borderRadius: '4px',
    border: '1px solid',
    borderColor: symphony.palette.primary,
    '&:hover': {
      borderColor: symphony.palette.B800,
      backgroundColor: symphony.palette.B50,
    },
  },
}));

const CheckListTableDefinition = ({items}: Props) => {
  const classes = useStyles();
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);

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

  const checklistTypes = useMemo(
    () =>
      Object.keys(CHECKLIST_ITEM_DEFINITION_TYPES).map(type => ({
        key: type,
        label: CHECKLIST_ITEM_DEFINITION_TYPES[type].description,
        value: type,
      })),
    [],
  );
  const checklistItems = items.map((checkListItem, i) => (
    <DraggableTableRow id={checkListItem.id} index={i} key={i}>
      <TableCell className={classes.cell} size="small" component="div">
        <FormField>
          <Select
            className={classes.selectMenu}
            options={checklistTypes}
            selectedValue={checkListItem.type}
            onChange={value =>
              dispatch({
                type: 'EDIT_ITEM',
                value: {
                  ...checkListItem,
                  type: value,
                },
              })
            }
          />
        </FormField>
      </TableCell>
      <TableCell component="div">
        <CheckListItemDefinition
          item={checkListItem}
          onChange={item =>
            dispatch({
              type: 'EDIT_ITEM',
              value: item,
            })
          }
        />
      </TableCell>
      <TableCell className={classes.iconCell} align="right" component="div">
        <FormAction>
          <Button
            skin="primary"
            variant="text"
            onClick={() =>
              dispatch({
                type: 'REMOVE_ITEM',
                itemId: checkListItem.id,
              })
            }>
            <DeleteIcon />
          </Button>
        </FormAction>
      </TableCell>
    </DraggableTableRow>
  ));

  return (
    <div className={classes.container}>
      <Table component="div" className={classes.table}>
        <TableHead component="div">{tableHeader}</TableHead>
        <DroppableTableBody
          onDragEnd={positionChange =>
            dispatch({
              type: 'CHANGE_ITEM_POSITION',
              sourceIndex: positionChange.source.index,
              destinationIndex: positionChange.destination.index,
            })
          }>
          {checklistItems}
        </DroppableTableBody>
      </Table>
      <FormAction>
        <Button
          className={classes.addButton}
          color="primary"
          variant="text"
          onClick={() => dispatch({type: 'ADD_ITEM'})}>
          {fbt(
            'Add Item',
            'Caption of the Add Checklist Item button (under the checklist table)',
          )}
        </Button>
      </FormAction>
    </div>
  );
};

export default CheckListTableDefinition;
