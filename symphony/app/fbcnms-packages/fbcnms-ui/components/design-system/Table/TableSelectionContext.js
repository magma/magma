/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {SelectionCallbackType} from './Table';
import type {SelectionType} from '../Checkbox/Checkbox';

import * as React from 'react';
import emptyFunction from '../../../../fbcnms-util/emptyFunction';
import {useContext, useMemo} from 'react';

export type TableSelectionContextValue = {
  selectedIds: Array<string | number>,
  selectionMode: 'all' | 'none' | 'some',
  changeRowSelection: (
    id: string | number,
    selection: SelectionType,
    isExclusive?: boolean,
  ) => void,
  changeHeaderSelectionMode: (selection: SelectionType) => void,
};

const TableSelectionContext = React.createContext<TableSelectionContextValue>({
  selectedIds: [],
  selectionMode: 'none',
  changeRowSelection: emptyFunction,
  changeHeaderSelectionMode: emptyFunction,
});

type Props = {
  children: React.Node,
  allIds: Array<string | number>,
  selectedIds: Array<string | number>,
  onSelectionChanged?: SelectionCallbackType,
};

export const TableSelectionContextProvider = ({
  selectedIds,
  allIds,
  children,
  onSelectionChanged,
}: Props) => {
  const selectionMode = useMemo(() => {
    if (selectedIds.length === 0) {
      return 'none';
    }

    return allIds.every(id => selectedIds.includes(id)) ? 'all' : 'some';
  }, [allIds, selectedIds]);
  return (
    <TableSelectionContext.Provider
      value={{
        selectedIds: selectedIds ?? [],
        changeRowSelection: (id, selection, isExclusive) => {
          if (!onSelectionChanged) {
            return;
          }
          const newTableSelection =
            isExclusive === true
              ? [id]
              : selection === 'unchecked'
              ? selectedIds.filter(idItem => idItem !== id)
              : [...selectedIds, id];
          onSelectionChanged(newTableSelection, 'single_item_toggled', {
            id,
            change: selection,
          });
        },
        changeHeaderSelectionMode: selection => {
          onSelectionChanged &&
            onSelectionChanged(
              selection === 'checked' ? [...allIds] : [],
              selection === 'checked' ? 'all' : 'none',
            );
        },
        selectionMode,
      }}>
      {children}
    </TableSelectionContext.Provider>
  );
};

export function useSelection() {
  return useContext(TableSelectionContext);
}

export default TableSelectionContext;
