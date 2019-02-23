/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * 
 * @format
 */

'use strict';

const autoprefixer = require('autoprefixer');
const paths = require('fbcnms-webpack-config/paths');
const webpack = require('webpack');

module.exports = baseConfig => {
  baseConfig.module.rules = [
    {
      test: /\.(js|jsx|mjs)$/,
      include: [paths.appSrc, paths.packagesDir],
      exclude: /node_modules/,
      loader: require.resolve('babel-loader'),
      options: {
        configFile: '../../babel.config.js',
        // This is a feature of `babel-loader` for webpack (not Babel
        // itself). It enables caching results in
        // ./node_modules/.cache/babel-loader/ directory for faster
        // rebuilds.
        cacheDirectory: true,
      },
    },
  ];
  return baseConfig;
};