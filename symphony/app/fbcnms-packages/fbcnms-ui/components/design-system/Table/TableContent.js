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
import TableRowCheckbox from './TableRowCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {makeStyles} from '@material-ui/styles';
import {useTable} from './TableContext';
import {useTableCommonStyles} from './TableCommons';

const useStyles = makeStyles(() => ({
  row: {
    backgroundColor: symphony.palette.white,
    '&$bands:nth-child(odd)': {
      backgroundColor: symphony.palette.background,
    },
    '&$border': {
      borderTop: '1px solid',
      borderColor: symphony.palette.separatorLight,
    },
  },
  bands: {},
  border: {},
  none: {},
  checkBox: {
    width: '24px',
    paddingLeft: '8px',
  },
}));

export type RowsSeparationTypes = 'bands' | 'border' | 'none';

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
  rowsSeparator?: RowsSeparationTypes,
  dataRowClassName?: string,
  cellClassName?: string,
};

const TableContent = <T>(props: Props<T>) => {
  const {
    data,
    columns,
    dataRowClassName,
    cellClassName,
    rowsSeparator = 'bands',
  } = props;
  const classes = useStyles();
  const commonClasses = useTableCommonStyles();
  const {showSelection} = useTable();

  return (
    <tbody>
      {data.map((d, rowIndex) => (
        <tr
          key={`row_${rowIndex}`}
          className={classNames(
            classes.row,
            dataRowClassName,
            classes[rowsSeparator],
          )}>
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
                className={classNames(commonClasses.cell, cellClassName)}>
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
