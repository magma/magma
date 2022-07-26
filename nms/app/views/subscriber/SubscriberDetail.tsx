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

import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DataGrid from '../../components/DataGrid';
import EventsTable from '../../views/events/EventsTable';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '../../components/LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberChart from './SubscriberChart';
import SubscriberContext from '../../context/SubscriberContext';
import SubscriberDetailConfig from './SubscriberDetailConfig';
import TopBar from '../../components/TopBar';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../api/useMagmaAPI';

import {Navigate, Route, Routes, useParams} from 'react-router-dom';
import {REFRESH_INTERVAL} from '../../context/AppContext';
import {SubscriberJsonConfig} from './SubscriberDetailConfig';
import {Theme} from '@material-ui/core/styles';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useContext, useState} from 'react';
import {useInterval} from '../../hooks';
import type {DataRows} from '../../components/DataGrid';
import type {Subscriber, SubscriberState} from '../../../generated';

const useStyles = makeStyles<Theme>(theme => ({
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
  const [subscriberConfig, setSubscriberConfig] = useState({} as Subscriber);
  const {isLoading} = useMagmaAPI(
    MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdGet,
    {
      networkId: networkId,
      subscriberId: subscriberId,
    },
    useCallback(
      (response: Subscriber) => {
        setSubscriberConfig(response);

        if (!ctx.state[subscriberId]) {
          void ctx.setState?.('', undefined, {
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
  const [refresh, setRefresh] = useState(false);
  const ctx = useContext(SubscriberContext);
  const subscriberInfo: Subscriber = ctx.state?.[subscriberId];
  useInterval(
    () => ctx.refetchSessionState(subscriberId),
    refresh ? REFRESH_INTERVAL : null,
  );

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

        <Status
          sessionState={ctx.sessionState}
          subscriberInfo={subscriberInfo}
        />
      </Grid>
    </Grid>
  );
}

function Overview() {
  const classes = useStyles();
  const params = useParams();
  const subscriberId: string = nullthrows(params.subscriberId);
  const ctx = useContext(SubscriberContext);

  if (!ctx.state?.[subscriberId]) {
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
          <EventsTable eventStream="SUBSCRIBER" tags={subscriberId} sz="md" />
        </Grid>
      </Grid>
    </div>
  );
}

function Info(props: {subscriberInfo: Subscriber}) {
  const kpiData: Array<DataRows> = [
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
  sessionState?: SubscriberState;
  subscriberInfo: Subscriber;
};

function Status(props: statusProps) {
  const featureUnsupported = 'Unsupported';
  const statusUnknown = 'Unknown';

  const gwId: string =
    props.sessionState?.directory?.location_history?.[0] ?? statusUnknown;

  const kpiData: Array<DataRows> = [
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
