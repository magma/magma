/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import APMetrics from './APMetrics';
import AppBar from '@material-ui/core/AppBar';
import CWFNetworkMetrics from './CWFNetworkMetrics';
import IMSIMetrics from './IMSIMetrics';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {Redirect, Route, Switch} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  bar: {
    backgroundColor: theme.palette.blueGrayDark,
  },
  tabs: {
    flex: 1,
    color: 'white',
  },
}));

export default function() {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, location} = useRouter();

  const currentTab = findIndex(['ap', 'network', 'subscribers'], route =>
    location.pathname.startsWith(match.url + '/' + route),
  );

  return (
    <>
      <AppBar position="static" color="default" className={classes.bar}>
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="inherit"
          className={classes.tabs}>
          <Tab component={NestedRouteLink} label="Access Points" to="/ap" />
          <Tab component={NestedRouteLink} label="Network" to="/network" />
          <Tab
            component={NestedRouteLink}
            label="Subscribers"
            to="/subscribers"
          />
        </Tabs>
      </AppBar>
      <Switch>
        <Route path={relativePath('/ap')} component={APMetrics} />
        <Route path={relativePath('/network')} component={CWFNetworkMetrics} />
        <Route path={relativePath('/subscribers')} component={IMSIMetrics} />
        <Redirect to={relativeUrl('/ap')} />
      </Switch>
    </>
  );
}
