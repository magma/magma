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

import * as hooks from '../../../components/context/RefreshContext';
import FEGGatewayContext from '../../../components/context/FEGGatewayContext';
import FEGGatewayDetailStatus from '../FEGGatewayDetailStatus';
import MagmaAPI from '../../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {FederationGatewayHealthStatus} from '../../../components/GatewayUtils';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {mockAPI} from '../../../util/TestUtils';
import {render, wait} from '@testing-library/react';
import type {
  FederationGateway,
  PromqlReturnObject,
} from '../../../../generated-ts';

jest.mock('../../../../app/hooks/useSnackbar');

const mockCheckinTime = new Date();

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
  status: {
    checkin_time: mockCheckinTime.getTime(),
    meta: {
      gps_latitude: '0',
      gps_longitude: '0',
      gps_connected: '0',
      enodeb_connected: '0',
      mme_connected: '0',
    },
    platform_info: {
      packages: [{version: '1.2'}],
    },
  },
};

const mockGw1: FederationGateway = {
  ...mockGw0,
  id: 'test_gw1',
  name: 'test_gateway1',
};

const fegGateways = {
  [mockGw0.id]: mockGw0,
  [mockGw1.id]: mockGw1,
};

const fegGatewaysHealth = {
  [mockGw0.id]: {status: 'HEALTHY'},
  [mockGw1.id]: {status: 'UNHEALTHY'},
} as Record<string, FederationGatewayHealthStatus>;

const mockCPUUsage: PromqlReturnObject = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        values: [
          ['1625239404', '3'],
          ['1625239419', '0'],
        ],
      },
    ],
  },
};

describe('<FEGGatewayDetailStatus />', () => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-argument
  jest.spyOn(hooks, 'useRefreshingContext').mockImplementation((() => ({
    fegGateways: fegGateways,
    health: fegGatewaysHealth,
    activeFegGatewayId: mockGw0.id,
  })) as any);

  beforeEach(() => {
    // called when getting the CPU Usage
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryRangeGet',
      mockCPUUsage,
    );
  });

  const Wrapper = () => (
    <MemoryRouter
      initialEntries={[
        '/nms/mynetwork/equipment/overview/gateway/test_feg_gw0/overview',
      ]}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGGatewayContext.Provider
            value={{
              state: fegGateways,
              setState: async () => {},
              health: fegGatewaysHealth,
              activeFegGatewayId: mockGw0.id,
            }}>
            <Routes>
              <Route
                path="/nms/:networkId/equipment/overview/gateway/:gatewayId/overview"
                element={<FEGGatewayDetailStatus refresh={true} />}
              />
            </Routes>
          </FEGGatewayContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders federation gateway status correctly', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    // verify gateway status
    expect(getByTestId('Health')).toHaveTextContent('Good');
    expect(getByTestId('Last Check in')).toHaveTextContent(
      mockCheckinTime.toLocaleString(),
    );
    expect(getByTestId('CPU Usage')).toHaveTextContent('3');
    expect(getByTestId('Event Aggregation')).toHaveTextContent('Enabled');
    expect(getByTestId('Log Aggregation')).toHaveTextContent('Disabled');
    expect(getByTestId('CPE Monitoring')).toHaveTextContent('Enabled');
  });
});
