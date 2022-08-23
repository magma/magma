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
import ApnContext from '../../../context/ApnContext';
import LteNetworkContext, {
  LteNetworkContextType,
} from '../../../context/LteNetworkContext';
import PolicyContext, {PolicyContextType} from '../../../context/PolicyContext';
import React from 'react';
import SubscriberDashboard from '../SubscriberOverview';
import SubscriberDetailConfig from '../SubscriberDetailConfig';
import defaultTheme from '../../../theme/default';
import {SubscriberContextProvider} from '../../../context/SubscriberContext';

import MagmaAPI from '../../../api/MagmaAPI';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {
  NetworkEpcConfigs,
  NetworkRanConfigs,
  PolicyRule,
  Subscriber,
} from '../../../../generated';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, render, waitFor, within} from '@testing-library/react';
import {forbiddenNetworkTypes} from '../SubscriberUtils';
import {mockAPI} from '../../../util/TestUtils';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';

jest.mock('../../../hooks/useSnackbar');

jest.setTimeout(30000);

const subscribersMock: Record<string, Subscriber> = {
  IMSI00000000001002: {
    name: 'subscriber0',
    active_apns: ['apn_0'],
    id: 'IMSI00000000001002',
    forbidden_network_types: forbiddenNetworkTypes,
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
      static_ips: {apn_0: '1.1.1.1'},
    },
  },
  IMSI00000000001003: {
    name: 'subscriber1',
    active_apns: [],
    id: 'IMSI00000000001003',
    forbidden_network_types: forbiddenNetworkTypes,
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
const policies: Record<string, PolicyRule> = {
  policy_0: {
    flow_list: [],
    id: 'policy_0',
    monitoring_key: '',
    priority: 1,
  },
  policy_1: {
    flow_list: [
      {
        action: 'PERMIT',
        match: {
          direction: 'UPLINK',
          ip_proto: 'IPPROTO_IP',
        },
      },
      {
        action: 'PERMIT',
        match: {
          direction: 'DOWNLINK',
          ip_proto: 'IPPROTO_IP',
        },
      },
    ],
    id: 'policy_1',
    monitoring_key: '',
    priority: 1,
  },
  policy_2: {
    flow_list: [],
    id: 'policy_2',
    monitoring_key: '',
    priority: 10,
  },
};

const apns = {
  apn_0: {
    apn_configuration: {
      ambr: {
        max_bandwidth_dl: 200000000,
        max_bandwidth_ul: 100000000,
      },
      qos_profile: {
        class_id: 9,
        preemption_capability: true,
        preemption_vulnerability: false,
        priority_level: 15,
      },
    },
    apn_name: 'apn_0',
  },
  apn_1: {
    apn_configuration: {
      ambr: {
        max_bandwidth_dl: 200000000,
        max_bandwidth_ul: 100000000,
      },
      qos_profile: {
        class_id: 9,
        preemption_capability: true,
        preemption_vulnerability: false,
        priority_level: 15,
      },
    },
    apn_name: 'apn_1',
  },
};
const testNetwork = {
  description: 'Test Network Description',
  id: 'test_network',
  name: 'Test Network',
  dns: {
    enable_caching: false,
    local_ttl: 0,
    records: [],
  },
};
const epc: NetworkEpcConfigs = {
  default_rule_id: 'default_rule_1',
  lte_auth_amf: 'gAA=',
  lte_auth_op: 'EREREREREREREREREREREQ==',
  mcc: '001',
  mnc: '01',
  network_services: ['dpi', 'policy_enforcement'],
  hss_relay_enabled: false,
  gx_gy_relay_enabled: false,
  sub_profiles: {
    additionalProp1: {
      max_dl_bit_rate: 20000000,
      max_ul_bit_rate: 100000000,
    },
    additionalProp2: {
      max_dl_bit_rate: 20000000,
      max_ul_bit_rate: 100000000,
    },
    additionalProp3: {
      max_dl_bit_rate: 20000000,
      max_ul_bit_rate: 100000000,
    },
  },
  mobility: {
    ip_allocation_mode: 'NAT',
    enable_static_ip_assignments: false,
    enable_multi_apn_ip_allocation: false,
  },
  tac: 1,
};

const ran: NetworkRanConfigs = {
  bandwidth_mhz: 20,
  tdd_config: {
    earfcndl: 44390,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
  },
};

describe('<AddSubscriberButton />', () => {
  beforeEach(() => {
    (useEnqueueSnackbar as jest.Mock).mockReturnValue(jest.fn());
    mockAPI(MagmaAPI.subscribers, 'lteNetworkIdSubscribersPost');
    mockAPI(MagmaAPI.subscribers, 'lteNetworkIdSubscribersSubscriberIdPut');

    mockAPI(MagmaAPI.subscribers, 'lteNetworkIdSubscribersGet', {
      subscribers: subscribersMock,
      next_page_token: 'foo',
      total_count: 42,
    });
    mockAPI(MagmaAPI.subscribers, 'lteNetworkIdSubscriberStateGet', {});
  });

  const AddWrapper = () => {
    const policyCtx: PolicyContextType = {
      state: policies,
      baseNames: {},
      qosProfiles: {},
      ratingGroups: {},
      setBaseNames: async () => {},
      setRatingGroups: async () => {},
      setQosProfiles: async () => {},
      setState: async () => {},
      refetch: () => {},
    };

    const apnCtx = {
      state: apns,
      setState: async () => {},
    };

    const networkCtx: LteNetworkContextType = {
      state: {
        ...testNetwork,
        cellular: {
          epc: epc,
          ran: ran,
        },
      },
      updateNetworks: async () => {},
    };

    return (
      <MemoryRouter initialEntries={['/nms/test/subscribers']} initialIndex={0}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={defaultTheme}>
            <ThemeProvider theme={defaultTheme}>
              <LteNetworkContext.Provider value={networkCtx}>
                <PolicyContext.Provider value={policyCtx}>
                  <ApnContext.Provider value={apnCtx}>
                    <SubscriberContextProvider networkId="test">
                      <Routes>
                        <Route
                          path="/nms/:networkId/subscribers/*"
                          element={<SubscriberDashboard />}
                        />
                      </Routes>
                    </SubscriberContextProvider>
                  </ApnContext.Provider>
                </PolicyContext.Provider>
              </LteNetworkContext.Provider>
            </ThemeProvider>
          </ThemeProvider>
        </StyledEngineProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    const policyCtx: PolicyContextType = {
      state: policies,
      baseNames: {},
      qosProfiles: {},
      ratingGroups: {},
      setBaseNames: async () => {},
      setRatingGroups: async () => {},
      setQosProfiles: async () => {},
      setState: async () => {},
      refetch: () => {},
    };

    const apnCtx = {
      state: apns,
      setState: async () => {},
    };

    const networkCtx = {
      state: {
        ...testNetwork,
        cellular: {
          epc: epc,
          ran: ran,
        },
      },
      updateNetworks: async () => {},
    };
    return (
      <MemoryRouter
        initialEntries={[
          '/nms/test/subscribers/overview/IMSI00000000001002/config',
        ]}
        initialIndex={0}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={defaultTheme}>
            <ThemeProvider theme={defaultTheme}>
              <LteNetworkContext.Provider value={networkCtx}>
                <PolicyContext.Provider value={policyCtx}>
                  <ApnContext.Provider value={apnCtx}>
                    <SubscriberContextProvider networkId="test">
                      <Routes>
                        <Route
                          path="/nms/:networkId/subscribers/overview/:subscriberId/config"
                          element={<SubscriberDetailConfig />}
                        />
                      </Routes>
                    </SubscriberContextProvider>
                  </ApnContext.Provider>
                </PolicyContext.Provider>
              </LteNetworkContext.Provider>
            </ThemeProvider>
          </ThemeProvider>
        </StyledEngineProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Add', async () => {
    const {
      getByTestId,
      queryByTestId,
      getByRole,
      findByRole,
      findByTestId,
      findByText,
    } = render(<AddWrapper />);

    expect(queryByTestId('addSubscriberDialog')).toBeNull();
    // Add Subscriber
    fireEvent.click(await findByText('Manage Subscribers'));
    fireEvent.click(await findByText('Add Subscribers'));

    expect(await findByTestId('addSubscriberDialog')).not.toBeNull();

    const detailsTable = getByTestId('subscriber-details-table');

    // first row is the header
    const rowHeader = within(detailsTable).getAllByRole('row', {hidden: true});
    expect(rowHeader[0]).toHaveTextContent('Subscriber Name');
    expect(rowHeader[0]).toHaveTextContent('IMSI');
    expect(rowHeader[0]).toHaveTextContent('Auth Key');
    expect(rowHeader[0]).toHaveTextContent('Auth OPC');
    expect(rowHeader[0]).toHaveTextContent('Service');
    expect(rowHeader[0]).toHaveTextContent('Data Plan');
    expect(rowHeader[0]).toHaveTextContent('Active APNs');
    expect(rowHeader[0]).toHaveTextContent('Active Policies');

    //Add subscriber
    fireEvent.click(await findByRole('button', {name: 'Add'}));

    const name = (await findByTestId('name')).firstChild;
    const IMSI = getByTestId('IMSI').firstChild;
    const authKey = getByTestId('authKey').firstChild;
    const authOpc = getByTestId('authOpc').firstChild;
    const service = getByTestId('service').firstChild;
    const dataPlan = getByTestId('dataPlan').firstChild;

    if (
      name instanceof HTMLInputElement &&
      IMSI instanceof HTMLInputElement &&
      authKey instanceof HTMLElement &&
      authOpc instanceof HTMLElement &&
      service instanceof HTMLElement &&
      dataPlan instanceof HTMLElement
    ) {
      fireEvent.change(name, {target: {value: 'IMSI00000000001004'}});
      fireEvent.change(IMSI, {target: {value: 'IMSI00000000001004'}});
      fireEvent.change(authKey, {
        target: {value: '8baf473f2f8fd09487cccbd7097c6862'},
      });
      fireEvent.change(authOpc, {
        target: {value: '8e27b6af0e692e750f32667a3b14605d'},
      });
    } else {
      throw 'invalid type';
    }

    // Add subscriber
    fireEvent.click(getByRole('button', {name: 'Save'}));

    // Verify new subscriber row before saving
    const rowItems = await within(detailsTable).findAllByRole('row', {
      hidden: true,
    });
    expect(rowItems[1]).toHaveTextContent('IMSI00000000001004');
    expect(rowItems[1]).toHaveTextContent('8baf473f2f8fd09487cccbd7097c6862');
    expect(rowItems[1]).toHaveTextContent('8e27b6af0e692e750f32667a3b14605d');
    expect(rowItems[1]).toHaveTextContent('ACTIVE');
    expect(rowItems[1]).toHaveTextContent('default');

    // Save subscriber
    fireEvent.click(getByTestId('saveSubscriber'));
    expect(
      MagmaAPI.subscribers.lteNetworkIdSubscribersPost,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      subscribers: [
        {
          active_apns: undefined,
          active_policies: undefined,
          id: 'IMSI00000000001004',
          lte: {
            auth_algo: 'MILENAGE',
            auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
            auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
            state: 'ACTIVE',
            sub_profile: 'default',
          },
          name: 'IMSI00000000001004',
        },
      ],
    });
  });

  it('Verify Subscriber edit', async () => {
    const {getByTestId, queryByTestId, findByTestId} = render(
      <DetailWrapper />,
    );
    expect(queryByTestId('editDialog')).toBeNull();

    // Edit tab 1 : subscriber info
    fireEvent.click(await findByTestId('subscriber'));
    expect(await findByTestId('editDialog')).not.toBeNull();

    const name = getByTestId('name').firstChild;

    if (name instanceof HTMLInputElement) {
      fireEvent.change(name, {target: {value: 'test_subscriber'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByTestId('subscriber-saveButton'));

    await waitFor(() => {
      expect(
        MagmaAPI.subscribers.lteNetworkIdSubscribersSubscriberIdPut,
      ).toHaveBeenCalledWith({
        networkId: 'test',
        subscriberId: 'IMSI00000000001002',
        subscriber: {
          active_apns: ['apn_0'],
          active_base_names: undefined,
          forbidden_network_types: forbiddenNetworkTypes,
          id: 'IMSI00000000001002',
          lte: {
            auth_algo: 'MILENAGE',
            auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
            auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
            state: 'ACTIVE',
            sub_profile: 'default',
          },
          name: 'test_subscriber',
          static_ips: {apn_0: '1.1.1.1'},
        },
      });
    });
  });
});
