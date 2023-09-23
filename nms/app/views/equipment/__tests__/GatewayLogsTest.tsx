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

import * as customHistogram from '../../../components/CustomMetrics';
import GatewayLogs from '../GatewayLogs';
import MagmaAPI from '../../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {AdapterDateFns} from '@mui/x-date-pickers/AdapterDateFns';
import {LocalizationProvider} from '@mui/x-date-pickers';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {mockAPI} from '../../../util/TestUtils';
import {render, waitFor} from '@testing-library/react';

jest.mock('../../../../app/hooks/useSnackbar');
jest.spyOn(customHistogram, 'default').mockImplementation(() => <></>);

const LogTableWrapper = () => (
  <MemoryRouter
    initialEntries={['/nms/mynetwork/gateway/mygateway/logs']}
    initialIndex={0}>
    <LocalizationProvider dateAdapter={AdapterDateFns}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <Routes>
            <Route
              path="/nms/:networkId/gateway/:gatewayId/logs"
              element={<GatewayLogs />}
            />
          </Routes>
        </ThemeProvider>
      </StyledEngineProvider>
    </LocalizationProvider>
  </MemoryRouter>
);

// This test is being skipped. Test failures needs to be investigated
// and fixed, see https://github.com/magma/magma/issues/15122 for details.
describe.skip('<GatewayLogs />', () => {
  const mockLogCount = 100;
  const mockLogs = [
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'nd6gqXIB736xuPLmCwoc',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'control_proxy',
        message: 'Message1',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'P6yfqXIBGyZqNMEmqkW-',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'magmad',
        message: 'Info:Message2',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'P6yfqXIBGyZqNMEmqkW-',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'magmad',
        message: 'Info:Message2',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'P6yfqXIBGyZqNMEmqkW-',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'magmad',
        message: 'Info:Message2',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'P6yfqXIBGyZqNMEmqkW-',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'magmad',
        message: 'Info:Message2',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
    {
      _index: 'magma-2020.06.12',
      _type: '_doc',
      _id: 'P6yfqXIBGyZqNMEmqkW-',
      _source: {
        time: 'Jun 12 17:42:08',
        ident: 'magmad',
        message: 'Info:Message2',
        '@timestamp': '2020-06-12T17:42:08.000000000+00:00',
        tag: 'gateway.syslog',
      },
    },
  ];
  beforeEach(() => {
    mockAPI(MagmaAPI.logs, 'networksNetworkIdLogsCountGet', mockLogCount);
    mockAPI(MagmaAPI.logs, 'networksNetworkIdLogsSearchGet', mockLogs);
  });

  it('verify gateway logs rendering', async () => {
    const {findAllByRole} = render(<LogTableWrapper />);

    await waitFor(() => {
      // can get called multiple times from the histogram component
      // as well
      expect(MagmaAPI.logs.networksNetworkIdLogsCountGet).toHaveBeenCalled();

      expect(
        MagmaAPI.logs.networksNetworkIdLogsSearchGet,
      ).toHaveBeenCalledTimes(1);
    });

    const rowItems = await findAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Date');
    expect(rowItems[0]).toHaveTextContent('Service');
    expect(rowItems[0]).toHaveTextContent('Type');
    expect(rowItems[0]).toHaveTextContent('Output');

    expect(rowItems[2]).toHaveTextContent(
      new Date('2020-06-12T17:42:08.000000000+00:00').toLocaleString(),
    );
    expect(rowItems[2]).toHaveTextContent('control_proxy');
    expect(rowItems[2]).toHaveTextContent('debug');
    expect(rowItems[2]).toHaveTextContent('Message1');

    expect(rowItems[2]).toHaveTextContent(
      new Date('2020-06-12T17:42:08.000000000+00:00').toLocaleString(),
    );
    expect(rowItems[3]).toHaveTextContent('magmad');
    expect(rowItems[3]).toHaveTextContent('info');
    expect(rowItems[3]).toHaveTextContent('Info:Message2');
    expect(rowItems[3]).toHaveTextContent(
      new Date('2020-06-12T17:42:08.000000000+00:00').toLocaleString(),
    );
  });
});
