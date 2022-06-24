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
import AddEditReceiver from './AddEditReceiver';
import Button from '@material-ui/core/Button';
import CircularProgress from '@material-ui/core/CircularProgress';
import GlobalConfig from './GlobalConfig';
import Grid from '@material-ui/core/Grid';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import SettingsIcon from '@material-ui/icons/Settings';
import SimpleTable, {LabelsCell} from '../../table/SimpleTable';
import TableActionDialog from '../../table/TableActionDialog';
import TableAddButton from '../../table/TableAddButton';
import {Theme} from '@material-ui/core/styles';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../../AlarmContext';
import {useNetworkId} from '../../hooks';
import {useSnackbars} from '../../../../../hooks/useSnackbar';

import {getErrorMessage} from '../../../../../util/ErrorUtils';
import type {AlertReceiver} from '../../AlarmAPIType';

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
}));

export default function Receivers() {
  const {apiUtil, alertManagerGlobalConfigEnabled} = useAlarmContext();
  const [isAddEditReceiver, setIsAddEditReceiver] = React.useState(false);
  const [isEditGlobalSettings, setIsEditGlobalSettings] = React.useState(false);
  const [isNewReceiver, setIsNewReceiver] = React.useState(false);
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [selectedRow, setSelectedRow] = React.useState<AlertReceiver | null>(
    null,
  );
  const menuAnchorEl = React.useRef<HTMLElement | null>(null);
  const [lastRefreshTime, setLastRefreshTime] = React.useState<string>(
    new Date().toLocaleString(),
  );
  const networkId = useNetworkId();
  const classes = useStyles();
  const snackbars = useSnackbars();

  const handleActionsMenuClose = React.useCallback(() => {
    setSelectedRow(null);
    menuAnchorEl.current = null;
    setIsMenuOpen(false);
  }, [menuAnchorEl, setIsMenuOpen, setSelectedRow]);

  const handleEdit = React.useCallback(() => {
    setIsAddEditReceiver(true);
    setIsNewReceiver(false);
    setIsMenuOpen(false);
  }, [setIsAddEditReceiver, setIsNewReceiver]);

  const handleDelete = React.useCallback(() => {
    async function makeRequest() {
      try {
        if (selectedRow) {
          await apiUtil.deleteReceiver({
            networkId,
            receiverName: selectedRow.name,
          });
          snackbars.success(`Successfully deleted receiver`);
          setIsMenuOpen(false);
        }
      } catch (error) {
        snackbars.error(
          `Unable to delete receiver: ${getErrorMessage(
            error,
          )}. Please try again.`,
        );
      } finally {
        setLastRefreshTime(new Date().toLocaleString());
      }
    }
    void makeRequest();
  }, [apiUtil, networkId, selectedRow, snackbars]);

  const handleViewDialogOpen = React.useCallback(() => {
    setIsDialogOpen(true);
    setIsMenuOpen(false);
  }, [setIsDialogOpen]);

  const handleViewDialogClose = React.useCallback(() => {
    setIsDialogOpen(false);
  }, [setIsDialogOpen]);

  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getReceivers,
    {networkId},
    lastRefreshTime,
  );

  if (error) {
    snackbars.error(`Unable to load receivers: ${getErrorMessage(error)}`);
  }

  const receiversData = response || [];

  if (isAddEditReceiver) {
    return (
      <AddEditReceiver
        receiver={selectedRow || newReceiver()}
        isNew={isNewReceiver}
        onExit={() => {
          setIsAddEditReceiver(false);
          setLastRefreshTime(new Date().toLocaleString());
        }}
      />
    );
  }

  if (isEditGlobalSettings) {
    return (
      <GlobalConfig
        onExit={() => {
          setIsEditGlobalSettings(false);
          setLastRefreshTime(new Date().toLocaleString());
        }}
      />
    );
  }

  return (
    <Grid className={classes.root} container spacing={2} direction="column">
      {alertManagerGlobalConfigEnabled === true && (
        <Grid item>
          <Button
            onClick={() => setIsEditGlobalSettings(true)}
            startIcon={<SettingsIcon />}
            variant="outlined">
            Settings
          </Button>
        </Grid>
      )}
      <Grid item>
        <>
          <SimpleTable
            onRowClick={row => setSelectedRow(row)}
            columnStruct={[
              {
                title: 'Name',
                field: 'name',
              },
              {
                title: 'Notifications',
                field: 'labels',
                render: row => {
                  const labels = getNotificationsSummary(row);
                  return <LabelsCell value={labels} />;
                },
              },
            ]}
            tableData={receiversData}
            dataTestId="receiver"
            menuItems={[
              {
                name: 'View',
                handleFunc: () => handleViewDialogOpen(),
              },
              {
                name: 'Edit',
                handleFunc: () => handleEdit(),
              },
              {
                name: 'Delete',
                handleFunc: () => handleDelete(),
              },
            ]}
          />
        </>
      </Grid>
      {isLoading && receiversData.length === 0 && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl.current}
        keepMounted
        open={isMenuOpen}
        onClose={handleActionsMenuClose}>
        <MenuItem onClick={handleViewDialogOpen}>View</MenuItem>
        <MenuItem onClick={handleEdit}>Edit</MenuItem>
        <MenuItem onClick={handleDelete}>Delete</MenuItem>
      </Menu>
      <TableActionDialog
        open={isDialogOpen}
        onClose={handleViewDialogClose}
        title={'View Receiver'}
        row={selectedRow || {}}
        showCopyButton={true}
        showDeleteButton={false}
      />
      <TableAddButton
        label="Add Receiver"
        onClick={() => {
          setIsNewReceiver(true);
          setIsAddEditReceiver(true);
          setSelectedRow(null);
        }}
        data-testid="add-receiver-button"
      />
    </Grid>
  );
}

function newReceiver(): AlertReceiver {
  return {
    name: '',
  };
}

function getNotificationsSummary(receiver: AlertReceiver) {
  const summary: Record<string, string> = {};
  const {
    slack_configs,
    email_configs,
    webhook_configs,
    pagerduty_configs,
    pushover_configs,
  } = receiver;
  if (slack_configs) {
    const channelNames = slack_configs.reduce<Array<string>>(
      (list, {channel}) => {
        if (channel != null && channel.trim() !== '') {
          list.push(channel.replace(/#/, ''));
        }
        return list;
      },
      [],
    );
    if (channelNames.length > 0) {
      summary['Slack Channels'] = channelNames.join(', ');
    } else {
      const configCount = slack_configs.length;
      summary['Slack'] = `${configCount} channel${configCount > 1 ? 's' : ''}`;
    }
  }
  if (email_configs) {
    const emailAddresses = email_configs.map(({to}) => to);
    summary['Emails'] = emailAddresses.join(', ');
  }
  if (webhook_configs) {
    const webhookUrls = webhook_configs.map(({url}) => {
      try {
        const parsed = new URL(url);
        const trimmedLength = 24;
        const hostAndPath = `${parsed.host}${parsed.pathname}`;
        return hostAndPath.length > trimmedLength
          ? `${hostAndPath.substring(0, trimmedLength)}...`
          : hostAndPath;
      } catch (e) {
        console.error(e);
        return url.slice(0, 8);
      }
    });
    summary['Webhook'] = webhookUrls.join(', ');
  }

  if (pagerduty_configs) {
    summary['PagerDuty'] = pagerduty_configs
      .map(conf => conf.severity)
      .join(',')
      .slice(0, 12);
  }

  if (pushover_configs) {
    summary['Pushover'] = pushover_configs
      .map(conf => conf.title || '')
      .join(', ')
      .slice(0, 12);
  }

  return summary;
}
