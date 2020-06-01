/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *

 * @format
 */

module.exports = {
  presets: ['@babel/preset-react', '@babel/preset-flow'],
  plugins: [
    'babel-plugin-lodash',
    '@babel/plugin-proposal-class-properties',
    '@babel/plugin-proposal-nullish-coalescing-operator',
    '@babel/plugin-proposal-optional-chaining',
    '@babel/plugin-transform-react-jsx',
  ],
  ignore: ['./node_modules', './lib'],
};
