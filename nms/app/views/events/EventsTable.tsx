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

import ActionTable, {TableRef} from '../../components/ActionTable';
import AutorefreshCheckbox from '../../components/AutorefreshCheckbox';
import CardTitleRow from '../../components/layout/CardTitleRow';
import EventChart from './EventChart';
import ExpandLess from '@material-ui/icons/ExpandLess';
import ExpandMore from '@material-ui/icons/ExpandMore';
import Grid from '@material-ui/core/Grid';
import MagmaAPI from '../../../api/MagmaAPI';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import React from 'react';
import Text from '../../theme/design-system/Text';
import moment from 'moment';
import nullthrows from '../../../shared/util/nullthrows';
import {DateTimePicker} from '@material-ui/pickers';
import {Event as MagmaEvent} from '../../../generated-ts';
import {MaterialTableProps} from '@material-table/core';
import {Theme} from '@material-ui/core/styles';
import {colors} from '../../theme/default';
import {getStep} from '../../components/CustomMetrics';
import {makeStyles} from '@material-ui/styles';
import {useCallback} from 'react';
import {useEffect, useMemo, useRef, useState} from 'react';
import {useParams} from 'react-router-dom';
import {useRefreshingDateRange} from '../../components/AutorefreshCheckbox';
import type {ActionQuery} from '../../components/ActionTable';

