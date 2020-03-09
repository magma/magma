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
import {useSelection} from './TableSelectionContext';
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
    '&$hoverHighlighting:hover': {
      cursor: 'pointer',
      backgroundColor: symphony.palette.D10,
    },
  },
  bands: {},
  border: {},
  none: {},
  hoverHighlighting: {},
  checkBox: {
    width: '24px',
    paddingLeft: '16px',
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
  const {selectedIds, changeRowSelection} = useSelection();

  return (
    <tbody>
      {data.map((d, rowIndex) => {
        const rowId = d.key ?? rowIndex;
        return (
          <tr
            key={`row_${rowIndex}`}
            onClick={() => {
              changeRowSelection(
                rowId,
                !selectedIds.includes(rowId) ? 'checked' : 'unchecked',
                true,
              );
            }}
            className={classNames(
              classes.row,
              dataRowClassName,
              classes[rowsSeparator],
              {[classes.hoverHighlighting]: showSelection},
            )}>
            {showSelection && (
              <td className={classes.checkBox}>
                <TableRowCheckbox id={rowId} />
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
        );
      })}
    </tbody>
  );
};

export default TableContent;
