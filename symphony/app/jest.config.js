/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

module.exports = {
  collectCoverageFrom: [
    '**/fbcnms-projects/**/*.js',
    '**/fbcnms-packages/**/*.js',
    '!**/__mocks__/**',
    '!**/__tests__/**',
    '!**/fbcnms-packages/fbcnms-ui/stories/**',
    '!**/thrift/gen-nodejs/**',
    '!**/node_modules/**',
    '!**/fbcnms-packages/fbcnms-test/**',
  ],

  coverageReporters: ['json', 'html'],
  modulePathIgnorePatterns: [],
  projects: [
    {
      name: 'server',
      testEnvironment: 'node',
      testMatch: [
        '<rootDir>/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/server/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-auth/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-express-middleware/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-platform-server/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/platform-server/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/workflows/**/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-packages/fbcnms-util/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
    },
    {
      moduleNameMapper: {
        '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
          '<rootDir>/__mocks__/fileMock.js',
        '\\.(css|less)$': 'identity-obj-proxy',
      },
      name: 'app',
      setupFiles: [
        '<rootDir>/fbcnms-packages/fbcnms-babel-register/polyfill.js',
      ],
      testEnvironment: 'jsdom',
      testMatch: [
        '<rootDir>/fbcnms-projects/**/app/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-ui/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-alarms/(components|hooks)/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-packages/fbcnms-util/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-mobileapp/**/__tests__/*.js',
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
    },
  ],
  testEnvironment: 'jsdom',
  testPathIgnorePatterns: ['/node_modules/'],
};
