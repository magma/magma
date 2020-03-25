/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CheckListItem} from './ChecklistItemsDialogMutateState';
import type {ChecklistItemsDialogMutateStateActionType} from './ChecklistItemsDialogMutateAction';
import type {ChecklistItemsDialogStateType} from './ChecklistItemsDialogMutateState';

import shortid from 'shortid';
import {reorder} from '../../draggable/DraggableUtils';

export function reducer(
  state: ChecklistItemsDialogStateType,
  action: ChecklistItemsDialogMutateStateActionType,
): ChecklistItemsDialogStateType {
  const {items} = state;
  switch (action.type) {
    case 'EDIT_ITEM':
      const itemIndex = items.findIndex(i => i.id === action.value.id);
      return {
        ...state,
        items: [
          ...items.slice(0, itemIndex),
          {
            ...items[itemIndex],
            ...action.value,
          },
          ...items.slice(itemIndex + 1),
        ],
      };
    case 'ADD_ITEM':
      const newId = shortid.generate();
      return {
        editedDefinitionId: newId,
        items: [
          ...items,
          {
            id: newId,
            title: '',
            type: 'simple',
            index: items.length,
          },
        ],
      };
    case 'CHANGE_ITEM_POSITION':
      return {
        ...state,
        items: reorder<CheckListItem>(
          items,
          action.sourceIndex,
          action.destinationIndex,
        ).map((item, index) => {
          return {
            ...item,
            index,
          };
        }),
      };
    case 'REMOVE_ITEM':
      const itemToRemoveIndex = items.findIndex(c => c.id === action.itemId);
      return {
        ...state,
        items: [
          ...items.slice(0, itemToRemoveIndex),
          ...items.slice(itemToRemoveIndex + 1, items.length),
        ],
      };
    case 'SET_EDITED_DEFINITION_ID':
      return {
        ...state,
        editedDefinitionId: action.itemId,
      };
    default:
      return state;
  }
}

export function getInitialState(
  initialItems: Array<CheckListItem>,
): ChecklistItemsDialogStateType {
  return {
    items: initialItems.slice().map(item => ({...item})),
    editedDefinitionId: null,
  };
}
