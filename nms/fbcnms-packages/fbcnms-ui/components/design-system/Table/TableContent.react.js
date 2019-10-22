/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {TableColumnType, TableRowDataType} from './Table.react';

import React from 'react';
import SymphonyTheme from '../../../theme/symphony';
import Text from '../../Text.react';
import {makeStyles} from '@material-ui/styles';

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
}));

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
};

const TableContent = <T>(props: Props<T>) => {
  const {data, columns} = props;
  const classes = useStyles();

  return (
    <tbody>
      {data.map((d, i) => (
        <tr key={`row_${i}`} className={classes.row}>
          {columns.map((col, i) => {
            const renderedCol = col.render(d);
            return (
              <td key={d.id ?? `col_${i}`} className={classes.cell}>
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
