/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TableColumnType} from './Table.react';

import React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import TableHeaderCheckbox from './TableHeaderCheckbox.react';
import Text from '../../Text.react';
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
}));

type Props<T> = {
  columns: Array<TableColumnType<T>>,
};

const TableHeader = <T>(props: Props<T>) => {
  const {columns} = props;
  const classes = useStyles();
  const {showSelection} = useTable();
  return (
    <thead className={classes.root}>
      <tr>
        {showSelection && (
          <th className={classes.checkBox}>
            <TableHeaderCheckbox />
          </th>
        )}
        {columns.map((col, i) => (
          <th key={`col_${i}`} className={classes.cell}>
            {typeof col.title === 'string' ? (
              <Text className={classes.cellText} variant="body2">
                {col.title}
              </Text>
            ) : (
              col.title
            )}
          </th>
        ))}
      </tr>
    </thead>
  );
};

export default TableHeader;
