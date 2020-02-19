/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

export type BinaryComparator = $Keys<typeof BINARY_COMPARATORS_MAP>;
export const BINARY_COMPARATORS_MAP = {
  '==': '==',
  '!=': '!=',
  '>': '>',
  '<': '<',
  '<=': '<=',
  '>=': '>=',
};
export const BINARY_COMPARATORS: Array<BinaryComparator> = Object.keys(
  BINARY_COMPARATORS_MAP,
);

export type BinaryArithmetic = $Keys<typeof BINARY_ARITHMETIC_OPS_MAP>;
export const BINARY_ARITHMETIC_OPS_MAP = {
  '+': '+',
  '-': '-',
  '*': '*',
  '/': '/',
  '%': '%',
  '^': '^',
};
const BINARY_ARITHMETIC_OPS = Object.keys(BINARY_ARITHMETIC_OPS_MAP);

export type BinaryLogical = $Keys<typeof BINARY_LOGIC_OPS_MAP>;
export const BINARY_LOGIC_OPS_MAP = {and: 'and', or: 'or', unless: 'unless'};
const BINARY_LOGIC_OPS = Object.keys(BINARY_LOGIC_OPS_MAP);

export const BINARY_OPERATORS = [
  ...BINARY_COMPARATORS,
  ...BINARY_ARITHMETIC_OPS,
  ...BINARY_LOGIC_OPS,
];
export type BinaryOperator =
  | BinaryComparator
  | BinaryArithmetic
  | BinaryLogical;

export type LabelOperator = '=' | '!=' | '=~' | '!~';
export const LABEL_OPERATORS = ['=', '!=', '=~', '!~'];

export type AggregationOperator = $Keys<typeof AGGREGATION_OPERATORS_MAP>;
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
  sum_over_time: 'sum_over_time',
  min_over_time: 'min_over_time',
  max_over_time: 'max_over_time',
  avg_over_time: 'avg_over_time',
  stddev_over_time: 'stddev_over_time',
  stdvar_over_time: 'stdvar_over_time',
  count_over_time: 'count_over_time',
  quantile_over_time: 'quantile_over_time',
  count_over_time: 'count_values',
};
export const AGGREGATION_OPERATORS: Array<string> = Object.keys(
  AGGREGATION_OPERATORS_MAP,
);

export type FunctionName = $Keys<typeof FUNCTION_NAMES_MAP>;
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
};
export const FUNCTION_NAMES: Array<string> = Object.keys(FUNCTION_NAMES_MAP);

export type ClauseOperator = $Keys<typeof CLAUSE_OPS>;
const CLAUSE_OPS_MAP = {
  by: 'by',
  on: 'on',
  unless: 'unless',
  without: 'without',
  ignoring: 'ignoring',
};
export const CLAUSE_OPS: Array<string> = Object.keys(CLAUSE_OPS_MAP);

export type GroupOperator = $Keys<typeof GROUP_OPS>;
const GROUP_OPS_MAP = {
  group_left: 'group_left',
  group_right: 'group_right',
};
export const GROUP_OPS: Array<string> = Object.keys(GROUP_OPS_MAP);

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
