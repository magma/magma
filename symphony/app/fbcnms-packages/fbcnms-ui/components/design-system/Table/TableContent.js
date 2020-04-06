/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {TRefFor} from '@fbcnms/ui/components/design-system/types/TRefFor.flow.js';
import type {TableColumnType, TableRowDataType} from './Table';

import * as React from 'react';
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
    borderLeft: `2px solid transparent`,
    '&$bands:nth-child(odd)': {
      backgroundColor: symphony.palette.background,
    },
    '&$border': {
      borderBottom: `1px solid ${symphony.palette.separatorLight}`,
    },
    '&$hoverHighlighting:hover': {
      cursor: 'pointer',
      '&$border': {
        backgroundColor: symphony.palette.D10,
      },
      '&$bands': {
        '& #column0 $textualCell': {
          color: symphony.palette.primary,
        },
      },
      '& $checkBox': {
        opacity: 1,
      },
    },
  },
  activeRow: {
    borderLeft: `2px solid ${symphony.palette.primary}`,
    '&:not($bands)': {
      backgroundColor: symphony.palette.D10,
    },
  },
  bands: {},
  border: {},
  none: {},
  hoverHighlighting: {},
  checkBox: {
    width: '28px',
    paddingLeft: '12px',
  },
  textualCell: {},
}));

export const ROW_SEPARATOR_TYPES = {
  bands: 'bands',
  border: 'border',
  none: 'none',
};
export type RowsSeparationTypes = $Keys<typeof ROW_SEPARATOR_TYPES>;

type Props<T> = {
  data: Array<TableRowDataType<T>>,
  columns: Array<TableColumnType<T>>,
  rowsSeparator?: RowsSeparationTypes,
  dataRowClassName?: string,
  cellClassName?: string,
  fwdRef?: TRefFor<HTMLElement>,
};

const TableContent = <T>(props: Props<T>) => {
  const {
    data,
    columns,
    dataRowClassName,
    cellClassName,
    rowsSeparator = ROW_SEPARATOR_TYPES.bands,
    fwdRef,
  } = props;
  const classes = useStyles();
  const commonClasses = useTableCommonStyles();
  const {showSelection, clickableRows} = useTable();
  const {activeId, setActiveId} = useSelection();

  return (
    <tbody ref={fwdRef}>
      {data.map((d, rowIndex) => {
        const rowId = d.key ?? rowIndex;
        return (
          <tr
            key={`row_${rowIndex}`}
            onClick={() => {
              if (setActiveId == null) {
                return;
              }
              const newActiveId = rowId !== activeId ? rowId : null;
              setActiveId(newActiveId);
            }}
            className={classNames(
              classes.row,
              dataRowClassName,
              classes[rowsSeparator],
              {
                [classes.hoverHighlighting]: clickableRows,
                [classes.activeRow]: rowId === activeId,
              },
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
                  id={`column${colIndex}`}
                  className={classNames(
                    commonClasses.cell,
                    col.className,
                    cellClassName,
                  )}>
                  <Text
                    className={classes.textualCell}
                    useEllipsis={true}
                    variant="body2">
                    {renderedCol}
                  </Text>
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
