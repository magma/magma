/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

'use strict';

const webpackConfig = require('fbcnms-webpack-config/dev-webpack');
const paths = require('fbcnms-webpack-config/paths');

module.exports = webpackConfig.createDevWebpackConfig({
  projectName: 'inventory',
  extraPaths: [paths.resolveApp('../magmalte')],
  entry: {
    master: [paths.resolveApp('app/master.js')],
    onboarding: [paths.resolveApp('app/onboarding.js')],
  },
  hot: false,
});
