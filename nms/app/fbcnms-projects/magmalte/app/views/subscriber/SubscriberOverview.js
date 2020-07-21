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
import AddSubscriberButton from './SubscriberAddDialog';
import AppBar from '@material-ui/core/AppBar';
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PeopleIcon from '@material-ui/icons/People';
import React from 'react';
import SubscriberContext from '../../components/context/SubscriberContext';
import SubscriberDetail from './SubscriberDetail';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {CardTitleRow} from '../../components/layout/CardTitleRow';
import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

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
  const [subscriberMap, setSubscriberMap] = useState({});
  const {isLoading} = useMagmaAPI(
    MagmaV1API.getLteByNetworkIdSubscribers,
    {
      networkId: networkId,
    },
    useCallback(response => setSubscriberMap(response), []),
  );

  const updateSubscriberMap = async (key: string, val: subscriber) => {
    if (key in subscriberMap) {
      await MagmaV1API.putLteByNetworkIdSubscribersBySubscriberId({
        networkId: networkId,
        subscriber: val,
        subscriberId: key,
      });
    } else {
      await MagmaV1API.postLteByNetworkIdSubscribers({
        networkId: networkId,
        subscriber: val,
      });
    }
    setSubscriberMap({...subscriberMap, [key]: val});
  };

  if (isLoading) {
    return <LoadingFiller />;
  }

  return (
    <Switch>
      <Route
        path={relativePath('/overview/:subscriberId')}
        render={() => {
          return (
            <SubscriberContext.Provider
              value={{
                state: subscriberMap ?? {},
                setState: updateSubscriberMap,
              }}>
              <SubscriberDetail subscriberMap={subscriberMap} />
            </SubscriberContext.Provider>
          );
        }}
      />

      <Route
        path={relativePath('/overview')}
        render={() => {
          return (
            <SubscriberContext.Provider
              value={{
                state: subscriberMap ?? {},
                setState: updateSubscriberMap,
              }}>
              <SubscriberDashboardInternal subscriberMap={subscriberMap} />
            </SubscriberContext.Provider>
          );
        }}
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
        <Grid container direction="row" justify="flex-end" alignItems="center">
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
                  Primary Action
                </Button>
              </Grid>
            </Grid>
          </Grid>
        </Grid>
      </AppBar>

      <div className={classes.dashboardRoot}>
        <Grid container spacing={4}>
          <Grid item xs={12}>
            <Grid container>
              <Grid item xs={6}>
                <CardTitleRow icon={PeopleIcon} label={TITLE} />
              </Grid>
              <Grid container item xs={6} justify="flex-end">
                <AddSubscriberButton />
              </Grid>
            </Grid>

            {subscriberMap ? (
              <ActionTable
                data={Object.keys(subscriberMap).map((imsi: string) => {
                  const subscriberInfo = subscriberMap[imsi];
                  return {
                    name: subscriberInfo.name ?? imsi,
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
