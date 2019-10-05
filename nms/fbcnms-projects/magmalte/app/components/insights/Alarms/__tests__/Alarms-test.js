/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import 'jest-dom/extend-expect';
import * as React from 'react';
import Alarms from '../Alarms';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {Route} from 'react-router-dom';
import {SnackbarProvider} from 'notistack';
import {cleanup, render} from '@testing-library/react';

jest.mock('../../../../common/useMagmaAPI');
jest.mock('@fbcnms/ui/hooks/useSnackbar');
const useMagmaAPI = require('../../../../common/useMagmaAPI');
const useSnackbar = require('@fbcnms/ui/hooks/useSnackbar');
const useMagmaAPIMock = jest
  .spyOn(useMagmaAPI, 'default')
  .mockReturnValue({response: []});
jest.mock('@fbcnms/ui/hooks/useRouter');
const useRouter = require('@fbcnms/ui/hooks/useRouter');

jest.spyOn(useRouter, 'default').mockReturnValue({
  match: {params: {networkId: ''}, path: '/', url: '/'},
  relativePath: p => p,
});

const Wrapper = (props: {route: string, children: React.Node}) => (
  <MemoryRouter initialEntries={[props.route || '/alarms']} initialIndex={0}>
    <MuiThemeProvider theme={defaultTheme}>
      <MuiStylesThemeProvider theme={defaultTheme}>
        <SnackbarProvider>{props.children}</SnackbarProvider>
      </MuiStylesThemeProvider>
    </MuiThemeProvider>
  </MemoryRouter>
);

afterEach(() => {
  cleanup();
  useMagmaAPIMock.mockClear();
});

describe('react router tests', () => {
  test('/alerts renders the firing alerts panel', () => {
    const {getByTestId} = render(
      <Wrapper route={'/alarms'}>
        <Route path="/alarms" component={Alarms} />,
      </Wrapper>,
    );
    // assert that the top level firing alerts header is visible
    expect(getByTestId('firing-alerts')).toBeInTheDocument();
  });

  test('/alerts/new_alert renders AddEditAlert', () => {
    const {getByText} = render(
      <Wrapper route={'/alarms/new_alert'}>
        <Route path="/alarms" component={Alarms} />,
      </Wrapper>,
    );
    expect(getByText('New Alert')).toBeInTheDocument();
  });

  test('/alerts/edit_alerts renders EditAllAlerts', () => {
    const {getByText} = render(
      <Wrapper route={'/alarms/edit_alerts'}>
        <Route path="/alarms" component={Alarms} />,
      </Wrapper>,
    );
    expect(getByText('Edit Alerts')).toBeInTheDocument();
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

    const {getByText} = render(
      <Wrapper route={'/alarms'}>
        <Alarms />
      </Wrapper>,
    );

    expect(getByText('<<TEST ALERT>>')).toBeInTheDocument();
    expect(getByText('<<TEST DESCRIPTION>>')).toBeInTheDocument();
    expect(getByText('<<TEST TEAM>>')).toBeInTheDocument();
  });

  test('if an error occurs while loading alerts, enqueues an error snackbar', () => {
    useMagmaAPIMock.mockReturnValue({
      error: {message: 'an error occurred'},
    });
    const enqueueSnackbarMock = jest.fn();
    jest
      .spyOn(useSnackbar, 'useEnqueueSnackbar')
      .mockReturnValueOnce(enqueueSnackbarMock);

    render(
      <Wrapper route={'/alarms'}>
        <Alarms />
      </Wrapper>,
    );

    expect(enqueueSnackbarMock).toHaveBeenCalled();
  });
});
