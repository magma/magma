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

const ONE_MEGABYTE = 1000000;
export const DEFAULT_DATA_PLAN_ID = 'default';
export const BITRATE_MULTIPLIER = ONE_MEGABYTE;
export const DATA_PLAN_UNLIMITED_RATES = {
  max_bandwidth_dl: 200 * ONE_MEGABYTE,
  max_bandwidth_ul: 100 * ONE_MEGABYTE,
};
