/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItem} from '../checkListCategory/ChecklistItemsDialogMutateState';

import BasicCheckListItemFilling from './BasicCheckListItemFilling';
import FreeTextCheckListItemFilling from './FreeTextCheckListItemFilling';
import MultipleChoiceCheckListItemFilling from './MultipleChoiceCheckListItemFilling';
import React from 'react';

export const CHECKLIST_ITEM_FILLING_TYPES = {
  simple: {
    component: BasicCheckListItemFilling,
  },
  string: {
    component: FreeTextCheckListItemFilling,
  },
  enum: {
    component: MultipleChoiceCheckListItemFilling,
  },
};

export const GetValidChecklistItemType = (
  type: string,
): 'simple' | 'string' | 'enum' | null => {
  if (type === 'simple' || type === 'string' || type === 'enum') {
    return type;
  }

  return null;
};

type Props = {
  item: CheckListItem,
  onChange?: (updatedChecklistItemFilling: CheckListItem) => void,
};

const CheckListItemFilling = (props: Props) => {
  const {item} = props;

  const itemTypeKey = item && GetValidChecklistItemType(item.type);
  const itemType = itemTypeKey && CHECKLIST_ITEM_FILLING_TYPES[itemTypeKey];
  const CheckListItemFillingComponent = itemType && itemType.component;
  if (!CheckListItemFillingComponent) {
    return null;
  }

  const checkListItemFillingComponentProps = {
    ...props,
    checkListItem: props.item,
  };

  return (
    <CheckListItemFillingComponent {...checkListItemFillingComponentProps} />
  );
};

export default CheckListItemFilling;
