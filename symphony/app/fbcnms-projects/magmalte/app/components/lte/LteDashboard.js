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

import AppBar from '@material-ui/core/AppBar';
import DashboardAlertTable from '../DashboardAlertTable';
import DashboardKPIs from '../DashboardKPIs';
import EnodebKPIs from '../EnodebKPIs';
import EventAlertChart from '../EventAlertChart';
import GatewayKPIs from '../GatewayKPIs';
import Grid from '@material-ui/core/Grid';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import {DateTimePicker} from '@material-ui/pickers';
import {GpsFixed, NetworkCheck, People} from '@material-ui/icons';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: `0 ${theme.spacing(5)}px`,
  },
  tabs: {
    color: 'white',
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
  input: {
    color: 'white',
    backgroundColor: '#545F77',
    border: 'none',
    borderRadius: '4px',
    textAlign: 'center',
    padding: `${theme.spacing(1)}px 0`,
  },
  cardTitle: {
    marginBottom: theme.spacing(1),
  },
  cardTitleIcon: {
    marginRight: theme.spacing(1),
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
}));

function LteDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();

  // datetime picker
  const [startDate, setStartDate] = useState(moment().subtract(3, 'days'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Dashboard
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container direction="row" justify="flex-end" alignItems="center">
          <Grid item xs={6}>
            <Tabs
              value={0}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Network"
                component={NestedRouteLink}
                label={<DashboardTabLabel label="Network" />}
                to="/network"
                className={classes.tab}
              />
              <Tab
                key="Subscribers"
                component={NestedRouteLink}
                label={<DashboardTabLabel label="Subscribers" />}
                to="#"
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
                <Text color="light">Filter By Date</Text>
              </Grid>
              <DateTimePicker
                autoOk
                inputVariant="outlined"
                maxDate={endDate}
                disableFuture
                value={startDate}
                inputProps={{className: classes.input}}
                onChange={setStartDate}
              />
              <Grid item>
                <Text color="light">to</Text>
              </Grid>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                disableFuture
                value={endDate}
                inputProps={{className: classes.input}}
                onChange={setEndDate}
              />
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/network')}
          render={props => (
            <LteNetworkDashboard {...props} startEnd={[startDate, endDate]} />
          )}
        />
        <Redirect to={relativeUrl('/network')} />
      </Switch>
    </>
  );
}

function LteNetworkDashboard({startEnd}: {startEnd: [moment, moment]}) {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <EventAlertChart startEnd={startEnd} />
        </Grid>

        <Grid item xs={12}>
          <DashboardAlertTable />
        </Grid>

        <Grid item xs={12}>
          <DashboardKPIs />
        </Grid>

        {/* <Grid item xs={6}>
          <Paper>
            <GatewayKPIs />
          </Paper>
        </Grid>

        <Grid item xs={6}>
          <Paper elevation={2}>
            <EnodebKPIs />
          </Paper>
        </Grid>

        <Grid item xs={12}>
          <Text>
            <GpsFixed /> Events (388)
          </Text>
          <Paper>
            <div className={classes.contentPlaceholder}>
              Events Table Goes Here
            </div>
          </Paper>
        </Grid>*/}
      </Grid>
    </div>
  );
}

function DashboardTabLabel(props) {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      {props.label === 'Subscribers' ? (
        <People className={classes.tabIconLabel} />
      ) : props.label === 'Network' ? (
        <NetworkCheck className={classes.tabIconLabel} />
      ) : null}
      {props.label}
    </div>
  );
}

export default LteDashboard;
