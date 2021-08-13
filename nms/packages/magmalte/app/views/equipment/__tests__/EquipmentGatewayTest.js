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
import 'jest-dom/extend-expect';
import Gateway from '../EquipmentGateway';
import GatewayContext from '../../../components/context/GatewayContext';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';
import type {lte_gateway, promql_return_object} from '@fbcnms/magma-api';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);

const mockCheckinMetric: promql_return_object = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        values: [['1588898968.042', '6']],
      },
    ],
  },
};

const mockGw0: lte_gateway = {
  id: 'test_gw0',
  name: 'test_gateway0',
  description: 'test_gateway0',
  tier: 'default',
  device: {
    key: {key: '', key_type: 'SOFTWARE_ECDSA_SHA256'},
    hardware_id: '',
  },
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
    tier: 'tier2',
  },
  connected_enodeb_serials: [],
  cellular: {
    epc: {
      ip_block: '192.168.0.1/24',
      nat_enabled: false,
      sgi_management_iface_static_ip: '1.1.1.1/24',
      sgi_management_iface_vlan: '100',
    },
    ran: {
      pci: 620,
      transmit_enabled: true,
    },
  },
  status: {
    checkin_time: 0,
    meta: {
      gps_latitude: '0',
      gps_longitude: '0',
      gps_connected: '0',
      enodeb_connected: '0',
      mme_connected: '0',
    },
  },
};

const mockKPIMetric: promql_return_object = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        value: ['1588898968.042', '6'],
      },
      {
        metric: {},
        value: ['1588898968.042', '8'],
      },
    ],
  },
};

const currTime = Date.now();

describe('<Gateway />', () => {
  beforeEach(() => {
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockCheckinMetric,
    );
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQuery.mockResolvedValue(
      mockKPIMetric,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const mockGw1 = Object.assign({}, mockGw0);
  const mockGw2 = Object.assign({}, mockGw0);
  mockGw1.id = 'test_gw1';
  mockGw1.name = 'test_gateway1';
  mockGw1.connected_enodeb_serials = ['xxx', 'yyy'];

  mockGw2.id = 'test_gw2';
  mockGw2.name = 'test_gateway2';
  mockGw2.connected_enodeb_serials = ['xxx'];
  mockGw2.status = {
    checkin_time: currTime,
  };
  const lteGateways = {
    test1: mockGw0,
    test2: mockGw1,
    test3: mockGw2,
  };

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateway']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <GatewayContext.Provider
            value={{
              state: lteGateways,
              setState: async _ => {},
              updateGateway: async _ => {},
            }}>
            <Route
              path="/nms/:networkId/gateway/"
              render={props => <Gateway {...props} lteGateways={lteGateways} />}
            />
          </GatewayContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    const {getByTestId, getAllByRole, getAllByTitle} = render(<Wrapper />);
    await wait();

    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange,
    ).toHaveBeenCalledTimes(1);

    expect(
      MagmaAPIBindings.getNetworksByNetworkIdPrometheusQuery,
    ).toHaveBeenCalledTimes(3);

    // verify KPI metrics
    expect(getByTestId('Max Latency')).toHaveTextContent('8');
    expect(getByTestId('Min Latency')).toHaveTextContent('6');
    expect(getByTestId('Avg Latency')).toHaveTextContent('7');
    expect(getByTestId('% Healthy Gateways')).toHaveTextContent('33.33');

    const rowItems = await getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('ID');
    expect(rowItems[0]).toHaveTextContent('enodeBs');
    expect(rowItems[0]).toHaveTextContent('Subscribers');
    expect(rowItems[0]).toHaveTextContent('Health');
    expect(rowItems[0]).toHaveTextContent('Check In Time');

    expect(rowItems[1]).toHaveTextContent('test_gw0');
    expect(rowItems[1]).toHaveTextContent('test_gateway0');
    expect(rowItems[1]).toHaveTextContent('0');
    expect(rowItems[1]).toHaveTextContent('Bad');
    expect(rowItems[1]).toHaveTextContent('-');

    expect(rowItems[2]).toHaveTextContent('test_gw1');
    expect(rowItems[2]).toHaveTextContent('test_gateway1');
    expect(rowItems[2]).toHaveTextContent('2');
    expect(rowItems[2]).toHaveTextContent('Bad');
    expect(rowItems[2]).toHaveTextContent('-');

    expect(rowItems[3]).toHaveTextContent('test_gw2');
    expect(rowItems[3]).toHaveTextContent('test_gateway2');
    expect(rowItems[3]).toHaveTextContent('1');
    expect(rowItems[3]).toHaveTextContent('Good');
    expect(rowItems[3]).toHaveTextContent(
      new Date(currTime).toLocaleDateString(),
    );

    // click the actions button for gateway 0
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
  });
});
