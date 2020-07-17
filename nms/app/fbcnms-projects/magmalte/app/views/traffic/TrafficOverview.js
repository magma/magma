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
import ApnOverview from './ApnOverview';
import AppBar from '@material-ui/core/AppBar';
import Grid from '@material-ui/core/Grid';
import LibraryBooksIcon from '@material-ui/icons/LibraryBooks';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import PolicyOverview from './PolicyOverview';
import React from 'react';
import RssFeedIcon from '@material-ui/icons/RssFeed';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import Text from '@fbcnms/ui/components/design-system/Text';

import {Redirect, Route, Switch} from 'react-router-dom';
import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const POLICY_TITLE = 'Policies';
const APN_TITLE = 'APNs';
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
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

export default function TrafficDashboard() {
  const classes = useStyles();
  const {relativePath, relativeUrl} = useRouter();
  const [tabPos, setTabPos] = React.useState(0);

  return (
    <>
      <div className={classes.topBar}>
        <Text color="light" weight="medium">
          Traffic
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
                key="Policy"
                component={NestedRouteLink}
                label={<PolicyDashboardTabLabel />}
                to="/policy"
                className={classes.tab}
              />
              <Tab
                key="APNs"
                component={NestedRouteLink}
                label={<APNDashboardTabLabel />}
                to="/apn"
                className={classes.tab}
              />
            </Tabs>
          </Grid>
        </Grid>
      </AppBar>

      <Switch>
        <Route path={relativePath('/policy')} component={PolicyOverview} />
        <Route path={relativePath('/apn')} component={ApnOverview} />
        <Redirect to={relativeUrl('/policy')} />
      </Switch>
    </>
  );
}

function PolicyDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <LibraryBooksIcon className={classes.tabIconLabel} /> {POLICY_TITLE}
    </div>
  );
}

function APNDashboardTabLabel() {
  const classes = useStyles();

  return (
    <div className={classes.tabLabel}>
      <RssFeedIcon className={classes.tabIconLabel} /> {APN_TITLE}
    </div>
  );
}
