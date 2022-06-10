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

import AddEditGatewayButton from '../GatewayDetailConfigEdit';
import ApnContext from '../../../components/context/ApnContext';
import GatewayConfig from '../GatewayDetailConfig';
import GatewayContext from '../../../components/context/GatewayContext';
import LteNetworkContext from '../../../components/context/LteNetworkContext';
import MagmaAPI from '../../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {DynamicServices} from '../../../components/GatewayUtils';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {
  SetGatewayState,
  UpdateGateway,
  UpdateGatewayProps,
} from '../../../state/lte/EquipmentState';
import {fireEvent, render, wait} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import {useState} from 'react';
import type {Apn, LteGateway, LteNetwork} from '../../../../generated-ts';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

const mockGw0: LteGateway = {
  apn_resources: {},
  id: ' testGatewayId0',
  name: ' testGateway0',
  description: ' testGateway0',
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
  connected_enodeb_serials: [],
  cellular: {
    epc: {
      ip_block: '192.168.0.1/24',
      nat_enabled: false,
      sgi_management_iface_static_ip: '1.1.1.1/24',
      sgi_management_iface_vlan: '100',
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
  checked_in_recently: false,
};

const mockGw1: LteGateway = {
  apn_resources: {},
  id: ' testGatewayId1',
  name: ' testGateway1',
  description: ' testGateway1',
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
  connected_enodeb_serials: [],
  cellular: {
    epc: {
      ip_block: '192.168.0.1/24',
      nat_enabled: false,
      sgi_management_iface_static_ip: '1.1.1.1/24',
      sgi_management_iface_vlan: '100',
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
  checked_in_recently: false,
};

const mockNw: LteNetwork = {
  cellular: {
    epc: {
      default_rule_id: 'default_rule_1',
      gx_gy_relay_enabled: false,
      hss_relay_enabled: false,
      lte_auth_amf: 'gAA=',
      lte_auth_op: '=',
      mcc: '001',
      mnc: '01',
      mobility: {
        ip_allocation_mode: 'DHCP_BROADCAST',
        enable_multi_apn_ip_allocation: true,
      },
      network_services: ['policy_enforcement'],
      tac: 1,
    },
    ran: {
      bandwidth_mhz: 20,
      tdd_config: {
        earfcndl: 44590,
        special_subframe_pattern: 7,
        subframe_assignment: 2,
      },
    },
  },
  description: 'magma appliance',
  dns: {
    enable_caching: false,
    local_ttl: 0,
  },
  features: {
    features: {
      placeholder: 'true',
    },
  },
  id: '1dev_agw',
  name: '1dev_agw',
};

const mockApns: Record<string, Apn> = {
  'oai.ipv4': {
    apn_configuration: {
      ambr: {max_bandwidth_dl: 1000000, max_bandwidth_ul: 1000000},
      qos_profile: {
        class_id: 9,
        preemption_capability: false,
        preemption_vulnerability: false,
        priority_level: 15,
      },
    },
    apn_name: 'oai.ipv4',
  },
};

describe('<AddEditGatewayButton />', () => {
  beforeEach(() => {
    (useEnqueueSnackbar as jest.Mock).mockReturnValue(jest.fn());
  });

  const AddWrapper = () => {
    const [lteGateways, setLteGateways] = useState<Record<string, LteGateway>>({
      testGatewayId0: mockGw0,
    });
    return (
      <MemoryRouter initialEntries={['/nms/test/gateway']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <ApnContext.Provider
              value={{
                state: mockApns,
                setState: async () => {},
              }}>
              <LteNetworkContext.Provider
                value={{
                  state: mockNw,
                  updateNetworks: async () => {},
                }}>
                <GatewayContext.Provider
                  value={{
                    state: lteGateways,
                    setState: async (key, value?) =>
                      SetGatewayState({
                        lteGateways: lteGateways,
                        setLteGateways: setLteGateways,
                        networkId: 'test',
                        key: key,
                        value: value,
                      }),
                    updateGateway: props =>
                      UpdateGateway({
                        networkId: 'test',
                        setLteGateways,
                        ...props,
                      } as UpdateGatewayProps),
                  }}>
                  <Routes>
                    <Route
                      path="/nms/:networkId/gateway"
                      element={
                        <AddEditGatewayButton
                          title="Add Gateway"
                          isLink={false}
                        />
                      }
                    />
                  </Routes>
                </GatewayContext.Provider>
              </LteNetworkContext.Provider>
            </ApnContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    const [lteGateways, setLteGateways] = useState<Record<string, LteGateway>>({
      testGatewayId0: mockGw0,
    });
    return (
      <MemoryRouter
        initialEntries={['/nms/test/gateway/testGatewayId0/config']}
        initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <GatewayContext.Provider
              value={{
                state: lteGateways,
                setState: async () => {},
                updateGateway: props =>
                  UpdateGateway({
                    networkId: 'test',
                    setLteGateways: setLteGateways,
                    ...props,
                  } as UpdateGatewayProps),
              }}>
              <Routes>
                <Route
                  path="/nms/:networkId/gateway/:gatewayId/config"
                  element={<GatewayConfig />}
                />
              </Routes>
            </GatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Gateway Add', async () => {
    jest
      .spyOn(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysPost')
      .mockImplementation();
    jest
      .spyOn(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGatewayIdMagmadPut')
      .mockImplementation();
    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularEpcPut',
      )
      .mockImplementation();
    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularDnsPut',
      )
      .mockImplementation();
    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularRanPut',
      )
      .mockImplementation();
    jest
      .spyOn(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGatewayIdPut')
      .mockImplementation();
    jest
      .spyOn(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGatewayIdCellularPut')
      .mockImplementation();

    const {getByTestId, getByText, queryByTestId} = render(<AddWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Add Gateway'));
    await wait();

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    const gatewayID = getByTestId('id').firstChild;
    const gatewayName = getByTestId('name').firstChild;
    const hwId = getByTestId('hardwareId').firstChild;
    const version = getByTestId('version').firstChild;
    const description = getByTestId('description').firstChild;
    const challengeKey = getByTestId('challengeKey').firstChild;

    // test adding an existing gateway
    if (gatewayID instanceof HTMLInputElement) {
      fireEvent.change(gatewayID, {target: {value: 'testGatewayId0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(getByTestId('configEditError')).toHaveTextContent(
      'Gateway testGatewayId0 already exists',
    );

    if (
      gatewayID instanceof HTMLInputElement &&
      gatewayName instanceof HTMLInputElement &&
      hwId instanceof HTMLInputElement &&
      version instanceof HTMLInputElement &&
      challengeKey instanceof HTMLInputElement &&
      description instanceof HTMLInputElement
    ) {
      fireEvent.change(gatewayID, {target: {value: 'testGatewayID1'}});
      fireEvent.change(gatewayName, {target: {value: 'testGatewayName'}});
      fireEvent.change(description, {
        target: {value: 'Test Gateway Description'},
      });
      fireEvent.change(challengeKey, {target: {value: 'testChallenge'}});
      fireEvent.change(hwId, {target: {value: 'testHwId'}});
      fireEvent.change(version, {target: {value: '1.0'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(MagmaAPI.lteGateways.lteNetworkIdGatewaysPost).toHaveBeenCalledWith({
      gateway: {
        apn_resources: {},
        id: 'testGatewayID1',
        name: 'testGatewayName',
        cellular: {
          epc: {
            dns_primary: '',
            dns_secondary: '',
            ip_block: '192.168.128.0/24',
            nat_enabled: true,
            sgi_management_iface_gw: '',
            sgi_management_iface_static_ip: '',
            sgi_management_iface_vlan: '',
          },
          ran: {
            pci: 260,
            transmit_enabled: true,
          },
        },
        checked_in_recently: false,
        connected_enodeb_serials: [],
        description: 'Test Gateway Description',
        device: {
          hardware_id: 'testHwId',
          key: {
            key: 'testChallenge',
            key_type: 'SOFTWARE_ECDSA_SHA256',
          },
        },

        magmad: {
          autoupgrade_enabled: true,
          autoupgrade_poll_interval: 60,
          checkin_interval: 60,
          checkin_timeout: 30,
          dynamic_services: [
            DynamicServices.EVENTD,
            DynamicServices.TD_AGENT_BIT,
          ],
        },
        status: {
          platform_info: {
            packages: [
              {
                version: '1.0',
              },
            ],
          },
        },
        tier: 'default',
      },
      networkId: 'test',
    });

    // mock adding test gatewayID1 to ensure we invoke the put method
    mockAPI(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGet', {
      testGatewayID1: mockGw1,
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    // Verify Dynamic Services Edit
    const monitordService = getByTestId('monitordService').firstChild;
    if (
      monitordService instanceof HTMLElement &&
      monitordService.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(monitordService.firstChild);
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdMagmadPut,
    ).toHaveBeenCalledWith({
      gatewayId: 'testGatewayID1',
      networkId: 'test',
      magmad: {
        autoupgrade_enabled: true,
        autoupgrade_poll_interval: 60,
        checkin_interval: 60,
        checkin_timeout: 30,
        dynamic_services: [
          DynamicServices.EVENTD,
          DynamicServices.TD_AGENT_BIT,
          DynamicServices.MONITORD,
        ],
        logging: {
          aggregation: {
            target_files_by_tag: {
              mme: 'var/log/mme.log',
            },
          },
          log_level: 'DEBUG',
        },
      },
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    // Verify EPC Edit
    const natEnabled = getByTestId('natEnabled').firstChild;
    const gwSgiIpv6 = getByTestId('gwSgiIpv6').firstChild;
    const sgiStaticIpv6 = getByTestId('sgiStaticIpv6').firstChild;
    const ipv6Block = getByTestId('ipv6Block').firstChild;
    if (
      natEnabled instanceof HTMLElement &&
      natEnabled.firstChild instanceof HTMLElement &&
      gwSgiIpv6 instanceof HTMLInputElement &&
      sgiStaticIpv6 instanceof HTMLInputElement &&
      ipv6Block instanceof HTMLInputElement
    ) {
      fireEvent.click(natEnabled.firstChild);
      fireEvent.change(gwSgiIpv6, {
        target: {value: '2001:4860:4860:0:0:0:0:1'},
      });
      fireEvent.change(sgiStaticIpv6, {
        target: {value: '2001:4860:4860:0:0:0:0:8888'},
      });
      fireEvent.change(ipv6Block, {
        target: {value: 'fdee:5:6c::/48'},
      });
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularEpcPut,
    ).toHaveBeenCalledWith({
      gatewayId: 'testGatewayID1',
      networkId: 'test',
      config: {
        ip_block: '192.168.128.0/24',
        ipv6_block: 'fdee:5:6c::/48',
        nat_enabled: false,
        dns_primary: '',
        dns_secondary: '',
        sgi_management_iface_gw: '',
        sgi_management_iface_static_ip: '',
        sgi_management_iface_vlan: '',
        sgi_management_iface_ipv6_gw: '2001:4860:4860:0:0:0:0:1',
        sgi_management_iface_ipv6_addr: '2001:4860:4860:0:0:0:0:8888',
      },
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    // Verify RAN Edit
    let pci = getByTestId('pci').firstChild;
    if (pci instanceof HTMLInputElement) {
      expect(pci.disabled).toBe(false);
    } else {
      throw 'invalid type';
    }

    const enbDhcpService = getByTestId('enbDhcpService').firstChild;
    if (
      enbDhcpService instanceof HTMLElement &&
      enbDhcpService.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(enbDhcpService.firstChild);
    } else {
      throw 'invalid type';
    }
    await wait();

    pci = getByTestId('pci').firstChild;
    expect(pci).toBeDisabled();

    const registeredEnodeb = getByTestId('registeredEnodeb').firstChild;
    expect(registeredEnodeb).not.toHaveAttribute('aria-disabled');

    fireEvent.click(getByText('Save And Continue'));
    await wait();
    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularDnsPut,
    ).toHaveBeenCalledWith({
      config: {
        dhcp_server_enabled: false,
        enable_caching: false,
        local_ttl: 0,
        records: [],
      },
      gatewayId: 'testGatewayID1',
      networkId: 'test',
    });
    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularRanPut,
    ).toHaveBeenCalledWith({
      config: {
        pci: 260,
        transmit_enabled: true,
      },
      gatewayId: 'testGatewayID1',
      networkId: 'test',
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).not.toBeNull();

    // Verify Apn Resources Edit
    expect(queryByTestId('apnResourcesAdd')).not.toBeNull();
    const apnResourcesAdd = queryByTestId('apnResourcesAdd');
    if (!apnResourcesAdd) {
      throw new Error('apn resources add button unexpected null');
    }
    fireEvent.click(apnResourcesAdd);
    const apnID = getByTestId('apnID').firstChild;
    const vlanID = getByTestId('vlanID').firstChild;

    if (apnID instanceof HTMLInputElement) {
      fireEvent.change(apnID, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }

    if (vlanID instanceof HTMLInputElement) {
      fireEvent.change(vlanID, {target: {value: '1'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdPut,
    ).toHaveBeenCalledWith({
      gateway: {
        apn_resources: {'': {apn_name: '', id: '1', vlan_id: 1}},
        cellular: {
          dns: {
            dhcp_server_enabled: false,
            enable_caching: false,
            local_ttl: 0,
            records: [],
          },
          epc: {
            ip_block: '192.168.128.0/24',
            ipv6_block: 'fdee:5:6c::/48',
            nat_enabled: false,
            dns_primary: '',
            dns_secondary: '',
            sgi_management_iface_gw: '',
            sgi_management_iface_static_ip: '',
            sgi_management_iface_vlan: '',
            sgi_management_iface_ipv6_gw: '2001:4860:4860:0:0:0:0:1',
            sgi_management_iface_ipv6_addr: '2001:4860:4860:0:0:0:0:8888',
          },
          ran: {pci: 260, transmit_enabled: true},
        },
        checked_in_recently: false,
        connected_enodeb_serials: [],
        description: 'Test Gateway Description',
        device: {
          hardware_id: 'testHwId',
          key: {key: 'testChallenge', key_type: 'SOFTWARE_ECDSA_SHA256'},
        },
        id: 'testGatewayID1',
        magmad: {
          autoupgrade_enabled: true,
          autoupgrade_poll_interval: 60,
          checkin_interval: 60,
          checkin_timeout: 30,
          dynamic_services: [
            DynamicServices.EVENTD,
            DynamicServices.TD_AGENT_BIT,
            DynamicServices.MONITORD,
          ],
          logging: {
            aggregation: {
              target_files_by_tag: {
                mme: 'var/log/mme.log',
              },
            },
            log_level: 'DEBUG',
          },
        },
        name: 'testGatewayName',

        status: {
          platform_info: {
            packages: [
              {
                version: '1.0',
              },
            ],
          },
        },
        tier: 'default',
      },
      gatewayId: 'testGatewayID1',
      networkId: 'test',
    });
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).not.toBeNull();
    // Verify Header Enrichment Edit
    const HeEnabled = getByTestId('enableHE').firstChild;
    if (
      HeEnabled instanceof HTMLElement &&
      HeEnabled.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(HeEnabled.firstChild);
    } else {
      throw 'invalid type';
    }
    expect(queryByTestId('encryptionEdit')).toBeNull();
    const encryptionEnabled = getByTestId('enableEncryption').firstChild;
    if (
      encryptionEnabled instanceof HTMLElement &&
      encryptionEnabled.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(encryptionEnabled.firstChild);
    } else {
      throw 'invalid type';
    }
    await wait();
    // Encryption fields are visible if encryption is enabled
    expect(queryByTestId('encryptionEdit')).not.toBeNull();

    fireEvent.click(getByText('Save And Close'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPut,
    ).toHaveBeenCalledWith({
      config: {
        dns: {
          dhcp_server_enabled: false,
          enable_caching: false,
          local_ttl: 0,
          records: [],
        },
        epc: {
          ip_block: '192.168.128.0/24',
          ipv6_block: 'fdee:5:6c::/48',
          nat_enabled: false,
          dns_primary: '',
          dns_secondary: '',
          sgi_management_iface_gw: '',
          sgi_management_iface_static_ip: '',
          sgi_management_iface_vlan: '',
          sgi_management_iface_ipv6_gw: '2001:4860:4860:0:0:0:0:1',
          sgi_management_iface_ipv6_addr: '2001:4860:4860:0:0:0:0:8888',
        },
        ran: {
          pci: 260,
          transmit_enabled: true,
        },
        he_config: {
          enable_encryption: true,
          enable_header_enrichment: true,
          he_encoding_type: 'BASE64',
          he_encryption_algorithm: 'RC4',
          he_hash_function: 'MD5',
          encryption_key: '',
        },
      },
      gatewayId: 'testGatewayID1',
      networkId: 'test',
    });
  });

  it('Verify Gateway Ran Edit', async () => {
    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularDnsPut',
      )
      .mockImplementation();

    const {getByTestId, getByText, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).not.toBeNull();

    let pci = getByTestId('pci').firstChild;
    if (pci instanceof HTMLInputElement) {
      expect(pci.disabled).toBe(false);
    } else {
      throw 'invalid type';
    }

    const enbDhcpService = getByTestId('enbDhcpService').firstChild;
    if (
      enbDhcpService instanceof HTMLElement &&
      enbDhcpService.firstChild instanceof HTMLElement
    ) {
      fireEvent.click(enbDhcpService.firstChild);
    } else {
      throw 'invalid type';
    }
    await wait();

    pci = getByTestId('pci').firstChild;
    if (pci instanceof HTMLInputElement) {
      expect(pci.disabled).toBe(true);
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await wait();
    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularDnsPut,
    ).toHaveBeenCalledWith({
      config: {
        dhcp_server_enabled: false,
        enable_caching: false,
        local_ttl: 0,
        records: [],
      },
      gatewayId: ' testGatewayId0',
      networkId: 'test',
    });
  });
});
