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
import type {DataRows} from '../../components/DataGrid';
import type {subscriber} from '@fbcnms/magma-api';

import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DataGrid from '../../components/DataGrid';
import EventsTable from '../../views/events/EventsTable';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberChart from './SubscriberChart';
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetailConfig from './SubscriberDetailConfig';
import TopBar from '../../components/TopBar';
import nullthrows from '@fbcnms/util/nullthrows';

import {
  REFRESH_INTERVAL,
  useRefreshingContext,
} from '../../components/context/RefreshContext';
import {Redirect, Route, Switch} from 'react-router-dom';
import {SubscriberJsonConfig} from './SubscriberDetailConfig';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
    color: colors.primary.white,
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: colors.primary.white,
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '16px 0 16px 0',
    display: 'flex',
    alignItems: 'center',
  },
  tabIconLabel: {
    marginRight: '8px',
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
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  // TODO: Remove this once event table has been added
  paper: {
    textAlign: 'center',
    padding: theme.spacing(10),
  },
}));

export default function SubscriberDetail() {
  const {relativePath, relativeUrl, match} = useRouter();
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const ctx = useContext(SubscriberContext);
  // TODO: render a "Not found" component if the IMSI is not found
  const subscriberInfo = ctx.state?.[subscriberId] || {};

  return (
    <>
      <TopBar
        header={`Equipment/${subscriberInfo.name ?? subscriberId}`}
        tabs={
          !Object.keys(subscriberInfo).length
            ? [
                {
                  label: 'Event',
                  to: '/event',
                  icon: MyLocationIcon,
                },
              ]
            : [
                {
                  label: 'Overview',
                  to: '/overview',
                  icon: DashboardIcon,
                },
                {
                  label: 'Event',
                  to: '/event',
                  icon: MyLocationIcon,
                },
                {
                  label: 'Config',
                  to: '/config',
                  icon: SettingsIcon,
                },
              ]
        }
      />

      <Switch>
        <Route
          path={relativePath('/config/json')}
          render={() => <SubscriberJsonConfig />}
        />
        <Route
          path={relativePath('/config')}
          render={() => <SubscriberDetailConfig />}
        />
        <Route path={relativePath('/overview')} render={() => <Overview />} />
        <Route
          path={relativePath('/event')}
          render={() => (
            <EventsTable
              sz="lg"
              eventStream="SUBSCRIBER"
              isAutoRefreshing={true}
              tags={subscriberId}
            />
          )}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}
function StatusInfo() {
  const {match} = useRouter();
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const networkId: string = nullthrows(match.params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();

  const ctx = useRefreshingContext({
    context: SubscriberContext,
    networkId: networkId,
    type: 'subscriber',
    interval: REFRESH_INTERVAL,
    enqueueSnackbar: enqueueSnackbar,
    refresh: true,
    id: subscriberId,
  });
  const subscriberInfo = ctx.state?.[subscriberId];

  return (
    <Grid container spacing={4}>
      <Grid item xs={12} md={6}>
        <CardTitleRow icon={PersonIcon} label="Subscriber" />
        <Info subscriberInfo={subscriberInfo} />
      </Grid>
      <Grid item xs={12} md={6}>
        <CardTitleRow icon={GraphicEqIcon} label="Status" />
        <Status subscriberInfo={subscriberInfo} />
      </Grid>
    </Grid>
  );
}

function Overview() {
  const classes = useStyles();
  const {match} = useRouter();
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const ctx = useContext(SubscriberContext);
  const subscriberInfo = ctx.state?.[subscriberId];
  if (!subscriberInfo) {
    return null;
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <StatusInfo />
        </Grid>
        <Grid item xs={12}>
          <SubscriberChart />
        </Grid>
        <Grid item xs={12}>
          <EventsTable
            eventStream="SUBSCRIBER"
            tags={subscriberInfo.id}
            sz="md"
          />
        </Grid>
      </Grid>
    </div>
  );
}

function Info(props: {subscriberInfo: subscriber}) {
  const kpiData: DataRows[] = [
    [
      {
        value: props.subscriberInfo.name ?? props.subscriberInfo.id,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'IMSI',
        value: props.subscriberInfo.id,
        statusCircle: false,
      },
      {
        category: 'Service',
        value: props.subscriberInfo.lte.state,
        statusCircle: true,
        status: props.subscriberInfo.lte.state === 'ACTIVE',
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}

function Status({subscriberInfo}: {subscriberInfo: subscriber}) {
  const featureUnsupported = 'Unsupported';
  const statusUnknown = 'Unknown';
  const gwId =
    subscriberInfo?.state?.directory?.location_history?.[0] ?? statusUnknown;

  const kpiData: DataRows[] = [
    [
      {
        category: 'Gateway ID',
        value: gwId,
        statusCircle: false,
        tooltip: 'latest gateway connected to the subscriber',
      },
      {
        category: 'eNodeB SN',
        value: featureUnsupported,
        statusCircle: false,
        tooltip: 'not supported',
      },
    ],
    [
      {
        category: 'Connection Status',
        value: statusUnknown,
        statusCircle: false,
        tooltip: 'not supported',
      },
      {
        category: 'UE Latency',
        value: subscriberInfo.monitoring?.icmp?.latency_ms ?? statusUnknown,
        unit: subscriberInfo.monitoring?.icmp?.latency_ms ? 'ms' : '',
        statusCircle: false,
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}
