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
import AddEditGatewayPoolButton from '../GatewayPoolEdit';
import GatewayContext from '../../../context/GatewayContext';
import GatewayPools from '../EquipmentGatewayPools';
import MagmaAPI from '../../../api/MagmaAPI';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {GatewayPoolsContextProvider} from '../../../context/GatewayPoolsContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {StyledEngineProvider, ThemeProvider} from '@mui/material/styles';
import {fireEvent, waitFor} from '@testing-library/react';
import {mockAPI, mockAPIOnce} from '../../../util/TestUtils';
import {render} from '../../../util/TestingLibrary';
import {useEnqueueSnackbar} from '../../../hooks/useSnackbar';
import type {LteGateway} from '../../../../generated';

jest.mock('axios');
jest.mock('../../../hooks/useSnackbar');

const lteNetworkIdGatewayPoolsGetMock = {
  pool1: {
    config: {mme_group_id: 1},
    gateway_ids: ['gw1', 'gw2', 'gw5'],
    gateway_pool_id: 'pool1',
    gateway_pool_name: 'pool_1',
  },
  pool2: {
    config: {mme_group_id: 2},
    gateway_ids: ['gw3', 'gw4'],
    gateway_pool_id: 'pool2',
    gateway_pool_name: 'pool_2',
  },
  pool3: {
    config: {mme_group_id: 3},
    gateway_ids: [],
    gateway_pool_id: 'pool3',
    gateway_pool_name: 'pool_3',
  },
};

const lteNetworkIdGatewaysGatewayIdCellularPoolingGetMock = [
  // gw1
  {
    gateway_pool_id: 'pool1',
    mme_code: 1,
    mme_relative_capacity: 255,
  },
  // gw2
  {
    gateway_pool_id: 'pool1',
    mme_code: 2,
    mme_relative_capacity: 255,
  },
  // gw5
  {
    gateway_pool_id: 'pool1',
    mme_code: 3,
    mme_relative_capacity: 1,
  },
  // gw3
  {
    gateway_pool_id: 'pool2',
    mme_code: 1,
    mme_relative_capacity: 255,
  },
  // gw4
  {
    gateway_pool_id: 'pool2',
    mme_code: 2,
    mme_relative_capacity: 1,
  },
];

const networkId = 'test';
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

const mockGw1 = Object.assign({}, mockGw0);
mockGw1.id = 'testGatewayId1';
mockGw1.name = 'testGateway1';

const lteGateways = {
  testGatewayId0: mockGw0,
  testGatewayId1: mockGw1,
};
describe('<GatewayPools />', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.lteNetworks,
      'lteNetworkIdGatewayPoolsGet',
      lteNetworkIdGatewayPoolsGetMock,
    );
    for (const element of lteNetworkIdGatewaysGatewayIdCellularPoolingGetMock) {
      mockAPIOnce(
        MagmaAPI.lteGateways,
        'lteNetworkIdGatewaysGatewayIdCellularPoolingGet',
        [element],
      );
    }
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/test/pools']} initialIndex={0}>
      <StyledEngineProvider injectFirst>
        <ThemeProvider theme={defaultTheme}>
          <GatewayPoolsContextProvider networkId={networkId}>
            <Routes>
              <Route path="/nms/:networkId/pools/" element={<GatewayPools />} />
            </Routes>
          </GatewayPoolsContextProvider>
        </ThemeProvider>
      </StyledEngineProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    jest
      .spyOn(
        MagmaAPI.lteNetworks,
        'lteNetworkIdGatewayPoolsGatewayPoolIdDelete',
      )
      .mockImplementation();

    const {findAllByRole, findByText, getByText, openActionsTableMenu} = render(
      <Wrapper />,
    );

    const rowItems = await findAllByRole('row');

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
    await openActionsTableMenu(2);
    fireEvent.click(getByText('Remove'));
    expect(
      await findByText('Are you sure you want to delete pool3?'),
    ).toBeInTheDocument();
    fireEvent.click(getByText('Confirm'));

    await waitFor(() =>
      expect(
        MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdDelete,
      ).toHaveBeenCalledWith({
        networkId: 'test',
        gatewayPoolId: 'pool3',
      }),
    );
  });

  it('verify gateway pool edit', async () => {
    mockAPI(MagmaAPI.lteNetworks, 'lteNetworkIdGatewayPoolsGatewayPoolIdPut');

    const {
      queryByTestId,
      getByTestId,
      getByText,
      openActionsTableMenu,
    } = render(<Wrapper />);

    await openActionsTableMenu(2);
    fireEvent.click(getByText('Edit'));

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('PrimaryEdit')).toBeNull();
    expect(queryByTestId('SecondaryEdit')).toBeNull();

    const name = getByTestId('name').firstChild as HTMLInputElement;
    const poolId = getByTestId('poolId').firstChild as HTMLInputElement;
    const mmeGroupId = getByTestId('mmeGroupId').firstChild as HTMLInputElement;

    expect(poolId).toBeDisabled();

    fireEvent.change(name, {target: {value: 'foo'}});
    fireEvent.change(mmeGroupId, {target: {value: '4'}});

    fireEvent.click(getByText('Save And Continue'));

    await waitFor(() => {
      expect(
        MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsGatewayPoolIdPut,
      ).toHaveBeenCalledWith({
        networkId,
        gatewayPoolId: 'pool3',
        hAGatewayPool: {
          config: {mme_group_id: 4},
          gateway_pool_id: 'pool3',
          gateway_pool_name: 'foo',
        },
      });
    });
  });
});

