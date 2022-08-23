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
import AccessAlarmIcon from '@mui/icons-material/AccessAlarm';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@mui/icons-material/CellWifi';
import DashboardAlertTable from '../../components/DashboardAlertTable';
import DashboardIcon from '@mui/icons-material/Dashboard';
import Dialog from '@mui/material/Dialog';
import DialogContent from '@mui/material/DialogContent';
import DialogTitle from '../../theme/design-system/DialogTitle';
import EventsTable from '../../views/events/EventsTable';
import GatewayConfig, {GatewayJsonConfig} from './GatewayDetailConfig';
import GatewayConfigYml from './GatewayYMLConfig';
import GatewayContext from '../../context/GatewayContext';
import GatewayDetailEnodebs from './GatewayDetailEnodebs';
import GatewayDetailStatus from './GatewayDetailStatus';
import GatewayDetailSubscribers from './GatewayDetailSubscribers';
import GatewayLogs from './GatewayLogs';
import GatewaySummary from './GatewaySummary';
import GraphicEqIcon from '@mui/icons-material/GraphicEq';
import Grid from '@mui/material/Grid';
import ListAltIcon from '@mui/icons-material/ListAlt';
import MenuButton from '../../components/MenuButton';
import MenuItem from '@mui/material/MenuItem';
import MyLocationIcon from '@mui/icons-material/MyLocation';
import PeopleIcon from '@mui/icons-material/People';
import React, {useContext, useState} from 'react';
import SettingsIcon from '@mui/icons-material/Settings';
import SettingsInputAntennaIcon from '@mui/icons-material/SettingsInputAntenna';
import Tab from '@mui/material/Tab';
import Tabs from '@mui/material/Tabs';
import Text from '../../theme/design-system/Text';
import TopBar from '../../components/TopBar';
import nullthrows from '../../../shared/util/nullthrows';
import withAlert from '../../components/Alert/withAlert';
import {
  GenericCommandControls,
  PingCommandControls,
  TroubleshootingControl,
} from '../../components/GatewayCommandFields';
import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {RunGatewayCommands} from './RunGatewayCommands';
import {colors} from '../../theme/default';
import {makeStyles} from '@mui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';

import {Theme} from '@mui/material/styles';
import {getErrorMessage} from '../../util/ErrorUtils';
import type {LteGateway} from '../../../generated';
import type {WithAlert} from '../../components/Alert/withAlert';

const useStyles = makeStyles<Theme>(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
  },
}));
type CommandProps = {
  gatewayID: string;
  open: boolean;
  onClose: () => void;
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
        onChange={(_, v) => setTabPos(v as boolean)}
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
    void props.confirm(warnMsg).then(async confirmed => {
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
        enqueueSnackbar(getErrorMessage(e), {
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
  gwInfo: LteGateway;
  refresh: boolean;
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

  const filter = (
    refresh: boolean,
    setRefresh: React.Dispatch<React.SetStateAction<boolean>>,
  ) => {
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
            <Grid item xs={12}>
              <CardTitleRow icon={CellWifiIcon} label={gatewayId} />
              <GatewaySummary gwInfo={gwInfo} />
            </Grid>
            <Grid item xs={12}>
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
