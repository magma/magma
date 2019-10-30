/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import AddEditAlert from './AddEditAlert';
import AlarmsHeader from './AlarmsHeader';
import AlarmsTable from './AlarmsTable';
import Button from '@material-ui/core/Button';
import EditAllAlerts from './EditAllAlerts';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import NestedRouteLink from '@fbcnms/ui/components/NestedRouteLink';
import React from 'react';

import useMagmaAPI from '../../../common/useMagmaAPI';
import {Route, Switch} from 'react-router-dom';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
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

  const {isLoading, error, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdAlerts,
    {networkId: match.params.networkId},
    undefined, // onResponse
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

  const alerts = response || [];
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
        onRefreshClick={refreshTime => setLastRefreshTime(refreshTime)}
        data-testid="firing-alerts">
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
