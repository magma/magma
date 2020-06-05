/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import ActionTable from '../../components/ActionTable';
import Grid from '@material-ui/core/Grid';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import SettingsInputAntennaIcon from '@material-ui/icons/SettingsInputAntenna';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: theme.palette.magmalte.background,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: theme.palette.magmalte.appbar,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: 'white',
  },
  tab: {
    fontSize: '18px',
    textTransform: 'none',
  },
  tabLabel: {
    padding: '20px 0 20px 0',
  },
  tabIconLabel: {
    verticalAlign: 'middle',
    margin: '0 5px 3px 0',
  },
  // TODO: remove this when we actually fill out the grid sections
  contentPlaceholder: {
    padding: '50px 0',
  },
  paper: {
    height: 100,
    padding: theme.spacing(10),
    textAlign: 'center',
    color: theme.palette.text.secondary,
  },
  formControl: {
    margin: theme.spacing(1),
    minWidth: 120,
  },
}));

type LogRowType = {
  date: Date,
  service: string,
  logType: string,
  output: string,
};

const getLogType = (msg: string): string => {
  for (const typ of [
    'emerg',
    'alert',
    'crit',
    'err',
    'warning',
    'notice',
    'info',
    'debug',
  ]) {
    if (msg.toLowerCase().startsWith(typ)) {
      return typ;
    }
  }
  return 'debug';
};

export default function GatewayLogs() {
  const classes = useStyles();

  return (
    <div className={classes.dashboardRoot}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          Log Bar Chart
        </Grid>
        <Grid item xs={12}>
          <LogTable />
        </Grid>
      </Grid>
    </div>
  );
}

function LogTable() {
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);

  const {response: gatewayLogs} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdLogs,
    {
      networkId: networkId,
      filters: `gateway_id:${gatewayId}`,
    },
  );
  if (!gatewayLogs) {
    return null;
  }

  const logRows: Array<LogRowType> = gatewayLogs.map(elastic_hit => {
    const src = elastic_hit._source;
    const date = new Date(src['@timestamp'] ?? 0);
    const msg = src['message'];
    return {
      date: date,
      service: src['ident'],
      logType: getLogType(msg),
      output: msg,
    };
  });
  return (
    <ActionTable
      titleIcon={SettingsInputAntennaIcon}
      title={'EnodeBs'}
      data={logRows}
      columns={[
        {title: 'Date', field: 'date', type: 'datetime'},
        {title: 'Service', field: 'service'},
        {title: 'Type', field: 'logType'},
        {title: 'Output', field: 'output'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
      }}
    />
  );
}
