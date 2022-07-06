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
 *
 * @format
 */

const sharedPresets = [
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
  '@babel/preset-react',
];

module.exports = {
  ignore: [
    filename => {
      if (filename.indexOf('fbcnms') >= 0) {
        return false;
      } else if (filename.indexOf('node_modules') >= 0) {
        return true;
      }
      return false;
    },
  ],
  presets: [...sharedPresets, '@babel/preset-flow'],
  plugins: [
    'babel-plugin-lodash',
    '@babel/plugin-proposal-class-properties',
    '@babel/plugin-proposal-nullish-coalescing-operator',
    '@babel/plugin-proposal-optional-chaining',
    '@babel/plugin-transform-react-jsx',
  ],
  env: {
    test: {
      sourceMaps: 'both',
    },
  },
  overrides: [
    {
      test: /\.tsx?$/,
      presets: [...sharedPresets, '@babel/preset-typescript'],
    },
  ],
};
