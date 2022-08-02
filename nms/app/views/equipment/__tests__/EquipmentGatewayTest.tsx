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
import Gateway from '../EquipmentGateway';
import GatewayContext from '../../../context/GatewayContext';
import MagmaAPI from '../../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {mockAPI} from '../../../util/TestUtils';
import {render} from '../../../util/TestingLibrary';
import {waitFor} from '@testing-library/react';
import type {LteGateway, PromqlReturnObject} from '../../../../generated';

jest.mock('../../../hooks/useSnackbar');

const mockCheckinMetric: PromqlReturnObject = {
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

const mockGw0: LteGateway = {
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
  checked_in_recently: false,
};

const mockKPIMetric: PromqlReturnObject = {
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
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryRangeGet',
      mockCheckinMetric,
    );

    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryGet',
      mockKPIMetric,
    );
  });

  const mockGw1: LteGateway = {
    ...mockGw0,
    id: 'test_gw1',
    name: 'test_gateway1',
    connected_enodeb_serials: ['xxx', 'yyy'],
  };
  const mockGw2: LteGateway = {
    ...mockGw0,
    id: 'test_gw2',
    name: 'test_gateway2',
    checked_in_recently: true,
    connected_enodeb_serials: ['xxx'],
    status: {...mockGw0.status, checkin_time: currTime},
  };
  const lteGateways = {
    test1: mockGw0,
    test2: mockGw1,
    test3: mockGw2,
  };

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateway']} initialIndex={0}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <GatewayContext.Provider
            value={{
              state: lteGateways,
              setState: async () => {},
              updateGateway: async () => {},
              refetch: () => {},
            }}>
            <Routes>
              <Route path="/nms/:networkId/gateway/" element={<Gateway />} />
            </Routes>
          </GatewayContext.Provider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    const {
      findByTestId,
      getByTestId,
      getAllByRole,
      openActionsTableMenu,
    } = render(<Wrapper />);

    await waitFor(() =>
      expect(
        MagmaAPI.metrics.networksNetworkIdPrometheusQueryRangeGet,
      ).toHaveBeenCalledTimes(1),
    );
    expect(
      MagmaAPI.metrics.networksNetworkIdPrometheusQueryGet,
    ).toHaveBeenCalledTimes(3);

    // verify KPI metrics
    expect(getByTestId('Max Latency')).toHaveTextContent('8');
    expect(getByTestId('Min Latency')).toHaveTextContent('6');
    expect(getByTestId('Avg Latency')).toHaveTextContent('7');
    expect(getByTestId('% Healthy Gateways')).toHaveTextContent('33.33');

    const rowItems = getAllByRole('row');

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
    await openActionsTableMenu(0);
    expect(await findByTestId('actions-menu')).toBeVisible();
  });
});
