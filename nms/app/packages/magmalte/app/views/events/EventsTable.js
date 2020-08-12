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
import type {event as MagmaEvent} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import CardTitleRow from '../../components/layout/CardTitleRow';
import EventChart from './EventChart';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import React from 'react';
import Text from '../../theme/design-system/Text';
import moment from 'moment';
import nullthrows from '@fbcnms/util/nullthrows';

import {DateTimePicker} from '@material-ui/pickers';
import {colors} from '../../theme/default';
import {getStep} from '../../components/CustomHistogram';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useMemo, useRef, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';

const useStyles = makeStyles(theme => ({
  header: {
    margin: '10px',
    display: 'flex',
    justifyContent: 'space-between',
  },
  buttons: {
    display: 'flex',
    justifyContent: 'flex-end',
    flexDirection: 'row',
  },
  paper: {
    margin: theme.spacing(3),
  },
  importButton: {
    marginRight: '8px',
  },
  eventDetailTable: {
    // maxWidth: <value>, //TODO: This should be set to the parent table size
    width: '100%',
    padding: theme.spacing(1),
  },
  eventDetailLabel: {
    verticalAlign: 'top',
    fontWeight: 'bold',
  },
  eventDetailValue: {
    overflowWrap: 'break-word',
    maxWidth: '60vw', //TODO: Remove this when sizing added to `eventDetailTable`.
  },
  dashboardRoot: {
    margin: theme.spacing(5),
  },
  dateTimeText: {
    color: colors.primary.comet,
  },
}));

function getEventDescription(event) {
  switch (event.event_type) {
    case 'processed_updates':
      return 'Updates streamed from orchestrator were processed by the gateway';
    case 'updated_stored_mconfig':
      return "The gateway's stored mconfig was updated from the orchestrator";
    case 'session_created':
      return 'Subscriber session was created';
    case 'session_terminated':
      return 'Subscriber session was terminated';
    default:
      return event.event_type;
  }
}

const streamNameMagmad = 'magmad';
const streamNameSessiond = 'sessiond';

export type magmaEventStream = 'NETWORK' | 'GATEWAY' | 'SUBSCRIBER';
type EventRowType = {
  ts: string,
  eventType: string,
  eventDescription: string,
  value: {},
  hardwareID: string,
  tag: string,
};

type EventDescriptionProps = {
  rowData: EventRowType,
};

function ExpandEvent(props: EventDescriptionProps) {
  const classes = useStyles();
  const eventDetails = {
    hardware_id: props.rowData.hardwareID,
    tag: props.rowData.tag,
  };

  if (props.rowData.value) {
    for (const [key, value] of Object.entries(props.rowData.value)) {
      eventDetails[key] = value;
    }
  }

  return (
    <table className={classes.eventDetailTable}>
      <tbody>
        {Object.entries(eventDetails).map((entry, i) => (
          <tr key={i}>
            <td className={classes.eventDetailLabel}>{entry[0]}: </td>
            <td className={classes.eventDetailValue}>
              {JSON.stringify(entry[1])}
            </td>
          </tr>
        ))}
      </tbody>
    </table>
  );
}

function handleEventQuery(
  networkId,
  streams,
  tags,
  q,
  from,
  start,
  end,
  enqueueSnackbar,
) {
  return new Promise(async (resolve, reject) => {
    try {
      const eventCount = await MagmaV1API.getEventsByNetworkIdAboutCount({
        networkId: networkId,
        streams: streams,
        tags: tags,
        events: q.search !== '' ? q.search : undefined,
        from,
        start: start.toISOString(),
        end: end.toISOString(),
      });

      const eventResp = await MagmaV1API.getEventsByNetworkId({
        networkId: networkId,
        streams: streams,
        tags: tags,
        events: q.search !== '' ? q.search : undefined,
        from: (q.page * q.pageSize).toString(),
        size: q.pageSize.toString(),
        start: start.toISOString(),
        end: end.toISOString(),
      });
      const page =
        eventCount < q.page * q.pageSize ? eventCount / q.pageSize : q.page;

      // flowlint-next-line unclear-type:off
      const unfiltered: Array<MagmaEvent> = (eventResp: any);
      const data = unfiltered.map(event => {
        return {
          ts: event.timestamp,
          eventType: event.event_type,
          eventDescription: getEventDescription(event),
          value: event.value,
          hardwareID: event.hardware_id,
          tag: event.tag,
        };
      });
      resolve({
        data: data,
        page: page,
        totalCount: eventCount,
      });
    } catch (e) {
      enqueueSnackbar(e, {variant: 'error'});
      reject(e);
    }
  });
}

