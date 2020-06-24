/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import InventorySuspense from '../../common/InventorySuspense';
import ProjectComparisonView from '../projects/ProjectComparisonView';
import React, {useContext} from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import WorkOrderComparisonView from './WorkOrderComparisonView';
import WorkOrderConfigure from './WorkOrderConfigure';
import {DialogShowingContextProvider} from '@fbcnms/ui/components/design-system/Dialog/DialogShowingContext';
import {Redirect, Route, Switch} from 'react-router-dom';
import {RelayEnvironmentProvider} from 'react-relay/hooks';
import {WorkOrdersNavListItems} from './WorkOrdersNavListItems';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {useMainContext} from '../../components/MainContext';
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
  const {tabs} = useContext(AppContext);
  const relativeUrl = useRelativeUrl();
  const {integrationUserDefinition} = useMainContext();

  return (
    <InventorySuspense isTopLevel={true}>
      <DialogShowingContextProvider>
        <div className={classes.root}>
          <AppSideBar
            mainItems={<WorkOrdersNavListItems />}
            projects={getProjectLinks(tabs, integrationUserDefinition)}
            showSettings={true}
            user={integrationUserDefinition}
          />
          <AppContent>
            <RelayEnvironmentProvider environment={RelayEnvironment}>
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
            </RelayEnvironmentProvider>
          </AppContent>
        </div>
      </DialogShowingContextProvider>
    </InventorySuspense>
  );
}

export default () => {
  return (
    <ApplicationMain>
      <WorkOrdersMain />
    </ApplicationMain>
  );
};
