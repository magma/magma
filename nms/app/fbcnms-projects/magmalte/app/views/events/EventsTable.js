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

import type {event as MagmaEvent} from '@fbcnms/magma-api';

import ActionTable from '../../components/ActionTable';
import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MyLocationIcon from '@material-ui/icons/MyLocation';
import React from 'react';
import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

import {filter, flatMap, map, some} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
import {useRouter} from '@fbcnms/ui/hooks';
import {useState} from 'react';

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
  eventDetailValue: {
    'max-width': '500px',
    overflow: 'scroll',
  },
  dashboardRoot: {
    margin: theme.spacing(5),
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

export const magmaEventTypes = Object.freeze({
  NETWORK: 1,
  GATEWAY: 2,
  SUBSCRIBER: 3,
});

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
  const [expanded, setExpanded] = useState(false);
  if (props.rowData.value) {
    for (const [key, value] of Object.entries(props.rowData.value)) {
      eventDetails[key] = value;
    }
  }
  return (
    <ExpansionPanel
      elevation={0}
      expanded={expanded}
      onChange={() => setExpanded(!expanded)}>
      <ExpansionPanelSummary
        expandIcon={<ExpandMoreIcon />}
        aria-controls="panel1bh-content">
        {Object.keys(eventDetails).join(', ')}
      </ExpansionPanelSummary>
      <ExpansionPanelDetails>
        <table>
          <tbody>
            {Object.entries(eventDetails).map((entry, i) => (
              <tr key={i}>
                <td>{entry[0]}: </td>
                <td className={classes.eventDetailValue}>
                  {JSON.stringify(entry[1])}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </ExpansionPanelDetails>
    </ExpansionPanel>
  );
}

export default function EventsTable({
  eventTypes,
  eventKey,
  sz,
}: {
  eventTypes: number,
  eventKey?: string,
  sz: 'sm' | 'md' | 'lg',
}) {
  const classes = useStyles();
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const magmadEvents = useMagmaAPI(
    MagmaV1API.getEventsByNetworkIdByStreamName,
    {
      networkId: nullthrows(match.params.networkId),
      streamName: streamNameMagmad,
    },
    undefined,
  );
  const sessiondEvents = useMagmaAPI(
    MagmaV1API.getEventsByNetworkIdByStreamName,
    {
      networkId: nullthrows(match.params.networkId),
      streamName: streamNameSessiond,
    },
    undefined,
  );

  let renderedEvents = [];

  switch (eventTypes) {
    case magmaEventTypes.NETWORK:
      renderedEvents = [magmadEvents, sessiondEvents];
      break;
    case magmaEventTypes.GATEWAY:
      renderedEvents = [magmadEvents, sessiondEvents];
      break;
    case magmaEventTypes.SUBSCRIBER:
      renderedEvents = [sessiondEvents];
      break;
  }

  if (some(renderedEvents, s => s.isLoading)) {
    return <LoadingFiller />;
  }

  if (some(renderedEvents, s => s.isLoading)) {
    return <LoadingFiller />;
  }

  const errored = filter(renderedEvents, s => s.error);
  const errors = map(errored, e => e.error.message);
  errors.forEach(err => enqueueSnackbar(err, {variant: 'error'}));

  const loaded = filter(renderedEvents, s => s.response);
  const unfiltered: Array<MagmaEvent> = flatMap(loaded, s => s.response);
  let events = unfiltered;
  if (eventTypes === magmaEventTypes.GATEWAY) {
    events = filter(unfiltered, s => s.hardware_id === eventKey);
  }

  const eventRows: Array<EventRowType> = events.map(event => {
    return {
      ts: event.timestamp,
      eventType: event.event_type,
      eventDescription: getEventDescription(event),
      value: event.value,
      hardwareID: event.hardware_id,
      tag: event.tag,
    };
  });

  return (
    <>
      {sz === 'sm' && (
        <ActionTable
          title=""
          data={eventRows}
          columns={[
            {title: 'Timestamp', field: 'ts', type: 'datetime'},
            {title: 'EventType', field: 'eventType'},
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [5],
            toolbar: false,
          }}
        />
      )}
      {sz === 'md' && (
        <ActionTable
          title="Events"
          titleIcon={MyLocationIcon}
          data={eventRows}
          columns={[
            {title: 'Timestamp', field: 'ts', type: 'datetime'},
            {title: 'Event Type', field: 'eventType'},
            {title: 'Event Description', field: 'eventDescription'},
            {
              title: 'More Details',
              field: 'eventDescription',
              render: rowData => <ExpandEvent rowData={rowData} />,
            },
          ]}
          options={{
            actionsColumnIndex: -1,
            pageSizeOptions: [10, 20],
          }}
        />
      )}
      {sz === 'lg' && (
        <div className={classes.dashboardRoot}>
          <ActionTable
            title="Events"
            titleIcon={MyLocationIcon}
            data={eventRows}
            columns={[
              {title: 'Timestamp', field: 'ts', type: 'datetime'},
              {title: 'Event Type', field: 'eventType'},
              {title: 'Event Description', field: 'eventDescription'},
              {
                title: 'More Details',
                field: 'eventDescription',
                render: rowData => <ExpandEvent rowData={rowData} />,
              },
            ]}
            options={{
              actionsColumnIndex: -1,
              pageSizeOptions: [10, 20],
            }}
          />
        </div>
      )}
    </>
  );
}
