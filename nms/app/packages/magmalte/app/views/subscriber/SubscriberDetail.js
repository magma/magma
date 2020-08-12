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

import AppBar from '@material-ui/core/AppBar';
import CardTitleRow from '../../components/layout/CardTitleRow';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DataGrid from '../../components/DataGrid';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EventsTable from '../../views/events/EventsTable';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SubscriberDetailConfig from './SubscriberDetailConfig';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import {SubscriberJsonConfig} from './SubscriberDetailConfig';

import {DetailTabItems, GetCurrentTabPos} from '../../components/TabUtils.js';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const CHART_TITLE = 'Data Usage';
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

export default function SubscriberDetail(props: {
  subscriberMap: ?{[string]: subscriber},
}) {
  const classes = useStyles();
  const {relativePath, relativeUrl, match} = useRouter();
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const subscriberInfo = props.subscriberMap?.[subscriberId];
  if (!subscriberInfo) {
    return null;
  }

  return (
    <>
      <div className={classes.topBar}>
        <Text variant="body2">Subscriber/{subscriberId}</Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
          <Grid item xs={12}>
            <Tabs
              value={GetCurrentTabPos(match.url, DetailTabItems)}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Overview"
                component={NestedRouteLink}
                label={<OverviewTabLabel />}
                to="/overview"
                className={classes.tab}
              />
              <Tab
                key="Event"
                component={NestedRouteLink}
                label={<EventTabLabel />}
                to="/event"
                className={classes.tab}
              />
              <Tab
                key="Config"
                component={NestedRouteLink}
                label={<ConfigTabLabel />}
                to="/config"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
        </Grid>
      </AppBar>
      <Switch>
        <Route
          path={relativePath('/config/json')}
          render={() => <SubscriberJsonConfig />}
        />
        <Route
          path={relativePath('/config')}
          render={() => (
            <SubscriberDetailConfig subscriberInfo={subscriberInfo} />
          )}
        />
        <Route
          path={relativePath('/overview')}
          render={() => <Overview subscriberInfo={subscriberInfo} />}
        />
        <Route
          path={relativePath('/event')}
          render={() => (
            <EventsTable
              sz="lg"
              eventStream="SUBSCRIBER"
              tags={subscriberInfo.id}
            />
          )}
        />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function Overview(props: {subscriberInfo: subscriber}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow icon={PersonIcon} label="Subscriber" />
              <Info subscriberInfo={props.subscriberInfo} />
            </Grid>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <Status subscriberInfo={props.subscriberInfo} />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <DateTimeMetricChart
            title={CHART_TITLE}
            queries={[
              `ue_traffic{IMSI="${props.subscriberInfo.id}",direction="down"}`,
              `ue_traffic{IMSI="${props.subscriberInfo.id}",direction="up"}`,
            ]}
            legendLabels={['Download', 'Upload']}
          />
        </Grid>
        <Grid item xs={12}>
          <EventsTable
            eventStream="SUBSCRIBER"
            tags={props.subscriberInfo.id}
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
        value: props.subscriberInfo.id,
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

function Status() {
  const featureUnsupported = 'Unsupported';
  const statusUnknown = 'Unknown';

  const kpiData: DataRows[] = [
    [
      {
        category: 'Gateway ID',
        value: featureUnsupported,
        statusCircle: false,
      },
      {
        category: 'eNodeB SN',
        value: featureUnsupported,
        statusCircle: false,
      },
    ],
    [
      {
        category: 'Connection Status',
        value: statusUnknown,
        statusCircle: false,
      },
      {
        category: 'UE Latency',
        value: statusUnknown,
        statusCircle: false,
      },
    ],
  ];

  return <DataGrid data={kpiData} />;
}

function OverviewTabLabel() {
  const classes = useStyles();
  return (
    <div className={classes.tabLabel}>
      <DashboardIcon className={classes.tabIconLabel} /> Overview
    </div>
  );
}

function ConfigTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <SettingsIcon className={classes.tabIconLabel} /> Config
    </div>
  );
}

function EventTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <MyLocationIcon className={classes.tabIconLabel} /> Event
    </div>
  );
}
