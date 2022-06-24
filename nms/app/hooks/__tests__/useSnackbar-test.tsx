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

import * as notistack from 'notistack';
import React from 'react';

import {renderHook} from '@testing-library/react-hooks';
import {useSnackbar} from '../index';

jest.mock('@material-ui/core/Slide', () => () => <div />);
jest.mock('notistack');

it('renders without crashing', () => {
  const mockEnqueueSnackbar = jest.fn().mockReturnValue('key');
  const mockCloseSnackbar = jest.fn();
  jest.spyOn(notistack, 'useSnackbar').mockImplementation(() => ({
    enqueueSnackbar: mockEnqueueSnackbar,
    closeSnackbar: mockCloseSnackbar,
  }));

  const {rerender} = renderHook(
    message => useSnackbar(message, {variant: 'error'}, true),
    {initialProps: 'test1'},
  );

  expect(mockEnqueueSnackbar).toHaveBeenCalledTimes(1);
  expect(mockCloseSnackbar).toHaveBeenCalledTimes(0);

  rerender('test2');
  expect(mockEnqueueSnackbar).toHaveBeenCalledTimes(2);
  expect(mockCloseSnackbar).toHaveBeenCalledTimes(0);
});

it('dismisses previous', () => {
  const mockEnqueueSnackbar = jest.fn().mockReturnValue('key');
  const mockCloseSnackbar = jest.fn();
  jest.spyOn(notistack, 'useSnackbar').mockImplementation(() => ({
    enqueueSnackbar: mockEnqueueSnackbar,
    closeSnackbar: mockCloseSnackbar,
  }));

  const {rerender} = renderHook(
    message => useSnackbar(message, {variant: 'error'}, true, true),
    {initialProps: 'test1'},
  );

  expect(mockEnqueueSnackbar).toHaveBeenCalledTimes(1);
  expect(mockCloseSnackbar).toHaveBeenCalledTimes(0);

  rerender('test2');

  expect(mockEnqueueSnackbar).toHaveBeenCalledTimes(2);
  expect(mockCloseSnackbar).toHaveBeenCalledTimes(1);
});
