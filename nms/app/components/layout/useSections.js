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
// $FlowFixMe migrated to typescript
import type {NetworkContextType} from '../context/NetworkContext';
// $FlowFixMe migrated to typescript
import type {NetworkType} from '../../../shared/types/network';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {SectionsConfigs} from './Section';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AppContext from '../../../app/components/context/AppContext';
import MagmaV1API from '../../../generated/WebClient';
// $FlowFixMe migrated to typescript
import NetworkContext from '../context/NetworkContext';
import {
  CWF,
  FEG,
  LTE,
  coalesceNetworkType,
  // $FlowFixMe migrated to typescript
} from '../../../shared/types/network';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {getCWFSections} from '../cwf/CWFSections';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {getFEGSections} from '../feg/FEGSections';
import {getLteSections} from '../lte/LteSections';
import {useContext, useEffect, useState} from 'react';

export default function useSections(): SectionsConfigs {
  const {networkId} = useContext<NetworkContextType>(NetworkContext);
  const {isFeatureEnabled} = useContext(AppContext);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const alertsEnabled = isFeatureEnabled('alerts');

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
      return getCWFSections();
    case FEG:
      return getFEGSections();
    case LTE:
    default: {
      return getLteSections(alertsEnabled);
    }
  }
}
