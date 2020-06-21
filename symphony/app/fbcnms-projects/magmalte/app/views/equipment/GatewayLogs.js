/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import type {ActionQuery} from '../../components/ActionTable';

import ActionTable from '../../components/ActionTable';
import Grid from '@material-ui/core/Grid';
import LogChart from './GatewayLogChart';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';

import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';

// elastic search pagination through 'from' mechanism has a 10000 row limit
// we have to use a different mechanism in case we want to go higher, we should
// use search_after
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-request-search-after.html
const MAX_PAGE_ROW_COUNT = 10000;

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
  button: {
    margin: theme.spacing(1),
  },
}));

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
          <LogChart />
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

  return (
    <ActionTable
      title={'Logs'}
      data={(query: ActionQuery) =>
        new Promise((resolve, reject) => {
          const countReq = MagmaV1API.getNetworksByNetworkIdLogsCount({
            networkId: networkId,
            filters: `gateway_id:${gatewayId}`,
            simpleQuery: query.search,
          });

          const searchReq = MagmaV1API.getNetworksByNetworkIdLogsSearch({
            networkId: networkId,
            filters: `gateway_id:${gatewayId}`,
            from: (query.page * query.pageSize).toString(),
            size: query.pageSize.toString(),
            simpleQuery: query.search,
          });

          Promise.all([countReq, searchReq])
            .then(([countResp, searchResp]) => {
              let gatewayLogCount = countResp;
              if (gatewayLogCount > MAX_PAGE_ROW_COUNT) {
                gatewayLogCount = MAX_PAGE_ROW_COUNT;
              }
              const logRows = searchResp.filter(Boolean).map(elastic_hit => {
                const src = elastic_hit._source;
                const date = new Date(src['@timestamp'] ?? 0);
                const msg = src['message'];
                return {
                  date: date,
                  service: src['ident'],
                  logType: getLogType(msg ?? ''),
                  output: msg,
                };
              });
              resolve({
                data: logRows,
                page: query.page,
                totalCount: gatewayLogCount,
              });
            })
            .catch(err => reject(err));
        })
      }
      columns={[
        {title: 'Date', field: 'date', type: 'datetime'},
        {title: 'Service', field: 'service'},
        {title: 'Type', field: 'logType'},
        {title: 'Output', field: 'output'},
      ]}
      options={{
        actionsColumnIndex: -1,
        pageSizeOptions: [5, 10],
        exportButton: true,
        exportAllData: true,
      }}
    />
  );
}
