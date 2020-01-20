/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as Moo from 'moo';

import {Range} from './PromQL';

import {
  AGGREGATION_OPERATORS,
  BINARY_OPERATORS,
  FUNCTION_NAMES,
  LABEL_OPERATORS,
} from './PromQLTypes';

const lexerRules = {
  WS: new RegExp(`[ \\t]+`),
  lBrace: '{',
  rBrace: '}',
  lParen: '(',
  rParen: ')',
  lBracket: '[',
  rBracket: ']',
  comma: ',',
  range: {
    match: new RegExp(`[0-9]+[smhdwy]`),
    value: s =>
      new Range(
        parseInt(s.substring(0, s.length - 1), 10),
        s.substring(s.length - 1),
      ),
  },
  scalar: {
    match: new RegExp(`[-+]?[0-9]*\\.?[0-9]+`),
    value: s => parseFloat(s),
  }, // matches floating point and integers
  aggOp: AGGREGATION_OPERATORS,
  functionName: FUNCTION_NAMES,
  binOp: BINARY_OPERATORS,
  labelOp: LABEL_OPERATORS,
  word: new RegExp(`\\w+`),
  string: {match: new RegExp(`"[^"]*"`), value: s => s.slice(1, -1)}, // strip quotes from string
};

export type Token = {
  value: string,
  type: TokenType,
};

type TokenType = $Keys<typeof lexerRules>;

export const lexer = Moo.compile(lexerRules);
// Ignore whitespace tokens
lexer.next = (next => () => {
  let tok;
  while ((tok = next.call(lexer)) && tok.type === 'WS') {}
  return tok;
})(lexer.next);

export function Tokenize(input: string): Array<Token> {
  lexer.reset(input);

  const tokens = [];
  let token;
  while ((token = lexer.next())) {
    tokens.push({value: token.value, type: token.type});
  }
  return tokens;
}
