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

import AccessAlarmIcon from '@material-ui/icons/AccessAlarm';
// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import EventsTable from '../../views/events/EventsTable';
import FEGClusterStatus from '../equipment/FEGClusterStatus';
import FEGGatewayConnectionStatus from './FEGGatewayConnectionStatus';
import FEGGatewayContext from '../../components/context/FEGGatewayContext';
import FEGGatewayDetailConfig from './FEGGatewayDetailConfig';
import FEGGatewayDetailStatus from './FEGGatewayDetailStatus';
import FEGGatewayDetailSubscribers from './FEGGatewayDetailSubscribers';
import FEGGatewaySummary from './FEGGatewaySummary';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import Tooltip from '@material-ui/core/Tooltip';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';

import {EVENT_STREAM} from '../../views/events/EventsTable';
import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useContext, useState} from 'react';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
}));

/**
 * Returns the gateway detail page of the federation network.
 * It consists of a gateway overview component to display
 * the useful informations about the federation gateway selected
 * and a top bar to navigate through different pages.
 */
export default function FEGGatewayDetail() {
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);

  return (
    <>
      <TopBar
        header={`Equipment/${gatewayId}`}
        tabs={[
          {
            label: 'Overview',
            to: 'overview',
            icon: DashboardIcon,
          },
          {
            label: 'Event',
            to: 'event',
            icon: MyLocationIcon,
          },
          {
            label: 'Logs',
            to: 'logs',
            icon: ListAltIcon,
          },
          {
            label: 'Alerts',
            to: 'alert',
            icon: AccessAlarmIcon,
          },
          {
            label: 'Config',
            to: 'config',
            icon: SettingsIcon,
          },
        ]}
      />

      <Routes>
        <Route path="overview" element={<FEGGatewayOverview />} />
        <Route path="config" element={<FEGGatewayDetailConfig />} />
        <Route index element={<Navigate to="overview" replace />} />
      </Routes>
    </>
  );
}

/**
 * Returns the gateway information, table of events coming from
 * the gateway, its status, the cluster status of the gateways
 * in the network, and the connected subscribers.
 */
function FEGGatewayOverview() {
  const classes = useStyles();
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);
  const gwCtx = useContext(FEGGatewayContext);
  const gwInfo = gwCtx.state[gatewayId];
  const [refresh, setRefresh] = useState(true);
  const [refreshSubscribers, setRefreshSubscribers] = useState(false);

  const filter = (refresh: boolean, setRefresh) => {
    return (
      <AutorefreshCheckbox
        autorefreshEnabled={refresh}
        onToggle={() => setRefresh(current => !current)}
      />
    );
  };

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4} direction="column">
            <Grid item xs={12} alignItems="center">
              <CardTitleRow icon={CellWifiIcon} label={gatewayId} />
              <FEGGatewaySummary gwInfo={gwInfo} />
            </Grid>
            <Grid item xs={12} alignItems="center">
              <CardTitleRow icon={MyLocationIcon} label="Events" />
              <EventsTable
                eventStream={EVENT_STREAM.GATEWAY}
                hardwareId={gwInfo.device?.hardware_id}
                sz="sm"
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4} direction="column">
            <Grid item>
              <CardTitleRow
                icon={GraphicEqIcon}
                label="Status"
                filter={() => filter(refresh, setRefresh)}
              />
              <FEGGatewayDetailStatus refresh={refresh} />
            </Grid>
            <Grid item>
              <CardTitleRow icon={GraphicEqIcon} label="Connection Status" />
              <FEGGatewayConnectionStatus />
            </Grid>
            <Grid item>
              <FEGClusterStatus />
            </Grid>
            <Grid item>
              <Tooltip
                title="List of subscriber sessions under the networks serviced by
              this federation network">
                <Grid>
                  <CardTitleRow
                    icon={PeopleIcon}
                    label="Managed Subscribers"
                    filter={() =>
                      filter(refreshSubscribers, setRefreshSubscribers)
                    }
                  />
                </Grid>
              </Tooltip>
              <FEGGatewayDetailSubscribers
                gwInfo={gwInfo}
                refresh={refreshSubscribers}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}
