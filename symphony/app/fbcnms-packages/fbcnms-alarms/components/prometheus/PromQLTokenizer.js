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
  SyntaxError,
} from './PromQLTypes';

type LexerRules = {[string]: LexerRule | $ReadOnlyArray<LexerRule>};
type LexerRule = string | RegExp | ComplexRule;
type ComplexRule = {match: RegExp, value: string => string | number | Range};

const lexerRules: LexerRules = {
  WS: /[ \t]+/,
  lBrace: '{',
  rBrace: '}',
  lParen: '(',
  rParen: ')',
  lBracket: '[',
  rBracket: ']',
  comma: ',',
  range: {
    match: /[0-9]+[smhdwy]/,
    value: s =>
      new Range(
        parseInt(s.substring(0, s.length - 1), 10),
        s.substring(s.length - 1),
      ),
  },
  scalar: [
    {
      // binary integers
      match: /0b[01]+/,
      value: s => Number.parseInt(s.substring(2), 2),
    },
    {
      // octal integers
      // in accordance with https://golang.org/pkg/strconv/#ParseInt spec
      match: /0o[0-7]+/,
      value: s => Number.parseInt(s.substring(2), 8),
    },
    {
      // hexadecimal integers
      match: /0x[0-9a-fA-F]+/,
      value: s => Number.parseInt(s.substring(2), 16),
    },
    {
      // decimal floats and integers
      // TODO remove sign from lexer
      // it can only be correcly processed by parser
      match: /[-+]?[0-9]*\.?[0-9]+(?:e|E[0-9]+)?/,
      value: s => Number.parseFloat(s),
    },
  ],
  aggOp: AGGREGATION_OPERATORS,
  functionName: FUNCTION_NAMES,
  binOp: BINARY_OPERATORS,
  labelOp: LABEL_OPERATORS,
  word: /\w+/,
  string: [
    {
      // double-quoted string with no escape sequences;
      // shortcut for performance
      match: /"[^"\\]*"/,
      value: s => s.slice(1, -1),
    },
    {
      // single-quoted string with no escape sequences;
      // shortcut for performance
      match: /'[^'\\]*'/,
      value: s => s.slice(1, -1),
    },
    {
      // back-ticked string, raw, no escape sequences
      match: /`[^`]*`/,
      value: s => s.slice(1, -1),
    },
    {
      // double-quoted string with escape sequences
      match: /"[^"\\]*(?:\\.[^"\\]*)*"/,
      value: s => unescapeString(s.slice(1, -1), `"`),
    },
    {
      // single-quoted string with escape sequences
      match: /'[^'\\]*(?:\\.[^'\\]*)*'/,
      value: s => unescapeString(s.slice(1, -1), `'`),
    },
  ],
};

/**
 * Loosely ported from
 * https://github.com/prometheus/prometheus/blob/46c52607611992aeee631a1e19f053d886ca34d4/util/strutil/quote.go
 * @param s string to unescape; must be stripped of opening and closing quotes.
 * @param quote type of quotes that the string was enclosed in;
 * it is needed to determine how to unescape \' and \".
 * @return unescaped string; throws an error if the string is malformed.
 */
function unescapeString(s: string, quote: '"' | "'"): string {
  let result = '';
  let i = 0;

  while (i < s.length) {
    const {char, newIndex} = unescapeCharAt(s, i, quote);
    result += char;
    i = newIndex;
  }

  return result;
}

type CharAndIndex = {char: string, newIndex: number};

function unescapeCharAt(s: string, i: number, quote: '"' | "'"): CharAndIndex {
  let currentChar = s.charAt(i);

  if (currentChar != '\\') {
    return {char: currentChar, newIndex: i + 1};
  }

  if (i === s.length - 1) {
    throw new SyntaxError(unterminatedEscape);
  }

  const currentIndex = i + 1;
  currentChar = s.charAt(currentIndex);

  if (currentChar === quote || currentChar === '\\') {
    return {char: currentChar, newIndex: currentIndex + 1};
  }

  if (currentChar in simpleUnescaper) {
    return {char: simpleUnescaper[currentChar], newIndex: currentIndex + 1};
  }

  const hexUnescaper = hexUnescapers[currentChar];
  if (hexUnescaper !== undefined) {
    return hexUnescaper(s, currentIndex + 1);
  }

  return unescapeOctAt(s, currentIndex);
}

const simpleUnescaper = {
  a: '\u0007', // bell
  b: '\b',
  f: '\f',
  n: '\n',
  r: '\r',
  t: '\t',
  v: '\v',
};

const hexUnescapers = {
  x: (s, i) => unescapeHexAt(s, i, 2),
  u: (s, i) => unescapeHexAt(s, i, 4),
  U: (s, i) => unescapeHexAt(s, i, 8),
};

function unescapeHexAt(s: string, i: number, width: 2 | 4 | 8): CharAndIndex {
  if (i >= s.length - width) {
    throw new SyntaxError(unterminatedEscape);
  }
  const hex = s.substring(i, i + width);
  // this will happily parse e.g. 'fz' into 15
  // this bug is left here intentionally, to balance correctness vs simplicty
  const codePoint = Number.parseInt(hex, 16);

  if (isNaN(codePoint)) {
    throw new SyntaxError(`${unterminatedEscape} (${hex})`);
  }

  return {char: String.fromCodePoint(codePoint), newIndex: i + width};
}

function unescapeOctAt(s: string, i: number): CharAndIndex {
  if (i >= s.length - 3) {
    throw new SyntaxError(unterminatedEscape);
  }

  const oct = s.substring(i, i + 3);
  // this will happily parse e.g. '2T' or '29' into 2
  // this bug is left here intentionally, to balance correctness vs simplicty
  const codePoint = Number.parseInt(oct, 8);

  if (isNaN(codePoint)) {
    throw new SyntaxError(`${unterminatedEscape} (${oct})`);
  }

  return {char: String.fromCodePoint(codePoint), newIndex: i + 3};
}

const unterminatedEscape = 'Unterminated escape sequence';

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
