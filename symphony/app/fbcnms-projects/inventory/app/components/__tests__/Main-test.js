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
import Main from '../Main';
import React from 'react';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {cleanup, render} from '@testing-library/react';

jest.mock('mapbox-gl', () => ({
  Map: () => ({}),
}));
jest.mock('../map/MapView', () => () => <div>Im the Map!</div>);
jest.mock('../projects/ProjectsMap', () => () => (
  <div>Im the ProjectsMap!</div>
));
jest.mock('@fbcnms/magmalte/app/components/MapView', () => () => (
  <div>Im the Map!</div>
));

jest.mock('@fbcnms/magmalte/app/components/Main', () => ({
  __esModule: true,
  default: () => <div>MagmaMain</div>,
}));
jest.mock('../Inventory', () => ({
  __esModule: true,
  default: () => <div>Inventory</div>,
}));
jest.mock('../automation/Automation', () => ({
  __esModule: true,
  default: () => <div>Automation</div>,
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
    path: '/nms',
    text: 'MagmaMain',
  },
  {
    section: 'inventory',
    path: '/inventory',
    text: 'Inventory',
  },
  {
    section: 'automation',
    path: '/automation',
    text: 'Automation',
  },
];

testCases.forEach(testCase => {
  it(`renders for ${testCase.path} path`, () => {
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

    const {getByText} = render(
      <Wrapper path={testCase.path}>
        <Main />
      </Wrapper>,
    );

    expect(getByText(testCase.text)).toBeInTheDocument();
  });
});
