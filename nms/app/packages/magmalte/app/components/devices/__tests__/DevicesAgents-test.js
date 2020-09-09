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

import DevicesAgents from '../DevicesAgents';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import {MemoryRouter, Route, Switch} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';

import 'jest-dom/extend-expect';
import MagmaAPIBindings from '@fbcnms/magma-api';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, render, wait} from '@testing-library/react';

import {RAW_AGENT} from '../test/DevicesMock';

jest.mock('@fbcnms/magma-api');

const Wrapper = () => (
  <MemoryRouter initialEntries={['/nms/network1/agents/']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>
          <Switch>
            <Route
              path="/nms/:networkId/agents"
              render={() => <DevicesAgents />}
            />
          </Switch>
        </SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<DevicesAgents />', () => {
  beforeEach(() => {
    MagmaAPIBindings.getSymphonyByNetworkIdAgents.mockResolvedValueOnce({
      [RAW_AGENT.id]: RAW_AGENT,
    });
  });

  it('renders', async () => {
    const {getByText} = render(<Wrapper />);

    await wait();

    expect(MagmaAPIBindings.getSymphonyByNetworkIdAgents).toHaveBeenCalledTimes(
      1,
    );
    expect(getByText('Configure Agents')).toBeInTheDocument();
    expect(
      getByText('faceb00c-face-b00c-face-000c2940b2bf'),
    ).toBeInTheDocument();
    expect(getByText('fbbosfbcdockerengine')).toBeInTheDocument();
  });
});