describe('<AddEditGatewayPoolButton />', () => {
  let lteNetworkIdGatewayPoolsGatewayPoolIdGet: jest.SpyInstance;

  beforeEach(() => {
    (useEnqueueSnackbar as jest.Mock).mockReturnValue(jest.fn());
    mockAPI(MagmaAPI.lteNetworks, 'lteNetworkIdGatewayPoolsPost');

    mockAPI(
      MagmaAPI.lteGateways,
      'lteNetworkIdGatewaysGatewayIdCellularPoolingPut',
    );

    lteNetworkIdGatewayPoolsGatewayPoolIdGet = mockAPI(
      MagmaAPI.lteNetworks,
      'lteNetworkIdGatewayPoolsGatewayPoolIdGet',
      {
        config: {mme_group_id: 4},
        gateway_ids: [],
        gateway_pool_id: 'pool4',
        gateway_pool_name: 'pool4',
      },
    );
  });
  const AddWrapper = () => {
    return (
      <MemoryRouter initialEntries={['/nms/test/pools']} initialIndex={0}>
        <StyledEngineProvider injectFirst>
          <ThemeProvider theme={defaultTheme}>
            <ThemeProvider theme={defaultTheme}>
              <GatewayContext.Provider
                value={{
                  state: lteGateways,
                  setState: async () => {},
                  updateGateway: async () => {},
                  refetch: () => {},
                }}>
                <GatewayPoolsContextProvider networkId={networkId}>
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
                </GatewayPoolsContextProvider>
              </GatewayContext.Provider>
            </ThemeProvider>
          </ThemeProvider>
        </StyledEngineProvider>
      </MemoryRouter>
    );
  };

  it('verify gateway pool add', async () => {
    const networkId = 'test';
    const {findByText, queryByTestId, getByTestId, getByText} = render(
      <AddWrapper />,
    );
    expect(queryByTestId('gatewayPoolEditDialog')).toBeNull();
    fireEvent.click(await findByText('Add Gateway Pool'));

    // check if only first tab (config) is active
    expect(queryByTestId('configEdit')).not.toBeNull();
    expect(queryByTestId('PrimaryEdit')).toBeNull();
    expect(queryByTestId('SecondaryEdit')).toBeNull();

    const name = getByTestId('name').firstChild as HTMLInputElement;
    const poolId = getByTestId('poolId').firstChild as HTMLInputElement;
    const mmeGroupId = getByTestId('mmeGroupId').firstChild as HTMLInputElement;

    fireEvent.change(name, {target: {value: 'pool_4'}});
    fireEvent.change(poolId, {target: {value: 'pool4'}});
    fireEvent.change(mmeGroupId, {target: {value: '4'}});

    fireEvent.click(getByText('Save And Continue'));

    const newGatewayPool = {
      config: {mme_group_id: 4},
      gateway_pool_id: 'pool4',
      gateway_pool_name: 'pool_4',
    };

    await waitFor(() =>
      expect(
        MagmaAPI.lteNetworks.lteNetworkIdGatewayPoolsPost,
      ).toHaveBeenCalledWith({
        networkId,
        hAGatewayPool: newGatewayPool,
      }),
    );

    await waitFor(() =>
      lteNetworkIdGatewayPoolsGatewayPoolIdGet.mockResolvedValue({
        data: {
          config: {mme_group_id: 4},
          gateway_ids: [],
          gateway_pool_id: 'pool4',
          gateway_pool_name: 'pool4',
        },
      }),
    );

    // check if only second tab (PrimaryEdit) is active
    expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('PrimaryEdit')).not.toBeNull();
    expect(queryByTestId('SecondaryEdit')).toBeNull();

    const mmeCode = getByTestId('mmeCode').firstChild as HTMLInputElement;
    const PrimaryId = getByTestId('gwIdPrimary').firstChild as HTMLElement;

    fireEvent.mouseDown(PrimaryId);
    fireEvent.click(await findByText('testGatewayId0'));
    fireEvent.change(mmeCode, {target: {value: '4'}});

    fireEvent.click(getByText('Save And Continue'));
    await waitFor(() =>
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
      }),
    );
    lteNetworkIdGatewayPoolsGatewayPoolIdGet.mockResolvedValue({
      data: {
        config: {mme_group_id: 4},
        gateway_ids: ['testGatewayId0'],
        gateway_pool_id: 'pool4',
        gateway_pool_name: 'pool4',
      },
    }),
      // check if only third tab (SecondaryEdit) is active
      expect(queryByTestId('configEdit')).toBeNull();
    expect(queryByTestId('PrimaryEdit')).toBeNull();
    expect(queryByTestId('SecondaryEdit')).not.toBeNull();

    const mmeCodeSecondary = getByTestId('mmeCode')
      .firstChild as HTMLInputElement;
    const secondaryId = getByTestId('gwIdSecondary').firstChild as HTMLElement;

    fireEvent.mouseDown(secondaryId);
    fireEvent.click(await findByText('testGatewayId1'));
    fireEvent.change(mmeCodeSecondary, {target: {value: '5'}});

    fireEvent.click(getByText('Save'));

    await waitFor(() =>
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
      }),
    );
  });
});
