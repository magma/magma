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

import FEGEquipmentGateway from '../FEGEquipmentGateway';
import MagmaAPI from '../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import moment from 'moment';
import {AxiosResponse} from 'axios';
import {FEGGatewayContextProvider} from '../../../context/FEGGatewayContext';
import {FederationGatewaysApiFegNetworkIdGatewaysGatewayIdHealthStatusGetRequest} from '../../../../generated';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {mockAPI, mockAPIOnce} from '../../../util/TestUtils';
import {render, wait, waitFor} from '@testing-library/react';
import type {
  Csfb,
  FederationGateway,
  FederationGatewayHealthStatus,
  FederationNetworkClusterStatus,
  Gx,
  Gy,
  PromqlReturnObject,
  S6a,
  Swx,
} from '../../../../generated';

jest.mock('../../../hooks/useSnackbar');

const mockGx: Gx = {
  server: {
    address: '174.16.1.14:3868',
  },
};

const mockGy: Gy = {
  server: {
    address: '174.18.1.0:3868',
  },
};

const mockSwx: Swx = {
  server: {
    address: '174.18.1.0:3869',
  },
};

const mockS6a: S6a = {
  server: {
    address: '174.18.1.0:2000',
  },
};

const mockCsfb: Csfb = {
  client: {
    server_address: '174.18.1.0:2200',
  },
};

const mockGw0: FederationGateway = {
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
    swx: {
      hlr_plmn_ids: [],
      server: {
        protocol: 'tcp',
      },
      servers: [],
    },
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

const lastFalloverTimeResponse1 = moment().unix();

const lastFalloverTimeResponse2 = moment().unix();

const lastFalloverTime = `${moment.unix(lastFalloverTimeResponse2).calendar()}`;

const mockFalloverStatus: PromqlReturnObject = {
  status: 'success',
  data: {
    resultType: 'matrix',
    result: [
      {
        metric: {},
        value: ['1625079646.98', `${lastFalloverTimeResponse1}`],
      },
      {
        metric: {},
        value: ['1625079647.98', `${lastFalloverTimeResponse2}`],
      },
    ],
  },
};

const mockGw1: FederationGateway = {
  ...mockGw0,
  id: 'test_gw1',
  name: 'test_gateway1',
  federation: {
    aaa_server: {},
    eap_aka: {},
    health: {},
    hss: {},
    served_network_ids: [],
    gx: {},
    gy: {},
    swx: mockSwx,
    s6a: {},
    csfb: {},
  },
};

const fegGateways = {
  [mockGw0.id]: mockGw0,
  [mockGw1.id]: mockGw1,
};

const mockHealthyGatewayStatus: FederationGatewayHealthStatus = {
  description: '',
  status: 'HEALTHY',
};

const mockUnhealthyGatewayStatus: FederationGatewayHealthStatus = {
  description: '',
  status: 'UNHEALTHY',
};

const mockClusterStatus: FederationNetworkClusterStatus = {
  active_gateway: mockGw0.id,
};

describe('<FEGEquipmentGateway />', () => {
  beforeEach(() => {
    // gateway context gets list of federation gateways
    mockAPI(
      MagmaAPI.federationGateways,
      'fegNetworkIdGatewaysGet',
      fegGateways,
    );
    // gateway context gets health status of the gateways
    jest
      .spyOn(
        MagmaAPI.federationGateways,
        'fegNetworkIdGatewaysGatewayIdHealthStatusGet',
      )
      .mockImplementation(
        (
          req: FederationGatewaysApiFegNetworkIdGatewaysGatewayIdHealthStatusGetRequest,
        ) => {
          if (req.gatewayId == mockGw0.id) {
            // only gateway 0 is healthy
            return Promise.resolve({
              data: mockHealthyGatewayStatus,
            } as AxiosResponse);
          }
          return Promise.resolve({
            data: mockUnhealthyGatewayStatus,
          } as AxiosResponse);
        },
      );
    // gateway context gets the active gateway id
    mockAPI(
      MagmaAPI.federationNetworks,
      'fegNetworkIdClusterStatusGet',
      mockClusterStatus,
    );

    // called by gateway checkin chart
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryRangeGet',
      mockCheckinMetric,
    );

    // called when getting max latency
    mockAPIOnce(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryGet',
      mockKPIMetric,
    );
    // called when getting min latency
    mockAPIOnce(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryGet',
      mockKPIMetric,
    );
    // called when getting avg latency
    mockAPIOnce(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryGet',
      mockKPIMetric,
    );

    // called when getting the last fallover time
    mockAPI(
      MagmaAPI.metrics,
      'networksNetworkIdPrometheusQueryGet',
      mockFalloverStatus,
    );
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateway']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGGatewayContextProvider networkId="mynetwork">
            <Routes>
              <Route
                path="/nms/:networkId/gateway/"
                element={<FEGEquipmentGateway />}
              />
            </Routes>
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

  it('renders federation gateway table correctly', async () => {
    const {getByTestId, getAllByRole, queryByTestId} = render(<Wrapper />);
    await waitFor(() => {
      const rowItems = getAllByRole('row');
      // first row is the header
      expect(rowItems[0]).toHaveTextContent('Name');
      expect(rowItems[0]).toHaveTextContent('Primary');
      expect(rowItems[0]).toHaveTextContent('Health');
      expect(rowItems[0]).toHaveTextContent('Gx');
      expect(rowItems[0]).toHaveTextContent('Gy');
      expect(rowItems[0]).toHaveTextContent('SWx');
      expect(rowItems[0]).toHaveTextContent('S6a');
      expect(rowItems[0]).toHaveTextContent('CSFB');

      expect(rowItems[1]).toHaveTextContent('test_gateway');
      expect(getByTestId('test_feg_gw0 is primary')).toBeVisible();
      expect(rowItems[1]).toHaveTextContent('Good');
      expect(rowItems[1]).toHaveTextContent('174.16.1.14:3868');
      expect(rowItems[1]).toHaveTextContent('174.18.1.0:3868');
      expect(rowItems[1]).toHaveTextContent('-');
      expect(rowItems[1]).toHaveTextContent('174.18.1.0:2000');
      expect(rowItems[1]).toHaveTextContent('174.18.1.0:2200');

      expect(rowItems[2]).toHaveTextContent('test_gateway1');
      expect(queryByTestId('test_gw1 is primary')).toBeNull();
      expect(rowItems[2]).toHaveTextContent('Bad');
      expect(rowItems[2]).toHaveTextContent('-');
      expect(rowItems[2]).toHaveTextContent('-');
      expect(rowItems[2]).toHaveTextContent('174.18.1.0:3869');
      expect(rowItems[2]).toHaveTextContent('-');
      expect(rowItems[2]).toHaveTextContent('-');
    });
  });

  it('renders cluster status correctly', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    // verify health status of the primary and secondary gateways
    expect(getByTestId('Primary Health')).toHaveTextContent('Good');
    expect(getByTestId('Secondary Health')).toHaveTextContent('Bad');
    // verify that primary/active gateway's name is rendered
    expect(getByTestId('Primary Gateway Name')).toHaveTextContent(
      'test_gateway',
    );
    // verify that correct fallover time is displayed
    expect(getByTestId('Last Fallover Time')).toHaveTextContent(
      lastFalloverTime,
    );
  });
});
