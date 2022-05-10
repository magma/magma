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

import type {EmbeddedData, User} from '../../../shared/types/embeddedData';
import type {FeatureID} from '../../../fbc_js_core/types/features';
import type {SSOSelectedType} from '../../../fbc_js_core/types/auth';
import type {Tab} from '../../../fbc_js_core/types/tabs';

import * as React from 'react';
import {noop} from 'lodash';

export type AppContextType = {
  csrfToken: ?string,
  version: ?string,
  networkIds: string[],
  tabs: $ReadOnlyArray<Tab>,
  user: User,
  showExpandButton: () => void,
  hideExpandButton: () => void,
  isOrganizations: boolean,
  isFeatureEnabled: FeatureID => boolean,
  isTabEnabled: Tab => boolean,
  ssoEnabled: boolean,
  ssoSelectedType: SSOSelectedType,
  hasAccountSettings: boolean,
};

const appContextDefaults = {
  csrfToken: null,
  version: null,
  networkIds: [],
  tabs: [],
  user: {tenant: '', email: '', isSuperUser: false, isReadOnlyUser: false},
  showExpandButton: noop,
  hideExpandButton: noop,
  isFeatureEnabled: () => false,
  isTabEnabled: () => false,
  ssoEnabled: false,
  ssoSelectedType: 'none',
  hasAccountSettings: false,
};

// $FlowFixMe[prop-missing]
const AppContext = React.createContext<AppContextType>(appContextDefaults);

type Props = {|
  children: React.Node,
  isOrganizations?: boolean,
  networkIDs?: string[],
|};

export function AppContextProvider(props: Props) {
  const config: {appData: EmbeddedData} = window.CONFIG;
  const {appData} = config;
  const value = {
    ...appContextDefaults,
    ...appData,
    hasAccountSettings: !appData.ssoEnabled,
    isOrganizations: !!props.isOrganizations, // is organizations management aka. the host site
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
