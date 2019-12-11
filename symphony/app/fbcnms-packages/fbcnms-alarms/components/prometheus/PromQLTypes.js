/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

export type BinaryComparator = '==' | '!=' | '>' | '<' | '<=' | '>=';
export type BinaryArithmetic = '+' | '-' | '*' | '/' | '%' | '^';
export type BinaryLogical = 'and' | 'or' | 'unless';

export type BinaryOperator =
  | BinaryComparator
  | BinaryArithmetic
  | BinaryLogical;

export type AggregationOperator =
  | 'sum'
  | 'min'
  | 'max'
  | 'avg'
  | 'stddev'
  | 'stdvar'
  | 'count'
  | 'quantile'
  | 'bottomk'
  | 'topk'
  | 'sum_over_time'
  | 'min_over_time'
  | 'max_over_time'
  | 'avg_over_time'
  | 'stddev_over_time'
  | 'stdvar_over_time'
  | 'count_over_time'
  | 'quantile_over_time'
  | 'count_values';

export type LabelOperator = '=' | '!=' | '=~' | '!~';

export type FunctionName =
  | 'abs'
  | 'absent'
  | 'ceil'
  | 'changes'
  | 'clamp_max'
  | 'clamp_min'
  | 'day_of_month'
  | 'day_of_week'
  | 'days_in_month'
  | 'delta'
  | 'deriv'
  | 'exp'
  | 'floor'
  | 'histogram_quantile'
  | 'holt_winters'
  | 'hour'
  | 'idelta'
  | 'increase'
  | 'irate'
  | 'label_join'
  | 'label_replace'
  | 'ln'
  | 'log2'
  | 'log10'
  | 'minute'
  | 'month'
  | 'predict_linear'
  | 'rate'
  | 'resets'
  | 'round'
  | 'scalar'
  | 'sort'
  | 'sort_desc'
  | 'sqrt'
  | 'time'
  | 'timestamp'
  | 'vector'
  | 'year';
