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

const webpackConfig = require('fbcnms-webpack-config/production-webpack');
const paths = require('fbcnms-webpack-config/paths');

module.exports = webpackConfig.createProductionWebpackConfig({
  projectName: 'inventory',
  devtool: 'source-map',
  extraPaths: [paths.resolveApp('../magmalte')],
  entry: {
    master: [paths.resolveApp('app/master.js')],
    onboarding: [paths.resolveApp('app/onboarding.js')],
  },
});
