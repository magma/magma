/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {TabProps} from '@fbcnms/ui/components/design-system/Tabs/TabsBar';

import AppContext from '@fbcnms/ui/context/AppContext';
import EquipmentPortTypes from '../components/configure/EquipmentPortTypes';
import EquipmentTypes from '../components/configure/EquipmentTypes';
import InventoryErrorBoundary from '../common/InventoryErrorBoundary';
import InventorySuspense from '../common/InventorySuspense';
import LocationTypes from '../components/configure/LocationTypes';
import React, {useContext, useEffect, useMemo, useState} from 'react';
import ServiceTypes from '../components/configure/ServiceTypes';
import TabsBar from '@fbcnms/ui/components/design-system/Tabs/TabsBar';
import fbt from 'fbt';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useHistory, useLocation} from 'react-router';
import {useRelativeUrl} from '@fbcnms/ui/hooks/useRouter';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
    flexDirection: 'column',
    height: '100vh',
    transform: 'translateZ(0)',
  },
  tabs: {
    backgroundColor: 'white',
    borderBottom: `1px ${theme.palette.grey[200]} solid`,
    minHeight: '60px',
    overflow: 'visible',
  },
  tabContainer: {
    width: '250px',
  },
  tabsRoot: {
    top: 0,
    left: 0,
    right: 0,
    height: '60px',
  },
}));

type RouteTab = {
  id: string,
  tab: TabProps,
  path: string,
};

export default function Configure() {
  const relativeUrl = useRelativeUrl();
  const history = useHistory();
  const location = useLocation();
  const classes = useStyles();
  const servicesEnabled = useContext(AppContext).isFeatureEnabled('services');
  const tabBars: Array<RouteTab> = useMemo(
    () => [
      {
        id: 'equipment_types',
        tab: {
          label: fbt('EQUIPMENT', ''),
        },
        path: 'equipment_types',
      },
      {
        id: 'location_types',
        tab: {
          label: fbt('LOCATIONS', ''),
        },
        path: 'location_types',
      },
      {
        id: 'port_types',
        tab: {
          label: fbt('PORTS', ''),
        },
        path: 'port_types',
      },
      ...(servicesEnabled
        ? [
            {
              id: 'service_types',
              tab: {
                label: fbt('SERVICES', ''),
              },
              path: 'service_types',
            },
          ]
        : []),
    ],
    [servicesEnabled],
  );
  const tabMatch = location.pathname.match(/([^\/]*)\/*$/);
  const tabIndex =
    tabMatch == null ? -1 : tabBars.findIndex(el => el.id === tabMatch[1]);
  const [activeTabBar, setActiveTabBar] = useState<number>(
    tabIndex !== -1 ? tabIndex : 0,
  );

  useEffect(() => {
    ServerLogger.info(LogEvents.CONFIGURE_TAB_NAVIGATION_CLICKED, {
      id: tabBars[activeTabBar].id,
    });
    history.push(`/inventory/configure/${tabBars[activeTabBar].path}`);
  }, [tabBars, activeTabBar, history]);

  return (
    <div className={classes.root}>
      <TabsBar
        spread={true}
        size="large"
        tabs={tabBars.map(tabBar => tabBar.tab)}
        activeTabIndex={activeTabBar}
        onChange={setActiveTabBar}
      />
      <InventoryErrorBoundary>
        <InventorySuspense>
          <Switch>
            <Route
              path={relativeUrl('/equipment_types')}
              component={EquipmentTypes}
            />
            <Route
              path={relativeUrl('/location_types')}
              component={LocationTypes}
            />
            <Route
              path={relativeUrl('/port_types')}
              component={EquipmentPortTypes}
            />
            <Route
              path={relativeUrl('/service_types')}
              component={ServiceTypes}
            />
            <Redirect
              from={relativeUrl('/')}
              to={relativeUrl('/equipment_types')}
            />
          </Switch>
        </InventorySuspense>
      </InventoryErrorBoundary>
    </div>
  );
}
