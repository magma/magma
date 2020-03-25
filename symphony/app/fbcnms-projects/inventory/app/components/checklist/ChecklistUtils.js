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

export const isChecklistItemDone = (item: CheckListItem): boolean => {
  switch (item.type) {
    case 'enum':
      return item.enumValues != null && item.enumValues.trim().length > 0;
    case 'simple':
      return item.checked === true;
    case 'string':
      return item.stringValue != null && item.stringValue.trim() !== '';
    default:
      throw new Error(
        `Invariant violation - checklist item type not found: ${item.type}`,
      );
  }
};

export const enumStringToArray = (enumString: ?string): Array<string> => {
  return enumString != null && enumString !== ''
    ? enumString.split(',')
    : ([]: Array<string>);
};
