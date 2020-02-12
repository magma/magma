/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import * as React from 'react';
import AddEditRule from './rules/AddEditRule';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import SimpleTable from './SimpleTable';
import TableActionDialog from './TableActionDialog';
import TableAddButton from './common/TableAddButton';
import axios from 'axios';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from './AlarmContext';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useLoadRules} from './hooks';
import {useRouter} from '@fbcnms/ui/hooks';
import type {ColumnData} from './SimpleTable';
import type {GenericRule} from './rules/RuleInterface';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
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

const PROMETHEUS_RULE_TYPE = 'prometheus';

export default function AlertRules<TRuleUnion>() {
  const {apiUtil, ruleMap} = useAlarmContext();
  const classes = useStyles();
  const enqueueSnackbar = useEnqueueSnackbar();
  const {match} = useRouter();
  const [lastRefreshTime, setLastRefreshTime] = React.useState(
    new Date().getTime().toString(),
  );
  const menuAnchorEl = React.useRef<?HTMLElement>(null);
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);
  const [
    selectedRow,
    setSelectedRow,
  ] = React.useState<?GenericRule<TRuleUnion>>(null);
  const [isNewAlert, setIsNewAlert] = React.useState(false);
  const [isAddEditAlert, setIsAddEditAlert] = React.useState(false);
  const [isViewAlertModalOpen, setIsViewAlertModalOpen] = React.useState(false);
  const [matchingAlertsCount, setMatchingAlertsCount] = React.useState<?number>(
    null,
  );

  const {rules, isLoading} = useLoadRules({
    ruleMap,
    lastRefreshTime,
  });

  const columnStruct = React.useMemo<
    Array<ColumnData<GenericRule<TRuleUnion>>>,
  >(
    () => [
      {
        title: 'name',
        getValue: x => x.name,
        renderFunc: (rule, classes) => {
          return (
            <>
              <div className={classes.titleCell}>{rule.name}</div>
              <div className={classes.secondaryItalicCell}>
                {rule.description}
              </div>
            </>
          );
        },
      },
      {
        title: 'severity',
        getValue: rule => rule.severity,
        render: 'severity',
      },
      {
        title: 'period',
        getValue: rule => rule.period,
      },
      {
        title: 'expression',
        getValue: rule => rule.expression,
      },
    ],
    [],
  );

  const handleActionsMenuOpen = React.useCallback(
    (row: GenericRule<TRuleUnion>, eventTarget: HTMLElement) => {
      setSelectedRow(row);
      menuAnchorEl.current = eventTarget;
      setIsMenuOpen(true);
    },
    [menuAnchorEl, setIsMenuOpen, setSelectedRow],
  );
  const loadMatchingAlerts = React.useCallback(async () => {
    try {
      // only show matching alerts for prometheus rules for now
      if (selectedRow && selectedRow.ruleType === PROMETHEUS_RULE_TYPE) {
        const response = await apiUtil.viewMatchingAlerts({
          networkId: match.params.networkId,
          expression: selectedRow.expression,
        });
        setMatchingAlertsCount(response.length);
      }
    } catch (error) {
      enqueueSnackbar('Could not load matching alerts for rule', {
        variant: 'error',
      });
    }
  }, [selectedRow, apiUtil, match.params.networkId, enqueueSnackbar]);
  const handleActionsMenuClose = React.useCallback(() => {
    setSelectedRow(null);
    menuAnchorEl.current = null;
    setIsMenuOpen(false);
  }, [menuAnchorEl, setIsMenuOpen, setSelectedRow]);
  const handleEdit = React.useCallback(() => {
    setIsAddEditAlert(true);
    setIsNewAlert(false);
  }, []);
  const handleView = React.useCallback(() => {
    loadMatchingAlerts();
    setIsViewAlertModalOpen(true);
  }, [loadMatchingAlerts]);
  const handleDelete = React.useCallback(async () => {
    try {
      if (selectedRow) {
        const cancelSource = axios.CancelToken.source();
        const {deleteRule} = ruleMap[selectedRow.ruleType];
        await deleteRule({
          networkId: match.params.networkId,
          ruleName: selectedRow.name,
          cancelToken: cancelSource.token,
        });
        enqueueSnackbar(`Successfully deleted alert rule`, {
          variant: 'success',
        });
      }
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
      setLastRefreshTime(new Date().toLocaleString());
      setIsMenuOpen(false);
    }
  }, [enqueueSnackbar, match.params.networkId, ruleMap, selectedRow]);

  const handleViewAlertModalClose = React.useCallback(() => {
    setIsViewAlertModalOpen(false);
    setMatchingAlertsCount(null);
  }, [setIsViewAlertModalOpen, setMatchingAlertsCount]);

  if (isAddEditAlert) {
    return (
      <AddEditRule
        initialConfig={selectedRow}
        isNew={isNewAlert}
        defaultRuleType={PROMETHEUS_RULE_TYPE}
        onExit={() => {
          setIsAddEditAlert(false);
          setLastRefreshTime(new Date().toLocaleString());
          handleActionsMenuClose();
        }}
      />
    );
  }

  return (
    <Grid className={classes.root}>
      <SimpleTable
        columnStruct={columnStruct}
        tableData={rules || []}
        onActionsClick={handleActionsMenuOpen}
      />
      {isLoading && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl.current}
        open={isMenuOpen}
        onClose={handleActionsMenuClose}>
        <MenuItem onClick={handleEdit}>Edit</MenuItem>
        <MenuItem onClick={handleView}>View</MenuItem>
        <MenuItem onClick={handleDelete}>Delete</MenuItem>
      </Menu>
      {selectedRow && (
        <TableActionDialog
          open={isViewAlertModalOpen}
          onClose={handleViewAlertModalClose}
          title={'View Alert Rule'}
          additionalContent={
            matchingAlertsCount !== null && (
              <span>
                This rule matches <strong>{matchingAlertsCount}</strong> active
                alarm(s).
              </span>
            )
          }
          row={selectedRow.rawRule}
          showCopyButton={true}
          showDeleteButton={true}
          onDelete={() => {
            handleViewAlertModalClose();
            return handleDelete();
          }}
          RowViewer={
            ruleMap && selectedRow
              ? ruleMap[selectedRow.ruleType].RuleViewer
              : undefined
          }
        />
      )}
      <TableAddButton
        onClick={() => {
          setIsNewAlert(true);
          setSelectedRow(null);
          setIsAddEditAlert(true);
        }}
        label="Add Alert"
        data-testid="add-edit-alert-button"
      />
    </Grid>
  );
}
