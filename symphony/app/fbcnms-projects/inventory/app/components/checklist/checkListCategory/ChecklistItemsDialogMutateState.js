/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {CheckListItemType} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';

export type CheckListItemBase = $ReadOnly<{|
  id: string,
  index?: ?number,
  type: CheckListItemType,
  title: string,
  helpText?: ?string,
|}>;

export type BasicCheckListItemData = {|
  checked?: ?boolean,
|};

export type EnumCheckListItemData = {|
  enumValues?: ?string,
|};

export type FreeTextCheckListItemData = {|
  stringValue?: ?string,
|};

export type CheckListItem = {|
  ...CheckListItemBase,
  ...BasicCheckListItemData,
  ...EnumCheckListItemData,
  ...FreeTextCheckListItemData,
|};

export type ChecklistItemsDialogStateType = Array<CheckListItem>;
