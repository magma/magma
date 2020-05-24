/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @noformat
 */

const fbStrict = require('eslint-config-fb-strict');

// This pattern will match these texts:
//   var Foo = require('Foo');
//   var Bar = require('Foo').Bar;
//   var BarFoo = require(Bar + 'Foo');
//   var {Bar, Foo} = require('Foo');
//   import type {Bar, Foo} from 'Foo';
//   import {Bar, Foo} from 'Foo';
//   } from 'Foo';
// Also supports 'let' and 'const'.
const variableNamePattern = String.raw`\s*[a-zA-Z_$][a-zA-Z_$\d]*\s*`;
const atLeastOneVariablePattern =
  '\\{?' + variableNamePattern +
  '(?:,' + variableNamePattern + ')*\\}?';
const importStatement = String.raw`^(?:var|let|const|import type|import)\s+` +
  atLeastOneVariablePattern;
const maxLenIgnorePattern =
  '(?:' + importStatement + '|\\})' +
  String.raw`\s*(?:=\s*require\(|from)[a-zA-Z_+./"'\s\d\-]+\)?[^;\n]*[;\n]`;
const path = require('path');
const {buildSchema, printSchema} = require('graphql');
const fs = require('fs');

let schemaPath = path.resolve(
  __dirname,
  '../../../../fbcode/fbc/symphony/graph/graphql/schema/symphony.graphql',
);
if (!fs.existsSync(schemaPath)) {
  schemaPath = path.resolve(
    __dirname,
    '../../../graph/graphql/schema/symphony.graphql',
  )
}
const schemaFile = fs.readFileSync(schemaPath, {encoding: 'utf8'});
const schemaObject = buildSchema(schemaFile);

const restrictedImportsRule = ['error',{
  'paths':[{
    'name': 'lodash-es',
    'message': 'Please use lodash directly.',
  }],
}];

module.exports = Object.assign({}, fbStrict, {
  env: {
    browser: true,
    es6: true,
    jasmine: true,
  },
  globals: {
    ArrayBufferView: false,
    Buffer: false,
    Class: false,
    FormData: true,
    Iterable: false,
    Iterator: false,
    IteratorResult: false,
    Promise: false,
    __BUNDLE_START_TIME__: false,
    __filename: false,
  },
  parser: 'babel-eslint',
  parserOptions: {
    'ecmaFeatures': {
      'jsx': true,
    },
    'ecmaVersion': 6,
    'sourceType': 'module',
  },
  plugins: [
    'flowtype',
    'header',
    'import',
    'jest',
    'lint',
    'node',
    'prettier',
    'react',
    'react-hooks',
    'relay',
    'sort-imports-es6-autofix',
    'graphql',
  ],
  rules: {
    'comma-dangle': ['warn', 'always-multiline'],
    'max-len': ['warn', {
      'ignorePattern': maxLenIgnorePattern,
      'ignoreStrings': true,
      'ignoreTemplateLiterals': true,
      'ignoreUrls': true,
    }],
    'no-alert': 'off',
    'no-console': ['warn', {allow: ['error', 'warn']}],
    'no-restricted-modules': restrictedImportsRule,
    'no-restricted-imports': restrictedImportsRule,
    'no-undef': 'error',
    'no-unused-vars': ['error', {
      'vars': 'all',
      'args': 'after-used',
      'argsIgnorePattern': '^_',
      'varsIgnorePattern': '^_',
    }],
    'no-var': 'error',
    'prefer-const': ['warn', {destructuring: 'all'}],
    'sort-keys': 'off',
    'no-warning-comments': [1, {'terms': ['$FlowFixMe'], 'location': 'anywhere'}],
    'strict': 'off',

    // Import Plugin
    // https://github.com/benmosher/eslint-plugin-import
    'import/default': 2,
    'import/export': 2,
    'import/named': 2,
    'import/namespace': 2,
    'import/no-unresolved': 2,

    'lint/cs-intent-use-injected-props': 'off',
    'lint/duplicate-class-function': 'off',
    'lint/flow-exact-props': 'off',
    'lint/flow-exact-state': 'off',
    'lint/flow-readonly-props': 'off',
    'lint/only-plain-ascii': 'off',
    'lint/react-avoid-set-state-with-potentially-stale-state': 'off',
    'lint/sort-keys-fixable': 'off',
    'lint/strictly-null': 'off',
    'lint/test-only-props': 'off',

    // Flow Plugin
    // The following rules are made available via `eslint-plugin-flowtype`
    'flowtype/define-flow-type': 1,
    'flowtype/no-weak-types': [1],
    'flowtype/use-flow-type': 1,
    // The following is disabled for many file types in overrides
    'flowtype/require-valid-file-annotation': [2, 'always'],

    // Node Plugin
    // https://github.com/mysticatea/eslint-plugin-node
    'node/no-missing-require': 2,

    // Prettier Plugin
    // https://github.com/prettier/eslint-plugin-prettier
    'prettier/prettier': 2,

    // React Plugin
    // https://github.com/yannickcr/eslint-plugin-react
    'react/display-name': 0,
    'react/jsx-boolean-value': 0,
    'react/jsx-no-comment-textnodes': 1,
    'react/jsx-no-duplicate-props': 2,
    'react/jsx-no-undef': 2,
    'react/jsx-sort-props': 0,
    'react/jsx-uses-react': 1,
    'react/jsx-uses-vars': 1,
    'react/no-did-mount-set-state': 1,
    'react/no-did-update-set-state': 1,
    'react/no-is-mounted': 'warn',
    'react/no-multi-comp': 0,
    'react/no-string-refs': 1,
    'react/no-unknown-property': 0,
    'react/prop-types': 0,
    'react/react-in-jsx-scope': 1,
    'react/self-closing-comp': 1,
    'react/wrap-multilines': 0,

    // React Hooks Plugin
    'react-hooks/rules-of-hooks': 'error',
    'react-hooks/exhaustive-deps': 'warn',

    // Relay Plugin
    'relay/unused-fields': 'off',

    // sort-imports autofix plugin (sort-imports doesnt autofix)
    'sort-imports-es6-autofix/sort-imports-es6': [2, {
      'ignoreCase': false,
      'ignoreMemberSort': false,
      'memberSyntaxSortOrder': ['none', 'all', 'single', 'multiple'],
    }],
    'graphql/no-deprecated-fields': [
      'error',
      {
        env: 'relay',
        schemaString: printSchema(schemaObject),
        tagName: 'graphql',
      },
    ],

    // Jest Plugin
    // The following rules are made available via `eslint-plugin-jest`.
    // 'jest/no-disabled-tests': 1,
    // 'jest/no-focused-tests': 1,
    // 'jest/no-identical-title': 1,
    // 'jest/valid-expect': 1,
  },
  'overrides': [
   {
     'files': [
       '**/*eslint*/*.js',
       '.eslintrc.js',
       'babel.config.js',
       'jest.config.js',
       '**/flow-typed/**/*.js',
       './babel.config.js',
     ],
     'rules': {
       'flowtype/require-valid-file-annotation': 'off',
     },
   },
   {
     'files': [
       '**/__tests__/*.js',
     ],
     'rules': {
       'no-warning-comments': [0],
     },
   },
   {
     'files': [
       'flow-typed/**/*.js',
     ],
     'rules': {
       'flowtype/no-weak-types': [0],
     },
   }],
   'extends': [
    'plugin:relay/recommended',
  ],
});
