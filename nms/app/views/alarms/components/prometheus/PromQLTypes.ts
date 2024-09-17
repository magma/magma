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

export type BinaryComparator = keyof typeof BINARY_COMPARATORS_MAP;
export const BINARY_COMPARATORS_MAP = {
  '==': '==',
  '!=': '!=',
  '>': '>',
  '<': '<',
  '<=': '<=',
  '>=': '>=',
};
export const BINARY_COMPARATORS = Object.keys(BINARY_COMPARATORS_MAP) as Array<
  BinaryComparator
>;

export type BinaryArithmetic = keyof typeof BINARY_ARITHMETIC_OPS_MAP;
export const BINARY_ARITHMETIC_OPS_MAP = {
  '+': '+',
  '-': '-',
  '*': '*',
  '/': '/',
  '%': '%',
  '^': '^',
};
export const BINARY_ARITHMETIC_OPS = Object.keys(
  BINARY_ARITHMETIC_OPS_MAP,
) as Array<BinaryArithmetic>;

export type BinarySet = 'and' | 'or' | 'unless';
export const BINARY_SET_OPS: Array<BinarySet> = ['and', 'or', 'unless'];

export const BINARY_OPERATORS = [
  ...BINARY_COMPARATORS,
  ...BINARY_ARITHMETIC_OPS,
  ...BINARY_SET_OPS,
];

export type LabelOperator = '=' | '!=' | '=~' | '!~';
export const LABEL_OPERATORS: Array<LabelOperator> = ['=', '!=', '=~', '!~'];

export type AggregationOperator = keyof typeof AGGREGATION_OPERATORS_MAP;
const AGGREGATION_OPERATORS_MAP = {
  sum: 'sum',
  min: 'min',
  max: 'max',
  avg: 'avg',
  stddev: 'stddev',
  stdvar: 'stdvar',
  count: 'count',
  count_values: 'count_values',
  quantile: 'quantile',
  bottomk: 'bottomk',
  topk: 'topk',
};
export const AGGREGATION_OPERATORS: Array<string> = Object.keys(
  AGGREGATION_OPERATORS_MAP,
);

export type FunctionName = keyof typeof FUNCTION_NAMES_MAP;
const FUNCTION_NAMES_MAP = {
  abs: 'abs',
  absent: 'absent',
  ceil: 'ceil',
  changes: 'changes',
  clamp_max: 'clamp_max',
  clamp_min: 'clamp_min',
  day_of_month: 'day_of_month',
  day_of_week: 'day_of_week',
  days_in_month: 'days_in_month',
  delta: 'deriv',
  exp: 'exp',
  floor: 'floor',
  histogram_quantile: 'histogram_quantile',
  holt_winters: 'holt_winters',
  hour: 'hour',
  idelta: 'idelta',
  increase: 'increase',
  irate: 'irate',
  label_join: 'label_join',
  label_replace: 'label_replace',
  ln: 'ln',
  log2: 'log2',
  log10: 'log10',
  minute: 'minute',
  month: 'month',
  predict_linear: 'predict_linear',
  rate: 'rate',
  resets: 'resets',
  round: 'round',
  scalar: 'scalar',
  sort: 'sort',
  sort_desc: 'sort_desc',
  sqrt: 'sqrt',
  time: 'time',
  timestamp: 'timestamp',
  vector: 'vector',
  year: 'year',
  // 'Over time' functions, operating on range-vectors.
  // They differ from aggregation operators:
  // - aggregations operate on instant vectors and `by`/`without` dimensions
  // - functions operate on range vectors and don't allow specifying dimensions
  sum_over_time: 'sum_over_time',
  min_over_time: 'min_over_time',
  max_over_time: 'max_over_time',
  avg_over_time: 'avg_over_time',
  stddev_over_time: 'stddev_over_time',
  stdvar_over_time: 'stdvar_over_time',
  quantile_over_time: 'quantile_over_time',
  count_over_time: 'count_values',
};
export const FUNCTION_NAMES: Array<string> = Object.keys(FUNCTION_NAMES_MAP);

export type AggrClauseType = 'by' | 'without';
export const AGGR_CLAUSE_TYPES: Array<AggrClauseType> = ['by', 'without'];

export type MatchClauseType = 'on' | 'ignoring';
export const MATCH_CLAUSE_TYPES: Array<MatchClauseType> = ['on', 'ignoring'];

export type GroupClauseType = 'group_left' | 'group_right';
export const GROUP_CLAUSE_TYPES: Array<GroupClauseType> = [
  'group_left',
  'group_right',
];

export class SyntaxError extends Error {
  constructor(message: string) {
    super(message);
    this.name = this.constructor.name;
    if (typeof Error.captureStackTrace === 'function') {
      Error.captureStackTrace(this, this.constructor);
    } else {
      this.stack = new Error(message).stack;
    }
  }
}
