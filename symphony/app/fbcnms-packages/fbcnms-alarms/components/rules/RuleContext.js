/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import type {RuleInterfaceMap} from './RuleInterface';

export type RuleContext<TRuleUnion> = {
  ruleMap: RuleInterfaceMap<TRuleUnion>,
  ruleType: string,
  selectRuleType: string => void,
};

const context = React.createContext<RuleContext<*>>({});

export default context;
