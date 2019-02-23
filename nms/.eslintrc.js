/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 * @noformat
 */
'use strict';

module.exports.extends = ['eslint-config-fbcnms'];
module.exports.overrides = [
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
    ],
  },
  {
    env: {
      node: true,
    },
    files: [
      '.eslintrc.js',
      'babel.config.js',
      'jest.config.js',
      'fbcnms-packages/eslint-config-fbcnms/**/*.js',
      'fbcnms-packages/fbcnms-auth/**/*.js',
      'fbcnms-packages/fbcnms-babel-register/**/*.js',
      'fbcnms-packages/fbcnms-express-middleware/**/*.js',
      'fbcnms-packages/fbcnms-logging/**/*.js',
      'fbcnms-packages/fbcnms-ui/stories/**/*.js',
      'fbcnms-packages/fbcnms-util/**/*.js',
      'fbcnms-packages/fbcnms-webpack-config/**/*.js',
      'fbcnms-projects/*/config/webpack.*.js',
      'fbcnms-projects/*/scripts/**/*.js',
      'fbcnms-projects/*/server/**/*.js',
    ],
    rules: {
      'no-console': 'off',
    },
  },
];
