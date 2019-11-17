/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FiringAlarm, Labels} from '../AlarmAPIType';

import AlertActionDialog from '../AlertActionDialog';
import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable, {SEVERITY} from '../SimpleTable';
import {get} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import type {ApiUtil} from '../AlarmsApi';

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
  filterLabels?: (labels: Labels, rule: FiringAlarm) => Labels,
};

export default function FiringAlerts(props: Props) {
  const {apiUtil, filterLabels} = props;
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [_currentAlert, setCurrentAlert] = useState<Object>({});
  const [_showAlertActionDialog, setShowAlertActionDialog] = useState<?'view'>(
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

  const {isLoading, error, response} = apiUtil.useAlarmsApi(
    apiUtil.viewFiringAlerts,
    {networkId: match.params.networkId},
    lastRefreshTime,
  );

  if (error) {
    enqueueSnackbar(
      `Unable to load firing alerts: ${
        error.response ? error.response.data.message : error.message
      }`,
      {variant: 'error'},
    );
  }

  const alertData = response
    ? response.map(alert => {
        let labels = alert.labels;
        if (labels && filterLabels) {
          labels = filterLabels(labels, alert);
        }
        return {
          name: alert.labels?.alertname,
          labels: labels ?? {},
          annotations: alert.annotations ?? {},
          rawData: alert,
        };
      })
    : [];

  // can we take the columnStruct and use it to edit existing + add new?
  // we have order and path to let us reconstruct the message,
  //    need to hide structs (state)
  return (
    <>
      <SimpleTable
        columnStruct={[
          {
            title: 'name',
            renderFunc: (data, classes) => {
              const entity = data.labels.entity || data.labels.nodeMac || null;
              const desc = data.annotations.description;
              return (
                <>
                  <div className={classes.titleCell}>{data.name}</div>
                  {entity && (
                    <div className={classes.secondaryCell}>{entity}</div>
                  )}
                  <div className={classes.secondaryItalicCell}>{desc}</div>
                </>
              );
            },
          },
          {title: 'severity', path: ['labels', 'severity'], render: 'severity'},
          {
            title: 'labels',
            path: ['labels'],
            hideFields: ['alertname', 'severity', 'team'],
          },
          {
            title: 'annotations',
            path: ['annotations'],
            hideFields: ['description'],
          },
        ]}
        tableData={alertData}
        sortFunc={alert =>
          get(
            SEVERITY,
            [get(alert, ['labels', 'severity']).toLowerCase(), 'index'],
            undefined,
          )
        }
        onActionsClick={(alert, target) => {
          setMenuAnchorEl(target);
          setCurrentAlert(alert);
        }}
        data-testid="firing-alerts"
      />
      {isLoading && alertData.length === 0 && (
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
        open={_showAlertActionDialog != null}
        onClose={() => onDialogAction(null)}
        title={'View Alert'}
        alertConfig={_currentAlert?.rawData || {}}
      />
    </>
  );
}
