/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

'use strict';

const webpackConfig = require('fbcnms-webpack-config/dev-webpack');
const paths = require('fbcnms-webpack-config/paths');

module.exports = webpackConfig.createDevWebpackConfig({
  hot: true,
  projectName: 'nms',
  entry: {
    master: paths.resolveApp('app/master.js'),
  },
});
