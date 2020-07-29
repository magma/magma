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

import Gateways from '../Gateways';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {MemoryRouter, Route, Switch} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import type {lte_gateway} from '@fbcnms/magma-api';

import 'jest-dom/extend-expect';
import MagmaAPIBindings from '@fbcnms/magma-api';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, fireEvent, render, wait} from '@testing-library/react';

const OFFLINE_GATEWAY: lte_gateway = {
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
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
    tier: 'tier2',
  },
  id: 'murt_usa',
  name: 'murt_test',
  description: 'hello I am a gateway',
  tier: 'default',
  device: {
    hardware_id: 'a935dd3f-efaa-435a-bcb6-8168d0caf333',
    key: {
      key:
        'MHYwEAYHKoZIzj0CAQYFK4EEACIDYgAESpHVXt266GW0WTKL0CvIvEFpECQL0rkHgs5Bc0efoSde01wuphb8tK1zL9t8rsVFlv2tyUHXeoJt7/AaEonGYOuEkbHocRy9LBAVue2sOFWrIhJvqieujrd15dLH1zBm',
      key_type: 'SOFTWARE_ECDSA_SHA256',
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

const Wrapper = () => (
  <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Switch>
            <Route path="/nms/:networkId" component={Gateways} />
          </Switch>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<Gateways />', () => {
  beforeEach(() => {
    axiosMock.get.mockResolvedValueOnce({
      data: [OFFLINE_GATEWAY],
    });
    MagmaAPIBindings.getLteByNetworkIdGateways.mockResolvedValue({
      murt_usa: OFFLINE_GATEWAY,
    });
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockResolvedValueOnce([
      'default',
    ]);
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  it('renders', async () => {
    const {getByText} = render(<Wrapper />);

    await wait();

    expect(MagmaAPIBindings.getLteByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(getByText('Configure Gateways')).toBeInTheDocument();
    expect(getByText('Hardware UUID')).toBeInTheDocument();
    expect(getByText('murt_test')).toBeInTheDocument();
  });

  it('shows dialog when Add Gateway is clicked', async () => {
    const {getByText} = render(<Wrapper />);

    await wait();

    fireEvent.click(getByText('Add Gateway'));

    await wait();

    expect(getByText('Gateway Name')).toBeInTheDocument();
    expect(getByText('Gateway ID')).toBeInTheDocument();
    expect(getByText('Challenge Key')).toBeInTheDocument();
    expect(getByText('Upgrade Tier')).toBeInTheDocument();
  });

  it('shows prompt when delete is clicked', async () => {
    MagmaAPIBindings.deleteLteByNetworkIdGatewaysByGatewayId.mockResolvedValueOnce(
      {},
    );

    const {getByText, getByTestId} = render(<Wrapper />);
    await wait();

    fireEvent.click(getByTestId('delete-gateway-icon'));

    expect(
      getByText('Are you sure you want to delete murt_test?'),
    ).toBeInTheDocument();

    // Confirm deletion
    fireEvent.click(getByText('Confirm'));
    await wait();
    expect(
      MagmaAPIBindings.deleteLteByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledTimes(1);

    axiosMock.delete.mockClear();
  });

  it('doesnt delete when cancel is clicked', async () => {
    axiosMock.delete.mockResolvedValueOnce({
      data: {success: true},
    });

    const {getByText, getByTestId} = render(<Wrapper />);
    await wait();

    fireEvent.click(getByTestId('delete-gateway-icon'));

    // Cancel deletion
    fireEvent.click(getByText('Cancel'));
    await wait();
    expect(axiosMock.delete).toHaveBeenCalledTimes(0);

    axiosMock.delete.mockClear();
  });

  it('shows dialog when edit is clicked', async () => {
    const {getByText, getByTestId} = render(<Wrapper />);

    await wait();

    fireEvent.click(getByTestId('edit-gateway-icon'));

    expect(getByText('Reboot Gateway')).toBeInTheDocument();
  });
});
