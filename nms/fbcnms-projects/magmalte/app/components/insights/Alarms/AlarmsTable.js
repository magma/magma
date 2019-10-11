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
  teamChip: {
    color: theme.palette.secondary.main,
  },
  alertName: {
    fontSize: 18,
    fontWeight: 500,
  },
  alertDescription: {
    fontStyle: 'italic',
    color: theme.palette.text.secondary,
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

export const SEVERITY_STYLE = {
  critical: 'redSeverityChip',
  major: 'orangeSeverityChip',
  minor: 'yellowSeverityChip',
  warning: 'greySeverityChip',
  info: 'greySeverityChip',
  notice: 'greySeverityChip',
};

type Props = {
  alertsColumnName: string,
  alertData: Array<{
    name: string,
    annotations: {[string]: string},
    labels: {[string]: string},
  }>,
  onActionsClick?: (alertName: string, target: HTMLElement) => void,
};

export default function AlarmsTable(props: Props) {
  const classes = useStyles();
  const {alertsColumnName, alertData, onActionsClick} = props;

  const rows = alertData.map(alert => {
    const {description, ...customAnnotations} = alert.annotations ?? {};
    const {severity, team} = alert.labels ?? {};
    return (
      <TableRow key={alert.name}>
        <BodyTableCell>
          <div className={classes.alertName}>{alert.name}</div>
          <div className={classes.alertDescription}>{description}</div>
        </BodyTableCell>
        <BodyTableCell>
          {severity in SEVERITY_STYLE && (
            <Chip
              classes={{
                outlined: classes[SEVERITY_STYLE[severity]],
              }}
              label={severity.toUpperCase()}
              variant="outlined"
            />
          )}
        </BodyTableCell>
        <BodyTableCell>
          {team && (
            <Chip
              classes={{outlinedPrimary: classes.teamChip}}
              label={team.toUpperCase()}
              color="primary"
              variant="outlined"
            />
          )}
        </BodyTableCell>
        <BodyTableCell>
          {Object.keys(customAnnotations).map(key => (
            <Chip
              key={key}
              className={classes.labelChip}
              label={`${key}:${customAnnotations[key]}`}
              size="small"
            />
          ))}
        </BodyTableCell>
        {onActionsClick && (
          <BodyTableCell>
            <Button
              variant="outlined"
              onClick={event => onActionsClick(alert.name, event.target)}>
              <MoreHorizIcon color="action" />
            </Button>
          </BodyTableCell>
        )}
      </TableRow>
    );
  });

  const columnNames = [alertsColumnName, 'severity', 'team', 'annotations'];
  if (onActionsClick) {
    columnNames.push('actions');
  }

  return (
    <div className={classes.body}>
      <Table>
        <TableHead>
          <TableRow>
            {columnNames.map((column, i) => (
              <HeadTableCell key={i}>{column}</HeadTableCell>
            ))}
          </TableRow>
        </TableHead>
        <TableBody>{rows}</TableBody>
      </Table>
    </div>
  );
}
