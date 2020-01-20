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
import SymphonyTheme from '../../../theme/symphony';
import TableHeaderCheckbox from './TableHeaderCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import {makeStyles} from '@material-ui/styles';
import {useTable} from './TableContext';

const useStyles = makeStyles(_theme => ({
  root: {
    backgroundColor: SymphonyTheme.palette.white,
  },
  cell: {
    padding: '4px 8px 4px 16px',
    minHeight: '40px',
    height: '40px',
  },
  cellText: {
    display: 'flex',
    justifyContent: 'flex-start',
    color: SymphonyTheme.palette.D400,
  },
  checkBox: {
    width: '24px',
    paddingLeft: '8px',
  },
  cellContent: {
    display: 'flex',
    alignItems: 'center',
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
      color: SymphonyTheme.palette.primary,
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
};

const TableHeader = <T>({onSortClicked, columns}: Props<T>) => {
  const classes = useStyles();
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
    <thead className={classes.root}>
      <tr>
        {showSelection && (
          <th className={classes.checkBox}>
            <TableHeaderCheckbox />
          </th>
        )}
        {columns.map(col => (
          <th
            key={col.key}
            className={classNames(classes.cell, {
              [classes.sortableCell]: col.sortable,
            })}
            onClick={() =>
              col.sortable && onSortClicked && onSortClicked(col.key)
            }>
            <div className={classes.cellContent}>
              {typeof col.title === 'string' ? (
                <Text className={classes.cellText} variant="body2">
                  {col.title}
                </Text>
              ) : (
                col.title
              )}
              {getSortIcon(col)}
            </div>
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
