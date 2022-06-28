/*
 * Copyright 2022 The Magma Authors.
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
import Button from '@material-ui/core/Button';
import FormControl from '@material-ui/core/FormControl';
import Grid from '@material-ui/core/Grid';
import ListIcon from '@material-ui/icons/ListAlt';
import MenuItem from '@material-ui/core/MenuItem';
import OutlinedInput from '@material-ui/core/OutlinedInput';
import React, {useCallback, useRef, useState} from 'react';
import Select from '@material-ui/core/Select';
import moment from 'moment';
import nullthrows from '../../../shared/util/nullthrows';
import {KeyboardDateTimePicker} from '@material-ui/pickers';
import {isFinite} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useParams} from 'react-router-dom';

import ActionTable, {TableRef} from '../../components/ActionTable';
import AutorefreshCheckbox, {
  useRefreshingDateRange,
} from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import MagmaAPI from '../../../api/MagmaAPI';
import Text from '../../theme/design-system/Text';
import {REFRESH_INTERVAL} from '../../components/context/RefreshContext';
import {Theme} from '@material-ui/core/styles/createTheme';
import {colors} from '../../theme/default';

const useStyles = makeStyles<Theme>(theme => ({
  root: {
    margin: theme.spacing(3),
    flexGrow: 1,
  },
  logsDirectionFilter: {
    width: theme.spacing(19),
  },
  filterText: {
    color: colors.primary.comet,
  },
  filterActionsWrapper: {
    paddingRight: theme.spacing(4),
  },
  filterActions: {
    textAlign: 'right',
  },
}));

type LogsDirectionNullable = 'SAS' | 'DP' | 'CBSD' | null;

type LogsDirectionFilterProps = {
  value: LogsDirectionNullable;
  onChange: (value: LogsDirectionNullable) => void;
  selectProps: {
    'data-testid': string;
  };
};

function LogsDirectionFilter({
  value,
  onChange,
  selectProps,
}: LogsDirectionFilterProps) {
  const classes = useStyles();

  let parsedValue;
  switch (value) {
    case 'SAS':
      parsedValue = 1;
      break;
    case 'DP':
      parsedValue = 2;
      break;
    case 'CBSD':
      parsedValue = 3;
      break;
    default:
      parsedValue = 0;
      break;
  }

  return (
    <FormControl>
      <Select
        className={classes.logsDirectionFilter}
        value={parsedValue}
        onChange={({target}) => {
          let newValue: LogsDirectionNullable;
          switch (parseInt(target.value as string)) {
            case 1:
              newValue = 'SAS';
              break;
            case 2:
              newValue = 'DP';
              break;
            case 3:
              newValue = 'CBSD';
              break;
            default:
              newValue = null;
              break;
          }
          onChange(newValue);
        }}
        input={<OutlinedInput />}
        {...selectProps}>
        <MenuItem value={0}>Any</MenuItem>
        <MenuItem value={1}>SAS</MenuItem>
        <MenuItem value={2}>DP</MenuItem>
        <MenuItem value={3}>CBSD</MenuItem>
      </Select>
    </FormControl>
  );
}

function LogsList() {
  const classes = useStyles();

  const params = useParams();
  const networkId: string = nullthrows(params.networkId);

  const tableRef: TableRef = useRef();

  const [isAutoRefreshing, setIsAutoRefreshing] = useState(true);
  const {startDate, endDate, setStartDate, setEndDate} = useRefreshingDateRange(
    isAutoRefreshing,
    REFRESH_INTERVAL,
    () => {
      tableRef.current && tableRef.current.onQueryChange();
    },
  );
  const [serialNumber, setSerialNumber] = useState<string>('');
  const [fccId, setFccId] = useState<string>('');
  const [from, setFrom] = useState<LogsDirectionNullable>(null);
  const [to, setTo] = useState<LogsDirectionNullable>(null);
  const [responseCode, setResponseCode] = useState<string | null>(null);
  const [logName, setLogName] = useState<string>('');

  const getDataFn = useCallback(
    async (query: {page: number; pageSize: number}) => {
      const responseCodeParsed = parseInt(responseCode!);
      const response = (
        await MagmaAPI.logs.dpNetworkIdLogsGet({
          networkId,
          offset: query.page * query.pageSize,
          limit: query.pageSize,
          begin: startDate?.toISOString(),
          end: endDate?.toISOString(),
          serialNumber: serialNumber || undefined,
          fccId: fccId || undefined,
          type: logName || undefined,
          responseCode: isFinite(responseCodeParsed)
            ? responseCodeParsed
            : undefined,
          from: from || undefined,
          to: to || undefined,
        })
      ).data;

      const totalCount = response?.total_count || 0;

      const logsMapped = response?.logs?.length
        ? response?.logs.map(item => {
            return {
              body: item.body,
              fccId: item.fcc_id,
              from: item.from,
              to: item.to,
              serialNumber: item.serial_number,
              time: moment(item.time)?.toLocaleString(),
              type: item.type,
            };
          })
        : [];

      return {
        data: logsMapped,
        page: query.page,
        totalCount,
      };
    },
    [
      networkId,
      startDate,
      endDate,
      serialNumber,
      fccId,
      from,
      to,
      responseCode,
      logName,
    ],
  );

  const onSearch = useCallback(() => {
    tableRef.current && tableRef.current.onQueryChange();
  }, []);

  return (
    <div className={classes.root}>
      <Grid container justify="space-between" spacing={3}>
        <Grid item xs={12}>
          <CardTitleRow
            key="title"
            icon={ListIcon}
            label={`Domain Proxy Logs`}
          />

          <Grid container spacing={1} alignItems="center">
            <Grid item xs={1}>
              <Text variant="body2" className={classes.filterText}>
                Filter by:
              </Text>
            </Grid>

            <Grid
              item
              xs={3}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  Serial Number
                </Text>
              </Grid>
              <Grid item>
                <OutlinedInput
                  type="string"
                  inputProps={{
                    'data-testid': 'serial-number-input',
                  }}
                  value={serialNumber}
                  onChange={({target}) => setSerialNumber(target.value)}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={3}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  FCC ID
                </Text>
              </Grid>
              <Grid item>
                <OutlinedInput
                  inputProps={{
                    'data-testid': 'fcc-id-input',
                  }}
                  type="string"
                  value={fccId}
                  onChange={({target}) => setFccId(target.value)}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={2}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  From
                </Text>
              </Grid>
              <Grid item>
                <LogsDirectionFilter
                  value={from}
                  onChange={newValue => setFrom(newValue)}
                  selectProps={{
                    'data-testid': 'logs-direction-from-input',
                  }}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={3}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  Start Date
                </Text>
              </Grid>
              <Grid item>
                <KeyboardDateTimePicker
                  inputProps={{
                    'data-testid': 'start-date-input',
                  }}
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  maxDate={endDate}
                  disableFuture
                  value={startDate}
                  onChange={newValue => setStartDate(newValue as moment.Moment)}
                  format="yyyy/MM/DD HH:mm"
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={4}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  Response Code
                </Text>
              </Grid>
              <Grid item>
                <OutlinedInput
                  inputProps={{
                    'data-testid': 'response-code-input',
                  }}
                  type="number"
                  value={responseCode}
                  onChange={({target}) => setResponseCode(target.value)}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={3}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  Log Name
                </Text>
              </Grid>
              <Grid item>
                <OutlinedInput
                  inputProps={{
                    'data-testid': 'log-name-input',
                  }}
                  type="string"
                  value={logName}
                  onChange={({target}) => setLogName(target.value)}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={2}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  To
                </Text>
              </Grid>
              <Grid item>
                <LogsDirectionFilter
                  selectProps={{
                    'data-testid': 'logs-direction-to-input',
                  }}
                  value={to}
                  onChange={newValue => setTo(newValue)}
                />
              </Grid>
            </Grid>

            <Grid
              item
              xs={3}
              container
              spacing={1}
              alignItems="center"
              justify="flex-end">
              <Grid item>
                <Text variant="body3" className={classes.filterText}>
                  End Date
                </Text>
              </Grid>
              <Grid item>
                <KeyboardDateTimePicker
                  inputProps={{
                    'data-testid': 'end-date-input',
                  }}
                  autoOk
                  variant="inline"
                  inputVariant="outlined"
                  disableFuture
                  value={endDate}
                  onChange={newValue => setEndDate(newValue as moment.Moment)}
                  format="yyyy/MM/DD HH:mm"
                />
              </Grid>
            </Grid>
          </Grid>

          <Grid container spacing={1} className={classes.filterActionsWrapper}>
            <Grid item xs={12} className={classes.filterActions}>
              <AutorefreshCheckbox
                autorefreshEnabled={isAutoRefreshing}
                onToggle={() => setIsAutoRefreshing(!isAutoRefreshing)}
              />

              <Button
                variant="contained"
                color="primary"
                onClick={onSearch}
                data-testid="search-button">
                Search
              </Button>
            </Grid>
          </Grid>

          <ActionTable
            tableRef={tableRef}
            columns={[
              {
                title: 'From',
                field: 'from',
              },
              {
                title: 'To',
                field: 'to',
              },
              {
                title: 'Log Name',
                field: 'type',
              },
              {
                title: 'Message',
                field: 'body',
              },
              {
                title: 'Serial Number',
                field: 'serialNumber',
              },
              {
                title: 'FCC ID',
                field: 'fccId',
              },
              {
                title: 'Time',
                field: 'time',
              },
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSize: 20,
              pageSizeOptions: [20, 60, 100],
              search: false,
              sorting: false,
            }}
            data={getDataFn}
          />
        </Grid>
      </Grid>
    </div>
  );
}

export default LogsList;
