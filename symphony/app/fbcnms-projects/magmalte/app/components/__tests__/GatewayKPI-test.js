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
import GatewayKPIs from '../GatewayKPIs';
import MagmaAPIBindings from '@fbcnms/magma-api';
import React from 'react';
import axiosMock from 'axios';
import {MemoryRouter, Route} from 'react-router-dom';
import {cleanup, render, wait} from '@testing-library/react';
import type {lte_gateway} from '@fbcnms/magma-api';

afterEach(cleanup);

const mockSt: lte_gateway = {
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

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

describe('<GatewaysKPIs />', () => {
  beforeEach(() => {
    const mockUpSt = Object.assign({}, mockSt);
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
      test1: mockSt,
      test2: mockSt,
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
