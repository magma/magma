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
import type {lte_gateway} from '../../../../generated/MagmaAPIBindings';

// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import AddEditGatewayPoolButton from '../GatewayPoolEdit';
// $FlowFixMe migrated to typescript
import GatewayContext from '../../../components/context/GatewayContext';
import GatewayPools from '../EquipmentGatewayPools';
// $FlowFixMe migrated to typescript
import GatewayPoolsContext from '../../../components/context/GatewayPoolsContext';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import defaultTheme from '../../../theme/default';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import MagmaAPI from '../../../../api/MagmaAPI';

import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {
  SetGatewayPoolsState,
  UpdateGatewayPoolRecords,
  // $FlowFixMe migrated to typescript
} from '../../../state/lte/EquipmentState';
import {fireEvent, render, wait} from '@testing-library/react';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import {useState} from 'react';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

const gwPoolStateMock = {
  pool1: {
    gatewayPool: {
      config: {mme_group_id: 1},
      gateway_ids: ['gw1', 'gw2', 'gw5'],
      gateway_pool_id: 'pool1',
      gateway_pool_name: 'pool_1',
    },
    gatewayPoolRecords: [
      {
        gateway_pool_id: 'pool1',
        gateway_id: 'gw1',
        mme_code: 1,
        mme_relative_capacity: 255,
      },
      {
        gateway_pool_id: 'pool1',
        gateway_id: 'gw2',
        mme_code: 2,
        mme_relative_capacity: 255,
      },
      {
        gateway_pool_id: 'pool1',
        gateway_id: 'gw5',
        mme_code: 3,
        mme_relative_capacity: 1,
      },
    ],
  },
  pool2: {
    gatewayPool: {
      config: {mme_group_id: 2},
      gateway_ids: ['gw3', 'gw4'],
      gateway_pool_id: 'pool2',
      gateway_pool_name: 'pool_2',
    },
    gatewayPoolRecords: [
      {
        gateway_pool_id: 'pool2',
        gateway_id: 'gw3',
        mme_code: 1,
        mme_relative_capacity: 255,
      },
      {
        gateway_pool_id: 'pool2',
        gateway_id: 'gw4',
        mme_code: 2,
        mme_relative_capacity: 1,
      },
    ],
  },
  pool3: {
    gatewayPool: {
      config: {mme_group_id: 3},
      gateway_ids: [],
      gateway_pool_id: 'pool3',
      gateway_pool_name: 'pool_3',
    },
    gatewayPoolRecords: [],
  },
};
const networkId = 'test';
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
  checked_in_recently: false,
};

const mockGw1 = Object.assign({}, mockGw0);
mockGw1.id = 'testGatewayId1';
mockGw1.name = 'testGateway1';

