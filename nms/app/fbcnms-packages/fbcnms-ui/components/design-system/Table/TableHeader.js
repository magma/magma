/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TableRowDataType} from './Table';
import type {TableSortOrders, TableSortSettings} from './TableContext';

import * as React from 'react';
import ArrowDownIcon from '../Icons/ArrowDown';
import ArrowUpIcon from '../Icons/ArrowUp';
import Draggable from 'react-draggable';
import TableHeaderCheckbox from './TableHeaderCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {TABLE_SORT_ORDER, useTable} from './TableContext';
import {makeStyles} from '@material-ui/styles';
import {useCallback, useMemo, useState} from 'react';
import {useTableCommonStyles} from './TableCommons';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    borderLeft: `2px solid transparent`,
  },
  cellText: {
    justifyContent: 'flex-start',
    color: symphony.palette.D400,
  },
  checkBox: {
    width: '28px',
    paddingLeft: '12px',
  },
  cellContent: {
    display: 'flex',
    alignItems: 'center',
    color: symphony.palette.D400,
    overflow: 'hidden',
    ...symphony.typography.body2,
  },
  sortIcon: {
    width: '24px',
    height: '24px',
    display: 'flex',
    justifyContent: 'center',
    alignItems: 'center',
  },
  sortableCell: {
    cursor: 'pointer',
    '&:hover $cellText': {
      color: symphony.palette.primary,
    },
    '&:hover $hidden': {
      visibility: 'visible',
    },
  },
  hidden: {
    visibility: 'hidden',
  },
  dragHandleArea: {
    willChange: 'transform',
    position: 'absolute',
    right: 0,
    bottom: 8,
    top: 8,
    width: '4px',
    zIndex: 2,
    cursor: 'ew-resize',
    display: 'flex',
    justifyContent: 'center',
  },
  isOnColumnResize: {},
  dragHandle: {
    backgroundColor: 'transparent',
    width: '2px',
    height: '100%',
    zIndex: 2,
  },
  dragHandleAreaActive: {
    '& $dragHandle': {
      backgroundColor: symphony.palette.primary,
    },
  },
  headerCell: {
    position: 'relative',
    boxSizing: 'border-box',
    '&$isOnColumnResize': {
      cursor: 'ew-resize',
    },
    '&:not($isOnColumnResize)': {
      '&:hover $dragHandle, $dragHandleAreaActive $dragHandle': {
        backgroundColor: symphony.palette.primary,
      },
    },
  },
}));

export type TableColumnType<T> = $ReadOnly<{|
  key: string,
  title: React.Node,
  titleClassName?: ?string,
  render: (rowData: TableRowDataType<T>) => React.Node,
  tooltip?: ?(rowData: TableRowDataType<T>) => ?(string | number),
  className?: ?string,
  getSortingValue?: ?(rowData: TableRowDataType<T>) => ?(string | number),
  hidden?: boolean,
  /* either pixels or ratio (e.g. 0.33) */
  width?: number,
|}>;

export type TableHeaderData<T> = $ReadOnly<{|
  columns: Array<TableColumnType<T>>,
  onSortChanged?: ?(newSortSettings: TableSortSettings) => void,
|}>;

type Props<T> = $ReadOnly<{|
  ...TableHeaderData<T>,
  cellClassName?: string,
  paddingRight?: ?number,
|}>;

const TableHeader = <T>({
  columns,
  onSortChanged,
  cellClassName,
  paddingRight,
}: Props<T>) => {
  const classes = useStyles();
  const commonClasses = useTableCommonStyles();
  const {
    settings: {showSelection, sort, columnWidths, resizableColumns},
    setSortSettings,
    changeColumnWidthByDelta,
    width: tableWidth,
  } = useTable();
  const [draggingColumnKey, setDraggingColumnKey] = useState(null);

  const shownColumns = useMemo(() => columns.filter(col => !col.hidden), [
    columns,
  ]);

  const getSortIcon = useCallback(
    col => {
      if (!col.getSortingValue) {
        return null;
      }

      return (
        <div
          className={classNames(classes.sortIcon, {
            [classes.hidden]: col.key != sort?.columnKey,
          })}>
          {sort?.order === TABLE_SORT_ORDER.descending ? (
            <ArrowUpIcon />
          ) : (
            <ArrowDownIcon />
          )}
        </div>
      );
    },
    [classes.hidden, classes.sortIcon, sort],
  );

  const handleSortChange = useCallback(
    newSortingColumnKey => {
      const newSortingOrder: TableSortOrders =
        sort?.columnKey === newSortingColumnKey &&
        sort?.order === TABLE_SORT_ORDER.ascending
          ? 'descending'
          : TABLE_SORT_ORDER.ascending;
      const newSortSettings = {
        columnKey: newSortingColumnKey,
        order: newSortingOrder,
      };
      setSortSettings(newSortSettings);
      if (onSortChanged) {
        onSortChanged(newSortSettings);
      }
    },
    [onSortChanged, setSortSettings, sort],
  );

  return (
    <thead className={classes.root} style={{paddingRight: paddingRight || 0}}>
      <tr>
        {showSelection && (
          <th className={classes.checkBox}>
            <TableHeaderCheckbox />
          </th>
        )}
        {shownColumns.map((col, index) => (
          <th
            key={col.key}
            className={classNames(
              commonClasses.cell,
              classes.headerCell,
              col.titleClassName,
              cellClassName,
              {
                [classes.sortableCell]: col.getSortingValue != null,
                [classes.isOnColumnResize]: draggingColumnKey != null,
              },
            )}
            style={{
              width:
                tableWidth != null && columnWidths
                  ? columnWidths[index].width
                  : undefined,
            }}>
            <div
              className={classes.cellContent}
              onClick={
                col.getSortingValue != null
                  ? () => handleSortChange(col.key)
                  : undefined
              }>
              <Text
                className={classes.cellText}
                variant="body2"
                useEllipsis={true}>
                {col.title}
              </Text>
              {getSortIcon(col)}
            </div>
            {resizableColumns && index !== shownColumns.length - 1 && (
              <Draggable
                defaultClassName={classNames(classes.dragHandleArea)}
                defaultClassNameDragging={classes.dragHandleAreaActive}
                axis="x"
                position={{x: 0}}
                onStart={() => setDraggingColumnKey(col.key)}
                onStop={() => setDraggingColumnKey(null)}
                onDrag={(event, {deltaX, node}) => {
                  const rect = node.getBoundingClientRect();
                  if (
                    deltaX === 0 ||
                    (deltaX < 0 && event.clientX > rect.left) ||
                    (deltaX > 0 && event.clientX < rect.right)
                  ) {
                    return;
                  }
                  changeColumnWidthByDelta(index, deltaX);
                }}>
                <div>
                  <div className={classes.dragHandle} />
                </div>
              </Draggable>
            )}
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
