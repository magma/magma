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

export type TableContextValue = {
  showSelection: boolean,
  clickableRows: boolean,
};

const TableContext = React.createContext<TableContextValue>({
  showSelection: false,
  clickableRows: false,
});

export function useTable() {
  return React.useContext(TableContext);
}

export default TableContext;
