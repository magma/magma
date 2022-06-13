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
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DashboardAlertTable from '../../components/DashboardAlertTable';
import DashboardIcon from '@material-ui/icons/Dashboard';
import Dialog from '@material-ui/core/Dialog';
import DialogContent from '@material-ui/core/DialogContent';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DialogTitle from '../../theme/design-system/DialogTitle';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import EventsTable from '../../views/events/EventsTable';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GatewayConfig, {GatewayJsonConfig} from './GatewayDetailConfig';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GatewayConfigYml from './GatewayYMLConfig';
// $FlowFixMe migrated to typescript
import GatewayContext from '../../components/context/GatewayContext';
import GatewayDetailEnodebs from './GatewayDetailEnodebs';
// $FlowFixMe migrated to typescript
import GatewayDetailStatus from './GatewayDetailStatus';
import GatewayDetailSubscribers from './GatewayDetailSubscribers';
import GatewayLogs from './GatewayLogs';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GatewaySummary from './GatewaySummary';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
// $FlowFixMe migrated to typescript
import MenuButton from '../../components/MenuButton';
import MenuItem from '@material-ui/core/MenuItem';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PeopleIcon from '@material-ui/icons/People';
import React, {useContext, useState} from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import withAlert from '../../components/Alert/withAlert';
import {
  GenericCommandControls,
  PingCommandControls,
  TroubleshootingControl,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from '../../components/GatewayCommandFields';
import {Navigate, Route, Routes, useParams} from 'react-router-dom';
// $FlowFixMe migrated to typescript
import {RunGatewayCommands} from '../../state/lte/EquipmentState';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {WithAlert} from '../../components/Alert/withAlert';
import type {lte_gateway} from '../../../generated/MagmaAPIBindings';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    color: colors.primary.white,
  },
}));
type CommandProps = {
  gatewayID: string,
  open: boolean,
  onClose: () => void,
};

function GatewayCommandDialog(props: CommandProps) {
  return (
    <Dialog open={props.open} onClose={props.onClose} scroll="body">
      <DialogTitle onClose={props.onClose} label={'Run Gateway Commands'} />
      <DialogContent>
        <PingCommandControls gatewayID={props.gatewayID} />
        <GenericCommandControls gatewayID={props.gatewayID} />
      </DialogContent>
    </Dialog>
  );
}

const AGGREGATION_TITLE = 'Aggregation';

function TroubleshootingDialog(props: CommandProps) {
  const [tabPos, setTabPos] = useState(false);
  const classes = useStyles();

  return (
    <Dialog
      fullWidth={true}
      maxWidth="md"
      open={props.open}
      onClose={props.onClose}
      scroll="body">
      <DialogTitle onClose={props.onClose} label={'Troubleshoot Gateway'} />
      <Tabs
        value={tabPos}
        onChange={(_, v) => setTabPos(v)}
        indicatorColor="primary"
        className={classes.tabBar}>
        <Tab key={AGGREGATION_TITLE} label={AGGREGATION_TITLE} />
      </Tabs>
      <DialogContent>
        <TroubleshootingControl gatewayID={props.gatewayID} />
      </DialogContent>
    </Dialog>
  );
}

