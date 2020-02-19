/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as PromQL from '../PromQL';
import {Parse} from '../PromQLParser';
import {Tokenize} from '../PromQLTokenizer';

class ErrorMatcher {
  messageRegex: ?RegExp;
  constructor(messageRegex: ?RegExp) {
    this.messageRegex = messageRegex;
  }
}

function expectSyntaxError(msg: ?RegExp): ErrorMatcher {
  return new ErrorMatcher(msg);
}

const testCases = [
  [
    'double-quoted string',
    `"this is \\" a \' string with \\\\, \\t, \\u0100, \\100, \\xA5, and \\n."`,
    [
      {
        value: `this is " a ' string with \\, \t, \u0100, @, \xA5, and \n.`,
        type: 'string',
      },
    ],
    null,
  ],
  [
    'single-quoted string',
    `'this is \\' a \" string with \\\\, \\t, \\u0100, \\100, \\xA5, and \\n.'`,
    [
      {
        value: `this is ' a " string with \\, \t, \u0100, @, \xA5, and \n.`,
        type: 'string',
      },
    ],
    null,
  ],
  [
    'back-ticked string',
    '`this is a string with no escaping: \\n \\" \\t \\u0100`',
    [
      {
        value: `this is a string with no escaping: \\n \\" \\t \\u0100`,
        type: 'string',
      },
    ],
    null,
  ],
  [
    'escaped string - with parser',
    '{l="\\"esc\\""}',
    [
      {value: '{', type: 'lBrace'},
      {value: 'l', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '"esc"', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector('', new PromQL.Labels().addEqual('l', '"esc"')),
  ],
  [
    'malformed double-quoted string',
    `"forgot to escape " the quote"`,
    expectSyntaxError(),
    expectSyntaxError(),
  ],
  [
    'malformed single-quoted string',
    `'forgot to escpape ' the quote'`,
    expectSyntaxError(),
    expectSyntaxError(),
  ],
  [
    'marlformed back-ticked string',
    '`forgot to escape ` the quote`',
    expectSyntaxError(),
    expectSyntaxError(),
  ],
  [
    'unknown escape sequence',
    `"I had 99 problems \\w"`,
    expectSyntaxError(/unterminated escape/i),
    expectSyntaxError(/unterminated escape/i),
  ],
  [
    'malformed \\x escape sequence',
    `"\\xa "`,
    expectSyntaxError(/unterminated escape/i),
    expectSyntaxError(/unterminated escape/i),
  ],
  [
    'marformed \\u escape sequence',
    `"\\u010 "`,
    expectSyntaxError(/unterminated escape/i),
    expectSyntaxError(/unterminated escape/i),
  ],
  [
    'malformed \\U escape sequence',
    `"\\U0100 "`,
    expectSyntaxError(/unterminated escape/i),
    expectSyntaxError(/unterminated escape/i),
  ],
  [
    'single metric selector',
    'metric',
    [{value: 'metric', type: 'word'}],
    new PromQL.InstantSelector('metric'),
  ],
  [
    'empty label selector',
    '{}',
    [{value: '{', type: 'lBrace'}, {value: '}', type: 'rBrace'}],
    new PromQL.InstantSelector('', new PromQL.Labels()),
  ],
  [
    'whitespace',
    'metric and\tmetric',
    [
      {value: 'metric', type: 'word'},
      {value: 'and', type: 'binOp'},
      {value: 'metric', type: 'word'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      'and',
    ),
  ],
  [
    'just label selector',
    `{code="500"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector('', new PromQL.Labels().addEqual('code', '500')),
  ],
  [
    'label selector',
    `metric{code="500"}`,
    [
      {value: 'metric', type: 'word'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      'metric',
      new PromQL.Labels().addEqual('code', '500'),
    ),
  ],
  [
    'multiple selectors',
    `metric{code="500",label="value"}`,
    [
      {value: 'metric', type: 'word'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'label', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: 'value', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      'metric',
      new PromQL.Labels().addEqual('code', '500').addEqual('label', 'value'),
    ),
  ],
  [
    '> operator',
    `metric > metric`,
    [
      {value: 'metric', type: 'word'},
      {value: '>', type: 'binOp'},
      {value: 'metric', type: 'word'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      '>',
    ),
  ],
  [
    '>= operator',
    `metric >= metric`,
    [
      {value: 'metric', type: 'word'},
      {value: '>=', type: 'binOp'},
      {value: 'metric', type: 'word'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      '>=',
    ),
  ],
  [
    'label list (e.g. by (label1, label2) clause)',
    `by (label1, label2)`,
    [
      {value: 'by', type: 'word'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'word'},
      {value: ',', type: 'comma'},
      {value: 'label2', type: 'word'},
      {value: ')', type: 'rParen'},
    ],
    null,
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
    new PromQL.AggregationOperation('sum', [
      new PromQL.InstantSelector('metric'),
    ]),
  ],
  [
    'aggregation by',
    `sum(metric) by (label)`,
    [
      {value: 'sum', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'metric', type: 'word'},
      {value: ')', type: 'rParen'},
      {value: 'by', type: 'word'},
      {value: '(', type: 'lParen'},
      {value: 'label', type: 'word'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.AggregationOperation(
      'sum',
      [new PromQL.InstantSelector('metric')],
      new PromQL.Clause('by', ['label']),
    ),
  ],
  [
    'aggregation by (clause preceding',
    `sum by (label) (metric)`,
    [
      {value: 'sum', type: 'aggOp'},
      {value: 'by', type: 'word'},
      {value: '(', type: 'lParen'},
      {value: 'label', type: 'word'},
      {value: ')', type: 'rParen'},
      {value: '(', type: 'lParen'},
      {value: 'metric', type: 'word'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.AggregationOperation(
      'sum',
      [new PromQL.InstantSelector('metric')],
      new PromQL.Clause('by', ['label']),
    ),
  ],
  [
    'simple function',
    `rate(1)`,
    [
      {value: 'rate', type: 'functionName'},
      {value: '(', type: 'lParen'},
      {value: 1, type: 'scalar'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.Function('rate', [new PromQL.Scalar(1)]),
  ],
  [
    'function with string parameter',
    `count_values("version", build_version)`,
    [
      {value: 'count_values', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'version', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'build_version', type: 'word'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.AggregationOperation('count_values', [
      new PromQL.String('version'),
      new PromQL.InstantSelector('build_version'),
    ]),
  ],
  ['binary integer scalar', '0b101010', [{value: 42, type: 'scalar'}], null],
  [
    'octal integer scalar',
    '0o33653337357',
    [{value: 3735928559, type: 'scalar'}],
    null,
  ],
  ['decimal integer scalar', '1337', [{value: 1337, type: 'scalar'}], null],
  [
    'hexadecimal integer scalar',
    '0xfaceb00c',
    [{value: 4207849484, type: 'scalar'}],
    null,
  ],
  [
    'floating point scalar',
    `vector(-1.234)`,
    [
      {value: 'vector', type: 'functionName'},
      {value: '(', type: 'lParen'},
      {value: -1.234, type: 'scalar'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.Function('vector', [new PromQL.Scalar(-1.234)]),
  ],
  [
    'time duration',
    `[5m]`,
    [
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
    ],
    null,
  ],
  [
    'long duration',
    `[50d]`,
    [
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(50, 'd'), type: 'range'},
      {value: ']', type: 'rBracket'},
    ],
    null,
  ],
  [
    'range selector',
    `metric[50d]`,
    [
      {value: 'metric', type: 'word'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(50, 'd'), type: 'range'},
      {value: ']', type: 'rBracket'},
    ],
    new PromQL.RangeSelector(
      new PromQL.InstantSelector('metric'),
      new PromQL.Range(50, 'd'),
    ),
  ],
  [
    'aggregated threshold',
    `avg(rate(http_status{code="500"}[5m])) > 5`,
    [
      {value: 'avg', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'rate', type: 'functionName'},
      {value: '(', type: 'lParen'},
      {value: 'http_status', type: 'word'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
      {value: ')', type: 'rParen'},
      {value: ')', type: 'rParen'},
      {value: '>', type: 'binOp'},
      {value: 5, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.AggregationOperation('avg', [
        new PromQL.Function('rate', [
          new PromQL.RangeSelector(
            new PromQL.InstantSelector(
              'http_status',
              new PromQL.Labels().addEqual('code', '500'),
            ),
            new PromQL.Range(5, 'm'),
          ),
        ]),
      ]),
      new PromQL.Scalar(5),
      '>',
    ),
  ],
  [
    'aggregated threshold with by clause',
    `avg(rate(http_status{code="500"}[5m])) by (region, code) > 5`,
    [
      {value: 'avg', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'rate', type: 'functionName'},
      {value: '(', type: 'lParen'},
      {value: 'http_status', type: 'word'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'word'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
      {value: ')', type: 'rParen'},
      {value: ')', type: 'rParen'},
      {value: 'by', type: 'word'},
      {value: '(', type: 'lParen'},
      {value: 'region', type: 'word'},
      {value: ',', type: 'comma'},
      {value: 'code', type: 'word'},
      {value: ')', type: 'rParen'},
      {value: '>', type: 'binOp'},
      {value: 5, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.AggregationOperation(
        'avg',
        [
          new PromQL.Function('rate', [
            new PromQL.RangeSelector(
              new PromQL.InstantSelector(
                'http_status',
                new PromQL.Labels().addEqual('code', '500'),
              ),
              new PromQL.Range(5, 'm'),
            ),
          ]),
        ],
        new PromQL.Clause('by', ['region', 'code']),
      ),
      new PromQL.Scalar(5),
      '>',
    ),
  ],
  [
    'instant selector offset',
    `http_requests_total offset 5m`,
    [
      {value: 'http_requests_total', type: 'word'},
      {value: 'offset', type: 'word'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
    ],
    new PromQL.InstantSelector(
      'http_requests_total',
      new PromQL.Labels(),
      new PromQL.Range(5, 'm'),
    ),
  ],
  [
    'binary operation with match clause',
    `metric / on (label1,label2) metric2`,
    [
      {value: 'metric', type: 'word'},
      {value: '/', type: 'binOp'},
      {value: 'on', type: 'word'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'word'},
      {value: ',', type: 'comma'},
      {value: 'label2', type: 'word'},
      {value: ')', type: 'rParen'},
      {value: 'metric2', type: 'word'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric2'),
      '/',
      new PromQL.VectorMatchClause(
        new PromQL.Clause('on', ['label1', 'label2']),
      ),
    ),
  ],
  [
    'metric that starts with clause operator name',
    `bytes_received > 0`,
    [
      {value: 'bytes_received', type: 'word'},
      {value: '>', type: 'binOp'},
      {value: 0, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('bytes_received'),
      new PromQL.Scalar(0),
      '>',
    ),
  ],
  [
    'simple equality expression',
    `up == 0`,
    [
      {value: 'up', type: 'word'},
      {value: '==', type: 'binOp'},
      {value: 0, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('up'),
      new PromQL.Scalar(0),
      '==',
    ),
  ],
];

describe('Tokenize', () => {
  test.each(testCases)('%s', (name, input, expectedTokens, _) => {
    if (expectedTokens instanceof ErrorMatcher) {
      expect(() => Tokenize(input)).toThrowError(expectedTokens.messageRegex);
    } else {
      expect(Tokenize(input)).toEqual(expectedTokens);
    }
  });
});

describe('Parser', () => {
  test.each(testCases)('%s', (name, input, _, expected) => {
    if (expected instanceof ErrorMatcher) {
      expect(() => Parse(input)).toThrowError(expected.messageRegex);
    } else if (expected !== null) {
      expect(Parse(input)).toEqual(expected);
    }
  });
});
