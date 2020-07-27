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
 * @flow strict-local
 * @format
 */
'use strict';

import type {FeatureID} from '@fbcnms/types/features';
import type {SSOSelectedType} from '@fbcnms/types/auth';
import type {Tab} from '@fbcnms/types/tabs';

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';

export type User = {
  tenant: string,
  email: string,
  isSuperUser: boolean,
  isReadOnlyUser: boolean,
};

export type AppContextType = {
  csrfToken: ?string,
  version: ?string,
  networkIds: string[],
  tabs: $ReadOnlyArray<Tab>,
  user: User,
  showExpandButton: () => void,
  hideExpandButton: () => void,
  isFeatureEnabled: FeatureID => boolean,
  isTabEnabled: Tab => boolean,
  ssoEnabled: boolean,
  ssoSelectedType: SSOSelectedType,
};

export type AppContextAppData = {|
  csrfToken: string,
  tabs: $ReadOnlyArray<Tab>,
  user: User,
  enabledFeatures: FeatureID[],
  ssoEnabled: boolean,
  ssoSelectedType: SSOSelectedType,
  csvCharset: ?string,
|};

const appContextDefaults = {
  csrfToken: null,
  version: null,
  networkIds: [],
  tabs: [],
  user: {tenant: '', email: '', isSuperUser: false, isReadOnlyUser: false},
  showExpandButton: emptyFunction,
  hideExpandButton: emptyFunction,
  isFeatureEnabled: () => false,
  isTabEnabled: () => false,
  ssoEnabled: false,
  ssoSelectedType: 'none',
};

const AppContext = React.createContext<AppContextType>(appContextDefaults);

type Props = {|
  children: React.Node,
  networkIDs?: string[],
|};

export function AppContextProvider(props: Props) {
  const config: {appData: AppContextAppData} = window.CONFIG;
  const {appData} = config;
  const value = {
    ...appContextDefaults,
    ...appData,
    networkIds: props.networkIDs || [],
    isTabEnabled: (tab: Tab): boolean => {
      return appData.tabs?.indexOf(tab) !== -1;
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
