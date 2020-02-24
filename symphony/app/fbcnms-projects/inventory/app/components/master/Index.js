/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext, {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import AssignmentIcon from '@material-ui/icons/Assignment';
import CloudMetrics from './CloudMetrics';
import Features from './Features';
import FlagIcon from '@material-ui/icons/Flag';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import OrganizationEdit from './OrganizationEdit';
import Organizations from './Organizations';
import Paper from '@material-ui/core/Paper';
import PeopleIcon from '@material-ui/icons/People';
import React, {useContext} from 'react';
import SecuritySettings from '@fbcnms/magmalte/app/components/SecuritySettings';
import ShowChartIcon from '@material-ui/icons/ShowChart';
import UsersSettings from '@fbcnms/magmalte/app/components/UsersSettings';
import nullthrows from '@fbcnms/util/nullthrows';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  paper: {
    margin: theme.spacing(3),
    padding: theme.spacing(),
  },
}));

function NavItems() {
  const relativeUrl = useRelativeUrl();
  return (
    <>
      <NavListItem
        label="Organizations"
        path={relativeUrl('/organizations')}
        icon={<AssignmentIcon />}
      />
      <NavListItem
        label="Features"
        path={relativeUrl('/features')}
        icon={<FlagIcon />}
      />
      <NavListItem
        label="Metrics"
        path={relativeUrl('/metrics')}
        icon={<ShowChartIcon />}
      />
      <NavListItem
        label="Users"
        path={relativeUrl('/users')}
        icon={<PeopleIcon />}
      />
    </>
  );
}

function Master() {
  const classes = useStyles();
  const {user, ssoEnabled} = useContext(AppContext);
  const relativeUrl = useRelativeUrl();

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={<NavItems />}
        user={nullthrows(user)}
        showSettings={!ssoEnabled}
      />
      <AppContent>
        <Switch>
          <Route
            path={relativeUrl('/organizations/detail/:name')}
            component={OrganizationEdit}
          />
          <Route
            path={relativeUrl('/organizations')}
            component={Organizations}
          />
          <Route path={relativeUrl('/features')} component={Features} />
          <Route path={relativeUrl('/metrics')} component={CloudMetrics} />
          <Route path={relativeUrl('/users')} component={UsersSettings} />
          <Route
            path={relativeUrl('/settings')}
            render={() => (
              <Paper className={classes.paper}>
                <SecuritySettings />
              </Paper>
            )}
          />
          <Redirect to={relativeUrl('/organizations')} />
        </Switch>
      </AppContent>
    </div>
  );
}

const Index = () => {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <Master />
      </AppContextProvider>
    </ApplicationMain>
  );
};

export default () => <Route path="/master" component={Index} />;
