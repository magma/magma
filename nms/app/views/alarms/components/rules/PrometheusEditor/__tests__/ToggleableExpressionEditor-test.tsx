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

import * as PromQL from '../../../prometheus/PromQL';

import {thresholdToPromQL} from '../ToggleableExpressionEditor';

import type {ThresholdExpression} from '../ToggleableExpressionEditor';

type ToPromQLTestCase = {
  expression: ThresholdExpression;
  expectedPromQL: string;
};

test('correctly converts a ThresholdExpression to PromQL', () => {
  const testCases: Array<ToPromQLTestCase> = [
    {
      expression: {
        metricName: 'test',
        comparator: new PromQL.BinaryComparator('<'),
        filters: new PromQL.Labels(),
        value: 7,
      },
      expectedPromQL: 'test < 7',
    },
    {
      expression: {
        metricName: 'test',
        comparator: new PromQL.BinaryComparator('>'),
        filters: new PromQL.Labels()
          .addEqual('label1', 'val1')
          .addEqual('label2', 'val2'),
        value: 10,
      },
      expectedPromQL: 'test{label1="val1",label2="val2"} > 10',
    },
    {
      expression: {
        metricName: 'test',
        comparator: new PromQL.BinaryComparator('>'),
        filters: new PromQL.Labels()
          .addRegex('label1', 'val1')
          .addRegex('label2', 'val2'),
        value: 10,
      },
      expectedPromQL: 'test{label1=~"val1",label2=~"val2"} > 10',
    },
  ];

  testCases.forEach(test => {
    expect(thresholdToPromQL(test.expression)).toBe(test.expectedPromQL);
  });
});
