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
import type {NetworkContextType} from '../context/NetworkContext';
// $FlowFixMe migrated to typescript
import type {NetworkType} from '../../../shared/types/network';
import type {SectionsConfigs} from './Section';

import AppContext from '../../../app/components/context/AppContext';
import NetworkContext from '../context/NetworkContext';
import {
  CWF,
  FEG,
  LTE,
  coalesceNetworkType,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from '../../../shared/types/network';

import MagmaAPI from '../../../api/MagmaAPI';
import {getCWFSections} from '../cwf/CWFSections';
import {getFEGSections} from '../feg/FEGSections';
import {getLteSections} from '../lte/LteSections';
import {useContext, useEffect, useState} from 'react';

export default function useSections(): SectionsConfigs {
  const {networkId} = useContext<NetworkContextType>(NetworkContext);
  const {isFeatureEnabled} = useContext(AppContext);
  const [networkType, setNetworkType] = useState<NetworkType | null>(null);
  const alertsEnabled = isFeatureEnabled('alerts');

  useEffect(() => {
    const fetchNetworkType = async () => {
      if (networkId) {
        const networkType = (
          await MagmaAPI.networks.networksNetworkIdTypeGet({
            networkId,
          })
        ).data;
        setNetworkType(coalesceNetworkType(networkId, networkType));
      }
    };

    void fetchNetworkType();
  }, [networkId]);

  if (!networkId || networkType === null) {
    return [null, []];
  }

  switch (networkType) {
    case CWF:
      return getCWFSections();
    case FEG:
      return getFEGSections();
    case LTE:
    default: {
      return getLteSections(alertsEnabled);
    }
  }
}
