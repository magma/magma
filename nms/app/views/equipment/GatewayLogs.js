/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow strict-local
 * @format
 */
import type {ActionQuery} from '../../components/ActionTable';

import ActionTable from '../../components/ActionTable';
// $FlowFixMe migrated to typescript
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import Button from '@material-ui/core/Button';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import ListAltIcon from '@material-ui/icons/ListAlt';
import LogChart from './GatewayLogChart';
import MagmaV1API from '../../../generated/WebClient';
import React from 'react';
import Text from '../../theme/design-system/Text';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import {CsvBuilder} from 'filefy';
import {DateTimePicker} from '@material-ui/pickers';
import {colors} from '../../theme/default';
import {getStep} from '../../components/CustomMetrics';
import {makeStyles} from '@material-ui/styles';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../app/hooks/useSnackbar';
import {useMemo, useRef, useState} from 'react';
import {useParams} from 'react-router-dom';
// $FlowFixMe migrated to typescript
import {useRefreshingDateRange} from '../../components/AutorefreshCheckbox';

// elastic search pagination through 'from' mechanism has a 10000 row limit
// we have to use a different mechanism in case we want to go higher, we should
// use search_after
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-request-search-after.html
const MAX_PAGE_ROW_COUNT = 10000;
const EXPORT_DELIMITER = ',';
const LOG_COLUMNS = [
  {
    title: 'Date',
    field: 'date',
    type: 'datetime',
    width: 200,
    filtering: false,
  },
  {title: 'Service', field: 'service', width: 200},
  {title: 'Tag', field: 'tag', width: 200},
  {title: 'Type', field: 'logType', width: 200, filtering: false},
  {title: 'Output', field: 'output', filtering: false},
];

const useStyles = makeStyles(theme => ({
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  dateTimeText: {
    color: colors.primary.comet,
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

function buildQueryFilters(q: ActionQuery, gatewayId: string) {
  const logQuery = {
    simpleQuery: q.search ?? '',
    fields: undefined,
    filters: undefined,
  };
  const filters = [`gateway_id:${gatewayId}`];
  if (q.search === '' && q.filters.length === 1) {
    // for this case we can do a regex search
    logQuery.simpleQuery = q.filters[0].value;
    logQuery.fields = q.filters[0].column.field === 'service' ? 'ident' : 'tag';
  } else if (q.filters.length > 0) {
    q.filters.forEach((filter, _) => {
      switch (filter.column.field) {
        case 'service':
          filters.push(`ident:${filter.value}`);
          break;
        case 'tag':
          filters.push(`tag:${filter.value}`);
          break;
      }
    });
  }
  logQuery.filters = filters.join(',');
  return logQuery;
}

async function searchLogs(networkId, gatewayId, from, size, start, end, q) {
  const logs = await MagmaV1API.getNetworksByNetworkIdLogsSearch({
    networkId: networkId,
    from: from.toString(),
    size: size.toString(),
    start: start.toISOString(),
    end: end.toISOString(),
    ...buildQueryFilters(q, gatewayId),
  });

  return logs.filter(Boolean).map(elastic_hit => {
    const src = elastic_hit._source;
    const date = new Date(src['@timestamp'] ?? 0);
    const msg = src['message'];
    return {
      date: date,
      service: src['ident'] ?? '-',
      tag: src['tag'] ?? '-',
      logType: getLogType(msg ?? ''),
      output: msg,
    };
  });
}

async function exportLogs(
  networkId,
  gatewayId,
  from,
  size,
  start,
  end,
  q,
  enqueueSnackbar,
) {
  try {
    const logRows = await searchLogs(
      networkId,
      gatewayId,
      0,
      MAX_PAGE_ROW_COUNT,
      start,
      end,
      q,
    );
    const data = logRows.map(rowData =>
      LOG_COLUMNS.map(columnDef => rowData[columnDef.field]),
    );
    const currTs = Date.now();
    new CsvBuilder(`logs_${currTs}.csv`)
      .setDelimeter(EXPORT_DELIMITER)
      .setColumns(LOG_COLUMNS.map(columnDef => columnDef.title))
      .addRows(data)
      .exportFile();
  } catch (e) {
    enqueueSnackbar(e?.message ?? 'error retrieving logs', {variant: 'error'});
  }
}

function handleLogQuery(networkId, gatewayId, from, size, start, end, q) {
  return new Promise(async (resolve, reject) => {
    try {
      const countReq = MagmaV1API.getNetworksByNetworkIdLogsCount({
        networkId: networkId,
        start: start.toISOString(),
        end: end.toISOString(),
        ...buildQueryFilters(q, gatewayId),
      });

      const searchReq = searchLogs(
        networkId,
        gatewayId,
        (q.page * q.pageSize).toString(),
        q.pageSize.toString(),
        start,
        end,
        q,
      );

      const [countResp, searchResp] = await Promise.all([countReq, searchReq]);
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
    } catch (e) {
      reject(e?.message ?? 'error retrieving logs');
    }
  });
}

export default function GatewayLogs() {
  const classes = useStyles();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const gatewayId: string = nullthrows(params.gatewayId);
  const [logCount, setLogCount] = useState(0);
  const [actionQuery, setActionQuery] = useState<ActionQuery>({});
  const enqueueSnackbar = useEnqueueSnackbar();
  const tableRef = useRef(null);
  const [isAutoRefreshing, setIsAutoRefreshing] = useState(true);
  const {startDate, endDate, setStartDate, setEndDate} = useRefreshingDateRange(
    isAutoRefreshing,
    30000,
    () => {
      tableRef.current && tableRef.current.onQueryChange();
    },
  );

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

  function LogsFilter() {
    return (
      <>
        <Grid
          container
          justifyContent="flex-end"
          alignItems="center"
          spacing={1}>
          <Grid item>
            <Text variant="body3" className={classes.dateTimeText}>
              Filter By Date
            </Text>
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
                setIsAutoRefreshing(false);
              }}
            />
          </Grid>
          <Grid item>
            <Text variant="body3" className={classes.dateTimeText}>
              To
            </Text>
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
                setIsAutoRefreshing(false);
              }}
            />
          </Grid>
          <Grid item>
            <Button
              variant="contained"
              color="primary"
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
        <Grid
          container
          justifyContent="flex-end"
          alignItems="center"
          spacing={1}>
          <Grid item>
            <AutorefreshCheckbox
              autorefreshEnabled={isAutoRefreshing}
              onToggle={() => setIsAutoRefreshing(current => !current)}
            />
          </Grid>
        </Grid>
      </>
    );
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow
            icon={ListAltIcon}
            label={`Logs (${logCount})`}
            filter={LogsFilter}
          />
          <LogChart {...startEnd} setLogCount={setLogCount} />
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
              );
            }}
            columns={LOG_COLUMNS}
            options={{
              filtering: true,
              actionsColumnIndex: -1,
              pageSize: 10,
              pageSizeOptions: [10, 20],
            }}
          />
        </Grid>
      </Grid>
    </div>
  );
}
