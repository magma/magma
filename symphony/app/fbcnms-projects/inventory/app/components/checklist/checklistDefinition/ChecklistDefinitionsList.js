/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItem} from '../checkListCategory/ChecklistItemsDialogMutateState';

import CheckListItemDefinition from './CheckListItemDefinition';
import ChecklistItemsDialogMutateDispatchContext from '../checkListCategory/ChecklistItemsDialogMutateDispatchContext';
import React, {useContext} from 'react';
import classNames from 'classnames';
import symphony from '@fbcnms/ui/theme/symphony';
import {DragDropContext, Droppable} from 'react-beautiful-dnd';
import {Draggable} from 'react-beautiful-dnd';
import {ReorderIcon} from '@fbcnms/ui/components/design-system/Icons';
import {makeStyles} from '@material-ui/styles';
import {sortByIndex} from '../../draggable/DraggableUtils';
import {useFormContext} from '../../../common/FormContext';

type Props = {
  items: Array<CheckListItem>,
  editedDefinitionId: ?string,
};

const useStyles = makeStyles(() => ({
  itemsList: {
    paddingTop: '16px',
    position: 'relative',
  },
  listItem: {
    display: 'flex',
    flexDirection: 'row',
    alignItems: 'center',
    '&:not(:last-child)': {
      marginBottom: '8px',
    },
  },
  dragIndicatorIcon: {
    cursor: 'grab',
    fill: symphony.palette.D300,
  },
  dragHandle: {
    position: 'absolute',
    left: '-24px',
  },
  hiddenDragHandle: {
    display: 'none',
  },
}));

const ChecklistDefinitionsList = ({items, editedDefinitionId}: Props) => {
  const classes = useStyles();
  const dispatch = useContext(ChecklistItemsDialogMutateDispatchContext);
  const form = useFormContext();

  const checklistItems = items.sort(sortByIndex).map(item => {
    return (
      <Draggable
        key={item.id}
        draggableId={item.id}
        index={item.index}
        isDragDisabled={
          form.alerts.missingPermissions.detected ||
          item.id !== editedDefinitionId
        }>
        {(provided, snapshot) => (
          <div
            className={classes.listItem}
            ref={provided.innerRef}
            {...provided.draggableProps}>
            <div
              {...provided.dragHandleProps}
              className={classNames(classes.dragHandle, {
                [classes.hiddenDragHandle]: item.id !== editedDefinitionId,
              })}>
              <ReorderIcon className={classes.dragIndicatorIcon} />
            </div>
            <CheckListItemDefinition
              item={item}
              editedDefinitionId={
                snapshot.isDragging ? null : editedDefinitionId
              }
              onChange={item =>
                dispatch({
                  type: 'EDIT_ITEM',
                  value: item,
                })
              }
            />
          </div>
        )}
      </Draggable>
    );
  });

  return (
    <div className={classes.itemsList}>
      <DragDropContext
        onDragEnd={result => {
          dispatch({
            type: 'CHANGE_ITEM_POSITION',
            sourceIndex: result.source.index,
            destinationIndex: result.destination.index,
          });
        }}>
        <Droppable droppableId="checklist_definitions_droppable">
          {provided => (
            <div ref={provided.innerRef} {...provided.droppableProps}>
              {checklistItems}
              {provided.placeholder}
            </div>
          )}
        </Droppable>
      </DragDropContext>
    </div>
  );
};

export default ChecklistDefinitionsList;
