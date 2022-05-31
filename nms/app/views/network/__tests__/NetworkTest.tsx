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

import ApnContext from '../../../components/context/ApnContext';
import EnodebContext, {
  EnodebContextType,
} from '../../../components/context/EnodebContext';
import FEGNetworkContext from '../../../components/context/FEGNetworkContext';
import FEGNetworkDashboard from '../FEGNetworkDashboard';
import GatewayContext, {
  GatewayContextType,
} from '../../../components/context/GatewayContext';
import LteNetworkContext, {
  LteNetworkContextType,
  UpdateNetworkContextProps,
} from '../../../components/context/LteNetworkContext';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import NetworkDashboard from '../NetworkDashboard';
import PolicyContext, {
  PolicyContextType,
} from '../../../components/context/PolicyContext';
import React from 'react';
import SubscriberContext, {
  SubscriberContextType,
} from '../../../components/context/SubscriberContext';
import axiosMock, {AxiosResponse} from 'axios';
import defaultTheme from '../../../theme/default';

import {CoreNetworkTypes} from '../../subscriber/SubscriberUtils';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {
  UpdateNetworkProps,
  UpdateNetworkState,
} from '../../../state/lte/NetworkState';
import {fireEvent, render, wait} from '@testing-library/react';

import MagmaAPI from '../../../../api/MagmaAPI';
import axios from 'axios';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import type {FegNetwork, NetworkEpcConfigs} from '../../../../generated-ts';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

const forbiddenNetworkTypes = (Object.keys(CoreNetworkTypes) as Array<
  keyof typeof CoreNetworkTypes
>).map(key => CoreNetworkTypes[key]);

