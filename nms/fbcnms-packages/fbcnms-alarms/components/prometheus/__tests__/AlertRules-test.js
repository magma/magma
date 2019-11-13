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
import AlertRules from '../AlertRules';
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import type {ApiUtil} from '../../AlarmsApi';

jest.mock('@fbcnms/ui/hooks/useSnackbar');
jest.mock('@fbcnms/ui/hooks/useRouter');

afterEach(() => {
  cleanup();
  jest.clearAllMocks();
});

const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);
jest
  .spyOn(require('@fbcnms/ui/hooks/useRouter'), 'default')
  .mockReturnValue({match: {params: {networkId: 'test'}}});

/**
 * I don't understand how to properly type this mock so using any for now.
 * The consuming code is all strongly typed, this shouldn't be much of an issue.
 */
// eslint-disable-next-line flowtype/no-weak-types
const useMagmaAPIMock = jest.fn<any, any>(() => ({
  isLoading: false,
  response: [],
  error: null,
}));
const apiMock = jest.fn();

// TextField select is difficult to test so replace it with an Input
jest.mock('@material-ui/core/TextField', () => {
  const Input = require('@material-ui/core/Input').default;
  return ({children: _, InputProps: __, label, ...props}) => (
    <label>
      {label}
      <Input {...props} />
    </label>
  );
});

const axiosMock = jest
  .spyOn(require('axios'), 'default')
  .mockImplementation(jest.fn(() => Promise.resolve({data: {}})));

function Wrapper(props: {route?: string, children: React.Node}) {
  return (
    <MemoryRouter initialEntries={[props.route || '/']} initialIndex={0}>
      <MuiThemeProvider theme={defaultTheme}>
        <MuiStylesThemeProvider theme={defaultTheme}>
          <SnackbarProvider>{props.children}</SnackbarProvider>
        </MuiStylesThemeProvider>
      </MuiThemeProvider>
    </MemoryRouter>
  );
}

const commonProps = {
  apiUrls: mockApiUrls(),
  apiUtil: mockApiUtil(),
};

test('renders rules returned by api', () => {
  useMagmaAPIMock.mockReturnValueOnce({
    response: mockRules(),
    error: null,
    isLoading: false,
  });
  const {getByText} = render(
    <Wrapper>
      <AlertRules {...commonProps} />
    </Wrapper>,
  );
  expect(getByText('<<test>>')).toBeInTheDocument();
  expect(getByText('up == 0')).toBeInTheDocument();
});

test('clicking the add alert icon displays the AddEditAlert view', () => {
  useMagmaAPIMock.mockReturnValueOnce({
    response: mockRules(),
    error: null,
    isLoading: false,
  });
  const {queryByTestId, getByTestId} = render(
    <Wrapper>
      <AlertRules {...commonProps} />
    </Wrapper>,
  );
  expect(queryByTestId('add-edit-alert')).not.toBeInTheDocument();
  // click the add alert rule fab
  act(() => {
    fireEvent.click(getByTestId('add-edit-alert-button'));
  });
  expect(getByTestId('add-edit-alert')).toBeInTheDocument();
});

test('clicking close button when AddEditAlert is open closes the panel', () => {
  useMagmaAPIMock.mockReturnValueOnce({
    response: [],
    error: null,
    isLoading: false,
  });
  const {queryByTestId, getByTestId, getByText} = render(
    <Wrapper>
      <AlertRules {...commonProps} />
    </Wrapper>,
  );
  expect(queryByTestId('add-edit-alert')).not.toBeInTheDocument();
  // click the add alert rule fab
  act(() => {
    fireEvent.click(getByTestId('add-edit-alert-button'));
  });
  expect(getByTestId('add-edit-alert')).toBeInTheDocument();
  act(() => {
    fireEvent.click(getByText(/close/i));
  });
  expect(queryByTestId('add-edit-alert')).not.toBeInTheDocument();
});

test.todo(
  'clicking the "edit" button in the table menu opens AddEditAlert for that alert',
);

/**
 * Test AlertRules' integration with AddEditAlert. It passes in a specific
 * columnStruct object so we need to test that this works properly.
 */
describe('AddEditAlert > Prometheus Editor', () => {
  test('Filling the form and clicking Add will post to the endpoint', async () => {
    useMagmaAPIMock.mockReturnValueOnce({
      response: mockRules(),
      error: null,
      isLoading: false,
    });
    const {getByText, getByTestId, getByLabelText} = render(
      <Wrapper>
        <AlertRules {...commonProps} />
      </Wrapper>,
    );
    act(() => {
      fireEvent.click(getByTestId('add-edit-alert-button'));
    });
    act(() => {
      fireEvent.change(getByLabelText(/rule name/i), {
        target: {value: '<<ALERTNAME>>'},
      });
    });
    act(() => {
      fireEvent.change(getByLabelText(/severity/i), {
        target: {value: 'minor'},
      });
    });
    act(() => {
      fireEvent.change(getByLabelText(/duration/i), {
        target: {value: '1'},
      });
    });
    act(() => {
      fireEvent.change(getByLabelText(/unit/i), {
        target: {value: 'm'},
      });
    });
    act(() => {
      fireEvent.change(getByLabelText(/expression/i), {
        target: {value: 'vector(1)'},
      });
    });
    // This triggers an async call so must be awaited
    await act(async () => {
      fireEvent.click(getByText(/add/i));
    });
    expect(apiMock).toHaveBeenLastCalledWith({
      networkId: 'test',
      rule: {
        alert: '<<ALERTNAME>>',
        labels: {
          severity: 'minor',
        },
        for: '1m',
        expr: 'vector(1)',
      },
    });
  });

  test('a snackbar is enqueued if adding a rule fails', async () => {
    const enqueueMock = jest.fn();
    jest
      .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
      .mockReturnValue(enqueueMock);

    axiosMock.mockRejectedValueOnce({
      response: {
        status: 500,
        data: {message: 'an error message'},
      },
    });
    const {getByText, getByTestId} = render(
      <Wrapper>
        <AlertRules {...commonProps} />
      </Wrapper>,
    );
    act(() => {
      fireEvent.click(getByTestId('add-edit-alert-button'));
    });

    await act(async () => {
      fireEvent.click(getByText(/add/i));
    });

    expect(enqueueMock).toHaveBeenCalled();
  });
});

function mockRules() {
  return [
    {
      alert: '<<test>>',
      annotations: {},
      labels: {},
      expr: 'up == 0',
      for: '1m',
    },
  ];
}

function mockApiUrls() {
  return {
    viewFiringAlerts: () => '/viewFiringAlerts',
    alertConfig: () => '/alertConfig',
    updateAlertConfig: () => '/updateAlertConfig',
    bulkAlertConfig: () => '/bulkAlertConfig',
    receiverConfig: () => '/receiverConfig',
    // get count of matching metrics,
    viewMatchingAlerts: () => '/viewMatchingAlerts',
    receiverUpdate: () => '/receiverUpdate',
    routeConfig: () => '/routeConfig',
    viewSilences: () => '/viewSilences',
    viewRoutes: () => '/viewRoutes',
    viewReceivers: () => '/viewReceivers',
  };
}

function mockApiUtil(): ApiUtil {
  return {
    useAlarmsApi: useMagmaAPIMock,
    viewFiringAlerts: apiMock,
    viewMatchingAlerts: apiMock,
    createAlertRule: apiMock,
    editAlertRule: apiMock,
    getAlertRules: apiMock,
    deleteAlertRule: apiMock,
    getSuppressions: apiMock,
    getReceivers: apiMock,
    getRoutes: apiMock,
  };
}
