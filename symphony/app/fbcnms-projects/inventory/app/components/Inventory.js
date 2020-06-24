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
import Configure from '../pages/Configure';
import EquipmentComparisonView from './comparison_view/EquipmentComparisonView';
import ExpandButtonContext from './context/ExpandButtonContext';
import Inventory from '../pages/Inventory';
import InventoryComparisonView from './comparison_view/InventoryComparisonView';
import InventorySuspense from '../common/InventorySuspense';
import LocationsMap from './map/LocationsMap';
import MainNavListItems from './MainNavListItems';
import React, {useCallback, useContext, useEffect, useState} from 'react';
import ServicesMain from './services/ServicesMain';
import {DialogShowingContextProvider} from '@fbcnms/ui/components/design-system/Dialog/DialogShowingContext';
import {Redirect, Route, Switch} from 'react-router-dom';
import {getProjectLinks} from '@fbcnms/projects/projects';
import {makeStyles} from '@material-ui/styles';
import {setLoggerUser} from '../common/LoggingUtils';
import {shouldShowSettings} from '@fbcnms/magmalte/app/components/Settings';
import {useMainContext} from './MainContext';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
  },
}));

function Index() {
  const classes = useStyles();
  const {isExpandButtonShown, expand, collapse, isExpanded} = useContext(
    ExpandButtonContext,
  );

  const {tabs, ssoEnabled, isFeatureEnabled} = useContext(AppContext);
  const relativeUrl = useRelativeUrl();
  const {location} = useRouter();
  const {integrationUserDefinition} = useMainContext();

  const multiSubjectReports = isFeatureEnabled('multi_subject_reports');

  return (
    <div className={classes.root}>
      <AppSideBar
        useExpandButton={location.pathname.includes('inventory/inventory')}
        expanded={isExpanded}
        showExpandButton={isExpandButtonShown}
        onExpandClicked={() => (isExpanded ? collapse() : expand())}
        mainItems={<MainNavListItems />}
        projects={getProjectLinks(tabs, integrationUserDefinition)}
        showSettings={shouldShowSettings({
          isSuperUser: integrationUserDefinition.isSuperUser,
          ssoEnabled,
        })}
        user={integrationUserDefinition}
      />
      <AppContent>
        <Switch>
          <Route path={relativeUrl('/configure')} component={Configure} />
          <Route path={relativeUrl('/inventory')} component={Inventory} />
          <Route path={relativeUrl('/map')} component={LocationsMap} />
          <Route
            path={relativeUrl('/search')}
            component={
              multiSubjectReports
                ? InventoryComparisonView
                : EquipmentComparisonView
            }
          />
          <Route path={relativeUrl('/services')} component={ServicesMain} />
          <Redirect exact from="/" to={relativeUrl('/inventory')} />
          <Redirect exact from="/inventory" to={relativeUrl('/inventory')} />
        </Switch>
      </AppContent>
    </div>
  );
}

export default function IndexWrapper() {
  useEffect(() => setLoggerUser(window.CONFIG.appData.user), []);
  const [isExpanded, setIsExpanded] = useState(true);
  const expand = useCallback(() => setIsExpanded(true), []);
  const collapse = useCallback(() => setIsExpanded(false), []);

  const [isExpandButtonShown, showExpandButton] = useState(false);
  return (
    <ApplicationMain>
      <ExpandButtonContext.Provider
        value={{
          showExpandButton: () => showExpandButton(true),
          hideExpandButton: () => showExpandButton(false),
          expand: expand,
          collapse: collapse,
          isExpanded,
          isExpandButtonShown,
        }}>
        <InventorySuspense isTopLevel={true}>
          <DialogShowingContextProvider>
            <Index />
          </DialogShowingContextProvider>
        </InventorySuspense>
      </ExpandButtonContext.Provider>
    </ApplicationMain>
  );
}
