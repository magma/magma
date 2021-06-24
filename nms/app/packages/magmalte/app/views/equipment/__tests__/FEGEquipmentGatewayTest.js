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
import FEGEquipmentGateway from '../FEGEquipmentGateway';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';
import {FEGGatewayContextProvider} from '../../../components/feg/FEGContext';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, render, wait} from '@testing-library/react';
import type {
  federation_gateway,
  federation_gateway_health_status,
  federation_network_cluster_status,
  promql_return_object,
} from '@fbcnms/magma-api';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);

const mockGw0: federation_gateway = {
  id: 'test_feg_gw0',
  name: 'test_gateway',
  description: 'hello I am a federated gateway',
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
  federation: {
    aaa_server: {},
    eap_aka: {
      plmn_ids: [],
    },
    gx: {
      server: {
        protocol: 'tcp',
      },
      servers: [],
      virtual_apn_rules: [],
    },
    gy: {
      server: {
        protocol: 'tcp',
      },
      servers: [],
      virtual_apn_rules: [],
    },
    health: {
      health_services: [],
    },
    hss: {},
    s6a: {
      plmn_ids: [],
      server: {
        protocol: 'tcp',
      },
    },
    served_network_ids: [],
    swx: {
      hlr_plmn_ids: [],
      server: {
        protocol: 'tcp',
      },
      servers: [],
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

const mockGw1: federation_gateway = {
  ...mockGw0,
  id: 'test_gw1',
  name: 'test_gateway1',
};

const fegGateways = {
  [mockGw0.id]: mockGw0,
  [mockGw1.id]: mockGw1,
};

const mockHealthyGatewayStatus: federation_gateway_health_status = {
  description: '',
  status: 'HEALTHY',
};

const mockUnhealthyGatewayStatus: federation_gateway_health_status = {
  description: '',
  status: 'UNHEALTHY',
};

const mockClusterStatus: federation_network_cluster_status = {
  active_gateway: mockGw0.id,
};

describe('<FEGEquipmentGateway />', () => {
  beforeEach(() => {
    // gateway context gets list of federation gateways
    MagmaAPIBindings.getFegByNetworkIdGateways.mockResolvedValue(fegGateways);
    // gateway context gets health status of the gateways
    MagmaAPIBindings.getFegByNetworkIdGatewaysByGatewayIdHealthStatus.mockImplementation(
      req => {
        if (req.gatewayId == mockGw0.id) {
          // only gateway 0 is healthy
          return mockHealthyGatewayStatus;
        }
        return mockUnhealthyGatewayStatus;
      },
    );
    // gateway context gets the active gateway id
    MagmaAPIBindings.getFegByNetworkIdClusterStatus.mockResolvedValue(
      mockClusterStatus,
    );
    // called by gateway checkin chart
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQueryRange.mockResolvedValue(
      mockCheckinMetric,
    );
    // called when getting max, min and average latency
    MagmaAPIBindings.getNetworksByNetworkIdPrometheusQuery.mockResolvedValue(
      mockKPIMetric,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateway']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGGatewayContextProvider networkId="mynetwork" networkType="FEG">
            <Route
              path="/nms/:networkId/gateway/"
              render={props => <FEGEquipmentGateway {...props} />}
            />
          </FEGGatewayContextProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders federation gateway KPIs correctly', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    // verify KPI metrics
    expect(getByTestId('Max Latency')).toHaveTextContent('8');
    expect(getByTestId('Min Latency')).toHaveTextContent('6');
    expect(getByTestId('Avg Latency')).toHaveTextContent('7');
    expect(getByTestId('Federation Gateway Count')).toHaveTextContent('2');
    expect(getByTestId('Healthy Federation Gateway Count')).toHaveTextContent(
      '1',
    );
    expect(getByTestId('% Healthy Gateways')).toHaveTextContent('50');
  });
});
