/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */
import 'jest-dom/extend-expect';
import * as React from 'react';
import AddEditRule from '../AddEditRule';
import RuleEditorBase from '../RuleEditorBase';
import nullthrows from '@fbcnms/alarms/util/nullthrows';
import {act, cleanup, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil, renderAsync} from '../../../test/testHelpers';
import {assertType} from '@fbcnms/alarms/util/assert';
import {mockPrometheusRule} from '../../../test/testData';
import {toBaseFields} from '../PrometheusEditor/PrometheusEditor';
import type {AlertConfig} from '../../AlarmAPIType';
import type {RuleEditorProps} from '../RuleInterface';

const mockRuleMap = {
  mock: {
    RuleEditor: MockRuleEditor,
    deleteRule: jest.fn(),
    getRules: jest.fn(),
    friendlyName: 'mock',
  },
};
const {apiUtil, AlarmsWrapper} = alarmTestUtil({ruleMap: mockRuleMap});

const commonProps = {
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
      .spyOn(apiUtil, 'getReceivers')
      .mockImplementation(() => [{name: 'test_receiver'}]);
    jest.spyOn(apiUtil, 'getRouteTree').mockReturnValue({
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
      <AlarmsWrapper>
        <AddEditRule
          {...commonProps}
          initialConfig={mockPrometheusRule({
            name: 'TESTRULE',
            ruleType: 'mock',
          })}
        />
      </AlarmsWrapper>,
    );
    const select = getByLabelText(/send notification to/i);
    expect(select.textContent).toBe('test_receiver');
  });

  test('selecting a receiver sets the value in the select box', () => {
    mockUseAlarms();
    jest
      .spyOn(apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(apiUtil, 'getRouteTree').mockReturnValue({
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
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
    );
    act(() => {
      fireEvent.mouseDown(getByLabelText(/send notification to/i));
    });
    act(() => {
      fireEvent.click(getByText('new_receiver'));
    });

    const receiverInput = assertType(
      getByTestId('select-receiver-input'),
      HTMLInputElement,
    );
    expect(receiverInput.value).toBe('new_receiver');
  });

  test('setting a receiver adds a new route', async () => {
    mockUseAlarms();
    jest
      .spyOn(apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}]);
    jest.spyOn(apiUtil, 'getRouteTree').mockReturnValue({
      receiver: 'network_base_route',
      routes: [],
    });

    const editRouteTreeMock = jest.spyOn(apiUtil, 'editRouteTree');

    const {getByLabelText, getByText, getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
    );

    act(() => {
      fireEvent.mouseDown(getByLabelText(/send notification to/i));
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
      .spyOn(apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(apiUtil, 'getRouteTree').mockReturnValue({
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

    const editRouteTreeMock = jest.spyOn(apiUtil, 'editRouteTree');

    const {getByLabelText, getByText, getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
      {baseElement: nullthrows(document.body)},
    );

    act(() => {
      fireEvent.mouseDown(getByLabelText(/send notification to/i));
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
      .spyOn(apiUtil, 'getReceivers')
      .mockReturnValue([{name: 'test_receiver'}, {name: 'new_receiver'}]);
    jest.spyOn(apiUtil, 'getRouteTree').mockReturnValue({
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
    const editRouteTreeMock = jest.spyOn(apiUtil, 'editRouteTree');
    const {getByLabelText, getByText, getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
      {baseElement: nullthrows(document.body)},
    );

    act(() => {
      fireEvent.mouseDown(getByLabelText(/send notification to/i));
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

function MockRuleEditor(props: RuleEditorProps<AlertConfig>) {
  const {isNew, rule} = props;
  return (
    <RuleEditorBase
      isNew={isNew}
      onSave={jest.fn()}
      onExit={jest.fn()}
      onChange={jest.fn()}
      initialState={toBaseFields(rule)}>
      <span />
    </RuleEditorBase>
  );
}

function mockUseAlarms() {
  jest.spyOn(apiUtil, 'useAlarmsApi').mockImplementation((fn, params) => {
    return {
      response: fn(params),
    };
  });
}
