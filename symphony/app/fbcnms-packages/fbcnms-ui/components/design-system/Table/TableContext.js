/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import emptyFunction from '@fbcnms/util/emptyFunction';
import {createContext, useContext, useEffect, useState} from 'react';

export const TABLE_SORT_ORDER = {
  ascending: 'ascending',
  descending: 'descending',
};

export type TableSortOrders = $Keys<typeof TABLE_SORT_ORDER>;

export type TableSortSettings = $ReadOnly<{|
  columnKey: string,
  order: TableSortOrders,
|}>;

export type TableSettings = $ReadOnly<{|
  showSelection: boolean,
  clickableRows: boolean,
  sort?: ?TableSortSettings,
|}>;

export type TableContextValue = $ReadOnly<{|
  settings: TableSettings,
  setSortSettings: (?TableSortSettings) => void,
|}>;

const TableContext = createContext<TableContextValue>({
  settings: {
    showSelection: false,
    clickableRows: false,
  },
  setSortSettings: emptyFunction,
});

export function useTable() {
  return useContext(TableContext);
}

type Props = $ReadOnly<{|
  settings: TableSettings,
  children: React.Node,
|}>;

export function TableContextProvider(props: Props) {
  const {children, settings} = props;
  const [sortSettings, setSortSettings] = useState<?TableSortSettings>(null);
  useEffect(() => setSortSettings(settings.sort), [settings.sort]);

  return (
    <TableContext.Provider
      value={{
        settings: {
          ...settings,
          sort: sortSettings,
        },
        setSortSettings,
      }}>
      {children}
    </TableContext.Provider>
  );
}

export default TableContext;
