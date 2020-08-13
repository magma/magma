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

import NetworkContext from '../context/NetworkContext';
import {
  CWF,
  FEG,
  LTE,
  RHINO,
  SYMPHONY,
  THIRD_PARTY,
  WIFI,
  coalesceNetworkType,
} from '@fbcnms/types/network';
import type {NetworkContextType} from '../context/NetworkContext';
import type {NetworkType} from '@fbcnms/types/network';
import type {SectionsConfigs} from '../layout/Section';

import AppContext from '@fbcnms/ui/context/AppContext';
import {useContext, useEffect, useState} from 'react';

import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import {getCWFSections} from '../cwf/CWFSections';
import {getDevicesSections} from '@fbcnms/magmalte/app/components/devices/DevicesSections';
import {getFEGSections} from '../feg/FEGSections';
import {getLteSections} from '@fbcnms/magmalte/app/components/lte/LteSections';
import {getMeshSections} from '@fbcnms/magmalte/app/components/wifi/WifiSections';
import {getRhinoSections} from '@fbcnms/magmalte/app/components/rhino/RhinoSections';

export default function useSections(): SectionsConfigs {
  const {networkId} = useContext<NetworkContextType>(NetworkContext);
  const {isFeatureEnabled} = useContext(AppContext);
  const [networkType, setNetworkType] = useState<?NetworkType>(null);
  const alertsEnabled = isFeatureEnabled('alerts');

  useEffect(() => {
    if (networkId) {
      MagmaV1API.getNetworksByNetworkIdType({networkId}).then(networkType =>
        setNetworkType(coalesceNetworkType(networkId, networkType)),
      );
    }
  }, [networkId]);

  if (!networkId || networkType === null) {
    return [null, []];
  }

  switch (networkType) {
    case WIFI:
      return getMeshSections(alertsEnabled);
    case SYMPHONY:
    case THIRD_PARTY:
      return getDevicesSections(alertsEnabled);
    case RHINO:
      return getRhinoSections();
    case CWF:
      return getCWFSections();
    case FEG:
      return getFEGSections();
    case LTE:
    default:
      return getLteSections();
  }
}
