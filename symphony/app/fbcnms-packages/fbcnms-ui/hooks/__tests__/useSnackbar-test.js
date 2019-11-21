/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React from 'react';
import renderer from 'react-test-renderer';
import {SnackbarProvider} from 'notistack';

import {useSnackbar} from '../../hooks';

jest.mock('@material-ui/core/Slide', () => () => <div />);

it('renders without crashing', () => {
  const tree = renderer.create(
    <SnackbarProvider
      maxSnack={3}
      autoHideDuration={10000}
      anchorOrigin={{
        vertical: 'bottom',
        horizontal: 'right',
      }}>
      <Test />
    </SnackbarProvider>,
  );

  // needed to trigger `useEffect()` hook
  tree.update();
});

const Test = () => {
  useSnackbar('Error', {variant: 'error'}, true);
  return <div />;
};
