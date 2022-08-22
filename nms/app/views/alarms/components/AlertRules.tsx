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
 */
import * as React from 'react';
import AddEditRule from './rules/AddEditRule';
import Button from '@mui/material/Button';
import CardTitleRow from '../../../components/layout/CardTitleRow';
import CircularProgress from '@mui/material/CircularProgress';
import Grid from '@mui/material/Grid';
import SeverityIndicator from './severity/SeverityIndicator';
import SimpleTable, {SimpleTableProps} from './table/SimpleTable';
import TableActionDialog from './table/TableActionDialog';
import axios from 'axios';
import nullthrows from '../../../../shared/util/nullthrows';

import {NotificationsActive} from '@mui/icons-material';
import {PROMETHEUS_RULE_TYPE, useAlarmContext} from './AlarmContext';
import {Parse} from './prometheus/PromQLParser';
import {Theme} from '@mui/material/styles';
import {colors} from '../../../theme/default';
import {getErrorMessage} from '../../../util/ErrorUtils';
import {makeStyles} from '@mui/styles';
import {triggerAlertSync} from '../../../util/SyncAlerts';
import {useEnqueueSnackbar, useSnackbars} from '../../../hooks/useSnackbar';
import {useLoadRules} from './hooks';
import {useParams} from 'react-router-dom';
import type {GenericRule} from './rules/RuleInterface';

const useStyles = makeStyles<Theme>(theme => ({
  root: {
    paddingTop: theme.spacing(4),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
  emptyRulesTitle: {
    color: colors.primary.comet,
  },
  emptyRulesDescription: {
    color: colors.primary.comet,
    fontSize: '12px',
    marginBottom: '8px',
  },
}));

function AlertRulesFilter(props: {
  onAddRuleClick: () => void;
  refresh: () => void;
}) {
  const params = useParams();
  const networkId = nullthrows(params.networkId);
  const enqueueSnackbar = useEnqueueSnackbar();
  return (
    <Grid container justifyContent="center" alignItems="center" spacing={2}>
      <Grid item>
        <Button
          variant="outlined"
          color="primary"
          onClick={() => {
            void triggerAlertSync(networkId, enqueueSnackbar);
            props.refresh();
          }}>
          Sync Predefined Alerts
        </Button>
      </Grid>
      <Grid item>
        <Button
          data-testid="add-edit-alert-button"
          variant="contained"
          color="primary"
          onClick={() => props.onAddRuleClick()}>
          Create Custom Rule
        </Button>
      </Grid>
    </Grid>
  );
}

function AlertRulesEmpty(props: {
  onAddRuleClick: () => void;
  refresh: () => void;
}) {
  const classes = useStyles();

  return (
    <Grid container direction="column" alignItems="center">
      <div className={classes.emptyRulesTitle}>No Rules Added</div>
      <div className={classes.emptyRulesDescription}>
        Find out about possible issues in the network by syncing predefined
        alerts or creating custom rules.
      </div>
      <AlertRulesFilter
        onAddRuleClick={() => props.onAddRuleClick()}
        refresh={() => props.refresh()}
      />
    </Grid>
  );
}

export default function AlertRules<TRuleUnion>() {
  const {apiUtil, ruleMap} = useAlarmContext();
  const snackbars = useSnackbars();
  const classes = useStyles();
  const params = useParams();
  const [lastRefreshTime, setLastRefreshTime] = React.useState(
    new Date().getTime().toString(),
  );
  const [selectedRow, setSelectedRow] = React.useState<
    GenericRule<TRuleUnion> | null | undefined
  >(null);
  const [isNewAlert, setIsNewAlert] = React.useState(false);
  const [isAddEditAlert, setIsAddEditAlert] = React.useState(false);
  const [isViewAlertModalOpen, setIsViewAlertModalOpen] = React.useState(false);
  const [matchingAlertsCount, setMatchingAlertsCount] = React.useState<
    number | null | undefined
  >(null);
  const {rules, isLoading} = useLoadRules({
    ruleMap,
    lastRefreshTime,
  });

  const loadMatchingAlerts = React.useCallback(async () => {
    try {
      // only show matching alerts for prometheus rules for now
      if (selectedRow && selectedRow.ruleType === PROMETHEUS_RULE_TYPE) {
        const response = (
          await apiUtil.viewMatchingAlerts({
            networkId: params.networkId!,
            expression: selectedRow.expression,
          })
        ).data;
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
    void loadMatchingAlerts();
    setIsViewAlertModalOpen(true);
  }, [loadMatchingAlerts]);
  const handleDelete = React.useCallback(async () => {
    try {
      if (selectedRow) {
        const cancelSource = axios.CancelToken.source();
        const {deleteRule} = ruleMap[selectedRow.ruleType];
        await deleteRule({
          networkId: params.networkId!,
          ruleName: selectedRow.name,
          cancelToken: cancelSource.token,
        });
        snackbars.success(`Successfully deleted alert rule`);
      }
    } catch (error) {
      snackbars.error(
        `Unable to delete alert rule: ${getErrorMessage(
          error,
        )}. Please try again.`,
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
    () =>
      [
        {
          title: 'Name',
          field: 'name',
        },
        {
          title: 'Severity',
          field: 'severity',
          render: (currRow: GenericRule<any>) => (
            <SeverityIndicator severity={currRow.severity} />
          ),
        },
        {
          title: 'Fire Alert When',
          field: 'fireAlertWhen',
          render: (currRow: GenericRule<any>) => {
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
      ] as SimpleTableProps<GenericRule<any>>['columnStruct'],
    [],
  );

  const openAddRule = () => {
    setIsNewAlert(true);
    setSelectedRow(null);
    setIsAddEditAlert(true);
  };
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
      <CardTitleRow
        label="Alert rules"
        icon={NotificationsActive}
        filter={() => (
          <AlertRulesFilter
            onAddRuleClick={() => openAddRule()}
            refresh={() => setLastRefreshTime(new Date().toLocaleString())}
          />
        )}
      />
      <SimpleTable
        localization={{
          body: {
            emptyDataSourceMessage: (
              <AlertRulesEmpty
                onAddRuleClick={() => openAddRule()}
                refresh={() => setLastRefreshTime(new Date().toLocaleString())}
              />
            ),
          },
        }}
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
              void handleDelete();
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
    </Grid>
  );
}
