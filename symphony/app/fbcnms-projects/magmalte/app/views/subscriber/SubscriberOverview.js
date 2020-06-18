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
import ActionTable from '../../components/ActionTable';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

const TITLE = 'Subscribers';

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
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function SubscriberDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          {TITLE}
        </Text>
      </div>

      <AppBar position="static" color="default" className={classes.tabBar}>
        <Grid container>
          <Grid item xs={6}>
            <Tabs
              value={0}
              indicatorColor="primary"
              TabIndicatorProps={{style: {height: '5px'}}}
              textColor="inherit"
              className={classes.tabs}>
              <Tab
                key="Subscribers"
                component={NestedRouteLink}
                label={<SubscriberTabLabel />}
                to="/subscribersv2"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
          <Grid
            container
            item
            xs={6}
            justify="flex-end"
            alignItems="center"
            spacing={2}>
            <Grid item>
              <Button color="secondary" variant="contained">
                Secondary Action
              </Button>
            </Grid>
            <Grid item>
              <Button color="primary" variant="contained">
                Primary Action
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route
          path={relativePath('/subscriber')}
          component={SubscriberDashboardInternal}
        />
        <Redirect to={relativeUrl('/subscriber')} />
      </Switch>
    </>
  );
}

type SubscriberRowType = {
  name: string,
  imsi: string,
  service: string,
  currentUsage: string,
  dailyAvg: string,
  lastReportedTime: Date,
};

function SubscriberDashboardInternal() {
  const classes = useStyles();
  const {history, relativeUrl, match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});

  const {response, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {
      networkId: networkId,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  const subscriberRows: Array<SubscriberRowType> = response
    ? Object.keys(response).map((imsi: string) => {
        const subscriberInfo = response[imsi];
        return {
          name: subscriberInfo.id,
          imsi: imsi,
          service: subscriberInfo.lte.state,
          currentUsage: '0',
          dailyAvg: '0',
          lastReportedTime: new Date(
            subscriberInfo.monitoring?.icmp?.last_reported_time ?? 0,
          ),
        };
      })
    : [];
  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Text key="title">
            <PeopleIcon /> {TITLE}
          </Text>
        </Grid>

        <Grid item xs={12}>
          <ActionTable
            data={subscriberRows}
            columns={[
              {title: 'Name', field: 'name'},
              {title: 'IMSI', field: 'imsi'},
              {title: 'Service', field: 'service'},
              {title: 'Current Usage', field: 'currentUsage'},
              {title: 'Daily Average', field: 'dailyAvg'},
              {
                title: 'Last Reported Time',
                field: 'lastReportedTime',
                type: 'datetime',
              },
            ]}
            handleCurrRow={(row: SubscriberRowType) => setCurrRow(row)}
            menuItems={[
              {
                name: 'View',
                handleFunc: () => {
                  history.push(relativeUrl('/' + currRow.imsi));
                },
              },
              {name: 'Edit'},
              {name: 'Remove'},
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}

function SubscriberTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <PeopleIcon className={classes.tabIconLabel} /> {TITLE}
    </div>
  );
}
