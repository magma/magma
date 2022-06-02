/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import PrometheusEditor from './PrometheusEditor';

import {PROMETHEUS_RULE_TYPE} from '../../AlarmContext';
import type {AlertConfig} from '../../AlarmAPIType';
import type {ApiUtil} from '../../AlarmsApi';
import type {GenericRule, RuleInterfaceMap} from '../RuleInterface';

export default function getPrometheusRuleInterface({
  apiUtil,
}: {
  apiUtil: ApiUtil;
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
