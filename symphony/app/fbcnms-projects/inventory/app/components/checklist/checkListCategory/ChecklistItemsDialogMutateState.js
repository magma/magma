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
  CheckListItemEnumSelectionMode,
  CheckListItemType,
  WorkOrderDetails_workOrder,
  YesNoResponse,
} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';

export type CheckListItemDefinition = $ReadOnly<{|
  id: string,
  title: string,
  type: CheckListItemType,
  index?: ?number,
  enumValues?: ?string,
  enumSelectionMode?: ?CheckListItemEnumSelectionMode,
  helpText?: ?string,
|}>;

export type CheckListItemBase = $ReadOnly<{|
  id: string,
  index?: ?number,
  type: CheckListItemType,
  title: string,
  helpText?: ?string,
|}>;

export type BasicCheckListItemData = $ReadOnly<{|
  checked?: ?boolean,
|}>;

export type EnumCheckListItemData = $ReadOnly<{|
  enumValues?: ?string,
  selectedEnumValues?: ?string,
  enumSelectionMode?: ?CheckListItemEnumSelectionMode,
|}>;

export type FreeTextCheckListItemData = $ReadOnly<{|
  stringValue?: ?string,
|}>;

export type CheckListItemFile = $ReadOnly<{|
  id: string,
  storeKey: string,
  fileName: string,
  sizeInBytes?: number,
  modificationTime?: number,
  uploadTime?: number,
  annotation?: ?string,
|}>;

export type CheckListItemPendingFile = $ReadOnly<{|
  id: string,
  name: string,
  progress: number,
|}>;

export type FilesCheckListItemData = $ReadOnly<{|
  files?: ?Array<CheckListItemFile>,
  pendingFiles?: ?Array<CheckListItemPendingFile>,
|}>;

export type YesNoCheckListItemData = $ReadOnly<{|
  yesNoResponse?: ?YesNoResponse,
|}>;

export type CellScanCheckListItemData = {|
  +cellData?: ?$ElementType<
    $ElementType<
      $ElementType<
        $ElementType<
          $ElementType<WorkOrderDetails_workOrder, 'checkListCategories'>,
          number,
        >,
        'checkList',
      >,
      number,
    >,
    'cellData',
  >,
|};

export type WifiScanCheckListItemData = {|
  +wifiData?: ?$ElementType<
    $ElementType<
      $ElementType<
        $ElementType<
          $ElementType<WorkOrderDetails_workOrder, 'checkListCategories'>,
          number,
        >,
        'checkList',
      >,
      number,
    >,
    'wifiData',
  >,
|};

export type CheckListItem = {|
  ...CheckListItemBase,
  ...BasicCheckListItemData,
  ...EnumCheckListItemData,
  ...FreeTextCheckListItemData,
  ...FilesCheckListItemData,
  ...YesNoCheckListItemData,
  ...CellScanCheckListItemData,
  ...WifiScanCheckListItemData,
|};

export type ChecklistItemsDialogStateType = $ReadOnly<{
  items: Array<CheckListItem>,
  editedDefinitionId: ?string,
}>;
