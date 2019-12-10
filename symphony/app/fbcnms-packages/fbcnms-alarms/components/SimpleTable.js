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
import Button from '@material-ui/core/Button';
import Chip from '@material-ui/core/Chip';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import grey from '@material-ui/core/colors/grey';
import orange from '@material-ui/core/colors/orange';
import red from '@material-ui/core/colors/red';
import yellow from '@material-ui/core/colors/yellow';
import {makeStyles} from '@material-ui/styles';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  body: {
    padding: theme.spacing(3),
  },
  labelChip: {
    backgroundColor: theme.palette.grey[50],
    color: theme.palette.secondary.main,
    margin: '5px',
  },
  titleCell: {
    fontSize: 18,
    fontWeight: 500,
    marginBottom: 2,
  },
  secondaryCell: {
    color: theme.palette.text.secondary,
  },
  secondaryItalicCell: {
    fontStyle: 'italic',
    color: theme.palette.text.secondary,
  },
  secondaryChip: {
    color: theme.palette.secondary.main,
  },
  redSeverityChip: {
    color: theme.palette.secondary.main,
    border: `1px solid ${red.A400}`,
  },
  orangeSeverityChip: {
    color: theme.palette.secondary.main,
    border: `1px solid ${orange.A400}`,
  },
  yellowSeverityChip: {
    color: theme.palette.secondary.main,
    border: `1px solid ${yellow.A400}`,
  },
  greySeverityChip: {
    color: theme.palette.secondary.main,
    border: `1px solid ${grey[500]}`,
  },
  ellipsisChip: {
    display: 'block',
    maxWidth: 256,
    overflow: 'hidden',
    textOverflow: 'ellipsis',
    whiteSpace: 'nowrap',
  },
}));

const HeadTableCell = withStyles({
  root: {
    borderBottom: 'none',
    fontSize: '14px',
    color: 'black',
    textTransform: 'uppercase',
  },
})(TableCell);

const BodyTableCell = withStyles({
  root: {
    borderBottom: 'none',
  },
})(TableCell);

const SEVERITY = {
  critical: {style: 'redSeverityChip'},
  major: {style: 'orangeSeverityChip'},
  minor: {style: 'yellowSeverityChip'},
  warning: {style: 'yellowSeverityChip'},
  info: {style: 'greySeverityChip'},
  notice: {style: 'greySeverityChip'},
};

type RenderCellProps<TRow> = {
  row: TRow,
  rowIdx: number,
  column: ColumnData<TRow>,
  columnIdx: number,
  classes: {[string]: string},
};

function RenderCell<TRow>({
  row,
  rowIdx,
  column,
  columnIdx,
  classes,
}: RenderCellProps<TRow>) {
  const commonProps = {
    classes: classes,
    columnIdx: columnIdx,
    cellKey: `${rowIdx}_${columnIdx}`,
  };

  if (typeof column.renderFunc === 'function') {
    return (
      <CustomCell {...commonProps} value={column.renderFunc(row, classes)} />
    );
  } else {
    /**
     * Since column.render is the discriminator the ColumnData disjoint union,
     * getValue needs to be called individually inside each conditional to work
     * properly. Flow can't know which type getValue will return until after its
     * type has been speciallized by checking against the column.render
     * property.
     */
    if (column.render === 'severity') {
      return <SeverityCell {...commonProps} value={column.getValue(row)} />;
    } else if (column.render === 'multipleGroups') {
      return <MultiGroupsCell {...commonProps} value={column.getValue(row)} />;
    } else if (column.render === 'chip') {
      return <ChipCell {...commonProps} value={column.getValue(row)} />;
    } else if (column.render === 'labels') {
      return (
        <LabelsCell
          {...commonProps}
          value={column.getValue(row)}
          hideFields={column.hideFields}
        />
      );
    } else if (column.render === 'list') {
      return (
        <TextCell {...commonProps} value={column.getValue(row).join(', ')} />
      );
    } else {
      return <TextCell {...commonProps} value={column.getValue(row)} />;
    }
  }
}

function CustomCell({value}: CellProps<React.Node>) {
  return (
    <BodyTableCell>
      <div>{value}</div>
    </BodyTableCell>
  );
}

export type CellProps<TValue> = {
  value: TValue,
  cellKey: string,
  columnIdx: number,
  classes: {[string]: string},
};

type GroupsList = Array<Labels>;

function MultiGroupsCell({value, classes, columnIdx}: CellProps<GroupsList>) {
  return (
    <BodyTableCell>
      {value.map((cellValue, idx) => (
        <div
          key={idx}
          className={
            columnIdx === 0 ? classes.titleCell : classes.secondaryCell
          }>
          {Object.keys(cellValue).map(keyName => (
            <Chip
              key={keyName}
              classes={{label: classes.ellipsisChip}}
              className={classes.labelChip}
              label={
                <span>
                  <em>{keyName}</em>={renderLabelValue(cellValue[keyName])}
                </span>
              }
              size="small"
            />
          ))}
        </div>
      ))}
    </BodyTableCell>
  );
}

const renderLabelValue = labelValue => {
  if (typeof labelValue === 'boolean') {
    return labelValue ? 'true' : 'false';
  }
  return labelValue;
};

