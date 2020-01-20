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
import MagmaAPIBindings from '@fbcnms/magma-api';
import Main from '../Main';
import React from 'react';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {cleanup, render, wait} from '@testing-library/react';

jest.mock('@fbcnms/magma-api');
jest.mock('mapbox-gl', () => ({
  Map: () => ({}),
}));
jest.mock('../MapView', () => () => <div>Im the Map!</div>);

jest.mock('../main/Index', () => ({
  __esModule: true,
  default: () => <div>Index</div>,
}));

const Wrapper = props => (
  <MemoryRouter initialEntries={[props.path]} initialIndex={0}>
    <AppContextProvider>{props.children}</AppContextProvider>
  </MemoryRouter>
);

afterEach(cleanup);

const testCases = [
  {
    section: 'nms',
    path: '/nms/mynetwork',
    text: 'Index',
  },
  {
    section: 'admin',
    path: '/admin/settings',
    testId: 'change-password-title',
  },
];

testCases.forEach(testCase => {
  beforeEach(() => {
    MagmaAPIBindings.getNetworks.mockResolvedValueOnce(['mynetwork']);
  });

  afterEach(() => {
    jest.resetAllMocks();
  });

  it(`renders for ${testCase.path} path`, async () => {
    global.CONFIG = {
      appData: {
        enabledFeatures: [],
        tabs: ['nms', 'inventory'],
        user: {
          isSuperUser: false,
        },
      },
      MAPBOX_ACCESS_TOKEN: '',
    };

    const {getByTestId, getByText} = render(
      <Wrapper path={testCase.path}>
        <Main />
      </Wrapper>,
    );

    await wait();

    if (testCase.text) {
      expect(getByText(testCase.text)).toBeInTheDocument();
    } else {
      expect(getByTestId(testCase.testId)).toBeInTheDocument();
    }
  });
});
