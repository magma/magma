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

import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import React from 'react';
import WifiDeviceDialog from '../WifiDeviceDialog';
import {MemoryRouter, Route, Switch} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';

import MagmaAPIBindings from '@fbcnms/magma-api';
import defaultTheme from '@fbcnms/ui/theme/default';

import {cleanup, fireEvent, render, wait} from '@testing-library/react';

import {RAW_GATEWAY} from '../test/GatewayMock';

jest.mock('axios');
jest.mock('@fbcnms/magma-api');
jest.mock('@fbcnms/ui/hooks/useSnackbar');

const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

const Wrapper = props => (
  <MemoryRouter
    initialEntries={['/dialog/network1/mesh1/' + props.deviceID]}
    initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <Switch>
          <Route
            path="/dialog/:networkId/:meshID/:deviceID"
            render={() => (
              <WifiDeviceDialog
                title="Edit Device Dialog"
                onSave={props.onSave}
                onCancel={props.onCancel}
              />
            )}
          />
          <Route
            path="/dialog/:networkId/:meshID/"
            render={() => (
              <WifiDeviceDialog
                title="Add Device Dialog"
                onSave={props.onSave}
                onCancel={props.onCancel}
              />
            )}
          />
        </Switch>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(cleanup);

describe('<WifiDeviceDialog />', () => {
  beforeEach(() => {
    // eslint-disable-next-line max-len
    MagmaAPIBindings.getWifiByNetworkIdGatewaysByGatewayId.mockResolvedValueOnce(
      RAW_GATEWAY,
    );
  });

  it('no deviceID shows Add Device', async () => {
    const onSave = jest.fn(() => {});
    const onCancel = jest.fn(() => {});

    const {getByText} = render(
      <Wrapper onSave={onSave} onCancel={onCancel} deviceID="" />,
    );
    expect(getByText('Add Device').textContent).toContain('Add Device');
  });

  it('tabs show when deviceID is provided', async () => {
    const onSave = jest.fn(() => {});
    const onCancel = jest.fn(() => {});

    const {getByText} = render(
      <Wrapper onSave={onSave} onCancel={onCancel} deviceID="device1" />,
    );

    await wait();

    expect(
      MagmaAPIBindings.getWifiByNetworkIdGatewaysByGatewayId,
    ).toHaveBeenCalledTimes(1);

    fireEvent.click(getByText('Controller'));
    expect(getByText('Autoupgrade Enabled').textContent).toContain(
      'Autoupgrade Enabled',
    );

    fireEvent.click(getByText('Hardware'));
    expect(getByText('HW ID').textContent).toContain('HW ID');

    fireEvent.click(getByText('Command'));
    expect(getByText('Reboot Device').textContent).toContain('Reboot Device');
  });

  it('cancel button calls onCancel', async () => {
    const onSave = jest.fn(() => {});
    const onCancel = jest.fn(() => {});

    const {getByText} = render(
      <Wrapper onSave={onSave} onCancel={onCancel} deviceID="device1" />,
    );

    await wait();

    fireEvent.click(getByText('Cancel'));

    expect(onCancel.mock.calls.length).toBe(1);
    expect(onSave.mock.calls.length).toBe(0);
  });
});
