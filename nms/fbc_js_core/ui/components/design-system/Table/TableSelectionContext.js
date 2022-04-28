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

import type {
  ActiveCallbackType,
  NullableTableRowId,
  SelectionCallbackType,
  TableRowId,
} from './Table';
import type {SelectionType} from '../Checkbox/Checkbox';

import * as React from 'react';
import emptyFunction from '../../../../../fbc_js_core/util/emptyFunction';
import {useContext, useMemo} from 'react';

export type TableSelectionContextValue = {
  activeId: string | number | null,
  setActiveId?: ?(id: NullableTableRowId) => void,
  selectedIds: $ReadOnlyArray<TableRowId>,
  selectionMode: 'all' | 'none' | 'some',
  changeRowSelection: (
    id: TableRowId,
    selection: SelectionType,
    isExclusive?: boolean,
  ) => void,
  changeHeaderSelectionMode: (selection: SelectionType) => void,
};

const TableSelectionContext = React.createContext<TableSelectionContextValue>({
  activeId: null,
  setActiveId: emptyFunction,
  selectedIds: [],
  selectionMode: 'none',
  changeRowSelection: emptyFunction,
  changeHeaderSelectionMode: emptyFunction,
});

type Props = {
  children: React.Node,
  allIds: $ReadOnlyArray<TableRowId>,
  activeId?: NullableTableRowId,
  onActiveChanged?: ActiveCallbackType,
  selectedIds: $ReadOnlyArray<TableRowId>,
  onSelectionChanged?: SelectionCallbackType,
};

export const TableSelectionContextProvider = ({
  activeId = null,
  selectedIds,
  allIds,
  children,
  onSelectionChanged,
  onActiveChanged,
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
        activeId: activeId,
        setActiveId: onActiveChanged,
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
