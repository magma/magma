/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TableColumnType} from './Table';

import ArrowDownIcon from '../Icons/ArrowDown';
import ArrowUpIcon from '../Icons/ArrowUp';
import React from 'react';
import TableHeaderCheckbox from './TableHeaderCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useTable} from './TableContext';
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

type Props<T> = {
  columns: Array<TableColumnType<T>>,
  onSortClicked?: (colKey: string) => void,
  cellClassName?: string,
  paddingRight?: ?number,
};

const TableHeader = <T>({
  onSortClicked,
  columns,
  cellClassName,
  paddingRight,
}: Props<T>) => {
  const classes = useStyles();
  const commonClasses = useTableCommonStyles();
  const {showSelection} = useTable();

  const getSortIcon = col => {
    if (!col.sortable) {
      return null;
    }

    return (
      <div
        className={classNames(classes.sortIcon, {
          [classes.hidden]: !Boolean(col.sortDirection),
        })}>
        {col.sortDirection === 'asc' ? <ArrowUpIcon /> : <ArrowDownIcon />}
      </div>
    );
  };

  return (
    <thead className={classes.root} style={{paddingRight: paddingRight || 0}}>
      <tr>
        {showSelection && (
          <th className={classes.checkBox}>
            <TableHeaderCheckbox />
          </th>
        )}
        {columns.map(col => (
          <th
            key={col.key}
            className={classNames(
              commonClasses.cell,
              col.titleClassName,
              cellClassName,
              {
                [classes.sortableCell]: col.sortable,
              },
            )}
            onClick={() =>
              col.sortable && onSortClicked && onSortClicked(col.key)
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
