/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ChecklistItemsDialogMutateStateActionType} from './ChecklistItemsDialogMutateAction';

import React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';

type Dispatch<A> = A => void;

type ChecklistItemsDialogMutateDispatchContextDispatcher = Dispatch<ChecklistItemsDialogMutateStateActionType>;

export default (React.createContext<ChecklistItemsDialogMutateDispatchContextDispatcher>(
  emptyFunction,
): React$Context<ChecklistItemsDialogMutateDispatchContextDispatcher>);
