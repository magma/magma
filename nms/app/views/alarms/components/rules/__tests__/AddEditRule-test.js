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

import * as React from 'react';
import AddEditRule from '../AddEditRule';
import RuleEditorBase from '../RuleEditorBase';
// $FlowFixMe migrated to typescript
import nullthrows from '../../../../../../shared/util/nullthrows';
import {act, fireEvent, render} from '@testing-library/react';
import {alarmTestUtil, renderAsync} from '../../../test/testHelpers';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
import {mockPrometheusRule} from '../../../test/testData';
import {toBaseFields} from '../PrometheusEditor/PrometheusEditor';
// $FlowFixMe migrated to typescript
import type {AlertConfig} from '../../AlarmAPIType';
// $FlowFixMe[cannot-resolve-module] for TypeScript migration
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

describe('Receiver select', () => {
  function assertType<T, I>(value: ?T, shouldBe: Class<I>): I {
    if (value instanceof shouldBe) {
      return value;
    }
    // $FlowFixMe: shouldBe.name does exist
    throw new Error('value is not of type ' + shouldBe.name);
  }

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
    const {getByTestId} = await renderAsync(
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
    const select = getByTestId('select-receiver');
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
    const {getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
    );
    const selectReceiver = getByTestId('select-receiver-input');

    act(() => {
      fireEvent.change(selectReceiver, {target: {value: 'new_receiver'}});
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

    const {getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
    );

    const selectReceiver = getByTestId('select-receiver-input');

    act(() => {
      fireEvent.change(selectReceiver, {target: {value: 'test_receiver'}});
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

    const {getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
      {baseElement: nullthrows(document.body)},
    );

    const selectReceiver = getByTestId('select-receiver-input');

    act(() => {
      fireEvent.change(selectReceiver, {target: {value: 'new_receiver'}});
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
    const {getByTestId} = render(
      <AlarmsWrapper>
        <AddEditRule {...commonProps} />
      </AlarmsWrapper>,
      {baseElement: nullthrows(document.body)},
    );

    const selectReceiver = getByTestId('select-receiver-input');

    act(() => {
      // select option None
      fireEvent.change(selectReceiver, {target: {value: ''}});
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
