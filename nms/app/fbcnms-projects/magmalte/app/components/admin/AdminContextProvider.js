/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import * as React from 'react';
import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import {AppContextProvider} from '@fbcnms/ui/context/AppContext';

import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';

export default function AdminContextProvider(props: {children: React.Node}) {
  const {error, isLoading, response} = useMagmaAPI(MagmaV1API.getNetworks, {});

  if (isLoading) {
    return <LoadingFiller />;
  }

  const networkIds = error || !response ? ['mpk_test'] : response.sort();

  return (
    <AppContextProvider networkIDs={networkIds}>
      {props.children}
    </AppContextProvider>
  );
}
