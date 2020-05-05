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
import * as React from 'react';
import SnackbarItem from '../components/SnackbarItem';
import {useCallback, useEffect, useState} from 'react';
import {useSnackbar as useNotistackSnackbar} from 'notistack';
import type {EnqueueSnackbarOptions} from 'notistack';
import type {Variants} from 'notistack';

type AllowedConfig = $Exact<{variant?: Variants} & EnqueueSnackbarOptions>;
export default function useSnackbar(
  message: string,
  config: AllowedConfig,
  show: boolean,
  dismissPrevious?: boolean,
) {
  const {enqueueSnackbar, closeSnackbar} = useNotistackSnackbar();
  const stringConfig = JSON.stringify(config);
  const [snackbarKey, setSnackbarKey] = useState(null);
  useEffect(() => {
    if (show) {
      const config: AllowedConfig = JSON.parse(stringConfig);
      const k = enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem
            id={key}
            message={message}
            variant={config.variant ?? 'success'}
          />
        ),
        ...config,
      });
      if (dismissPrevious) {
        snackbarKey != null && closeSnackbar(snackbarKey);
        setSnackbarKey(k);
      }
    }
    /*eslint-disable react-hooks/exhaustive-deps*/
  }, [
    // we shouldn't add snackbarKey
    // to the dependency list otherwise it'd creeate an infinite recursion
    closeSnackbar,
    dismissPrevious,
    enqueueSnackbar,
    message,
    show,
    stringConfig,
  ]);
  /*eslint-enable react-hooks/exhaustive-deps*/
}

export function useEnqueueSnackbar() {
  const {enqueueSnackbar} = useNotistackSnackbar();
  return useCallback(
    (message: string, config: Object) =>
      enqueueSnackbar(message, {
        children: key => (
          <SnackbarItem
            id={key}
            message={message}
            variant={config.variant ?? 'success'}
          />
        ),
        ...config,
      }),
    [enqueueSnackbar],
  );
}
