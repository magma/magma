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
import type {WithAlert} from '@fbcnms/ui/components/Alert/withAlert';

import AccessAlarmIcon from '@material-ui/icons/AccessAlarm';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardAlertTable from '../../components/DashboardAlertTable';
import DashboardIcon from '@material-ui/icons/Dashboard';
import EventsTable from '../../views/events/EventsTable';
import GatewayConfig from './GatewayDetailConfig';
import GatewayContext from '../../components/context/GatewayContext';
import GatewayDetailEnodebs from './GatewayDetailEnodebs';
import GatewayDetailStatus from './GatewayDetailStatus';
import GatewayDetailSubscribers from './GatewayDetailSubscribers';
import GatewayLogs from './GatewayLogs';
import GatewaySummary from './GatewaySummary';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';
import withAlert from '@fbcnms/ui/components/Alert/withAlert';

import {GatewayJsonConfig} from './GatewayDetailConfig';
import {Redirect, Route, Switch} from 'react-router-dom';
import {RunGatewayCommands} from '../../state/lte/EquipmentState';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
  paper: {
    textAlign: 'center',
    padding: theme.spacing(10),
  },
}));

function GatewayRebootButtonInternal(props: WithAlert) {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const enqueueSnackbar = useEnqueueSnackbar();

  const handleClick = () => {
    props
      .confirm(`Are you sure you want to reboot ${gatewayId}?`)
      .then(async confirmed => {
        if (!confirmed) {
          return;
        }
        try {
          await RunGatewayCommands({networkId, gatewayId, command: 'reboot'});
          enqueueSnackbar('gateway reboot triggered successfully', {
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
    <Button
      variant="contained"
      className={classes.appBarBtn}
      onClick={handleClick}>
      Reboot
    </Button>
  );
}
const GatewayRebootButton = withAlert(GatewayRebootButtonInternal);

export function GatewayDetail() {
  const {relativePath, relativeUrl, match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const gwCtx = useContext(GatewayContext);

  return (
    <>
      <TopBar
        header={`Equipment/${gatewayId}`}
        tabs={[
          {
            label: 'Overview',
            to: '/overview',
            icon: DashboardIcon,
            filters: <GatewayRebootButton />,
          },
          {
            label: 'Event',
            to: '/event',
            icon: MyLocationIcon,
            filters: <GatewayRebootButton />,
          },
          {
            label: 'Logs',
            to: '/logs',
            icon: ListAltIcon,
            filters: <GatewayRebootButton />,
          },
          {
            label: 'Alerts',
            to: '/alert',
            icon: AccessAlarmIcon,
            filters: <GatewayRebootButton />,
          },
          {
            label: 'Config',
            to: '/config',
            icon: SettingsIcon,
            filters: <GatewayRebootButton />,
          },
        ]}
      />

      <Switch>
        <Route
          path={relativePath('/config/json')}
          component={GatewayJsonConfig}
        />
        <Route path={relativePath('/config')} component={GatewayConfig} />
        <Route
          path={relativePath('/event')}
          render={() => (
            <EventsTable
              eventStream="GATEWAY"
              tags={gwCtx.state[gatewayId].device.hardware_id}
              sz="lg"
            />
          )}
        />
        <Route
          path={relativePath('/alert')}
          render={() => (
            <DashboardAlertTable labelFilters={{gatewayID: gatewayId}} />
          )}
        />
        <Route path={relativePath('/overview')} component={GatewayOverview} />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function GatewayOverview() {
  const classes = useStyles();
  const {match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const ctx = useContext(GatewayContext);
  const gwInfo = ctx.state[gatewayId];

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
                tags={gwInfo.device.hardware_id}
                sz="sm"
              />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12} md={6}>
          <Grid container spacing={4} direction="column">
            <Grid item>
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <GatewayDetailStatus gwInfo={gwInfo} />
            </Grid>
            <Grid item>
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label="Connected eNodeBs"
              />
              <GatewayDetailEnodebs gwInfo={gwInfo} />
            </Grid>
            <Grid item>
              <CardTitleRow icon={PeopleIcon} label="Subscribers" />
              <GatewayDetailSubscribers gwInfo={gwInfo} />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
}

export default GatewayDetail;
