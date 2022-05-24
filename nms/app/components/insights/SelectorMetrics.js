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
 * @flow strict-local
 * @format
 */

import type {MetricGraphConfig} from '../insights/Metrics';

// $FlowFixMe migrated to typescript
import LoadingFiller from '../LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import Metrics from '../insights/Metrics';
import React from 'react';
import {Route, Routes, useNavigate, useParams} from 'react-router-dom';

// $FlowFixMe migrated to typescript
import nullthrows from '../../../shared/util/nullthrows';
import useMagmaAPI from '../../../api/useMagmaAPIFlow';
import {useCallback, useState} from 'react';
import {useSnackbar} from '../../../app/hooks';

export default function (props: {
  configs: MetricGraphConfig[],
  selectorKey: string,
}) {
  const navigate = useNavigate();
  const params = useParams();

  const [allMetrics, setAllMetrics] = useState();
  const [selectedItem, setSelectedItem] = useState('');

  const {error, isLoading} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdPrometheusSeries,
    {networkId: nullthrows(params.networkId)},
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

  const metrics = (
    <Metrics
      configs={props.configs}
      onSelectorChange={(e, value) => navigate(value)}
      selectors={Object.keys(allMetrics)}
      defaultSelector={selectedItem}
      selectorName={props.selectorKey}
    />
  );

  return (
    <Routes>
      <Route path=":selectedID" element={metrics} />
      <Route index element={metrics} />
    </Routes>
  );
}
