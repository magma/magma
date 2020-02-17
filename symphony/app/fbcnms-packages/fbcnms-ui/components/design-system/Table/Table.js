/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {SelectionType} from '../Checkbox/Checkbox';

import * as React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import TableContent from './TableContent';
import TableContext from './TableContext';
import TableHeader from './TableHeader';
import classNames from 'classnames';
import {TableSelectionContextProvider} from './TableSelectionContext';
import {makeStyles} from '@material-ui/styles';
import {useMemo} from 'react';

const useStyles = makeStyles(_theme => ({
  table: {
    width: '100%',
    boxShadow: SymphonyTheme.shadows.DP1,
    borderRadius: '4px',
    borderCollapse: 'collapse',
  },
}));

export type TableRowDataType<T> = {key?: string} & T;

export type TableColumnType<T> = {
  key: string,
  title: React.Node | string,
  render: (rowData: TableRowDataType<T>) => React.Node | string,
  sortable?: boolean,
  sortDirection?: 'asc' | 'desc',
};

export type TableSelectionType = 'all' | 'none' | 'single_item_toggled';

export type SelectionCallbackType = (
  selectedIds: Array<string | number>,
  selection: TableSelectionType,
  toggledItem?: ?{id: string | number, change: SelectionType},
) => void;

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
  showSelection?: boolean,
  className?: string,
  dataRowClassName?: string,
  selectedIds?: Array<string | number>,
  onSelectionChanged?: SelectionCallbackType,
  onSortClicked?: (colKey: string) => void,
};

const Table = <T>(props: Props<T>) => {
  const {
    className,
    data,
    selectedIds,
    showSelection,
    onSelectionChanged,
    columns,
    onSortClicked,
  } = props;
  const classes = useStyles();

  const renderChildren = () => (
    <table className={classNames(classes.table, className)}>
      <TableHeader columns={columns} onSortClicked={onSortClicked} />
      <TableContent
        columns={columns}
        data={data}
        dataRowClassName={props.dataRowClassName}
      />
    </table>
  );

  const allIds = useMemo(() => data.map((d, i) => d.key ?? i), [data]);
  return (
    <TableContext.Provider value={{showSelection: showSelection ?? false}}>
      {showSelection ? (
        <TableSelectionContextProvider
          allIds={allIds}
          selectedIds={selectedIds ?? []}
          onSelectionChanged={onSelectionChanged}>
          {renderChildren()}
        </TableSelectionContextProvider>
      ) : (
        renderChildren()
      )}
    </TableContext.Provider>
  );
};

export default Table;
