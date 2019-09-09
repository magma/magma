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

import React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';
import type {FeatureID} from '@fbcnms/types/features';

export type User = {
  tenant: string,
  email: string,
  isSuperUser: boolean,
};

type Context = {
  csrfToken: ?string,
  version: ?string,
  networkIds: string[],
  tabs: string[],
  user: User,
  showExpandButton: () => void,
  hideExpandButton: () => void,
  enabledFeatures: FeatureID[],
};

export default React.createContext<Context>({
  csrfToken: null,
  version: null,
  networkIds: [],
  tabs: [],
  user: {tenant: '', email: '', isSuperUser: false},
  showExpandButton: emptyFunction,
  hideExpandButton: emptyFunction,
  enabledFeatures: [],
});
