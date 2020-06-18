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
  CheckListItemDefinition,
} from './checkListCategory/ChecklistItemsDialogMutateState';

export type ChecklistCategoryDefinition = $ReadOnly<{|
  id: string,
  title: string,
  description: ?string,
  checkList: Array<CheckListItemDefinition>,
|}>;

export type ChecklistCategory = $ReadOnly<{|
  id: string,
  key?: string,
  title: string,
  description?: ?string,
  checkList: Array<CheckListItem>,
|}>;

export type ChecklistCategoriesStateType = Array<ChecklistCategory>;
