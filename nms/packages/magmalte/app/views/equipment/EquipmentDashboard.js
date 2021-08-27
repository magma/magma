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

import AddEditEnodeButton from './EnodebDetailConfigEdit';
import AddEditGatewayButton from './GatewayDetailConfigEdit';
import AddEditGatewayPoolButton from './GatewayPoolEdit';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import Enodeb from './EquipmentEnodeb';
import EnodebDetail from './EnodebDetailMain';
import Gateway from './EquipmentGateway';
import GatewayDetail from './GatewayDetailMain';
import GatewayPools from './EquipmentGatewayPools';
import Grid from '@material-ui/core/Grid';
import GroupWorkIcon from '@material-ui/icons/GroupWork';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import TopBar from '../../components/TopBar';
import UpgradeButton from './UpgradeTiersDialog';

import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

function EquipmentDashboard() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <Switch>
        <Route
          path={relativePath('/overview/gateway/:gatewayId')}
          component={GatewayDetail}
        />
        <Route
          path={relativePath('/overview/enodeb/:enodebSerial')}
          component={EnodebDetail}
        />
        <Route
          path={relativePath('/overview')}
          component={EquipmentDashboardInternal}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function EquipmentDashboardInternal() {
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <TopBar
        header="Equipment"
        tabs={[
          {
            label: 'Gateways',
            to: '/gateway',
            icon: CellWifiIcon,
            filters: (
              <Grid
                container
                justify="flex-end"
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
            to: '/enodeb',
            icon: SettingsInputAntennaIcon,
            filters: <AddEditEnodeButton title="Add New" isLink={false} />,
          },
          {
            label: 'Gateway Pools',
            to: '/pools',
            icon: GroupWorkIcon,
            filters: (
              <AddEditGatewayPoolButton title="Add New" isLink={false} />
            ),
          },
        ]}
      />
      <Switch>
        <Route path={relativePath('/gateway')} component={Gateway} />
        <Route path={relativePath('/enodeb')} component={Enodeb} />
        <Route path={relativePath('/pools')} component={GatewayPools} />
        <Redirect to={relativeUrl('/gateway')} />
      </Switch>
    </>
  );
}

export default EquipmentDashboard;
