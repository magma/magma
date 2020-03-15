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
    |};

export type ChecklistItemsDialogMutateDispatcher = (
  action: ChecklistItemsDialogMutateStateActionType,
) => void;
