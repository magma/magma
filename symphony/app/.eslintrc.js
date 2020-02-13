/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

// enforces copyright header to be present in every file
// eslint-disable-next-line max-len
const openSourcePattern = /\*\n \* Copyright 20\d{2}-present Facebook\. All Rights Reserved\.\n \*\n \* This source code is licensed under the BSD-style license found in the\n \* LICENSE file in the root directory of this source tree\.\n \*\n/;

module.exports.extends = ['eslint-config-fbcnms'];
module.exports.overrides = [
  {
    files: ['*'],
    rules: {
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
    files: [
      'fbcnms-packages/**/*.js',
      'fbcnms-projects/inventory/**/*.js',
      'fbcnms-projects/magmalte/**/*.js',
      'fbcnms-projects/platform-server/**/*.js',
    ],
    rules: {
      'header/header': [2, 'block', {pattern: openSourcePattern}],
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
      'fbcnms-packages/eslint-config-fbcnms/**/*.js',
      'fbcnms-packages/fbcnms-auth/**/*.js',
      'fbcnms-packages/fbcnms-babel-register/**/*.js',
      'fbcnms-packages/fbcnms-express-middleware/**/*.js',
      'fbcnms-packages/fbcnms-logging/**/*.js',
      'fbcnms-packages/fbcnms-magma-api/**/*.js',
      'fbcnms-packages/fbcnms-platform-server/**/*.js',
      'fbcnms-packages/fbcnms-relay/**/*.js',
      'fbcnms-packages/fbcnms-sequelize-models/**/*.js',
      'fbcnms-packages/fbcnms-ui/stories/**/*.js',
      'fbcnms-packages/fbcnms-util/**/*.js',
      'fbcnms-packages/fbcnms-i18n/**/*.js',
      'fbcnms-packages/fbcnms-webpack-config/**/*.js',
      'fbcnms-projects/*/config/webpack.*.js',
      'fbcnms-projects/*/scripts/**/*.js',
      'fbcnms-projects/*/server/**/*.js',
      'fbcnms-projects/platform-server/**/*.js',
      'scripts/fb/fbt/*.js',
    ],
    rules: {
      'no-console': 'off',
    },
  },
  {
    files: ['**/tgnms/**/*.js'],
    rules: {
      // tgnms doesn't want this because there's too many errors
      'flowtype/no-weak-types': 'off',
      'flowtype/require-valid-file-annotation': 'off',
    },
  },
];
module.exports.settings = {
  react: {
    version: 'detect',
  },
};
