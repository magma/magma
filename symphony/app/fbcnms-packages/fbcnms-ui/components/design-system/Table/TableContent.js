/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TableColumnType, TableRowDataType} from './Table';

import React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import TableRowCheckbox from './TableRowCheckbox';
import Text from '../Text';
import {makeStyles} from '@material-ui/styles';
import {useTable} from './TableContext';

const useStyles = makeStyles(_theme => ({
  row: {
    backgroundColor: SymphonyTheme.palette.background,
    '&:nth-child(even)': {
      backgroundColor: SymphonyTheme.palette.white,
    },
  },
  cell: {
    padding: '4px 8px 4px 16px',
    minHeight: '40px',
    height: '40px',
  },
  checkBox: {
    width: '24px',
    paddingLeft: '8px',
  },
}));

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
};

const TableContent = <T>(props: Props<T>) => {
  const {data, columns} = props;
  const classes = useStyles();
  const {showSelection} = useTable();

  return (
    <tbody>
      {data.map((d, rowIndex) => (
        <tr key={`row_${rowIndex}`} className={classes.row}>
          {showSelection && (
            <td className={classes.checkBox}>
              <TableRowCheckbox id={d.key ?? rowIndex} />
            </td>
          )}
          {columns.map((col, colIndex) => {
            const renderedCol = col.render(d);
            return (
              <td
                key={`col_${colIndex}_${d.key ?? rowIndex}`}
                className={classes.cell}>
                {typeof renderedCol === 'string' ? (
                  <Text variant="body2">{renderedCol}</Text>
                ) : (
                  renderedCol
                )}
              </td>
            );
          })}
        </tr>
      ))}
    </tbody>
  );
};

export default TableContent;
