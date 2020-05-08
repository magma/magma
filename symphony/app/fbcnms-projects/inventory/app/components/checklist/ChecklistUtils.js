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
  CheckListCategoryInput,
  CheckListItemInput,
} from '../../mutations/__generated__/AddWorkOrderMutation.graphql';
import type {CheckListItem} from './checkListCategory/ChecklistItemsDialogMutateState';
import type {CheckListItemType} from '../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';
import type {ChecklistCategoriesStateType} from './ChecklistCategoriesMutateState';

import {isTempId} from '../../common/EntUtils';

export const getValidChecklistItemType = (
  type: CheckListItemType,
): CheckListItemType => {
  switch (type) {
    case 'simple':
      return 'simple';
    case 'string':
      return 'string';
    case 'enum':
      return 'enum';
    case 'files':
      return 'files';
    case 'yes_no':
      return 'yes_no';
    case 'cell_scan':
      return 'cell_scan';
    case 'wifi_scan':
      return 'wifi_scan';
    default:
      throw new Error(
        `Invariant violation - checklist item type not found: ${type}`,
      );
  }
};

export const isChecklistItemDone = (item: CheckListItem): boolean => {
  switch (item.type) {
    case 'enum':
      return item.enumValues != null && item.enumValues.trim().length > 0;
    case 'simple':
      return item.checked === true;
    case 'string':
      return item.stringValue != null && item.stringValue.trim() !== '';
    case 'files':
      return item.files != null && item.files.length > 0;
    case 'yes_no':
      return item.yesNoResponse != null;
    case 'cell_scan':
      return item.cellData != null;
    case 'wifi_scan':
      return item.wifiData != null;
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

export const convertChecklistCategoriesStateToInput = (
  checkListCategories: ChecklistCategoriesStateType,
): Array<CheckListCategoryInput> => {
  return checkListCategories.map(category => {
    const checkList: Array<CheckListItemInput> = category.checkList.map(
      item => {
        return {
          id: isTempId(item.id) ? undefined : item.id,
          title: item.title,
          type: item.type,
          index: item.index,
          helpText: item.helpText,
          enumValues: item.enumValues,
          selectedEnumValues: item.selectedEnumValues,
          enumSelectionMode: item.enumSelectionMode,
          stringValue: item.stringValue,
          checked: item.checked,
          yesNoResponse: item.yesNoResponse,
          files: item.files?.map(file => ({
            id: isTempId(file.id) ? undefined : file.id,
            storeKey: file.storeKey,
            fileName: file.fileName,
            sizeInBytes: file.sizeInBytes,
            modificationTime: file.modificationTime,
            uploadTime: file.uploadTime,
            fileType: 'FILE',
          })),
        };
      },
    );
    return {
      id: isTempId(category.id) ? undefined : category.id,
      title: category.title,
      description: category.description,
      checkList,
    };
  });
};
