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

export type PowerSearchContextValue = {
  bookmarkName: ?string,
  setBookmark: string => void,
};

const PowerSearchContext = React.createContext<PowerSearchContextValue>({
  bookmarkName: null,
  setBookmark: () => {},
});

export function usePowerSearch() {
  return React.useContext(PowerSearchContext);
}

export default PowerSearchContext;
