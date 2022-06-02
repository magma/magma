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
 * @flow
 * @format
 */

import * as React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import PrometheusEditor from '../PrometheusEditor';
import {alarmTestUtil} from '../../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {parseTimeString} from '../PrometheusEditor';
import {render} from '@testing-library/react';

import type {AlarmsWrapperProps} from '../../../../test/testHelpers';
// $FlowFixMe migrated to typescript
import type {AlertConfig} from '../../../AlarmAPIType';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {ApiUtil} from '../../../AlarmsApi';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {GenericRule} from '../../RuleInterface';

// TextField select is difficult to test so replace it with an Input
jest.mock('@material-ui/core/TextField', () => {
  const Input = require('@material-ui/core/Input').default;
  return ({children: _, InputProps: __, label, ...props}) => (
    <label>
      {label}
      <Input {...props} />
    </label>
  );
});

describe('PrometheusEditor', () => {
  let AlarmsWrapper: React.ComponentType<$Shape<AlarmsWrapperProps>>;
  let apiUtil: ApiUtil;

  beforeEach(() => {
    ({apiUtil, AlarmsWrapper} = alarmTestUtil());
  });

  const commonProps = {
    onRuleUpdated: () => {},
    onExit: () => {},
    isNew: false,
    onRuleSaved: jest.fn(),
  };

  test('editing a threshold alert opens the PrometheusEditor with the threshold expression editor enabled', async () => {
    jest.spyOn(apiUtil, 'getMetricSeries').mockResolvedValue([]);
    const testThresholdRule: GenericRule<AlertConfig> = {
      severity: '',
      ruleType: '',
      rawRule: {alert: '', expr: 'metric > 123'},
      period: '',
      name: '',
      description: '',
      expression: 'metric > 123',
    };
    const {getByDisplayValue} = render(
      <AlarmsWrapper thresholdEditorEnabled={true}>
        <PrometheusEditor {...commonProps} rule={testThresholdRule} />
      </AlarmsWrapper>,
    );
    expect(getByDisplayValue('metric')).toBeInTheDocument();
    expect(getByDisplayValue('123')).toBeInTheDocument();
  });

  test('editing a non-threshold alert opens the PrometheusEditor with the advanced editor enabled', async () => {
    const testThresholdRule: GenericRule<AlertConfig> = {
      severity: '',
      ruleType: '',
      rawRule: {alert: '', expr: 'vector(1)'},
      period: '',
      name: '',
      description: '',
      expression: 'vector(1)',
    };
    const {getByDisplayValue} = render(
      <AlarmsWrapper thresholdEditorEnabled={true}>
        <PrometheusEditor {...commonProps} rule={testThresholdRule} />
      </AlarmsWrapper>,
    );
    expect(getByDisplayValue('vector(1)')).toBeInTheDocument();
  });

  describe('Duration Parser', () => {
    const testCases = [
      ['empty input', '', {hours: 0, minutes: 0, seconds: 0}],
      ['out of order units', '1s2m3h', {hours: 0, minutes: 0, seconds: 0}],
      ['all units', '1h2m3s', {hours: 1, minutes: 2, seconds: 3}],
      ['hour', '1h', {hours: 1, minutes: 0, seconds: 0}],
      ['minute', '1m', {hours: 0, minutes: 1, seconds: 0}],
      ['second', '1s', {hours: 0, minutes: 0, seconds: 1}],
    ];
    test.each(testCases)('%s', (name, input, expectedDuration) => {
      expect(parseTimeString(input)).toEqual(expectedDuration);
    });
  });
});
