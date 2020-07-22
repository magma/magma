/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {$AxiosError} from 'axios';

import React from 'react';

import {useSnackbar} from '@fbcnms/ui/hooks';

export default function NetworkError({error}: {error: $AxiosError<string>}) {
  let errorMessage = error.message;
  if (error.response && error.response.status >= 400) {
    errorMessage = error.response?.statusText;
  }
  useSnackbar(
    'Unable to communicate with magma controller: ' + errorMessage,
    {variant: 'error'},
    !!error,
  );
  return <div />;
}
