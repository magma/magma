/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import AlertDetailsPane from './AlertDetails/AlertDetailsPane';
import CircularProgress from '@material-ui/core/CircularProgress';
import Grid from '@material-ui/core/Grid';
import React from 'react';
import SimpleTable, {toLabels} from '../table/SimpleTable';
import Slide from '@material-ui/core/Slide';
import {SEVERITY} from '../severity/Severity';
import {get} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useAlarmContext} from '../AlarmContext';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

import type {FiringAlarm} from '../AlarmAPIType';

const useStyles = makeStyles(theme => ({
  root: {
    padding: theme.spacing(4),
  },
  loading: {
    display: 'flex',
    height: '100%',
    alignItems: 'center',
    justifyContent: 'center',
  },
}));

export default function FiringAlerts() {
  const {apiUtil, filterLabels} = useAlarmContext();
  const [selectedRow, setSelectedRow] = useState<?FiringAlarm>(null);
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

  const showRowDetailsPane = React.useCallback(
    (row: FiringAlarm) => {
      setSelectedRow(row);
    },
    [setSelectedRow],
  );
  const hideDetailsPane = React.useCallback(() => {
    setSelectedRow(null);
  }, [setSelectedRow]);

  if (error) {
    enqueueSnackbar(
      `Unable to load firing alerts. ${
        error.response ? error.response.data.message : error.message || ''
      }`,
      {variant: 'error'},
    );
  }

  const alertData: Array<FiringAlarm> = response
    ? response.map(alert => {
        let labels = alert.labels;
        if (labels && filterLabels) {
          labels = filterLabels(labels);
        }
        return {
          ...alert,
          labels,
        };
      })
    : [];

  return (
    <Grid className={classes.root} container spacing={2}>
      <Grid item xs={selectedRow ? 8 : 12}>
        <SimpleTable
          onRowClick={showRowDetailsPane}
          columnStruct={[
            {
              title: 'name',
              getValue: x => x.labels?.alertname,
              renderFunc: (data, classes) => {
                const entity =
                  data.labels.entity || data.labels.nodeMac || null;
                const desc = data?.annotations?.description ?? '';
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
          data-testid="firing-alerts"
        />
        {isLoading && alertData.length === 0 && (
          <div className={classes.loading}>
            <CircularProgress />
          </div>
        )}
      </Grid>
      <Slide direction="left" in={!!selectedRow}>
        <Grid item xs={4}>
          {selectedRow && (
            <AlertDetailsPane alert={selectedRow} onClose={hideDetailsPane} />
          )}
        </Grid>
      </Slide>
    </Grid>
  );
}
