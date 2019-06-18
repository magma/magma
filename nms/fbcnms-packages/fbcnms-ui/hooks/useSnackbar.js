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
import {useEffect} from 'react';
import {useSnackbar as useNotistackSnackbar} from 'notistack';
import * as React from 'react';
import SnackbarItem from '../components/SnackbarItem.react';

export default function useSnackbar(
  message: string,
  config: any,
  show: boolean,
) {
  const {enqueueSnackbar} = useNotistackSnackbar();
  const stringConfig = JSON.stringify(config);
  useEffect(() => {
    if (show) {
      enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem
            id={key}
            message={message}
            variant={config.variant ?? 'success'}
          />
        ),
        ...config,
      });
    }
  }, [message, show, stringConfig]);
}

export function useEnqueueSnackbar() {
  const {enqueueSnackbar} = useNotistackSnackbar();
  return (message: string, config: Object) =>
    enqueueSnackbar(message, {
      children: key => (
        <SnackbarItem
          id={key}
          message={message}
          variant={config.variant ?? 'success'}
        />
      ),
      ...config,
    });
}
