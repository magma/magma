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

import FEGGatewayContext from '../../../components/context/FEGGatewayContext';
import FEGGatewayDetailConfig from '../FEGGatewayDetailConfig';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
// $FlowFixMe migrated to typescript
import {SetGatewayState} from '../../../state/feg/EquipmentState';
import {fireEvent, render, wait} from '@testing-library/react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../../app/hooks/useSnackbar';
import {useState} from 'react';
import type {
  csfb,
  federation_gateway,
  gx,
  gy,
  s6a,
  swx,
} from '../../../../generated/MagmaAPIBindings';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../../app/hooks/useSnackbar');
const mockGx: gx = {
  server: {
    address: '174.16.1.14:3868',
    dest_host: 'magma.magma.com',
    dest_realm: 'magma.com',
    host: '',
    realm: '',
    local_address: '',
    product_name: 'magma',
    protocol: 'tcp',
    disable_dest_host: false,
  },
  virtual_apn_rules: [],
};

const mockGy: gy = {
  server: {
    address: '174.18.1.0:3868',
    dest_host: '',
    dest_realm: '',
    host: 'localhost',
    realm: 'test',
    local_address: '',
    product_name: '',
    protocol: 'tcp',
    disable_dest_host: false,
  },
  init_method: 2,
  virtual_apn_rules: [],
};

const mockSwx: swx = {
  server: {
    address: '174.18.1.0:3869',
    dest_host: '',
    dest_realm: '',
    host: '',
    realm: '',
    local_address: ':3809',
    product_name: '',
    protocol: 'tcp',
    disable_dest_host: false,
  },
};

const mockS6a: s6a = {
  server: {
    address: '174.18.1.0:2000',
    dest_host: '',
    dest_realm: '',
    host: '',
    realm: '',
    local_address: '',
    product_name: '',
    protocol: 'tcp',
    disable_dest_host: false,
  },
  plmn_ids: [],
};

const mockCsfb: csfb = {
  client: {
    server_address: '174.18.1.0:2200',
    local_address: ':3440',
  },
};

