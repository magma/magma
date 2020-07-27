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
 * @flow
 * @format
 */

import 'jest-dom/extend-expect';
import FBCMobileAppQRCode from '../FBCMobileAppQRCode';
import React from 'react';
import {TestWrapper} from '../testHelpers/index';
import {act, cleanup, render} from '@testing-library/react';

jest.mock('@fbcnms/ui/hooks');
const {useAxios}: any = require('@fbcnms/ui/hooks');

jest.mock('../generateQRCode');
const {default: generateQRCode}: any = require('../generateQRCode');

afterEach(cleanup);

test('makes request to endpoint when rendered', () => {
  useAxios.mockImplementation(() => ({
    data: {test: 'test'},
  }));
  render(<FBCMobileAppQRCode endpoint="/qrtest" />, {wrapper: TestWrapper});
  expect(useAxios).toHaveBeenCalledWith({
    method: 'GET',
    url: '/qrtest',
  });
});

test('shows loading text when loading', () => {
  useAxios.mockImplementation(() => ({
    isLoading: true,
    response: null,
  }));
  const {getByTestId} = render(<FBCMobileAppQRCode endpoint="/qrtest" />, {
    wrapper: TestWrapper,
  });
  expect(getByTestId('loading')).toBeInTheDocument();
});

test('shows error text if axios returns an error', () => {
  useAxios.mockImplementation(() => ({
    error: new Error(),
    isLoading: false,
    response: null,
  }));
  const {getByTestId} = render(<FBCMobileAppQRCode endpoint="/qrtest" />, {
    wrapper: TestWrapper,
  });
  expect(getByTestId('error-message')).toBeInTheDocument();
});

test('renders a qr code in an image tag', async () => {
  useAxios.mockImplementation(() => ({
    isLoading: false,
    response: {
      data: {
        test: 'test',
      },
    },
  }));
  generateQRCode.mockResolvedValue('data:image/jpeg;base64:123');

  let result: any = {};
  await act(async () => {
    result = await render(<FBCMobileAppQRCode endpoint="/qrtest" />, {
      wrapper: TestWrapper,
    });
  });
  const qrCode = result.container.querySelector('img[src]');
  expect(qrCode).toBeInTheDocument();
  expect(qrCode.src).toBe('data:image/jpeg;base64:123');
});