function GatewayMenuInternal(props: WithAlert) {
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const gatewayId: string = nullthrows(params.gatewayId);
  const [gatewayCommandOpen, setGatewayCommandOpen] = useState(false);
  const [troubleshootingDialogOpen, setTroubleshootingDialogOpen] = useState(
    false,
  );
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleGatewayMenuClick = (
    command: 'reboot' | 'ping' | 'restartServices' | 'generic',
    warnMsg: string,
  ) => {
    props.confirm(warnMsg).then(async confirmed => {
      if (!confirmed) {
        return;
      }
      try {
        await RunGatewayCommands({
          networkId,
          gatewayId,
          command: command,
        });
        enqueueSnackbar('command triggered successfully', {
          variant: 'success',
        });
      } catch (e) {
        enqueueSnackbar(e.response?.data?.message ?? e.message, {
          variant: 'error',
        });
      }
    });
  };

  return (
    <div>
      <GatewayCommandDialog
        gatewayID={gatewayId}
        open={gatewayCommandOpen}
        onClose={() => setGatewayCommandOpen(false)}
      />
      <TroubleshootingDialog
        gatewayID={gatewayId}
        open={troubleshootingDialogOpen}
        onClose={() => setTroubleshootingDialogOpen(false)}
      />

      <MenuButton label="Actions" size="small">
        <MenuItem
          data-testid="gatewayReboot"
          onClick={() =>
            handleGatewayMenuClick(
              'reboot',
              `Are you sure you want to reboot ${gatewayId}?`,
            )
          }>
          <Text variant="body2">Reboot</Text>
        </MenuItem>
        <MenuItem
          data-testid="gatewayRestartServices"
          onClick={() =>
            handleGatewayMenuClick(
              'restartServices',
              `Are you sure you want to restart all services on ${gatewayId}?`,
            )
          }>
          <Text variant="body2">Restart Services</Text>
        </MenuItem>
        <MenuItem onClick={() => setGatewayCommandOpen(true)}>
          <Text variant="body2">Command</Text>
        </MenuItem>
        <MenuItem onClick={() => setTroubleshootingDialogOpen(true)}>
          <Text variant="body2">Troubleshoot</Text>
        </MenuItem>
      </MenuButton>
    </div>
  );
}
const GatewayMenu = withAlert(GatewayMenuInternal);

export type GatewayDetailType = {
  gwInfo: lte_gateway,
  refresh: boolean,
};
export function GatewayDetail() {
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);
  const gwCtx = useContext(GatewayContext);

  return (
    <>
      <TopBar
        header={`Equipment/${gatewayId}`}
        tabs={[
          {
            label: 'Overview',
            to: 'overview',
            icon: DashboardIcon,
            filters: <GatewayMenu />,
          },
          {
            label: 'Event',
            to: 'event',
            icon: MyLocationIcon,
            filters: <GatewayMenu />,
          },
          {
            label: 'Logs',
            to: 'logs',
            icon: ListAltIcon,
            filters: <GatewayMenu />,
          },
          {
            label: 'Alerts',
            to: 'alert',
            icon: AccessAlarmIcon,
            filters: <GatewayMenu />,
          },
          {
            label: 'Config',
            to: 'config',
            icon: SettingsIcon,
            filters: <GatewayMenu />,
          },
          {
            label: 'Services Config',
            to: 'services',
            icon: SettingsIcon,
            filters: <GatewayMenu />,
          },
        ]}
      />

      <Routes>
        <Route path="/config/json" element={<GatewayJsonConfig />} />
        <Route path="/config" element={<GatewayConfig />} />
        <Route
          path="/event"
          element={
            <EventsTable
              eventStream="GATEWAY"
              hardwareId={gwCtx.state[gatewayId].device?.hardware_id}
              sz="lg"
              isAutoRefreshing={true}
            />
          }
        />
        <Route
          path="/alert"
          element={
            <DashboardAlertTable labelFilters={{gatewayID: gatewayId}} />
          }
        />
        <Route path="overview" element={<GatewayOverview />} />
        <Route path="/logs" element={<GatewayLogs />} />
        <Route path="/services" element={<GatewayConfigYml />} />
        <Route index element={<Navigate to="overview" replace />} />
      </Routes>
    </>
  );
}

function GatewayOverview() {
  const classes = useStyles();
  const params = useParams();
  const gatewayId: string = nullthrows(params.gatewayId);
  const gwCtx = useContext(GatewayContext);
  const gwInfo = gwCtx.state[gatewayId];
  const [refresh, setRefresh] = useState(true);
  const [refreshEnodebs, setRefreshEnodebs] = useState(false);
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
              <GatewaySummary gwInfo={gwInfo} />
            </Grid>
            <Grid item xs={12} alignItems="center">
              <CardTitleRow icon={MyLocationIcon} label="Events" />
              <EventsTable
                eventStream="GATEWAY"
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
              <GatewayDetailStatus refresh={refresh} />
            </Grid>
            <Grid item>
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label="Connected eNodeBs"
                filter={() => filter(refreshEnodebs, setRefreshEnodebs)}
              />
              <GatewayDetailEnodebs gwInfo={gwInfo} refresh={refreshEnodebs} />
            </Grid>
            <Grid item>
              <CardTitleRow
                icon={PeopleIcon}
                label="Subscribers"
                filter={() => filter(refreshSubscribers, setRefreshSubscribers)}
              />
              <GatewayDetailSubscribers
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

export default GatewayDetail;
