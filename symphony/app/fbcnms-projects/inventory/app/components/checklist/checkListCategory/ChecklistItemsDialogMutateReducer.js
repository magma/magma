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
  CheckListItem,
  CheckListItemPendingFile,
} from './ChecklistItemsDialogMutateState';
import type {ChecklistItemsDialogMutateStateActionType} from './ChecklistItemsDialogMutateAction';
import type {ChecklistItemsDialogStateType} from './ChecklistItemsDialogMutateState';

import shortid from 'shortid';
import {reorder} from '../../draggable/DraggableUtils';

function editItem(
  state: ChecklistItemsDialogStateType,
  updatedItem: CheckListItem,
): ChecklistItemsDialogStateType {
  const {items} = state;
  const itemIndex = items.findIndex(i => i.id === updatedItem.id);
  return {
    ...state,
    items: [
      ...items.slice(0, itemIndex),
      {
        ...items[itemIndex],
        ...updatedItem,
      },
      ...items.slice(itemIndex + 1),
    ],
  };
}

export function reducer(
  state: ChecklistItemsDialogStateType,
  action: ChecklistItemsDialogMutateStateActionType,
): ChecklistItemsDialogStateType {
  const {items} = state;
  switch (action.type) {
    case 'EDIT_ITEM':
      return editItem(state, action.value);
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
    case 'EDIT_ITEM_PENDING_FILE':
      const item = items.find(i => i.id === action.itemId);
      if (item == null) {
        return state;
      }
      let pendingFiles: Array<CheckListItemPendingFile> =
        item.pendingFiles ?? [];
      const updatedPendingFile = action.file;
      const pendingFileIndex = pendingFiles.findIndex(
        f => f.id === action.file.id,
      );
      if (updatedPendingFile.progress === 100) {
        return editItem(state, {
          ...item,
          pendingFiles: pendingFiles.filter(f => f.id !== action.file.id),
        });
      }

      if (pendingFileIndex < 0) {
        pendingFiles = [...pendingFiles, updatedPendingFile];
      } else {
        pendingFiles = [
          ...pendingFiles.slice(0, pendingFileIndex),
          updatedPendingFile,
          ...pendingFiles.slice(pendingFileIndex + 1),
        ];
      }
      return editItem(state, {
        ...item,
        pendingFiles,
      });
    case 'ADD_ITEM_FILE':
      const addItemFileItem = items.find(i => i.id === action.itemId);
      if (addItemFileItem == null) {
        return state;
      }
      return editItem(state, {
        ...addItemFileItem,
        files: [...(addItemFileItem.files ?? []), action.file],
      });
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
