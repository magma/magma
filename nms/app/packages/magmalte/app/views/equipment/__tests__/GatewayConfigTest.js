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
import type {lte_gateway} from '@fbcnms/magma-api';

import 'jest-dom/extend-expect';

import AddEditGatewayButton from '../GatewayDetailConfigEdit';
import GatewayConfig from '../GatewayDetailConfig';
import GatewayContext from '../../../components/context/GatewayContext';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default.js';

import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {
  SetGatewayState,
  UpdateGateway,
} from '../../../state/lte/EquipmentState';
import {cleanup, fireEvent, render, wait} from '@testing-library/react';
import {useState} from 'react';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
afterEach(cleanup);
const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

const mockGw0: lte_gateway = {
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
    tier: 'tier2',
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
};

describe('<AddEditGatewayButton />', () => {
  afterEach(() => {
    MagmaAPIBindings.postLteByNetworkIdGateways.mockClear();
    MagmaAPIBindings.putLteByNetworkIdGatewaysByGatewayIdCellularDns.mockClear();
  });

  const AddWrapper = () => {
    const [lteGateways, setLteGateways] = useState({testGatewayId0: mockGw0});
    return (
      <MemoryRouter initialEntries={['/nms/test/gateway']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
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
                updateGateway: async () => {},
              }}>
              <Route
                path="/nms/:networkId/gateway"
                render={props => (
                  <AddEditGatewayButton
                    {...props}
                    title="Add Gateway"
                    isLink={false}
                  />
                )}
              />
            </GatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  const DetailWrapper = () => {
    const [lteGateways, setLteGateways] = useState({testGatewayId0: mockGw0});
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
                  }),
              }}>
              <Route
                path="/nms/:networkId/gateway/:gatewayId/config"
                render={props => <GatewayConfig {...props} />}
              />
            </GatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('Verify Gateway Add', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<AddWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByText('Add Gateway'));
    await wait();

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('aggregationEdit')).toBeNull();
    expect(queryByTestId('ranEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();

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
    expect(MagmaAPIBindings.postLteByNetworkIdGateways).toHaveBeenCalledWith({
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
          dynamic_services: [],
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

  it('Verify Gateway Ran Edit', async () => {
    const {getByTestId, getByText, queryByTestId} = render(<DetailWrapper />);
    await wait();
    expect(queryByTestId('editDialog')).toBeNull();

    fireEvent.click(getByTestId('ranEditButton'));
    await wait();

    expect(queryByTestId('infoEdit')).toBeNull();
    expect(queryByTestId('epcEdit')).toBeNull();
    expect(queryByTestId('aggregationEdit')).toBeNull();
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
      MagmaAPIBindings.putLteByNetworkIdGatewaysByGatewayIdCellularDns,
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
