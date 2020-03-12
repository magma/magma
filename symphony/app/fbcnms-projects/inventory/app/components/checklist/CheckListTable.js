/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ChecklistItemsDialogStateType} from './checkListCategory/ChecklistItemsDialogMutateState';

import CheckListTableDefinition from './checklistDefinition/CheckListTableDefinition';
import CheckListTableFilling from './checklistFilling/CheckListTableFilling';
import React from 'react';
import {CHECKLIST_ITEM_TYPES} from './CheckListItem';
import {sortByIndex} from '../draggable/DraggableUtils';

type Props = {
  items: ChecklistItemsDialogStateType,
  onDesignMode?: boolean,
};

const CheckListTable = (props: Props) => {
  const checkListTableItems = Array.prototype.filter
    .call(props.items || [], item =>
      CHECKLIST_ITEM_TYPES.hasOwnProperty(item.type),
    )
    .sort(sortByIndex);

  const CheckListTableComponent = props.onDesignMode
    ? CheckListTableDefinition
    : CheckListTableFilling;

  const checkListTableProps = {
    ...props,
    items: checkListTableItems,
  };

  return <CheckListTableComponent {...checkListTableProps} />;
};

export default CheckListTable;
