/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import Button from '@material-ui/core/Button';
import React, {useState} from 'react';
import {STORY_CATEGORIES} from '../storybookUtils';
import {storiesOf} from '@storybook/react';
import {useSnackbar} from '../../hooks';
import type {Variants} from 'notistack';

type Props = {
  message: string,
  variant: Variants,
};

const SnackbarTrigger = ({message, variant}: Props) => {
  const [isError, setError] = useState(false);
  useSnackbar(message, {variant: variant}, isError);
  return <Button onClick={() => setError(true)}>Display snackbar!</Button>;
};

storiesOf(`${STORY_CATEGORIES.MUI_COMPONENTS}/Snackbar`, module)
  .add('error', () => (
    <SnackbarTrigger message="Wow, much error" variant="error" />
  ))
  .add('success', () => (
    <SnackbarTrigger message="Wow, much success" variant="success" />
  ))
  .add('default', () => (
    <SnackbarTrigger message="Wow, much default" variant="default" />
  ))
  .add('info', () => (
    <SnackbarTrigger message="Wow, much info" variant="info" />
  ))
  .add('warning', () => (
    <SnackbarTrigger message="Wow, much warning" variant="warning" />
  ))
  .add('long error', () => (
    <SnackbarTrigger
      message="Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua."
      variant="error"
    />
  ));
