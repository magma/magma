/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *

 * @format
 */

module.exports = {
  ignore: [
    filename => {
      if (filename.indexOf('fbcnms') >= 0) {
        return false;
      } else if (filename.indexOf('magmalte') >= 0) {
        return false;
      } else if (filename.indexOf('node_modules') >= 0) {
        return true;
      }
      return false;
    },
  ],
  presets: [
    [
      '@babel/preset-env',
      {
        targets: {
          node: 'current',
          chrome: '58',
        },
        corejs: 3,
        useBuiltIns: 'entry',
      },
    ],
    '@babel/preset-flow',
    '@babel/preset-react',
  ],
  plugins: [
    'babel-plugin-lodash',
    'babel-plugin-relay',
    '@babel/plugin-proposal-class-properties',
    '@babel/plugin-proposal-nullish-coalescing-operator',
    '@babel/plugin-proposal-optional-chaining',
    '@babel/plugin-transform-react-jsx',
    'babel-plugin-fbt',
    'babel-plugin-fbt-runtime',
  ],
  env: {
    test: {
      sourceMaps: 'both',
    },
  },
};
