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
  CheckListItemFile,
} from './ChecklistItemsDialogMutateState';
import type {CheckListItemPendingFile} from './ChecklistItemsDialogMutateState';

export type ChecklistItemsDialogMutateStateActionType =
  | {|
      type: 'EDIT_ITEM',
      value: CheckListItem,
    |}
  | {|
      type: 'ADD_ITEM',
    |}
  | {|
      type: 'CHANGE_ITEM_POSITION',
      sourceIndex: number,
      destinationIndex: number,
    |}
  | {|
      type: 'REMOVE_ITEM',
      itemId: string,
    |}
  | {|
      type: 'SET_EDITED_DEFINITION_ID',
      itemId: string,
    |}
  | {|
      type: 'EDIT_ITEM_PENDING_FILE',
      itemId: string,
      file: CheckListItemPendingFile,
    |}
  | {|
      type: 'ADD_ITEM_FILE',
      itemId: string,
      file: CheckListItemFile,
    |};

export type ChecklistItemsDialogMutateDispatcher = (
  action: ChecklistItemsDialogMutateStateActionType,
) => void;
