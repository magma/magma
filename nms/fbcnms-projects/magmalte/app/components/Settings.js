/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppBar from '@material-ui/core/AppBar';
import AppContext from '@fbcnms/ui/context/AppContext';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React, {useContext} from 'react';
import SecuritySettings from './SecuritySettings';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Redirect, Route, Switch} from 'react-router-dom';
import {findIndex} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  tabs: {
    flex: 1,
  },
  paper: {
    margin: theme.spacing(3),
  },
}));

export default function Settings(props: {isSuperUser?: boolean}) {
  const classes = useStyles();
  const {match, relativePath, relativeUrl, location} = useRouter();
  const {user} = useContext(AppContext);
  let {isSuperUser} = props;
  if (isSuperUser === undefined) {
    isSuperUser = user.isSuperUser;
  }

  const currentTab = findIndex(['security', 'users'], route =>
    location.pathname.startsWith(match.url + '/' + route),
  );

  return (
    <Paper className={classes.paper} elevation={2}>
      <AppBar position="static" color="default">
        <Tabs
          value={currentTab !== -1 ? currentTab : 0}
          indicatorColor="primary"
          textColor="primary"
          className={classes.tabs}>
          <Tab component={NestedRouteLink} label="Security" to="/security/" />
          {isSuperUser && (
            <Tab component={NestedRouteLink} label="Users" to="/users/" />
          )}
        </Tabs>
      </AppBar>
      <Switch>
        <Route path={relativePath('/security')} component={SecuritySettings} />
        {isSuperUser && (
          <Route
            path={relativePath('/users')}
            render={() => <Redirect to="/admin/users" />}
          />
        )}
        <Redirect to={relativeUrl('/security')} />
      </Switch>
    </Paper>
  );
}
