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
 *
 * @format
 */

// enforces copyright header to be present in every file
// eslint-disable-next-line max-len
const openSourcePattern = /\*\n \* Copyright 2020 The Magma Authors\.\n \*\n \* This source code is licensed under the BSD-style license found in the\n \* LICENSE file in the root directory of this source tree\.\n \*\n \* Unless required by applicable law or agreed to in writing, software\n \* distributed under the License is distributed on an "AS IS" BASIS,\n \* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied\.\n \* See the License for the specific language governing permissions and\n \* limitations under the License\.\n \*\n/;
// eslint-disable-next-line max-len
const newOpenSourcePattern = /Copyright 2020 The Magma Authors\./;
const combinedOpenSourcePattern = new RegExp(
  '(' + newOpenSourcePattern.source + ')|(' + openSourcePattern.source + ')',
);

module.exports.extends = ['eslint-config-fbcnms'];
module.exports.overrides = [
  {
    files: ['*'],
    rules: {
      'react-hooks/exhaustive-deps': 'error',
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
    },
    files: [
      '**/__mocks__/**/*.js',
      '**/__tests__/**/*.js',
      '**/tests/*.js',
      'testHelpers.js',
      'testData.js',
    ],
  },
  {
    files: ['packages/**/*.js'],
    rules: {
      'header/header': [2, 'block', {pattern: combinedOpenSourcePattern}],
    },
  },
  {
    env: {
      node: true,
    },
    files: [
      '.eslintrc.js',
      'babel.config.js',
      'jest.config.js',
      'jest.*.config.js',
      'packages/fbcnms-magma-api/**/*.js',
      'packages/magmalte/config/*.js',
      'packages/magmalte/scripts/**/*.js',
      'packages/magmalte/server/**/*.js',
      'packages/magmalte/grafana/**/*.js',
    ],
    rules: {
      'no-console': 'off',
    },
  },
];
module.exports.settings = {
  react: {
    version: 'detect',
  },
};
