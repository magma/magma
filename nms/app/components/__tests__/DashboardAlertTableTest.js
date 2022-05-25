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
import DashboardAlertTable from '../DashboardAlertTable';
import MagmaAPIBindings from '../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, wait} from '@testing-library/react';
import type {
  gettable_alert,
  prom_firing_alert,
} from '../../../generated/MagmaAPIBindings';

const tbl_alert: gettable_alert = {
  name: 'null_receiver',
};
const mockAlertSt: Array<prom_firing_alert> = [
  {
    annotations: {
      description: 'TestMetric1 Description',
      summary: 'TestMetric1 Minor Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert1',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'host',
      networkID: 'test',
      severity: 'critical',
    },
  },
  {
    annotations: {
      description: 'TestMetric2 Description',
      summary: 'TestMetric2 Major Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert2',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'host',
      networkID: 'test',
      severity: 'major',
    },
  },
  {
    annotations: {
      description: 'TestMetric3 Description',
      summary: 'TestMetric3 Critical Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert3',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'host',
      networkID: 'test',
      severity: 'minor',
    },
  },
  {
    annotations: {
      description: 'TestMetric4 Description',
      summary: 'TestMetric1 Other Alert',
    },
    endsAt: '2020-05-14T18:55:25.844Z',
    fingerprint: '0de443c4dd7af53e',
    receivers: tbl_alert,
    startsAt: '2020-05-14T18:27:25.844Z',
    status: {inhibitedBy: [], silencedBy: [], state: 'active'},
    updatedAt: '2020-05-14T18:52:31.971Z',
    generatorURL:
      'http://1521e855c607:9090/graph?g0.expr=TestMetric9+%3E+1\u0026g0.tab=1',
    labels: {
      alertname: 'TestAlert4',
      instance: '192.168.0.124:2112',
      job: 'myapp',
      monitor: 'host',
      networkID: 'test',
      severity: 'normal',
    },
  },
];

jest.mock('axios');
jest.mock('../../../generated/MagmaAPIBindings');
jest.mock('../../../app/hooks/useSnackbar');

describe('<DashboardAlertTable />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getNetworksByNetworkIdAlerts.mockResolvedValue(
      mockAlertSt,
    );
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <Routes>
            <Route path="/nms/:networkId/*" element={<DashboardAlertTable />} />
          </Routes>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    const {getByText, getAllByRole} = render(<Wrapper />);
    await wait();
    expect(MagmaAPIBindings.getNetworksByNetworkIdAlerts).toHaveBeenCalledTimes(
      1,
    );

    // get all rows
    const rowItems = await getAllByRole('row');
    // check if the default is critical alert sections
    expect(rowItems[1]).toHaveTextContent('TestAlert1');
    fireEvent.click(getByText('Critical(1)'));
    expect(rowItems[1]).toHaveTextContent('TestAlert1');

    fireEvent.click(getByText('Major(1)'));
    const rowItems2 = await getAllByRole('row');
    expect(rowItems2[1]).toHaveTextContent('TestAlert2');

    fireEvent.click(getByText('Minor(1)'));
    const rowItems3 = await getAllByRole('row');
    expect(rowItems3[1]).toHaveTextContent('TestAlert3');

    fireEvent.click(getByText('Other(1)'));
    const rowItems4 = await getAllByRole('row');
    expect(rowItems4[1]).toHaveTextContent('TestAlert4');

    expect(getByText('Alerts (4)')).toBeInTheDocument();
  });
});
