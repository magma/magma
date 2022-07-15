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
 */

import * as React from 'react';
import LabelsEditor from '../LabelsEditor';
import {act, fireEvent, render} from '@testing-library/react';

const commonProps = {
  labels: {},
  onChange: jest.fn(),
};

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

test('spaces in the key field are replaced by underscores', () => {
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
      target: {value: 'testKey1 '},
    });
  });
  act(() => {
    fireEvent.change(getByDisplayValue('testKey2'), {
      target: {value: 'testKey2 edited'},
    });
  });
  expect(commonProps.onChange).toHaveBeenCalledWith({
    testKey1_: 'testVal1',
    testKey2_edited: 'testVal2',
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
