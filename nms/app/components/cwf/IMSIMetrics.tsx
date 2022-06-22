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
 */

import LoadingFiller from '../LoadingFiller';
import Metrics from '../insights/Metrics';
import React from 'react';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';
import type {MetricGraphConfig} from '../insights/Metrics';

import MagmaAPI from '../../../api/MagmaAPI';
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPI';

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

export default function () {
  const params = useParams();
  const navigate = useNavigate();

  const {response, error, isLoading} = useMagmaAPI(
    MagmaAPI.metrics.networksNetworkIdPrometheusSeriesGet,
    {
      networkId: nullthrows(params.networkId),
    },
  );
  if (!response || error || isLoading) {
    return <LoadingFiller />;
  }

  const imsiSet = new Set<string>();
  response.forEach(item => {
    if (item.imsi) {
      imsiSet.add(item.imsi);
    }
  });
  const allIMSIs = [...imsiSet];
  const metrics = (
    <Metrics
      configs={IMSI_CONFIGS}
      onSelectorChange={(_, value) => navigate(value)}
      selectors={allIMSIs}
      defaultSelector={allIMSIs[0]}
      selectorName={'imsi'}
      renderOptionOverride={option => <ImsiAndIPMenuItem imsi={option} />}
    />
  );

  return (
    <Routes>
      <Route path=":selectedID" element={metrics} />
      <Route index element={metrics} />
    </Routes>
  );
}

function ImsiAndIPMenuItem(props: {imsi: string}) {
  const params = useParams();
  // The directory record endpoint requires that "IMSI" be prepended
  // to imsi number. Some metric series might have that on their label.
  const queryIMSI = props.imsi.startsWith('IMSI')
    ? props.imsi
    : 'IMSI' + props.imsi;
  const {response} = useMagmaAPI(
    MagmaAPI.carrierWifiNetworks
      .cwfNetworkIdSubscribersSubscriberIdDirectoryRecordGet,
    {networkId: params.networkId!, subscriberId: queryIMSI},
  );

  const ipv4 = response?.ipv4_addr;
  return (
    <>
      ipv4 ? `${props.imsi} : ${ipv4}` : props.imsi
    </>
  );
}
