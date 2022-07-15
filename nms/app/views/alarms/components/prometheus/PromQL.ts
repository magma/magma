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

import type {
  AggrClauseType,
  AggregationOperator,
  BinaryArithmetic,
  BinarySet,
  FunctionName,
  GroupClauseType,
  LabelOperator,
  MatchClauseType,
  BinaryComparator as SimpleBinaryComparator,
} from './PromQLTypes';

type Value = string | number;

export interface Expression<T extends Value> {
  selectorName?: string | null | undefined;
  value?: T;
  op?: string;
  toPromQL(): string;
}

export class Function implements Expression<Value> {
  name: FunctionName;
  arguments: Array<Expression<Value>>;

  constructor(name: FunctionName, args: Array<Expression<Value>>) {
    this.name = name;
    this.arguments = args;
  }

  toPromQL(): string {
    return (
      `${this.name}(` +
      this.arguments.map(arg => arg.toPromQL()).join(',') +
      ')'
    );
  }
}

export class InstantSelector implements Expression<any> {
  selectorName: string | null | undefined;
  labels: Labels | null | undefined;
  offset: Range | null | undefined;

  constructor(
    selectorName: string | null | undefined,
    labels?: Labels | null,
    offset?: Range | null,
  ) {
    this.selectorName = selectorName;
    this.labels = labels || new Labels();
    this.offset = offset;
  }

  toPromQL(): string {
    return (
      (this.selectorName || '') +
      (this.labels ? this.labels.toPromQL() : '') +
      (this.offset ? ' offset ' + this.offset.toString() : '')
    );
  }

  setOffset(offset: Range) {
    this.offset = offset;
    return this;
  }
}

export class RangeSelector extends InstantSelector {
  range: Range;

  constructor(selector: InstantSelector, range: Range) {
    super(selector.selectorName, selector.labels);
    this.range = range;
  }

  toPromQL(): string {
    return `${super.toPromQL()}[${this.range.toString()}]`;
  }
}

export class Range {
  unit: string;
  value: number;

  constructor(value: number, unit: string) {
    this.unit = unit;
    this.value = value;
  }

  toString(): string {
    return `${this.value}${this.unit}`;
  }
}

/**
 * The modifier methods of Labels mutate the underlying data
 * and return `this` to enable chaining on constructors.
 */
export class Labels {
  labels: Array<Label>;
  constructor(labels?: Array<Label> | null) {
    this.labels = labels || [];
  }

  toPromQL(): string {
    if (this.labels.length === 0) {
      return '';
    }
    return '{' + this.labels.map(label => label.toString()).join(',') + '}';
  }

  addLabel(name: string, value: string, operator: LabelOperator) {
    this.labels.push(new Label(name, value, operator));
  }

  addEqual(name: string, value: string): Labels {
    this.labels.push(new Label(name, value, '='));
    return this;
  }

  addNotEqual(name: string, value: string): Labels {
    this.labels.push(new Label(name, value, '!='));
    return this;
  }

  addRegex(name: string, value: string): Labels {
    this.labels.push(new Label(name, value, '=~'));
    return this;
  }

  addNotRegex(name: string, value: string): Labels {
    this.labels.push(new Label(name, value, '!~'));
    return this;
  }

  setIndex(
    i: number,
    name: string,
    value: string,
    operator?: LabelOperator | null,
  ): Labels {
    if (i >= 0 && i < this.len()) {
      this.labels[i].name = name;
      this.labels[i].value = value;
      this.labels[i].operator = operator || this.labels[i].operator;
    }
    return this;
  }

  remove(i: number): Labels {
    if (i >= 0 && i < this.len()) {
      this.labels.splice(i, 1);
    }
    return this;
  }

  removeByName(name: string): Labels {
    this.labels = this.labels.filter(label => label.name !== name);
    return this;
  }

  len(): number {
    return this.labels.length;
  }

  copy(): Labels {
    const ret = new Labels();
    this.labels.forEach(label => {
      ret.addLabel(label.name, label.value, label.operator);
    });
    return ret;
  }
}

export class Label {
  name: string;
  value: string;
  operator: LabelOperator;

