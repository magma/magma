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
import type {lte_gateway} from '@fbcnms/magma-api';

import AccessAlarmIcon from '@material-ui/icons/AccessAlarm';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import CellWifiIcon from '@material-ui/icons/CellWifi';
import DashboardIcon from '@material-ui/icons/Dashboard';
import GatewayLogs from './GatewayLogs';
import GatewaySummary from './GatewaySummary';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import ListAltIcon from '@material-ui/icons/ListAlt';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';

import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: 'white',
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
}));

export function GatewayDetail({
  lteGateways,
}: {
  lteGateways: {[string]: lte_gateway},
}) {
  const classes = useStyles();
  const [tabPos, setTabPos] = React.useState(0);
  const {relativePath, relativeUrl, match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const gwInfo = lteGateways[gatewayId];
  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Equipment/{gatewayId}
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container>
          <Grid item xs={6}>
            <Tabs
              value={tabPos}
              onChange={(_, v) => setTabPos(v)}
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
                key="Log"
                component={NestedRouteLink}
                label={<LogTabLabel />}
                to="/logs"
                className={classes.tab}
              />
              <Tab
                key="Alert"
                component={NestedRouteLink}
                label={<AlertTabLabel />}
                to="/alert"
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
          <Grid item xs={6}>
            <Grid container justify="flex-end" alignItems="center" spacing={2}>
              <Grid item>
                <Text color="light">Secondary Action</Text>
              </Grid>
              <Grid item>
                <Button color="primary" variant="contained">
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
          render={() => <GatewayOverview gwInfo={gwInfo} />}
        />
        <Route path={relativePath('/logs')} component={GatewayLogs} />
        <Redirect to={relativeUrl('/overview')} />
      </Switch>
    </>
  );
}

function GatewayOverview({gwInfo}: {gwInfo: lte_gateway}) {
  const classes = useStyles();
  const {match} = useRouter();
  const gatewayId: string = nullthrows(match.params.gatewayId);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={6}>
          <Grid item xs={12}>
            <Text>
              <CellWifiIcon /> {gatewayId}
            </Text>
            <GatewaySummary gwInfo={gwInfo} />
          </Grid>
          <Grid item xs={12}>
            <Text>
              <MyLocationIcon /> Events
            </Text>
            <Paper className={classes.paper}>Event Information</Paper>
          </Grid>
        </Grid>

        <Grid container spacing={3} alignItems="stretch" item xs={6}>
          <Grid item xs={12}>
            <Text>
              <GraphicEqIcon />
              Status
            </Text>
            <Paper className={classes.paper}>Status KPI Tray</Paper>
          </Grid>
          <Grid item xs={12}>
            <Text>
              <SettingsInputAntennaIcon /> Connected eNodeBs
            </Text>
            <Paper className={classes.paper}>
              Connected eNodeB Information
            </Paper>
          </Grid>
          <Grid item xs={12}>
            <Text>
              <PeopleIcon /> Subscribers
            </Text>
            <Paper className={classes.paper}>Subscribers data</Paper>
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

function LogTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <ListAltIcon className={classes.tabIconLabel} /> Logs
    </div>
  );
}

function AlertTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <AccessAlarmIcon className={classes.tabIconLabel} /> Alerts
    </div>
  );
}

export default GatewayDetail;
