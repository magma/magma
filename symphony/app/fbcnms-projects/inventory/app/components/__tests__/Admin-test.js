/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import 'jest-dom/extend-expect';
import Admin from '../admin/Admin';
import React from 'react';
import {MemoryRouter} from 'react-router-dom';
import {Route} from 'react-router-dom';
import {cleanup, render} from '@testing-library/react';

jest.mock(
  '@fbcnms/magmalte/app/components/admin/AdminContextProvider',
  () => props => <div>{props.children}</div>,
);

const Wrapper = props => (
  <MemoryRouter initialEntries={[props.path]} initialIndex={1}>
    <Route path="/admin">{props.children}</Route>
  </MemoryRouter>
);

afterEach(cleanup);

global.CONFIG = {
  appData: {
    enabledFeatures: [],
    tabs: ['admin'],
    user: {
      isSuperUser: false,
    },
  },
  MAPBOX_ACCESS_TOKEN: '',
};

test('renders /settings', () => {
  const {getByTestId} = render(
    <Wrapper path={'/admin/settings'}>
      <Admin />
    </Wrapper>,
  );

  expect(getByTestId('change-password-title')).toBeInTheDocument();
});
