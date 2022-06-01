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

import * as React from 'react';
import type {ApiRequest} from '../AlarmsApi';
import type {FiringAlarm} from '../AlarmAPIType';

export type RuleEditorProps<TRule> = {
  rule: GenericRule<TRule> | undefined | null;
  // invoked when rule is modified locally
  onRuleUpdated: (rule: GenericRule<TRule>) => void;
  onExit: () => void;
  isNew: boolean;
};

export type RuleViewerProps = {row?: any};

export type AlertViewerProps = {
  alert: FiringAlarm;
};

/**
 * Rules should be mapped to the generic rule for rendering in tables and
 * passing to shared components. The raw rule should be passed to
 * ruletype-specific components such as PrometheusEditor.
 */
export type GenericRule<TRule> = {
  severity: string;
  name: string;
  description: string;
  period: string;
  expression: string;
  /**
   * Type of rule, used for selecting which rule editor to use. Must exist as
   * a key inside the RuleInterfaceMap passed to the Alarms component
   */
  ruleType: string;
  // The original rule which this generic rule is mapped from
  rawRule: TRule;
};

export type RuleInterface<TRule> = {
  friendlyName: string;
  /**
   * Component to create and edit an alerting rule.
   */
  RuleEditor: React.ComponentType<RuleEditorProps<TRule>>;
  /**
   * Component to be rendered inside of the "view alert" modal. Use this to
   * display information specific to the rule type. If this is not provided,
   * the alert will be rendered as json in a pre tag
   */
  AlertViewer?: React.ComponentType<AlertViewerProps>;
  RuleViewer?: React.ComponentType<RuleViewerProps>;
  /**
   * Retrieve all rules for this rule type
   */
  getRules: (req: ApiRequest) => Promise<Array<GenericRule<TRule>>>;
  deleteRule: (req: {ruleName: string} & ApiRequest) => Promise<void>;
};

export type RuleInterfaceMap<TUnion> = Record<string, RuleInterface<TUnion>>;
