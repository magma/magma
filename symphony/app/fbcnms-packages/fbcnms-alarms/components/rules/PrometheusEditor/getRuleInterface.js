/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import PrometheusEditor from './PrometheusEditor';

import type {AlertConfig} from '../../AlarmAPIType';
import type {ApiUtil} from '../../AlarmsApi';
import type {GenericRule, RuleInterfaceMap} from '../../rules/RuleInterface';

const PROMETHEUS_RULE_TYPE = 'prometheus';
export default function getPrometheusRuleInterface({
  apiUtil,
}: {
  apiUtil: ApiUtil,
}): RuleInterfaceMap<AlertConfig> {
  return {
    [PROMETHEUS_RULE_TYPE]: {
      friendlyName: PROMETHEUS_RULE_TYPE,
      RuleEditor: PrometheusEditor,
      /**
       * Get alert rules from backend and map to generic
       */
      getRules: async req => {
        const rules = await apiUtil.getAlertRules(req);
        return rules.map<GenericRule<AlertConfig>>(rule => ({
          name: rule.alert,
          description: rule.annotations?.description || '',
          severity: rule.labels?.severity || '',
          period: rule.for || '',
          expression: rule.expr,
          ruleType: PROMETHEUS_RULE_TYPE,
          rawRule: rule,
        }));
      },
      deleteRule: params => apiUtil.deleteAlertRule(params),
    },
  };
}
