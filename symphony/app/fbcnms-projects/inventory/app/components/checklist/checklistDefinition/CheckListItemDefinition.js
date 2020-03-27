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

import CheckListItemCollapsedDefinition from './CheckListItemCollapsedDefinition';
import React from 'react';
import {CheckListItemConfigs} from '../checkListCategory/CheckListItemConsts';
import {getValidChecklistItemType} from '../ChecklistUtils';

export type CheckListItemDefinitionProps = {
  item: CheckListItem,
  editedDefinitionId: ?string,
  onChange?: (updatedChecklistItemDefinition: CheckListItem) => void,
};

const CheckListItemDefinition = (props: CheckListItemDefinitionProps) => {
  const {item, editedDefinitionId} = props;

  const itemTypeKey = item && getValidChecklistItemType(item.type);
  const itemType = itemTypeKey && CheckListItemConfigs[itemTypeKey];
  const CheckListItemDefinitionComponent =
    itemType && itemType.definitionComponent;
  if (!CheckListItemDefinitionComponent) {
    return null;
  }

  if (item.id !== editedDefinitionId) {
    return <CheckListItemCollapsedDefinition item={item} />;
  }

  return <CheckListItemDefinitionComponent {...props} />;
};

export default CheckListItemDefinition;
