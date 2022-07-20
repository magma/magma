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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
import SubscriberContext from '../../../components/context/SubscriberContext';
import SubscriberDashboard from '../SubscriberOverview';
import defaultTheme from '../../../theme/default';

import MagmaAPI from '../../../api/MagmaAPI';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {Subscriber} from '../../../../generated';
import {fireEvent, render, waitFor} from '@testing-library/react';
import {forbiddenNetworkTypes} from '../SubscriberUtils';
import {mockAPI} from '../../../util/TestUtils';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

const subscribers: Record<string, Subscriber> = {
  IMSI0000000000: {
    name: 'subscriber0',
    active_apns: ['oai.ipv4'],
    forbidden_network_types: forbiddenNetworkTypes,
    id: 'IMSI0000000000',
    lte: {
      auth_algo: 'MILENAGE',
      auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
      auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
      state: 'ACTIVE',
      sub_profile: 'default',
    },
    config: {
      forbidden_network_types: forbiddenNetworkTypes,
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'ACTIVE',
        sub_profile: 'default',
      },
    },
  },
  IMSI0000000001: {
    name: 'subscriber1',
    active_apns: ['oai.ipv4'],
    forbidden_network_types: forbiddenNetworkTypes,
    id: 'IMSI0000000001',
    lte: {
      auth_algo: 'MILENAGE',
      auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
      auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
      state: 'INACTIVE',
      sub_profile: 'default',
    },
    config: {
      forbidden_network_types: forbiddenNetworkTypes,
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'INACTIVE',
        sub_profile: 'default',
      },
    },
  },
};

describe('<SubscriberDashboard />', () => {
  beforeEach(() => {
    mockAPI(MagmaAPI.subscribers, 'lteNetworkIdSubscribersGet', {
      subscribers,
      next_page_token: '',
      total_count: 2,
    });

    mockAPI(MagmaAPI.networks, 'networksGet', []);
    mockAPI(MagmaAPI.networks, 'networksNetworkIdTypeGet');
  });

  const Wrapper = () => {
    const subscriberCtx = {
      state: subscribers,
      forbiddenNetworkTypes: {},
      gwSubscriberMap: {},
      sessionState: {},
      totalCount: 2,
      refetchSessionState: () => {},
    };

    return (
      <MemoryRouter
        initialEntries={['/nms/test/subscribers/overview']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <NetworkContext.Provider
              value={{
                networkId: 'test',
              }}>
              <SubscriberContext.Provider value={subscriberCtx}>
                <Routes>
                  <Route
                    path="/nms/:networkId/subscribers/overview/*"
                    element={<SubscriberDashboard />}
                  />
                </Routes>
              </SubscriberContext.Provider>
            </NetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Dashboard', async () => {
    const {getAllByTitle, getAllByRole, getByTestId} = render(<Wrapper />);

    await waitFor(() => {
      const rowItems = getAllByRole('row');
      // first row is the header
      expect(rowItems[0]).toHaveTextContent('Name');
      expect(rowItems[0]).toHaveTextContent('IMSI');
      expect(rowItems[0]).toHaveTextContent('Service');
      expect(rowItems[0]).toHaveTextContent('Current Usage');
      expect(rowItems[0]).toHaveTextContent('Daily Average');
      expect(rowItems[0]).toHaveTextContent('Last Reported Time');

      expect(rowItems[1]).toHaveTextContent('subscriber0');
      expect(rowItems[1]).toHaveTextContent('IMSI0000000000');
      expect(rowItems[1]).toHaveTextContent('ACTIVE');
      expect(rowItems[1]).toHaveTextContent('0');

      expect(rowItems[2]).toHaveTextContent('subscriber1');
      expect(rowItems[2]).toHaveTextContent('IMSI0000000001');
      expect(rowItems[2]).toHaveTextContent('INACTIVE');
      expect(rowItems[2]).toHaveTextContent('0');
    });

    // click the actions button for subscriber0
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[0]);
    await waitFor(() => {
      expect(getByTestId('actions-menu')).toBeVisible();
    });
  });

  it('subscribers can be searched by IMSI', async () => {
    const geSubscribersBySubscriberId = mockAPI(
      MagmaAPI.subscribers,
      'lteNetworkIdSubscribersSubscriberIdGet',
      subscribers['IMSI0000000001'],
    );

    const {findByPlaceholderText, findAllByRole} = render(<Wrapper />);

    const searchInput = await findByPlaceholderText(
      'Search IMSI001011234560000',
    );

    fireEvent.change(searchInput, {
      target: {value: 'IMSI0000000001'},
    });

    await waitFor(() =>
      expect(geSubscribersBySubscriberId).toBeCalledWith({
        networkId: 'test',
        subscriberId: 'IMSI0000000001',
      }),
    );

    const rowItems = await findAllByRole('row');

    expect(rowItems[1]).toHaveTextContent('subscriber1');
    expect(rowItems[1]).toHaveTextContent('IMSI0000000001');
    expect(rowItems[1]).toHaveTextContent('INACTIVE');
    expect(rowItems[1]).toHaveTextContent('0');
  });
});
