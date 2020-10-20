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

import type {prom_alert_config} from '@fbcnms/magma-api';

export const xwfmAlerts: {[string]: prom_alert_config} = {
  'Test Auto Alert': {
    alert: 'Test Auto Alert',
    expr: 'test_metric > 0',
    labels: {severity: 'minor'},
    annotations: {description: 'A test of the automatic alert system'},
  },
};
