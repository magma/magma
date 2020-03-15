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

export type ChecklistCategory = {|
  id: string,
  key?: string,
  title: string,
  description?: ?string,
  checkList: Array<CheckListItem>,
|};

export type ChecklistCategoriesStateType = Array<ChecklistCategory>;
