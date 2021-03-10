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
import EnodebContext from '../context/EnodebContext';
import EnodebKPIs from '../EnodebKPIs';
import GatewayContext from '../context/GatewayContext';
import GatewayKPIs from '../GatewayKPIs';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../theme/default';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, render, wait} from '@testing-library/react';
import type {enodeb_state, lte_gateway} from '@fbcnms/magma-api';

afterEach(cleanup);

const mockGwSt: lte_gateway = {
  id: 'test_gw1',
  name: 'test_gateway',
  description: 'hello I am a gateway',
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
      nat_enabled: true,
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

const mockEnbSt: enodeb_state = {
  enodeb_configured: true,
  enodeb_connected: true,
  fsm_state: '',
  gps_connected: true,
  gps_latitude: '',
  gps_longitude: '',
  mme_connected: true,
  opstate_enabled: true,
  ptp_connected: true,
  rf_tx_desired: true,
  rf_tx_on: true,
  ip_address: '192.168.1.254',
};

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

describe('<GatewaysKPIs />', () => {
  const Wrapper = () => {
    const mockUpSt = Object.assign({}, mockGwSt);
    mockUpSt['status'] = {
      checkin_time: Date.now(),
      meta: {
        gps_latitude: '0',
        gps_longitude: '0',
        gps_connected: '0',
        enodeb_connected: '0',
        mme_connected: '0',
      },
    };
    const gatewayCtx = {
      state: {
        test1: mockGwSt,
        test2: mockGwSt,
        test3: mockUpSt,
      },
      setState: async () => {},
      updateGateway: async () => {},
    };

    return (
      <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <GatewayContext.Provider value={gatewayCtx}>
              <Route path="/nms/:networkId" component={GatewayKPIs} />
            </GatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };
  it('renders', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();

    expect(getByTestId('Connected')).toHaveTextContent('1');
    expect(getByTestId('Disconnected')).toHaveTextContent('2');
  });
});

describe('<EnodebKPIs />', () => {
  const mockEnbNotTxSt = Object.assign({}, mockEnbSt);
  mockEnbNotTxSt.rf_tx_on = false;
  const enbInfo = {
    test1: {
      enb: {
        name: 'test1',
        serial: 'test1',
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: true,
        },
      },
      enb_state: mockEnbSt,
    },
    test2: {
      enb: {
        name: 'test2',
        serial: 'test2',
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: true,
        },
      },
      enb_state: mockEnbSt,
    },
    test3: {
      enb: {
        name: 'test3',
        serial: 'test3',
        config: {
          cell_id: 0,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          transmit_enabled: true,
        },
      },
      enb_state: mockEnbNotTxSt,
    },
  };
  const enodebCtx = {
    state: {enbInfo},
    setState: async () => {},
  };

  const Wrapper = () => {
    return (
      <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <EnodebContext.Provider value={enodebCtx}>
              <Route path="/nms/:networkId" component={EnodebKPIs} />
            </EnodebContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };
  it('renders', async () => {
    const {getByTestId} = render(<Wrapper />);
    expect(getByTestId('Total')).toHaveTextContent('3');
    expect(getByTestId('Transmitting')).toHaveTextContent('2');
  });
});
