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

import ActionTable from '../../components/ActionTable';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SubscriberDetail from './SubscriberDetail';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
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
}));

export default function SubscriberDashboard() {
  const {match, relativePath, relativeUrl} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const {response: subscriberMap, isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {
      networkId: networkId,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <Switch>
      <Route
        path={relativePath('/overview/:subscriberId')}
        render={() => <SubscriberDetail subscriberMap={subscriberMap} />}
      />

      <Route
        path={relativePath('/overview')}
        render={() => (
          <SubscriberDashboardInternal subscriberMap={subscriberMap} />
        )}
      />
      <Redirect to={relativeUrl('/overview')} />
    </Switch>
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

function SubscriberDashboardInternal({
  subscriberMap,
}: {
  subscriberMap: ?{[string]: subscriber} | void,
}) {
  const classes = useStyles();
  const {history, relativeUrl} = useRouter();
  const [currRow, setCurrRow] = useState<SubscriberRowType>({});

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
              <Button className={classes.appBarBtn}>Secondary Action</Button>
            </Grid>
            <Grid item>
              <Button className={classes.appBarBtnSecondary}>
                Primary Action
              </Button>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>
      <div className={classes.dashboardRoot}>
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Text key="title">
              <PeopleIcon /> {TITLE}
            </Text>
          </Grid>

          <Grid item xs={12}>
            {subscriberMap ? (
              <ActionTable
                data={Object.keys(subscriberMap).map((imsi: string) => {
                  const subscriberInfo = subscriberMap[imsi];
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
                })}
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
            ) : (
              '<Text>No Subscribers Found</Text>'
            )}
          </Grid>
        </Grid>
      </div>
    </>
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
