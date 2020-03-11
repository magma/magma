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
import ProjectComparisonView from '../projects/ProjectComparisonView';
import React, {useContext} from 'react';
import WorkOrderComparisonView from './WorkOrderComparisonView';
import WorkOrderConfigure from './WorkOrderConfigure';
import nullthrows from '@fbcnms/util/nullthrows';
import {Redirect, Route, Switch} from 'react-router-dom';
import {WorkOrdersNavListItems} from './WorkOrdersNavListItems';
import {getProjectLinks} from '@fbcnms/magmalte/app/common/projects';
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

function WorkOrdersMain() {
  const classes = useStyles();
  const {tabs, user} = useContext(AppContext);
  const relativeUrl = useRelativeUrl();

  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={<WorkOrdersNavListItems />}
        projects={getProjectLinks(tabs, user)}
        showSettings={true}
        user={nullthrows(user)}
      />
      <AppContent>
        <Switch>
          <Route
            path={relativeUrl('/search')}
            component={WorkOrderComparisonView}
          />
          <Route
            path={relativeUrl('/projects/search')}
            component={ProjectComparisonView}
          />
          <Route
            path={relativeUrl('/configure')}
            component={WorkOrderConfigure}
          />
          <Redirect to={relativeUrl('/search')} />
        </Switch>
      </AppContent>
    </div>
  );
}

export default () => {
  return (
    <ApplicationMain>
      <AppContextProvider>
        <WorkOrdersMain />
      </AppContextProvider>
    </ApplicationMain>
  );
};
