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

// $FlowFixMe migrated to typescript
import FEGNetworkContext from '../../../components/context/FEGNetworkContext';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import FEGServicingAccessGatewaysTable from '../FEGServicingAccessGatewayTable';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';
import type {
  feg_lte_network,
  feg_network,
  lte_gateway,
} from '../../../../generated/MagmaAPIBindings';

const mockFegLteNetworks: Array<string> = [
  'test_network1',
  'test_network2',
  'test_network3',
];

const mockFegLteNetwork: feg_lte_network = {
  cellular: {
    epc: {
      gx_gy_relay_enabled: false,
      hss_relay_enabled: false,
      lte_auth_amf: '',
      lte_auth_op: '',
      mcc: '',
      mnc: '',
      tac: 1,
    },
    ran: {
      bandwidth_mhz: 20,
    },
  },
  description: 'I am a test federated lte network',
  dns: {
    enable_caching: false,
    local_ttl: 0,
  },
  federation: {
    feg_network_id: 'mynetwork',
  },
  id: 'test_network1',
  name: 'test network1',
};

const mockGw1: lte_gateway = {
  id: 'test_gw1',
  name: 'test gateway1',
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
  checked_in_recently: false,
};

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../../app/hooks/useSnackbar');

describe('<ServicingAccessGatewaysInfo />', () => {
  const testNetwork: feg_network = {
    description: 'Test Network Description',
    federation: {
      aaa_server: {},
      eap_aka: {},
      gx: {},
      gy: {},
      health: {},
      hss: {},
      s6a: {},
      served_network_ids: [],
      swx: {},
    },
    id: 'mynetwork',
    name: 'Test Network',
    dns: {
      enable_caching: false,
      local_ttl: 0,
      records: [],
    },
  };
  const mockFegLteNetwork2 = {
    ...mockFegLteNetwork,
    federation: {feg_network_id: ''},
    id: 'test_network2',
    name: 'test network2',
  };
  const mockFegLteNetwork3 = {
    ...mockFegLteNetwork,
    id: 'test_network3',
    name: 'test network3',
  };
  const mockGw2 = {...mockGw1, id: 'test_gw2', name: 'test gateway2'};
  const mockGw3 = {
    ...mockGw1,
    id: 'test_gw3',
    name: 'test gateway3',
    checked_in_recently: true,
    status: {checkin_time: Date.now()},
  };
  beforeEach(() => {
    MagmaAPIBindings.getFegLte.mockResolvedValue(mockFegLteNetworks);
    MagmaAPIBindings.getFegLteByNetworkId
      .mockReturnValueOnce(mockFegLteNetwork)
      .mockReturnValueOnce(mockFegLteNetwork2)
      .mockResolvedValue(mockFegLteNetwork3);
    MagmaAPIBindings.getLteByNetworkIdGateways
      .mockReturnValueOnce({
        [mockGw1.id]: mockGw1,
        [mockGw2.id]: mockGw2,
      })
      .mockResolvedValue({[mockGw3.id]: mockGw3});
  });

  const Wrapper = () => {
    const networkCtx = {
      state: {
        ...testNetwork,
      },
      updateNetworks: async _ => {},
    };
    return (
      <MemoryRouter
        initialEntries={['/nms/mynetwork/network']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <FEGNetworkContext.Provider value={networkCtx}>
              <Routes>
                <Route
                  path="/nms/:networkId/network/"
                  element={<FEGServicingAccessGatewaysTable />}
                />
              </Routes>
            </FEGNetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };
  it('renders serviced access gateway table correctly', async () => {
    const {getAllByRole} = render(<Wrapper />);
    await wait();
    //first get list of feg_lte networks
    expect(MagmaAPIBindings.getFegLte).toHaveBeenCalledTimes(1);
    //get info about each feg_lte network
    expect(MagmaAPIBindings.getFegLteByNetworkId).toHaveBeenCalledTimes(3);
    expect(MagmaAPIBindings.getFegLteByNetworkId).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork.id,
    });
    expect(MagmaAPIBindings.getFegLteByNetworkId).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork2.id,
    });
    expect(MagmaAPIBindings.getFegLteByNetworkId).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork3.id,
    });
    //only 2 of the 3 feg_lte networks are serviced by current network
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(2);
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork.id,
    });
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork3.id,
    });
    const rowItems = await getAllByRole('row');
    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Access Network');
    expect(rowItems[0]).toHaveTextContent('Access Gateway Id');
    expect(rowItems[0]).toHaveTextContent('Access Gateway Name');
    expect(rowItems[0]).toHaveTextContent('Access Gateway Health');
    //only network 1(with 2 gateways) and network 3 (1 gateway) are serviced
    expect(rowItems[1]).toHaveTextContent('test network1');
    expect(rowItems[1]).toHaveTextContent('test_gw1');
    expect(rowItems[1]).toHaveTextContent('test gateway1');
    expect(rowItems[1]).toHaveTextContent('Bad');
    expect(rowItems[2]).toHaveTextContent('test network1');
    expect(rowItems[2]).toHaveTextContent('test_gw2');
    expect(rowItems[2]).toHaveTextContent('test gateway2');
    expect(rowItems[2]).toHaveTextContent('Bad');
    expect(rowItems[3]).toHaveTextContent('test network3');
    expect(rowItems[3]).toHaveTextContent('test_gw3');
    expect(rowItems[3]).toHaveTextContent('test gateway3');
    expect(rowItems[3]).toHaveTextContent('Good');
  });
});
