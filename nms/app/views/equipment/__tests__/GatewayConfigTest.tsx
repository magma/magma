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
import ApnContext from '../../../context/ApnContext';
import GatewayConfig from '../GatewayDetailConfig';
import LteNetworkContext from '../../../context/LteNetworkContext';
import MagmaAPI from '../../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {DynamicServices} from '../../../components/GatewayUtils';
import {GatewayContextProvider} from '../../../context/GatewayContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, render, waitFor, within} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import type {Apn, LteGateway, LteNetwork} from '../../../../generated';

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
    ngc: {
      amf_default_sd: 'AFAF',
      amf_default_sst: 2,
      amf_name: 'amf.example.org',
      amf_pointer: '1F',
      amf_region_id: 'C1',
      amf_set_id: '2A1',
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
    mockAPI(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGet', {
      testGatewayId0: mockGw0,
    });
  });

  const AddWrapper = () => {
    return (
      <MemoryRouter initialEntries={['/nms/test/gateway']} initialIndex={0}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={defaultTheme}>
            <ThemeProvider theme={defaultTheme}>
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
                  <GatewayContextProvider networkId="test">
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
                  </GatewayContextProvider>
                </LteNetworkContext.Provider>
              </ApnContext.Provider>
            </ThemeProvider>
          </ThemeProvider>
        </StyledEngineProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    return (
      <MemoryRouter
        initialEntries={['/nms/test/gateway/testGatewayId0/config']}
        initialIndex={0}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={defaultTheme}>
            <ThemeProvider theme={defaultTheme}>
              <GatewayContextProvider networkId={'test'}>
                <Routes>
                  <Route
                    path="/nms/:networkId/gateway/:gatewayId/config"
                    element={<GatewayConfig />}
                  />
                </Routes>
              </GatewayContextProvider>
            </ThemeProvider>
          </ThemeProvider>
        </StyledEngineProvider>
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

    const {
      getByTestId,
      getByText,
      queryByTestId,
      findByTestId,
      findByText,
    } = render(<AddWrapper />);
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(await findByText('Add Gateway'));

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    // test adding an existing gateway
    changeInput('id', 'testGatewayId0', getByTestId);

    fireEvent.click(getByText('Save And Continue'));
    expect(await findByTestId('configEditError')).toHaveTextContent(
      'Gateway testGatewayId0 already exists',
    );

    changeInput('id', 'testGatewayID1', getByTestId);
    changeInput('name', 'testGatewayName', getByTestId);
    changeInput('hardwareId', 'testHwId', getByTestId);
    changeInput('challengeKey', 'testChallenge', getByTestId);
    changeInput('description', 'Test Gateway Description', getByTestId);
    changeInput('version', '1.0', getByTestId);

    fireEvent.click(getByText('Save And Continue'));
    await waitFor(() => {
      expect(
        MagmaAPI.lteGateways.lteNetworkIdGatewaysPost,
      ).toHaveBeenCalledWith({
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
    if (monitordService instanceof HTMLInputElement) {
      fireEvent.click(monitordService);
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await waitFor(() => {
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
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).not.toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();

    // Verify EPC Edit
    const natEnabled = getByTestId('natEnabled').firstChild;
    if (natEnabled instanceof HTMLInputElement) {
      fireEvent.click(natEnabled);
    } else {
      throw 'invalid type';
    }
    changeInput('gwSgiIpv6', '2001:4860:4860:0:0:0:0:1', getByTestId);
    changeInput('sgiStaticIpv6', '2001:4860:4860:0:0:0:0:8888', getByTestId);
    changeInput('ipv6Block', 'fdee:5:6c::/48', getByTestId);
    fireEvent.click(getByText('Save And Continue'));

    await waitFor(() => {
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
    if (enbDhcpService instanceof HTMLInputElement) {
      fireEvent.click(enbDhcpService);
    } else {
      throw 'invalid type';
    }

    pci = (await findByTestId('pci')).firstChild;
    expect(pci).toBeDisabled();

    const registeredEnodeb = getByTestId('registeredEnodeb').firstChild;
    expect(registeredEnodeb).not.toHaveAttribute('aria-disabled');

    fireEvent.click(getByText('Save And Continue'));
    await waitFor(() => {
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
    changeInput('apnID', '1', getByTestId);
    changeInput('vlanID', '1', getByTestId);

    fireEvent.click(getByText('Save And Continue'));

    await waitFor(() => {
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
    });

    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('apnResourcesEdit')).toBeNull();
    expect(queryByTestId('headerEnrichmentEdit')).not.toBeNull();
    // Verify Header Enrichment Edit
    const HeEnabled = getByTestId('enableHE').firstChild;
    if (HeEnabled instanceof HTMLInputElement) {
      fireEvent.click(HeEnabled);
    } else {
      throw 'invalid type';
    }
    expect(queryByTestId('encryptionEdit')).toBeNull();
    const encryptionEnabled = getByTestId('enableEncryption').firstChild;
    if (encryptionEnabled instanceof HTMLInputElement) {
      fireEvent.click(encryptionEnabled);
    } else {
      throw 'invalid type';
    }
    // Encryption fields are visible if encryption is enabled
    expect(await findByTestId('encryptionEdit')).not.toBeNull();

    fireEvent.click(getByText('Save And Continue'));

    await waitFor(() => {
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

    changeInput('amfName', 'myamf.test', getByTestId);
    changeInput('amfPointer', '2C', getByTestId);
    changeInput('amfRegionID', 'D5', getByTestId);
    changeInput('amfSetID', '1CA', getByTestId);
    changeInput('amfDefaultSliceServiceType', '42', getByTestId);
    changeInput('amfDefaultSliceDescriptor', 'ABCD', getByTestId);

    fireEvent.click(getByText('Save And Close'));

    await waitFor(() => {
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
          ngc: {
            amf_name: 'myamf.test',
            amf_pointer: '2C',
            amf_region_id: 'D5',
            amf_set_id: '1CA',
            amf_default_sst: 42,
            amf_default_sd: 'ABCD',
          },
        },
        gatewayId: 'testGatewayID1',
        networkId: 'test',
      });
    });
  });

  it('Verify Gateway Ran Edit', async () => {
    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularDnsPut',
      )
      .mockImplementation();

    const {getByTestId, getByText, queryByTestId, findByTestId} = render(
      <DetailWrapper />,
    );
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(await findByTestId('ranEditButton'));

    expect(await findByTestId('ranEdit')).not.toBeNull();
    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('dynamicServicesEdit')).toBeNull();

    let pci = getByTestId('pci').firstChild;
    if (pci instanceof HTMLInputElement) {
      expect(pci.disabled).toBe(false);
    } else {
      throw 'invalid type';
    }

    const enbDhcpService = getByTestId('enbDhcpService').firstChild;
    if (enbDhcpService instanceof HTMLInputElement) {
      fireEvent.click(enbDhcpService);
    } else {
      throw 'invalid type';
    }

    pci = (await findByTestId('pci')).firstChild;
    if (pci instanceof HTMLInputElement) {
      expect(pci.disabled).toBe(true);
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save'));
    await waitFor(() => {
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

  it('Verify Gateway NGC Config', async () => {
    jest
      .spyOn(MagmaAPI.lteGateways, 'lteNetworkIdGatewaysGatewayIdCellularPut')
      .mockImplementation();

    const {findByTestId, getByTestId, getByText} = render(<DetailWrapper />);
    const ngcConfig = await findByTestId('ngc-config');

    const nameCell = within(ngcConfig).getByTestId('Name');
    within(nameCell).getByText('amf.example.org');
    const pointerCell = within(ngcConfig).getByTestId('Pointer');
    within(pointerCell).getByText('1F');
    const regionCell = within(ngcConfig).getByTestId('Region ID');
    within(regionCell).getByText('C1');
    const setCell = within(ngcConfig).getByTestId('Set ID');
    within(setCell).getByText('2A1');
    const serviceTypeCell = within(ngcConfig).getByTestId(
      'Default Slice Service Type',
    );
    within(serviceTypeCell).getByText('2');
    const descriptorCell = within(ngcConfig).getByTestId(
      'Default Slice Descriptor',
    );
    within(descriptorCell).getByText('AFAF');

    fireEvent.click(await findByTestId('ngcEditButton'));

    changeInput('amfName', 'test.com', getByTestId);
    changeInput('amfPointer', '', getByTestId);
    changeInput('amfRegionID', '', getByTestId);
    changeInput('amfSetID', '', getByTestId);
    changeInput('amfDefaultSliceServiceType', '12', getByTestId);
    changeInput('amfDefaultSliceDescriptor', '', getByTestId);

    fireEvent.click(getByText('Save'));

    await waitFor(() => {
      expect(
        MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPut,
      ).toHaveBeenCalledWith({
        config: {
          ...mockGw0.cellular,
          ngc: {
            amf_name: 'test.com',
            amf_pointer: undefined,
            amf_region_id: undefined,
            amf_set_id: undefined,
            amf_default_sst: 12,
            amf_default_sd: undefined,
          },
        },
        gatewayId: mockGw0.id,
        networkId: 'test',
      });
    });
  });
});

function changeInput(
  testID: string,
  newValue: string,
  getByTestId: (testID: string) => HTMLElement,
) {
  const input = getByTestId(testID).firstChild;

  if (input instanceof HTMLInputElement) {
    fireEvent.change(input, {target: {value: newValue}});
  } else {
    throw 'invalid type';
  }
}
