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

import * as PromQL from '../PromQL';
import {Parse} from '../PromQLParser';
import {Tokenize} from '../PromQLTokenizer';

class ErrorMatcher {
  messageRegex: RegExp | undefined;
  constructor(messageRegex: RegExp | undefined) {
    this.messageRegex = messageRegex;
  }
}

function expectSyntaxError(msg?: RegExp): ErrorMatcher {
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
      {value: 'l', type: 'identifier'},
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
    [{value: 'metric', type: 'identifier'}],
    new PromQL.InstantSelector('metric'),
  ],
  [
    'metric with the same name as registered function',
    'absent',
    [{value: 'absent', type: 'identifier'}],
    new PromQL.InstantSelector('absent'),
  ],
  [
    'standalone aggregation is not allowed',
    'avg',
    [{value: 'avg', type: 'aggOp'}],
    expectSyntaxError(/malformed promql/i),
  ],
  [
    'metric selector with leading colon',
    `:metric:name`,
    [
      {value: ':', type: 'colon'},
      {value: 'metric', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: 'name', type: 'identifier'},
    ],
    new PromQL.InstantSelector(':metric:name'),
  ],
  [
    'metric selector with trailing colon',
    `metric:name:`,
    [
      {value: 'metric', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: 'name', type: 'identifier'},
      {value: ':', type: 'colon'},
    ],
    new PromQL.InstantSelector('metric:name:'),
  ],
  [
    'metric selector with colons in name - complex',
    `:some:metric::name:{label=~"value"}`,
    [
      {value: ':', type: 'colon'},
      {value: 'some', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: 'metric', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: ':', type: 'colon'},
      {value: 'name', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: '{', type: 'lBrace'},
      {value: 'label', type: 'identifier'},
      {value: '=~', type: 'labelOp'},
      {value: 'value', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      ':some:metric::name:',
      new PromQL.Labels().addRegex('label', 'value'),
    ),
  ],
  [
    'empty label selector',
    '{}',
    [
      {value: '{', type: 'lBrace'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector('', new PromQL.Labels()),
  ],
  [
    'whitespace',
    'metric and\tmetric',
    [
      {value: 'metric', type: 'identifier'},
      {value: 'and', type: 'setOp'},
      {value: 'metric', type: 'identifier'},
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
      {value: 'code', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector('', new PromQL.Labels().addEqual('code', '500')),
  ],
  [
    'label name equal to registered function',
    `{abs="-1"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'abs', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '-1', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector('', new PromQL.Labels().addEqual('abs', '-1')),
  ],
  [
    'label names equal to keywords',
    `{by="this magic",group_left="is allowed",unless="I failed"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'by', type: 'aggClause'},
      {value: '=', type: 'labelOp'},
      {value: 'this magic', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'group_left', type: 'groupClause'},
      {value: '=', type: 'labelOp'},
      {value: 'is allowed', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'unless', type: 'setOp'},
      {value: '=', type: 'labelOp'},
      {value: 'I failed', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      '',
      new PromQL.Labels()
        .addEqual('by', 'this magic')
        .addEqual('group_left', 'is allowed')
        .addEqual('unless', 'I failed'),
    ),
  ],
  [
    'label names with colons are not allowed',
    `{some:label="value"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'some', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: 'label', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: 'value', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    expectSyntaxError(),
  ],
  [
    'label selector',
    `metric{code="500"}`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'label', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
      {value: '>', type: 'binComp'},
      {value: 'metric', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      new PromQL.BinaryComparator('>'),
    ),
  ],
  [
    '>= operator',
    `metric >= metric`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '>=', type: 'binComp'},
      {value: 'metric', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      new PromQL.BinaryComparator('>='),
    ),
  ],
  [
    '> operator bool mode',
    'metric > bool metric2',
    [
      {value: 'metric', type: 'identifier'},
      {value: '>', type: 'binComp'},
      {value: 'bool', type: 'identifier'},
      {value: 'metric2', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric2'),
      new PromQL.BinaryComparator('>').makeBoolean(),
    ),
  ],
  [
    '!= operator as vector comparator',
    `metric != metric`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '!=', type: 'neq'},
      {value: 'metric', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric'),
      new PromQL.BinaryComparator('!='),
    ),
  ],
  [
    '!= operator as label matcher',
    `{status!="500"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'status', type: 'identifier'},
      {value: '!=', type: 'neq'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      '',
      new PromQL.Labels().addNotEqual('status', '500'),
    ),
  ],
  [
    'label list (e.g. by (label1, label2) clause)',
    `by (label1, label2)`,
    [
      {value: 'by', type: 'aggClause'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'identifier'},
      {value: ',', type: 'comma'},
      {value: 'label2', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: 'by', type: 'aggClause'},
      {value: '(', type: 'lParen'},
      {value: 'label', type: 'identifier'},
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
      {value: 'by', type: 'aggClause'},
      {value: '(', type: 'lParen'},
      {value: 'label', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: '(', type: 'lParen'},
      {value: 'metric', type: 'identifier'},
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
      {value: 'rate', type: 'identifier'},
      {value: '(', type: 'lParen'},
      {value: 1, type: 'scalar'},
      {value: ')', type: 'rParen'},
    ],
    new PromQL.Function('rate', [new PromQL.Scalar(1)]),
  ],
  [
    'unknown function',
    `kolmogorovComplexity(input)`,
    [
      {value: 'kolmogorovComplexity', type: 'identifier'},
      {value: '(', type: 'lParen'},
      {value: 'input', type: 'identifier'},
      {value: ')', type: 'rParen'},
    ],
    expectSyntaxError(/unknown function/i),
  ],
  [
    'function names with colons are not allowed',
    `some:func(1)`,
    [
      {value: 'some', type: 'identifier'},
      {value: ':', type: 'colon'},
      {value: 'func', type: 'identifier'},
      {value: '(', type: 'lParen'},
      {value: 1, type: 'scalar'},
      {value: ')', type: 'rParen'},
    ],
    expectSyntaxError(),
  ],
  [
    'function with string parameter',
    `count_values("version", build_version)`,
    [
      {value: 'count_values', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'version', type: 'string'},
      {value: ',', type: 'comma'},
      {value: 'build_version', type: 'identifier'},
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
      {value: 'vector', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
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
      {value: 'rate', type: 'identifier'},
      {value: '(', type: 'lParen'},
      {value: 'http_status', type: 'identifier'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
      {value: ')', type: 'rParen'},
      {value: ')', type: 'rParen'},
      {value: '>', type: 'binComp'},
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
      new PromQL.BinaryComparator('>'),
    ),
  ],
  [
    'aggregated threshold with by clause',
    `avg(rate(http_status{code="500"}[5m])) by (region, code) > 5`,
    [
      {value: 'avg', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'rate', type: 'identifier'},
      {value: '(', type: 'lParen'},
      {value: 'http_status', type: 'identifier'},
      {value: '{', type: 'lBrace'},
      {value: 'code', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '500', type: 'string'},
      {value: '}', type: 'rBrace'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(5, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
      {value: ')', type: 'rParen'},
      {value: ')', type: 'rParen'},
      {value: 'by', type: 'aggClause'},
      {value: '(', type: 'lParen'},
      {value: 'region', type: 'identifier'},
      {value: ',', type: 'comma'},
      {value: 'code', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: '>', type: 'binComp'},
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
      new PromQL.BinaryComparator('>'),
    ),
  ],
  [
    'instant selector offset',
    `http_requests_total offset 5m`,
    [
      {value: 'http_requests_total', type: 'identifier'},
      {value: 'offset', type: 'identifier'},
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
      {value: 'metric', type: 'identifier'},
      {value: '/', type: 'arithmetic'},
      {value: 'on', type: 'matchClause'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'identifier'},
      {value: ',', type: 'comma'},
      {value: 'label2', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: 'metric2', type: 'identifier'},
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
    'binary operation with group clause',
    `metric + on(label1) group_left(label2) metric2`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '+', type: 'arithmetic'},
      {value: 'on', type: 'matchClause'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: 'group_left', type: 'groupClause'},
      {value: '(', type: 'lParen'},
      {value: 'label2', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: 'metric2', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric2'),
      '+',
      new PromQL.VectorMatchClause(
        new PromQL.Clause('on', ['label1']),
        new PromQL.Clause('group_left', ['label2']),
      ),
    ),
  ],
  [
    'binary operation with group clause without labels',
    `metric + on(label1) group_left metric2`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '+', type: 'arithmetic'},
      {value: 'on', type: 'matchClause'},
      {value: '(', type: 'lParen'},
      {value: 'label1', type: 'identifier'},
      {value: ')', type: 'rParen'},
      {value: 'group_left', type: 'groupClause'},
      {value: 'metric2', type: 'identifier'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('metric'),
      new PromQL.InstantSelector('metric2'),
      '+',
      new PromQL.VectorMatchClause(
        new PromQL.Clause('on', ['label1']),
        new PromQL.Clause('group_left'),
      ),
    ),
  ],
  [
    'metric that starts with clause operator name',
    `bytes_received > 0`,
    [
      {value: 'bytes_received', type: 'identifier'},
      {value: '>', type: 'binComp'},
      {value: 0, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('bytes_received'),
      new PromQL.Scalar(0),
      new PromQL.BinaryComparator('>'),
    ),
  ],
  [
    'simple equality expression',
    `up == 0`,
    [
      {value: 'up', type: 'identifier'},
      {value: '==', type: 'binComp'},
      {value: 0, type: 'scalar'},
    ],
    new PromQL.BinaryOperation(
      new PromQL.InstantSelector('up'),
      new PromQL.Scalar(0),
      new PromQL.BinaryComparator('=='),
    ),
  ],
  [
    'comments',
    `metric_name # comment with ### symbols`,
    [{value: 'metric_name', type: 'identifier'}],
    new PromQL.InstantSelector('metric_name'),
  ],
  [
    'octothorpes in strings are not comments',
    `{fragment="#index"}`,
    [
      {value: '{', type: 'lBrace'},
      {value: 'fragment', type: 'identifier'},
      {value: '=', type: 'labelOp'},
      {value: '#index', type: 'string'},
      {value: '}', type: 'rBrace'},
    ],
    new PromQL.InstantSelector(
      '',
      new PromQL.Labels().addEqual('fragment', '#index'),
    ),
  ],
  [
    'simple subquery',
    `metric[1h:10m]`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(1, 'h'), type: 'range'},
      {value: ':', type: 'colon'},
      {value: new PromQL.Range(10, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
    ],
    new PromQL.SubQuery(
      new PromQL.InstantSelector('metric'),
      new PromQL.Range(1, 'h'),
      new PromQL.Range(10, 'm'),
    ),
  ],
  [
    'simple subquery wihtout step',
    `metric[1h:]`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(1, 'h'), type: 'range'},
      {value: ':', type: 'colon'},
      {value: ']', type: 'rBracket'},
    ],
    new PromQL.SubQuery(
      new PromQL.InstantSelector('metric'),
      new PromQL.Range(1, 'h'),
    ),
  ],
  [
    'simple subquery with offset',
    `metric[1h:] offset 1d`,
    [
      {value: 'metric', type: 'identifier'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(1, 'h'), type: 'range'},
      {value: ':', type: 'colon'},
      {value: ']', type: 'rBracket'},
      {value: 'offset', type: 'identifier'},
      {value: new PromQL.Range(1, 'd'), type: 'range'},
    ],
    new PromQL.SubQuery(
      new PromQL.InstantSelector('metric'),
      new PromQL.Range(1, 'h'),
    ).withOffset(new PromQL.Range(1, 'd')),
  ],
  [
    'complex subquery',
    `avg(metric[1h:10m] offset 1d)[10m:]`,
    [
      {value: 'avg', type: 'aggOp'},
      {value: '(', type: 'lParen'},
      {value: 'metric', type: 'identifier'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(1, 'h'), type: 'range'},
      {value: ':', type: 'colon'},
      {value: new PromQL.Range(10, 'm'), type: 'range'},
      {value: ']', type: 'rBracket'},
      {value: 'offset', type: 'identifier'},
      {value: new PromQL.Range(1, 'd'), type: 'range'},
      {value: ')', type: 'rParen'},
      {value: '[', type: 'lBracket'},
      {value: new PromQL.Range(10, 'm'), type: 'range'},
      {value: ':', type: 'colon'},
      {value: ']', type: 'rBracket'},
    ],
    new PromQL.SubQuery(
      new PromQL.AggregationOperation('avg', [
        new PromQL.SubQuery(
          new PromQL.InstantSelector('metric'),
          new PromQL.Range(1, 'h'),
          new PromQL.Range(10, 'm'),
        ).withOffset(new PromQL.Range(1, 'd')),
      ]),
      new PromQL.Range(10, 'm'),
    ),
  ],
] as const;

describe('Tokenize', () => {
  it.each(testCases)('%s', (name, input, expectedTokens, _) => {
    if (expectedTokens instanceof ErrorMatcher) {
      expect(() => Tokenize(input)).toThrowError(expectedTokens.messageRegex);
    } else {
      expect(Tokenize(input)).toEqual(expectedTokens);
    }
  });
});

describe('Parser', () => {
  it.each(testCases)('%s', (name, input, _, expected) => {
    if (expected instanceof ErrorMatcher) {
      expect(() => Parse(input)).toThrowError(expected.messageRegex);
    } else if (expected !== null) {
      expect(Parse(input)).toEqual(expected);
    }
  });
});