const lteGateways = {
  testGatewayId0: mockGw0,
  testGatewayId1: mockGw1,
};
describe('<GatewayPools />', () => {
  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/pools']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <GatewayPoolsContext.Provider
            value={{
              state: gwPoolStateMock,
              setState: (key, value?) =>
                SetGatewayPoolsState({
                  networkId,
                  gatewayPools: gwPoolStateMock,
                  setGatewayPools: () => {},
                  key,
                  value,
                }),
              updateGatewayPoolRecords: async _ => {},
            }}>
            <Routes>
              <Route path="/nms/:networkId/pools/" element={<GatewayPools />} />
            </Routes>
          </GatewayPoolsContext.Provider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    jest
      .spyOn(
        MagmaAPI.lteNetworks,
        'lteNetworkIdGatewayPoolsGatewayPoolIdDelete',
      )
      .mockImplementation();

    const {getByTestId, getAllByRole, getByText, getAllByTitle} = render(
      <Wrapper />,
    );
    await wait();

    const rowItems = await getAllByRole('row');

    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('ID');
    expect(rowItems[0]).toHaveTextContent('MME Group ID');
    expect(rowItems[0]).toHaveTextContent('Primary Gateway');
    expect(rowItems[0]).toHaveTextContent('Secondary Gateway');

    expect(rowItems[1]).toHaveTextContent('pool_1');
    expect(rowItems[1]).toHaveTextContent('pool1');
    expect(rowItems[1]).toHaveTextContent('1');
    expect(rowItems[1]).toHaveTextContent('gw1gw2');
    expect(rowItems[1]).toHaveTextContent('gw5');

    expect(rowItems[2]).toHaveTextContent('pool_2');
    expect(rowItems[2]).toHaveTextContent('pool2');
    expect(rowItems[2]).toHaveTextContent('2');
    expect(rowItems[2]).toHaveTextContent('gw3');
    expect(rowItems[2]).toHaveTextContent('gw4');

    expect(rowItems[3]).toHaveTextContent('pool_3');
    expect(rowItems[3]).toHaveTextContent('pool3');
    expect(rowItems[3]).toHaveTextContent('3');
    expect(rowItems[3]).toHaveTextContent('-');
    expect(rowItems[3]).toHaveTextContent('-');

    // delete gateway pool3
    const actionList = getAllByTitle('Actions');
    expect(getByTestId('actions-menu')).not.toBeVisible();
    fireEvent.click(actionList[2]);
    await wait();
    expect(getByTestId('actions-menu')).toBeVisible();
    fireEvent.click(getByText('Remove'));
    await wait();
    expect(
      getByText('Are you sure you want to delete pool3?'),
    ).toBeInTheDocument();
    fireEvent.click(getByText('Confirm'));
    await wait();

    expect(
      MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdDelete,
    ).toHaveBeenCalledWith({
      networkId: 'test',
      gatewayPoolId: 'pool3',
    });
  });
});

