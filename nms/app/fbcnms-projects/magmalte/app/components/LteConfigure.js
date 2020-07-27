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

import Apn from './network/Apn';
import Configure from './network/Configure';
import DataPlanConfig from './network/DataPlanConfig';
import NetworkConfig from './network/NetworkConfig';
import PoliciesConfig from './network/PoliciesConfig';
import React from 'react';
import UpgradeConfig from './network/UpgradeConfig';

export default function LteConfigure() {
  const tabs = [
    {
      component: DataPlanConfig,
      label: 'Data Plans',
      path: 'dataplans',
    },
    {
      component: Apn,
      label: 'APN Configuration',
      path: 'apns',
    },
    {
      component: NetworkConfig,
      label: 'Network Configuration',
      path: 'network',
    },
    {
      component: UpgradeConfig,
      label: 'Upgrades',
      path: 'upgrades',
    },
    {
      component: PoliciesConfig,
      label: 'Policies',
      path: 'policies',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}
