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
import DataPlanConfig from './DataPlanConfig';
import MagmaTopBar from '../MagmaTopBar';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import NetworkConfig from './NetworkConfig';
import Paper from '@material-ui/core/Paper';
import PoliciesConfig from './PoliciesConfig';
import React from 'react';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import UpgradeConfig from './UpgradeConfig';

import nullthrows from '@fbcnms/util/nullthrows';
import {Route, Switch, withRouter} from 'react-router-dom';
import {findIndex} from 'lodash';

import {withStyles} from '@material-ui/core/styles';

const styles = theme => ({
  paper: {
    margin: theme.spacing.unit * 3,
  },
  tabs: {
    flex: 1,
  },
});

type Props = WithStyles & ContextRouter & {};

type State = {
  currentTab: number,
};

type TabRoutes = {
  component: ComponentType<any>,
  label: string,
  path: string,
}[];

const getTabRoutes = (_networkID): TabRoutes => [
  {
    component: DataPlanConfig,
    label: 'Data Plans',
    path: 'dataplans',
  },
  {
    component: NetworkConfig,
    label: 'Network Configuration',
    path: 'network',
  },
  {
    component: UpgradeConfig,
    label: 'Upgrades',
    path: 'upgrades',
  },
  {
    component: PoliciesConfig,
    label: 'Policies',
    path: 'policies',
  },
];

class Configure extends React.Component<Props, State> {
  constructor(props: Props) {
    super(props);

    const networkID = nullthrows(props.match.params.networkId);
    const tabRoutes = getTabRoutes(networkID);

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
    const {classes, match} = this.props;
    const {currentTab} = this.state;
    const tabRoutes = getTabRoutes(nullthrows(match.params.networkId));
    return (
      <>
        <MagmaTopBar title="Configure" />
        <Paper className={this.props.classes.paper}>
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
      </>
    );
  }
}

export default withStyles(styles)(withRouter(Configure));
