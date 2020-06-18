/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow strict-local
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import Metrics from '../insights/Metrics';
import React from 'react';
import {Route} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '@fbcnms/ui/magma/useMagmaAPI';
import {useCallback, useState} from 'react';
import {useRouter} from '@fbcnms/ui/hooks';
import {useSnackbar} from '@fbcnms/ui/hooks';

export default function(props: {
  configs: MetricGraphConfig[],
  selectorKey: string,
}) {
  const {history, relativePath, relativeUrl, match} = useRouter();

  const [allMetrics, setAllMetrics] = useState();
  const [selectedItem, setSelectedItem] = useState('');

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusSeries,
    {networkId: nullthrows(match.params.networkId)},
    useCallback(
      response => {
        const metricsByDeviceID = {};
        response.forEach(item => {
          if (item[props.selectorKey]) {
            metricsByDeviceID[item[props.selectorKey]] =
              metricsByDeviceID[item[props.selectorKey]] || new Set();
            metricsByDeviceID[item[props.selectorKey]].add(item.__name__);
          }
        });
        setSelectedItem(Object.keys(metricsByDeviceID)[0]);
        setAllMetrics(metricsByDeviceID);
      },
      [props.selectorKey],
    ),
  );

  useSnackbar('Error fetching subscribers', {variant: 'error'}, error);

  if (error || isLoading || !allMetrics) {
    return <LoadingFiller />;
  }

  return (
    <Route
      path={relativePath('/:selectedID?')}
      render={() => (
        <Metrics
          configs={props.configs}
          onSelectorChange={(e, value) =>
            history.push(relativeUrl(`/${value}`))
          }
          selectors={Object.keys(allMetrics)}
          defaultSelector={selectedItem}
          selectorName={props.selectorKey}
        />
      )}
    />
  );
}
