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
import AddEditAlert from '../AddEditAlert';
import RuleEditorBase from '../RuleEditorBase';
import {SymphonyWrapper} from '@fbcnms/test/testHelpers';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import {mockApiUtil, renderAsync} from '../../../test/testHelpers';
import {mockPrometheusRule} from '../../../test/data';

import type {RuleEditorProps} from '../../RuleInterface';

const commonProps = {
  apiUtil: mockApiUtil(),
  ruleMap: {
    mock: {
      RuleEditor: MockRuleEditor,
      deleteRule: jest.fn(),
      getRules: jest.fn(),
      friendlyName: 'mock',
    },
  },
  isNew: false,
  initialConfig: mockPrometheusRule({
    name: 'TESTRULE',
    ruleType: 'mock',
  }),
  onExit: jest.fn(),
};

afterEach(() => {
  cleanup();
  jest.resetAllMocks();
});

describe('Receiver select', () => {
  test('a rule with a receiver selected sets the receiver select value', async () => {
    mockUseAlarms();
    jest
      .spyOn(commonProps.apiUtil, 'getReceivers')
      .mockImplementation(() => [{name: 'test_receiver'}]);
    jest.spyOn(commonProps.apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [
        {
          receiver: 'test_receiver',
          match: {
            alertname: 'TESTRULE',
          },
        },
      ],
    });
    const {getByLabelText} = await renderAsync(
      <SymphonyWrapper>
        <AddEditAlert
          {...commonProps}
          initialConfig={mockPrometheusRule({
            name: 'TESTRULE',
            ruleType: 'mock',
          })}
        />
      </SymphonyWrapper>,
    );
    const select = getByLabelText(/send notification to/i);
    expect(select.textContent).toBe('test_receiver');
  });

  test('selecting a receiver sets the value in the select box', () => {
    mockUseAlarms();
    jest
      .spyOn(commonProps.apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(commonProps.apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [
        {
          receiver: 'test_receiver',
          match: {
            alertname: 'TESTRULE',
          },
        },
      ],
    });
    const {getByTestId, getByLabelText, getByText} = render(
      <SymphonyWrapper>
        <AddEditAlert {...commonProps} />
      </SymphonyWrapper>,
    );
    act(() => {
      fireEvent.click(getByLabelText(/send notification to/i));
    });
    act(() => {
      fireEvent.click(getByText('new_receiver'));
    });
    expect(getByTestId('select-receiver-input').value).toBe('new_receiver');
  });

  test('setting a receiver adds a new route', async () => {
    mockUseAlarms();
    jest
      .spyOn(commonProps.apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}]);
    jest.spyOn(commonProps.apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [],
    });

    const editRouteTreeMock = jest.spyOn(commonProps.apiUtil, 'editRouteTree');

    const {getByLabelText, getByText, getByTestId} = render(
      <SymphonyWrapper>
        <AddEditAlert {...commonProps} />
      </SymphonyWrapper>,
    );

    act(() => {
      fireEvent.click(getByLabelText(/send notification to/i));
    });
    act(() => {
      fireEvent.click(getByText('test_receiver'));
    });
    await act(async () => {
      fireEvent.submit(getByTestId('editor-form'));
    });
    expect(editRouteTreeMock).toHaveBeenCalledWith({
      networkId: undefined,
      route: {
        receiver: 'network_base_route',
        routes: [
          {
            receiver: 'test_receiver',
            match: {
              alertname: 'TESTRULE',
            },
          },
        ],
      },
    });
  });
  test('selecting a new receiver updates an existing route', async () => {
    mockUseAlarms();
    jest
      .spyOn(commonProps.apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(commonProps.apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [
        {
          receiver: 'test_receiver',
          match: {
            alertname: 'TESTRULE',
          },
        },
      ],
    });

    const editRouteTreeMock = jest.spyOn(commonProps.apiUtil, 'editRouteTree');

    const {getByLabelText, getByText, getByTestId} = render(
      <SymphonyWrapper>
        <AddEditAlert {...commonProps} />
      </SymphonyWrapper>,
      {baseElement: document.body},
    );

    act(() => {
      fireEvent.click(getByLabelText(/send notification to/i));
    });
    await act(async () => {
      fireEvent.click(getByText('new_receiver'));
    });
    await act(async () => {
      fireEvent.submit(getByTestId('editor-form'));
    });

    expect(editRouteTreeMock).toHaveBeenCalledWith({
      networkId: undefined,
      route: {
        receiver: 'network_base_route',
        routes: [
          {
            receiver: 'new_receiver',
            match: {
              alertname: 'TESTRULE',
            },
          },
        ],
      },
    });
  });
  test('un-selecting receiver removes the existing route', async () => {
    mockUseAlarms();
    jest
      .spyOn(commonProps.apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(commonProps.apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [
        {
          receiver: 'test_receiver',
          match: {
            alertname: 'TESTRULE',
          },
        },
      ],
    });
    const editRouteTreeMock = jest.spyOn(commonProps.apiUtil, 'editRouteTree');
    const {getByLabelText, getByText, getByTestId} = render(
      <SymphonyWrapper>
        <AddEditAlert {...commonProps} />
      </SymphonyWrapper>,
      {baseElement: document.body},
    );

    act(() => {
      fireEvent.click(getByLabelText(/send notification to/i));
    });
    await act(async () => {
      fireEvent.click(getByText('None'));
    });
    await act(async () => {
      fireEvent.submit(getByTestId('editor-form'));
    });

    expect(editRouteTreeMock).toHaveBeenCalledWith({
      networkId: undefined,
      route: {
        receiver: 'network_base_route',
        routes: [],
      },
    });
  });
});

function MockRuleEditor(props: RuleEditorProps<{}>) {
  const {apiUtil, isNew, rule} = props;
  return (
    <RuleEditorBase
      apiUtil={apiUtil}
      isNew={isNew}
      onSave={jest.fn()}
      onExit={jest.fn()}
      onChange={jest.fn()}
      rule={rule}>
      <span />
    </RuleEditorBase>
  );
}

function mockUseAlarms() {
  jest
    .spyOn(commonProps.apiUtil, 'useAlarmsApi')
    .mockImplementation((fn, params) => {
      return {
        response: fn(params),
      };
    });
}
