/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import BasicCheckListItemDefinition from './checklistDefinition/BasicCheckListItemDefinition';
import BasicCheckListItemFilling from './checklistFilling/BasicCheckListItemFilling';
import FreeTextCheckListItemDefinition from './checklistDefinition/FreeTextCheckListItemDefinition';
import FreeTextCheckListItemFilling from './checklistFilling/FreeTextCheckListItemFilling';
import React from 'react';
import fbt from 'fbt';
import {createFragmentContainer, graphql} from 'react-relay';
import type {CheckListItem_item} from './__generated__/CheckListItem_item.graphql';

export const CHECKLIST_ITEM_TYPES = {
  simple: {
    description: fbt(
      'Check when complete',
      'Description of a simple checklist item (`mark when done` like)',
    ),
    component: {
      design: BasicCheckListItemDefinition,
      filling: BasicCheckListItemFilling,
    },
  },
  string: {
    description: fbt(
      'Free text',
      'Description of a free text checklist item (e.g. `enter details here`)',
    ),
    component: {
      design: FreeTextCheckListItemDefinition,
      filling: FreeTextCheckListItemFilling,
    },
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
  item: CheckListItem_item,
  designMode?: boolean,
  onChange?: (updatedChecklistItem: CheckListItem_item) => void,
};

const CheckListItem = (props: Props) => {
  const {item, designMode} = props;

  const itemTypeKey = item && GetValidChecklistItemType(item.type);
  const itemType = itemTypeKey && CHECKLIST_ITEM_TYPES[itemTypeKey];
  const itemComponents = itemType && itemType.component;
  const CheckListItemComponent =
    itemComponents &&
    (designMode ? itemComponents.design : itemComponents.filling);
  if (!CheckListItemComponent) {
    return null;
  }

  const checkListItemComponentProps = {
    ...props,
    checkListItem: props.item,
  };

  return <CheckListItemComponent {...checkListItemComponentProps} />;
};

export default createFragmentContainer(CheckListItem, {
  item: graphql`
    fragment CheckListItem_item on CheckListItem {
      id
      title
      type
      index
      helpText
      enumValues
      stringValue
      checked
    }
  `,
});
