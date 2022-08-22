/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

function isWebpack(caller) {
  return !!(caller && caller.name === 'babel-loader');
}

module.exports = api => {
  const enableReactRefresh = api.caller(isWebpack) && !api.env('production');

  return {
    ignore: ['node_modules'],
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
      // Enable development transform of React with new automatic runtime
      [
        '@babel/preset-react',
        {development: !api.env('production'), runtime: 'automatic'},
      ],
      '@babel/preset-typescript',
    ],
    plugins: [
      'babel-plugin-lodash',
      '@babel/plugin-proposal-class-properties',
      '@babel/plugin-proposal-nullish-coalescing-operator',
      '@babel/plugin-proposal-optional-chaining',
      '@babel/plugin-transform-react-jsx',
      ...(enableReactRefresh ? ['react-refresh/babel'] : []),
    ],
    env: {
      test: {
        sourceMaps: 'both',
      },
    },
  };
};
