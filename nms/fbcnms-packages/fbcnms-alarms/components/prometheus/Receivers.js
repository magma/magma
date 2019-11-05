/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {AlertReceiver} from '../AlarmAPIType';

import AlertActionDialog from '../AlertActionDialog';
import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable from '../SimpleTable';

import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';
import type {ApiUrls} from '../ApiUrls';

const useStyles = makeStyles({
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
});

type Props = {
  apiUrls: ApiUrls,
};

export default function Receivers(props: Props) {
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentAlert, setCurrentAlert] = useState<Object>({});
  const [showAlertActionDialog, setShowAlertActionDialog] = useState<?'view'>(
    null,
  );
  const [lastRefreshTime, _setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDialogAction = args => {
    setShowAlertActionDialog(args);
    setMenuAnchorEl(null);
  };
  const receiverImportantLabels = (receiver: AlertReceiver) => {
    if (receiver.slack_configs) {
      const slackConfig = receiver.slack_configs[0];
      // api_url = blah
      const {api_url, channel, text, title} = slackConfig;
      return {api_url, channel, text, title};
    }
  };

  const {isLoading, error, response} = useAxios<null, Array<AlertReceiver>>({
    method: 'get',
    url: props.apiUrls.viewReceivers(match),
    cacheCounter: lastRefreshTime,
  });

  if (error) {
    enqueueSnackbar(
      `Unable to load receivers: ${
        error.response ? error.response.data.message : error.message
      }`,
      {variant: 'error'},
    );
  }

  const receiversList = response?.data || [];

  const receiversData = receiversList.map(receiver => {
    return {
      name: receiver.name,
      type: 'slack',
      labels: receiverImportantLabels(receiver),
    };
  });
  // show alarms table
  // many structures to support, slack, pagerduty. show a name + most important info as labels?
  return (
    <>
      <SimpleTable
        tableData={receiversData}
        onActionsClick={(alert, target) => {
          setMenuAnchorEl(target);
          setCurrentAlert(alert);
        }}
        columnStruct={[
          {title: 'name', path: ['name']},
          {title: 'type', path: ['type'], render: 'chip'},
          {title: 'labels', path: ['labels']},
        ]}
      />
      {isLoading && receiversData.length === 0 && (
        <div className={classes.loading}>
          <CircularProgress />
        </div>
      )}
      <Menu
        anchorEl={menuAnchorEl}
        keepMounted
        open={Boolean(menuAnchorEl)}
        onClose={() => setMenuAnchorEl(null)}>
        <MenuItem onClick={() => onDialogAction('view')}>View</MenuItem>
      </Menu>
      <AlertActionDialog
        open={showAlertActionDialog != null}
        onClose={() => onDialogAction(null)}
        title={'View Alert'}
        alertConfig={currentAlert || {}}
        showCopyButton={true}
        showDeleteButton={false}
      />
    </>
  );
}
