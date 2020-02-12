/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {AlertConfig} from '../components/AlarmAPIType';
import type {GenericRule} from '../components/rules/RuleInterface';

export function mockPrometheusRule(merge?: $Shape<GenericRule<AlertConfig>>) {
  return {
    name: '<<test>>',
    severity: 'info',
    description: '<<test description>>',
    expression: 'up == 0',
    period: '1m',
    ruleType: 'prometheus',
    rawRule: {
      alert: '<<test>>',
      labels: {
        severity: 'info',
      },
      expr: 'up == 0',
      for: '1m',
    },
    ...(merge || {}),
  };
}
