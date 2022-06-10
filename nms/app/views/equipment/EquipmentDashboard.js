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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AddEditEnodeButton from './EnodebDetailConfigEdit';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import AddEditGatewayPoolButton from './GatewayPoolEdit';
import Cbsds from '../domain-proxy/Cbsds';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodeb from './EquipmentEnodeb';
import EnodebDetail from './EnodebDetailMain';
import Gateway from './EquipmentGateway';
import GatewayDetail from './GatewayDetailMain';
import GatewayPools from './EquipmentGatewayPools';
import Grid from '@material-ui/core/Grid';
import GroupWorkIcon from '@material-ui/icons/GroupWork';
import RadioIcon from '@material-ui/icons/Radio';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
import UpgradeButton from './UpgradeTiersDialog';
import {AddEditCbsdButton} from '../domain-proxy/CbsdEdit';

import {Navigate, Route, Routes} from 'react-router-dom';

function EquipmentDashboard() {
  return (
    <Routes>
      <Route
        path="/overview/gateway/:gatewayId/*"
        element={<GatewayDetail />}
      />
      <Route
        path="/overview/enodeb/:enodebSerial/*"
        element={<EnodebDetail />}
      />
      <Route path="/overview/*" element={<EquipmentDashboardInternal />} />
      <Route index element={<Navigate to="overview" replace />} />
    </Routes>
  );
}

function EquipmentDashboardInternal() {
  return (
    <>
      <TopBar
        header="Equipment"
        tabs={[
          {
            label: 'Gateways',
            to: 'gateway',
            icon: CellWifiIcon,
            filters: (
              <Grid
                container
                justifyContent="flex-end"
                alignItems="center"
                spacing={2}>
                <Grid item>
                  <UpgradeButton />
                </Grid>
                <Grid item>
                  <AddEditGatewayButton title="Add New" isLink={false} />
                </Grid>
              </Grid>
            ),
          },
          {
            label: 'eNodeB',
            to: 'enodeb',
            icon: SettingsInputAntennaIcon,
            filters: <AddEditEnodeButton title="Add New" isLink={false} />,
          },
          {
            label: 'CBSDs',
            to: 'cbsds',
            icon: RadioIcon,
            filters: <AddEditCbsdButton title="Add New" isLink={false} />,
          },
          {
            label: 'Gateway Pools',
            to: 'pools',
            icon: GroupWorkIcon,
            filters: (
              <AddEditGatewayPoolButton title="Add New" isLink={false} />
            ),
          },
        ]}
      />
      <Routes>
        <Route path="/gateway" element={<Gateway />} />
        <Route path="/enodeb" element={<Enodeb />} />
        <Route path="/cbsds" element={<Cbsds />} />
        <Route path="/pools/*" element={<GatewayPools />} />
        <Route index element={<Navigate to="gateway" replace />} />
      </Routes>
    </>
  );
}

export default EquipmentDashboard;
