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
import {cleanup} from '@testing-library/react';
import {act as hooksAct, renderHook} from '@testing-library/react-hooks';
import {useForm, useLoadRules} from '../hooks';
import type {GenericRule, RuleInterface} from '../rules/RuleInterface';

jest.useFakeTimers();
afterEach(() => {
  cleanup();
  jest.clearAllMocks();
});

const enqueueSnackbarMock = jest.fn();
jest
  .spyOn(require('@fbcnms/ui/hooks/useSnackbar'), 'useEnqueueSnackbar')
  .mockReturnValue(enqueueSnackbarMock);

jest.spyOn(require('@fbcnms/ui/hooks/useRouter'), 'default').mockReturnValue({
  match: {
    params: {
      networkId: 'test',
    },
  },
});

describe('useLoadRules hook', () => {
  test('calls all getRules functions and merges their results', async () => {
    // return 2 rules from prometheus and one from events
    const prometheusMock = jest.fn(() =>
      Promise.resolve([mockRule(), mockRule()]),
    );
    const eventsMock = jest.fn(() => Promise.resolve([mockRule()]));

    const ruleMap = {
      prometheus: mockRuleInterface({getRules: prometheusMock}),
      events: mockRuleInterface({getRules: eventsMock}),
    };
    const {result} = await renderHookAsync(() =>
      useLoadRules({ruleMap, lastRefreshTime: ''}),
    );

    expect(prometheusMock).toHaveBeenCalled();
    expect(eventsMock).toHaveBeenCalled();
    expect(result.current.rules.length).toBe(3);
  });

  test('if a call errors, a snackbar is enqueued', async () => {
    jest.spyOn(console, 'error').mockImplementationOnce(jest.fn());
    const prometheusMock = jest.fn(() => Promise.resolve([]));
    const eventsMock = jest.fn(() =>
      Promise.reject(new Error('cannot load events')),
    );
    const ruleMap = {
      prometheus: mockRuleInterface({getRules: prometheusMock}),
      events: mockRuleInterface({getRules: eventsMock}),
    };
    await renderHookAsync(() => useLoadRules({ruleMap, lastRefreshTime: ''}));
    expect(prometheusMock).toHaveBeenCalled();
    expect(eventsMock).toHaveBeenCalled();
    expect(enqueueSnackbarMock).toHaveBeenCalled();
  });

  test('if a call is cancelled or errors, other calls still complete', async () => {
    jest.spyOn(console, 'error').mockImplementationOnce(jest.fn());
    const prometheusMock = jest.fn(() =>
      Promise.resolve([mockRule(), mockRule()]),
    );
    const eventsMock = jest.fn(() =>
      Promise.reject(new Error('cannot load events')),
    );
    const ruleMap = {
      prometheus: mockRuleInterface({getRules: prometheusMock}),
      events: mockRuleInterface({getRules: eventsMock}),
    };
    const {result} = await renderHookAsync(() =>
      useLoadRules({ruleMap, lastRefreshTime: ''}),
    );
    expect(prometheusMock).toHaveBeenCalled();
    expect(eventsMock).toHaveBeenCalled();
    expect(result.current.rules).toHaveLength(2);
  });
});

