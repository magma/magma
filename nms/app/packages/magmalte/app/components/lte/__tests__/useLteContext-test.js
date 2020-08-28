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
 * @flow
 * @format
 */
import type {enodeb, enodeb_state, lte_gateway} from '@fbcnms/magma-api';

import MagmaAPIBindings from '@fbcnms/magma-api';

import {CWF, LTE} from '@fbcnms/types/network';
import {act, renderHook} from '@testing-library/react-hooks';
import {useLteContext} from '../LteSections';

const enqueueSnackbarMock = jest.fn();
jest.mock('@fbcnms/magma-api');
jest.mock('mapbox-gl', () => {});
jest.mock('@fbcnms/ui/insights/map/MapView', () => {});
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

const mockGwSt: lte_gateway = {
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
};

const mockEnbAll: {[string]: enodeb} = {
  test1: {
    name: 'test1',
    serial: 'test1',
    config: {
      cell_id: 0,
      device_class: 'Baicells Nova-233 G2 OD FDD',
      transmit_enabled: true,
    },
  },
  test2: {
    name: 'test2',
    serial: 'test2',
    config: {
      cell_id: 0,
      device_class: 'Baicells Nova-233 G2 OD FDD',
      transmit_enabled: true,
    },
  },
  test3: {
    name: 'test3',
    serial: 'test3',
    config: {
      cell_id: 0,
      device_class: 'Baicells Nova-233 G2 OD FDD',
      transmit_enabled: true,
    },
  },
};

const mockEnbSt: enodeb_state = {
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
};

