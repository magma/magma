/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CheckListItem} from './checkListCategory/ChecklistItemsDialogMutateState';

export type ChecklistCategoriesMutateStateActionType =
  | {|
      type: 'UPDATE_CATEGORY_TITLE',
      categoryId: string,
      value: string,
    |}
  | {|
      type: 'UPDATE_CATEGORY_DESCRIPTION',
      categoryId: string,
      value: string,
    |}
  | {|
      type: 'UPDATE_CATEGORY_CHECKLIST',
      categoryId: string,
      value: Array<CheckListItem>,
    |}
  | {|
      type: 'ADD_CATEGORY',
    |}
  | {|
      type: 'REMOVE_CATEGORY',
      categoryId: string,
    |};
