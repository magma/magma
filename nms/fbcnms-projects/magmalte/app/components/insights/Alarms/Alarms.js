/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AlertConfig} from './AlarmAPIType';

import AddEditAlert from './AddEditAlert';
import Button from '@material-ui/core/Button';
import Chip from '@material-ui/core/Chip';
import IconButton from '@material-ui/core/IconButton';
import LinearProgress from '@material-ui/core/LinearProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import MoreHorizIcon from '@material-ui/icons/MoreHoriz';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';
import RefreshIcon from '@material-ui/icons/Refresh';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableHead from '@material-ui/core/TableHead';
import TableRow from '@material-ui/core/TableRow';
import Tooltip from '@material-ui/core/Tooltip';
import Typography from '@material-ui/core/Typography';
import ViewDeleteAlertDialog from './ViewDeleteAlertDialog';

import grey from '@material-ui/core/colors/grey';
import orange from '@material-ui/core/colors/orange';
import red from '@material-ui/core/colors/red';
import yellow from '@material-ui/core/colors/yellow';

import axios from 'axios';
import {MagmaAlarmAPIUrls} from './AlarmAPI';
import {Route, Switch} from 'react-router-dom';
import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';
import {withStyles} from '@material-ui/core/styles';

const useStyles = makeStyles(theme => ({
  header: {
    padding: theme.spacing(3),
    display: 'flex',
    justifyContent: 'space-between',
    backgroundColor: 'white',
    borderBottom: `1px solid ${theme.palette.divider}`,
  },
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

export default function Alarms() {
  const {match, relativePath, history} = useRouter();
  const onExit = () => history.push(`${match.url}/`);
  return (
    <>
      <Switch>
        <Route
          path={relativePath('/new_alert')}
          render={() => <AddEditAlert onExit={onExit} />}
        />
        <Route path={match.path} render={() => <AlarmsTable />} />
      </Switch>
    </>
  );
}

function AlarmsTable() {
  const [menuAnchorEl, setMenuAnchorEl] = useState(null);
  const [currentAlert, setCurrentAlert] = useState<AlertConfig>({
    alert: '',
    expr: '',
  });
  const [showViewDeleteDialog, setShowViewDeleteDialog] = useState<?(
    | 'view'
    | 'delete'
  )>(null);
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDialogClose = () => {
    setShowViewDeleteDialog(null);
    setMenuAnchorEl(null);
  };

  const onDelete = () => {
    axios
      .delete(MagmaAlarmAPIUrls.alertConfig(match), {
        params: {alert_name: currentAlert.alert},
      })
      .then(() =>
        enqueueSnackbar(`Successfully deleted alert`, {
          variant: 'success',
        }),
      )
      .catch(error =>
        enqueueSnackbar(`Unable to delete alert: ${error}. Please try again.`, {
          variant: 'error',
        }),
      )
      .finally(() => {
        onDialogClose();
        setLastRefreshTime(new Date().toLocaleString());
      });
  };

  const {isLoading, error, response} = useAxios<null, Array<AlertConfig>>({
    method: 'get',
    url: MagmaAlarmAPIUrls.alertConfig(match),
    cacheCounter: lastRefreshTime,
  });

  if (error) {
    enqueueSnackbar('Unable to load alerts', {variant: 'error'});
  }

  const alerts = response?.data || [];

  const rows = alerts.map(alert => {
    const {description, ...customAnnotations} = alert.annotations ?? {};
    const {severity, team} = alert.labels ?? {};
    return (
      <TableRow key={alert.alert}>
        <BodyTableCell>
          <div className={classes.alertName}>{alert.alert}</div>
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
        <BodyTableCell>
          <Button
            variant="outlined"
            onClick={event => {
              setMenuAnchorEl(event.target);
              setCurrentAlert(alert);
            }}>
            <MoreHorizIcon color="action" />
          </Button>
        </BodyTableCell>
      </TableRow>
    );
  });

  return (
    <div>
      <div className={classes.header}>
        <Typography variant="h5">Alerts</Typography>
        <div>
          <Tooltip title={'Last refreshed: ' + lastRefreshTime}>
            <IconButton
              color="inherit"
              onClick={() => setLastRefreshTime(new Date().toLocaleString())}
              disabled={isLoading}>
              <RefreshIcon />
            </IconButton>
          </Tooltip>
          <NestedRouteLink to={'/new_alert'}>
            <Button variant="contained" color="primary">
              New Alert
            </Button>
          </NestedRouteLink>
        </div>
      </div>
      {isLoading ? <LinearProgress /> : null}
      <div className={classes.body}>
        <Table>
          <TableHead>
            <TableRow>
              <HeadTableCell>ALL ALERTS</HeadTableCell>
              <HeadTableCell>SEVERITY</HeadTableCell>
              <HeadTableCell>TEAM</HeadTableCell>
              <HeadTableCell>ANNOTATIONS</HeadTableCell>
              <HeadTableCell>ACTIONS</HeadTableCell>
            </TableRow>
          </TableHead>
          <TableBody>{rows}</TableBody>
        </Table>
      </div>
      <Menu
        anchorEl={menuAnchorEl}
        keepMounted
        open={Boolean(menuAnchorEl)}
        onClose={() => setMenuAnchorEl(null)}>
        <MenuItem onClick={() => setShowViewDeleteDialog('view')}>
          View
        </MenuItem>
        <MenuItem onClick={() => setShowViewDeleteDialog('delete')}>
          Delete
        </MenuItem>
      </Menu>
      <ViewDeleteAlertDialog
        open={showViewDeleteDialog != null}
        onClose={onDialogClose}
        onDelete={onDelete}
        alertConfig={currentAlert}
        deletingAlert={showViewDeleteDialog === 'delete'}
      />
    </div>
  );
}
