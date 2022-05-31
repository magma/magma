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
 */

import React from 'react';
import type {
  AlertConfig,
  FiringAlarm,
  Labels,
  // $FlowFixMe migrated to typescript
} from '../components/AlarmAPIType';
import type {
  AlertViewerProps,
  GenericRule,
  RuleEditorProps,
  RuleInterface,
  RuleViewerProps,
} from '../components/rules/RuleInterface';

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

export function mockAlert(merge?: $Shape<FiringAlarm>): FiringAlarm {
  const {labels, annotations, ...otherFields} = merge || {};
  const defaultLabels: Labels = {alertname: 'test', severity: 'NOTICE'};
  const defaultAnnotations: Labels = {description: 'test description'};
  return {
    startsAt: '2020-02-10T21:09:12Z',
    endsAt: '',
    fingerprint: '',
    receivers: [],
    status: {inhibitedBy: [], silencedBy: [], state: ''},
    labels: {
      ...defaultLabels,
      ...(labels || {}),
    },
    annotations: {
      ...defaultAnnotations,
      ...(annotations || {}),
    },
    ...(otherFields: $Shape<
      $Rest<FiringAlarm, {|labels: Labels, annotations: Labels|}>,
    >),
  };
}

export function mockRuleInterface<TRule>(
  overrides?: $Shape<RuleInterface<TRule>>,
): RuleInterface<TRule> {
  const {
    friendlyName,
    AlertViewer,
    RuleEditor,
    RuleViewer,
    deleteRule,
    getRules,
  } = overrides || {};
  return {
    friendlyName: friendlyName ?? 'mock rule',
    RuleEditor:
      RuleEditor ??
      function (_props: RuleEditorProps<TRule>) {
        return <span />;
      },
    RuleViewer:
      RuleViewer ??
      function (_props: RuleViewerProps<TRule>) {
        return <span />;
      },
    AlertViewer:
      AlertViewer ??
      function (_props: AlertViewerProps) {
        return <span />;
      },
    deleteRule: deleteRule ?? jest.fn(() => Promise.resolve()),
    getRules: getRules ?? jest.fn(() => Promise.resolve([])),
  };
}
