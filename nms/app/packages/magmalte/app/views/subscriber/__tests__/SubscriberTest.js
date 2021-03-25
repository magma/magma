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
import * as hooks from '../../../components/context/RefreshContext';

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkContext from '../../../components/context/NetworkContext';
import React from 'react';
import SubscriberContext from '../../../components/context/SubscriberContext';
import SubscriberDashboard from '../SubscriberOverview';
import defaultTheme from '../../../theme/default.js';

import {FEG_LTE} from '@fbcnms/types/network';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);
const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

describe('<SubscriberDashboard />', () => {
  const subscribers = {
    IMSI0000000000: {
      name: 'subscriber0',
      active_apns: ['oai.ipv4'],
      id: 'IMSI0000000000',
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
    IMSI0000000001: {
      name: 'subscriber1',
      active_apns: ['oai.ipv4'],
      id: 'IMSI0000000001',
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
  };

  const sessions = {
    IMSI0000000000: {
      subscriber_state: {
        apn_0: [
          {
            active_policy_rules: [
              {
                id: 'policy_0',
                priority: 2,
                flow_list: [
                  {
                    match: {
                      direction: 'UPLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                  {
                    match: {
                      direction: 'DOWNLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                ],
                tracking_type: 'NO_TRACKING',
              },
              {
                id: 'policy_01',
                priority: 2,
                flow_list: [
                  {
                    match: {
                      direction: 'UPLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                  {
                    match: {
                      direction: 'DOWNLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                ],
                tracking_type: 'NO_TRACKING',
              },
            ],
            lifecycle_state: 'SESSION_ACTIVE',
            session_id: 'IMSI0000000000-120333',
            active_duration_sec: 7,
            msisdn: '',
            apn: 'apn_0',
            session_start_time: 1605281201,
            ipv4: '192.168.128.217',
          },
          {
            active_policy_rules: [],
            lifecycle_state: 'SESSION_TERMINATED',
            session_id: 'IMSI0000000000-120335',
            active_duration_sec: 2,
            msisdn: '',
            apn: 'apn_0',
            session_start_time: 1605281100,
            ipv4: '192.168.128.217',
          },
        ],
        apn_1: [
          {
            active_policy_rules: [
              {
                id: 'policy_1',
                priority: 2,
                flow_list: [
                  {
                    match: {
                      direction: 'UPLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                  {
                    match: {
                      direction: 'DOWNLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                ],
                tracking_type: 'NO_TRACKING',
              },
            ],
            lifecycle_state: 'SESSION_ACTIVE',
            session_id: 'IMSI0000000000-120337',
            active_duration_sec: 7,
            msisdn: '',
            apn: 'apn_1',
            session_start_time: 1605281209,
            ipv4: '192.168.128.218',
          },
        ],
      },
    },
    IMSI0000000001: {
      subscriber_state: {
        apn_3: [
          {
            active_policy_rules: [],
            lifecycle_state: 'SESSION_TERMINATING',
            session_id: 'IMSI0000000001-120345',
            active_duration_sec: 1,
            msisdn: '',
            apn: 'apn_3',
            session_start_time: 1605281500,
            ipv4: '192.168.128.217',
          },
        ],
      },
    },
    IMSI0000000002: {
      subscriber_state: {
        apn_4: [
          {
            active_policy_rules: [
              {
                id: 'policy_2',
                priority: 2,
                flow_list: [
                  {
                    match: {
                      direction: 'UPLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                  {
                    match: {
                      direction: 'DOWNLINK',
                      ip_proto: 'IPPROTO_IP',
                    },
                    action: 'PERMIT',
                  },
                ],
                tracking_type: 'NO_TRACKING',
              },
            ],
            lifecycle_state: 'SESSION_ACTIVE',
            session_id: 'IMSI0000000002-120347',
            active_duration_sec: 27,
            msisdn: '',
            apn: 'apn_4',
            session_start_time: 1605281600,
            ipv4: '192.168.128.219',
          },
        ],
        apn_5: [
          {
            active_policy_rules: [],
            lifecycle_state: 'SESSION_TERMINATING',
            session_id: 'IMSI0000000002-120348',
            active_duration_sec: 1,
            msisdn: '',
            apn: 'apn_5',
            session_start_time: 1605281560,
            ipv4: '192.168.128.217',
          },
        ],
        apn_6: [
          {
            active_policy_rules: [],
            lifecycle_state: 'SESSION_TERMINATED',
            session_id: 'IMSI0000000002-120349',
            active_duration_sec: 1,
            msisdn: '',
            apn: 'apn_6',
            session_start_time: 1605281590,
            ipv4: '192.168.128.217',
          },
        ],
      },
    },
  };

  const Wrapper = ({networkType}) => {
    const subscriberCtx = {
      state: subscribers,
      gwSubscriberMap: {},
      sessionState: sessions,
    };

    jest
      .spyOn(hooks, 'useRefreshingContext')
      .mockImplementation(() => subscriberCtx);

    return (
      <MemoryRouter
        initialEntries={['/nms/test/subscribers/overview']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <NetworkContext.Provider
              value={{
                networkId: 'test',
                networkType: networkType,
              }}>
              <SubscriberContext.Provider value={subscriberCtx}>
                <Route
                  path="/nms/:networkId/subscribers/overview"
                  component={SubscriberDashboard}
                />
              </SubscriberContext.Provider>
            </NetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Dashboard', async () => {
    const {
      getAllByTitle,
      getAllByTestId,
      getByTestId,
      getAllByRole,
      getByText,
    } = render(<Wrapper networkType={FEG_LTE} />);
    await wait();
    const rowItems = await getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('IMSI');
    expect(rowItems[0]).toHaveTextContent('Service');
    expect(rowItems[0]).toHaveTextContent('Current Usage');
    expect(rowItems[0]).toHaveTextContent('Daily Average');
    expect(rowItems[0]).toHaveTextContent('Last Reported Time');
    expect(rowItems[0]).toHaveTextContent('Active APNs');
    expect(rowItems[0]).toHaveTextContent('Active Sessions');

    expect(rowItems[1]).toHaveTextContent('subscriber0');
    expect(rowItems[1]).toHaveTextContent('IMSI0000000000');
    expect(rowItems[1]).toHaveTextContent('ACTIVE');
    expect(rowItems[1]).toHaveTextContent('0');
    expect(rowItems[1]).toHaveTextContent('2');
    expect(rowItems[1]).toHaveTextContent('apn_0,apn_1');
    expect(rowItems[1]).toHaveTextContent('-');

    expect(rowItems[2]).toHaveTextContent('subscriber1');
    expect(rowItems[2]).toHaveTextContent('IMSI0000000001');
    expect(rowItems[2]).toHaveTextContent('INACTIVE');
    expect(rowItems[2]).toHaveTextContent('0');
    expect(rowItems[2]).toHaveTextContent('1');
    expect(rowItems[2]).toHaveTextContent('apn_3');
    expect(rowItems[2]).toHaveTextContent('-');

    expect(rowItems[3]).toHaveTextContent('IMSI0000000002');
    expect(rowItems[3]).toHaveTextContent('0');
    expect(rowItems[3]).toHaveTextContent('1');
    expect(rowItems[3]).toHaveTextContent('apn_4,apn_5,apn_6');
    expect(rowItems[3]).toHaveTextContent('192.168.128.219');

    // click the actions button for subscriber0
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[0]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();

    const details = getAllByTestId('details');
    fireEvent.click(details[0]);
    await wait();
    expect(getByText('APN Name')).toBeVisible();
    expect(getByText('Session ID')).toBeVisible();
    expect(getByText('State')).toBeVisible();
    expect(getByText('Active Duration')).toBeVisible();
    expect(getByText('Active Policy IDs')).toBeVisible();
  });
});
