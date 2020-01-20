/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';

export type TableContextValue = {
  showSelection: boolean,
};

const TableContext = React.createContext<TableContextValue>({
  showSelection: false,
});

export function useTable() {
  return React.useContext(TableContext);
}

export default TableContext;
