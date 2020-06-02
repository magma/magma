/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Card from '@material-ui/core/Card';
import React from 'react';
import TabbedTable from './TabbedTable';
import Text from '@fbcnms/ui/components/design-system/Text';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import Grid from '@material-ui/core/Grid';
import {Alarm} from '@material-ui/icons';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import type {RowData} from './TabbedTable';
import type {prom_firing_alert} from '@fbcnms/magma-api';

const useStyles = makeStyles(theme => ({
  cardTitle: {
    marginBottom: theme.spacing(1),
  },
  cardTitleIcon: {
    marginRight: theme.spacing(1),
  },
}));

type AlertTable = {[string]: Array<RowData>};

type Severity = 'Critical' | 'Major' | 'Minor' | 'Other';
const severityMap: {[string]: Severity} = {
  critical: 'Critical',
  page: 'Critical',
  warn: 'Major',
  major: 'Major',
  minor: 'Minor',
};

export default function() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);

  const {isLoading, response} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdAlerts,
    {
      networkId,
    },
  );

  if (isLoading) {
    return <LoadingFiller />;
  }

  const data: AlertTable = {Critical: [], Major: [], Minor: [], Other: []};
  if (!response) {
    return null;
  }

  const alerts: Array<prom_firing_alert> = response;
  alerts.forEach(alert => {
    const labelInfo = {
      job: alert.labels['job'] || '',
      instance: alert.labels['instance'] || '',
    };

    const timingInfo = {
      startsAt: alert.startsAt || '',
      endsAt: alert.endsAt || '',
      updatedAt: alert.updatedAt || '',
    };

    const sev: Severity = severityMap[alert.labels['severity']] || 'Other';

    data[sev].push({
      name: alert.labels['alertname'],
      cols: [
        JSON.stringify(labelInfo),
        JSON.stringify(alert.annotations),
        JSON.stringify(alert.status),
        JSON.stringify(timingInfo),
      ],
    });
  });

  return (
    <>
      {/* TODO: Can come back and make this a reusable component for other cards */}
      <Grid container alignItems="center" className={classes.cardTitle}>
        <Alarm className={classes.cardTitleIcon} />
        <Text>Alerts ({alerts.length})</Text>
      </Grid>
      {/* <Card> */}
      <TabbedTable data={data} />
      {/* </Card> */}
    </>
  );
}