describe('useForm hook', () => {
  test('formState contains form initial state', () => {
    const {result} = renderHook(() =>
      useForm({
        initialState: {
          test: 1,
          nested: {x: 1, y: 2},
        },
      }),
    );
    expect(result.current.formState).toMatchObject({
      test: 1,
      nested: {x: 1, y: 2},
    });
  });
  test('addListItem immutably adds a list item', () => {
    const initialArray = [{test: 1}];
    // this one shouldn't be modified or reassigned
    const untouchedArray = [];
    const {result} = renderHook(() =>
      useForm({
        initialState: {
          list: initialArray,
          list_untouched: untouchedArray,
        },
      }),
    );
    expect(result.current.formState.list.length).toBe(1);
    expect(result.current.formState.list).toBe(initialArray);
    hooksAct(() => {
      result.current.addListItem('list', {test: 2});
    });
    /**
     * ensure that list length changed, that equality was broken,
     * and that the new item was added
     */
    expect(result.current.formState.list.length).toBe(2);
    expect(result.current.formState.list).not.toBe(initialArray);
    expect(result.current.formState.list).toMatchObject([{test: 1}, {test: 2}]);
    // nothing else should've been modified
    expect(result.current.formState.list_untouched).toBe(untouchedArray);
  });
  test('updateListItem immutably updates a list item', () => {
    const initialArray = [{test: 1}, {test: 2}];
    // this one shouldn't be modified or reassigned
    const untouchedArray = [];
    const {result} = renderHook(() =>
      useForm({
        initialState: {
          list: initialArray,
          list_untouched: untouchedArray,
        },
      }),
    );
    expect(result.current.formState.list.length).toBe(2);
    expect(result.current.formState.list).toBe(initialArray);
    expect(result.current.formState.list).toMatchObject([{test: 1}, {test: 2}]);
    hooksAct(() => {
      result.current.updateListItem('list', 1, {test: 4});
    });
    /**
     * ensure that list length changed, that equality was broken,
     * and that the new item was added
     */
    expect(result.current.formState.list.length).toBe(2);
    expect(result.current.formState.list).not.toBe(initialArray);
    expect(result.current.formState.list).toMatchObject([{test: 1}, {test: 4}]);
    // nothing else should've been modified
    expect(result.current.formState.list_untouched).toBe(untouchedArray);
  });

  test('addListItem will create a list if it does not exist already', () => {
    const state: {optionalList?: Array<{test: number}>} = {};
    const {result} = renderHook(() => useForm({initialState: state}));
    expect(result.current.formState.optionalList).not.toBeDefined();
    hooksAct(() => {
      result.current.addListItem('optionalList', {test: 1});
    });
    expect(result.current.formState.optionalList).toBeDefined();
    expect(result.current.formState.optionalList).toMatchObject([{test: 1}]);
  });

  test('editing fields modifying lists should not conflict', () => {
    const state = {
      text: '',
      list: [],
    };
    const {result} = renderHook(() => useForm({initialState: state}));
    const textEventHandler = result.current.handleInputChange(val => ({
      text: val,
    }));
    hooksAct(() => {
      textEventHandler({target: {value: 'test text'}});
    });
    hooksAct(() => {
      result.current.addListItem('list', {test: 1});
    });
    expect(result.current.formState).toMatchObject({
      text: 'test text',
      list: [{test: 1}],
    });
  });

  test('editing a field calls onFormUpdated', () => {
    const state = {
      text: '',
      list: [],
    };
    const onFormUpdatedMock = jest.fn();
    const {result} = renderHook(() =>
      useForm({initialState: state, onFormUpdated: onFormUpdatedMock}),
    );
    const textEventHandler = result.current.handleInputChange(val => ({
      text: val,
    }));
    hooksAct(() => {
      textEventHandler({target: {value: 'test text'}});
    });
    expect(onFormUpdatedMock).toHaveBeenCalledWith({
      list: [],
      text: 'test text',
    });
  });

  test('calling updateFormState calls onFormUpdated', () => {
    const state = {
      text: '',
      list: [],
    };
    const onFormUpdatedMock = jest.fn();
    const {result} = renderHook(() =>
      useForm({initialState: state, onFormUpdated: onFormUpdatedMock}),
    );

    hooksAct(() => {
      result.current.updateFormState({
        text: 'test text',
      });
    });
    expect(onFormUpdatedMock).toHaveBeenCalledWith({
      list: [],
      text: 'test text',
    });
  });
});

function mockRuleInterface(
  merge?: $Shape<RuleInterface<{}>>,
): RuleInterface<{}> {
  return Object.assign(
    {
      friendlyName: '',
      RuleEditor: MockComponent,
      RuleViewer: MockComponent,
      getRules: _ => Promise.resolve([]),
      deleteRule: _ => Promise.resolve(),
    },
    merge || {},
  );
}

function mockRule(): GenericRule<{}> {
  return {
    severity: '',
    name: '',
    description: '',
    period: '',
    expression: '',
    ruleType: '',
    rawRule: {},
  };
}

function MockComponent() {
  return <span />;
}

// eslint-disable-next-line flowtype/no-weak-types
async function renderHookAsync(renderFn): any {
  let response;
  await hooksAct(async () => {
    response = await renderHook(renderFn);
  });
  return response;
}
