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

module.exports = webpackConfig.createDevWebpackConfig({
  hot: true,
  projectName: 'nms',
});
