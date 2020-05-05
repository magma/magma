/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable, {toLabels} from '../table/SimpleTable';
import TableActionDialog from '../table/TableActionDialog';
import useRouter from '@fbcnms/alarms/hooks/useRouter';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useEnqueueSnackbar} from '@fbcnms/alarms/hooks/useSnackbar';
import {useState} from 'react';

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

export default function Suppressions() {
  const {apiUtil} = useAlarmContext();
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentRow, setCurrentRow] = useState<{}>({});
  const [showDialog, setShowDialog] = useState(false);
  const [lastRefreshTime, _setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const [_isAddEditAlert, _setIsAddEditAlert] = useState<boolean>(false);
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.getSuppressions,
    {networkId: match.params.networkId},
    lastRefreshTime,
  );

  if (error) {
    enqueueSnackbar(
      `Unable to load suppressions: ${
        error.response ? error.response.data.message : error.message
      }`,
      {variant: 'error'},
    );
  }

  const silencesList = response || [];

  return (
    <>
      <SimpleTable
        tableData={silencesList}
        onActionsClick={(alert, target) => {
          setMenuAnchorEl(target);
          setCurrentRow(alert);
        }}
        columnStruct={[
          {title: 'name', getValue: row => row.comment || ''},
          {title: 'active', getValue: row => row.status?.state ?? ''},
          {title: 'created by', getValue: row => row.createdBy},
          {
            title: 'matchers',
            getValue: row =>
              row.matchers
                ? row.matchers.map(matcher => toLabels(matcher))
                : [],
            render: 'multipleGroups',
          },
        ]}
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
        <MenuItem onClick={() => setShowDialog(true)}>View</MenuItem>
      </Menu>
      <TableActionDialog
        open={showDialog}
        onClose={() => setShowDialog(false)}
        title={'View Suppression'}
        row={currentRow || {}}
        showCopyButton={true}
        showDeleteButton={false}
      />
    </>
  );
}
