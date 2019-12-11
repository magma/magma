/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import {Tokenize} from '../PromQLTokenizer';

describe('Tokenize', () => {
  const testCases = [
    ['single metric selector', 'metric', [{value: 'metric', type: 'word'}]],
    [
      'whitespace',
      'metric and\tmetric',
      [
        {value: 'metric', type: 'word'},
        {value: 'and', type: 'binOp'},
        {value: 'metric', type: 'word'},
      ],
    ],
    [
      'label selector',
      `metric{code="500"}`,
      [
        {value: 'metric', type: 'word'},
        {value: '{', type: 'lBrace'},
        {value: 'code', type: 'word'},
        {value: '=', type: 'labelOp'},
        {value: '"500"', type: 'string'},
        {value: '}', type: 'rBrace'},
      ],
    ],
    [
      'multiple selectors',
      `metric{code="500",label="value"}`,
      [
        {value: 'metric', type: 'word'},
        {value: '{', type: 'lBrace'},
        {value: 'code', type: 'word'},
        {value: '=', type: 'labelOp'},
        {value: '"500"', type: 'string'},
        {value: ',', type: 'comma'},
        {value: 'label', type: 'word'},
        {value: '=', type: 'labelOp'},
        {value: '"value"', type: 'string'},
        {value: '}', type: 'rBrace'},
      ],
    ],
    [
      '> operator',
      `metric > metric`,
      [
        {value: 'metric', type: 'word'},
        {value: '>', type: 'binOp'},
        {value: 'metric', type: 'word'},
      ],
    ],
    [
      '>= operator',
      `metric >= metric`,
      [
        {value: 'metric', type: 'word'},
        {value: '>=', type: 'binOp'},
        {value: 'metric', type: 'word'},
      ],
    ],
    [
      'label list (e.g. by (label1, label2) clause)',
      `by (label1, label2)`,
      [
        {value: 'by', type: 'clauseOp'},
        {value: '(', type: 'lParen'},
        {value: 'label1', type: 'word'},
        {value: ',', type: 'comma'},
        {value: 'label2', type: 'word'},
        {value: ')', type: 'rParen'},
      ],
    ],
    [
      'simple aggregation',
      `sum(metric)`,
      [
        {value: 'sum', type: 'aggOp'},
        {value: '(', type: 'lParen'},
        {value: 'metric', type: 'word'},
        {value: ')', type: 'rParen'},
      ],
    ],
    [
      'simple function',
      `rate(1)`,
      [
        {value: 'rate', type: 'functionName'},
        {value: '(', type: 'lParen'},
        {value: '1', type: 'scalar'},
        {value: ')', type: 'rParen'},
      ],
    ],
    [
      'floating point scalar',
      `vector(-1.234)`,
      [
        {value: 'vector', type: 'functionName'},
        {value: '(', type: 'lParen'},
        {value: '-1.234', type: 'scalar'},
        {value: ')', type: 'rParen'},
      ],
    ],
    [
      'time duration',
      `[5m]`,
      [
        {value: '[', type: 'lBracket'},
        {value: '5m', type: 'duration'},
        {value: ']', type: 'rBracket'},
      ],
    ],
    [
      'long duration',
      `[50d]`,
      [
        {value: '[', type: 'lBracket'},
        {value: '50d', type: 'duration'},
        {value: ']', type: 'rBracket'},
      ],
    ],
    [
      'aggregated threshold',
      `avg(rate(http_status{code="500"}[5m])) by (region) > 5`,
      [
        {value: 'avg', type: 'aggOp'},
        {value: '(', type: 'lParen'},
        {value: 'rate', type: 'functionName'},
        {value: '(', type: 'lParen'},
        {value: 'http_status', type: 'word'},
        {value: '{', type: 'lBrace'},
        {value: 'code', type: 'word'},
        {value: '=', type: 'labelOp'},
        {value: '"500"', type: 'string'},
        {value: '}', type: 'rBrace'},
        {value: '[', type: 'lBracket'},
        {value: '5m', type: 'duration'},
        {value: ']', type: 'rBracket'},
        {value: ')', type: 'rParen'},
        {value: ')', type: 'rParen'},
        {value: 'by', type: 'clauseOp'},
        {value: '(', type: 'lParen'},
        {value: 'region', type: 'word'},
        {value: ')', type: 'rParen'},
        {value: '>', type: 'binOp'},
        {value: '5', type: 'scalar'},
      ],
    ],
  ];

  test.each(testCases)('%s', (_, input, expectedTokens) => {
    expect(Tokenize(input)).toEqual(expectedTokens);
  });
});
