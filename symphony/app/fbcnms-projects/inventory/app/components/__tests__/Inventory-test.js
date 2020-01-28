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
import Inventory from '../Inventory';
import React from 'react';
import {MemoryRouter} from 'react-router-dom';
import {Route} from 'react-router-dom';
import {cleanup, render} from '@testing-library/react';

jest.mock('../../pages/Configure', () => () => <div>ConfigurePage</div>);
jest.mock('../../pages/Inventory', () => () => <div>InventoryPage</div>);
jest.mock('../map/LocationsMap', () => () => <div>LocationsMap</div>);

const Wrapper = props => (
  <MemoryRouter initialEntries={[props.path]} initialIndex={1}>
    <Route path="/inventory">{props.children}</Route>
  </MemoryRouter>
);

afterEach(cleanup);

global.CONFIG = {
  appData: {
    enabledFeatures: [],
    tabs: ['inventory'],
    user: {
      isSuperUser: false,
    },
  },
  MAPBOX_ACCESS_TOKEN: '',
};

test('renders /configure', () => {
  const {getByText} = render(
    <Wrapper path={'/inventory/configure'}>
      <Inventory />
    </Wrapper>,
  );

  expect(getByText('ConfigurePage')).toBeInTheDocument();
});

test('renders /settings', () => {
  const {getByTestId} = render(
    <Wrapper path={'/inventory/settings'}>
      <Inventory />
    </Wrapper>,
  );

  expect(getByTestId('change-password-title')).toBeInTheDocument();
});

test('renders /inventory', () => {
  const {getByText} = render(
    <Wrapper path={'/inventory/inventory'}>
      <Inventory />
    </Wrapper>,
  );

  expect(getByText('InventoryPage')).toBeInTheDocument();
});

test('renders /map', () => {
  const {getByText} = render(
    <Wrapper path={'/inventory/map'}>
      <Inventory />
    </Wrapper>,
  );

  expect(getByText('LocationsMap')).toBeInTheDocument();
});
