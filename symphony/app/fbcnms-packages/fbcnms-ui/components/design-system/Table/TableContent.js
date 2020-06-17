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
import type {TableColumnType} from './TableHeader';
import type {TableRowDataType} from './Table';

import * as React from 'react';
import TableRowCheckbox from './TableRowCheckbox';
import Text from '../Text';
import classNames from 'classnames';
import symphony from '../../../theme/symphony';
import {TABLE_SORT_ORDER, useTable} from './TableContext';
import {makeStyles} from '@material-ui/styles';
import {sortMixed} from '../../../utils/displayUtils';
import {useEffect, useState} from 'react';
import {useSelection} from './TableSelectionContext';
import {useTableCommonStyles} from './TableCommons';

const useStyles = makeStyles(() => ({
  row: {
    backgroundColor: symphony.palette.white,
    borderLeft: `2px solid transparent`,
    '&$bands:nth-child(odd)': {
      backgroundColor: symphony.palette.background,
    },
    '&$border:not(:last-child)': {
      borderBottom: `1px solid ${symphony.palette.separatorLight}`,
    },
    '&$hoverHighlighting:hover:not($disabled)': {
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
  disabled: {},
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
    paddingTop: '7px',
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

type Props<T> = $ReadOnly<{|
  data: $ReadOnlyArray<TableRowDataType<T>>,
  columns: $ReadOnlyArray<TableColumnType<T>>,
  rowsSeparator?: RowsSeparationTypes,
  dataRowClassName?: string,
  cellClassName?: string,
  fwdRef?: TRefFor<HTMLElement>,
|}>;

const requiredOnTopValue = row => (row.alwaysShowOnTop === true ? 1 : 0);

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
  const {settings} = useTable();
  const {activeId, setActiveId} = useSelection();

  const [sortedData, setSortedData] = useState<Array<TableRowDataType<T>>>([]);

  useEffect(() => {
    const sortSettings = settings.sort;
    const sortingColumn =
      sortSettings && columns.find(col => col.key == sortSettings.columnKey);
    const getSortingValue = sortingColumn?.getSortingValue;
    const sortingFactor =
      sortSettings?.order === TABLE_SORT_ORDER.descending ? -1 : 1;

    setSortedData(
      data.slice().sort((row1, row2) => {
        const topRowsSortingValue =
          requiredOnTopValue(row2) - requiredOnTopValue(row1);
        if (topRowsSortingValue !== 0) {
          return topRowsSortingValue;
        }
        if (getSortingValue != null) {
          return (
            sortMixed(getSortingValue(row1), getSortingValue(row2)) *
            sortingFactor
          );
        }
        return 0;
      }),
    );
  }, [columns, data, settings.sort]);

  return (
    <tbody ref={fwdRef}>
      {sortedData.map((d, rowIndex) => {
        const rowId = d.key ?? rowIndex;
        return (
          <tr
            key={`row_${rowIndex}`}
            title={d.tooltip}
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
              d.className,
              classes[rowsSeparator],
              {
                [classes.hoverHighlighting]: settings.clickableRows,
                [classes.activeRow]: rowId === activeId,
                [classes.disabled]: d.disabled,
              },
            )}>
            {settings.showSelection && (
              <td className={classes.checkBox}>
                {d.disabled !== true ? <TableRowCheckbox id={rowId} /> : null}
              </td>
            )}
            {columns
              .filter(col => !col.hidden)
              .map((col, colIndex) => {
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
                      color="inherit"
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