export default function EventsTable({
  eventStream,
  tags,
  sz,
}: {
  eventStream: magmaEventStream,
  tags?: string,
  sz: 'sm' | 'md' | 'lg',
}) {
  const classes = useStyles();
  const [startDate, setStartDate] = useState(moment().subtract(3, 'hours'));
  const [endDate, setEndDate] = useState(moment());
  const [eventCount, setEventCount] = useState(0);
  const tableRef = useRef(null);
  const enqueueSnackbar = useEnqueueSnackbar();
  const {match} = useRouter();
  const networkId = nullthrows(match.params.networkId);
  const streams =
    eventStream === 'SUBSCRIBER'
      ? streamNameSessiond
      : streamNameMagmad + ',' + streamNameSessiond;

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

  function DateFilter() {
    return (
      <Grid container justify="flex-end" alignItems="center" spacing={1}>
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
              tableRef.current && tableRef.current.onQueryChange();
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
              tableRef.current && tableRef.current.onQueryChange();
            }}
          />
        </Grid>
      </Grid>
    );
  }

  return (
    <>
      {sz === 'sm' && (
        <ActionTable
          title=""
          tableRef={tableRef}
          data={(query: ActionQuery) => {
            return handleEventQuery(
              networkId,
              streams,
              tags,
              query,
              0,
              startDate,
              endDate,
              enqueueSnackbar,
            );
          }}
          columns={[
            {title: 'Timestamp', field: 'ts', type: 'datetime', width: 200},
            {title: 'EventType', field: 'eventType'},
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [5],
            toolbar: false,
          }}
          detailPanel={[
            {
              icon: ExpandMore,
              openIcon: ExpandLess,
              render: rowData => {
                return <ExpandEvent rowData={rowData} />;
              },
            },
          ]}
        />
      )}
      {sz === 'md' && (
        <ActionTable
          tableRef={tableRef}
          data={(query: ActionQuery) => {
            return handleEventQuery(
              networkId,
              streams,
              tags,
              query,
              0,
              startDate,
              endDate,
              enqueueSnackbar,
            );
          }}
          columns={[
            {title: 'Timestamp', field: 'ts', type: 'datetime', width: 200},
            {title: 'Event Type', field: 'eventType', width: 200},
            {title: 'Event Description', field: 'eventDescription'},
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [10, 20],
          }}
          detailPanel={[
            {
              icon: ExpandMore,
              openIcon: ExpandLess,
              render: rowData => {
                return <ExpandEvent rowData={rowData} />;
              },
            },
          ]}
        />
      )}
      {sz === 'lg' && (
        <div className={classes.dashboardRoot}>
          <Grid container spacing={4}>
            <Grid item xs={12}>
              <CardTitleRow
                icon={MyLocationIcon}
                label={`Events (${eventCount})`}
                filter={DateFilter}
              />
              <EventChart
                {...startEnd}
                setEventCount={setEventCount}
                streams={streams}
                tags={tags ?? ''}
              />
            </Grid>
            <Grid item xs={12}>
              <ActionTable
                tableRef={tableRef}
                toolbar={{
                  searchTooltip: 'Search Event Types',
                }}
                data={(query: ActionQuery) => {
                  return handleEventQuery(
                    networkId,
                    streams,
                    tags,
                    query,
                    0,
                    startDate,
                    endDate,
                    enqueueSnackbar,
                  );
                }}
                columns={[
                  {
                    title: 'Timestamp',
                    field: 'ts',
                    type: 'datetime',
                    width: 200,
                  },
                  {title: 'Event Type', field: 'eventType', width: 200},
                  {title: 'Event Description', field: 'eventDescription'},
                ]}
                options={{
                  actionsColumnIndex: -1,
                  pageSizeOptions: [10, 20],
                }}
                detailPanel={[
                  {
                    icon: ExpandMore,
                    openIcon: ExpandLess,
                    render: rowData => {
                      return <ExpandEvent rowData={rowData} />;
                    },
                  },
                ]}
              />
            </Grid>
          </Grid>
        </div>
      )}
    </>
  );
}