  constructor(name: string, value: string, operator: LabelOperator) {
    this.name = name;
    this.value = value;
    this.operator = operator;
  }

  toString(): string {
    return `${this.name}${this.operator}"${this.value}"`;
  }
}

export class Scalar implements Expression<number> {
  value: number;

  constructor(value: number) {
    this.value = value;
  }

  toPromQL(): string {
    return this.value.toString();
  }
}

export class BinaryOperation implements Expression<string | number> {
  lh: Expression<string | number>;
  rh: Expression<string | number>;
  operator: BinaryOperator;
  clause: VectorMatchClause | null | undefined;

  constructor(
    lh: Expression<string | number>,
    rh: Expression<string | number>,
    operator: BinaryOperator,
    clause?: VectorMatchClause | null,
  ) {
    this.lh = lh;
    this.rh = rh;
    this.operator = operator;
    this.clause = clause;
  }

  toPromQL(): string {
    return (
      `${this.lh.toPromQL()} ${this.operator.toString()} ` +
      (this.clause ? this.clause.toString() + ' ' : '') +
      `${this.rh.toPromQL()}`
    );
  }
}

export type BinaryOperator = BinaryArithmetic | BinarySet | BinaryComparator;

export class BinaryComparator {
  op: SimpleBinaryComparator;
  boolMode: boolean;

  constructor(op: SimpleBinaryComparator) {
    this.op = op;
    this.boolMode = false;
  }

  makeBoolean() {
    this.boolMode = true;
    return this;
  }

  makeRegular() {
    this.boolMode = false;
    return this;
  }

  toString(): string {
    return this.boolMode ? `${this.op} bool` : this.op;
  }
}

export class VectorMatchClause {
  matchClause: Clause<MatchClauseType>;
  groupClause: Clause<GroupClauseType> | null | undefined;

  constructor(
    matchClause: Clause<MatchClauseType>,
    groupClause?: Clause<GroupClauseType> | null,
  ) {
    this.matchClause = matchClause;
    this.groupClause = groupClause;
  }

  toString(): string {
    return (
      this.matchClause.toString() +
      (this.groupClause ? ' ' + this.groupClause.toString() : '')
    );
  }
}

export type ClauseType = AggrClauseType | MatchClauseType | GroupClauseType;
export class Clause<C extends ClauseType> {
  operator: C;
  labelList: Array<string>;

  constructor(operator: C, labelList: Array<string> = []) {
    this.operator = operator;
    this.labelList = labelList;
  }

  toString(): string {
    return (
      this.operator +
      (this.labelList.length > 0 ? ` (${this.labelList.join(',')})` : '')
    );
  }
}

export class AggregationOperation implements Expression<Value> {
  name: AggregationOperator;
  parameters: Array<Expression<Value>>;
  clause: Clause<AggrClauseType> | null | undefined;

  constructor(
    name: AggregationOperator,
    parameters: Array<Expression<Value>>,
    clause?: Clause<AggrClauseType> | null,
  ) {
    this.name = name;
    this.parameters = parameters;
    this.clause = clause;
  }

  toPromQL(): string {
    return (
      `${this.name}(` +
      this.parameters.map(param => param.toPromQL()).join(',') +
      ')' +
      (this.clause ? ' ' + this.clause.toString() : '')
    );
  }
}

export class String implements Expression<string> {
  value: string;

  constructor(value: string) {
    this.value = value;
  }

  toPromQL(): string {
    return `"${this.value}"`;
  }
}

export class SubQuery implements Expression<Value> {
  expr: Expression<Value>;
  range: Range;
  resolution: Range | null | undefined;
  offset: Range | null | undefined;

  constructor(
    expr: Expression<Value>,
    range: Range,
    resolution?: Range | null,
    offset?: Range | null,
  ) {
    this.expr = expr;
    this.range = range;
    this.resolution = resolution;
    this.offset = offset;
  }

  withOffset(offset: Range) {
    this.offset = offset;
    return this;
  }

  toPromQL(): string {
    const maybeStep = this.resolution != null ? this.resolution.toString() : '';
    return (
      `${this.expr.toPromQL()}[${this.range.toString()}:${maybeStep}]` +
      (this.offset ? ' offset ' + this.offset.toString() : '')
    );
  }
}
