/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

// https://github.com/iamhosseindhv/notistack/pull/17
import {SnackbarContextNext} from 'notistack/build/SnackbarContext';
import {useContext, useEffect} from 'react';

export default function useSnackbar(
  message: string,
  config: any,
  show: boolean,
) {
  const enqueueSnackbar = useEnqueueSnackbar();
  const stringConfig = JSON.stringify(config);
  useEffect(() => {
    if (show) {
      enqueueSnackbar(message, config);
    }
  }, [message, show, stringConfig]);
}

export function useEnqueueSnackbar() {
  return useContext(SnackbarContextNext).handleEnqueueSnackbar;
}
