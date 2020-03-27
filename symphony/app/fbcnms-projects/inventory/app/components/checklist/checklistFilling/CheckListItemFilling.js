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

import React from 'react';
import {CheckListItemConfigs} from '../checkListCategory/CheckListItemConsts';
import {getValidChecklistItemType} from '../ChecklistUtils';

export type CheckListItemFillingProps = {
  item: CheckListItem,
  onChange?: (updatedChecklistItemFilling: CheckListItem) => void,
};

const CheckListItemFilling = (props: CheckListItemFillingProps) => {
  const {item} = props;

  const itemTypeKey = item && getValidChecklistItemType(item.type);
  const itemType = itemTypeKey && CheckListItemConfigs[itemTypeKey];
  const CheckListItemFillingComponent = itemType && itemType.fillingComponent;
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
