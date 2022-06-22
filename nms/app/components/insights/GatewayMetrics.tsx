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

import type {MetricGraphConfig} from './Metrics';

import LoadingFiller from '../LoadingFiller';
import Metrics from './Metrics';
import React from 'react';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';

import MagmaAPI from '../../../api/MagmaAPI';
import useMagmaAPI from '../../../api/useMagmaAPI';
import {useSnackbar} from '../../hooks';

export default function (props: {configs: Array<MetricGraphConfig>}) {
  const navigate = useNavigate();
  const params = useParams();

  const {error, isLoading, response: selectors} = useMagmaAPI(
    MagmaAPI.gateways.networksNetworkIdGatewaysGet,
    {
      networkId: params.networkId!,
    },
  );

  useSnackbar('Error fetching devices', {variant: 'error'}, !!error);

  if (error || isLoading || !selectors) {
    return <LoadingFiller />;
  }

  const gatewayNames = Object.keys(selectors);
  const defaultGateway = gatewayNames[0];

  const metrics = (
    <Metrics
      configs={props.configs}
      onSelectorChange={(e, value) => {
        navigate(value);
      }}
      selectors={gatewayNames}
      defaultSelector={defaultGateway}
      selectorName={'gatewayID'}
    />
  );

  return (
    <Routes>
      <Route path=":selectedID" element={metrics} />
      <Route index element={metrics} />
    </Routes>
  );
}
