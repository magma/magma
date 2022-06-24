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
 */

import ErrorLayout from './main/ErrorLayout';
import Index, {ROOT_PATHS} from './main/Index';
import IndexWithoutNetwork from './IndexWithoutNetwork';
import NetworkError from './main/NetworkError';
import NoNetworksMessage from './NoNetworksMessage';
import React from 'react';
import {AppContextProvider} from './context/AppContext';
import {
  Navigate,
  Route,
  Routes,
  useLocation,
  useParams,
} from 'react-router-dom';

import LoadingFiller from './LoadingFiller';
import MagmaAPI from '../../api/MagmaAPI';
import useMagmaAPI from '../../api/useMagmaAPI';
import {sortBy} from 'lodash';

export const NO_NETWORK_MESSAGE =
  'You currently do not have access to any networks. Please contact your system administrator to be added';

function Nms({networkIds}: {networkIds: Array<string>}) {
  const {networkId} = useParams();

  if (networkIds.length > 0 && !networkId) {
    return <Navigate to={`/nms/${networkIds[0]}`} replace />;
  }

  const hasNoNetworks = networkIds.length === 0 && !ROOT_PATHS.has(networkId!);

  // If it's a superuser and there are no networks, prompt them to create a
  // network
  if (hasNoNetworks && window.CONFIG.appData.user.isSuperUser) {
    return <Navigate to="/admin/networks" replace />;
  }

  // If it's a regular user and there are no networks, then they likely dont
  // have access.
  if (hasNoNetworks && !window.CONFIG.appData.user.isSuperUser) {
    return (
      <AppContextProvider>
        <ErrorLayout>
          <NoNetworksMessage>{NO_NETWORK_MESSAGE}</NoNetworksMessage>
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

function NoNetworkFallback({networkIds}: {networkIds: Array<string>}) {
  const location = useLocation();
  if (networkIds.length > 0) {
    return (
      <Navigate to={`/nms/${networkIds[0]}${location.pathname}`} replace />
    );
  }

  return (
    <AppContextProvider>
      <IndexWithoutNetwork />
    </AppContextProvider>
  );
}

export default () => {
  const {response, error} = useMagmaAPI(MagmaAPI.networks.networksGet, {});
  const networkIds = sortBy(response, [n => n.toLowerCase()]);

  if (error) {
    return (
      <AppContextProvider>
        <ErrorLayout>
          <NetworkError error={error} />
        </ErrorLayout>
      </AppContextProvider>
    );
  }

  return !response ? (
    <LoadingFiller />
  ) : (
    <Routes>
      <Route
        path="/nms/:networkId/*"
        element={<Nms networkIds={networkIds} />}
      />
      <Route path="/nms" element={<Nms networkIds={networkIds} />} />
      <Route
        path="/*"
        element={<NoNetworkFallback networkIds={networkIds} />}
      />
    </Routes>
  );
};
