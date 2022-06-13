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

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import * as customHistogram from '../../../components/CustomMetrics';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import GatewayLogs from '../GatewayLogs';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MomentUtils from '@date-io/moment';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiPickersUtilsProvider} from '@material-ui/pickers';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../../app/hooks/useSnackbar');
jest.spyOn(customHistogram, 'default').mockImplementation(() => <></>);
const LogTableWrapper = () => (
  <MemoryRouter
    initialEntries={['/nms/mynetwork/gateway/mygateway/logs']}
    initialIndex={0}>
    <MuiPickersUtilsProvider utils={MomentUtils}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <Routes>
            <Route
              path="/nms/:networkId/gateway/:gatewayId/logs"
              element={<GatewayLogs />}
            />
          </Routes>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MuiPickersUtilsProvider>
  </MemoryRouter>
);

describe('<GatewayLogs />', () => {
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
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdLogsCount.mockResolvedValue(
      mockLogCount,
    );

    MagmaAPIBindings.getNetworksByNetworkIdLogsSearch.mockResolvedValue(
      mockLogs,
    );
  });

  it('verify gateway logs rendering', async () => {
    const {getAllByRole} = render(<LogTableWrapper />);
    await wait();
    const rowItems = getAllByRole('row');

    // can get called multiple times from the histogram component
    // as well
    expect(MagmaAPIBindings.getNetworksByNetworkIdLogsCount).toHaveBeenCalled();

    expect(
      MagmaAPIBindings.getNetworksByNetworkIdLogsSearch,
    ).toHaveBeenCalledTimes(1);

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
