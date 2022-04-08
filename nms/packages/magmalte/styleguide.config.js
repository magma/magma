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
 * @flow strict-local
 * @format
 */

require('@fbcnms/babel-register');

module.exports = {
  webpackConfig: require('./config/webpack.config.js'),
  components: [
    'app/components/**/*.js',
    'app/components/*.js',
    'app/views/**/*.js',
    'app/views/*.js',
  ],
  ignore: [
    // TODO: Figure out way to pass through env variable MAPBOX_ACCESS_TOKEN,
    //       which is causing a lot of components to crash, and necessitates
    //       ignoring a number of components.
    //       See the browser console after bringing up styleguidist for error
    'app/components/insights/map/*.js',
    'app/components/insights/*.js',
    'app/components/layout/Section.js',
    'fbc_js_core/**/*.js',
    // 'fbc_js_core/alarms/**/*.js',
    // 'fbc_js_core/alarms/*.js',
    // 'fbc_js_core/auth/**/*.js',
    // 'fbc_js_core/auth/*.js',
    // 'fbc_js_core/express_middleware/**/*.js',
    // 'fbc_js_core/platform_server/**/*.js',
    // 'fbc_js_core/sequelize_models/**/*.js',
    // 'fbc_js_core/sequelize_models/*.js',
    // TODO: Fix this ui package
    // 'fbc_js_core/ui/**/*.js',
    '**/__tests__/*.js',
    '**/__tests__/**/*.js',
  ],
};

/*
Referencing the TODO earlier:

Cannot read properties of undefined (reading 'MAPBOX_ACCESS_TOKEN')
    at eval (webpack-internal:///./app/components/insights/map/styles.js:24:78)
    at Module../app/components/insights/map/styles.js (main.bundle.js:1994:1)
    at __webpack_require__ (main.bundle.js:854:30)
    at fn (main.bundle.js:151:20)
    at eval (webpack-internal:///./app/components/insights/map/MapView.js:8:65)
    at Module../app/components/insights/map/MapView.js (main.bundle.js:1982:1)
    at __webpack_require__ (main.bundle.js:854:30)
    at fn (main.bundle.js:151:20)
    at eval (webpack-internal:///./app/components/insights/Insights.js:9:70)
    at Module../app/components/insights/Insights.js (main.bundle.js:1898:1)
*/