const useStyles = makeStyles<Theme>(theme => ({
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

function getEventDescription(event: {event_type: string}) {
  switch (event.event_type) {
    case 'processed_updates':
      return 'Updates streamed from orchestrator were processed by the gateway';
    case 'updated_stored_mconfig':
      return "The gateway's stored mconfig was updated from the orchestrator";
    case 'session_created':
      return 'Subscriber session was created';
    case 'session_terminated':
      return 'Subscriber session was terminated';
    case 'attach_success':
      return 'UE attaches successfully';
    case 'detach_success':
      return 'UE detaches successfully';
    default:
      return event.event_type;
  }
}

export type magmaEventStream = 'NETWORK' | 'GATEWAY' | 'SUBSCRIBER';
type EventRowType = {
  ts: string;
  eventType: string;
  eventDescription: string;
  value: Record<string, any>;
  hardwareId: string;
  tag: string;
};

export const EVENT_STREAM = {
  NETWORK: 'NETWORK',
  GATEWAY: 'GATEWAY',
  SUBSCRIBER: 'SUBSCRIBER',
} as const;

type EventDescriptionProps = {
  rowData: EventRowType;
};

function ExpandEvent(props: EventDescriptionProps) {
  const classes = useStyles();
  const eventDetails: Record<string, any> = {
    hardware_id: props.rowData.hardwareId,
    tag: props.rowData.tag,
  };

  if (props.rowData.value) {
    for (const [key, value] of Object.entries(props.rowData.value)) {
      // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
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

function buildEventQueryFromFilters(q: ActionQuery) {
  const queryFilters: Record<string, any> = {};

  if (q.filters !== undefined) {
    q.filters.forEach(filter => {
      switch (filter.column.field) {
        case 'streamName':
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          queryFilters['streams'] = filter.value;
          break;

        case 'eventType':
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          queryFilters['events'] = filter.value;
          break;

        case 'tag':
          // eslint-disable-next-line @typescript-eslint/no-unsafe-assignment
          queryFilters['tags'] = filter.value;
          break;
      }
    });
  }

  return queryFilters;
}

async function handleEventQuery(
  networkId: string,
  hardwareId: string | undefined,
  streams: string,
  tags: string,
  query: ActionQuery,
  from: number,
  start: moment.Moment,
  end: moment.Moment,
): Promise<{data: Array<EventRowType>; page: number; totalCount: number}> {
  const filters = buildEventQueryFromFilters(query);
  try {
    const eventCount = (
      await MagmaAPI.events.eventsNetworkIdAboutCountGet({
        networkId: networkId,
        streams: streams,
        hwIds: hardwareId,
        tags: tags,
        // TODO[ts-migration] "from" does not appear in API
        // @ts-ignore
        from,
        start: start.toISOString(),
        end: end.toISOString(),
        ...filters,
      })
    ).data;
    const eventResp = (
      await MagmaAPI.events.eventsNetworkIdGet({
        networkId: networkId,
        hwIds: hardwareId,
        streams: streams,
        tags: tags,
        from: (query.page * query.pageSize).toString(),
        size: query.pageSize.toString(),
        start: start.toISOString(),
        end: end.toISOString(),
        ...filters,
      })
    ).data;
    const page =
      eventCount < query.page * query.pageSize
        ? eventCount / query.pageSize
        : query.page;
    // TODO[ts-migration] There is a serious type mismatch here. Investigate!
    const unfiltered = (eventResp as unknown) as Array<MagmaEvent>;
    const data = unfiltered.map(event => {
      return {
        ts: event.timestamp,
        streamName: event.stream_name,
        eventType: event.event_type,
        eventDescription: getEventDescription(event),
        value: event.value,
        hardwareId: event.hardware_id,
        tag: event.tag,
      };
    });
    return {
      data: data,
      page: page,
      totalCount: eventCount,
    };
  } catch (error) {
    if (error instanceof Error && error.message) {
      throw error;
    } else {
      throw new Error('error retrieving events');
    }
  }
}

type EventTableProps = {
  eventStream: magmaEventStream;
  tags?: string;
  hardwareId?: string;
  sz: 'sm' | 'md' | 'lg';
  inStartDate?: moment.Moment;
  inEndDate?: moment.Moment;
  isAutoRefreshing?: boolean;
};

export default function EventsTable(props: EventTableProps) {
  const {hardwareId, eventStream, sz} = props;
  const classes = useStyles();
  const [eventCount, setEventCount] = useState(0);
  const tableRef: TableRef = useRef();
  const params = useParams();
  const networkId = nullthrows(params.networkId);
  const streams = '';

  const buildTags = (tags: string) => {
    let allTags = tags;
    const tagsDelimter = ',';

    if (eventStream == EVENT_STREAM.SUBSCRIBER) {
      // sessionD requires tag to include the prefix IMSI together with n digits but mme doesn't require the prefix IMSI
      allTags = [tags, tags.replace(/IMSI/, '')].join(tagsDelimter);
    }

    return allTags;
  };

  const tags = buildTags(props.tags ?? '');
  const [isAutoRefreshing, setIsAutoRefreshing] = useState(
    props.isAutoRefreshing ?? false,
  );
  const {startDate, endDate, setStartDate, setEndDate} = useRefreshingDateRange(
    isAutoRefreshing,
    30000,
    useCallback(() => {
      tableRef.current?.onQueryChange();
    }, []),
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
  useEffect(() => {
    if (props.inStartDate) {
      setStartDate(props.inStartDate);
    }
  }, [props.inStartDate, setStartDate]);

  useEffect(() => {
    if (props.inEndDate) {
      setEndDate(props.inEndDate);
    }
  }, [props.inEndDate, setEndDate]);

  let actionTableOptions: MaterialTableProps<EventRowType>['options'] = {
    actionsColumnIndex: -1,
    pageSize: 5,
    pageSizeOptions: [5, 10, 20],
    toolbar: false,
  };

  let actionColumns: MaterialTableProps<EventRowType>['columns'] = [
    {
      title: 'Timestamp',
      field: 'ts',
      type: 'datetime',
      width: 200,
      filtering: false,
    },
    {
      title: 'Stream Name',
      field: 'streamName',
      width: 200,
    },
    {
      title: 'Event Type',
      field: 'eventType',
      width: 200,
    },
  ];

  if (sz !== 'sm') {
    actionColumns = [
      ...actionColumns,
      {
        title: 'Tag',
        field: 'tag',
      },
      {
        title: 'Event Description',
        field: 'eventDescription',
        filtering: false,
      },
    ];
    actionTableOptions = {...actionTableOptions, pageSize: 10, filtering: true};
  }

  const actionTable = (
    <ActionTable
      tableRef={tableRef}
      data={(query: ActionQuery) => {
        return handleEventQuery(
          networkId,
          hardwareId,
          streams,
          tags,
          query,
          0,
          startDate,
          endDate,
        );
      }}
      columns={actionColumns}
      options={actionTableOptions}
      detailPanel={[
        {
          icon: ExpandMore,
          openIcon: ExpandLess,
          render: ({rowData}) => {
            return <ExpandEvent rowData={rowData} />;
          },
        },
      ]}
    />
  );

  if (sz === 'sm' || sz === 'md') {
    return actionTable;
  }

  return (
    <div className={classes.dashboardRoot}>
      <Grid container spacing={4}>
        <Grid item xs={12}>
          <CardTitleRow
            icon={MyLocationIcon}
            label={`Events (${eventCount})`}
            filter={() => (
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
                      onChange={date => {
                        setStartDate(date!);
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
                      onChange={date => {
                        setEndDate(date!);
                        setIsAutoRefreshing(false);
                      }}
                    />
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
            )}
          />
          <EventChart
            {...startEnd}
            setEventCount={setEventCount}
            streams={streams}
            tags={tags ?? ''}
          />
        </Grid>
        <Grid item xs={12}>
          {actionTable}
        </Grid>
      </Grid>
    </div>
  );
}
