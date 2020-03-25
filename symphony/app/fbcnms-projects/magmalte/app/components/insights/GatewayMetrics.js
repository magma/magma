/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {MetricGraphConfig} from './Metrics';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Metrics from './Metrics';
import React from 'react';
import {Route} from 'react-router-dom';

import useMagmaAPI from '../../common/useMagmaAPI';
import {useRouter} from '@fbcnms/ui/hooks';
import {useSnackbar} from '@fbcnms/ui/hooks';

export default function(props: {configs: MetricGraphConfig[]}) {
  const {history, relativePath, relativeUrl, match} = useRouter();

  const {error, isLoading, response: selectors} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdGateways,
    {networkId: match.params.networkId},
  );

  useSnackbar('Error fetching devices', {variant: 'error'}, error);

  if (error || isLoading || !selectors) {
    return <LoadingFiller />;
  }

  const gatewayNames = Object.keys(selectors);
  const defaultGateway = gatewayNames[0];

  return (
    <Route
      path={relativePath('/:selectedID?')}
      render={() => (
        <Metrics
          configs={props.configs}
          onSelectorChange={(e, value) => {
            history.push(relativeUrl(`/${value}`));
          }}
          selectors={gatewayNames}
          defaultSelector={defaultGateway}
          selectorName={'gatewayID'}
        />
      )}
    />
  );
}
