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

import FEGGatewayConnectionStatus from '../FEGGatewayConnectionStatus';
import FEGGatewayContext from '../../../context/FEGGatewayContext';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {FederationGatewayHealthStatus} from '../../../components/GatewayUtils';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {render} from '@testing-library/react';
import type {FederationGateway} from '../../../../generated';

const mockGw0: FederationGateway = {
  id: 'test_feg_gw0',
  name: 'test_gateway',
  description: 'hello I am a federated gateway',
  tier: 'default',
  device: {
    key: {key: '', key_type: 'SOFTWARE_ECDSA_SHA256'},
    hardware_id: 'c9439d30-61ef-46c7-93f2-e01fc131244d',
  },
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
    dynamic_services: ['monitord', 'eventd'],
  },
  federation: {
    aaa_server: {},
    eap_aka: {
      plmn_ids: [],
    },
    gx: {},
    gy: {},
    health: {
      health_services: [],
    },
    hss: {},
    s6a: {},
    served_network_ids: [],
    swx: {
      hlr_plmn_ids: [],
      server: {
        protocol: 'tcp',
      },
      servers: [],
    },
    csfb: {},
  },
};

const fegGateways = {
  [mockGw0.id]: mockGw0,
};

const fegGatewaysHealth = {
  [mockGw0.id]: {
    status: 'HEALTHY',
    service_status: {
      SESSION_PROXY: {health_status: 'UNHEALTHY', service_state: 'UNAVAILABLE'},
      SWX_PROXY: {health_status: 'UNHEALTHY', service_state: 'AVAILABLE'},
    },
  },
} as Record<string, FederationGatewayHealthStatus>;

describe('<FEGGatewayConnectionStatus />', () => {
  const Wrapper = () => (
    <MemoryRouter
      initialEntries={[
        '/nms/mynetwork/equipment/overview/gateway/test_feg_gw0/overview',
      ]}
      initialIndex={0}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <FEGGatewayContext.Provider
            value={{
              state: fegGateways,
              setState: async () => {},
              updateGateway: async () => {},
              refetch: () => {},
              health: fegGatewaysHealth,
              activeFegGatewayId: mockGw0.id,
            }}>
            <Routes>
              <Route
                path="/nms/:networkId/equipment/overview/gateway/:gatewayId/overview"
                element={<FEGGatewayConnectionStatus />}
              />
            </Routes>
          </FEGGatewayContext.Provider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );

  it('renders federation gateway connection status correctly', async () => {
    const {findByTestId, getByTestId} = render(<Wrapper />);
    // verify gateway connection status
    expect(await findByTestId('Gx/Gy Watchdog')).toHaveTextContent('N/A');
    expect(getByTestId('SWx Watchdog')).toHaveTextContent('Down');
    expect(getByTestId('S6a Watchdog')).toHaveTextContent('Not Enabled');
  });
});
