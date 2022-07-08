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
import {PromFiringAlert} from '../../../../generated-ts';
import type {AlertConfig, Labels} from '../components/AlarmAPIType';
import type {
  GenericRule,
  RuleInterface,
} from '../components/rules/RuleInterface';

export function mockPrometheusRule(merge?: Partial<GenericRule<AlertConfig>>) {
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

export function mockAlert(merge?: Partial<PromFiringAlert>): PromFiringAlert {
  const {labels, annotations, ...otherFields} = merge || {};
  const defaultLabels: Labels = {
    alertname: 'test',
    severity: 'NOTICE',
  };
  const defaultAnnotations: Labels = {
    description: 'test description',
  };
  return {
    annotations: {...defaultAnnotations, ...(annotations || {})},
    endsAt: '',
    fingerprint: '',
    labels: {...defaultLabels, ...(labels || {})},
    receivers: {name: 'foo'},
    startsAt: '2020-02-10T21:09:12Z',
    status: {
      inhibitedBy: [],
      silencedBy: [],
      state: '',
    },
    updatedAt: '',
    ...(otherFields as Partial<
      Omit<PromFiringAlert, 'labels' | 'annotations'>
    >),
  };
}

export function mockRuleInterface<TRule>(
  overrides?: Partial<RuleInterface<TRule>>,
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
      function () {
        return <span />;
      },
    RuleViewer:
      RuleViewer ??
      function () {
        return <span />;
      },
    AlertViewer:
      AlertViewer ??
      function () {
        return <span />;
      },
    deleteRule: deleteRule ?? jest.fn(() => Promise.resolve({data: undefined})),
    getRules: getRules ?? jest.fn(() => Promise.resolve([])),
  };
}
