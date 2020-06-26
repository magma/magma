/*
 * Copyright (c) Facebook, Inc. and its affiliates.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {EnodebInfo} from '../../components/lte/EnodebUtils';

import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DateTimeMetricChart from '../../components/DateTimeMetricChart';
import EnodebConfig from './EnodebDetailConfig';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import nullthrows from '@fbcnms/util/nullthrows';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';

import {CardTitleRow} from '../../components/layout/CardTitleRow';
import {colors, typography} from '../../theme/default';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {makeStyles} from '@material-ui/styles';
import {Redirect, Route, Switch} from 'react-router-dom';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
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
  appBarBtnSecondary: {
    color: colors.primary.white,
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    textAlign: 'center',
    padding: theme.spacing(10),
  },
  card: {
    variant: 'outlined',
  },
}));
const CHART_TITLE = 'Bandwidth Usage';

export function EnodebDetail({enbInfo}: {enbInfo: {[string]: EnodebInfo}}) {
  const classes = useStyles();
  const [tabPos, setTabPos] = React.useState(0);
  const {relativePath, relativeUrl, match} = useRouter();
  const enodebSerial: string = nullthrows(match.params.enodebSerial);
  return (
    <>
      <div className={classes.topBar}>
        <Text variant="body2">Equipment/{enbInfo[enodebSerial].enb.name}</Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
          <Grid item xs={6}>
            <Tabs
              value={tabPos}
              onChange={(event, v) => setTabPos(v)}
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
          <Grid
            item
            xs={6}
            direction="row"
            justify="flex-end"
            alignItems="center">
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                {/* TODO: these button styles need to be localized */}
                <Button variant="text" className={classes.appBarBtnSecondary}>
                  Secondary Action
                </Button>
              </Grid>
              <Grid item>
                <Button variant="contained" className={classes.appBarBtn}>
                  Reboot
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>
      <Switch>
        <Route
          path={relativePath('/overview')}
          render={() => <Overview enbInfo={enbInfo[enodebSerial]} />}
        />
        <Route
          path={relativePath('/config')}
          render={() => <EnodebConfig enbInfo={enbInfo[enodebSerial]} />}
        />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function Overview({enbInfo}: {enbInfo: EnodebInfo}) {
  const classes = useStyles();
  const perEnbMetricSupportAvailable = false;
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={12} md={6} alignItems="center">
              <CardTitleRow
                icon={SettingsInputAntennaIcon}
                label={enbInfo.enb.name}
              />
              <EnodebSummary enbInfo={enbInfo} />
            </Grid>

            <Grid item xs={6}>
              <CardTitleRow icon={GraphicEqIcon} label="Status" />
              <EnodebStatus enbInfo={enbInfo} />
            </Grid>
          </Grid>
        </Grid>
        <Grid item xs={12}>
          {perEnbMetricSupportAvailable ? (
            <DateTimeMetricChart
              title={CHART_TITLE}
              queries={[
                `sum(pdcp_user_plane_bytes_dl{service="enodebd"})/1000`,
                `sum(pdcp_user_plane_bytes_ul{service="enodebd"})/1000`,
              ]}
              legendLabels={['Download', 'Upload']}
            />
          ) : (
            <Paper className={classes.paper} elevation={0}>
              <Text variant="body2">
                Enodeb Throughput Chart Currently Unavailable
              </Text>
            </Paper>
          )}
        </Grid>
        <Grid item xs={12}>
          <Grid container spacing={4}>
            <Grid item xs={6}>
              <CardTitleRow icon={MyLocationIcon} label="Events" />
              <Paper className={classes.paper} elevation={0}>
                <Text variant="body2">Event Information</Text>
              </Paper>
            </Grid>

            <Grid item xs={6}>
              <CardTitleRow icon={PeopleIcon} label="Subscribers" />
              <Paper className={classes.paper} elevation={0}>
                <Text variant="body2">Subscribers data</Text>
              </Paper>
            </Grid>
          </Grid>
        </Grid>
      </Grid>
    </div>
  );
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

export default EnodebDetail;
