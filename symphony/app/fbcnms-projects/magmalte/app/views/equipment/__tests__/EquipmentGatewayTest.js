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
import Gateway, {DATE_TO_STRING_PARAMS} from '../EquipmentGateway';
import MagmaAPIBindings from '@fbcnms/magma-api';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {cleanup, render, wait} from '@testing-library/react';
import type {lte_gateway} from '@fbcnms/magma-api';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

afterEach(cleanup);

const mockGw0: lte_gateway = {
  id: 'test_gw0',
  name: 'test_gateway0',
  description: 'test_gateway0',
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
const currTime = Date.now();

describe('<Gateway />', () => {
  beforeEach(() => {
    // eslint-disable-next-line max-len
    const mockGw1 = Object.assign({}, mockGw0);
    const mockGw2 = Object.assign({}, mockGw0);
    mockGw1.id = 'test_gw1';
    mockGw1.name = 'test_gateway1';
    mockGw1.connected_enodeb_serials = ['xxx', 'yyy'];

    mockGw2.id = 'test_gw2';
    mockGw2.name = 'test_gateway2';
    mockGw2.connected_enodeb_serials = ['xxx'];
    mockGw2.status = {
      checkin_time: currTime,
    };
    MagmaAPIBindings.getLteByNetworkIdGateways.mockResolvedValue({
      test1: mockGw0,
      test2: mockGw1,
      test3: mockGw2,
    });
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateway']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <Route path="/nms/:networkId/gateway/" component={Gateway} />
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );

  it('renders', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    let testElem = getByTestId('gatewayInfo-0');
    expect(testElem).toHaveTextContent('test_gw0');
    expect(testElem).toHaveTextContent('test_gateway0');
    expect(testElem).toHaveTextContent('0');
    expect(testElem).toHaveTextContent('Bad');
    expect(testElem).toHaveTextContent(
      new Date(0).toLocaleDateString(...DATE_TO_STRING_PARAMS),
    );

    testElem = getByTestId('gatewayInfo-1');
    expect(testElem).toHaveTextContent('test_gw1');
    expect(testElem).toHaveTextContent('test_gateway1');
    expect(testElem).toHaveTextContent('2');
    expect(testElem).toHaveTextContent('Bad');
    expect(testElem).toHaveTextContent(
      new Date(0).toLocaleDateString(...DATE_TO_STRING_PARAMS),
    );

    testElem = getByTestId('gatewayInfo-2');
    expect(testElem).toHaveTextContent('test_gw2');
    expect(testElem).toHaveTextContent('test_gateway2');
    expect(testElem).toHaveTextContent('1');
    expect(testElem).toHaveTextContent('Good');
    expect(testElem).toHaveTextContent(
      new Date(currTime).toLocaleDateString(...DATE_TO_STRING_PARAMS),
    );
  });
});
