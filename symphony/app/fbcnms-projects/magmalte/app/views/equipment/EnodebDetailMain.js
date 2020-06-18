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
import EnodebConfig from './EnodebDetailConfig';
import EnodebThroughputChart from './EnodebThroughputChart';
import GatewayLogs from './GatewayLogs';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '../../theme/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';

import {colors} from '../../theme/default';
import {EnodebStatus, EnodebSummary} from './EnodebDetailSummaryStatus';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: colors.primary.white,
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
        <Grid container>
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
      <Grid container spacing={3} alignItems="stretch">
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid item xs={6}>
            <Text>
              <SettingsInputAntennaIcon /> {enbInfo.enb.name}
            </Text>
            <EnodebSummary enbInfo={enbInfo} />
          </Grid>

          <Grid item xs={6}>
            <Text>
              <GraphicEqIcon />
              Status
            </Text>
            <EnodebStatus enbInfo={enbInfo} />
          </Grid>
        </Grid>
        <Grid container item spacing={3} alignItems="stretch" xs={12}>
          <Grid item xs={12}>
            {perEnbMetricSupportAvailable ? (
              <EnodebThroughputChart
                title={CHART_TITLE}
                queries={[
                  `sum(pdcp_user_plane_bytes_dl{service="enodebd"})/1000`,
                  `sum(pdcp_user_plane_bytes_ul{service="enodebd"})/1000`,
                ]}
                legendLabels={['Download', 'Upload']}
              />
            ) : (
              <Paper className={classes.paper}>
                Enodeb Throughput Chart Currently Unavailable
              </Paper>
            )}
          </Grid>
        </Grid>
        <Grid container spacing={3} alignItems="stretch" item xs={12}>
          <Grid item xs={6}>
            <Text>
              <MyLocationIcon /> Events
            </Text>
            <Paper className={classes.paper}>Event Information</Paper>
          </Grid>

          <Grid item xs={6}>
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

export default EnodebDetail;
