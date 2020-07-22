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
import LabelsEditor from '../LabelsEditor';
import {act, cleanup, fireEvent, render} from '@testing-library/react';

const commonProps = {
  labels: {},
  onChange: jest.fn(),
};
afterEach(() => {
  cleanup();
  jest.resetAllMocks();
});

test('clicking the add button adds new textboxes', () => {
  const {getByTestId, queryAllByPlaceholderText} = render(
    <LabelsEditor {...commonProps} />,
  );
  expect(queryAllByPlaceholderText(/name/i).length).toBe(0);
  expect(queryAllByPlaceholderText(/value/i).length).toBe(0);
  act(() => {
    fireEvent.click(getByTestId('add-new-label'));
  });

  expect(queryAllByPlaceholderText(/name/i).length).toBe(1);
  expect(queryAllByPlaceholderText(/value/i).length).toBe(1);
  act(() => {
    fireEvent.click(getByTestId('add-new-label'));
  });
  expect(queryAllByPlaceholderText(/name/i).length).toBe(2);
  expect(queryAllByPlaceholderText(/value/i).length).toBe(2);
});

test('typing into a key field edits the key of a label', () => {
  const {getByDisplayValue} = render(
    <LabelsEditor
      {...commonProps}
      labels={{
        testKey1: 'testVal1',
        testKey2: 'testVal2',
      }}
    />,
  );

  act(() => {
    fireEvent.change(getByDisplayValue('testKey1'), {
      target: {value: 'testKey1-edited'},
    });
  });
  act(() => {
    fireEvent.change(getByDisplayValue('testKey2'), {
      target: {value: 'testKey2-edited'},
    });
  });
  expect(commonProps.onChange).toHaveBeenCalledWith({
    'testKey1-edited': 'testVal1',
    'testKey2-edited': 'testVal2',
  });
});
test('typing into a value field edits the value of a label', () => {
  const {getByDisplayValue} = render(
    <LabelsEditor
      {...commonProps}
      labels={{
        testKey1: 'testVal1',
        testKey2: 'testVal2',
      }}
    />,
  );

  act(() => {
    fireEvent.change(getByDisplayValue('testVal1'), {
      target: {value: 'testVal1-edited'},
    });
  });
  act(() => {
    fireEvent.change(getByDisplayValue('testVal2'), {
      target: {value: 'testVal2-edited'},
    });
  });
  expect(commonProps.onChange).toHaveBeenCalledWith({
    testKey1: 'testVal1-edited',
    testKey2: 'testVal2-edited',
  });
});
test('labels without a key are filtered out', () => {
  const {getByTestId, queryAllByPlaceholderText} = render(
    <LabelsEditor
      {...commonProps}
      labels={{
        testKey1: 'testVal1',
      }}
    />,
  );
  expect(queryAllByPlaceholderText(/name/i).length).toBe(1);
  act(() => {
    fireEvent.click(getByTestId('add-new-label'));
  });
  expect(commonProps.onChange).toHaveBeenLastCalledWith({
    testKey1: 'testVal1',
  });
  act(() => {
    fireEvent.click(getByTestId('add-new-label'));
  });
  expect(queryAllByPlaceholderText(/name/i).length).toBe(3);
  expect(commonProps.onChange).toHaveBeenLastCalledWith({
    testKey1: 'testVal1',
  });
});
test('clicking the reset button resets the label form to its initial state', () => {
  const {getByDisplayValue, getByLabelText, getByTestId} = render(
    <LabelsEditor
      {...commonProps}
      labels={{
        testKey1: 'testVal1',
      }}
    />,
  );
  act(() => {
    fireEvent.click(getByTestId('add-new-label'));
  });

  act(() => {
    fireEvent.change(getByDisplayValue('testKey1'), {
      target: {value: 'testKey1-edited'},
    });
  });

  // ensure state has been changed properly
  expect(commonProps.onChange).toHaveBeenLastCalledWith({
    'testKey1-edited': 'testVal1',
  });
  expect(getByDisplayValue('testKey1-edited')).not.toBeNull();
  // reset
  act(() => {
    fireEvent.click(getByLabelText(/Reset labels/i));
  });
  // ensure state is back to initial state
  expect(commonProps.onChange).toHaveBeenLastCalledWith({
    testKey1: 'testVal1',
  });
  expect(getByDisplayValue('testKey1')).not.toBeNull();
});
test('clicking the remove label button removes a label', () => {
  const {getByLabelText, queryByLabelText} = render(
    <LabelsEditor
      {...commonProps}
      labels={{
        testKey1: 'testVal1',
      }}
    />,
  );
  act(() => {
    fireEvent.click(getByLabelText(/remove label/i));
  });
  expect(queryByLabelText(/name/i)).toBeNull();
  expect(commonProps.onChange).toHaveBeenLastCalledWith({});
});
