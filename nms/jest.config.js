/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @format
 */

const esModules = ['lodash-es'].join('|');

module.exports = {
  collectCoverageFrom: [
    '**/fbcnms-projects/**/*.js',
    '**/fbcnms-packages/**/*.js',
    '!**/__mocks__/**',
    '!**/__tests__/**',
    '!**/thrift/gen-nodejs/**',
    '!**/node_modules/**',
  ],
  coverageReporters: ['json', 'html'],
  moduleNameMapper: {
    '\\.(jpg|jpeg|png|gif|eot|otf|webp|svg|ttf|woff|woff2|mp4|webm|wav|mp3|m4a|aac|oga)$':
      '<rootDir>/__mocks__/fileMock.js',
    '\\.(css|less)$': 'identity-obj-proxy',
  },
  modulePathIgnorePatterns: [],
  projects: [
    {
      name: 'server',
      testEnvironment: 'node',
      testMatch: [
        '<rootDir>/fbcnms-projects/**/server/**/__tests__/*.js',
        '<rootDir>/fbcnms-packages/fbcnms-i18n/**/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
      transformIgnorePatterns: [`/node_modules/(?!${esModules})`],
    },
    {
      name: 'app',
      testEnvironment: 'jsdom',
      testMatch: [
        '<rootDir>/fbcnms-projects/**/app/**/__tests__/*.js',
        // run app/server shared tests in both node and jsdom environments
        '<rootDir>/fbcnms-projects/**/shared/**/__tests__/*.js',
      ],
      transform: {
        '^.+\\.js$': 'babel-jest',
      },
      transformIgnorePatterns: [`/node_modules/(?!${esModules})`],
    },
  ],
  testEnvironment: 'jsdom',
  testPathIgnorePatterns: ['/node_modules/'],
};
