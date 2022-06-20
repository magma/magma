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
 */

// NOTE: Regex based routes for paths that don't require logged in user access
export default [
  /^\/favicon.ico$/,
  /^\/healthz$/,
  /^\/user\/login(\?.*)?/,
  /^\/user\/onboarding(\?.*)?/,
  /^\/([a-z_-]+\/)?static\/css/,
  /^\/([a-z_-]+\/)?static\/dist\//,
  /^\/([a-z_-]+\/)?static\/fonts/,
  /^\/([a-z_-]+\/)?static\/images/,
  /^\/([a-z_-]+\/)?user\/me$/,
  /^\/([a-z_-]+\/)?user\/login(\?.*)?$/,
  /^\/([a-z_-]+\/)?user\/login\/oidc$/,
  /^\/([a-z_-]+\/)?user\/login\/oidc\/callback/,
  /^\/([a-z_-]+\/)?user\/login\/saml$/,
  /^\/([a-z_-]+\/)?user\/login\/saml\/callback/,
  /^\/([a-z_-]+\/)?user\/logout$/,
  /^\/([a-z_-]+\/)?__webpack_hmr.js/,
  /^\/authconfig$/,
];
