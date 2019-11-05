/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AlertActionDialog from '../AlertActionDialog';
import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable from '../SimpleTable';
import type {AlertSuppression} from '../AlarmAPIType';

import {makeStyles} from '@material-ui/styles';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';
import type {ApiUrls} from '../ApiUrls';

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
}));

type Props = {
  apiUrls: ApiUrls,
};

export default function Suppressions(props: Props) {
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentAlert, setCurrentAlert] = useState<Object>({});
  const [showAlertActionDialog, setShowAlertActionDialog] = useState<?'view'>(
    null,
  );
  const [lastRefreshTime, _setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const [_isAddEditAlert, _setIsAddEditAlert] = useState<boolean>(false);
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const onDialogAction = args => {
    setShowAlertActionDialog(args);
    setMenuAnchorEl(null);
  };

  const {isLoading, error, response} = useAxios<null, Array<AlertSuppression>>({
    method: 'get',
    url: props.apiUrls.viewSilences(match),
    cacheCounter: lastRefreshTime,
  });

  if (error) {
    enqueueSnackbar(
      `Unable to load suppressions: ${
        error.response ? error.response.data.message : error.message
      }`,
      {variant: 'error'},
    );
  }

  const silencesList = response?.data || [];

  const columnStruct = [
    {title: 'name', path: ['comment']},
    {title: 'active', path: ['status', 'state']},
    {title: 'created by', path: ['createdBy']},
    {title: 'matchers', path: ['matchers'], render: 'multipleGroups'},
  ];

  return (
    <>
      <SimpleTable
        tableData={silencesList}
        onActionsClick={(alert, target) => {
          setMenuAnchorEl(target);
          setCurrentAlert(alert);
        }}
        columnStruct={columnStruct}
      />
      {isLoading && silencesList.length === 0 && (
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
