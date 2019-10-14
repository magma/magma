/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */
'use strict';

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';
import type {FeatureID} from '@fbcnms/types/features';

export type User = {
  tenant: string,
  email: string,
  isSuperUser: boolean,
};

export type AppContextType = {
  csrfToken: ?string,
  version: ?string,
  networkIds: string[],
  tabs: string[],
  user: User,
  showExpandButton: () => void,
  hideExpandButton: () => void,
  enabledFeatures: FeatureID[],
};

const AppContext = React.createContext<AppContextType>({
  csrfToken: null,
  version: null,
  networkIds: [],
  tabs: [],
  user: {tenant: '', email: '', isSuperUser: false},
  showExpandButton: emptyFunction,
  hideExpandButton: emptyFunction,
  enabledFeatures: [],
});

type Props = {|
  children: React.Node,
  networkIDs?: string[],
|};

export function AppContextProvider(props: Props) {
  const {appData} = window.CONFIG;
  const value = {
    ...appData,
    networkIds: props.networkIDs || [],
  };

  return (
    <AppContext.Provider value={value}>{props.children}</AppContext.Provider>
  );
}

export default AppContext;
