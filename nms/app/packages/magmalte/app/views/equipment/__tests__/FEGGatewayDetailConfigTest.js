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
import FEGGatewayContext from '../../../components/context/FEGGatewayContext';
import FEGGatewayDetailConfig from '../FEGGatewayDetailConfig';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, render, wait} from '@testing-library/react';
import type {
  csfb,
  federation_gateway,
  gx,
  gy,
  s6a,
  swx,
} from '@fbcnms/magma-api';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);

const mockGx: gx = {
  server: {
    address: '174.16.1.14:3868',
    dest_host: 'magma.magma.com',
    dest_realm: 'magma.com',
    product_name: 'magma',
  },
};

const mockGy: gy = {
  server: {
    address: '174.18.1.0:3868',
    host: 'localhost',
    realm: 'test',
  },
};

const mockSwx: swx = {
  server: {
    address: '174.18.1.0:3869',
    local_address: ':3809',
  },
};

const mockS6a: s6a = {
  server: {
    address: '174.18.1.0:2000',
    protocol: 'tcp',
  },
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
    eap_aka: {
      plmn_ids: [],
    },
    gx: mockGx,
    gy: mockGy,
    health: {
      health_services: [],
    },
    hss: {},
    s6a: mockS6a,
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
            <Route
              path="/nms/:networkId/equipment/overview/gateway/:gatewayId/config"
              render={props => <FEGGatewayDetailConfig {...props} />}
            />
          </FEGGatewayContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

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
});
