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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import defaultTheme from '../../theme/default';

import MagmaAPI from '../../../api/MagmaAPI';
import {GenericCommandResponse} from '../../../generated-ts';
import {MemoryRouter, Route, Routes} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {TroubleshootingControl} from '../GatewayCommandFields';
import {mockAPI, mockAPIError} from '../../util/TestUtils';
import {render, wait} from '@testing-library/react';

jest.mock('../../../generated/MagmaAPIBindings');
jest.mock('../../../app/hooks/useSnackbar');

const Wrapper = () => (
  <MemoryRouter initialEntries={['/nms/mynetwork']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <Routes>
          <Route
            path="/nms/:networkId"
            element={<TroubleshootingControl gatewayID={'test_gateway'} />}
          />
        </Routes>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

const validControlProxyContent = ({
  response: {
    returncode: 0,
    stderr: '',
    stdout: 'fluentd_address: fluentd.magma.io\nfluentd_port: 24224',
  },
} as unknown) as GenericCommandResponse;

describe('<verify successful aggregation validation/>', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.commands,
      'networksNetworkIdGatewaysGatewayIdCommandGenericPost',
      validControlProxyContent,
    );

    mockAPI(MagmaAPI.events, 'eventsNetworkIdAboutCountGet', 0);
  });

  it('', async () => {
    const {getByTestId, getAllByTestId} = render(<Wrapper />);
    await wait();
    const controProxyValidationContent = getByTestId(
      'Control Proxy Config Validation',
    );
    const apiValidationContent = getByTestId('API validation');

    expect(controProxyValidationContent).toHaveTextContent('Good');
    expect(getAllByTestId('fileContent')[0]).toHaveTextContent(
      'fluentd_address: fluentd.magma.io fluentd_port: 24224',
    );
    expect(apiValidationContent).toHaveTextContent('Good');
  });
});

describe('<verify control proxy validation failure/>', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.commands,
      'networksNetworkIdGatewaysGatewayIdCommandGenericPost',
      ({
        response: {
          returncode: 0,
          stderr: 'Error, file not found',
          stdout: '',
        },
      } as unknown) as GenericCommandResponse,
    );

    mockAPI(MagmaAPI.events, 'eventsNetworkIdAboutCountGet', 0);
  });

  it('', async () => {
    const {getByTestId, getAllByTestId} = render(<Wrapper />);
    await wait();
    const controProxyValidationContent = getByTestId(
      'Control Proxy Config Validation',
    );
    const apiValidationContent = getByTestId('API validation');

    expect(controProxyValidationContent).toHaveTextContent('Bad');
    expect(getAllByTestId('fileError')[0]).toHaveTextContent('file not found');
    expect(apiValidationContent).toHaveTextContent('Good');
  });
});

describe('<verify api validation failure/>', () => {
  beforeEach(() => {
    mockAPI(
      MagmaAPI.commands,
      'networksNetworkIdGatewaysGatewayIdCommandGenericPost',
      validControlProxyContent,
    );

    mockAPIError(MagmaAPI.events, 'eventsNetworkIdAboutCountGet');
  });

  it('', async () => {
    const {getByTestId} = render(<Wrapper />);
    await wait();
    const controProxyValidationContent = getByTestId(
      'Control Proxy Config Validation',
    );
    const apiValidationContent = getByTestId('API validation');

    expect(controProxyValidationContent).toHaveTextContent('Good');
    expect(apiValidationContent).toHaveTextContent('Bad');
  });
});
