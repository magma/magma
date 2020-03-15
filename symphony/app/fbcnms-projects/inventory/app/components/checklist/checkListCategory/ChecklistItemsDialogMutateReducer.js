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
  switch (action.type) {
    case 'EDIT_ITEM':
      const itemIndex = state.findIndex(i => i.id === action.value.id);
      return [
        ...state.slice(0, itemIndex),
        {
          ...state[itemIndex],
          ...action.value,
        },
        ...state.slice(itemIndex + 1),
      ];
    case 'ADD_ITEM':
      return [
        ...state,
        {
          id: shortid.generate(),
          title: '',
          type: 'simple',
          index: state.length,
        },
      ];
    case 'CHANGE_ITEM_POSITION':
      return reorder<CheckListItem>(
        state,
        action.sourceIndex,
        action.destinationIndex,
      ).map((item, index) => {
        return {
          ...item,
          index,
        };
      });
    case 'REMOVE_ITEM':
      const itemToRemoveIndex = state.findIndex(c => c.id === action.itemId);
      return [
        ...state.slice(0, itemToRemoveIndex),
        ...state.slice(itemToRemoveIndex + 1, state.length),
      ];
    default:
      return state;
  }
}

export function getInitialState(
  items: ChecklistItemsDialogStateType,
): ChecklistItemsDialogStateType {
  return items.slice().map(item => ({...item}));
}
