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

import React from 'react';
import type {ApiUtil} from './AlarmsApi';
import type {FiringAlarm, Labels} from './AlarmAPIType';
import type {RuleInterfaceMap} from './rules/RuleInterface';

export const PROMETHEUS_RULE_TYPE = 'prometheus';

export type AlarmContext = {
  apiUtil: ApiUtil;
  filterLabels?: (labels: Labels) => Labels;
  ruleMap: RuleInterfaceMap<any>;
  getAlertType?: GetAlertType | undefined | null;
  getNetworkId?: () => string;
  // feature flags
  thresholdEditorEnabled?: boolean;
  alertManagerGlobalConfigEnabled?: boolean;
};

/***
 * Determine the type of alert based on its labels/annotations. Since all
 * alerts come from alertmanager, regardless of source, we can only determine
 * the source by inspecting the labels/annotations.
 */
export type GetAlertType = (
  alert: FiringAlarm,
  ruleMap?: RuleInterfaceMap<any>,
) => keyof RuleInterfaceMap<any>;

const emptyApiUtil = {
  useAlarmsApi: () => ({
    response: null,
    error: new Error('not implemented'),
    isLoading: false,
  }),
  viewFiringAlerts: () => Promise.reject('not implemented'),
  getTroubleshootingLink: () => Promise.reject('not implemented'),
  viewMatchingAlerts: () => Promise.reject('not implemented'),
  createAlertRule: () => Promise.reject('not implemented'),
  editAlertRule: () => Promise.reject('not implemented'),
  getAlertRules: () => Promise.reject('not implemented'),
  deleteAlertRule: () => Promise.reject('not implemented'),
  getSuppressions: () => Promise.reject('not implemented'),
  createReceiver: () => Promise.reject('not implemented'),
  editReceiver: () => Promise.reject('not implemented'),
  getReceivers: () => Promise.reject('not implemented'),
  deleteReceiver: () => Promise.reject('not implemented'),
  getRouteTree: () => Promise.reject('not implemented'),
  editRouteTree: () => Promise.reject('not implemented'),
  getMetricNames: () => Promise.reject('not implemented'),
  getMetricSeries: () => Promise.reject('not implemented'),
  getGlobalConfig: () => Promise.reject('not implemented'),
  editGlobalConfig: () => Promise.reject('not implemented'),
  getTenants: () => Promise.reject('not implemented'),
  getAlertmanagerTenancy: () => Promise.reject('not implemented'),
  getPrometheusTenancy: () => Promise.reject('not implemented'),
};

const context = React.createContext<AlarmContext>({
  apiUtil: emptyApiUtil,
  filterLabels: x => x,
  ruleMap: {},
  getAlertType: () => PROMETHEUS_RULE_TYPE,
  thresholdEditorEnabled: false,
  alertManagerGlobalConfigEnabled: false,
});

export function useAlarmContext() {
  return React.useContext(context);
}

export default context;
