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
import Button from '@material-ui/core/Button';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import ListAltIcon from '@material-ui/icons/ListAlt';
import LogChart from './GatewayLogChart';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import React from 'react';
import Text from '@fbcnms/ui/components/design-system/Text';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {CsvBuilder} from 'filefy';
import {DateTimePicker} from '@material-ui/pickers';
import {colors} from '../../theme/default';
import {getStep} from '../../components/CustomHistogram';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useMemo, useRef, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

// elastic search pagination through 'from' mechanism has a 10000 row limit
// we have to use a different mechanism in case we want to go higher, we should
// use search_after
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-request-search-after.html
const MAX_PAGE_ROW_COUNT = 10000;
const EXPORT_DELIMITER = ',';
const LOG_COLUMNS = [
  {title: 'Date', field: 'date', type: 'datetime'},
  {title: 'Service', field: 'service'},
  {title: 'Type', field: 'logType'},
  {title: 'Output', field: 'output'},
];

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  topBar: {
    backgroundColor: colors.primary.mirage,
    padding: '20px 40px 20px 40px',
  },
  tabBar: {
    backgroundColor: colors.primary.brightGray,
    padding: '0 0 0 20px',
  },
  tabs: {
    color: colors.primary.white,
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

async function searchLogs(
  networkId,
  gatewayId,
  from,
  size,
  start,
  end,
  q,
  enqueueSnackbar,
) {
  const logs = await MagmaV1API.getNetworksByNetworkIdLogsSearch({
    networkId: networkId,
    filters: `gateway_id:${gatewayId}`,
    from: from.toString(),
    size: size.toString(),
    simpleQuery: q.search ?? '',
    start: start.toISOString(),
    end: end.toISOString(),
  })
    .then(searchResp => {
      const logs = searchResp.filter(Boolean).map(elastic_hit => {
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
      return logs;
    })
    .catch(err => {
      enqueueSnackbar('Error exporting logs ' + err, {
        variant: 'error',
      });
      return [];
    });
  return logs;
}

function exportLogs(
  networkId,
  gatewayId,
  from,
  size,
  start,
  end,
  q,
  enqueueSnackbar,
) {
  searchLogs(
    networkId,
    gatewayId,
    0,
    MAX_PAGE_ROW_COUNT,
    start,
    end,
    q,
    enqueueSnackbar,
  ).then(logRows => {
    const data = logRows.map(rowData =>
      LOG_COLUMNS.map(columnDef => rowData[columnDef.field]),
    );

    const currTs = Date.now();
    new CsvBuilder(`logs_${currTs}.csv`)
      .setDelimeter(EXPORT_DELIMITER)
      .setColumns(LOG_COLUMNS.map(columnDef => columnDef.title))
      .addRows(data)
      .exportFile();
  });
}

function handleLogQuery(
  networkId,
  gatewayId,
  from,
  size,
  start,
  end,
  q,
  enqueueSnackbar,
) {
  return new Promise((resolve, reject) => {
    const countReq = MagmaV1API.getNetworksByNetworkIdLogsCount({
      networkId: networkId,
      start: start.toISOString(),
      end: end.toISOString(),
      filters: `gateway_id:${gatewayId}`,
      simpleQuery: q.search,
    });

    const searchReq = searchLogs(
      networkId,
      gatewayId,
      (q.page * q.pageSize).toString(),
      q.pageSize.toString(),
      start,
      end,
      q,
      enqueueSnackbar,
    );

    Promise.all([countReq, searchReq])
      .then(([countResp, searchResp]) => {
        let gatewayLogCount = countResp;
        if (gatewayLogCount > MAX_PAGE_ROW_COUNT) {
          gatewayLogCount = MAX_PAGE_ROW_COUNT;
        }
        const page =
          gatewayLogCount < q.page * q.pageSize
            ? gatewayLogCount / q.pageSize
            : q.page;
        resolve({
          data: searchResp,
          page: page,
          totalCount: gatewayLogCount,
        });
      })
      .catch(err => reject(err));
  });
}

export default function GatewayLogs() {
  const classes = useStyles();
  const {match} = useRouter();
  const networkId: string = nullthrows(match.params.networkId);
  const gatewayId: string = nullthrows(match.params.gatewayId);
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [logCount, setLogCount] = useState(0);
  const [endDate, setEndDate] = useState(moment());
  const [actionQuery, setActionQuery] = useState<ActionQuery>({});
  const enqueueSnackbar = useEnqueueSnackbar();
  const tableRef = useRef(null);

  const startEnd = useMemo(() => {
    const [delta, unit, format] = getStep(startDate, endDate);
    return {
      start: startDate,
      end: endDate,
      delta: delta,
      unit: unit,
      format: format,
    };
  }, [startDate, endDate]);

  return (
    <div className={classes.dashboardRoot}>
      <Grid container align="top" alignItems="flex-start">
        <Grid container justify="space-between" spacing={3}>
          <Grid item xs={6}>
            <Text>
              <ListAltIcon />
              Logs({logCount})
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
                  onChange={val => {
                    setStartDate(val);
                    tableRef.current && tableRef.current.onQueryChange();
                  }}
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
                  onChange={val => {
                    setEndDate(val);
                    tableRef.current && tableRef.current.onQueryChange();
                  }}
                />
              </Grid>
              <Grid item>
                <Button
                  variant="contained"
                  color="primary"
                  className={classes.button}
                  startIcon={<LaunchIcon />}
                  onClick={() =>
                    exportLogs(
                      networkId,
                      gatewayId,
                      0,
                      MAX_PAGE_ROW_COUNT,
                      startDate,
                      endDate,
                      actionQuery,
                      enqueueSnackbar,
                    )
                  }>
                  Export
                </Button>
              </Grid>
            </Grid>
          </Grid>
          <Grid item xs={12}>
            <LogChart {...startEnd} setLogCount={setLogCount} />
          </Grid>
        </Grid>
        <Grid item xs={12}>
          <ActionTable
            title={'Logs'}
            tableRef={tableRef}
            data={(query: ActionQuery) => {
              setActionQuery(query);
              return handleLogQuery(
                networkId,
                gatewayId,
                0,
                MAX_PAGE_ROW_COUNT,
                startDate,
                endDate,
                query,
                enqueueSnackbar,
              );
            }}
            columns={LOG_COLUMNS}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [5, 10],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}
