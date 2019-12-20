/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import DevicesStatusTable from '../DevicesStatusTable';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {MemoryRouter, Route} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';

import 'jest-dom/extend-expect';
import MagmaAPIBindings from '@fbcnms/magma-api';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, render, wait} from '@testing-library/react';

import {RAW_DEVICES} from '../test/DevicesMock';

jest.mock('@fbcnms/magma-api');

const Wrapper = () => (
  <MemoryRouter initialEntries={['/nms/network1/agents/']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Route
            path="/nms/:networkId/agents"
            render={() => <DevicesStatusTable />}
          />
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<DevicesStatusTable />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getSymphonyByNetworkIdDevices.mockResolvedValueOnce(
      RAW_DEVICES,
    );
  });

  it('renders', async () => {
    const {getByText, getAllByText} = render(<Wrapper />);

    await wait();

    expect(
      MagmaAPIBindings.getSymphonyByNetworkIdDevices,
    ).toHaveBeenCalledTimes(1);

    // expected headers and titles and buttons
    expect(getByText('Devices')).toBeInTheDocument();
    expect(getByText('New Device')).toBeInTheDocument();
    expect(getByText('Name')).toBeInTheDocument();
    expect(getByText('State')).toBeInTheDocument();
    expect(getByText('Managing Agent')).toBeInTheDocument();
    expect(getByText('Actions')).toBeInTheDocument();

    // expected devices
    expect(getByText('ens_switch_1')).toBeInTheDocument();
    expect(getByText('localhost_snmpd')).toBeInTheDocument();
    expect(getByText('mikrotik')).toBeInTheDocument();
    expect(getByText('ping_fb_dns_from_lab')).toBeInTheDocument();
    expect(getByText('ping_fb_dns_ken_laptop')).toBeInTheDocument();
    expect(getByText('ping_google_ipv6')).toBeInTheDocument();
    expect(getByText('ping_google_ipv6_ken_laptop')).toBeInTheDocument();

    // a few devices are expected to be managed by this agent
    expect(getAllByText('fbbosfbcdockerengine')).toHaveLength(2);
  });
});
