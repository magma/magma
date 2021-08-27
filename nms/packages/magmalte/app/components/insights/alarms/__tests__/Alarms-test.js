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

import 'jest-dom/extend-expect';
import * as React from 'react';
import Alarms from '../Alarms';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MagmaAlarmsApiUtil} from '../../../../state/AlarmsApiUtil';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import {cleanup, render} from '@testing-library/react';

jest.mock('@fbcnms/ui/hooks/useSnackbar');
const useSnackbar = require('@fbcnms/ui/hooks/useSnackbar');
const snackbarsMock = {error: jest.fn(), success: jest.fn()};
jest
  .spyOn(useSnackbar, 'useSnackbars')
  .mockReturnValue(jest.fn(() => snackbarsMock));

const useSnackbarsMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useSnackbars')
  .mockReturnValue(useSnackbarsMock);
const useMagmaAPIMock = jest
  .spyOn(require('../../../../../api/useMagmaAPI'), 'default')
  .mockReturnValue({response: []});

const Wrapper = (props: {route: string, children: React.Node}) => (
  <MemoryRouter initialEntries={[props.route || '/alarms']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>{props.children}</SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

const AlarmsWrapper = () => <Alarms apiUtil={MagmaAlarmsApiUtil} />;

afterEach(() => {
  cleanup();
  useMagmaAPIMock.mockClear();
});

describe('react router tests', () => {
  test('/alerts renders the no alerts icon', () => {
    const {getByTestId} = render(
      <Wrapper route={'/alarms'}>
        <AlarmsWrapper />
      </Wrapper>,
    );

    // assert that the 'no alerts' icon is visible
    expect(getByTestId('no-alerts-icon')).toBeInTheDocument();
  });
});

describe('Firing Alerts', () => {
  test('renders currently firing alerts if api returns alerts', () => {
    useMagmaAPIMock.mockReturnValue({
      response: [
        {
          labels: {alertname: '<<TEST ALERT>>', team: '<<TEST TEAM>>'},
          annotations: {description: '<<TEST DESCRIPTION>>'},
        },
      ],
    });

    const {getByTestId, getByText} = render(
      <Wrapper route={'/alerts'}>
        <AlarmsWrapper />
      </Wrapper>,
    );

    // assert that the top level firing alerts header is visible
    expect(getByTestId('firing-alerts')).toBeInTheDocument();
    expect(getByText('<<TEST ALERT>>')).toBeInTheDocument();
    // TODO(andreilee): This has been removed
    // expect(getByText('<<TEST DESCRIPTION>>')).toBeInTheDocument();
  });

  // TODO(andreilee): Fix mock useSnackbars after migrating fbcnms/ui
  // test('if error occurs loading alerts, enqueues error snackbar', async () => {
  //   useMagmaAPIMock.mockReturnValueOnce({
  //     error: {message: 'an error occurred'},
  //   });

  //   const snackbarsMock = {error: jest.fn(), success: jest.fn()};
  //   jest
  //     .spyOn(useSnackbar, 'useSnackbars')
  //     .mockImplementation(jest.fn(() => snackbarsMock));

  //   render(
  //     <Wrapper route={'/alerts'}>
  //       <AlarmsWrapper />
  //     </Wrapper>,
  //   );

  //   await wait();

  //   expect(useSnackbarsMock).toHaveBeenCalled();
  //   expect(snackbarsMock.error).toHaveBeenCalled();
  // });
});
