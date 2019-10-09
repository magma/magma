/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {$AxiosError} from 'axios';

import Admin from './admin/Admin';
import AppContent from '@fbcnms/ui/components/layout/AppContent';
import AppContext from '@fbcnms/ui/context/AppContext';
import AppSideBar from '@fbcnms/ui/components/layout/AppSideBar.react';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import MagmaV1API from '../common/MagmaV1API';
import NetworkContext from './context/NetworkContext';
import NetworkSelector from './NetworkSelector.react';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage.react';
import React, {useContext} from 'react';
import SectionLinks from './layout/SectionLinks';
import SectionRoutes from './layout/SectionRoutes';
import VersionTooltip from './VersionTooltip';
import {Redirect, Route, Switch} from 'react-router-dom';

import useMagmaAPI from '../common/useMagmaAPI';
import {getProjectLinks} from '../common/projects';
import {makeStyles} from '@material-ui/styles';
import {sortBy} from 'lodash';
import {useRouter, useSnackbar} from '@fbcnms/ui/hooks';

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

function Index(props: {noAccess: boolean}) {
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
          {props.noAccess ? <NoNetworksMessage /> : <SectionRoutes />}
        </AppContent>
      </div>
    </NetworkContext.Provider>
  );
}

function NetworkError({error}: {error: $AxiosError<string>}) {
  const classes = useStyles();
  const {user, tabs} = useContext(AppContext);
  let errorMessage = error.message;
  if (error.response && error.response.status >= 400) {
    errorMessage = error.response?.statusText;
  }
  useSnackbar(
    'Unable to communicate with magma controller: ' + errorMessage,
    {variant: 'error'},
    !!error,
  );
  return (
    <div className={classes.root}>
      <AppSideBar
        mainItems={[]}
        secondaryItems={[]}
        projects={getProjectLinks(tabs, user)}
        user={user}
      />
      <AppContent>
        <div />
      </AppContent>
    </div>
  );
}

function Main() {
  const {match} = useRouter();
  const {response, error, isLoading} = useMagmaAPI(MagmaV1API.getNetworks, {});

  const networkIds = sortBy(response, [n => n.toLowerCase()]) || ['mpk_test'];
  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };

  if (error) {
    return (
      <ApplicationMain appContext={appContext}>
        <NetworkError error={error} />
      </ApplicationMain>
    );
  }

  if (networkIds.length > 0 && !match.params.networkId) {
    return <Redirect to={`/nms/${networkIds[0]}/map/`} />;
  }

  if (
    response &&
    networkIds.length === 0 &&
    window.CONFIG.appData.user.isSuperUser &&
    match.params.networkId !== 'network'
  ) {
    return <Redirect to="/nms/network/create" />;
  }

  return (
    <ApplicationMain appContext={appContext}>
      <Index noAccess={!isLoading && networkIds.length === 0} />
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
