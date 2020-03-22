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

import BasicCheckListItemDefinition from './BasicCheckListItemDefinition';
import CheckListItemCollapsedDefinition from './CheckListItemCollapsedDefinition';
import FreeTextCheckListItemDefinition from './FreeTextCheckListItemDefinition';
import React from 'react';
import fbt from 'fbt';

export const CHECKLIST_ITEM_DEFINITION_TYPES = {
  simple: {
    description: fbt(
      'Check when complete',
      'Description of a simple checklist item (`mark when done` like)',
    ),
    component: BasicCheckListItemDefinition,
  },
  string: {
    description: fbt(
      'Free text',
      'Description of a free text checklist item (e.g. `enter details here`)',
    ),
    component: FreeTextCheckListItemDefinition,
  },
};

export const GetValidChecklistItemType = (
  type: string,
): 'simple' | 'string' | null => {
  if (type === 'simple' || type === 'string') {
    return type;
  }

  return null;
};

type Props = {
  item: CheckListItem,
  editedDefinitionId: ?string,
  onChange?: (updatedChecklistItemDefinition: CheckListItem) => void,
};

const CheckListItemDefinition = (props: Props) => {
  const {item, editedDefinitionId} = props;

  const itemTypeKey = item && GetValidChecklistItemType(item.type);
  const itemType = itemTypeKey && CHECKLIST_ITEM_DEFINITION_TYPES[itemTypeKey];
  const CheckListItemDefinitionComponent = itemType && itemType.component;
  if (!CheckListItemDefinitionComponent) {
    return null;
  }

  if (item.id !== editedDefinitionId) {
    return <CheckListItemCollapsedDefinition item={item} />;
  }

  return <CheckListItemDefinitionComponent {...props} />;
};

export default CheckListItemDefinition;
