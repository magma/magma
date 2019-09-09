/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Admin from './admin/Admin';
import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar.react';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import NetworkContext from './context/NetworkContext';
import NetworkSelector from './NetworkSelector.react';
import React, {useContext} from 'react';
import SectionLinks from './layout/SectionLinks';
import SectionRoutes from './layout/SectionRoutes';
import VersionTooltip from './VersionTooltip';
import {Redirect, Route, Switch} from 'react-router-dom';

import {MagmaAPIUrls} from '../common/MagmaAPI';
import {getProjectLinks} from '../common/projects';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  root: {
    display: 'flex',
  },
  toolbarIcon: {
    display: 'flex',
    alignItems: 'center',
    justifyContent: 'flex-end',
    padding: '0 8px',
    ...theme.mixins.toolbar,
  },
}));

// These won't be considered networkIds
const ROOT_PATHS = new Set(['network']);

function Index() {
  const classes = useStyles();
  const {match} = useRouter();
  const {user, tabs} = useContext(AppContext);
  const networkId = ROOT_PATHS.has(match.params.networkId)
    ? null
    : match.params.networkId;

  return (
    <NetworkContext.Provider value={{networkId}}>
      <div className={classes.root}>
        <AppSideBar
          mainItems={[<SectionLinks key={1} />, <VersionTooltip key={2} />]}
          secondaryItems={[<NetworkSelector key={1} />]}
          projects={getProjectLinks(tabs, user)}
          user={user}
        />
        <AppContent>
          <SectionRoutes />
        </AppContent>
      </div>
    </NetworkContext.Provider>
  );
}

function Main() {
  const {match} = useRouter();
  const {response, error} = useAxios({
    method: 'get',
    url: MagmaAPIUrls.networks(),
  });

  const networkIds = sortBy(response?.data, [n => n.toLowerCase()]) || [
    'mpk_test',
  ];
  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };

  if (networkIds.length > 0 && !match.params.networkId) {
    return <Redirect to={`/nms/${networkIds[0]}/map/`} />;
  }

  if (
    response &&
    !error &&
    networkIds.length === 0 &&
    window.CONFIG.appData.user.isSuperUser &&
    match.params.networkId !== 'network'
  ) {
    return <Redirect to="/nms/network/create" />;
  }

  return (
    <ApplicationMain appContext={appContext}>
      <Index />
    </ApplicationMain>
  );
}

export default () => (
  <Switch>
    <Route path="/nms/:networkId" component={Main} />
    <Route path="/nms" component={Main} />
    <Route path="/admin" component={Admin} />
  </Switch>
);
