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
import GatewaySummary from '../GatewaySummary';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../../theme/default';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {render} from '@testing-library/react';
import type {LteGateway} from '../../../../generated-ts';

const mockGatewaySt: LteGateway = {
  cellular: {
    epc: {
      ip_block: '',
      nat_enabled: true,
    },
    ran: {
      pci: 260,
      transmit_enabled: true,
    },
  },
  connected_enodeb_serials: [],
  description: 'mpk_dogfooding',
  device: {
    hardware_id: 'e059637f-cd55-4109-816c-ce6ebc69020d',
    key: {
      key: '',
      key_type: 'SOFTWARE_ECDSA_SHA256',
    },
  },
  id: 'mpk_dogfooding_magma_1',
  magmad: {
    autoupgrade_enabled: true,
    autoupgrade_poll_interval: 301,
    checkin_interval: 60,
    checkin_timeout: 20,
  },
  name: 'team pod',
  status: {
    hardware_id: 'e059637f-cd55-4109-816c-ce6ebc69020d',
    platform_info: {
      packages: [
        {
          name: 'magma',
          version: '1.1.0-1590005479-e6e781a9',
        },
      ],
    },
  },
  tier: 'default',
  checked_in_recently: false,
};

describe('<GatewaySummary />', () => {
  it('renders', () => {
    const {container} = render(
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <GatewaySummary gwInfo={mockGatewaySt} />
        </MuiStylesThemeProvider>
      </MuiThemeProvider>,
    );
    expect(container).toHaveTextContent('mpk_dogfooding');
    expect(container).toHaveTextContent('1.1.0-1590005479-e6e781a9');
    expect(container).toHaveTextContent('e059637f-cd55-4109-816c-ce6ebc69020d');
    expect(container).toHaveTextContent('mpk_dogfooding_magma_1');
  });
});
