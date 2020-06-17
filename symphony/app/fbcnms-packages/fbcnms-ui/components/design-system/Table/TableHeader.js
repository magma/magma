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
import TableHeaderCheckbox from './TableHeaderCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {TABLE_SORT_ORDER, useTable} from './TableContext';
import {makeStyles} from '@material-ui/styles';
import {useCallback} from 'react';
import {useTableCommonStyles} from './TableCommons';

const useStyles = makeStyles(() => ({
  root: {
    backgroundColor: symphony.palette.white,
    borderLeft: `2px solid transparent`,
  },
  cellText: {
    display: 'flex',
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
}));

export type TableColumnType<T> = $ReadOnly<{|
  key: string,
  title: React.Node,
  titleClassName?: ?string,
  render: (rowData: TableRowDataType<T>) => React.Node,
  className?: ?string,
  getSortingValue?: ?(rowData: TableRowDataType<T>) => ?(string | number),
  hidden?: boolean,
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
  const {settings, setSortSettings} = useTable();

  const getSortIcon = useCallback(
    col => {
      if (!col.getSortingValue) {
        return null;
      }

      return (
        <div
          className={classNames(classes.sortIcon, {
            [classes.hidden]: col.key != settings.sort?.columnKey,
          })}>
          {settings.sort?.order === TABLE_SORT_ORDER.descending ? (
            <ArrowUpIcon />
          ) : (
            <ArrowDownIcon />
          )}
        </div>
      );
    },
    [classes.hidden, classes.sortIcon, settings.sort],
  );

  const handleSortChange = useCallback(
    newSortingColumnKey => {
      const newSortingOrder: TableSortOrders =
        settings.sort?.columnKey === newSortingColumnKey &&
        settings.sort?.order === TABLE_SORT_ORDER.ascending
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
    [onSortChanged, setSortSettings, settings.sort],
  );

  return (
    <thead className={classes.root} style={{paddingRight: paddingRight || 0}}>
      <tr>
        {settings.showSelection && (
          <th className={classes.checkBox}>
            <TableHeaderCheckbox />
          </th>
        )}
        {columns
          .filter(col => !col.hidden)
          .map(col => (
            <th
              key={col.key}
              className={classNames(
                commonClasses.cell,
                col.titleClassName,
                cellClassName,
                {
                  [classes.sortableCell]: col.getSortingValue != null,
                },
              )}
              onClick={
                col.getSortingValue != null
                  ? () => handleSortChange(col.key)
                  : undefined
              }>
              <div className={classes.cellContent}>
                <Text className={classes.cellText} variant="body2">
                  {col.title}
                </Text>
                {getSortIcon(col)}
              </div>
            </th>
          ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
