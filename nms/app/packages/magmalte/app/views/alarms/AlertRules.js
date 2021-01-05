/*
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
import type {ColumnData} from '@fbcnms/alarms/components/table/SimpleTable';
import type {GenericRule} from '@fbcnms/alarms/components/rules/RuleInterface';

import * as React from 'react';
import AddAlertIcon from '@material-ui/icons/AddAlert';
import AddEditRule from '@fbcnms/alarms/components/rules/AddEditRule';
import AppContext from '@fbcnms/ui/context/AppContext';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import SimpleTable from '@fbcnms/alarms/components/table/SimpleTable';
import TableActionDialog from '@fbcnms/alarms/components/table/TableActionDialog';
import TableAddButton from '@fbcnms/alarms/components/table/TableAddButton';
import axios from 'axios';
import nullthrows from '@fbcnms/util/nullthrows';

import {colors, typography} from '../../theme/default';
import {makeStyles} from '@material-ui/styles';
import {triggerAlertSync} from '../../state/SyncAlerts';
import {useAlarmContext} from '@fbcnms/alarms/components/AlarmContext';
import {useContext} from 'react';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useLoadRules} from '@fbcnms/alarms/components/hooks';
import {useRouter} from '@fbcnms/ui/hooks';

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
  appBarBtn: {
    color: colors.primary.white,
    background: colors.primary.comet,
    fontFamily: typography.button.fontFamily,
    fontWeight: typography.button.fontWeight,
    fontSize: typography.button.fontSize,
    lineHeight: typography.button.lineHeight,
    letterSpacing: typography.button.letterSpacing,

    '&:hover': {
      background: colors.primary.mirage,
    },
  },
}));

const PROMETHEUS_RULE_TYPE = 'prometheus';

export default function AlertRules<TRuleUnion>() {
  const {apiUtil, ruleMap} = useAlarmContext();
  const isSuperUser = useContext(AppContext).user.isSuperUser;
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
  const networkId = nullthrows(match.params.networkId);
  const {rules, isLoading} = useLoadRules({
    ruleMap,
    lastRefreshTime,
  });

  const SyncAlertsButton = () => {
    const classes = useStyles();
    const enqueueSnackbar = useEnqueueSnackbar();
    const predefinedAlertsTitle = 'Sync Predefined Alerts';

    return (
      <Button
        variant="contained"
        className={classes.appBarBtn}
        onClick={async () => {
          await triggerAlertSync(networkId, enqueueSnackbar);
          setLastRefreshTime(new Date().toLocaleString());
        }}>
        {predefinedAlertsTitle}
      </Button>
    );
  };

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
          networkId: networkId,
          expression: selectedRow.expression,
        });
        setMatchingAlertsCount(response.length);
      }
    } catch (error) {
      enqueueSnackbar('Could not load matching alerts for rule', {
        variant: 'error',
      });
    }
  }, [selectedRow, apiUtil, networkId, enqueueSnackbar]);
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
          networkId: networkId,
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
  }, [enqueueSnackbar, networkId, ruleMap, selectedRow]);

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
      <CardTitleRow
        icon={AddAlertIcon}
        label={'Alert Rules'}
        filter={isSuperUser ? SyncAlertsButton : undefined}
      />

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
