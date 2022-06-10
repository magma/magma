/*
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

import CellWifiIcon from '@material-ui/icons/CellWifi';
import FEGGateway from './FEGEquipmentGateway';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGGatewayDetail from './FEGGatewayDetailMain';
import React from 'react';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';

import {Navigate, Route, Routes} from 'react-router-dom';

/**
 * Returns the full equipment dashboard of the federation network.
 * It consists of an internal equipment dashboard component to display
 * the useful information about the federation gateways.
 */
function FEGEquipmentDashboard() {
  return (
    <>
      <Routes>
        <Route
          path="/overview/gateway/:gatewayId/*"
          element={<FEGGatewayDetail />}
        />
        <Route path="/overview/*" element={<EquipmentDashboardInternal />} />
        <Route index element={<Navigate to="overview" replace />} />
      </Routes>
    </>
  );
}
/**
 * It consists of a top bar to navigate and a federation gateway component
 * to provide information about the federation gateways.
 */
function EquipmentDashboardInternal() {
  return (
    <>
      <TopBar
        header="Equipment"
        tabs={[
          {
            label: 'Federation Gateways',
            to: 'gateway',
            icon: CellWifiIcon,
          },
        ]}
      />
      <Routes>
        <Route path="/gateway" element={<FEGGateway />} />
        <Route index element={<Navigate to="gateway" replace />} />
      </Routes>
    </>
  );
}

export default FEGEquipmentDashboard;
