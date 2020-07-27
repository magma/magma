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

const customWebpackConfigBuilder = require('./webpack.config.js').default;

module.exports = {
  stories: ['../**/*.stories.@(js|mdx)'],
  addons: [
    {
      name: '@storybook/addon-docs',
      options: {
        sourceLoaderOptions: null,
      },
    },
    '@storybook/addon-actions',
    '@storybook/addon-links',
  ],
  webpackFinal: config => customWebpackConfigBuilder({config}),
};
