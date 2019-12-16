/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {ComponentType} from 'react';
import type {ContextRouter} from 'react-router-dom';
import type {WithStyles} from '@material-ui/core';

import AppBar from '@material-ui/core/AppBar';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import Paper from '@material-ui/core/Paper';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';

import {Route, Switch, withRouter} from 'react-router-dom';
import {findIndex} from 'lodash';

import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  paper: {
    margin: theme.spacing(3),
  },
  tabs: {
    flex: 1,
  },
});

type Props = WithStyles<typeof styles> &
  ContextRouter & {
    tabRoutes: TabRoute[],
  };

type State = {
  currentTab: number,
};

type TabRoute = {
  component: ComponentType<any>,
  label: string,
  path: string,
};

class Configure extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);
    const {tabRoutes} = props;

    // Default to first page
    if (props.location.pathname.endsWith('/configure')) {
      props.history.push(`${props.match.url}/${tabRoutes[0].path}`);
    }

    const {pathname: currentPath} = props.location;
    const currentTab = findIndex(tabRoutes, route =>
      currentPath.startsWith(props.match.url + '/' + route.path),
    );

    this.state = {
      currentTab: currentTab !== -1 ? currentTab : 0,
    };
  }

  onTabChange = (event, currentTab: number) => {
    this.setState({currentTab});
  };

  render() {
    const {classes, match, tabRoutes} = this.props;
    const {currentTab} = this.state;
    return (
      <Paper className={this.props.classes.paper} elevation={2}>
        <AppBar position="static" color="default">
          <Tabs
            value={currentTab}
            indicatorColor="primary"
            textColor="primary"
            onChange={this.onTabChange}
            className={classes.tabs}>
            {tabRoutes.map((route, i) => (
              <Tab
                key={i}
                component={NestedRouteLink}
                label={route.label}
                to={route.path}
              />
            ))}
          </Tabs>
        </AppBar>
        <Switch>
          {tabRoutes.map((route, i) => (
            <Route
              key={i}
              path={`${match.path}/${route.path}`}
              component={route.component}
            />
          ))}
        </Switch>
      </Paper>
    );
  }
}

export default withStyles(styles)(withRouter(Configure));
