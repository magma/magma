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
 */

import ActionTable from '../../components/ActionTable';
import AutorefreshCheckbox, {
  useRefreshingDateRange,
} from '../../components/AutorefreshCheckbox';
import Button from '@material-ui/core/Button';
import CardTitleRow from '../../components/layout/CardTitleRow';
import Grid from '@material-ui/core/Grid';
import LaunchIcon from '@material-ui/icons/Launch';
import ListAltIcon from '@material-ui/icons/ListAlt';
import LogChart from './GatewayLogChart';
import MagmaAPI from '../../../api/MagmaAPI';
import React, {useMemo, useRef, useState} from 'react';
import Text from '../../theme/design-system/Text';
import nullthrows from '../../../shared/util/nullthrows';
import {CsvBuilder} from 'filefy';
import {DateTimePicker} from '@material-ui/pickers';
import {MaterialTableProps} from '@material-table/core';
import {OptionsObject} from 'notistack';
import {Theme} from '@material-ui/core/styles';
import {colors} from '../../theme/default';
import {getErrorMessage} from '../../util/ErrorUtils';
import {getStep} from '../../components/CustomMetrics';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '../../hooks/useSnackbar';
import {useParams} from 'react-router-dom';
import type {ActionQuery} from '../../components/ActionTable';

// elastic search pagination through 'from' mechanism has a 10000 row limit
// we have to use a different mechanism in case we want to go higher, we should
// use search_after
// https://www.elastic.co/guide/en/elasticsearch/reference/6.8/search-request-search-after.html
const MAX_PAGE_ROW_COUNT = 10000;
const EXPORT_DELIMITER = ',';
const LOG_COLUMNS = [
  // eslint-disable-next-line prettier/prettier
  {title: 'Date', field: 'date', type: 'datetime', width: 200, filtering: false,} as const,
  {title: 'Service', field: 'service', width: 200} as const,
  {title: 'Tag', field: 'tag', width: 200} as const,
  {title: 'Type', field: 'logType', width: 200, filtering: false} as const,
  {title: 'Output', field: 'output', filtering: false} as const,
];

const useStyles = makeStyles<Theme>(theme => ({
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

function buildQueryFilters(query: ActionQuery, gatewayId: string) {
  let simpleQuery = query.search ?? '';
  let fields: string | undefined;
  const filters = [`gateway_id:${gatewayId}`];

  if (query.search === '' && query.filters.length === 1) {
    // for this case we can do a regex search
    simpleQuery = query.filters[0].value as string;
    fields = query.filters[0].column.field === 'service' ? 'ident' : 'tag';
  } else if (query.filters.length > 0) {
    query.filters.forEach(filter => {
      switch (filter.column.field) {
        case 'service':
          filters.push(`ident:${filter.value as string}`);
          break;

        case 'tag':
          filters.push(`tag:${filter.value as string}`);
          break;
      }
    });
  }

  return {
    simpleQuery,
    fields,
    filters: filters.join(','),
  };
}

async function searchLogs(
  networkId: string,
  gatewayId: string,
  from: number,
  size: number,
  start: moment.Moment,
  end: moment.Moment,
  query: ActionQuery,
) {
  const logs = (
    await MagmaAPI.logs.networksNetworkIdLogsSearchGet({
      networkId: networkId,
      from: from.toString(),
      size: size.toString(),
      start: start.toISOString(),
      end: end.toISOString(),
      ...buildQueryFilters(query, gatewayId),
    })
  ).data;

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
  networkId: string,
  gatewayId: string,
  from: number,
  size: number,
  start: moment.Moment,
  end: moment.Moment,
  query: ActionQuery,
  enqueueSnackbar: (message: string, config: OptionsObject) => string | number,
) {
  try {
    const logRows = await searchLogs(
      networkId,
      gatewayId,
      0,
      MAX_PAGE_ROW_COUNT,
      start,
      end,
      query,
    );
    const data = logRows.map(rowData =>
      // TODO[ts-migration]: the addRows in line 187 requires a string array
      LOG_COLUMNS.map(columnDef => rowData[columnDef.field] as string),
    );
    const currTs = Date.now();
    new CsvBuilder(`logs_${currTs}.csv`)
      .setDelimeter(EXPORT_DELIMITER)
      .setColumns(LOG_COLUMNS.map(columnDef => columnDef.title))
      .addRows(data)
      .exportFile();
  } catch (error) {
    enqueueSnackbar(getErrorMessage(error) ?? 'error retrieving logs', {
      variant: 'error',
    });
  }
}

async function handleLogQuery(
  networkId: string,
  gatewayId: string,
  from: number,
  size: number,
  start: moment.Moment,
  end: moment.Moment,
  query: ActionQuery,
) {
  try {
    const countReq = (
      await MagmaAPI.logs.networksNetworkIdLogsCountGet({
        networkId: networkId,
        start: start.toISOString(),
        end: end.toISOString(),
        ...buildQueryFilters(query, gatewayId),
      })
    ).data;

    const searchReq = searchLogs(
      networkId,
      gatewayId,
      query.page * query.pageSize,
      query.pageSize,
      start,
      end,
      query,
    );

    const [countResp, searchResp] = await Promise.all([countReq, searchReq]);
    let gatewayLogCount = countResp;

    if (gatewayLogCount > MAX_PAGE_ROW_COUNT) {
      gatewayLogCount = MAX_PAGE_ROW_COUNT;
    }

    const page =
      gatewayLogCount < query.page * query.pageSize
        ? gatewayLogCount / query.pageSize
        : query.page;

    return {
      data: searchResp,
      page: page,
      totalCount: gatewayLogCount,
    };
  } catch (error) {
    throw getErrorMessage(error, 'error retrieving logs');
  }
}

export default function GatewayLogs() {
  const classes = useStyles();
  const params = useParams();
  const networkId: string = nullthrows(params.networkId);
  const gatewayId: string = nullthrows(params.gatewayId);
  const [logCount, setLogCount] = useState(0);
  const [actionQuery, setActionQuery] = useState<ActionQuery>(
    {} as ActionQuery,
  );
  const enqueueSnackbar = useEnqueueSnackbar();
  const tableRef: MaterialTableProps<any>['tableRef'] = useRef(null);
  const [isAutoRefreshing, setIsAutoRefreshing] = useState(true);
  const {startDate, endDate, setStartDate, setEndDate} = useRefreshingDateRange(
    isAutoRefreshing,
    30000,
    () => {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-call,@typescript-eslint/no-unsafe-member-access
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
                setStartDate(val!);
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
                setEndDate(val!);
                setIsAutoRefreshing(false);
              }}
            />
          </Grid>
          <Grid item>
            <Button
              variant="contained"
              color="primary"
              startIcon={<LaunchIcon />}
              onClick={() => {
                void exportLogs(
                  networkId,
                  gatewayId,
                  0,
                  MAX_PAGE_ROW_COUNT,
                  startDate,
                  endDate,
                  actionQuery,
                  enqueueSnackbar,
                );
              }}>
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
