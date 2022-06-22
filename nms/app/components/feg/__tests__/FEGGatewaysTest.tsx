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

import FEGGateways from '../FEGGateways';
import MagmaAPI from '../../../../api/MagmaAPI';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {FEGGatewayContextProvider} from '../FEGContext';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {fireEvent, render, wait} from '@testing-library/react';
import {mockAPI} from '../../../util/TestUtils';
import type {FederationGateway} from '../../../../generated-ts';

jest.mock('axios');
jest.mock('../../../../app/hooks/useSnackbar');

const mockGw0: FederationGateway = {
  id: 'test_feg_gw0',
  name: 'test_gateway0',
  description: 'hello I am a federated gateway',
  tier: 'default',
  device: {
    key: {key: '', key_type: 'SOFTWARE_ECDSA_SHA256'},
    hardware_id: 'c9439d30-61ef-46c7-93f2-e01fc144144d',
  },
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 300,
    checkin_interval: 60,
    checkin_timeout: 100,
  },
  federation: {
    aaa_server: {},
    eap_aka: {
      plmn_ids: [],
    },
    gx: {},
    gy: {},
    health: {
      health_services: [],
    },
    hss: {},
    s6a: {},
    served_network_ids: [],
    swx: {
      hlr_plmn_ids: [],
      server: {
        protocol: 'tcp',
      },
      servers: [],
    },
    csfb: {},
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

const mockGw1: FederationGateway = {
  ...mockGw0,
  id: 'test_gw1',
  name: 'test_gateway1',
  device: {
    ...mockGw0.device,
    key: {key: '', key_type: 'SOFTWARE_ECDSA_SHA256'},
    hardware_id: 'd1039d30-61ef-46c7-93f2-e01fc144144d',
  },
};

const fegGateways = {
  [mockGw0.id]: mockGw0,
  [mockGw1.id]: mockGw1,
};

describe('<FEGGatewaysTest />', () => {
  beforeEach(() => {
    // gateway context gets list of federation gateways
    mockAPI(
      MagmaAPI.federationGateways,
      'fegNetworkIdGatewaysGet',
      fegGateways,
    );
  });

  const Wrapper = () => (
    <MemoryRouter initialEntries={['/nms/mynetwork/gateways']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <FEGGatewayContextProvider networkId="mynetwork" networkType="feg">
            <Routes>
              <Route
                path="/nms/:networkId/gateways/"
                element={<FEGGateways />}
              />
            </Routes>
          </FEGGatewayContextProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
  it('test gateway table rendered correctly', async () => {
    const {getAllByRole} = render(<Wrapper />);
    await wait();
    const rowItems = getAllByRole('row');
    // first row is the header
    expect(rowItems[0]).toHaveTextContent('Name');
    expect(rowItems[0]).toHaveTextContent('Hardware UUID');

    expect(rowItems[1]).toHaveTextContent('test_gateway0');
    expect(rowItems[1]).toHaveTextContent(
      'c9439d30-61ef-46c7-93f2-e01fc144144d',
    );

    expect(rowItems[2]).toHaveTextContent('test_gateway1');
    expect(rowItems[2]).toHaveTextContent(
      'd1039d30-61ef-46c7-93f2-e01fc144144d',
    );
  });
  it('test gateway delete is working', async () => {
    jest.spyOn(
      MagmaAPI.federationGateways,
      'fegNetworkIdGatewaysGatewayIdDelete',
    );

    const {getByTestId, getByText} = render(<Wrapper />);
    await wait();
    fireEvent.click(getByTestId(`delete ${mockGw0.id}`));
    await wait();
    fireEvent.click(getByText('Confirm'));
    await wait();
    // make sure gateway was deleted
    expect(
      MagmaAPI.federationGateways.fegNetworkIdGatewaysGatewayIdDelete,
    ).toHaveBeenCalledWith({
      networkId: 'mynetwork',
      gatewayId: mockGw0.id,
    });
  });
});
