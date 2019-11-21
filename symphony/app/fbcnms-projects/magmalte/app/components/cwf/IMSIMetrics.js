/**
 * Copyright 2004-present Facebook. All Rights Reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * @flow
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';

import LoadingFiller from '@fbcnms/ui/components/LoadingFiller';
import MagmaV1API from '@fbcnms/magma-api/client/WebClient';
import MenuItem from '@material-ui/core/MenuItem';
import Metrics from '../insights/Metrics';
import React from 'react';
import {Route} from 'react-router-dom';

import nullthrows from '@fbcnms/util/nullthrows';
import useMagmaAPI from '../../common/useMagmaAPI';
import {useRouter} from '@fbcnms/ui/hooks';

const IMSI_CONFIGS: Array<MetricGraphConfig> = [
  {
    label: 'Traffic In',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: imsi => `sum(octets_in{imsi="${imsi}"})`,
      },
    ],
  },
  {
    label: 'Traffic Out',
    basicQueryConfigs: [],
    customQueryConfigs: [
      {
        resolveQuery: imsi => `sum(octets_out{imsi="${imsi}"})`,
      },
    ],
  },
  {
    label: 'Throughput In',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `avg(rate(octets_in{imsi="${imsi}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Throughput Out',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `avg(rate(octets_out{imsi="${imsi}"}[5m]))`,
      },
    ],
  },
  {
    label: 'Active Sessions',
    basicQueryConfigs: [],
    unit: '',
    customQueryConfigs: [
      {
        resolveQuery: imsi => `active_sessions{imsi="${imsi}"}`,
      },
    ],
  },
];

export default function() {
  const {history, relativePath, relativeUrl, match} = useRouter();

  const {response, error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusSeries,
    {networkId: nullthrows(match.params.networkId)},
  );
  if (!response || error || isLoading) {
    return <LoadingFiller />;
  }

  const imsiSet = new Set();
  response.forEach(item => {
    if (item.imsi) {
      imsiSet.add(item.imsi);
    }
  });
  const allIMSIs = [...imsiSet];

  const imsiMenuItems = allIMSIs.map(imsi => (
    <MenuItem value={imsi} key={imsi}>
      <ImsiAndIPMenuItem imsi={imsi} />
    </MenuItem>
  ));

  return (
    <Route
      path={relativePath('/:selectedID?')}
      render={() => (
        <Metrics
          configs={IMSI_CONFIGS}
          onSelectorChange={({target}) => {
            history.push(relativeUrl(`/${target.value}`));
          }}
          menuItemOverrides={imsiMenuItems}
          selectors={allIMSIs}
          defaultSelector={allIMSIs[0]}
          selectorName={'imsi'}
        />
      )}
    />
  );
}

function ImsiAndIPMenuItem(props: {imsi: string}) {
  const {match} = useRouter();
  const {response} = useMagmaAPI(
    MagmaV1API.getCwfByNetworkIdSubscribersBySubscriberIdDirectoryRecord,
    {networkId: match.params.networkId, subscriberId: props.imsi},
  );

  const ipv4 = response?.ipv4_addr;
  return ipv4 ? `${props.imsi} : ${ipv4}` : props.imsi;
}
