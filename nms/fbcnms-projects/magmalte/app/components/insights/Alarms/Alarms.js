/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {FiringMagmaAlarm} from './AlarmAPIType';

import AddEditAlert from './AddEditAlert';
import AlarmsHeader from './AlarmsHeader';
import AlarmsTable from './AlarmsTable';
import Button from '@material-ui/core/Button';
import EditAllAlerts from './EditAllAlerts';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';

import {MagmaAlarmAPIUrls} from './AlarmAPI';
import {Route, Switch} from 'react-router-dom';
import {useAxios, useRouter} from '@fbcnms/ui/hooks';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useState} from 'react';

export default function Alarms() {
  const {match, relativePath, history} = useRouter();
  return (
    <>
      <Switch>
        <Route
          path={relativePath('/new_alert')}
          render={() => (
            <AddEditAlert
              onExit={() => history.push(`${match.url}/edit_alerts`)}
            />
          )}
        />
        <Route
          path={relativePath('/edit_alerts')}
          render={() => (
            <EditAllAlerts
              onFiringAlerts={() => history.push(`${match.url}/`)}
              onNewAlert={() => history.push(`${match.url}/new_alert`)}
            />
          )}
        />
        <Route path={match.path} render={() => <FiringAlerts />} />
      </Switch>
    </>
  );
}

function FiringAlerts() {
  const [lastRefreshTime, setLastRefreshTime] = useState<string>(
    new Date().toLocaleString(),
  );
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();

  const {isLoading, error, response} = useAxios<null, Array<FiringMagmaAlarm>>({
    method: 'get',
    url: MagmaAlarmAPIUrls.viewFiringAlerts(match),
    cacheCounter: lastRefreshTime,
  });

  if (error) {
    enqueueSnackbar('Unable to load firing alerts', {variant: 'error'});
  }

  const alerts = response?.data || [];

  const alertData = alerts.map(alert => {
    return {
      name: alert.labels?.alertname,
      labels: alert.labels ?? {},
      annotations: alert.annotations ?? {},
    };
  });

  return (
    <>
      <AlarmsHeader
        title="Firing Alerts"
        isLoading={isLoading}
        lastRefreshTime={lastRefreshTime}
        onRefreshClick={refreshTime => setLastRefreshTime(refreshTime)}>
        <NestedRouteLink to={'/edit_alerts'}>
          <Button variant="contained" color="secondary">
            Edit Alerts
          </Button>
        </NestedRouteLink>
        <NestedRouteLink to={'/new_alert'}>
          <Button variant="contained" color="primary">
            New Alert
          </Button>
        </NestedRouteLink>
      </AlarmsHeader>
      <AlarmsTable alertsColumnName="Firing Alerts" alertData={alertData} />
    </>
  );
}