describe('<NetworkDashboard />', () => {
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

  const ran = {
    bandwidth_mhz: 20,
    tdd_config: {
      earfcndl: 44390,
      special_subframe_pattern: 7,
      subframe_assignment: 2,
    },
  };

  const gateways = {
    test_gateway1: {
      id: 'test_gw1',
      name: 'test_gateway',
      description: 'hello I am a gateway',
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
      checked_in_recently: false,
      connected_enodeb_serials: [],
      cellular: {
        epc: {
          ip_block: '192.168.0.1/24',
          nat_enabled: true,
        },
        ran: {
          pci: 620,
          transmit_enabled: true,
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
    },
  };

  const enbInfo = {
    '120200020718CJP0013': {
      enb: {
        attached_gateway_id: 'mpk_dogfooding_tiplab_1',
        config: {
          bandwidth_mhz: 10,
          cell_id: 6553601,
          device_class: 'Baicells Nova-233 G2 OD FDD',
          earfcndl: 9410,
          pci: 60,
          tac: 6,
          transmit_enabled: true,
        },
        name: '',
        serial: '120200020718CJP0013',
      },
      enb_state: {
        enodeb_configured: true,
        enodeb_connected: true,
        fsm_state: '',
        gps_connected: true,
        gps_latitude: '',
        gps_longitude: '',
        mme_connected: true,
        opstate_enabled: true,
        ptp_connected: true,
        rf_tx_desired: true,
        rf_tx_on: true,
        ip_address: '192.168.1.254',
      },
    },
  };

  const policies = {
    test1: {
      flow_list: [],
      id: 'test',
      priority: 10,
      redirect: {
        address_type: 'IPv4',
        server_address: 'http://localhost:8080',
        support: 'ENABLED',
      },
    },
    test2: {
      flow_list: [],
      id: 'test',
      priority: 10,
      redirect: {
        address_type: 'IPv4',
        server_address: 'http://localhost:8080',
        support: 'ENABLED',
      },
    },
  };

  const subscribers = {
    IMSI00000000001002: {
      active_apns: ['oai.ipv4'],
      forbidden_network_types: forbiddenNetworkTypes,
      id: 'IMSI722070171001002',
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
  };

  const apns = {
    internet: {
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
      apn_name: 'internet',
    },
    'oai.ipv4': {
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
      apn_name: 'oai.ipv4',
    },
  };

  beforeEach(() => {
    (useEnqueueSnackbar as jest.MockedFunction<
      typeof useEnqueueSnackbar
    >).mockReturnValue(jest.fn());

    (axiosMock as jest.Mocked<typeof axios>).post.mockImplementation(() =>
      Promise.resolve({data: {success: true}}),
    );
    jest.spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdGet').mockImplementation();
    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdPut')
      .mockResolvedValue({data: {success: true}} as AxiosResponse);

    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdCellularEpcPut')
      .mockResolvedValue({data: {success: true}} as AxiosResponse);
    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdCellularRanPut')
      .mockResolvedValue({data: {success: true}} as AxiosResponse);
    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdDnsPut')
      .mockResolvedValue({data: {success: true}} as AxiosResponse);
    jest
      .spyOn(MagmaAPI.networks, 'networksGet')
      .mockResolvedValue({data: []} as AxiosResponse);
  });

  const Wrapper = () => {
    const apnCtx = {
      state: apns,
      setState: async () => {},
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
    } as PolicyContextType;
    const enodebCtx = {
      state: {enbInfo},
      setState: async () => {},
    } as EnodebContextType;

    const gatewayCtx = {
      state: gateways,
      setState: async () => {},
      updateGateway: async () => {},
    } as GatewayContextType;

    const subscriberCtx = {
      state: subscribers,
      forbidden_network_types: subscribers,
      totalCount: 1,
      forbiddenNetworkTypes: {},
      gwSubscriberMap: {},
      sessionState: {},
    } as SubscriberContextType;

    const networkCtx = {
      state: {
        ...testNetwork,

        cellular: {
          epc: epc,
          ran: ran,
        },
      },
      updateNetworks: async (props: UpdateNetworkContextProps) => {
        return UpdateNetworkState({
          setLteNetwork: () => {},
          refreshState: testNetwork.id === props.networkId,
          ...props,
        } as UpdateNetworkProps); // TODO[TS-migration] Broken LteNetworkContext type
      },
    } as LteNetworkContextType;

    return (
      <MemoryRouter initialEntries={['/nms/test/network']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <LteNetworkContext.Provider value={networkCtx}>
              <PolicyContext.Provider value={policyCtx}>
                <ApnContext.Provider value={apnCtx}>
                  <GatewayContext.Provider value={gatewayCtx}>
                    <EnodebContext.Provider value={enodebCtx}>
                      <SubscriberContext.Provider value={subscriberCtx}>
                        <Routes>
                          <Route
                            path="/nms/:networkId/network/*"
                            element={<NetworkDashboard />}
                          />
                        </Routes>
                      </SubscriberContext.Provider>
                    </EnodebContext.Provider>
                  </GatewayContext.Provider>
                </ApnContext.Provider>
              </PolicyContext.Provider>
            </LteNetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Network Dashboard', async () => {
    const {getByTestId, getByLabelText} = render(<Wrapper />);
    await wait();

    const info = getByTestId('info');
    expect(info).toHaveTextContent('Test Network');
    expect(info).toHaveTextContent('test_network');
    expect(info).toHaveTextContent('Test Network Description');

    const ran = getByTestId('ran');
    expect(ran).toHaveTextContent('20');
    expect(ran).toHaveTextContent('TDD');
    expect(ran).toHaveTextContent('44390');
    expect(ran).toHaveTextContent('7');
    expect(ran).toHaveTextContent('2');

    const epc = getByTestId('epc');
    expect(epc).toHaveTextContent('Enabled');
    expect(epc).toHaveTextContent('001');
    expect(epc).toHaveTextContent('01');
    expect(epc).toHaveTextContent('1');

    let epcPasswordInputElement = getByTestId('LTE Auth AMF obscure')
      .firstChild;
    if (
      epcPasswordInputElement instanceof HTMLInputElement &&
      epcPasswordInputElement.value &&
      epcPasswordInputElement.type
    ) {
      expect(epcPasswordInputElement.value).toBe('gAA=');
      expect(epcPasswordInputElement.type).toBe('password');
    } else {
      throw 'unexpected types';
    }

    fireEvent.click(getByLabelText('toggle password visibility'));
    await wait();
    epcPasswordInputElement = getByTestId('LTE Auth AMF obscure').firstChild;
    if (
      epcPasswordInputElement instanceof HTMLInputElement &&
      epcPasswordInputElement.value &&
      epcPasswordInputElement.type
    ) {
      expect(epcPasswordInputElement.type).toBe('text');
    }

    // verify KPI tray
    expect(getByTestId('Gateways')).toHaveTextContent('1');
    expect(getByTestId('eNodeBs')).toHaveTextContent('1');
    expect(getByTestId('Subscribers')).toHaveTextContent('1');
    expect(getByTestId('Policies')).toHaveTextContent('2');
    expect(getByTestId('APNs')).toHaveTextContent('2');
  });

  it('Verify Network Add', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Add Network'));
    await wait();

    // check if only first tab (network) is active
    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    let netIdField = getByTestId('networkID').firstChild;
    let netNameField = getByTestId('networkName').firstChild;
    let netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      fireEvent.change(netIdField, {target: {value: 'testNetworkID'}});
      fireEvent.change(netNameField, {target: {value: 'Test LTE Network'}});
      fireEvent.change(netDescField, {
        target: {value: 'LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(axiosMock.post).toHaveBeenCalledWith('/nms/network/create', {
      networkID: 'testNetworkID',
      data: {
        name: 'Test LTE Network',
        description: 'LTE test network description',
        networkType: 'lte',
      },
    });

    // now tab should move to epc edit component
    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).not.toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    // switch tab to network and verify editing of recently created network
    fireEvent.click(getByTestId('networkTab'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    netIdField = getByTestId('networkID').firstChild;
    netNameField = getByTestId('networkName').firstChild;
    netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      expect(netIdField.value).toBe('testNetworkID');
      // networkID shouldn't be editable
      expect(netIdField.disabled).toBe(true);
      expect(netNameField.value).toBe('Test LTE Network');
      expect(netDescField.value).toBe('LTE test network description');

      fireEvent.change(netDescField, {
        target: {value: 'New LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(MagmaAPI.lteNetworks.lteNetworkIdPut).toHaveBeenCalledWith({
      networkId: 'testNetworkID',
      lteNetwork: {
        name: 'Test LTE Network',
        description: 'New LTE test network description',
        id: 'testNetworkID',
      },
    });

    // verify adding EPC parameters
    const mncField = getByTestId('mnc').firstChild;
    const mccField = getByTestId('mcc').firstChild;
    const tacField = getByTestId('tac').firstChild;
    if (
      mncField instanceof HTMLInputElement &&
      mccField instanceof HTMLInputElement &&
      tacField instanceof HTMLInputElement
    ) {
      fireEvent.change(mncField, {target: {value: '02'}});
      fireEvent.change(mccField, {target: {value: '003'}});
      fireEvent.change(tacField, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPI.lteNetworks.lteNetworkIdCellularEpcPut,
    ).toHaveBeenCalledWith({
      config: {
        cloud_subscriberdb_enabled: false,
        default_rule_id: 'default_rule_1',
        lte_auth_amf: 'gAA=',
        lte_auth_op: 'EREREREREREREREREREREQ==',
        mcc: '003',
        mnc: '02',
        network_services: ['policy_enforcement'],
        hss_relay_enabled: false,
        gx_gy_relay_enabled: false,
        sub_profiles: {},
        mobility: {
          ip_allocation_mode: 'NAT',
          enable_static_ip_assignments: false,
          enable_multi_apn_ip_allocation: false,
        },
        tac: 1,
      },
      networkId: 'testNetworkID',
    });

    // now save and continue should move to Ran component
    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    const specialSubframePattern = getByTestId('specialSubframePattern')
      .firstChild;
    const subframeAssignment = getByTestId('subframeAssignment').firstChild;
    if (
      earfcndl instanceof HTMLElement &&
      subframeAssignment instanceof HTMLElement &&
      specialSubframePattern instanceof HTMLElement
    ) {
      fireEvent.change(earfcndl, {target: {value: '44000'}});
      fireEvent.change(specialSubframePattern, {target: {value: '8'}});
      fireEvent.change(subframeAssignment, {target: {value: '2'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Add Network'));
    await wait();
    expect(
      MagmaAPI.lteNetworks.lteNetworkIdCellularRanPut,
    ).toHaveBeenCalledWith({
      config: {
        bandwidth_mhz: 20,
        fdd_config: undefined,
        tdd_config: {
          earfcndl: 44000,
          special_subframe_pattern: 8,
          subframe_assignment: 2,
        },
      },
      networkId: 'testNetworkID',
    });
    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledTimes(0);
  });

  it('Verify Network Edit Info', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('infoEditButton'));
    await wait();

    // check if first tab (network) is active
    expect(queryByTestId('networkInfoEdit')).not.toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    const netIdField = getByTestId('networkID').firstChild;
    const netNameField = getByTestId('networkName').firstChild;
    const netDescField = getByTestId('networkDescription').firstChild;

    if (
      netIdField instanceof HTMLInputElement &&
      netNameField instanceof HTMLInputElement &&
      netDescField instanceof HTMLTextAreaElement
    ) {
      expect(netIdField.value).toBe('test_network');

      // networkID shouldn't be editable
      expect(netIdField.disabled).toBe(true);
      expect(netNameField.value).toBe('Test Network');
      expect(netDescField.value).toBe('Test Network Description');

      fireEvent.change(netDescField, {
        target: {value: 'Edit LTE test network description'},
      });
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(MagmaAPI.lteNetworks.lteNetworkIdPut).toHaveBeenCalledWith({
      networkId: 'test_network',
      lteNetwork: {
        ...testNetwork,
        description: 'Edit LTE test network description',
        cellular: {
          epc: epc,
          ran: ran,
        },
      },
    });
    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledTimes(1);
  });

  it('Verify Network Edit EPC', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('epcEditButton'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).not.toBeNull();
    expect(queryByTestId('networkRanEdit')).toBeNull();

    const mncField = getByTestId('mnc').firstChild;
    if (mncField instanceof HTMLInputElement) {
      fireEvent.change(mncField, {target: {value: '03'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    await wait();

    expect(
      MagmaAPI.lteNetworks.lteNetworkIdCellularEpcPut,
    ).toHaveBeenCalledWith({
      config: {...epc, mnc: '03'},
      networkId: 'test_network',
    });
    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledTimes(1);
  });

  it('Verify Network Edit Ran', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<Wrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(queryByTestId('networkInfoEdit')).toBeNull();
    expect(queryByTestId('networkEpcEdit')).toBeNull();
    expect(queryByTestId('networkRanEdit')).not.toBeNull();

    const earfcndl = getByTestId('earfcndl').firstChild;
    if (earfcndl instanceof HTMLElement) {
      fireEvent.change(earfcndl, {target: {value: '40000'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    await wait();

    expect(
      MagmaAPI.lteNetworks.lteNetworkIdCellularRanPut,
    ).toHaveBeenCalledWith({
      config: {
        ...ran,
        tdd_config: {
          ...ran.tdd_config,
          earfcndl: 40000,
        },
      },
      networkId: 'test_network',
    });
    expect(MagmaAPI.lteNetworks.lteNetworkIdGet).toHaveBeenCalledTimes(1);
  });
});

describe('<FEGNetworkDashboard />', () => {
  const testNetwork: FegNetwork = {
    description: 'Test Network Description',
    federation: {
      aaa_server: {},
      eap_aka: {},
      gx: {},
      gy: {},
      health: {},
      hss: {},
      s6a: {},
      served_network_ids: ['terravm2_inbound_agw_network'],
      served_nh_ids: ['terravm2_feg_network', 'terravm3_feg_network'],
      swx: {},
    },
    id: 'test_network',
    name: 'Test Network',
    dns: {
      enable_caching: false,
      local_ttl: 0,
      records: [],
    },
  };

  beforeEach(() => {
    MagmaAPIBindings.getNetworks.mockImplementation(() => Promise.resolve([]));
  });

  const Wrapper = () => {
    const networkCtx = {
      state: {
        ...testNetwork,
      },
      updateNetworks: async () => {},
    };
    return (
      <MemoryRouter
        initialEntries={['/nms/test_network/network']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <FEGNetworkContext.Provider value={networkCtx}>
              <Routes>
                <Route
                  path="/nms/:networkId/network/*"
                  element={<FEGNetworkDashboard />}
                />
              </Routes>
            </FEGNetworkContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };
  it('Verify Network Info shown correctly', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();

    const info = getByTestId('feg_info');
    expect(info).toHaveTextContent('Test Network');
    expect(info).toHaveTextContent('test_network');
    expect(info).toHaveTextContent('Test Network Description');
    expect(info).toHaveTextContent('terravm2_inbound_agw_network');
    expect(info).toHaveTextContent('terravm2_feg_network');
    expect(info).toHaveTextContent('terravm3_feg_network');
  });
});
