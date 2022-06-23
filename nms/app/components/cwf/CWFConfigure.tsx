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

import Configure from '../network/Configure';
import LoadingFiller from '../LoadingFiller';
import MagmaAPI from '../../../api/MagmaAPI';
import PoliciesConfig from '../network/PoliciesConfig';
import React from 'react';
import UpgradeConfig from '../network/UpgradeConfig';
import useMagmaAPI from '../../../api/useMagmaAPI';
import {useParams} from 'react-router-dom';

export default function CWFConfigure() {
  const tabs = [
    {
      component: UpgradeConfig,
      label: 'Upgrades',
      path: 'upgrades',
    },
    {
      component: CWFPolicies,
      label: 'Policies',
      path: 'policies',
    },
  ];
  return <Configure tabRoutes={tabs} />;
}

function CWFPolicies() {
  const {networkId} = useParams();

  const {response, isLoading} = useMagmaAPI(
    MagmaAPI.carrierWifiNetworks.cwfNetworkIdGet,
    {
      networkId: networkId!,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <PoliciesConfig mirrorNetwork={response?.federation?.feg_network_id} />
  );
}
