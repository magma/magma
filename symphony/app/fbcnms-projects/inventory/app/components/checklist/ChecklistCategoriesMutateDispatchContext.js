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

import React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';

type Dispatch<A> = A => void;

type ChecklistCategoriesMutateDispatchContextDispatcher = Dispatch<ChecklistCategoriesMutateStateActionType>;

export default (React.createContext<ChecklistCategoriesMutateDispatchContextDispatcher>(
  emptyFunction,
): React$Context<ChecklistCategoriesMutateDispatchContextDispatcher>);