describe('<AddEditGatewayPoolButton />', () => {
  let lteNetworkIdGatewayPoolsGatewayPoolIdGet;

  beforeEach(() => {
    (useEnqueueSnackbar: JestMockFn<
      Array<empty>,
      $Call<typeof useEnqueueSnackbar>,
    >).mockReturnValue(jest.fn());
    jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdGatewayPoolsPost')
      .mockImplementation();

    jest
      .spyOn(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularPoolingPut',
      )
      .mockResolvedValue({data: undefined}); // TODO[TS-migration] What is the real response?

    lteNetworkIdGatewayPoolsGatewayPoolIdGet = jest
      .spyOn(MagmaAPI.lteNetworks, 'lteNetworkIdGatewayPoolsGatewayPoolIdGet')
      .mockResolvedValue({
        data: {
          config: {mme_group_id: 4},
          gateway_ids: [],
          gateway_pool_id: 'pool4',
          gateway_pool_name: 'pool4',
        },
      });
  });
  const AddWrapper = () => {
    const [gwPoolsState, setGatewayPoolsState] = useState({});
    return (
      <MemoryRouter initialEntries={['/nms/test/pools']} initialIndex={0}>
        <MuiThemeProvider theme={defaultTheme}>
          <MuiStylesThemeProvider theme={defaultTheme}>
            <GatewayContext.Provider
              value={{
                state: lteGateways,
                setState: async () => {},
                updateGateway: async () => {},
              }}>
              <GatewayPoolsContext.Provider
                value={{
                  state: {},
                  setState: async (key, value?, resources?) =>
                    SetGatewayPoolsState({
                      networkId,
                      gatewayPools: gwPoolsState,
                      setGatewayPools: setGatewayPoolsState,
                      key,
                      value,
                      resources,
                    }),
                  updateGatewayPoolRecords: (key, value?, resources) =>
                    UpdateGatewayPoolRecords({
                      networkId,
                      gatewayPools: gwPoolsState,
                      setGatewayPools: setGatewayPoolsState,
                      key,
                      value,
                      resources,
                    }),
                }}>
                <Routes>
                  <Route
                    path="/nms/:networkId/pools"
                    element={
                      <AddEditGatewayPoolButton
                        title="Add Gateway Pool"
                        isLink={false}
                      />
                    }
                  />
                </Routes>
              </GatewayPoolsContext.Provider>
            </GatewayContext.Provider>
          </MuiStylesThemeProvider>
        </MuiThemeProvider>
      </MemoryRouter>
    );
  };

  it('verify gateway pool add', async () => {
    const networkId = 'test';
    const {queryByTestId, getByTestId, getByText} = render(<AddWrapper />);
    await wait();

    expect(queryByTestId('gatewayPoolEditDialog')).toBeNull();
    fireEvent.click(getByText('Add Gateway Pool'));
    await wait();

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('PrimaryEdit')).toBeNull();
    expect(queryByTestId('SecondaryEdit')).toBeNull();

    const name = getByTestId('name').firstChild;
    const poolId = getByTestId('poolId').firstChild;
    const mmeGroupId = getByTestId('mmeGroupId').firstChild;

    if (
      name instanceof HTMLInputElement &&
      poolId instanceof HTMLInputElement &&
      mmeGroupId instanceof HTMLInputElement
    ) {
      fireEvent.change(name, {target: {value: 'pool_4'}});
      fireEvent.change(poolId, {target: {value: 'pool4'}});
      fireEvent.change(mmeGroupId, {target: {value: '4'}});
    } else {
      throw 'invalid type';
    }

    fireEvent.click(getByText('Save And Continue'));
    await wait();

    const newGatewayPool = {
      config: {mme_group_id: 4},
      gateway_pool_id: 'pool4',
      gateway_pool_name: 'pool_4',
    };

    expect(
      MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsPost,
    ).toHaveBeenCalledWith({
      networkId,
      hAGatewayPool: newGatewayPool,
    });

    await wait();

    // $FlowFixMe
    lteNetworkIdGatewayPoolsGatewayPoolIdGet.mockResolvedValue({
      data: {
        config: {mme_group_id: 4},
        gateway_ids: [],
        gateway_pool_id: 'pool4',
        gateway_pool_name: 'pool4',
      },
    });

    await wait();
    // check if only second tab (PrimaryEdit) is active
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('PrimaryEdit')).not.toBeNull();
    expect(queryByTestId('SecondaryEdit')).toBeNull();

    const mmeCode = getByTestId('mmeCode').firstChild;
    const PrimaryId = getByTestId('gwIdPrimary').firstChild;

    if (
      mmeCode instanceof HTMLInputElement &&
      PrimaryId instanceof HTMLElement
    ) {
      fireEvent.mouseDown(PrimaryId);
      await wait();
      fireEvent.click(getByText('testGatewayId0'));
      fireEvent.change(mmeCode, {target: {value: '4'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save And Continue'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPoolingPut,
    ).toHaveBeenCalledWith({
      networkId: networkId,
      gatewayId: 'testGatewayId0',
      resource: [
        {
          gateway_pool_id: 'pool4',
          mme_code: 4,
          mme_relative_capacity: 255,
        },
      ],
    });

    await wait();
    // $FlowFixMe
    lteNetworkIdGatewayPoolsGatewayPoolIdGet.mockResolvedValue({
      data: {
        config: {mme_group_id: 4},
        gateway_ids: ['testGatewayId0'],
        gateway_pool_id: 'pool4',
        gateway_pool_name: 'pool4',
      },
    });

    await wait();
    // check if only third tab (SecondaryEdit) is active
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('PrimaryEdit')).toBeNull();
    expect(queryByTestId('SecondaryEdit')).not.toBeNull();

    const mmeCodeSecondary = getByTestId('mmeCode').firstChild;
    const secondaryId = getByTestId('gwIdSecondary').firstChild;

    if (
      mmeCodeSecondary instanceof HTMLInputElement &&
      secondaryId instanceof HTMLElement
    ) {
      fireEvent.mouseDown(secondaryId);
      await wait();
      fireEvent.click(getByText('testGatewayId1'));
      fireEvent.change(mmeCodeSecondary, {target: {value: '5'}});
    } else {
      throw 'invalid type';
    }
    fireEvent.click(getByText('Save'));
    await wait();

    expect(
      MagmaAPI.lteGateways.lteNetworkIdGatewaysGatewayIdCellularPoolingPut,
    ).toHaveBeenCalledWith({
      networkId: networkId,
      gatewayId: 'testGatewayId1',
      resource: [
        {
          gateway_pool_id: 'pool4',
          mme_code: 5,
          mme_relative_capacity: 1,
        },
      ],
    });
  });
});
