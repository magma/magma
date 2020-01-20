/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import CircularProgress from '@material-ui/core/CircularProgress';
import Menu from '@material-ui/core/Menu';
import MenuItem from '@material-ui/core/MenuItem';
import React from 'react';
import SimpleTable, {toLabels} from '../SimpleTable';
import TableActionDialog from '../TableActionDialog';
import {SEVERITY} from '../Severity';
import {get} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';
import type {ApiUtil} from '../AlarmsApi';
import type {FiringAlarm, Labels} from '../AlarmAPIType';

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
  filterLabels?: (labels: Labels, alert: FiringAlarm) => Labels,
};

export default function FiringAlerts(props: Props) {
  const {apiUtil, filterLabels} = props;
  const [menuAnchorEl, setMenuAnchorEl] = useState<?HTMLElement>(null);
  const [currentAlert, setCurrentAlert] = useState<?FiringAlarm>(null);
  const [showDialog, setShowDialog] = useState(false);
  const [lastRefreshTime, _setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

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

  const alertData: Array<FiringAlarm> = response
    ? response.map(alert => {
        let labels = alert.labels;
        if (labels && filterLabels) {
          labels = filterLabels(labels, alert);
        }
        return {
          ...alert,
          labels,
        };
      })
    : [];

  return (
    <>
      <SimpleTable
        columnStruct={[
          {
            title: 'name',
            getValue: x => x.labels?.alertname,
            renderFunc: (data, classes) => {
              const entity = data.labels.entity || data.labels.nodeMac || null;
              const desc = data.annotations.description;
              return (
                <>
                  <div className={classes.titleCell}>
                    {data.labels?.alertname}
                  </div>
                  {entity && (
                    <div className={classes.secondaryCell}>{entity}</div>
                  )}
                  <div className={classes.secondaryItalicCell}>{desc}</div>
                </>
              );
            },
          },
          {
            title: 'severity',
            getValue: x => x.labels?.severity,
            render: 'severity',
          },
          {
            title: 'labels',
            getValue: x => toLabels(x.labels),
            render: 'labels',
            hideFields: ['alertname', 'severity', 'team'],
          },
          {
            title: 'annotations',
            getValue: x => toLabels(x.annotations),
            render: 'labels',
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
        <MenuItem onClick={() => setShowDialog(true)}>View</MenuItem>
      </Menu>
      <TableActionDialog
        open={showDialog}
        onClose={() => setShowDialog(false)}
        title={'View Alert'}
        row={currentAlert}
      />
    </>
  );
}
