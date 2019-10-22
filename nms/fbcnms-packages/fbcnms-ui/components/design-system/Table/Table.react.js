/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import * as React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import TableContent from './TableContent.react';
import TableContext from './TableContext';
import TableHeader from './TableHeader.react';
import {makeStyles} from '@material-ui/styles';

const useStyles = makeStyles(_theme => ({
  table: {
    width: '100%',
    boxShadow: SymphonyTheme.shadows.DP1,
    borderRadius: '4px',
    borderCollapse: 'collapse',
  },
}));

export type TableRowDataType<T> = {id?: string} & T;

export type TableColumnType<T> = {
  title: React.Node | string,
  render: (rowData: TableRowDataType<T>) => React.Node,
};

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
};

const Table = <T>(props: Props<T>) => {
  const {columns, data} = props;
  const classes = useStyles();
  return (
    <TableContext.Provider value={{}}>
      <table className={classes.table}>
        <TableHeader columns={columns} />
        <TableContent columns={columns} data={data} />
      </table>
    </TableContext.Provider>
  );
};

export default Table;
