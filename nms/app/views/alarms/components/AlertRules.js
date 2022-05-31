/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import * as React from 'react';
import AddEditRule from './rules/AddEditRule';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import SeverityIndicator from './severity/SeverityIndicator';
import SimpleTable from './table/SimpleTable';
import TableActionDialog from './table/TableActionDialog';
import TableAddButton from './table/TableAddButton';
import axios from 'axios';
import {Parse} from './prometheus/PromQLParser';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from './AlarmContext';
import {useLoadRules} from './hooks';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useSnackbars} from '../../../hooks/useSnackbar';

import {useParams} from 'react-router-dom';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import type {GenericRule} from './rules/RuleInterface';

const useStyles = makeStyles(theme => ({
  root: {
    paddingTop: theme.spacing(4),
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
  const snackbars = useSnackbars();
  const classes = useStyles();
  const params = useParams();
  const [lastRefreshTime, setLastRefreshTime] = React.useState(
    new Date().getTime().toString(),
  );
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

  const loadMatchingAlerts = React.useCallback(async () => {
    try {
      // only show matching alerts for prometheus rules for now
      if (selectedRow && selectedRow.ruleType === PROMETHEUS_RULE_TYPE) {
        const response = await apiUtil.viewMatchingAlerts({
          networkId: params.networkId,
          expression: selectedRow.expression,
        });
        setMatchingAlertsCount(response.length);
      }
    } catch (error) {
      snackbars.error('Could not load matching alerts for rule');
    }
  }, [selectedRow, apiUtil, params.networkId, snackbars]);
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
          networkId: params.networkId,
          ruleName: selectedRow.name,
          cancelToken: cancelSource.token,
        });
        snackbars.success(`Successfully deleted alert rule`);
      }
    } catch (error) {
      snackbars.error(
        `Unable to delete alert rule: ${
          error.response ? error.response?.data?.message : error.message
        }. Please try again.`,
      );
    } finally {
      setLastRefreshTime(new Date().toLocaleString());
    }
  }, [params.networkId, ruleMap, selectedRow, snackbars]);

  const handleViewAlertModalClose = React.useCallback(() => {
    setIsViewAlertModalOpen(false);
    setMatchingAlertsCount(null);
  }, [setIsViewAlertModalOpen, setMatchingAlertsCount]);
  const columns = React.useMemo(
    () => [
      {
        title: 'Name',
        field: 'name',
      },
      {
        title: 'Severity',
        field: 'severity',
        render: currRow => <SeverityIndicator severity={currRow.severity} />,
      },
      {
        title: 'Fire Alert When',
        field: 'fireAlertWhen',
        render: currRow => {
          try {
            const exp = Parse(currRow.expression);
            if (exp) {
              const metricName = exp.lh.selectorName?.toUpperCase() || '';
              const operator = exp.operator?.toString() || '';
              const value = exp.rh.value?.toString() || '';
              return `${metricName} ${operator} ${value} for ${currRow.period}`;
            }
          } catch {}
          return 'error';
        },
      },
      {
        title: 'Description',
        field: 'description',
      },
    ],
    [],
  );

  if (isAddEditAlert) {
    return (
      <AddEditRule
        initialConfig={selectedRow}
        isNew={isNewAlert}
        defaultRuleType={PROMETHEUS_RULE_TYPE}
        onExit={() => {
          setIsAddEditAlert(false);
          setLastRefreshTime(new Date().toLocaleString());
        }}
      />
    );
  }

  return (
    <Grid className={classes.root}>
      <SimpleTable
        onRowClick={row => setSelectedRow(row)}
        columnStruct={columns}
        tableData={rules || []}
        dataTestId="alert-rules"
        menuItems={[
          {
            name: 'View',
            handleFunc: () => handleView(),
          },
          {
            name: 'Edit',
            handleFunc: () => handleEdit(),
          },
          {
            name: 'Delete',
            handleFunc: () => {
              handleDelete();
            },
          },
        ]}
      />
      {isLoading && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}

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
