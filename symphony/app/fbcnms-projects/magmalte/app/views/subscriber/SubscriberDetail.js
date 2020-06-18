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
import type {subscriber} from '../../../../../fbcnms-packages/fbcnms-magma-api';

import AppBar from '@material-ui/core/AppBar';
import Card from '@material-ui/core/Card';
import CardHeader from '@material-ui/core/CardHeader';
import DashboardIcon from '@material-ui/icons/Dashboard';
import DeviceStatusCircle from '@fbcnms/ui/components/icons/DeviceStatusCircle';
import GraphicEqIcon from '@material-ui/icons/GraphicEq';
import Grid from '@material-ui/core/Grid';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import PersonIcon from '@material-ui/icons/Person';
import React from 'react';
import SettingsIcon from '@material-ui/icons/Settings';
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
  card: {
    variant: 'outlined',
  },
}));

export default function SubscriberDetail(props: {
  subscriberMap: ?{[string]: subscriber},
}) {
  const classes = useStyles();
  const [tabPos, setTabPos] = React.useState(0);
  const {relativePath, relativeUrl, match} = useRouter();
  const subscriberId: string = nullthrows(match.params.subscriberId);
  const subscriberInfo = props.subscriberMap?.[subscriberId];
  if (!subscriberInfo) {
    return null;
  }

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Subscriber/{subscriberId}
        </Text>
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
        </Grid>
      </AppBar>
      <Switch>
        <Route
          path={relativePath('/overview')}
          render={() => <Overview subscriberInfo={subscriberInfo} />}
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
      <Grid container spacing={3}>
        <Grid container item xs={12} spacing={3}>
          <Grid item xs={6}>
            <Text>
              <PersonIcon />
              Subscriber
            </Text>
            <Info subscriberInfo={props.subscriberInfo} />
          </Grid>
          <Grid item xs={6}>
            <Text>
              <GraphicEqIcon />
              Status
            </Text>
            <Status subscriberInfo={props.subscriberInfo} />
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <Paper>
            <div className={classes.contentPlaceholder}>Data Usage Chart</div>
          </Paper>
        </Grid>
        <Grid item xs={12}>
          <Paper>
            <div className={classes.contentPlaceholder}>
              Events Table Goes Here
            </div>
          </Paper>
        </Grid>
      </Grid>
    </div>
  );
}

function Info(props: {subscriberInfo: subscriber}) {
  return (
    <Grid container>
      <Grid item xs={12}>
        <Card variant={'outlined'}>
          <CardHeader
            titleTypographyProps={{variant: 'caption'}}
            subheaderTypographyProps={{variant: 'body1'}}
            data-testid="Name"
            subheader={props.subscriberInfo.id}
          />
        </Card>
      </Grid>
      <Grid container item xs={12}>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              data-testid="IMSI"
              title="IMSI"
              subheader={props.subscriberInfo.id}
            />
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              data-testid="service"
              title="Service"
              subheader={
                <>
                  <DeviceStatusCircle
                    isGrey={false}
                    isActive={props.subscriberInfo.lte.state === 'ACTIVE'}
                  />
                  <Text>{props.subscriberInfo.lte.state}</Text>
                </>
              }
            />
          </Card>
        </Grid>
      </Grid>
    </Grid>
  );
}

function Status(props: {subscriberInfo: subscriber}) {
  const featureUnsupported = 'Unsupported';
  const statusUnknown = 'Unknown';
  return (
    <Grid container>
      <Grid container item xs={12}>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              title="Gateway ID"
              subheader={featureUnsupported}
            />
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              title="eNodeB SN"
              subheader={featureUnsupported}
            />
          </Card>
        </Grid>
      </Grid>
      <Grid container item xs={12}>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              title="Connection Status"
              subheader={statusUnknown}
            />
          </Card>
        </Grid>
        <Grid item xs={6}>
          <Card variant={'outlined'}>
            <CardHeader
              titleTypographyProps={{variant: 'caption'}}
              subheaderTypographyProps={{variant: 'body1'}}
              data-testid="UE Latency"
              title="UE Latency"
              subheader={
                props.subscriberInfo.monitoring?.icmp?.latency_ms ??
                statusUnknown
              }
            />
          </Card>
        </Grid>
      </Grid>
    </Grid>
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
