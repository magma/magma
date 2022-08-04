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

import * as Moo from 'moo';

import {Range} from './PromQL';

import {
  AGGREGATION_OPERATORS,
  AGGR_CLAUSE_TYPES,
  BINARY_ARITHMETIC_OPS,
  BINARY_COMPARATORS,
  BINARY_SET_OPS,
  GROUP_CLAUSE_TYPES,
  LABEL_OPERATORS,
  MATCH_CLAUSE_TYPES,
  SyntaxError,
} from './PromQLTypes';

// TODO[TS-migration] as unknown as string casts because Moo.Rules says value must return a string
const lexerRules: Moo.Rules = {
  WS: /[ \t]+/,
  lBrace: '{',
  rBrace: '}',
  lParen: '(',
  rParen: ')',
  lBracket: '[',
  rBracket: ']',
  comma: ',',
  colon: ':',
  range: {
    match: /[0-9]+[smhdwy]/,
    value: (s: string) =>
      (new Range(
        parseInt(s.substring(0, s.length - 1), 10),
        s.substring(s.length - 1),
      ) as unknown) as string,
  },
  scalar: [
    {
      // binary integers
      match: /0b[01]+/,
      value: (s: string) =>
        (Number.parseInt(s.substring(2), 2) as unknown) as string,
    },
    {
      // octal integers
      // in accordance with https://golang.org/pkg/strconv/#ParseInt spec
      match: /0o[0-7]+/,
      value: (s: string) =>
        (Number.parseInt(s.substring(2), 8) as unknown) as string,
    },
    {
      // hexadecimal integers
      match: /0x[0-9a-fA-F]+/,
      value: s => (Number.parseInt(s.substring(2), 16) as unknown) as string,
    },
    {
      // decimal floats and integers
      // TODO remove sign from lexer
      // it can only be correcly processed by parser
      match: /[-+]?[0-9]*\.?[0-9]+(?:e|E[0-9]+)?/,
      value: s => (Number.parseFloat(s) as unknown) as string,
    },
  ],
  // `!=` needs explicit token because it is ambiguous:
  // can mean either vector comparator or label matcher.
  // Must be declared above binComp and labelOp, because their definitions
  // include `!=`, too.
  neq: '!=',
  binComp: BINARY_COMPARATORS,
  arithmetic: BINARY_ARITHMETIC_OPS,
  labelOp: LABEL_OPERATORS,
  // Allows greedy-matching identifiers, e.g.
  // `by` will be emitted as %clauseOp, but
  // `byteCount` will be emitted as %identifier
  identifier: {
    match: /\w+/,
    type: Moo.keywords({
      aggOp: AGGREGATION_OPERATORS,
      aggClause: AGGR_CLAUSE_TYPES,
      groupClause: GROUP_CLAUSE_TYPES,
      matchClause: MATCH_CLAUSE_TYPES,
      setOp: BINARY_SET_OPS,
    }),
  },
  string: [
    {
      // double-quoted string with no escape sequences;
      // shortcut for performance
      match: /"[^"\\]*"/,
      value: (s: string) => s.slice(1, -1),
    },
    {
      // single-quoted string with no escape sequences;
      // shortcut for performance
      match: /'[^'\\]*'/,
      value: (s: string) => s.slice(1, -1),
    },
    {
      // back-ticked string, raw, no escape sequences
      match: /`[^`]*`/,
      value: (s: string) => s.slice(1, -1),
    },
    {
      // double-quoted string with escape sequences
      match: /"[^"\\]*(?:\\.[^"\\]*)*"/,
      value: (s: string) => unescapeString(s.slice(1, -1), `"`),
    },
    {
      // single-quoted string with escape sequences
      match: /'[^'\\]*(?:\\.[^'\\]*)*'/,
      value: (s: string) => unescapeString(s.slice(1, -1), `'`),
    },
  ],
  // Comments must be stripped by tokenzier,
  // and therefore will not be present in the AST.
  comment: /#[^\n]*/,
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

type CharAndIndex = {char: string; newIndex: number};

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
    return {
      char: simpleUnescaper[currentChar as keyof typeof simpleUnescaper],
      newIndex: currentIndex + 1,
    };
  }

  if (currentChar in hexUnescapers) {
    return hexUnescapers[currentChar as keyof typeof hexUnescapers](
      s,
      currentIndex + 1,
    );
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
  x: (s: string, i: number) => unescapeHexAt(s, i, 2),
  u: (s: string, i: number) => unescapeHexAt(s, i, 4),
  U: (s: string, i: number) => unescapeHexAt(s, i, 8),
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
  value: string;
  type: TokenType;
};

type TokenType = keyof typeof lexerRules;

export const lexer = Moo.compile(lexerRules);
// Ignore whitespace and comment tokens
lexer.next = (next => () => {
  let tok;
  while (
    (tok = next.call(lexer)) &&
    tok.type &&
    ['WS', 'comment'].includes(tok.type)
  ) {}
  return tok;
})(lexer.next);

export function Tokenize(input: string): Array<Token> {
  lexer.reset(input);

  const tokens: Array<Token> = [];
  let token;
  while ((token = lexer.next())) {
    tokens.push({value: token.value, type: token.type!});
  }
  return tokens;
}
