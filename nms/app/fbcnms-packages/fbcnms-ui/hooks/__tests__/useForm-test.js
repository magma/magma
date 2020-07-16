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
import {cleanup} from '@testing-library/react';
import {act as hooksAct, renderHook} from '@testing-library/react-hooks';
import {useForm} from '../../hooks';

jest.useFakeTimers();
afterEach(() => {
  cleanup();
  jest.clearAllMocks();
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
      textEventHandler({
        target: ({value: 'test text'}: $Shape<HTMLInputElement>),
      });
    });
    hooksAct(() => {
      result.current.addListItem('list', {test: 1});
    });
    expect(result.current.formState).toMatchObject({
      text: 'test text',
      list: [{test: 1}],
    });
  });

  type TestState = {|
    text: string,
    list: Array<{}>,
  |};
  test('editing a field calls onFormUpdated', () => {
    const state: TestState = {
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
      textEventHandler({
        target: ({value: 'test text'}: $Shape<HTMLInputElement>),
      });
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
