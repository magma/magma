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
import ListAltIcon from '@material-ui/icons/ListAlt';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {Bar} from 'react-chartjs-2';
import {DateTimePicker} from '@material-ui/pickers';
import {makeStyles} from '@material-ui/styles';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

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
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());

  return (
    <div className={classes.dashboardRoot}>
      <Grid container align="top" alignItems="flex-start">
        <Grid item xs={6}>
          <Text>
            <ListAltIcon />
            Logs
          </Text>
        </Grid>
        <Grid item xs={6}>
          <Grid container justify="flex-end" alignItems="center" spacing={1}>
            <Grid item>
              <Text>Filter By Date</Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                maxDate={endDate}
                disableFuture
                value={startDate}
                onChange={setStartDate}
              />
            </Grid>
            <Grid item>
              <Text>To</Text>
            </Grid>
            <Grid item>
              <DateTimePicker
                autoOk
                variant="inline"
                inputVariant="outlined"
                disableFuture
                value={endDate}
                onChange={setEndDate}
              />
            </Grid>
          </Grid>
        </Grid>
      </Grid>
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
function LogChart() {
  const data = {
    labels: ['January', 'February', 'March', 'April', 'May', 'June', 'July'],
    datasets: [
      {
        label: 'Log Counts',
        backgroundColor: 'rgba(255,99,132,0.2)',
        borderColor: 'rgba(255,99,132,1)',
        borderWidth: 1,
        hoverBackgroundColor: 'rgba(255,99,132,0.4)',
        hoverBorderColor: 'rgba(255,99,132,1)',
        data: [65, 59, 80, 81, 56, 55, 40],
      },
    ],
  };
  return (
    <Bar
      data={data}
      options={{
        maintainAspectRatio: false,
        scaleShowValues: true,
        scales: {
          xAxes: [
            {
              gridLines: {
                display: false,
              },
              ticks: {
                maxTicksLimit: 10,
              },
            },
          ],
          yAxes: [
            {
              gridLines: {
                drawBorder: true,
              },
              ticks: {
                maxTicksLimit: 1,
              },
            },
          ],
        },
      }}
    />
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
