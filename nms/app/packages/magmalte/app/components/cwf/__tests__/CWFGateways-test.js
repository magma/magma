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

import {CWFGateways} from '../CWFGateways';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {MemoryRouter, Route, Switch} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import type {cwf_gateway} from '@fbcnms/magma-api';
import type {cwf_ha_pair} from '@fbcnms/magma-api';

import 'jest-dom/extend-expect';
import MagmaAPIBindings from '@fbcnms/magma-api';
import axiosMock from 'axios';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, getByTitle, render, wait, queryByTitle}
  from '@testing-library/react';

const CWF_HA_GATEWAY_1: cwf_gateway = {
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
    tier: 'default',
  },
  id: 'mock_cwf01',
  name: 'mock_cwf',
  description: 'mock gateway',
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
    checkin_time: 1,
    meta: {
      gps_latitude: '0',
      gps_longitude: '0',
      gps_connected: '0',
      enodeb_connected: '0',
      mme_connected: '0',
    },
  },
};

let CWF_HA_GATEWAY_2 = JSON.parse(JSON.stringify(CWF_HA_GATEWAY_1));
CWF_HA_GATEWAY_2.id = "mock_cwf02";
CWF_HA_GATEWAY_2.name = "mock_cwf2";
CWF_HA_GATEWAY_2.device.hardware_id = 'bb35dd3f-efaa-435a-bcb6-8168d0caf333';

const CWF_HA_PAIR: cwf_ha_pair = {
  config: {
    transport_virtual_ip: '10.10.10.12'
  },
  gateway_id_1: 'mock_cwf01',
  gateway_id_2: 'mock_cwf02',
  ha_pair_id: 'pair1',
  state: {
    ha_pair_status: {
      active_gateway: 'mock_cwf01'
    },
    gateway1_health: {
      status: 'HEALTHY',
      description: 'OK'
    },
    gateway2_health: {
      status: 'UNHEALTHY',
      description: 'Service restart'
    },
  }
};

jest.mock('axios');
jest.mock('@fbcnms/magma-api');

const Wrapper = () => (
  <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Switch>
            <Route path="/nms/:networkId" component={CWFGateways} />
          </Switch>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<CWFGateways />', () => {
  beforeEach(() => {
    axiosMock.get.mockResolvedValueOnce({
      data: [CWF_HA_GATEWAY_1, CWF_HA_GATEWAY_2],
    });
    MagmaAPIBindings.getCwfByNetworkIdGateways.mockResolvedValue({
      mock_cwf01: CWF_HA_GATEWAY_1,
      mock_cwf02: CWF_HA_GATEWAY_2,
    });
    MagmaAPIBindings.getCwfByNetworkIdHaPairs.mockResolvedValue({
      pair1: CWF_HA_PAIR,
    });
    MagmaAPIBindings.getNetworksByNetworkIdTiers.mockResolvedValueOnce([
      'default',
    ]);
  });

  afterEach(() => {
    axiosMock.get.mockClear();
  });

  it('renders', async () => {
    const {getAllByRole} = render(<Wrapper />);

    await wait();

    expect(MagmaAPIBindings.getCwfByNetworkIdGateways).toHaveBeenCalledTimes(1);
    expect(MagmaAPIBindings.getCwfByNetworkIdHaPairs).toHaveBeenCalledTimes(1);

    const rowItems = getAllByRole('row');
    expect(rowItems).toHaveLength(3);
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('Hardware UUID / GRE Key');

    expect(rowItems[1]).toHaveTextContent('mock_cwf');
    expect(rowItems[1]).toHaveTextContent(
        'a935dd3f-efaa-435a-bcb6-8168d0caf333'
    );
    const gatewayStatusCell = getByTitle(
        rowItems[1],
        'Last refreshed 12/31/1969, 7:00:00 PM'
    );
    expect(getByTitle(gatewayStatusCell, "Primary CWAG"))
        .toBeInTheDocument();

    expect(rowItems[2]).toHaveTextContent('mock_cwf2');
    expect(rowItems[2]).toHaveTextContent(
        'bb35dd3f-efaa-435a-bcb6-8168d0caf333',
    );
     const gatewayStatusCell2 = getByTitle(
        rowItems[2],
        'Last refreshed 12/31/1969, 7:00:00 PM'
    );
    expect(queryByTitle(gatewayStatusCell2, "Primary CWAG"))
        .not.toBeInTheDocument();
  });
});
