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
import type {DataRows} from '../../components/DataGrid';
import type {
  subscriber,
  subscriber_state,
} from '../../../generated/MagmaAPIBindings';

// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import DataGrid from '../../components/DataGrid';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import EventsTable from '../../views/events/EventsTable';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
// $FlowFixMe migrated to typescript
import LoadingFiller from '../../components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import SubscriberChart from './SubscriberChart';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetailConfig from './SubscriberDetailConfig';
// $FlowFixMe migrated to typescript
import TopBar from '../../components/TopBar';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';

import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {
  REFRESH_INTERVAL,
  useRefreshingContext,
  // $FlowFixMe[cannot-resolve-module] for TypeScript migration
} from '../../components/context/RefreshContext';
import {SubscriberJsonConfig} from './SubscriberDetailConfig';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useState} from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';

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
  const params = useParams();
  const subscriberId: string = nullthrows(params.subscriberId);
  const networkId: string = nullthrows(params.networkId);
  const ctx = useContext(SubscriberContext);
  const [subscriberConfig, setSubscriberConfig] = useState<subscriber>({});
  const {isLoading, response: _subscriberResponse} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribersBySubscriberId,
    {
      networkId: networkId,
      subscriberId: subscriberId,
    },
    useCallback(
      response => {
        setSubscriberConfig(response);
        if (!ctx.state[subscriberId]) {
          ctx.setState?.('', undefined, {
            ...ctx.state,
            [subscriberId]: response,
          });
        }
      },
      [ctx, subscriberId],
    ),
  );
  if (isLoading) {
    return <LoadingFiller />;
  }

  const subscriberInfo = ctx.state?.[subscriberId] || subscriberConfig;
  return (
    <>
      <TopBar
        header={`Subscriber/${subscriberInfo.name ?? subscriberId}`}
        tabs={
          !Object.keys(subscriberInfo).length
            ? [
                {
                  label: 'Event',
                  to: 'event',
                  icon: MyLocationIcon,
                },
              ]
            : [
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
                  label: 'Config',
                  to: 'config',
                  icon: SettingsIcon,
                },
              ]
        }
      />

      <Routes>
        <Route path="/config/json" element={<SubscriberJsonConfig />} />
        <Route path="/config" element={<SubscriberDetailConfig />} />
        <Route path="/overview" element={<Overview />} />
        <Route
          path="/event"
          element={
            <EventsTable
              sz="lg"
              eventStream="SUBSCRIBER"
              isAutoRefreshing={true}
              tags={subscriberId}
            />
          }
        />
        <Route index element={<Navigate to="overview" replace />} />
      </Routes>
    </>
  );
}
function StatusInfo() {
  const params = useParams();
  const subscriberId: string = nullthrows(params.subscriberId);
  const enqueueSnackbar = useEnqueueSnackbar();
  const [refresh, setRefresh] = useState(false);
  const ctx = useContext(SubscriberContext);
  // $FlowIgnore
  const subscriberInfo: subscriber = ctx.state?.[subscriberId];
  const networkId: string = nullthrows(params.networkId);
  const refreshingSessionState = useRefreshingContext({
    context: SubscriberContext,
    networkId,
    type: 'subscriber',
    interval: REFRESH_INTERVAL,
    enqueueSnackbar,
    refresh,
    id: subscriberId,
  });
  // $FlowIgnore
  const sessions: subscriber_state = refreshingSessionState.sessionState;
  function refreshFilter() {
    return (
      <AutorefreshCheckbox
        autorefreshEnabled={refresh}
        onToggle={() => setRefresh(current => !current)}
      />
    );
  }
  return (
    <Grid container spacing={4}>
      <Grid item xs={12} md={6}>
        <CardTitleRow icon={PersonIcon} label="Subscriber" />
        <Info subscriberInfo={subscriberInfo} />
      </Grid>
      <Grid item xs={12} md={6}>
        <CardTitleRow
          icon={GraphicEqIcon}
          label="Status"
          filter={() => refreshFilter()}
        />

        <Status sessionState={sessions} subscriberInfo={subscriberInfo} />
      </Grid>
    </Grid>
  );
}

function Overview() {
  const classes = useStyles();
  const params = useParams();
  const subscriberId: string = nullthrows(params.subscriberId);
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
type statusProps = {
  sessionState?: subscriber_state,
  subscriberInfo: subscriber,
};
function Status(props: statusProps) {
  const featureUnsupported = 'Unsupported';
  const statusUnknown = 'Unknown';

  const gwId =
    // $FlowIgnore
    props.sessionsState?.directory?.location_history?.[0] ?? statusUnknown;

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
        value:
          props.subscriberInfo.monitoring?.icmp?.latency_ms ?? statusUnknown,
        unit: props.subscriberInfo.monitoring?.icmp?.latency_ms ? 'ms' : '',
        statusCircle: false,
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}
