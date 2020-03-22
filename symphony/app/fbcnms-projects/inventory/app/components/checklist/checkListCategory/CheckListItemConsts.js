/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {CheckListItemType} from '../../work_orders/__generated__/WorkOrderDetails_workOrder.graphql';
import type {SvgIconStyleProps} from '@fbcnms/ui/components/design-system/Icons/SvgIcon';

import * as React from 'react';
import {
  ChecklistCheckIcon,
  TextIcon,
} from '@fbcnms/ui/components/design-system/Icons';

export const CheckListItemIcons: {
  [CheckListItemType]: React.ComponentType<SvgIconStyleProps>,
} = {
  simple: ChecklistCheckIcon,
  string: TextIcon,
};
