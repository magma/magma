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
import AppContext from '@fbcnms/ui/context/AppContext';
import ApplicationMain from '@fbcnms/ui/components/ApplicationMain';
import ErrorLayout from './main/ErrorLayout';
import Index, {ROOT_PATHS} from './main/Index';
import MagmaV1API from '../common/MagmaV1API';
import NetworkError from './main/NetworkError';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage.react';
import React from 'react';
import {Redirect, Route, Switch} from 'react-router-dom';

import useMagmaAPI from '../common/useMagmaAPI';
import {sortBy} from 'lodash';
import {useRouter} from '@fbcnms/ui/hooks';

function Main() {
  const {match} = useRouter();
  const {response, error} = useMagmaAPI(MagmaV1API.getNetworks, {});

  const networkIds = sortBy(response, [n => n.toLowerCase()]) || ['mpk_test'];
  const appContext = {
    ...window.CONFIG.appData,
    networkIds,
  };

  if (error) {
    return (
      <AppContext.Provider value={appContext}>
        <ErrorLayout>
          <NetworkError error={error} />
        </ErrorLayout>
      </AppContext.Provider>
    );
  }

  if (networkIds.length > 0 && !match.params.networkId) {
    return <Redirect to={`/nms/${networkIds[0]}/map/`} />;
  }

  const hasNoNetworks =
    response &&
    networkIds.length === 0 &&
    !ROOT_PATHS.has(match.params.networkId);

  // If it's a superuser and there are no networks, prompt them to create a
  // network
  if (hasNoNetworks && window.CONFIG.appData.user.isSuperUser) {
    return <Redirect to="/nms/network/create" />;
  }

  // If it's a regular user and there are no networks, then they likely dont
  // have access.
  if (hasNoNetworks && !window.CONFIG.appData.user.isSuperUser) {
    return (
      <AppContext.Provider value={appContext}>
        <ErrorLayout>
          <NoNetworksMessage />
        </ErrorLayout>
      </AppContext.Provider>
    );
  }

  return (
    <AppContext.Provider value={appContext}>
      <Index />
    </AppContext.Provider>
  );
}

export default () => (
  <ApplicationMain>
    <Switch>
      <Route path="/nms/:networkId" component={Main} />
      <Route path="/nms" component={Main} />
      <Route path="/admin" component={Admin} />
    </Switch>
  </ApplicationMain>
);
