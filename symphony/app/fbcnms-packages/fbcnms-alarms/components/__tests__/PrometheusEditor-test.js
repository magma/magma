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
import MuiStylesThemeProvider from '@material-ui/styles/ThemeProvider';
import PrometheusEditor from '../PrometheusEditor';
import defaultTheme from '@fbcnms/ui/theme/default';
import {MemoryRouter} from 'react-router-dom';
import {MuiThemeProvider} from '@material-ui/core/styles';
import {SnackbarProvider} from 'notistack';
import {cleanup, render} from '@testing-library/react';
import {mockApiUtil} from '../../test/testHelpers';
import type {AlertConfig} from '../AlarmAPIType';
import type {GenericRule} from '../RuleInterface';

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
  apiUtil: mockApiUtil(),
  onRuleUpdated: () => {},
  onExit: () => {},
  isNew: false,
  thresholdEditorEnabled: true,
};

test('editing a threshold alert opens the PrometheusEditor with the threshold expression editor enabled', async () => {
  const testThresholdRule: GenericRule<AlertConfig> = {
    severity: '',
    ruleType: '',
    rawRule: {alert: '', expr: 'metric > 123'},
    period: '',
    name: '',
    description: '',
    expression: 'metric > 123',
  };
  const {getByDisplayValue} = render(
    <Wrapper>
      <PrometheusEditor {...commonProps} rule={testThresholdRule} />
    </Wrapper>,
  );
  expect(getByDisplayValue('metric')).toBeInTheDocument();
  expect(getByDisplayValue('123')).toBeInTheDocument();
});

test('editing a non-threshold alert opens the PrometheusEditor with the advanced editor enabled', async () => {
  const nonThresholdRule: GenericRule<AlertConfig> = {
    severity: '',
    ruleType: '',
    rawRule: {alert: '', expr: 'vector(1)'},
    period: '',
    name: '',
    description: '',
    expression: 'vector(1)',
  };
  const {getByDisplayValue} = render(
    <Wrapper>
      <PrometheusEditor {...commonProps} rule={nonThresholdRule} />
    </Wrapper>,
  );
  expect(getByDisplayValue('vector(1)')).toBeInTheDocument();
});
