/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AlarmsHeader from './AlarmsHeader';
import AlarmsTable from './AlarmsTable';
import Button from '@material-ui/core/Button';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import ViewDeleteAlertDialog from './ViewDeleteAlertDialog';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../../common/useMagmaAPI';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

type Props = {
  onFiringAlerts: () => void,
  onNewAlert: () => void,
};

export default function EditAllAlerts(props: Props) {
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentAlert, setCurrentAlert] = useState<string>('');
  const [showViewDeleteDialog, setShowViewDeleteDialog] = useState<?(
    | 'view'
    | 'delete'
  )>(null);
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDialogClose = () => {
    setShowViewDeleteDialog(null);
    setMenuAnchorEl(null);
  };

  const onDelete = () => {
    MagmaV1API.deleteNetworksByNetworkIdPrometheusAlertConfig({
      networkId: nullthrows(match.params.networkId),
      alertName: currentAlert,
    })
      .then(() =>
        enqueueSnackbar(`Successfully deleted alert`, {
          variant: 'success',
        }),
      )
      .catch(error =>
        enqueueSnackbar(
          `Unable to delete alert: ${
            error.response ? error.response.data.message : error.message
          }. Please try again.`,
          {
            variant: 'error',
          },
        ),
      )
      .finally(() => {
        onDialogClose();
        setLastRefreshTime(new Date().toLocaleString());
      });
  };

  const {isLoading, error, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusAlertConfig,
    {networkId: nullthrows(match.params.networkId)},
    undefined, // onResponse
    lastRefreshTime,
  );

  if (error) {
    enqueueSnackbar(`Unable to load alerts: ${error.response.data.message}`, {
      variant: 'error',
    });
  }

  const alerts = response || [];
  const alertData = alerts.map(alert => {
    return {
      name: alert.alert,
      annotations: alert.annotations ?? {},
      labels: alert.labels ?? {},
    };
  });

  return (
    <>
      <AlarmsHeader
        title="Edit Alerts"
        isLoading={isLoading}
        lastRefreshTime={lastRefreshTime}
        onRefreshClick={refreshTime => setLastRefreshTime(refreshTime)}>
        <Button
          variant="contained"
          color="secondary"
          onClick={props.onFiringAlerts}>
          Firing Alerts
        </Button>
        <Button variant="contained" color="primary" onClick={props.onNewAlert}>
          New Alert
        </Button>
      </AlarmsHeader>
      <AlarmsTable
        alertsColumnName="All Alerts"
        alertData={alertData}
        onActionsClick={(alertName, target) => {
          setMenuAnchorEl(target);
          setCurrentAlert(alertName);
        }}
      />
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
        alertConfig={
          alerts.find(alert => alert.alert === currentAlert) ?? {
            alert: '',
            expr: '',
          }
        }
        deletingAlert={showViewDeleteDialog === 'delete'}
      />
    </>
  );
}
