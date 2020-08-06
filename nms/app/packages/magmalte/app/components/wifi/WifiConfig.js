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

import Configure from '@fbcnms/magmalte/app/components/network/Configure';
import React from 'react';
import WifiNetworkConfig from './WifiNetworkConfig';

export default function WifiConfig() {
  const tabs = [
    {
      component: WifiNetworkConfig,
      label: 'Network Configuration',
      path: '',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}
