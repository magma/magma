/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {RowsSeparationTypes} from './TableContent';
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

const useStyles = makeStyles(() => ({
  table: {
    width: '100%',
    borderCollapse: 'collapse',
  },
  standalone: {
    boxShadow: SymphonyTheme.shadows.DP1,
    borderRadius: '4px',
  },
  embedded: {
    '& $cell': {
      '&:first-child': {
        paddingLeft: '0px',
      },
      '&:last-child': {
        paddingRight: '0px',
      },
    },
  },
  cell: {},
}));

export type TableRowDataType<T> = {key?: string, ...T};

export type TableColumnType<T> = {
  key: string,
  title: React.Node | string,
  render: (rowData: TableRowDataType<T>) => React.Node | string,
  sortable?: boolean,
  sortDirection?: 'asc' | 'desc',
};

export type TableSelectionType = 'all' | 'none' | 'single_item_toggled';

export type TableRowId = string | number;
export type SelectionCallbackType = (
  selectedIds: Array<TableRowId>,
  selection: TableSelectionType,
  toggledItem?: ?{id: TableRowId, change: SelectionType},
) => void;

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
  showSelection?: boolean,
  className?: string,
  variant?: 'standalone' | 'embedded',
  dataRowsSeparator?: RowsSeparationTypes,
  dataRowClassName?: string,
  selectedIds?: Array<TableRowId>,
  onSelectionChanged?: SelectionCallbackType,
  onSortClicked?: (colKey: string) => void,
};

const Table = <T>(props: Props<T>) => {
  const {
    className,
    variant = 'standalone',
    data,
    selectedIds,
    showSelection,
    onSelectionChanged,
    columns,
    onSortClicked,
    dataRowClassName,
    dataRowsSeparator,
  } = props;
  const classes = useStyles();

  const renderChildren = () => (
    <table className={classNames(classes.table, className, classes[variant])}>
      <TableHeader
        columns={columns}
        onSortClicked={onSortClicked}
        cellClassName={classes.cell}
      />
      <TableContent
        columns={columns}
        data={data}
        dataRowClassName={dataRowClassName}
        rowsSeparator={dataRowsSeparator}
        cellClassName={classes.cell}
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
