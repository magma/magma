/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AppContext from '@fbcnms/ui/context/AppContext';
import AssignmentIcon from '@material-ui/icons/Assignment';
import LinearScaleIcon from '@material-ui/icons/LinearScale';
import MapIcon from '@material-ui/icons/Map';
import NavListItem from '@fbcnms/ui/components/NavListItem';
import React, {useContext} from 'react';
import SearchIcon from '@material-ui/icons/Search';
import ViewListIcon from '@material-ui/icons/ViewList';
import {LogEvents, ServerLogger} from '../common/LoggingUtils';
import {useRouter} from '@fbcnms/ui/hooks';

export default function MainNavListItems() {
  const {relativeUrl} = useRouter();
  const servicesEnabled = useContext(AppContext).isFeatureEnabled('services');
  return [
    <NavListItem
      key={1}
      label="Search"
      path={relativeUrl('/search')}
      icon={<SearchIcon />}
      onClick={() => ServerLogger.info(LogEvents.SEARCH_NAV_CLICKED)}
    />,
    <NavListItem
      key={2}
      label="Locations"
      path={relativeUrl('/inventory')}
      icon={<ViewListIcon />}
      onClick={() => ServerLogger.info(LogEvents.INVENTORY_NAV_CLICKED)}
    />,
    <NavListItem
      key={3}
      label="Map"
      path={relativeUrl('/map')}
      icon={<MapIcon />}
      onClick={() => ServerLogger.info(LogEvents.MAP_NAV_CLICKED)}
    />,
    <NavListItem
      key={4}
      label="Catalog"
      path={relativeUrl('/configure')}
      icon={<AssignmentIcon />}
      onClick={() => ServerLogger.info(LogEvents.CONFIGURE_NAV_CLICKED)}
    />,
    <NavListItem
      key={5}
      label="Services"
      path={relativeUrl('/services')}
      icon={<LinearScaleIcon />}
      onClick={() => ServerLogger.info(LogEvents.SERVICES_NAV_CLICKED)}
      hidden={!servicesEnabled}
    />,
  ];
}
