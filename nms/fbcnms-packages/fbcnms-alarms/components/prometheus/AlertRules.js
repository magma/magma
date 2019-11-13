/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddEditAlert from '../AddEditAlert';
import AddIcon from '@material-ui/icons/Add';
import AlertActionDialog from '../AlertActionDialog';
import CircularProgress from '@material-ui/core/CircularProgress';
import Fab from '@material-ui/core/Fab';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable, {SEVERITY} from '../SimpleTable';
import {get} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

import HelpIcon from '@material-ui/icons/Help';
import IconButton from '@material-ui/core/IconButton';
import Tooltip from '@material-ui/core/Tooltip';

import type {AlertConfig} from '../AlarmAPIType';
import type {ApiUtil} from '../AlarmsApi';

type Props = {
  apiUtil: ApiUtil,
};

const useStyles = makeStyles(theme => ({
  addButton: {
    position: 'fixed',
    bottom: 0,
    right: 0,
    margin: theme.spacing(2),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
  helpButton: {
    color: 'black',
  },
}));

export default function AlertRules(props: Props) {
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentAlert, setCurrentAlert] = useState<?AlertConfig>(null);
  const [showAlertActionDialog, setShowAlertActionDialog] = useState<?(
    | 'view'
    | 'delete'
    | 'edit'
  )>(null);
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const [isNewAlert, setIsNewAlert] = React.useState(false);
  const [isAddEditAlert, setIsAddEditAlert] = useState<boolean>(false);
  const [matchingAlarmsCount, setMatchingAlarmsCount] = useState<?number>(null);
  const {match} = useRouter();
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDialogEdit = () => {
    setIsNewAlert(false);
    setIsAddEditAlert(true);
    setMenuAnchorEl(null);
  };

  const onDialogAction = args => {
    // Fetch number of matching alarms to display in modal
    if (args && currentAlert && currentAlert.expr) {
      props.apiUtil
        .viewMatchingAlerts({
          networkId: match.params.networkId,
          expression: currentAlert.expr,
        })
        .then(response => {
          setMatchingAlarmsCount(response.length);
        });
    } else {
      setMatchingAlarmsCount(null);
    }

    setShowAlertActionDialog(args);
    setMenuAnchorEl(null);
  };

  const onDelete = async () => {
    try {
      await props.apiUtil.deleteAlertRule({
        networkId: match.params.networkId,
        ruleName: currentAlert?.alert || '',
      });
      enqueueSnackbar(`Successfully deleted alert rule`, {
        variant: 'success',
      });
    } catch (error) {
      enqueueSnackbar(
        `Unable to delete alert rule: ${
          error.response ? error.response?.data?.message : error.message
        }. Please try again.`,
        {
          variant: 'error',
        },
      );
    } finally {
      onDialogAction(null);
      setLastRefreshTime(new Date().toLocaleString());
    }
  };

  const {isLoading, error, response} = props.apiUtil.useAlarmsApi(
    props.apiUtil.getAlertRules,
    {networkId: match.params.networkId},
    lastRefreshTime,
  );

  if (error) {
    enqueueSnackbar(
      `Unable to load alert rules: ${error.response?.data?.message}`,
      {variant: 'error'},
    );
  }

  const alertData: Array<AlertConfig> = response || [];
  const columnStruct = [
    {
      title: 'name',
      path: ['alert'],
      renderFunc: (data: AlertConfig, classes) => {
        const desc = data.annotations?.description;
        return (
          <>
            <div className={classes.titleCell}>{data.alert}</div>
            <div className={classes.secondaryItalicCell}>{desc}</div>
          </>
        );
      },
    },
    {
      title: 'severity',
      path: ['labels', 'severity'],
      validOptions: Object.keys(SEVERITY),
      render: 'severity',
    },
    {
      title: 'period',
      path: ['for'],
      tooltip: (
        <Tooltip
          title={
            'Enter the amount of time the alert expression needs to be ' +
            'true for before the alert fires.'
          }
          placement="right">
          <HelpIcon />
        </Tooltip>
      ),
    },
    {
      title: 'expression',
      path: ['expr'],
      tooltip: (
        <Tooltip
          title={
            'To learn more about how to write alert expressions, click ' +
            'on the help icon to open the prometheus querying basics guide.'
          }
          placement="right">
          <IconButton
            className={classes.helpButton}
            href="https://prometheus.io/docs/prometheus/latest/querying/basics/"
            target="_blank"
            size="small">
            <HelpIcon />
          </IconButton>
        </Tooltip>
      ),
    },
  ];

  if (isAddEditAlert) {
    return (
      <AddEditAlert
        apiUtil={props.apiUtil}
        columnStruct={columnStruct}
        initialConfig={currentAlert}
        isNew={isNewAlert}
        onExit={() => {
          setIsAddEditAlert(false);
          setLastRefreshTime(new Date().toLocaleString());
        }}
        rule={currentAlert}
      />
    );
  }

  return (
    <>
      <SimpleTable
        columnStruct={columnStruct}
        tableData={alertData}
        sortFunc={alert =>
          get(
            SEVERITY,
            [get(alert, ['labels', 'severity']).toLowerCase(), 'index'],
            undefined,
          )
        }
        onActionsClick={(alert, target) => {
          setMenuAnchorEl(target);
          setCurrentAlert(alert);
        }}
      />
      {isLoading && alertData.length === 0 && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl}
        keepMounted
        open={Boolean(menuAnchorEl)}
        onClose={() => setMenuAnchorEl(null)}>
        <MenuItem onClick={() => onDialogEdit()}>Edit</MenuItem>
        <MenuItem onClick={() => onDialogAction('view')}>View</MenuItem>
        <MenuItem onClick={() => onDelete()}>Delete</MenuItem>
      </Menu>
      <AlertActionDialog
        open={showAlertActionDialog != null}
        onClose={() => onDialogAction(null)}
        title={
          showAlertActionDialog === 'delete'
            ? 'Delete Alert Rule'
            : 'View Alert Rule'
        }
        additionalContent={
          matchingAlarmsCount !== null && (
            <span>
              This rule matches <strong>{matchingAlarmsCount}</strong> active
              alarm(s).
            </span>
          )
        }
        alertConfig={currentAlert || {}}
        showCopyButton={showAlertActionDialog !== 'delete'}
        showDeleteButton={showAlertActionDialog === 'delete'}
        onDelete={onDelete}
      />
      <Fab
        className={classes.addButton}
        color="primary"
        onClick={() => {
          setIsNewAlert(true);
          setCurrentAlert(null);
          setIsAddEditAlert(true);
        }}
        aria-label="Add Alert"
        data-testid="add-edit-alert-button">
        <AddIcon />
      </Fab>
    </>
  );
}
