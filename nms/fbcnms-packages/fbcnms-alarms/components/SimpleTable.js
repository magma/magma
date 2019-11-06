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

import {clone, get} from 'lodash';
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

export const SEVERITY = {
  critical: {index: 1, style: 'redSeverityChip'},
  major: {index: 2, style: 'orangeSeverityChip'},
  minor: {index: 3, style: 'yellowSeverityChip'},
  warning: {index: 4, style: 'yellowSeverityChip'},
  info: {index: 5, style: 'greySeverityChip'},
  notice: {index: 6, style: 'greySeverityChip'},
};

const renderCell = (
  tableRow: Object,
  tableIdx: number,
  column: ColumnData,
  columnIdx: number,
  classes: any,
) => {
  let cellValue = tableRow;
  let renderFunc;

  if (column.renderFunc) {
    cellValue = column.renderFunc(tableRow, classes);
    renderFunc = renderCellCustomFunc;
  } else {
    cellValue = get(tableRow, column.path);
    if (column.render === 'severity') {
      renderFunc = renderCellSeverity;
    } else if (column.render === 'multipleGroups') {
      renderFunc = renderMultiLabels;
    } else if (column.render === 'chip') {
      renderFunc = renderCellChip;
    } else if (typeof cellValue === 'object') {
      renderFunc = renderCellObj;
      if (column.hideFields != undefined && column.hideFields?.length) {
        cellValue = clone(cellValue);
        // $FlowFixMe flow thinks hideFields won't be defined
        column.hideFields.forEach(key => delete cellValue[key]);
      }
    } else {
      renderFunc = renderCellString;
    }
  }
  return renderFunc(cellValue, classes, columnIdx, `${tableIdx}_${columnIdx}`);
};

const renderCellCustomFunc = (
  cellValue: Object,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={cellKey}>
    <div>{cellValue}</div>
  </BodyTableCell>
);

const renderMultiLabels = (
  cellValueList: any,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={`cell_${cellKey}`}>
    {cellValueList.map(cellValue => (
      <div
        key={`cell_div_${cellKey}`}
        className={columnIdx === 0 ? classes.titleCell : classes.secondaryCell}>
        {Object.keys(cellValue).map(keyName => (
          <Chip
            key={`cell_chip_${cellKey}_${keyName}`}
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

const renderLabelValue = labelValue => {
  if (typeof labelValue === 'boolean') {
    return labelValue ? 'true' : 'false';
  }
  return labelValue;
};

const renderCellObj = (
  cellValue: Object,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={cellKey}>
    <div
      className={
        columnIdx === 0 ? classes.titleCell : classes.secondaryItalicCell
      }>
      {Object.keys(cellValue).map(keyName => (
        <Chip
          key={`${cellKey}_${keyName}`}
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
  </BodyTableCell>
);

const renderCellString = (
  cellValue: Object,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={cellKey}>
    <div
      className={
        columnIdx === 0 ? classes.titleCell : classes.secondaryItalicCell
      }>
      {cellValue}
    </div>
  </BodyTableCell>
);

const renderCellSeverity = (
  cellValue: Object,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={cellKey}>
    {cellValue && cellValue.toLowerCase() in SEVERITY && (
      <Chip
        key={`${cellKey}_severity`}
        classes={{
          outlined: classes[SEVERITY[cellValue.toLowerCase()].style],
          label: classes.ellipsisChip,
        }}
        label={cellValue.toUpperCase()}
        variant="outlined"
      />
    )}
  </BodyTableCell>
);

const renderCellChip = (
  cellValue: Object,
  classes,
  columnIdx,
  cellKey,
): React.Element<any> => (
  <BodyTableCell key={cellKey}>
    {cellValue && (
      <Chip
        key={`${cellKey}_chip`}
        classes={{outlinedPrimary: classes.secondaryChip}}
        label={cellValue.toUpperCase()}
        color="primary"
        variant="outlined"
      />
    )}
  </BodyTableCell>
);

export type ColumnData = {
  title: string,
  path?: Array<string>,
  hideFields?: Array<string>,
  render?: string,
  // valid drop-down options list
  validOptions?: Array<string>,
  renderFunc?: (tableRow: any, classes: any) => React.Element<any>,
  tooltip?: React.Node,
};

type Props = {
  columnStruct: Array<ColumnData>,
  tableData: Array<Object>,
  onActionsClick?: (alert: Object, target: HTMLElement) => void,
  sortFunc?: (alert: Object, alert2: Object) => number,
};

export default function SimpleTable(props: Props) {
  const classes = useStyles();
  const {
    columnStruct,
    tableData,
    onActionsClick,
    sortFunc: _sortFunc,
    ...extraProps
  } = props;

  const data = tableData;

  const rows = data.map((tableRow: Object, tableIdx: number) => {
    const rowKey = JSON.stringify(tableRow || {});
    return (
      <TableRow key={rowKey}>
        {columnStruct.map((column, columnIdx) =>
          renderCell(tableRow, tableIdx, column, columnIdx, classes),
        )}

        {onActionsClick && (
          <BodyTableCell>
            <Button
              variant="outlined"
              onClick={event => onActionsClick(tableRow, event.target)}>
              <MoreHorizIcon color="action" />
            </Button>
          </BodyTableCell>
        )}
      </TableRow>
    );
  });
  if (onActionsClick) {
    columnStruct.push({title: 'actions'});
  }

  return (
    <div className={classes.body} {...extraProps}>
      <Table>
        <TableHead>
          <TableRow>
            {columnStruct.map((column, idx) => (
              <HeadTableCell key={'row' + idx}>{column.title}</HeadTableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>{rows}</TableBody>
      </Table>
    </div>
  );
}
