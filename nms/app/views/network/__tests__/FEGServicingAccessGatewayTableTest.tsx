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

import FEGNetworkContext from '../../../components/context/FEGNetworkContext';
import FEGServicingAccessGatewaysTable from '../FEGServicingAccessGatewayTable';
import MagmaAPI from '../../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';
import {AxiosResponse} from 'axios';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render, wait} from '@testing-library/react';
import type {
  FegLteNetwork,
  FegNetwork,
  LteGateway,
} from '../../../../generated-ts';

const mockFegLteNetworks: Array<string> = [
  'test_network1',
  'test_network2',
  'test_network3',
];

const mockFegLteNetwork1: FegLteNetwork = {
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

const mockGw1: LteGateway = {
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

//jest.mock('axios');
jest.mock('../../../../app/hooks/useSnackbar');

describe('<ServicingAccessGatewaysInfo />', () => {
  const testNetwork: FegNetwork = {
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
    ...mockFegLteNetwork1,
    federation: {FegNetwork_id: ''},
    id: 'test_network2',
    name: 'test network2',
  };
  const mockFegLteNetwork3 = {
    ...mockFegLteNetwork1,
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
    jest
      .spyOn(MagmaAPI.federatedLTENetworks, 'fegLteGet')
      .mockResolvedValue({data: mockFegLteNetworks} as AxiosResponse);
    jest
      .spyOn(MagmaAPI.federatedLTENetworks, 'fegLteNetworkIdGet')
      .mockResolvedValueOnce({data: mockFegLteNetwork1} as AxiosResponse)
      .mockResolvedValueOnce({data: mockFegLteNetwork2} as AxiosResponse)
      .mockResolvedValue({data: mockFegLteNetwork3} as AxiosResponse);
    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdGatewayPoolsGet')
      .mockResolvedValueOnce({
        data: {
          [mockGw1.id]: mockGw1,
          [mockGw2.id]: mockGw2,
        },
      } as AxiosResponse)
      .mockResolvedValue({
        data: {[mockGw3.id]: mockGw3},
      } as AxiosResponse);
  });

  const Wrapper = () => {
    const networkCtx = {
      state: {...testNetwork},
      updateNetworks: async () => {},
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
    expect(MagmaAPI.federatedLTENetworks.fegLteGet).toHaveBeenCalledTimes(1);
    //get info about each feg_lte network
    expect(
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet,
    ).toHaveBeenCalledTimes(3);
    expect(
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet,
    ).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork1.id,
    });
    expect(
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet,
    ).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork2.id,
    });
    expect(
      MagmaAPI.federatedLTENetworks.fegLteNetworkIdGet,
    ).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork3.id,
    });
    //only 2 of the 3 feg_lte networks are serviced by current network
    expect(
      MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGet,
    ).toHaveBeenCalledTimes(2);
    expect(
      MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGet,
    ).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork1.id,
    });
    expect(
      MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGet,
    ).toHaveBeenCalledWith({
      networkId: mockFegLteNetwork3.id,
    });
    const rowItems = getAllByRole('row');
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
