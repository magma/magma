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
 *
 * @flow strict-local
 * @format
 *
 */

import React from 'react';
import {PROMETHEUS_RULE_TYPE} from './rules/PrometheusEditor/getRuleInterface';
// $FlowFixMe migrated to typescript
import type {ApiUtil} from './AlarmsApi';
// $FlowFixMe migrated to typescript
import type {FiringAlarm, Labels} from './AlarmAPIType';
import type {RuleInterfaceMap} from './rules/RuleInterface';

export type AlarmContext = {|
  apiUtil: ApiUtil,
  filterLabels?: (labels: Labels) => Labels,
  ruleMap: RuleInterfaceMap<*>,
  getAlertType?: ?GetAlertType,
  getNetworkId?: () => string,
  // feature flags
  thresholdEditorEnabled?: boolean,
  alertManagerGlobalConfigEnabled?: boolean,
|};

/***
 * Determine the type of alert based on its labels/annotations. Since all
 * alerts come from alertmanager, regardless of source, we can only determine
 * the source by inspecting the labels/annotations.
 */
export type GetAlertType = (
  alert: FiringAlarm,
  ruleMap?: RuleInterfaceMap<mixed>,
) => $Keys<RuleInterfaceMap<mixed>>;

const emptyApiUtil = {
  useAlarmsApi: () => ({
    response: null,
    error: new Error('not implemented'),
    isLoading: false,
  }),
  viewFiringAlerts: (..._) => Promise.reject('not implemented'),
  getTroubleshootingLink: (..._) => Promise.reject('not implemented'),
  viewMatchingAlerts: (..._) => Promise.reject('not implemented'),
  createAlertRule: (..._) => Promise.reject('not implemented'),
  editAlertRule: (..._) => Promise.reject('not implemented'),
  getAlertRules: (..._) => Promise.reject('not implemented'),
  deleteAlertRule: (..._) => Promise.reject('not implemented'),
  getSuppressions: (..._) => Promise.reject('not implemented'),
  createReceiver: (..._) => Promise.reject('not implemented'),
  editReceiver: (..._) => Promise.reject('not implemented'),
  getReceivers: (..._) => Promise.reject('not implemented'),
  deleteReceiver: (..._) => Promise.reject('not implemented'),
  getRouteTree: (..._) => Promise.reject('not implemented'),
  editRouteTree: (..._) => Promise.reject('not implemented'),
  getMetricNames: (..._) => Promise.reject('not implemented'),
  getMetricSeries: (..._) => Promise.reject('not implemented'),
  getGlobalConfig: _ => Promise.reject('not implemented'),
  editGlobalConfig: _ => Promise.reject('not implemented'),
  getTenants: _ => Promise.reject('not implemented'),
  getAlertmanagerTenancy: _ => Promise.reject('not implemented'),
  getPrometheusTenancy: _ => Promise.reject('not implemented'),
};

const context = React.createContext<AlarmContext>({
  apiUtil: emptyApiUtil,
  filterLabels: x => x,
  ruleMap: {},
  getAlertType: _ => PROMETHEUS_RULE_TYPE,
  thresholdEditorEnabled: false,
  alertManagerGlobalConfigEnabled: false,
});

export function useAlarmContext() {
  return React.useContext(context);
}

export default context;
