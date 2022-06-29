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

// enforces copyright header to be present in every file
// eslint-disable-next-line max-len
const openSourcePattern = /\*\n \* Copyright \d{4} The Magma Authors\.\n \*\n \* This source code is licensed under the BSD-style license found in the\n \* LICENSE file in the root directory of this source tree\.\n \*\n \* Unless required by applicable law or agreed to in writing, software\n \* distributed under the License is distributed on an "AS IS" BASIS,\n \* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied\.\n \* See the License for the specific language governing permissions and\n \* limitations under the License\.\n \*\n/;
// eslint-disable-next-line max-len
const newOpenSourcePattern = /Copyright \d{4} The Magma Authors\./;
const combinedOpenSourcePattern = new RegExp(
  '(' + newOpenSourcePattern.source + ')|(' + openSourcePattern.source + ')',
);

const restrictedImportsRule = [
  'error',
  {
    paths: [
      {
        name: 'lodash-es',
        message: 'Please use lodash directly.',
      },
    ],
  },
];

module.exports = {
  extends: ['plugin:import/typescript'],
  env: {
    browser: true,
    es6: true,
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
    ecmaFeatures: {
      jsx: true,
    },
    ecmaVersion: 6,
    sourceType: 'module',
  },
  plugins: [
    'header',
    'import',
    'jest',
    'lint',
    'node',
    'prettier',
    'react',
    'react-hooks',
    'sort-imports-es6-autofix',
  ],
  settings: {
    react: {
      version: 'detect',
    },
    'import/extensions': ['.js', '.jsx', '.ts', '.tsx'],
  },
  rules: {
    'no-alert': 'off',
    'no-console': ['warn', {allow: ['error', 'warn']}],
    'no-restricted-modules': restrictedImportsRule,
    'no-restricted-imports': restrictedImportsRule,
    'no-undef': 'error',
    'no-unused-vars': [
      'error',
      {
        vars: 'all',
        args: 'after-used',
        argsIgnorePattern: '^_',
        varsIgnorePattern: '^_',
      },
    ],
    'no-var': 'error',
    'prefer-const': ['warn', {destructuring: 'all'}],
    'sort-keys': 'off',
    'no-warning-comments': 'off',
    strict: 'off',

    // Import Plugin
    // https://github.com/benmosher/eslint-plugin-import
    'import/default': 2,
    'import/export': 2,
    'import/named': 2,
    'import/namespace': 2,
    'import/no-unresolved': [2, {ignore: ['@material-ui/core/@@SvgIcon$']}],

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

    // Node Plugin
    // https://github.com/mysticatea/eslint-plugin-node
    'node/no-missing-require': 2,

    // Prettier Plugin
    // https://github.com/prettier/eslint-plugin-prettier
    'prettier/prettier': [
      2,
      {
        singleQuote: true,
        trailingComma: 'all',
        bracketSpacing: false,
        jsxBracketSameLine: true,
        parser: 'flow',
      },
    ],

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
    'react-hooks/exhaustive-deps': 'error',

    // sort-imports autofix plugin (sort-imports doesnt autofix)
    'sort-imports-es6-autofix/sort-imports-es6': [
      2,
      {
        ignoreCase: false,
        ignoreMemberSort: false,
        memberSyntaxSortOrder: ['none', 'all', 'single', 'multiple'],
      },
    ],

    // Jest Plugin
    // The following rules are made available via `eslint-plugin-jest`.
    // 'jest/no-disabled-tests': 1,
    // 'jest/no-focused-tests': 1,
    // 'jest/no-identical-title': 1,
    // 'jest/valid-expect': 1,
  },
  overrides: [
    {
      files: ['**/*.js'],
      excludedFiles: '**/sequelize_models/migrations/*.js',
      plugins: ['flowtype'],
      rules: {
        'header/header': [2, 'block', {pattern: combinedOpenSourcePattern}],

        'flowtype/define-flow-type': 1,
        'flowtype/no-weak-types': [1],
        'flowtype/use-flow-type': 1,
        // The following is disabled for many file types below
        'flowtype/require-valid-file-annotation': [2, 'always'],
      },
    },
    {
      files: ['**/*.ts', '**/*.tsx', '**/sequelize_models/migrations/*.js'],
      extends: [
        'plugin:@typescript-eslint/recommended',
        'plugin:@typescript-eslint/recommended-requiring-type-checking',
      ],
      parser: '@typescript-eslint/parser',
      plugins: ['@typescript-eslint'],
      parserOptions: {
        project: './tsconfig.json',
      },
      rules: {
        'header/header': [2, 'block', {pattern: combinedOpenSourcePattern}],
        'prettier/prettier': [
          2,
          {
            singleQuote: true,
            trailingComma: 'all',
            bracketSpacing: false,
            jsxBracketSameLine: true,
            parser: 'typescript',
          },
        ],
        '@typescript-eslint/array-type': [2, {default: 'generic'}],
        '@typescript-eslint/no-explicit-any': 'off',
        '@typescript-eslint/no-empty-function': 'off',
        '@typescript-eslint/ban-ts-comment': 'off',
        '@typescript-eslint/no-non-null-assertion': 'off',
        '@typescript-eslint/unbound-method': 'off',
        '@typescript-eslint/no-unused-vars': [2, {ignoreRestSiblings: true}],
      },
    },
    {
      env: {
        node: true,
      },
      files: [
        '.eslintrc.js',
        'babel.config.js',
        'babelRegister.js',
        'jest.config.js',
        'jest.*.config.js',
        'config/*',
        'scripts/**/*',
        'server/**/*',
        'shared/**/*',
        'grafana/**/*',
      ],
      rules: {
        'no-console': 'off',
      },
    },
    {
      files: [
        '**/*eslint*/*.js',
        '.eslintrc.js',
        'babel.config.js',
        'babelRegister.js',
        'jest.config.js',
        '**/flow-typed/**/*.js',
        './babel.config.js',
      ],
      rules: {
        'flowtype/require-valid-file-annotation': 'off',
      },
    },
    {
      files: ['**/__tests__/*.js'],
      rules: {
        'no-warning-comments': [0],
      },
    },
    {
      files: ['flow-typed/**/*.js'],
      rules: {
        'flowtype/no-weak-types': [0],
      },
    },
    {
      files: ['*.mdx'],
      extends: ['plugin:mdx/overrides'],
      rules: {
        'flowtype/require-valid-file-annotation': 'off',
        'prettier/prettier': [
          2,
          {
            parser: 'mdx',
          },
        ],
      },
    },
    {
      files: ['.eslintrc.js'],
      rules: {
        quotes: ['warn', 'single'],
      },
    },
    {
      env: {
        jest: true,
        node: true,
        'jest/globals': true,
      },
      files: [
        '**/__mocks__/**/*.js',
        '**/__tests__/**/*.js',
        '**/tests/*.js',
        'testHelpers.js',
        'testData.js',
      ],
    },
  ],
};
