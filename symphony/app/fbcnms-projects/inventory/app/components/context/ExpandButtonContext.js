/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
'use strict';

import React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';

type Context = {
  isExpandButtonShown: boolean,
  isExpanded: boolean,
  showExpandButton: () => void,
  hideExpandButton: () => void,
  expand: () => void,
  collapse: () => void,
};

export default React.createContext<Context>({
  isExpandButtonShown: false,
  isExpanded: false,
  showExpandButton: emptyFunction,
  hideExpandButton: emptyFunction,
  expand: emptyFunction,
  collapse: emptyFunction,
});
