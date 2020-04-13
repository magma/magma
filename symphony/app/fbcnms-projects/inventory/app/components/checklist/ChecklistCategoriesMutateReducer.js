/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ChecklistCategoriesMutateStateActionType} from './ChecklistCategoriesMutateAction';
import type {
  ChecklistCategoriesStateType,
  ChecklistCategory,
} from './ChecklistCategoriesMutateState';
import type {WorkOrderDetails_workOrder} from '../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';

import fbt from 'fbt';
import shortid from 'shortid';

export function getInitialState(
  categories: $ElementType<WorkOrderDetails_workOrder, 'checkListCategories'>,
): ChecklistCategoriesStateType {
  return categories.slice().map(category => ({
    id: category.id ?? shortid.generate(),
    title: category.title,
    description: category.description,
    checkList: category.checkList.slice().map(item => ({
      id: item.id,
      index: item.index,
      type: item.type,
      title: item.title,
      helpText: item.helpText,
      checked: item.checked,
      enumValues: item.enumValues,
      enumSelectionMode: !!item.enumSelectionMode
        ? item.enumSelectionMode
        : 'single',
      selectedEnumValues: item.selectedEnumValues,
      stringValue: item.stringValue,
      yesNoResponse: item.yesNoResponse,
      files: item.files.map(file => ({
        id: file.id,
        storeKey: file.storeKey ?? '',
        fileName: file.fileName,
      })),
      cellData: item.cellData,
      wifiData: item.wifiData,
    })),
  }));
}

function createNewCategory(): ChecklistCategory {
  return {
    id: shortid.generate(),
    title: `${fbt('New Category', 'Default name for checklist category')}`,
    description: '',
    checkList: [],
  };
}

function updateCategory<T: ChecklistCategoriesStateType>(
  state: T,
  categoryId: string,
  change,
): T {
  const categoryIndex = state.findIndex(c => c.id === categoryId);
  return [
    ...state.slice(0, categoryIndex),
    {
      ...state[categoryIndex],
      ...change,
    },
    ...state.slice(categoryIndex + 1),
  ];
}

export function reducer<T: ChecklistCategoriesStateType>(
  state: T,
  action: ChecklistCategoriesMutateStateActionType,
): T {
  switch (action.type) {
    case 'UPDATE_CATEGORY_TITLE':
      return updateCategory(state, action.categoryId, {
        title: action.value,
      });
    case 'UPDATE_CATEGORY_DESCRIPTION':
      return updateCategory(state, action.categoryId, {
        description: action.value,
      });
    case 'UPDATE_CATEGORY_CHECKLIST':
      return updateCategory(state, action.categoryId, {
        checkList: action.value,
      });
    case 'ADD_CATEGORY':
      return [...state, createNewCategory()];
    case 'REMOVE_CATEGORY':
      const categoryIndex = state.findIndex(c => c.id === action.categoryId);
      return [
        ...state.slice(0, categoryIndex),
        ...state.slice(categoryIndex + 1, state.length),
      ];
    default:
      return state;
  }
}
