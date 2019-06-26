/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import React, {useState} from 'react';
import Button from '@material-ui/core/Button';
import {storiesOf} from '@storybook/react';
import {useSnackbar} from '../hooks';

type Props = {
  message: string,
  variant: string,
};

const SnackbarTrigger = ({message, variant}: Props) => {
  const [isError, setError] = useState(false);
  useSnackbar(message, {variant: variant}, isError);
  return <Button onClick={() => setError(true)}>Display snackbar!</Button>;
};

storiesOf('Snackbar', module).add('error', () => (
  <SnackbarTrigger message="Wow, much error" variant="error" />
));

storiesOf('Snackbar', module).add('success', () => (
  <SnackbarTrigger message="Wow, much success" variant="success" />
));

storiesOf('Snackbar', module).add('long error', () => (
  <SnackbarTrigger
    message="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
    variant="error"
  />
));