const mockSubscribers = {
  IMSI00000000001002: {
    active_apns: ['oai.ipv4'],
    id: 'IMSI722070171001002',
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
};
const mockRan = {
  bandwidth_mhz: 20,
  tdd_config: {
    earfcndl: 44390,
    special_subframe_pattern: 7,
    subframe_assignment: 2,
  },
};
const mockTiers = ['default'];
const mockDefaultTier = {
  gateways: [],
  id: 'default',
  images: null,
  version: '',
};
const mockStableChannel = {
  id: 'stable',
  name: 'stable',
  supported_versions: ['0.3.44-1510352717-b7151784'],
};

const mockLteGateways = {
  test1: mockGwSt,
  test2: mockGwSt,
  test3: mockGwSt,
};

describe('use Lte Context testing', () => {
  afterEach(() => {
    MagmaAPIBindings.getLteByNetworkIdGateways.mockClear();
    MagmaAPIBindings.getLteByNetworkIdEnodebs.mockClear();

    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState.mockClear();
    MagmaAPIBindings.getLteByNetworkIdSubscribers.mockClear();
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockClear();
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockClear();
    MagmaAPIBindings.getNetworksByNetworkIdTiersByTierId.mockClear();
    MagmaAPIBindings.getChannelsByChannelId.mockClear();
  });

  test('verify lte context creation with context disabled', async () => {
    const {result} = renderHook(() => useLteContext('network1', LTE, true), {});
    expect(result.current).toBe(null);
  });

  test('verify lte context creation with invalid networkID', async () => {
    const {result} = renderHook(() => useLteContext('', LTE, true), {});
    expect(result.current).toBe(null);
  });

  test('verify lte context creation with non LTE network', async () => {
    const {result} = renderHook(() => useLteContext('network1', CWF, true), {});
    expect(result.current).toBe(null);
  });

  test('verify lte context creation with LTE network', async () => {
    MagmaAPIBindings.getLteByNetworkIdGateways.mockResolvedValue(
      mockLteGateways,
    );
    MagmaAPIBindings.getLteByNetworkIdEnodebs.mockResolvedValue(mockEnbAll);

    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState.mockResolvedValue(
      mockEnbSt,
    );
    MagmaAPIBindings.getLteByNetworkIdSubscribers.mockResolvedValue(
      mockSubscribers,
    );
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockResolvedValue(mockRan);
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockResolvedValue(mockTiers);
    MagmaAPIBindings.getNetworksByNetworkIdTiersByTierId.mockResolvedValue(
      mockDefaultTier,
    );
    MagmaAPIBindings.getChannelsByChannelId.mockResolvedValue(
      mockStableChannel,
    );

    const {result, waitForNextUpdate} = renderHook(
      () => useLteContext('network1', LTE, false),
      {},
    );
    await act(async () => {
      // State is updated when we wait for the update, so we need this wrapped
      // in act
      await waitForNextUpdate();
    });
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getLteByNetworkIdEnodebs).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState,
    ).toHaveBeenCalledTimes(3);
    expect(MagmaAPIBindings.getLteByNetworkIdSubscribers).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getNetworksByNetworkIdTiers).toHaveBeenCalledTimes(
      1,
    );
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdTiersByTierId,
    ).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getChannelsByChannelId).toHaveBeenCalledTimes(1);

    const enbInfo = {
      test1: {
        enb: mockEnbAll['test1'],
        enb_state: mockEnbSt,
      },
      test2: {
        enb: mockEnbAll['test2'],
        enb_state: mockEnbSt,
      },
      test3: {
        enb: mockEnbAll['test3'],
        enb_state: mockEnbSt,
      },
    };
    expect(result.current).toBeValid;
    if (result.current) {
      const {
        enodebCtx,
        gatewayCtx,
        gatewayTierCtx,
        subscriberCtx,
      } = result.current;
      expect(enodebCtx.state.enbInfo).toStrictEqual(enbInfo);
      expect(enodebCtx.lteRanConfigs).toStrictEqual(mockRan);
      expect(gatewayCtx.state).toStrictEqual(mockLteGateways);
      expect(subscriberCtx.state).toStrictEqual(mockSubscribers);
      expect(gatewayTierCtx.state.tiers).toStrictEqual({
        default: mockDefaultTier,
      });
      expect(gatewayTierCtx.state.supportedVersions).toStrictEqual(
        mockStableChannel.supported_versions,
      );
    }
  });

  test('verify lte context creation with LTE network with API errors', async () => {
    MagmaAPIBindings.getLteByNetworkIdGateways.mockRejectedValue(
      new Error('error'),
    );
    MagmaAPIBindings.getLteByNetworkIdEnodebs.mockResolvedValue(mockEnbAll);

    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState.mockRejectedValue(
      new Error('error'),
    );
    MagmaAPIBindings.getLteByNetworkIdSubscribers.mockResolvedValue(
      mockSubscribers,
    );
    MagmaAPIBindings.getLteByNetworkIdCellularRan.mockResolvedValue(mockRan);
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockResolvedValue(mockTiers);
    MagmaAPIBindings.getNetworksByNetworkIdTiersByTierId.mockResolvedValue(
      mockDefaultTier,
    );
    MagmaAPIBindings.getChannelsByChannelId.mockResolvedValue(
      mockStableChannel,
    );

    const {result, waitForNextUpdate} = renderHook(
      () => useLteContext('network1', LTE, false),
      {},
    );
    await act(async () => {
      // State is updated when we wait for the update, so we need this wrapped
      // in act
      await waitForNextUpdate();
    });
    // expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getLteByNetworkIdEnodebs).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState,
    ).toHaveBeenCalledTimes(3);
    expect(MagmaAPIBindings.getLteByNetworkIdSubscribers).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getLteByNetworkIdCellularRan).toHaveBeenCalledTimes(
      1,
    );
    expect(MagmaAPIBindings.getNetworksByNetworkIdTiers).toHaveBeenCalledTimes(
      1,
    );
    expect(
      MagmaAPIBindings.getNetworksByNetworkIdTiersByTierId,
    ).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getChannelsByChannelId).toHaveBeenCalledTimes(1);

    const enbInfo = {
      test1: {
        enb: mockEnbAll['test1'],
        enb_state: {},
      },
      test2: {
        enb: mockEnbAll['test2'],
        enb_state: {},
      },
      test3: {
        enb: mockEnbAll['test3'],
        enb_state: {},
      },
    };
    expect(result.current).toBeValid;
    if (result.current) {
      const {
        enodebCtx,
        gatewayCtx,
        gatewayTierCtx,
        subscriberCtx,
      } = result.current;
      expect(enodebCtx.state.enbInfo).toStrictEqual(enbInfo);
      expect(enodebCtx.lteRanConfigs).toStrictEqual(mockRan);
      expect(gatewayCtx.state).toStrictEqual({});
      expect(subscriberCtx.state).toStrictEqual(mockSubscribers);
      expect(gatewayTierCtx.state.tiers).toStrictEqual({
        default: mockDefaultTier,
      });
      expect(gatewayTierCtx.state.supportedVersions).toStrictEqual(
        mockStableChannel.supported_versions,
      );
    }
  });

  test('verify subscriber context', async () => {
    MagmaAPIBindings.getLteByNetworkIdSubscribers.mockResolvedValue(
      mockSubscribers,
    );

    const {result, waitForNextUpdate} = renderHook(
      () => useLteContext('network1', LTE, false),
      {},
    );
    await act(async () => {
      // State is updated when we wait for the update, so we need this wrapped
      // in act
      await waitForNextUpdate();
    });
    expect(MagmaAPIBindings.getLteByNetworkIdSubscribers).toHaveBeenCalledTimes(
      1,
    );

    if (!result.current) {
      throw 'result invalid';
    }
    let {subscriberCtx} = result.current;
    expect(subscriberCtx.state).toStrictEqual(mockSubscribers);

    // verify subscriber add
    const newSubscriber = {
      active_apns: ['oai.ipv4'],
      id: 'IMSI722070171001003',
      lte: {
        auth_algo: 'MILENAGE',
        auth_key: 'i69HPy+P0JSHzMvXCXxoYg==',
        auth_opc: 'jie2rw5pLnUPMmZ6OxRgXQ==',
        state: 'ACTIVE',
        sub_profile: 'default',
      },
    };
    MagmaAPIBindings.getLteByNetworkIdSubscribersBySubscriberId.mockResolvedValue(
      newSubscriber,
    );
    MagmaAPIBindings.postLteByNetworkIdSubscribers.mockResolvedValue({
      success: true,
    });

    // create subscriber
    subscriberCtx.setState?.('IMSI722070171001003', newSubscriber);
    await waitForNextUpdate();
    expect(
      MagmaAPIBindings.postLteByNetworkIdSubscribers,
    ).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPIBindings.getLteByNetworkIdSubscribersBySubscriberId,
    ).toHaveBeenCalledTimes(1);

    if (!result.current) {
      throw 'result invalid';
    }
    subscriberCtx = result.current.subscriberCtx;
    expect(subscriberCtx?.state['IMSI722070171001003']).toBe(newSubscriber);
    MagmaAPIBindings.getLteByNetworkIdSubscribersBySubscriberId.mockClear();

    // update subscriber
    newSubscriber.lte.state = 'INACTIVE';
    subscriberCtx.setState?.('IMSI722070171001003', newSubscriber);
    MagmaAPIBindings.putLteByNetworkIdSubscribersBySubscriberId.mockResolvedValue(
      {success: true},
    );
    MagmaAPIBindings.getLteByNetworkIdSubscribersBySubscriberId.mockResolvedValue(
      newSubscriber,
    );
    await waitForNextUpdate();
    expect(
      MagmaAPIBindings.putLteByNetworkIdSubscribersBySubscriberId,
    ).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPIBindings.getLteByNetworkIdSubscribersBySubscriberId,
    ).toHaveBeenCalledTimes(1);
    if (!result.current) {
      throw 'result invalid';
    }
    subscriberCtx = result.current.subscriberCtx;
    expect(subscriberCtx?.state['IMSI722070171001003']).toBe(newSubscriber);

    // delete subscriber
    subscriberCtx.setState?.('IMSI722070171001003');
    MagmaAPIBindings.deleteLteByNetworkIdSubscribersBySubscriberId.mockResolvedValue(
      {success: true},
    );
    await waitForNextUpdate();
    expect(
      MagmaAPIBindings.deleteLteByNetworkIdSubscribersBySubscriberId,
    ).toHaveBeenCalledTimes(1);
    if (!result.current) {
      throw 'result invalid';
    }
    subscriberCtx = result.current.subscriberCtx;
    expect('IMSI722070171001003' in subscriberCtx?.state).toBeFalse;
  });
});
