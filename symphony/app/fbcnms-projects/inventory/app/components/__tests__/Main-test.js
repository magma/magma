/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import Main from '../Main';
import React from 'react';
import renderer from 'react-test-renderer';
import {BrowserRouter} from 'react-router-dom';

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

it('renders without crashing', () => {
  global.CONFIG = {
    appData: {},
    MAPBOX_ACCESS_TOKEN: '',
  };
  renderer.create(
    <BrowserRouter>
      <Main user={{}} />
    </BrowserRouter>,
  );
});
