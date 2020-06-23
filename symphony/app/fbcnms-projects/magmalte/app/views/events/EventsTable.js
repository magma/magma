/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {event as MagmaEvent} from '@fbcnms/magma-api';

import ExpandMoreIcon from '@material-ui/icons/ExpandMore';
import ExpansionPanel from '@material-ui/core/ExpansionPanel';
import ExpansionPanelDetails from '@material-ui/core/ExpansionPanelDetails';
import ExpansionPanelSummary from '@material-ui/core/ExpansionPanelSummary';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Paper from '@material-ui/core/Paper';
import React, {useState} from 'react';
import Table from '@material-ui/core/Table';
import TableBody from '@material-ui/core/TableBody';
import TableCell from '@material-ui/core/TableCell';
import TableFooter from '@material-ui/core/TableFooter';
import TableHead from '@material-ui/core/TableHead';
import TablePagination from '@material-ui/core/TablePagination';
import TableRow from '@material-ui/core/TableRow';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {filter, flatMap, map, slice, some} from 'lodash';
import {makeStyles} from '@material-ui/styles';
import {useEnqueueSnackbar} from '@fbcnms/ui/hooks/useSnackbar';
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
  eventDetailValue: {
    'max-width': '500px',
    overflow: 'scroll',
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

export default function EventsTable({
  eventTypes,
  gatewayHardwareId,
}: {
  eventTypes: magmaEventType,
  gatewayHardwareId: string,
}) {
  const {match} = useRouter();
  const enqueueSnackbar = useEnqueueSnackbar();
  const [rowsPerPage, setRowsPerPage] = useState(10);
  const [page, setPage] = useState(0);
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

  const errored = filter(renderedEvents, s => s.error);
  const errors = map(errored, e => e.error.message);
  errors.forEach(err => enqueueSnackbar(err, {variant: 'error'}));

  const loaded = filter(renderedEvents, s => s.response);
  const unfiltered: Array<MagmaEvent> = flatMap(loaded, s => s.response);
  let events = unfiltered;
  if (eventTypes === magmaEventTypes.GATEWAY) {
    events = filter(unfiltered, s => s.hardware_id === gatewayHardwareId);
  }

  const eventsStartIndex = page * rowsPerPage;
  const eventsEndIndex = eventsStartIndex + rowsPerPage;
  const rows = map(
    slice(events, eventsStartIndex, eventsEndIndex),
    (event, index: number) => <EventTableRow key={index} event={event} />,
  );

  return (
    <Paper elevation={2}>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>Timestamp</TableCell>
            <TableCell>Event Type</TableCell>
            <TableCell>Description</TableCell>
            <TableCell>More details</TableCell>
            <TableCell />
          </TableRow>
        </TableHead>
        <TableBody>{rows}</TableBody>
        {events.length === 0 && (
          <TableFooter>
            <TableRow>
              <TableCell colSpan="3">No events found</TableCell>
            </TableRow>
          </TableFooter>
        )}
      </Table>
      <TablePagination
        rowsPerPageOptions={[10, 50, 100]}
        component="div"
        count={events.length}
        rowsPerPage={rowsPerPage}
        page={page}
        onChangePage={(_, newPage) => setPage(newPage)}
        onChangeRowsPerPage={e => {
          setRowsPerPage(parseInt(e.target.value, 10));
          setPage(0);
        }}
      />
    </Paper>
  );
}

type Props = {
  event: MagmaEvent,
};

function EventTableRow(props: Props) {
  const classes = useStyles();
  const {event} = props;
  const [expanded, setExpanded] = useState(false);

  const eventDetails = {
    hardware_id: event.hardware_id,
    tag: event.tag,
  };
  for (const [key, value] of Object.entries(event.value)) {
    eventDetails[key] = value;
  }

  return (
    <TableRow>
      <TableCell>{new Date(event.timestamp).toString()}</TableCell>
      <TableCell>{event.event_type}</TableCell>
      <TableCell>{getEventDescription(event)}</TableCell>
      <TableCell>
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
      </TableCell>
    </TableRow>
  );
}
