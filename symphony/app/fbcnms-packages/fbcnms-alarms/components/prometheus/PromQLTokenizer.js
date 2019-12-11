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

import {
  AGGREGATION_OPERATORS,
  BINARY_OPERATORS,
  CLAUSE_OPS,
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
  duration: new RegExp(`[0-9]+[smhdwy]`),
  scalar: new RegExp(`[-+]?[0-9]*\\.?[0-9]+`), // matches floating point and integers
  aggOp: AGGREGATION_OPERATORS,
  functionName: FUNCTION_NAMES,
  labelOp: LABEL_OPERATORS,
  binOp: BINARY_OPERATORS,
  clauseOp: CLAUSE_OPS,
  word: new RegExp(`\\w+`),
  string: new RegExp(`"[^"]*"`),
};

type Token = {
  value: string,
  type: TokenType,
};

type TokenType = $Keys<typeof lexerRules>;

const lexer = Moo.compile(lexerRules);
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
