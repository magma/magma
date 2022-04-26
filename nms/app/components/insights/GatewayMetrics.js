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

import type {MetricGraphConfig} from './Metrics';

import LoadingFiller from '../../../fbc_js_core/ui/components/LoadingFiller';
import MagmaV1API from '../../../generated/WebClient';
import Metrics from './Metrics';
import React from 'react';
import {Route} from 'react-router-dom';

import useMagmaAPI from '../../../api/useMagmaAPI';
import {useRouter} from '../../../fbc_js_core/ui/hooks';
import {useSnackbar} from '../../../fbc_js_core/ui/hooks';

export default function (props: {configs: MetricGraphConfig[]}) {
  const {history, relativePath, relativeUrl, match} = useRouter();

  const {error, isLoading, response: selectors} = useMagmaAPI(
    MagmaV1API.getNetworksByNetworkIdGateways,
    {
      networkId: match.params.networkId,
    },
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
