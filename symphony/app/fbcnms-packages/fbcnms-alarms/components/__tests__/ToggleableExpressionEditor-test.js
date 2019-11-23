/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
import 'jest-dom/extend-expect';
import {cleanup} from '@testing-library/react';
import {thresholdToPromQL} from '../ToggleableExpressionEditor';

import type {ThresholdExpression} from '../ToggleableExpressionEditor';

afterEach(() => {
  cleanup();
  jest.resetAllMocks();
});

type ToPromQLTestCase = {
  expression: ThresholdExpression,
  expectedPromQL: string,
};

test('correctly converts a ThresholdExpression to PromQL', () => {
  const testCases: Array<ToPromQLTestCase> = [
    {
      expression: {metricName: 'test', comparator: '<', filters: [], value: 7},
      expectedPromQL: 'test<7',
    },
    {
      expression: {
        metricName: 'test',
        comparator: '>',
        filters: [
          {name: 'label1', value: 'val1'},
          {name: 'label2', value: 'val2'},
        ],
        value: 10,
      },
      expectedPromQL: 'test{label1=~"^val1$",label2=~"^val2$",}>10',
    },
  ];

  testCases.forEach(test => {
    expect(thresholdToPromQL(test.expression)).toBe(test.expectedPromQL);
  });
});
