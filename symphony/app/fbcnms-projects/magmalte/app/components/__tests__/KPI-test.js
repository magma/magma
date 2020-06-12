/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import 'jest-dom/extend-expect';
import EnodebKPIs from '../EnodebKPIs';
import GatewayKPIs from '../GatewayKPIs';
import MagmaAPIBindings from '@fbcnms/magma-api';
import React from 'react';
import axiosMock from 'axios';
import {MemoryRouter, Route} from 'react-router-dom';
import {cleanup, render, wait} from '@testing-library/react';
import type {enodeb, enodeb_state, lte_gateway} from '@fbcnms/magma-api';

afterEach(cleanup);

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

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

describe('<GatewaysKPIs />', () => {
  beforeEach(() => {
    const mockUpSt = Object.assign({}, mockGwSt);
    mockUpSt['status'] = {
      checkin_time: Date.now(),
      meta: {
        gps_latitude: '0',
        gps_longitude: '0',
        gps_connected: '0',
        enodeb_connected: '0',
        mme_connected: '0',
      },
    };
    MagmaAPIBindings.getLteByNetworkIdGateways.mockResolvedValue({
      test1: mockGwSt,
      test2: mockGwSt,
      test3: mockUpSt,
    });
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
      <Route path="/nms/:networkId" component={GatewayKPIs} />
    </MemoryRouter>
  );
  it('renders', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();

    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(getByTestId('Connected')).toHaveTextContent('1');
    expect(getByTestId('Disconnected')).toHaveTextContent('2');
  });
});

describe('<EnodebKPIs />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getLteByNetworkIdEnodebs.mockResolvedValue(mockEnbAll);
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState.mockResolvedValue(
      mockEnbSt,
    );
    const mockEnbNotTxSt = Object.assign({}, mockEnbSt);
    mockEnbNotTxSt.rf_tx_on = false;
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState.mockReturnValueOnce(
      mockEnbNotTxSt,
    );
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
      <Route path="/nms/:networkId" component={EnodebKPIs} />
    </MemoryRouter>
  );
  it('renders', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();

    expect(MagmaAPIBindings.getLteByNetworkIdEnodebs).toHaveBeenCalledTimes(1);
    expect(
      MagmaAPIBindings.getLteByNetworkIdEnodebsByEnodebSerialState,
    ).toHaveBeenCalledTimes(3);
    expect(getByTestId('Total')).toHaveTextContent('3');
    expect(getByTestId('Transmitting')).toHaveTextContent('2');
  });
});
