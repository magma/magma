/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 * @flow
 * @format
 */

// https://github.com/iamhosseindhv/notistack/pull/17
import * as React from 'react';
// $FlowFixMe migrated to typescript
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
        content: key => (
          <SnackbarItem
            id={key}
            message={message}
            variant={config.variant ?? 'success'}
          />
        ),
        // $FlowFixMe[incompatible-use]
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
    (message: string, config: EnqueueSnackbarOptions) =>
      // $FlowFixMe[prop-missing]
      enqueueSnackbar(message, {
        content: key => (
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

export function useSnackbars() {
  const enqueueSnackbar = useEnqueueSnackbar();

  const successSnackbar = React.useCallback(
    (message: string) =>
      enqueueSnackbar(message, {
        variant: 'success',
      }),
    [enqueueSnackbar],
  );

  const errorSnackbar = React.useCallback(
    (message: string) => {
      enqueueSnackbar(message, {
        variant: 'error',
      });
    },
    [enqueueSnackbar],
  );

  const warningSnackbar = React.useCallback(
    (message: string) =>
      enqueueSnackbar(message, {
        variant: 'warning',
      }),
    [enqueueSnackbar],
  );

  const result = React.useMemo(
    () => ({
      success: successSnackbar,
      error: errorSnackbar,
      warning: warningSnackbar,
    }),
    [errorSnackbar, successSnackbar, warningSnackbar],
  );

  return result;
}
