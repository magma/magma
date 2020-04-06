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

import CheckListItemDefinitionBase from './CheckListItemDefinitionBase';
import React from 'react';

type Props = {
  item: CheckListItem,
  onChange?: (updatedItem: CheckListItem) => void,
};

const FreeTextCheckListItemDefinition = (props: Props) => {
  return <CheckListItemDefinitionBase {...props} />;
};

export default FreeTextCheckListItemDefinition;
