/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('../../common/RelayEnvironment');

import 'jest-dom/extend-expect';
import Main from '../Main';
import React from 'react';
import RelayEnvironment from '../../common/RelayEnvironment';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';
import {MemoryRouter} from 'react-router-dom';
import {MockPayloadGenerator} from 'relay-test-utils';
import {cleanup, render, waitForElement} from '@testing-library/react';

const MOCK_RESOLVER = {
  User() {
    return {
      id: '1',
      authID: 'mockuser@test.ing',
      email: 'mockuser@test.ing',
      firstName: 'mock',
      lastName: 'user',
    };
  },
  Permissions() {
    return {
      canWrite: 'true',
      adminPolicy: {
        access: {
          isAllowed: 'true',
        },
      },
    };
  },
};

jest.mock('mapbox-gl', () => ({
  Map: () => ({}),
}));
jest.mock('../map/MapView', () => () => <div>Im the Map!</div>);
jest.mock('../projects/ProjectsMap', () => () => (
  <div>Im the ProjectsMap!</div>
));
jest.mock('@fbcnms/ui/insights/map/MapView', () => () => (
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

    // $FlowFixMe (T62907961) Relay flow types
    RelayEnvironment.mock.queueOperationResolver(operation =>
      MockPayloadGenerator.generate(operation, MOCK_RESOLVER),
    );

    const {getByText} = render(
      <Wrapper path={testCase.path}>
        <Main />
      </Wrapper>,
    );

    await waitForElement(() => getByText(testCase.text), {timeout: 100});
    expect(getByText(testCase.text)).toBeInTheDocument();
  });
});
