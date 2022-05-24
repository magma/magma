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
import * as hooks from '../../../components/context/RefreshContext';

// $FlowFixMe migrated to typescript
import ApnContext from '../../../components/context/ApnContext';
import LteNetworkContext from '../../../components/context/LteNetworkContext';
import MagmaAPIBindings from '../../../../generated/MagmaAPIBindings';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import PolicyContext from '../../../components/context/PolicyContext';
import React from 'react';
// $FlowFixMe migrated to typescript
import SubscriberContext from '../../../components/context/SubscriberContext';
import SubscriberDashboard from '../SubscriberOverview';
import SubscriberDetailConfig from '../SubscriberDetailConfig';
import defaultTheme from '../../../theme/default.js';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {CoreNetworkTypes} from '../SubscriberUtils';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, wait} from '@testing-library/react';
import {setSubscriberState} from '../../../state/lte/SubscriberState';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import {useState} from 'react';

jest.mock('axios');
jest.mock('../../../../generated/MagmaAPIBindings.js');
jest.mock('../../../hooks/useSnackbar');

const forbiddenNetworkTypes = Object.keys(CoreNetworkTypes).map(
  key => CoreNetworkTypes[key],
);

const subscribersMock = {
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
const policies = {
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
const epc = {
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

const ran = {
  bandwidth_mhz: 20,
  tdd_config: {
    earfcndl: 44390,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
  },
};

describe('<AddSubscriberButton />', () => {
  beforeEach(() => {
    (useEnqueueSnackbar: JestMockFn<
      Array<empty>,
      $Call<typeof useEnqueueSnackbar>,
    >).mockReturnValue(jest.fn());
    MagmaAPIBindings.getLteByNetworkIdSubscriberConfigBaseNames.mockResolvedValue(
      [],
    );
    MagmaAPIBindings.getNetworks.mockResolvedValue([]);
  });

  const AddWrapper = () => {
    const [subscribers, setSubscribers] = useState(subscribersMock);
    const [sessionState, setSessionState] = useState({});
    const [forbiddenNetworkTypes, setForbiddenNetworkTypes] = useState({});

    const subscriberCtx = {
      state: subscribers,
      forbiddenNetworkTypes: forbiddenNetworkTypes,
      gwSubscriberMap: {},
      sessionState: sessionState,
      totalCount: 2,
      setState: async (key, value?) =>
        setSubscriberState({
          networkId: 'test',
          subscriberMap: subscribers,
          setSubscriberMap: setSubscribers,
          key: key,
          value: value,
          setSessionState: setSessionState,
          setForbiddenNetworkTypes: setForbiddenNetworkTypes,
        }),
    };
    const policyCtx = {
      state: policies,
      baseNames: {},
      qosProfiles: {},
      ratingGroups: {},
      setBaseNames: async () => {},
      setRatingGroups: async () => {},
      setQosProfiles: async () => {},
      setState: async () => {},
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

    jest
      .spyOn(hooks, 'useRefreshingContext')
      .mockImplementation(() => subscriberCtx);

    return (
      <MemoryRouter initialEntries={['/nms/test/subscribers']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <LteNetworkContext.Provider value={networkCtx}>
              <PolicyContext.Provider value={policyCtx}>
                <ApnContext.Provider value={apnCtx}>
                  <SubscriberContext.Provider value={subscriberCtx}>
                    <Routes>
                      <Route
                        path="/nms/:networkId/subscribers/*"
                        element={<SubscriberDashboard />}
                      />
                    </Routes>
                  </SubscriberContext.Provider>
                </ApnContext.Provider>
              </PolicyContext.Provider>
            </LteNetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    const [subscribers, setSubscribers] = useState(subscribersMock);
    const [sessionState, setSessionState] = useState({});
    const [forbiddenNetworkTypes, setForbiddenNetworkTypes] = useState({});
    const policyCtx = {
      state: policies,
      baseNames: {},
      qosProfiles: {},
      ratingGroups: {},
      setBaseNames: async () => {},
      setRatingGroups: async () => {},
      setQosProfiles: async () => {},
      setState: async () => {},
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
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <LteNetworkContext.Provider value={networkCtx}>
              <PolicyContext.Provider value={policyCtx}>
                <ApnContext.Provider value={apnCtx}>
                  <SubscriberContext.Provider
                    value={{
                      state: {
                        IMSI00000000001002: subscribers['IMSI00000000001002'],
                      },
                      gwSubscriberMap: {},
                      forbiddenNetworkTypes: forbiddenNetworkTypes,
                      sessionState: sessionState,
                      totalCount: 1,
                      setState: (key, value?) =>
                        setSubscriberState({
                          networkId: 'test',
                          subscriberMap: subscribers,
                          setSubscriberMap: setSubscribers,
                          setSessionState,
                          setForbiddenNetworkTypes,
                          key: key,
                          value: value,
                        }),
                    }}>
                    <Routes>
                      <Route
                        path="/nms/:networkId/subscribers/overview/:subscriberId/config"
                        element={<SubscriberDetailConfig />}
                      />
                    </Routes>
                  </SubscriberContext.Provider>
                </ApnContext.Provider>
              </PolicyContext.Provider>
            </LteNetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Subscribers Add', async () => {
    const {
      getByTestId,
      getByText,
      queryByTestId,
      getByTitle,
      getAllByRole,
    } = render(<AddWrapper />);
    await wait();

    expect(queryByTestId('addSubscriberDialog')).toBeNull();
    // Add Subscriber
    fireEvent.click(getByText('Manage Subscribers'));
    await wait();
    fireEvent.click(getByText('Add Subscribers'));
    await wait();
    expect(queryByTestId('addSubscriberDialog')).not.toBeNull();

    // first row is the header
    const rowHeader = await getAllByRole('row');
    expect(rowHeader[0]).toHaveTextContent('Subscriber Name');
    expect(rowHeader[0]).toHaveTextContent('IMSI');
    expect(rowHeader[0]).toHaveTextContent('Auth Key');
    expect(rowHeader[0]).toHaveTextContent('Auth OPC');
    expect(rowHeader[0]).toHaveTextContent('Service');
    expect(rowHeader[0]).toHaveTextContent('Data Plan');
    expect(rowHeader[0]).toHaveTextContent('Active APNs');
    expect(rowHeader[0]).toHaveTextContent('Active Policies');

    //Add subscriber
    fireEvent.click(getByTitle('Add'));
    await wait();
    const name = getByTestId('name').firstChild;
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
    fireEvent.click(getByTitle('Save'));
    await wait();

    // Verify new subscriber row before saving
    const rowItems = await getAllByRole('row');
    expect(rowItems[1]).toHaveTextContent('IMSI00000000001004');
    expect(rowItems[1]).toHaveTextContent('8baf473f2f8fd09487cccbd7097c6862');
    expect(rowItems[1]).toHaveTextContent('8e27b6af0e692e750f32667a3b14605d');
    expect(rowItems[1]).toHaveTextContent('ACTIVE');
    expect(rowItems[1]).toHaveTextContent('default');

    // Save subscriber
    fireEvent.click(getByTestId('saveSubscriber'));
    expect(MagmaAPIBindings.postLteByNetworkIdSubscribers).toHaveBeenCalledWith(
      {
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
      },
    );
  });

  it('Verify Subscriber edit', async () => {
    const {getByTestId, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    // Edit tab 1 : subscriber info
    fireEvent.click(getByTestId('subscriber'));
    await wait();
    expect(queryByTestId('editDialog')).not.toBeNull();

    const name = getByTestId('name').firstChild;

    if (name instanceof HTMLInputElement) {
      fireEvent.change(name, {target: {value: 'test_subscriber'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByTestId('subscriber-saveButton'));
    await wait();

    expect(
      MagmaAPIBindings.putLteByNetworkIdSubscribersBySubscriberId,
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
    // TODO: Test other tabs
  });
});
