/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

import Admin from './admin/Admin';
import ApplicationMain from './ApplicationMain';
import ErrorLayout from './main/ErrorLayout';
import Index, {ROOT_PATHS} from './main/Index';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NetworkError from './main/NetworkError';
import NoNetworksMessage from '@fbcnms/ui/components/NoNetworksMessage';
import React from 'react';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {Redirect, Route, Switch} from 'react-router-dom';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {sortBy} from 'lodash';
import {useRouter} from '@fbcnms/ui/hooks';

function Main() {
  const {match} = useRouter();
  const {response, error} = useMagmaAPI(MagmaV1API.getNetworks, {});

  const networkIds = sortBy(response, [n => n.toLowerCase()]) || ['mpk_test'];

  if (error) {
    return (
      <AppContextProvider>
        <ErrorLayout>
          <NetworkError error={error} />
        </ErrorLayout>
      </AppContextProvider>
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
    return <Redirect to="/admin/networks" />;
  }

  // If it's a regular user and there are no networks, then they likely dont
  // have access.
  if (hasNoNetworks && !window.CONFIG.appData.user.isSuperUser) {
    return (
      <AppContextProvider>
        <ErrorLayout>
          <NoNetworksMessage>
            You currently do not have access to any networks. Please contact
            your system administrator to be added
          </NoNetworksMessage>
        </ErrorLayout>
      </AppContextProvider>
    );
  }

  return (
    <AppContextProvider networkIDs={networkIds}>
      <Index />
    </AppContextProvider>
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