type Labels = {[string]: string | number | boolean};
function LabelsCell({
  value,
  classes,
  columnIdx,
  hideFields,
}: CellProps<Labels> & {hideFields?: Array<string>}) {
  const labels = React.useMemo(() => {
    if (!hideFields) {
      return value;
    }
    const filtered: Labels = {...value};
    // filter out all keys which are in the hideFields list
    hideFields.forEach(key => delete filtered[key]);
    return filtered;
  }, [value, hideFields]);
  return (
    <BodyTableCell>
      <div
        className={
          columnIdx === 0 ? classes.titleCell : classes.secondaryItalicCell
        }>
        {Object.keys(labels).map(keyName => (
          <Chip
            key={keyName}
            classes={{label: classes.ellipsisChip}}
            className={classes.labelChip}
            label={
              <span>
                <em>{keyName}</em>={renderLabelValue(labels[keyName])}
              </span>
            }
            size="small"
          />
        ))}
      </div>
    </BodyTableCell>
  );
}

function TextCell({value, classes, columnIdx}: CellProps<string>) {
  return (
    <BodyTableCell>
      <div
        className={
          columnIdx === 0 ? classes.titleCell : classes.secondaryItalicCell
        }>
        {value}
      </div>
    </BodyTableCell>
  );
}

function SeverityCell({value, classes}: CellProps<string>) {
  return (
    <BodyTableCell>
      {value && value.toLowerCase() in SEVERITY && (
        <Chip
          classes={{
            outlined: classes[SEVERITY[value.toLowerCase()].style],
            label: classes.ellipsisChip,
          }}
          label={value.toUpperCase()}
          variant="outlined"
          data-severity={value} // for testing
        />
      )}
    </BodyTableCell>
  );
}

function ChipCell({value, classes}: CellProps<string>) {
  return (
    <BodyTableCell>
      {value && (
        <Chip
          classes={{outlinedPrimary: classes.secondaryChip}}
          label={value.toUpperCase()}
          color="primary"
          variant="outlined"
          data-chip={value} // for testing
        />
      )}
    </BodyTableCell>
  );
}

type CommonColumnProps<TRow> = {|
  title: string,
  // DEPRECATED - use getValue instead
  path?: Array<string>,
  hideFields?: Array<string>,
  renderFunc?: (tableRow: TRow, classes: {[string]: string}) => React.Node,
  tooltip?: React.Node,
|};

// build up a disjoint union to handle all the renderer types
export type ColumnData<TRow> =
  | {
      getValue: (row: TRow) => Array<Labels>,
      render: 'multipleGroups',
      ...CommonColumnProps<TRow>,
    }
  | {
      getValue: (row: TRow) => Labels,
      render: 'labels',
      ...CommonColumnProps<TRow>,
    }
  | {
      getValue: (row: TRow) => string,
      render?: '' | 'chip',
      ...CommonColumnProps<TRow>,
    }
  | {
      getValue: (row: TRow) => string,
      render: 'severity',
      ...CommonColumnProps<TRow>,
    }
  | {
      getValue: (row: TRow) => Array<string>,
      render: 'list',
      ...CommonColumnProps<TRow>,
    };

type Props<TRow> = {
  columnStruct: Array<ColumnData<TRow>>,
  tableData: Array<TRow>,
  onActionsClick?: (row: TRow, target: HTMLElement) => void,
  sortFunc?: (row1: TRow, row2: TRow) => number,
};

export default function SimpleTable<T>(props: Props<T>) {
  const classes = useStyles();
  const {
    columnStruct,
    tableData,
    onActionsClick,
    sortFunc: _sortFunc,
    ...extraProps
  } = props;

  const data = tableData || [];

  const rows = data.map((row: T, rowIdx: number) => {
    const rowKey = JSON.stringify(row || {});
    return (
      <TableRow key={rowKey}>
        {columnStruct.map((column, columnIdx) => (
          <RenderCell
            row={row}
            rowIdx={rowIdx}
            column={column}
            columnIdx={columnIdx}
            classes={classes}
            key={`${rowIdx}_${columnIdx}`}
          />
        ))}

        {onActionsClick && (
          <BodyTableCell>
            <Button
              variant="outlined"
              onClick={event => onActionsClick(row, event.target)}
              aria-label="Action Menu">
              <MoreHorizIcon color="action" />
            </Button>
          </BodyTableCell>
        )}
      </TableRow>
    );
  });

  return (
    <div {...extraProps} className={classes.body}>
      <Table>
        <TableHead>
          <TableRow>
            {columnStruct
              .concat(onActionsClick ? [{title: 'actions'}] : [])
              .map((column, idx) => (
                <HeadTableCell key={'row' + idx}>{column.title}</HeadTableCell>
              ))}
          </TableRow>
        </TableHead>
        <TableBody>{rows}</TableBody>
      </Table>
    </div>
  );
}

export function toLabels(obj: {}): Labels {
  if (!obj) {
    return {};
  }
  return Object.keys(obj).reduce((map, key) => {
    map[key] = obj[key];
    return map;
  }, {});
}
