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
import AddEditReceiver from './AddEditReceiver';
import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import SimpleTable from '../../SimpleTable';
import TableActionDialog from '../../TableActionDialog';
import TableAddButton from '../../common/TableAddButton';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import type {AlertReceiver} from '../../AlarmAPIType';
import type {ApiUtil} from '../../AlarmsApi';

const useStyles = makeStyles({
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

type Props = {
  apiUtil: ApiUtil,
};

export default function Receivers(props: Props) {
  const {apiUtil} = props;
  const [isAddEditReceiver, setIsAddEditReceiver] = React.useState(false);
  const [isNewReceiver, setIsNewReceiver] = React.useState(false);
  const [isMenuOpen, setIsMenuOpen] = React.useState(false);
  const [isDialogOpen, setIsDialogOpen] = React.useState(false);
  const [selectedRow, setSelectedRow] = React.useState<?AlertReceiver>(null);
  const menuAnchorEl = React.useRef<?HTMLElement>(null);
  const [lastRefreshTime, setLastRefreshTime] = React.useState<string>(
    new Date().toLocaleString(),
  );
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const handleActionsMenuOpen = React.useCallback(
    (row: AlertReceiver, eventTarget: HTMLElement) => {
      setSelectedRow(row);
      menuAnchorEl.current = eventTarget;
      setIsMenuOpen(true);
    },
    [menuAnchorEl, setIsMenuOpen, setSelectedRow],
  );

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
            networkId: match.params.networkId,
            receiverName: selectedRow.name,
          });
          enqueueSnackbar(`Successfully deleted receiver`, {
            variant: 'success',
          });
          setIsMenuOpen(false);
        }
      } catch (error) {
        enqueueSnackbar(
          `Unable to delete receiver: ${
            error.response ? error.response?.data?.message : error.message
          }. Please try again.`,
          {
            variant: 'error',
          },
        );
      } finally {
        setLastRefreshTime(new Date().toLocaleString());
      }
    }
    makeRequest();
  }, [apiUtil, enqueueSnackbar, match.params.networkId, selectedRow]);

  const handleViewDialogOpen = React.useCallback(() => {
    setIsDialogOpen(true);
    setIsMenuOpen(false);
  }, [setIsDialogOpen]);

  const handleViewDialogClose = React.useCallback(() => {
    setIsDialogOpen(false);
  }, [setIsDialogOpen]);

  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getReceivers,
    {networkId: match.params.networkId},
    lastRefreshTime,
  );

  if (error) {
    enqueueSnackbar(
      `Unable to load receivers: ${
        error.response ? error.response.data.message : error.message
      }`,
      {variant: 'error'},
    );
  }

  const receiversData = response || [];

  if (isAddEditReceiver) {
    return (
      <AddEditReceiver
        receiver={selectedRow || newReceiver()}
        apiUtil={apiUtil}
        isNew={isNewReceiver}
        onExit={() => {
          setIsAddEditReceiver(false);
          setLastRefreshTime(new Date().toLocaleString());
        }}
      />
    );
  }

  return (
    <>
      <SimpleTable
        tableData={receiversData}
        onActionsClick={handleActionsMenuOpen}
        columnStruct={[
          {
            title: 'name',
            getValue: row => row.name,
          },
          {
            title: 'notifications',
            render: 'labels',
            getValue: getNotificationsSummary,
          },
        ]}
      />
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
          setIsMenuOpen(false);
        }}
        data-testid="add-receiver-button"
      />
    </>
  );
}

function newReceiver(): AlertReceiver {
  return {
    name: '',
  };
}

function getNotificationsSummary(receiver: AlertReceiver) {
  const summary = {};
  const {slack_configs, email_configs, webhook_configs} = receiver;
  if (slack_configs) {
    const channelNames = slack_configs.map(({channel}) => {
      if (!channel) {
        return '';
      }
      return channel.replace(/#/, '');
    });
    summary['Slack Channels'] = channelNames.join(', ');
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

  return summary;
}
