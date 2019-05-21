/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

jest.mock('mapbox-gl', () => ({
  Map: () => ({}),
}));
jest.mock('../MapView', () => () => <div>Im the Map!</div>);

import React from 'react';
import {BrowserRouter} from 'react-router-dom';
import Main from '../Main';
import renderer from 'react-test-renderer';

it('renders without crashing', () => {
  global.CONFIG = {
    appData: {},
    MAPBOX_ACCESS_TOKEN: 'mapbox-token',
  };
  renderer.create(
    <BrowserRouter>
      <Main user={{}} />
    </BrowserRouter>,
  );
});
