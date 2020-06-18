/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {ActionID} from '../../mutations/__generated__/AddActionsRuleMutation.graphql.js';

export type FilterData = string | string[];

export type ActionData = string | string[] | null;

// A user-configured TriggerFilter when creating a Rule
export type RuleFilter = {|
  filterID: string,
  operatorID: string,
  data: FilterData,
|};

// A user-configured action when creating a Rule
export type RuleAction = {|
  actionID: ActionID,
  data: ActionData,
|};
