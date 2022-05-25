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

import EventsTable from '../EventsTable';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
// $FlowFixMe migrated to typescript
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../hooks/useSnackbar');

const mockEvents = [
  {
    event_type: 'attach_success',
    hardware_id: 'fd9cb997-25b1-4d0a-b7d8-fea605558311',
    stream_name: 'mme',
    tag: '001011234560000',
    timestamp: '2021-07-24T12:01:36.863752622+00:00',
    value: {
      imsi: '001011234560000',
    },
  },
  {
    event_type: 'session_created',
    hardware_id: 'fd9cb997-25b1-4d0a-b7d8-fea605558311',
    stream_name: 'sessiond',
    tag: 'IMSI001011234560000',
    timestamp: '2021-07-24T12:05:36.978641891+00:00',
    value: {
      apn: 'oai.ipv4',
      charging_characteristics: '',
      imei:
        '\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0000\u0000',
      imsi: 'IMSI001011234560000',
      ip_addr: '192.168.130.10',
      ipv6_addr: '',
      mac_addr: '',
      msisdn: '',
      pdp_start_time: 1627128336,
      session_id: 'IMSI001011234560000-522417',
      spgw_ip: '10.22.10.92',
      user_location: ' 92 01 f1 10 00 01 00 f1 10 00 00 00 01',
    },
  },
  {
    event_type: 'session_terminated',
    hardware_id: 'fd9cb997-25b1-4d0a-b7d8-fea605558311',
    stream_name: 'sessiond',
    tag: 'IMSI001011234560000',
    timestamp: '2021-07-27T01:59:40.343828806+00:00',
    value: {
      apn: 'oai.ipv4',
      cause_for_rec_closing: 0,
      charging_characteristics: '',
      charging_rx: 0,
      charging_tx: 0,
      duration: 23353,
      imei:
        '\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0001\u0000\u0000',
      imsi: 'IMSI001011234560000',
      ip_addr: '192.168.156.10',
      ipv6_addr: '',
      list_of_service_data: [],
      mac_addr: '',
      monitoring_rx: 0,
      monitoring_tx: 0,
      msisdn: '',
      pdp_end_time: 1627351175,
      pdp_start_time: 1627327822,
      record_sequence_number: 1,
      session_id: 'IMSI001011234560000-522417',
      spgw_ip: '10.22.10.92',
      total_rx: 0,
      total_tx: 0,
      user_location: ' 92 01 f1 10 00 01 00 f1 10 00 00 00 01',
    },
  },
  {
    event_type: 'detach_success',
    hardware_id: 'fd9cb997-25b1-4d0a-b7d8-fea605558311',
    stream_name: 'mme',
    tag: '001011234560000',
    timestamp: '2021-07-27T08:34:02.531951334+00:00',
    value: {
      action: 'detach_accept_sent',
      imsi: '001011234560000',
    },
  },
];

describe('<EventsTable />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getEventsByNetworkIdAboutCount.mockResolvedValue(
      mockEvents.length,
    );
    MagmaAPIBindings.getEventsByNetworkId.mockResolvedValue(mockEvents);
  });

  const Wrapper = () => {
    return (
      <MemoryRouter
        initialEntries={[
          '/nms/test/subscribers/overview/config/IMSI0000000000/overview',
        ]}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <NetworkContext.Provider
              value={{
                networkId: 'test',
              }}>
              <Routes>
                <Route
                  path="/nms/:networkId/subscribers/overview/config/:subscriberId/overview"
                  element={
                    <EventsTable
                      eventStream={'SUBSCRIBER'}
                      tags={'IMSI001011234560000'}
                      sz={'md'}
                    />
                  }
                />
              </Routes>
            </NetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscriber Events Table', async () => {
    const {getAllByRole} = render(<Wrapper />);
    await wait();
    const mockQuery = {
      networkId: 'test',
      tags: 'IMSI001011234560000,001011234560000',
      streams: '',
      hwIds: undefined,
    };
    // verify that API is called with the correct tag
    expect(
      MagmaAPIBindings.getEventsByNetworkIdAboutCount,
    ).toHaveBeenCalledWith(expect.objectContaining(mockQuery));
    expect(MagmaAPIBindings.getEventsByNetworkId).toHaveBeenCalledWith(
      expect.objectContaining(mockQuery),
    );
    const rowItems = await getAllByRole('row');
    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Timestamp');
    expect(rowItems[0]).toHaveTextContent('Stream Name');
    expect(rowItems[0]).toHaveTextContent('Event Type');
    expect(rowItems[0]).toHaveTextContent('Tag');
    expect(rowItems[0]).toHaveTextContent('Event Description');

    // skipping second row because that is used for filtering events
    expect(rowItems[2]).toHaveTextContent(
      new Date(mockEvents[0].timestamp).toLocaleString(),
    );
    expect(rowItems[2]).toHaveTextContent('mme');
    expect(rowItems[2]).toHaveTextContent('attach_success');
    expect(rowItems[2]).toHaveTextContent('001011234560000');
    expect(rowItems[2]).toHaveTextContent('UE attaches successfully');

    expect(rowItems[3]).toHaveTextContent(
      new Date(mockEvents[1].timestamp).toLocaleString(),
    );
    expect(rowItems[3]).toHaveTextContent('sessiond');
    expect(rowItems[3]).toHaveTextContent('session_created');
    expect(rowItems[3]).toHaveTextContent('IMSI001011234560000');
    expect(rowItems[3]).toHaveTextContent('Subscriber session was created');

    expect(rowItems[4]).toHaveTextContent(
      new Date(mockEvents[2].timestamp).toLocaleString(),
    );
    expect(rowItems[4]).toHaveTextContent('sessiond');
    expect(rowItems[4]).toHaveTextContent('session_terminated');
    expect(rowItems[4]).toHaveTextContent('IMSI001011234560000');
    expect(rowItems[4]).toHaveTextContent('Subscriber session was terminated');

    expect(rowItems[5]).toHaveTextContent(
      new Date(mockEvents[3].timestamp).toLocaleString(),
    );
    expect(rowItems[5]).toHaveTextContent('mme');
    expect(rowItems[5]).toHaveTextContent('detach_success');
    expect(rowItems[5]).toHaveTextContent('001011234560000');
    expect(rowItems[5]).toHaveTextContent('UE detaches successfully');
  });
});
