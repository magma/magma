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
 * @flow
 * @format
 */
import type {NetworkContextType} from '../context/NetworkContext';
import type {NetworkType} from '@fbcnms/types/network';
import type {SectionsConfigs} from '../layout/Section';

import AppContext from '@fbcnms/ui/context/AppContext';
import MagmaV1API from '../../../generated/WebClient';
import NetworkContext from '../context/NetworkContext';
import {CWF, FEG, LTE, coalesceNetworkType} from '@fbcnms/types/network';

import {getCWFSections} from '../cwf/CWFSections';
import {getFEGSections} from '../feg/FEGSections';
import {getLteSections, getLteSectionsV2} from '../lte/LteSections';
import {useContext, useEffect, useState} from 'react';

export default function useSections(): SectionsConfigs {
  const {networkId} = useContext<NetworkContextType>(NetworkContext);
  const {user, isFeatureEnabled} = useContext(AppContext);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const alertsEnabled = isFeatureEnabled('alerts');
  const logsEnabled = isFeatureEnabled('logs');
  const dashboardV2Enabled = isFeatureEnabled('dashboard_v2');
  let dashboardV2EnabledFegCwf = false;

  // enable dashboard v2 for cwf and feg in test mode
  if (user && user.tenant !== '') {
    if (user.tenant.endsWith('-test') && dashboardV2Enabled) {
      dashboardV2EnabledFegCwf = true;
    }
  }
  useEffect(() => {
    const fetchNetworkType = async () => {
      if (networkId) {
        const networkType = await MagmaV1API.getNetworksByNetworkIdType({
          networkId,
        });
        setNetworkType(coalesceNetworkType(networkId, networkType));
      }
    };

    fetchNetworkType();
  }, [networkId]);

  if (!networkId || networkType === null) {
    return [null, []];
  }

  switch (networkType) {
    case CWF:
      return getCWFSections(dashboardV2EnabledFegCwf);
    case FEG:
      return getFEGSections(dashboardV2EnabledFegCwf);
    case LTE:
    default: {
      if (dashboardV2Enabled) {
        return getLteSectionsV2(alertsEnabled);
      }
      return getLteSections(alertsEnabled, logsEnabled);
    }
  }
}
