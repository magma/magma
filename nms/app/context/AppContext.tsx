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
'use strict';

import * as React from 'react';
import {NetworkId} from '../../shared/types/network';
import {noop, sortBy} from 'lodash';
import {useState} from 'react';
import type {EmbeddedData, User} from '../../shared/types/embeddedData';
import type {FeatureID} from '../../shared/types/features';
import type {SSOSelectedType} from '../../shared/types/auth';

export const REFRESH_INTERVAL = 30000;

export type AppContextType = {
  csrfToken: string | null | undefined;
  version: string | null | undefined;
  networkIds: Array<string>;
  addNetworkId: (id: string) => void;
  removeNetworkId: (id: NetworkId) => void;
  user: User;
  showExpandButton: () => void;
  hideExpandButton: () => void;
  isOrganizations: boolean;
  isFeatureEnabled: (feature: FeatureID) => boolean;
  ssoEnabled: boolean;
  ssoSelectedType: SSOSelectedType;
  hasAccountSettings: boolean;
};

const appContextDefaults: AppContextType = {
  csrfToken: null,
  version: null,
  networkIds: [],
  addNetworkId: () => {},
  removeNetworkId: () => {},
  user: {
    tenant: '',
    email: '',
    isSuperUser: false,
    isReadOnlyUser: false,
  },
  showExpandButton: noop,
  hideExpandButton: noop,
  isFeatureEnabled: () => false,
  isOrganizations: false,
  ssoEnabled: false,
  ssoSelectedType: 'none',
  hasAccountSettings: false,
};

const AppContext = React.createContext<AppContextType>(appContextDefaults);
type Props = {
  children: React.ReactNode;
  isOrganizations?: boolean;
  networkIDs?: Array<string>;
};

export function AppContextProvider(props: Props) {
  const config: {
    appData: EmbeddedData;
  } = window.CONFIG;

  const [networkIds, setNetworkIds] = useState(props.networkIDs || []);
  const {appData} = config;
  const value = {
    ...appContextDefaults,
    ...appData,
    hasAccountSettings: !appData.ssoEnabled,
    // is organizations management aka. the host site
    isOrganizations: !!props.isOrganizations,
    networkIds,
    addNetworkId: (id: NetworkId) => {
      setNetworkIds(currentIds =>
        sortBy([...currentIds, id], [n => n.toLowerCase()]),
      );
    },
    removeNetworkId: (idToRemove: NetworkId) => {
      setNetworkIds(currentIds => currentIds.filter(id => id !== idToRemove));
    },
    isFeatureEnabled: (featureID: FeatureID): boolean => {
      return appData.enabledFeatures.indexOf(featureID) !== -1;
    },
  };

  return (
    <AppContext.Provider value={value}>{props.children}</AppContext.Provider>
  );
}

export default AppContext;
