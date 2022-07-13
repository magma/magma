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

import FEGGatewayDetailSubscribers from '../FEGGatewayDetailSubscribers';
import FEGSubscriberContext from '../../../components/context/FEGSubscriberContext';
import MagmaAPI from '../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {AxiosResponse} from 'axios';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SubscriberId} from '../../../../shared/types/network';
import {render, wait} from '@testing-library/react';
import type {FederationGateway, SubscriberState} from '../../../../generated';

const mockSubscriberIds: Array<SubscriberId> = [
  'IMSI001011234565000',
  'IMSI001011234565001',
];
const mockSubscribers = [
  {
    name: 'subscriber0',
    active_apns: ['oai.ipv4'],
    id: mockSubscriberIds[0],
    lte: {
      auth_algo: 'MILENAGE',
      auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
      auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
      state: 'ACTIVE',
      sub_profile: 'default',
    },
    config: {
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'ACTIVE',
        sub_profile: 'default',
      },
    },
  },
  {
    name: 'subscriber1',
    active_apns: ['oai.ipv4'],
    id: mockSubscriberIds[1],
    lte: {
      auth_algo: 'MILENAGE',
      auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
      auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
      state: 'INACTIVE',
      sub_profile: 'default',
    },
    config: {
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'INACTIVE',
        sub_profile: 'default',
      },
    },
  },
];
const mockSubscriberSessionState0 = {
  [mockSubscriberIds[0]]: {
    directory: {},
  },
} as Record<SubscriberId, SubscriberState>;
const mockSubscriberSessionState1 = {
  [mockSubscriberIds[1]]: {
    mme: {
      accessRestrictionData: 32,
    },
  },
} as Record<SubscriberId, SubscriberState>;

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

describe('<FEGGatewayDetailSubscribers />', () => {
  beforeEach(() => {
    jest
      .spyOn(MagmaAPI.subscribers, 'lteNetworkIdSubscribersSubscriberIdGet')
      .mockResolvedValueOnce({data: mockSubscribers[0]} as AxiosResponse)
      .mockResolvedValue({data: mockSubscribers[1]} as AxiosResponse);
  });

  const Wrapper = () => (
    <MemoryRouter
      initialEntries={[
        '/nms/mynetwork/equipment/overview/gateway/test_feg_gw0/overview',
      ]}
      initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGSubscriberContext.Provider
            value={{
              refetch: () => {},
              sessionState: {
                feg_lte_network1: mockSubscriberSessionState0,
                feg_lte_network2: mockSubscriberSessionState1,
              },
              setSessionState: () => {},
            }}>
            <Routes>
              <Route
                path="/nms/:networkId/equipment/overview/gateway/:gatewayId/overview"
                element={
                  <FEGGatewayDetailSubscribers
                    refresh={false}
                    gwInfo={mockGw0}
                  />
                }
              />
            </Routes>
          </FEGSubscriberContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders gateway detail subscribers table correctly', async () => {
    const {getAllByRole} = render(<Wrapper />);
    await wait();
    // two subscribers in the network
    expect(
      MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdGet,
    ).toHaveBeenCalledTimes(2);

    const rowItems = getAllByRole('row');
    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('Subscriber ID');
    expect(rowItems[0]).toHaveTextContent('Service');
    expect(rowItems[1]).toHaveTextContent('subscriber0');
    expect(rowItems[1]).toHaveTextContent('IMSI001011234565000');
    expect(rowItems[1]).toHaveTextContent('ACTIVE');
    expect(rowItems[2]).toHaveTextContent('subscriber1');
    expect(rowItems[2]).toHaveTextContent('IMSI001011234565001');
    expect(rowItems[2]).toHaveTextContent('INACTIVE');
  });
});
