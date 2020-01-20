/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import grammar from './__generated__/PromQLGrammar.js';
import nearley from 'nearley';
import {Expression} from './PromQL';

export function Parser() {
  return new nearley.Parser(nearley.Grammar.fromCompiled(grammar));
}

export function Parse(input: ?string): ?Expression {
  if (!input) {
    throw 'empty input to parser';
  }
  const parser = Parser().feed(input);
  // parser returns array of all possible parsing trees, so access the first
  // element of results since this grammar should only produce 1 for each
  // input
  return parser.results[0];
}