const mockGw0: federation_gateway = {
  id: 'test_feg_gw0',
  name: 'test_gateway',
  description: 'hello I am a federated gateway',
  tier: 'default',
  device: {
    key: {key: '', key_type: 'SOFTWARE_ECDSA_SHA256'},
    hardware_id: 'c9439d30-61ef-46c7-93f2-e01fc131255d',
  },
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
    tier: 'tier2',
  },
  federation: {
    aaa_server: {},
    eap_aka: {},
    gx: mockGx,
    gy: mockGy,
    health: {},
    hss: {},
    s6a: mockS6a,
    s8: {
      apn_operator_suffix: '',
      local_address: '',
      pgw_address: '',
    },
    served_network_ids: [],
    swx: mockSwx,
    csfb: mockCsfb,
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

const fegGateways = {
  [mockGw0.id]: mockGw0,
};

const fegGatewaysHealth = {
  [mockGw0.id]: {status: 'HEALTHY'},
};

describe('<FEGGatewayDetailConfig />', () => {
  beforeEach(() => {
    // Mocking value because it is called by FEGGatewayDialogue / Edit Gateway Page
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockResolvedValue([]);
  });

  const Wrapper = () => (
    <MemoryRouter
      initialEntries={[
        '/nms/mynetwork/equipment/overview/gateway/test_feg_gw0/config',
      ]}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGGatewayContext.Provider
            value={{
              state: fegGateways,
              setState: async _ => {},
              health: fegGatewaysHealth,
              activeFegGatewayId: mockGw0.id,
            }}>
            <Routes>
              <Route
                path="/nms/:networkId/equipment/overview/gateway/:gatewayId/config"
                element={<FEGGatewayDetailConfig />}
              />
            </Routes>
          </FEGGatewayContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  const EditWrapper = () => {
    const [fegGateways, setFegGateways] = useState({[mockGw0.id]: mockGw0});
    const [fegGatewaysHealthStatus, setFegGatewaysHealthStatus] = useState(
      fegGatewaysHealth,
    );
    const [activeFegGatewayId, setActiveFegGatewayId] = useState(mockGw0.id);
    const enqueueSnackbar = useEnqueueSnackbar();
    return (
      <MemoryRouter
        initialEntries={[
          '/nms/mynetwork/equipment/overview/gateway/test_feg_gw0/config',
        ]}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <FEGGatewayContext.Provider
              value={{
                state: fegGateways,
                setState: async (key, value?, newState?) => {
                  return SetGatewayState({
                    networkId: 'mynetwork',
                    fegGateways,
                    fegGatewaysHealthStatus,
                    setFegGateways,
                    setFegGatewaysHealthStatus,
                    setActiveFegGatewayId,
                    key,
                    value,
                    newState,
                    enqueueSnackbar,
                  });
                },
                health: fegGatewaysHealthStatus,
                activeFegGatewayId: activeFegGatewayId,
              }}>
              <Routes>
                <Route
                  path="/nms/:networkId/equipment/overview/gateway/:gatewayId/config"
                  element={<FEGGatewayDetailConfig />}
                />
              </Routes>
            </FEGGatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };
  it('renders federation gateway configs correctly', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    // verify gateway info
    expect(getByTestId('Name')).toHaveTextContent('test_gateway');
    expect(getByTestId('Gateway ID')).toHaveTextContent('test_feg_gw0');
    expect(getByTestId('Hardware UUID')).toHaveTextContent(
      'c9439d30-61ef-46c7-93f2-e01fc131255d',
    );
    expect(getByTestId('Version')).toHaveTextContent('null');
    expect(getByTestId('Description')).toHaveTextContent(
      'hello I am a federated gateway',
    );
    // verify gx configurations
    expect(getByTestId('Gx')).toHaveTextContent('174.16.1.14:3868');
    expect(getByTestId('Gx')).toHaveTextContent('magma.magma.com');
    expect(getByTestId('Gx')).toHaveTextContent('magma.com');
    expect(getByTestId('Gx')).toHaveTextContent('magma');
    // verify gy configurations
    expect(getByTestId('Gy')).toHaveTextContent('74.18.1.0:3868');
    expect(getByTestId('Gy')).toHaveTextContent('localhost');
    expect(getByTestId('Gy')).toHaveTextContent('test');
    // verify swx configurations
    expect(getByTestId('SWx')).toHaveTextContent('174.18.1.0:3869');
    expect(getByTestId('SWx')).toHaveTextContent(':3809');
    // verify s6a configurations
    expect(getByTestId('S6a')).toHaveTextContent('174.18.1.0:2000');
    expect(getByTestId('S6a')).toHaveTextContent('tcp');
    // verify csfb configurations
    expect(getByTestId('CSFB')).toHaveTextContent('174.18.1.0:2200');
    expect(getByTestId('CSFB')).toHaveTextContent('3440');
  });
  it('verify gx edit is working', async () => {
    const {getByTestId, getByText} = render(<EditWrapper />);
    await wait();
    fireEvent.click(getByTestId('gxEditButton'));
    await wait();
    const address = getByTestId('address');
    const destHost = getByTestId('destinationHost');
    const destRealm = getByTestId('destRealm');
    fireEvent.change(address, {target: {value: '194.16.1.14:3868'}});
    fireEvent.change(destHost, {target: {value: 'abc.xyz.com'}});
    fireEvent.change(destRealm, {target: {value: 'xyz.com'}});
    fireEvent.click(getByText('Save'));
    await wait();
    // verify that address, dest_host, and dest_realm were edited
    expect(
      MagmaAPIBindings.putFegByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
      gateway: {
        ...mockGw0,
        federation: {
          ...mockGw0.federation,
          gx: {
            ...mockGw0.federation.gx,
            server: {
              ...mockGw0.federation.gx.server,
              address: '194.16.1.14:3868',
              dest_host: 'abc.xyz.com',
              dest_realm: 'xyz.com',
            },
          },
        },
      },
    });
  });
  it('verify gy edit is working', async () => {
    const {getByTestId, getByText} = render(<EditWrapper />);
    await wait();
    fireEvent.click(getByTestId('gyEditButton'));
    await wait();
    const address = getByTestId('address');
    const host = getByTestId('host');
    const realm = getByTestId('realm');
    fireEvent.change(address, {target: {value: '174.18.1.1:4868'}});
    fireEvent.change(host, {target: {value: '222.222.222.222'}});
    fireEvent.change(realm, {target: {value: 'test_realm'}});
    fireEvent.click(getByText('Save'));
    await wait();
    // verify that address, host, and realm were edited
    expect(
      MagmaAPIBindings.putFegByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
      gateway: {
        ...mockGw0,
        federation: {
          ...mockGw0.federation,
          gy: {
            ...mockGw0.federation.gy,
            server: {
              ...mockGw0.federation.gy.server,
              address: '174.18.1.1:4868',
              host: '222.222.222.222',
              realm: 'test_realm',
            },
          },
        },
      },
    });
  });
  it('verify swx edit is working', async () => {
    const {getByTestId, getByText} = render(<EditWrapper />);
    await wait();
    fireEvent.click(getByTestId('swxEditButton'));
    await wait();
    const address = getByTestId('address');
    const localAddress = getByTestId('localAddress');
    fireEvent.change(address, {target: {value: '174.58.1.0:3869'}});
    fireEvent.change(localAddress, {target: {value: ':4444'}});
    fireEvent.click(getByText('Save'));
    await wait();
    // verify that address and local_address were edited
    expect(
      MagmaAPIBindings.putFegByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
      gateway: {
        ...mockGw0,
        federation: {
          ...mockGw0.federation,
          swx: {
            ...mockGw0.federation.swx,
            server: {
              ...mockGw0.federation.swx.server,
              address: '174.58.1.0:3869',
              local_address: ':4444',
            },
          },
        },
      },
    });
  });
  it('verify s6a edit is working', async () => {
    const {getByTestId, getByText} = render(<EditWrapper />);
    await wait();
    fireEvent.click(getByTestId('s6aEditButton'));
    await wait();
    const protocol = getByTestId('protocol');
    fireEvent.change(protocol, {target: {value: 'sctp'}});
    fireEvent.click(getByText('Save'));
    await wait();
    // verify that protocol was edited
    expect(
      MagmaAPIBindings.putFegByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
      gateway: {
        ...mockGw0,
        federation: {
          ...mockGw0.federation,
          s6a: {
            ...mockGw0.federation.s6a,
            server: {
              ...mockGw0.federation.s6a.server,
              protocol: 'sctp',
            },
          },
        },
      },
    });
  });
  it('verify csfb edit is working', async () => {
    const {getByTestId, getByText} = render(<EditWrapper />);
    await wait();
    fireEvent.click(getByTestId('csfbEditButton'));
    await wait();
    const serverAddress = getByTestId('serverAddress');
    const localAddress = getByTestId('localAddress');
    fireEvent.change(serverAddress, {target: {value: '175.18.1.0:2200'}});
    fireEvent.change(localAddress, {target: {value: ':4400'}});
    fireEvent.click(getByText('Save'));
    await wait();
    // verify that server_address and local_address were edited
    expect(
      MagmaAPIBindings.putFegByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
      gateway: {
        ...mockGw0,
        federation: {
          ...mockGw0.federation,
          csfb: {
            ...mockGw0.federation.csfb,
            client: {
              server_address: '175.18.1.0:2200',
              local_address: ':4400',
            },
          },
        },
      },
    });
  });
});
