/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AppContext from '@fbcnms/ui/context/AppContext';
import EquipmentPortTypes from '../components/configure/EquipmentPortTypes';
import EquipmentTypes from '../components/configure/EquipmentTypes';
import LocationTypes from '../components/configure/LocationTypes';
import React, {useContext} from 'react';
import ServiceTypes from '../components/configure/ServiceTypes';
import Tab from '@material-ui/core/Tab';
import Tabs from '@material-ui/core/Tabs';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {Redirect, Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

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

function getTabURI(newValue: number) {
  switch (newValue) {
    case 0:
      return 'equipment_types';
    case 1:
      return 'location_types';
    case 2:
      return 'port_types';
    case 3:
      return 'service_types';
    default:
      return '';
  }
}

function getTabValue(tabURI: string) {
  switch (tabURI) {
    case 'equipment_types':
      return 0;
    case 'location_types':
      return 1;
    case 'port_types':
      return 2;
    case 'service_types':
      return 3;
    default:
      return 0;
  }
}

export default function Configure() {
  const {location, history, relativeUrl} = useRouter();
  const classes = useStyles();
  const servicesEnabled = useContext(AppContext).isFeatureEnabled('services');
  return (
    <div className={classes.root}>
      <Tabs
        className={classes.tabs}
        classes={{flexContainer: classes.tabsRoot}}
        value={getTabValue(location.pathname.match(/([^\/]*)\/*$/)[1])}
        onChange={(e, newValue) => {
          const tab = getTabURI(newValue);
          ServerLogger.info(LogEvents.CONFIGURE_TAB_NAVIGATION_CLICKED, {
            tab,
          });
          history.push(`/inventory/configure/${tab}`);
        }}
        indicatorColor="primary"
        textColor="primary">
        <Tab
          data-testid="configure-equipment-tab"
          classes={{root: classes.tabContainer}}
          label="Equipment"
        />
        <Tab
          data-testid="configure-locations-tab"
          classes={{root: classes.tabContainer}}
          label="Locations"
        />
        <Tab classes={{root: classes.tabContainer}} label="Ports" />
        {servicesEnabled && (
          <Tab classes={{root: classes.tabContainer}} label="Services" />
        )}
      </Tabs>
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
        <Route path={relativeUrl('/service_types')} component={ServiceTypes} />
        <Redirect
          from={relativeUrl('/')}
          to={relativeUrl('/equipment_types')}
        />
      </Switch>
    </div>
  );
}
