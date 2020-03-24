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
import symphony from '../../../theme/symphony';
import useVerticalScrollingEffect from '../hooks/useVerticalScrollingEffect';
import {TableSelectionContextProvider} from './TableSelectionContext';
import {makeStyles} from '@material-ui/styles';
import {useEffect, useMemo, useRef, useState} from 'react';

const borderRadius = '4px';
const useStyles = makeStyles(() => ({
  root: {
    display: 'flex',
  },
  standalone: {
    borderRadius: borderRadius,
    boxShadow: SymphonyTheme.shadows.DP1,
  },
  tableContainer: {
    display: 'flex',
    maxHeight: '100%',
    overflow: 'hidden',
    flexGrow: 0,
    flexBasis: '25%',
    borderRadius: borderRadius,
    '&$expanded': {
      flexGrow: 1,
    },
  },
  expanded: {},
  table: {
    display: 'flex',
    flexDirection: 'column',
    borderCollapse: 'collapse',
    overflow: 'hidden',
    '& tbody': {
      borderTop: `1px solid ${symphony.palette.separatorLight}`,
      overflowX: 'hidden',
      overflowY: 'auto',
    },
    '& tr': {
      tableLayout: 'fixed',
      display: 'table',
      width: '100%',
    },
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
  detailsCardContainer: {
    backgroundColor: SymphonyTheme.palette.white,
    borderLeft: '1px solid',
    borderColor: SymphonyTheme.palette.separatorLight,
    flexBasis: '10px',
    flexGrow: 1,
    borderTopRightRadius: borderRadius,
    borderBottomRightRadius: borderRadius,
  },
}));

export type TableRowDataType<T> = {key?: string, ...T};

export type TableColumnType<T> = {
  key: string,
  title: React.Node | string,
  titleClassName?: ?string,
  render: (rowData: TableRowDataType<T>) => React.Node | string,
  className?: ?string,
  sortable?: boolean,
  sortDirection?: 'asc' | 'desc',
};

export type TableSelectionType = 'all' | 'none' | 'single_item_toggled';

export type TableRowId = string | number;
export type NullableTableRowId = TableRowId | null;

export type SelectionCallbackType = (
  selectedIds: Array<TableRowId>,
  selection: TableSelectionType,
  toggledItem?: ?{id: TableRowId, change: SelectionType},
) => void;
export type ActiveCallbackType = (activeId: NullableTableRowId) => void;

/*
  detailsCard:
    When passed, will be shown on as part of the table content.
    Excepts for the first column, all columns will get hidden.
    The card will cover 75% of the table width.
*/
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
  activeRowId?: NullableTableRowId,
  onActiveRowIdChanged?: ActiveCallbackType,
  onSortClicked?: (colKey: string) => void,
  detailsCard?: ?React.Node,
};

const Table = <T>(props: Props<T>) => {
  const {
    className,
    variant = 'standalone',
    data,
    showSelection,
    activeRowId,
    onActiveRowIdChanged,
    selectedIds = [],
    onSelectionChanged,
    columns,
    onSortClicked,
    dataRowClassName,
    dataRowsSeparator,
    detailsCard,
  } = props;
  const classes = useStyles();
  const [dataColumns, setDataColumns] = useState([]);

  const [tableHeaderPaddingRight, setTableHeaderPaddingRight] = useState(0);
  const bodyRef = useRef(null);
  useVerticalScrollingEffect(
    bodyRef,
    scrollArgs => setTableHeaderPaddingRight(scrollArgs.scrollbarWidth),
    false,
  );

  useEffect(() => {
    setDataColumns(detailsCard == null ? columns : [columns[0]]);
  }, [detailsCard, columns]);

  const renderChildren = () => (
    <div className={classNames(classes.root, classes[variant])}>
      <div
        className={classNames(classes.tableContainer, {
          [classes.expanded]: !detailsCard,
        })}>
        <table className={classNames(classes.table, className)}>
          <TableHeader
            columns={dataColumns}
            onSortClicked={onSortClicked}
            cellClassName={classes.cell}
            paddingRight={tableHeaderPaddingRight}
          />
          <TableContent
            columns={dataColumns}
            data={data}
            dataRowClassName={dataRowClassName}
            rowsSeparator={dataRowsSeparator}
            cellClassName={classes.cell}
            fwdRef={bodyRef}
          />
        </table>
      </div>
      {detailsCard ? (
        <div className={classes.detailsCardContainer}>{detailsCard}</div>
      ) : null}
    </div>
  );

  const allIds = useMemo(() => data.map((d, i) => d.key ?? i), [data]);
  return (
    <TableContext.Provider value={{showSelection: showSelection ?? false}}>
      {showSelection ? (
        <TableSelectionContextProvider
          allIds={allIds}
          activeId={activeRowId}
          onActiveChanged={onActiveRowIdChanged}
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
