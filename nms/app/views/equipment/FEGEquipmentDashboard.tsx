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
 */

import CellWifiIcon from '@mui/icons-material/CellWifi';
import FEGClusterStatus from './FEGClusterStatus';
import FEGEquipmentGatewayKPIs from './FEGEquipmentGatewayKPIs';
import FEGGatewayDetail from './FEGGatewayDetailMain';
import FEGGatewayTable from './FEGGatewayTable';
import GatewayCheckinChart from './GatewayCheckinChart';
import Grid from '@mui/material/Grid';
import Paper from '@mui/material/Paper';
import React from 'react';
import TopBar from '../../components/TopBar';
import UpgradeButton from './UpgradeTiersDialog';

import {FEGAddGatewayButton} from '../../components/feg/FEGGatewayDialog';
import {Navigate, Route, Routes} from 'react-router-dom';
import {Theme} from '@mui/material/styles';
import {makeStyles} from '@mui/styles';

const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
  },
}));

function FEGEquipmentGateways() {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justifyContent="space-between" spacing={4}>
        <Grid item xs={12}>
          <GatewayCheckinChart />
        </Grid>
        <Grid item xs={12}>
          <Paper elevation={0}>
            <FEGEquipmentGatewayKPIs />
          </Paper>
        </Grid>
        <Grid item xs={12}>
          <FEGClusterStatus />
        </Grid>
        <Grid item xs={12}>
          <FEGGatewayTable />
        </Grid>
      </Grid>
    </div>
  );
}

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
                  <FEGAddGatewayButton />
                </Grid>
              </Grid>
            ),
          },
        ]}
      />
      <Routes>
        <Route path="/gateway" element={<FEGEquipmentGateways />} />
        <Route index element={<Navigate to="gateway" replace />} />
      </Routes>
    </>
  );
}

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

export default FEGEquipmentDashboard;
