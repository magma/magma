/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import React from 'react';
import emptyFunction from '../../../../../fbcnms-packages/fbcnms-util/emptyFunction';

export type Bookmark = {
  id: string,
  name: string,
};

export type PowerSearchContextValue = {
  bookmark: ?Bookmark,
  setBookmark: (?Bookmark) => void,
};

const PowerSearchContext = React.createContext<PowerSearchContextValue>({
  bookmark: null,
  setBookmark: emptyFunction,
});

export function usePowerSearch() {
  return React.useContext(PowerSearchContext);
}

export default PowerSearchContext;
