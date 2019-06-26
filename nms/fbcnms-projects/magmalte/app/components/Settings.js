/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {WithStyles} from '@material-ui/core';
import type {ContextRouter} from 'react-router-dom';

import React from 'react';
import AppBar from '@material-ui/core/AppBar';
import AppContext from './context/AppContext';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import Tabs from '@material-ui/core/Tabs';
import Tab from '@material-ui/core/Tab';
import UsersSettings from './UsersSettings';
import SecuritySettings from './SecuritySettings';
import MagmaTopBar from './MagmaTopBar';

import {findIndex} from 'lodash';
import {withStyles} from '@material-ui/core/styles';
import {Route, Switch, withRouter} from 'react-router-dom';

const styles = theme => ({
  tabs: {
    flex: 1,
  },
  paper: {
    margin: theme.spacing.unit * 3,
  },
});

type Props = ContextRouter & WithStyles;

type State = {
  currentTab: number,
};

class Settings extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    // Default to Security sub-section
    if (props.location.pathname.endsWith('/settings/')) {
      props.history.push(`${props.match.url}security/`);
    }

    const {pathname: currentPath} = props.location;
    const currentTab = findIndex(['security', 'users'], route =>
      currentPath.startsWith(props.match.url + '/' + route),
    );

    this.state = {
      currentTab: currentTab !== -1 ? currentTab : 0,
    };
  }

  render() {
    const {match, classes} = this.props;
    return (
      <AppContext.Consumer>
        {({user, networkIds}) => (
          <>
            <MagmaTopBar title="Settings" />
            <Paper className={this.props.classes.paper}>
              <AppBar position="static" color="default">
                <Tabs
                  value={this.state.currentTab}
                  indicatorColor="primary"
                  textColor="primary"
                  onChange={this.onTabChange}
                  className={classes.tabs}>
                  <Tab
                    component={NestedRouteLink}
                    label="Security"
                    to="/security/"
                  />
                  {user.isSuperUser && (
                    <Tab
                      component={NestedRouteLink}
                      label="Users"
                      to="/users/"
                    />
                  )}
                </Tabs>
              </AppBar>
              <Switch>
                <Route
                  path={`${match.path}/security`}
                  component={SecuritySettings}
                />
                {user.isSuperUser && (
                  <Route
                    path={`${match.path}/users`}
                    component={() => (
                      <UsersSettings allNetworkIDs={networkIds} />
                    )}
                  />
                )}
              </Switch>
            </Paper>
          </>
        )}
      </AppContext.Consumer>
    );
  }

  onTabChange = (_, currentTab: number) => this.setState({currentTab});
}

export default withStyles(styles)(withRouter(Settings));
